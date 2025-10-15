package campaigns

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"linkedin-mcp/internal/infrastructure/api"
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

	pagingNextKey       = "next"
	queryParamPageToken = "pageToken"
)

type Logger interface {
	Error(ctx context.Context, message string, tags map[string]string)
}

type Repository struct {
	client       api.Client
	queryBuilder *QueryBuilder
	logger       Logger
}

func NewRepository(client api.Client, queryBuilder *QueryBuilder, logger Logger) *Repository {
	return &Repository{
		client:       client,
		queryBuilder: queryBuilder,
		logger:       logger,
	}
}

func (r *Repository) SearchCampaigns(ctx context.Context, input SearchInput) (*SearchResult, error) {
	requestURL, headers := r.queryBuilder.BuildSearchCampaignsQuery(input)

	response, err := r.client.Get(ctx, requestURL, headers)
	if err != nil {
		r.logError(ctx, logMessageFailedRequest, map[string]string{
			logTagURL:   requestURL,
			logTagError: err.Error(),
		})
		return nil, fmt.Errorf(errFmtFailedRequest, err)
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

		trimmedBody := strings.TrimSpace(string(response.Body))
		if trimmedBody != "" {
			return nil, fmt.Errorf(errFmtLinkedInAPIErrorPlain, response.StatusCode, trimmedBody)
		}

		return nil, fmt.Errorf(errFmtLinkedInAPIError, response.StatusCode)
	}

	var liResp LinkedInResponse
	err = json.Unmarshal(response.Body, &liResp)
	if err != nil {
		r.logError(ctx, logMessageFailedDecodeResponse, map[string]string{
			logTagURL:   requestURL,
			logTagError: err.Error(),
		})
		return nil, fmt.Errorf(errFmtDecodeResponse, err)
	}

	result := &SearchResult{
		Elements: liResp.Elements,
	}

	if nextRaw, ok := liResp.Paging[pagingNextKey].(string); ok && nextRaw != "" {
		if u, err := url.Parse(nextRaw); err == nil {
			if token := u.Query().Get(queryParamPageToken); token != "" {
				result.Metadata.NextPageToken = token
			}
		}
	}

	return result, nil
}

func (r *Repository) logError(ctx context.Context, message string, tags map[string]string) {
	if r.logger == nil {
		return
	}

	r.logger.Error(ctx, message, tags)
}
