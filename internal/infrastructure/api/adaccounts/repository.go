package adaccounts

import (
	"context"
	"encoding/json"
	"fmt"

	"linkedin-mcp/internal/infrastructure/api"
)

type Repository struct {
	client       api.Client
	queryBuilder *QueryBuilder
}

func NewRepository(client api.Client, queryBuilder *QueryBuilder) *Repository {
	return &Repository{
		client:       client,
		queryBuilder: queryBuilder,
	}
}

func (r *Repository) SearchAdAccounts(ctx context.Context, input SearchInput) (*SearchResult, error) {
	requestURL, headers := r.queryBuilder.BuildSearchAdAccountsQuery(input)

	response, err := r.client.Get(ctx, requestURL, headers)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		var errBody map[string]any
		_ = json.Unmarshal(response.Body, &errBody)

		return nil, fmt.Errorf("linkedin api error: status %d, body: %v", response.StatusCode, errBody)
	}

	var liResp LinkedInResponse
	if err = json.Unmarshal(response.Body, &liResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	result := &SearchResult{
		Elements: liResp.Elements,
		Paging:   liResp.Paging,
	}

	return result, nil
}
