package adaccounts

import (
	"context"
	"encoding/json"
	"fmt"
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

func (r *Repository) SearchAdAccounts(ctx context.Context, input SearchInput) (*SearchResult, error) {
	requestURL, headers := r.queryBuilder.BuildSearchAdAccountsQuery(input)

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
	if err = json.Unmarshal(response.Body, &liResp); err != nil {
		r.logError(ctx, logMessageFailedDecodeResponse, map[string]string{
			logTagURL:   requestURL,
			logTagError: err.Error(),
		})
		return nil, fmt.Errorf(errFmtDecodeResponse, err)
	}

	result := &SearchResult{
		Elements: liResp.Elements,
		Paging:   liResp.Paging,
	}

	return result, nil
}

func (r *Repository) logError(ctx context.Context, message string, tags map[string]string) {
	if r.logger == nil {
		return
	}

	r.logger.Error(ctx, message, tags)
}
