package searchcreatives

import (
	"context"
	"fmt"
	"strings"

	"linkedin-mcp/internal/infrastructure/api/creatives"
	"linkedin-mcp/internal/infrastructure/tools/searchcreatives/dto"
	"linkedin-mcp/internal/infrastructure/tools/toolerrors"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const defaultPageSize = 100

const maxCreativesPageSize = 100

type Tool struct {
	repository *creatives.Repository
	connectURL string
}

func NewTool(repository *creatives.Repository, connectURL string) *Tool {
	return &Tool{repository: repository, connectURL: connectURL}
}

func (t *Tool) SearchCreatives(ctx context.Context, req *mcp.CallToolRequest, input dto.Input) (*mcp.CallToolResult, dto.Output, error) {
	result := &mcp.CallToolResult{}

	if err := validateInput(&input); err != nil {
		return result, dto.Output{}, fmt.Errorf("input validation failed: %w", err)
	}

	campaignURN, err := resolveCampaignURN(input.CampaignURN, input.CampaignID)
	if err != nil {
		return result, dto.Output{}, fmt.Errorf("input validation failed: %w", err)
	}

	pageSize := defaultPageSize
	if input.PageSize != nil {
		pageSize = *input.PageSize
	}

	var pageToken string
	if input.PageToken != nil {
		pageToken = *input.PageToken
	}

	searchInput := creatives.SearchInput{
		AccountID:    strings.TrimSpace(input.AccountID),
		CampaignURNs: []string{campaignURN},
		PageSize:     pageSize,
		PageToken:    strings.TrimSpace(pageToken),
		SortOrder:    strings.TrimSpace(input.SortOrder),
	}

	searchResult, err := t.repository.SearchCreatives(ctx, searchInput)
	if err != nil {
		return result, dto.Output{}, toolerrors.WrapToolExecutionError("search creatives", err, t.connectURL)
	}

	return result, dto.Output{
		Elements: searchResult.Elements,
		Paging:   searchResult.Paging,
	}, nil
}

func validateInput(input *dto.Input) error {
	input.AccountID = strings.TrimSpace(input.AccountID)
	if input.AccountID == "" {
		return fmt.Errorf("accountID is required")
	}
	if err := validateNumericID(input.AccountID); err != nil {
		return fmt.Errorf("accountID: %w", err)
	}

	if input.PageSize != nil {
		if *input.PageSize < 1 {
			return fmt.Errorf("pageSize must be at least 1")
		}
		if *input.PageSize > maxCreativesPageSize {
			return fmt.Errorf("pageSize cannot exceed %d for creatives search", maxCreativesPageSize)
		}
	}

	if o := strings.TrimSpace(input.SortOrder); o != "" && o != "ASCENDING" && o != "DESCENDING" {
		return fmt.Errorf("sortOrder must be ASCENDING or DESCENDING")
	}

	hasURN := strings.TrimSpace(input.CampaignURN) != ""
	hasID := strings.TrimSpace(input.CampaignID) != ""
	if !hasURN && !hasID {
		return fmt.Errorf("either campaignURN or campaignID is required")
	}
	if hasID && !hasURN {
		if err := validateNumericID(strings.TrimSpace(input.CampaignID)); err != nil {
			return fmt.Errorf("campaignID: %w", err)
		}
	}

	return nil
}

func resolveCampaignURN(urn, numericID string) (string, error) {
	urn = strings.TrimSpace(urn)
	if urn != "" {
		return urn, nil
	}
	id := strings.TrimSpace(numericID)
	if id == "" {
		return "", fmt.Errorf("campaign URN or id is required")
	}
	return fmt.Sprintf("urn:li:sponsoredCampaign:%s", id), nil
}

func validateNumericID(id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return fmt.Errorf("cannot be empty")
	}
	for _, r := range id {
		if r < '0' || r > '9' {
			return fmt.Errorf("must contain only digits")
		}
	}
	return nil
}
