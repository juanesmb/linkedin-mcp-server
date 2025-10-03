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
	q := url.Values{}
	q.Set("q", "search")

	// Arrays -> repeated query params with the same key
	for _, v := range input.CampaignGroupURNs {
		if v != "" {
			q.Add("search.campaignGroup.values", v)
		}
	}
	for _, v := range input.AssociatedEntityValues {
		if v != "" {
			q.Add("search.associatedEntity.values", v)
		}
	}
	for _, v := range input.CampaignURNs {
		if v != "" {
			q.Add("search.id.values", v)
		}
	}
	for _, v := range input.Status {
		if v != "" {
			q.Add("search.status.values", v)
		}
	}
	for _, v := range input.Type {
		if v != "" {
			q.Add("search.type.values", v)
		}
	}
	for _, v := range input.Name {
		if v != "" {
			q.Add("search.name.values", v)
		}
	}

	// Note: LinkedIn adCampaigns search does not support a test flag; ignore input.Test

	if input.SortOrder != "" {
		q.Set("sortOrder", input.SortOrder)
	}
	if input.PageSize > 0 {
		q.Set("pageSize", fmt.Sprintf("%d", input.PageSize))
	}
	if input.PageToken != "" {
		q.Set("pageToken", input.PageToken)
	}

	return q.Encode()
}
