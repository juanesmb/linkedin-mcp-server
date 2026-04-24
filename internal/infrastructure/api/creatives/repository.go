package creatives

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"linkedin-mcp/internal/infrastructure/api/gateway"
	"linkedin-mcp/internal/infrastructure/middleware"
)

const (
	logMessageFailedRequest        = "failed to make request"
	logMessageLinkedInAPIError     = "linkedin api responded with error"
	logMessageFailedDecodeResponse = "failed to decode response"

	logTagURL    = "url"
	logTagError  = "error"
	logTagStatus = "status"
	logTagBody   = "body"

	errFmtFailedRequest         = "failed to make request: %w"
	errFmtLinkedInAPIErrorJSON  = "linkedin api error: status %d, body: %v"
	errFmtLinkedInAPIErrorPlain = "linkedin api error: status %d, body: %s"
	errFmtLinkedInAPIError      = "linkedin api error: status %d"
	errFmtDecodeResponse        = "failed to decode response: %w"
)

type Logger interface {
	Error(ctx context.Context, message string, tags map[string]string)
}

type Repository struct {
	gatewayClient *gateway.Client
	queryBuilder  *QueryBuilder
	logger        Logger
}

func NewRepository(gatewayClient *gateway.Client, queryBuilder *QueryBuilder, logger Logger) *Repository {
	return &Repository{
		gatewayClient: gatewayClient,
		queryBuilder:  queryBuilder,
		logger:        logger,
	}
}

func (r *Repository) SearchCreatives(ctx context.Context, input SearchInput) (*SearchResult, error) {
	requestURL := r.queryBuilder.BuildSearchCreativesByCampaignsQuery(input)
	resourcePath, query, err := gateway.ParseLinkedInRESTProxyTarget(requestURL)
	if err != nil {
		return nil, fmt.Errorf("failed to build gateway proxy target: %w", err)
	}
	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("missing authenticated user in request context")
	}

	connectionResponse, err := r.gatewayClient.GetLinkedInConnection(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch LinkedIn connection state from gateway: %w", err)
	}
	if gateway.IsLinkedInNotConnectedResponse(connectionResponse) {
		return nil, gateway.ErrLinkedInNotConnected
	}
	if connectionResponse.StatusCode < 200 || connectionResponse.StatusCode >= 300 {
		return nil, fmt.Errorf("failed to fetch LinkedIn connection state from gateway: status %d", connectionResponse.StatusCode)
	}

	// LinkedIn Rest.li creatives criteria finder expects X-RestLi-Method: FINDER (see Microsoft Learn).
	response, err := r.gatewayClient.ProxyLinkedInOrRefresh(ctx, userID, resourcePath, query, map[string]string{
		"X-RestLi-Method": "FINDER",
	})
	if err != nil {
		r.logError(ctx, logMessageFailedRequest, map[string]string{
			logTagURL:   requestURL,
			logTagError: err.Error(),
		})
		return nil, fmt.Errorf(errFmtFailedRequest, err)
	}
	if gateway.IsLinkedInNotConnectedResponse(response) {
		return nil, gateway.ErrLinkedInNotConnected
	}
	if validationErr, ok := gateway.ParseLinkedInParamValidationResponse(response); ok {
		return nil, validationErr
	}

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		bodyString := strings.TrimSpace(string(response.Body))
		tags := map[string]string{
			logTagURL:    requestURL,
			logTagStatus: strconv.Itoa(response.StatusCode),
		}
		if bodyString != "" {
			tags[logTagBody] = bodyString
		}

		r.logError(ctx, logMessageLinkedInAPIError, tags)

		var errBody any
		if err := json.Unmarshal(response.Body, &errBody); err == nil {
			return nil, fmt.Errorf(errFmtLinkedInAPIErrorJSON, response.StatusCode, errBody)
		}

		if bodyString != "" {
			return nil, fmt.Errorf(errFmtLinkedInAPIErrorPlain, response.StatusCode, bodyString)
		}

		return nil, fmt.Errorf(errFmtLinkedInAPIError, response.StatusCode)
	}

	var liResp LinkedInListResponse
	if err := json.Unmarshal(response.Body, &liResp); err != nil {
		r.logError(ctx, logMessageFailedDecodeResponse, map[string]string{
			logTagURL:   requestURL,
			logTagError: err.Error(),
		})
		return nil, fmt.Errorf(errFmtDecodeResponse, err)
	}

	normalized := make([]NormalizedCreative, 0, len(liResp.Elements))
	for _, el := range liResp.Elements {
		normalized = append(normalized, NormalizeCreative(el))
	}

	return &SearchResult{
		Elements: normalized,
		Paging:   summarizeCreativesPagination(liResp.Metadata, input.PageSize),
	}, nil
}

func summarizeCreativesPagination(metadata map[string]any, requestedPageSize int) PagingSummary {
	out := PagingSummary{}
	if requestedPageSize > 0 {
		out.PageSize = requestedPageSize
	}
	if metadata != nil {
		if v, ok := metadata["nextPageToken"].(string); ok {
			out.NextPageToken = v
		}
	}
	return out
}

func (r *Repository) logError(ctx context.Context, message string, tags map[string]string) {
	if r.logger == nil {
		return
	}
	r.logger.Error(ctx, message, tags)
}
