package searchadaccounts

import (
	"context"
	"fmt"
	"strings"

	"linkedin-mcp/internal/infrastructure/api/adaccounts"
	"linkedin-mcp/internal/infrastructure/tools/searchadaccounts/dto"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type Tool struct {
	repository *adaccounts.Repository
}

func NewTool(repository *adaccounts.Repository) *Tool {
	return &Tool{repository: repository}
}

func (t *Tool) SearchAdAccounts(ctx context.Context, req *mcp.CallToolRequest, input dto.Input) (*mcp.CallToolResult, dto.Output, error) {
	result := &mcp.CallToolResult{}

	if err := t.validateInput(input); err != nil {
		return result, dto.Output{}, fmt.Errorf("input validation failed: %w", err)
	}

	searchInput := t.convertInput(input)

	searchResult, err := t.repository.SearchAdAccounts(ctx, searchInput)
	if err != nil {
		return result, dto.Output{}, fmt.Errorf("failed to search ad accounts: %w", err)
	}

	output := dto.Output{
		Elements: searchResult.Elements,
		Paging:   searchResult.Paging,
	}

	return result, output, nil
}

func (t *Tool) validateInput(input dto.Input) error {
	// Trim and validate accountID if present
	input.AccountID = strings.TrimSpace(input.AccountID)
	if input.AccountID != "" {
		if err := validateNumericID(input.AccountID); err != nil {
			return fmt.Errorf("accountID: %w", err)
		}
	}

	for i, id := range input.AccountIDs {
		trimmed := strings.TrimSpace(id)
		if trimmed == "" {
			return fmt.Errorf("accountIDs[%d] cannot be empty", i)
		}
		if err := validateNumericID(trimmed); err != nil {
			return fmt.Errorf("accountIDs[%d]: %w", i, err)
		}
		input.AccountIDs[i] = trimmed
	}

	for i, reference := range input.References {
		trimmed := strings.TrimSpace(reference)
		if trimmed == "" {
			return fmt.Errorf("references[%d] cannot be empty", i)
		}
		if !strings.HasPrefix(trimmed, "urn:li:organization:") && !strings.HasPrefix(trimmed, "urn:li:person:") {
			return fmt.Errorf("references[%d] must start with urn:li:organization: or urn:li:person:", i)
		}
		input.References[i] = trimmed
	}

	for i, name := range input.Names {
		trimmed := strings.TrimSpace(name)
		if trimmed == "" {
			return fmt.Errorf("names[%d] cannot be empty", i)
		}
		input.Names[i] = trimmed
	}

	validStatuses := map[string]bool{
		"DRAFT":    true,
		"ACTIVE":   true,
		"PAUSED":   true,
		"ARCHIVED": true,
	}
	for _, status := range input.Status {
		if !validStatuses[status] {
			return fmt.Errorf("invalid status: %s", status)
		}
	}

	if input.SortField != "" {
		validSortFields := map[string]bool{
			"id":               true,
			"name":             true,
			"createdTime":      true,
			"lastModifiedTime": true,
		}
		if !validSortFields[input.SortField] {
			return fmt.Errorf("invalid sortField: %s", input.SortField)
		}
	}

	if input.SortOrder != "" && input.SortOrder != "ASCENDING" && input.SortOrder != "DESCENDING" {
		return fmt.Errorf("sortOrder must be ASCENDING or DESCENDING")
	}

	if input.Count != nil {
		if *input.Count < 1 {
			return fmt.Errorf("count must be greater than 0")
		}
		if *input.Count > 1000 {
			return fmt.Errorf("count cannot exceed 1000")
		}
	}

	if input.Start != nil && *input.Start < 0 {
		return fmt.Errorf("start must be non-negative")
	}

	return nil
}

func (t *Tool) convertInput(input dto.Input) adaccounts.SearchInput {
	accountIDs := append([]string{}, input.AccountIDs...)
	if input.AccountID != "" {
		accountIDs = append(accountIDs, input.AccountID)
	}

	searchInput := adaccounts.SearchInput{
		Status:     input.Status,
		Test:       input.Test,
		AccountIDs: accountIDs,
		References: input.References,
		Names:      input.Names,
		SortField:  input.SortField,
		SortOrder:  input.SortOrder,
	}

	if input.Start != nil {
		searchInput.Start = *input.Start
	}

	if input.Count != nil {
		searchInput.Count = *input.Count
	}

	return searchInput
}

func validateNumericID(value string) error {
	if value == "" {
		return fmt.Errorf("value cannot be empty")
	}
	for _, r := range value {
		if r < '0' || r > '9' {
			return fmt.Errorf("value must contain only digits")
		}
	}
	return nil
}
