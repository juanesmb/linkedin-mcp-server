package metrics

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const schemaURL = "https://learn.microsoft.com/en-us/linkedin/marketing/integrations/ads-reporting/ads-reporting-schema?view=li-lms-2026-03#metrics-available"

type Resource struct{}

func NewResource() *Resource {
	return &Resource{}
}

func (r *Resource) ReadResource(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	if req.Params.URI != "linkedin://analytics/metrics" {
		return nil, fmt.Errorf("resource not found: %s", req.Params.URI)
	}

	payload := map[string]any{
		"source":  schemaURL,
		"purpose": "Canonical LinkedIn Ads Reporting metrics documentation.",
		"notes": []string{
			"LinkedIn documentation is the source of truth for available analytics metrics.",
			"Metric names and availability can change over time.",
			"Use exact metric field names from the linked table when building get_analytics requests.",
		},
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize LinkedIn analytics metrics resource payload: %w", err)
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
