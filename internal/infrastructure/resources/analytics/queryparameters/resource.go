package queryparameters

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const schemaURL = "https://learn.microsoft.com/en-us/linkedin/marketing/integrations/ads-reporting/ads-reporting-schema?view=li-lms-2026-03#analytics-finder-query-parameters"

// Resource handles LinkedIn analytics parameters as an MCP resource
type Resource struct{}

// NewResource creates a new Resource instance
func NewResource() *Resource {
	return &Resource{}
}

func (r *Resource) ReadResource(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	if req.Params.URI != "linkedin://analytics/parameters" {
		return nil, fmt.Errorf("resource not found: %s", req.Params.URI)
	}

	payload := map[string]any{
		"source":  schemaURL,
		"purpose": "Canonical LinkedIn Ads Reporting analytics finder query parameters documentation.",
		"notes": []string{
			"LinkedIn documentation is the source of truth for available analytics finder parameters.",
			"Supported pivots, sort fields, and parameter requirements can change over time.",
			"Use exact parameter names and constraints from the linked section before calling get_analytics.",
		},
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize LinkedIn analytics parameters resource payload: %w", err)
	}

	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{
			{
				URI:  req.Params.URI,
				Text: string(data),
			},
		},
	}, nil
}
