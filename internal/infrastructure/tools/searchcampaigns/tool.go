package searchcampaigns

import (
	"context"
	"fmt"
	"strings"

	"linkedin-mcp/internal/infrastructure/api/campaigns"
	"linkedin-mcp/internal/infrastructure/tools/searchcampaigns/dto"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type Tool struct {
	repository *campaigns.Repository
}

func NewTool(repository *campaigns.Repository) *Tool {
	return &Tool{
		repository: repository,
	}
}

func (t *Tool) SearchCampaigns(ctx context.Context, req *mcp.CallToolRequest, input dto.Input) (*mcp.CallToolResult, dto.Output, error) {
	result := &mcp.CallToolResult{}

	if err := t.validateInput(input); err != nil {
		return result, dto.Output{}, fmt.Errorf("input validation failed: %w", err)
	}

	searchInput := t.convertInput(input)

	searchResult, err := t.repository.SearchCampaigns(ctx, searchInput)
	if err != nil {
		return result, dto.Output{}, fmt.Errorf("failed to search campaigns: %w", err)
	}

	output := dto.Output{
		Elements: searchResult.Elements,
		Metadata: dto.Metadata{
			NextPageToken: searchResult.Metadata.NextPageToken,
		},
	}

	return result, output, nil
}

func (t *Tool) validateInput(input dto.Input) error {
	// Validate AccountID
	if input.AccountID == "" {
		return fmt.Errorf("accountID is required")
	}
	input.AccountID = strings.TrimSpace(input.AccountID)
	if input.AccountID == "" {
		return fmt.Errorf("accountID cannot be empty or whitespace only")
	}

	// Validate page size
	if input.PageSize < 0 {
		return fmt.Errorf("pageSize must be non-negative")
	}
	if input.PageSize > 1000 {
		return fmt.Errorf("pageSize cannot exceed 1000")
	}

	// Validate sort order
	if input.SortOrder != "" && input.SortOrder != "ASCENDING" && input.SortOrder != "DESCENDING" {
		return fmt.Errorf("sortOrder must be either ASCENDING or DESCENDING")
	}

	// Validate status values
	validStatuses := map[string]bool{
		"ACTIVE":           true,
		"PAUSED":           true,
		"ARCHIVED":         true,
		"COMPLETED":        true,
		"CANCELED":         true,
		"DRAFT":            true,
		"PENDING_DELETION": true,
		"REMOVED":          true,
	}
	for _, status := range input.Status {
		if !validStatuses[status] {
			return fmt.Errorf("invalid status: %s", status)
		}
	}

	// Validate type values
	validTypes := map[string]bool{
		"TEXT_AD":           true,
		"SPONSORED_UPDATES": true,
		"SPONSORED_INMAILS": true,
		"DYNAMIC":           true,
	}
	for _, campaignType := range input.Type {
		if !validTypes[campaignType] {
			return fmt.Errorf("invalid type: %s", campaignType)
		}
	}

	// Sanitize string inputs
	for i, urn := range input.CampaignGroupURNs {
		input.CampaignGroupURNs[i] = strings.TrimSpace(urn)
	}
	for i, urn := range input.CampaignURNs {
		input.CampaignURNs[i] = strings.TrimSpace(urn)
	}
	for i, name := range input.Name {
		input.Name[i] = strings.TrimSpace(name)
	}

	return nil
}

func (t *Tool) convertInput(input dto.Input) campaigns.SearchInput {
	var pageToken string
	if input.PageToken != nil {
		pageToken = *input.PageToken
	}

	return campaigns.SearchInput{
		AccountID:              input.AccountID,
		CampaignGroupURNs:      input.CampaignGroupURNs,
		AssociatedEntityValues: input.AssociatedEntityValues,
		CampaignURNs:           input.CampaignURNs,
		Status:                 input.Status,
		Type:                   input.Type,
		Name:                   input.Name,
		Test:                   input.Test,
		SortOrder:              input.SortOrder,
		PageSize:               input.PageSize,
		PageToken:              pageToken,
	}
}
