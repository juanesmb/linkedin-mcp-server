package queryparameters

import (
	"context"
	"embed"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

//go:embed dto/linkedin_analytics_parameters.json
var linkedinAnalyticsParametersFile embed.FS

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

	// Read the embedded JSON file
	data, err := linkedinAnalyticsParametersFile.ReadFile("dto/linkedin_analytics_parameters.json")
	if err != nil {
		return nil, fmt.Errorf("failed to read LinkedIn analytics parameters: %w", err)
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
