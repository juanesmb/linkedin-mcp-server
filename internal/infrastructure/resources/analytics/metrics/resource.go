package metrics

import (
	"context"
	"embed"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

//go:embed dto/linkedin_analytics_metrics.json
var linkedinAnalyticsMetricsFile embed.FS

type Resource struct{}

func NewResource() *Resource {
	return &Resource{}
}

func (r *Resource) ReadResource(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	if req.Params.URI != "linkedin://analytics/metrics" {
		return nil, fmt.Errorf("resource not found: %s", req.Params.URI)
	}

	// Read the embedded JSON file
	data, err := linkedinAnalyticsMetricsFile.ReadFile("dto/linkedin_analytics_metrics.json")
	if err != nil {
		return nil, fmt.Errorf("failed to read LinkedIn analytics metrics: %w", err)
	}

	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{
			{
				URI:      req.Params.URI,
				MIMEType: "application/json",
				Text:     string(data),
			},
		},
	}, nil
}
