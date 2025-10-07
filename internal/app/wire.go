package app

import (
	"linkedin-mcp/internal/infrastructure/api/campaigns"
	"linkedin-mcp/internal/infrastructure/http"
	"linkedin-mcp/internal/infrastructure/resources/analytics/metrics"
	"linkedin-mcp/internal/infrastructure/resources/analytics/queryparameters"
	"linkedin-mcp/internal/infrastructure/tools/searchcampaigns"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func initServer(configs Configs) *mcp.Server {
	server := mcp.NewServer(&mcp.Implementation{Name: "LinkedIn", Version: "v1.0.0"}, nil)

	searchCampaignsTool := initSearchCampaignsTool(configs)

	mcp.AddTool(server, &mcp.Tool{Name: "search_campaigns", Description: "Search for LinkedIn ad campaigns"}, searchCampaignsTool.SearchCampaigns)

	// Initialize and register the analytics parameters resource
	analyticsResource := initAnalyticsResource()
	server.AddResource(&mcp.Resource{
		URI:         "linkedin://analytics/parameters",
		Name:        "LinkedIn Analytics Query Parameters",
		Description: "JSON schema containing LinkedIn analytics API query parameters and their descriptions",
		MIMEType:    "application/json",
	}, analyticsResource.ReadResource)

	// Initialize and register the analytics metrics resource
	analyticsMetricsResource := initAnalyticsMetricsResource()
	server.AddResource(&mcp.Resource{
		URI:         "linkedin://analytics/metrics",
		Name:        "LinkedIn Analytics Metrics",
		Description: "JSON schema containing LinkedIn analytics API metrics and their descriptions",
		MIMEType:    "application/json",
	}, analyticsMetricsResource.ReadResource)

	return server
}

func initSearchCampaignsTool(configs Configs) *searchcampaigns.Tool {
	httpClient := http.NewClient(nil)

	queryBuilder := campaigns.NewQueryBuilder(configs.LinkedInConfigs.BaseURL,
		configs.LinkedInConfigs.AccountID,
		configs.LinkedInConfigs.Version,
		configs.LinkedInConfigs.AccessToken,
	)

	campaignsRepository := campaigns.NewRepository(httpClient, queryBuilder)

	return searchcampaigns.NewTool(campaignsRepository)
}

func initAnalyticsResource() *queryparameters.Resource {
	return queryparameters.NewResource()
}

func initAnalyticsMetricsResource() *metrics.Resource {
	return metrics.NewResource()
}
