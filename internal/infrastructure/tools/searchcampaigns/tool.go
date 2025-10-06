package searchcampaigns

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"linkedin-mcp/internal/infrastructure/tools/searchcampaigns/dto"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type Tool struct {
	accessToken string
	accountID   string
	baseURL     string
	version     string
}

type linkedInResponse struct {
	Elements []map[string]any `json:"elements"`
	Paging   map[string]any   `json:"paging"`
}

func NewTool(accessToken, accountID, baseURL, version string) *Tool {
	return &Tool{
		accessToken: accessToken,
		accountID:   accountID,
		baseURL:     baseURL,
		version:     version,
	}
}

func (t *Tool) SearchCampaigns(ctx context.Context, req *mcp.CallToolRequest, input dto.Input) (*mcp.CallToolResult, dto.Output, error) {
	result := &mcp.CallToolResult{}

	endpoint := fmt.Sprintf("%s/adAccounts/%s/adCampaigns", strings.TrimRight(t.baseURL, "/"), url.PathEscape(t.accountID))

	queryParams := makeQueryParams(input)

	fullURL := endpoint + "?" + queryParams

	reqHTTP, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)
	if err != nil {
		return result, dto.Output{}, err
	}
	reqHTTP.Header.Set("Authorization", t.accessToken)
	reqHTTP.Header.Set("LinkedIn-Version", t.version)
	reqHTTP.Header.Set("X-Restli-Protocol-Version", "2.0.0")
	reqHTTP.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Do(reqHTTP)
	if err != nil {
		return result, dto.Output{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var errBody map[string]any
		_ = json.NewDecoder(resp.Body).Decode(&errBody)

		return result, dto.Output{}, fmt.Errorf("linkedin api error: %s (%d) %v", resp.Status, resp.StatusCode, errBody)
	}

	var liResp linkedInResponse
	err = json.NewDecoder(resp.Body).Decode(&liResp)
	if err != nil {
		return result, dto.Output{}, err
	}

	output := dto.Output{
		Elements: liResp.Elements,
	}

	// Best-effort extraction of next page token if the server returned a next link with pageToken.
	if nextRaw, ok := liResp.Paging["next"].(string); ok && nextRaw != "" {
		if u, err := url.Parse(nextRaw); err == nil {
			if token := u.Query().Get("pageToken"); token != "" {
				output.Metadata.NextPageToken = token
			}
		}
	}

	return result, output, nil
}

func makeQueryParams(input dto.Input) string {
	// Build query parameters manually to avoid double-encoding the search parameter
	var params []string

	// Always include q=search
	params = append(params, "q=search")

	// Build single Rest.li-style composite search parameter, e.g.:
	// search=(type:(values:List(SPONSORED_UPDATES)),status:(values:List(ACTIVE)))
	var searchParts []string

	addList := func(field string, values []string) {
		cleaned := make([]string, 0, len(values))
		for _, v := range values {
			if v == "" {
				continue
			}
			cleaned = append(cleaned, v)
		}
		if len(cleaned) == 0 {
			return
		}
		searchParts = append(searchParts, fmt.Sprintf("%s:(values:List(%s))", field, strings.Join(cleaned, ",")))
	}

	addList("campaignGroup", input.CampaignGroupURNs)
	addList("associatedEntity", input.AssociatedEntityValues)
	addList("id", input.CampaignURNs)
	addList("status", input.Status)
	addList("type", input.Type)
	addList("name", input.Name)

	if input.Test != nil {
		searchParts = append(searchParts, fmt.Sprintf("test:%t", *input.Test))
	}

	if len(searchParts) > 0 {
		// Build the search parameter without URL encoding the Rest.li syntax
		searchParam := fmt.Sprintf("search=(%s)", strings.Join(searchParts, ","))
		params = append(params, searchParam)
	}

	if input.SortOrder != "" {
		params = append(params, fmt.Sprintf("sortOrder=%s", url.QueryEscape(input.SortOrder)))
	}
	if input.PageSize > 0 {
		params = append(params, fmt.Sprintf("pageSize=%d", input.PageSize))
	}
	if input.PageToken != "" {
		params = append(params, fmt.Sprintf("pageToken=%s", url.QueryEscape(input.PageToken)))
	}

	return strings.Join(params, "&")
}
