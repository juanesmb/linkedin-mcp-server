package app

import (
	_ "embed"
	"log"
	"strings"

	"linkedin-mcp/internal/infrastructure/api"
	adaccountsapi "linkedin-mcp/internal/infrastructure/api/adaccounts"
	"linkedin-mcp/internal/infrastructure/api/campaigns"
	"linkedin-mcp/internal/infrastructure/api/gateway"
	reportingapi "linkedin-mcp/internal/infrastructure/api/reporting"
	"linkedin-mcp/internal/infrastructure/http"
	infrastructurelog "linkedin-mcp/internal/infrastructure/log"
	locallogger "linkedin-mcp/internal/infrastructure/log/local"
	"linkedin-mcp/internal/infrastructure/resources/analytics/metrics"
	"linkedin-mcp/internal/infrastructure/resources/analytics/queryparameters"
	"linkedin-mcp/internal/infrastructure/tools/getanalytics"
	"linkedin-mcp/internal/infrastructure/tools/searchadaccounts"
	"linkedin-mcp/internal/infrastructure/tools/searchcampaigns"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

//go:embed instructions/server_instructions.md
var serverInstructions string

type Components struct {
	httpClient    api.Client
	gatewayClient *gateway.Client
	logger        infrastructurelog.Logger
}

func initServer(configs Configs, components Components) *mcp.Server {
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "LinkedIn",
		Version: "v1.0.0",
		Title:   "LinkedIn Advertising MCP server.",
	}, &mcp.ServerOptions{
		Instructions: loadServerInstructions(),
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "search_ad_accounts",
		Description: "Search for LinkedIn ad accounts without requiring an accountID argument.",
	}, initSearchAdAccountsTool(configs, components).SearchAdAccounts)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "search_campaigns",
		Description: "Search for LinkedIn ad campaigns. Requires the accountID argument.",
	}, initSearchCampaignsTool(configs, components).SearchCampaigns)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_analytics",
		Description: "Get LinkedIn ad analytics data. Requires accountID and should be used after reading analytics resources.",
	}, initReportingTool(configs, components).GetAnalytics)

	analyticsResource := initAnalyticsResource()
	server.AddResource(&mcp.Resource{
		URI:         "linkedin://analytics/parameters",
		Name:        "LinkedIn Analytics Query Parameters",
		Description: "Reference link to LinkedIn analytics query parameters documentation",
	}, analyticsResource.ReadResource)

	analyticsMetricsResource := initAnalyticsMetricsResource()
	server.AddResource(&mcp.Resource{
		URI:         "linkedin://analytics/metrics",
		Name:        "LinkedIn Analytics Metrics",
		Description: "Reference link to LinkedIn analytics metrics documentation",
	}, analyticsMetricsResource.ReadResource)

	return server
}

func initCommonComponents(configs Configs) Components {
	httpClient := http.NewClient(nil)
	return Components{
		httpClient:    httpClient,
		gatewayClient: gateway.NewClient(httpClient, configs.GatewayConfig.BaseURL, configs.GatewayConfig.InternalSecret),
		logger:        locallogger.NewLogger(),
	}
}

func initSearchCampaignsTool(configs Configs, components Components) *searchcampaigns.Tool {
	queryBuilder := campaigns.NewQueryBuilder(configs.LinkedInConfigs.BaseURL)

	campaignsRepository := campaigns.NewRepository(components.gatewayClient, queryBuilder, components.logger)

	return searchcampaigns.NewTool(campaignsRepository)
}

func initSearchAdAccountsTool(configs Configs, components Components) *searchadaccounts.Tool {
	queryBuilder := adaccountsapi.NewQueryBuilder(configs.LinkedInConfigs.BaseURL)

	repository := adaccountsapi.NewRepository(components.gatewayClient, queryBuilder, components.logger)

	return searchadaccounts.NewTool(repository)
}

func initReportingTool(configs Configs, components Components) *getanalytics.Tool {
	queryBuilder := reportingapi.NewQueryBuilder(configs.LinkedInConfigs.BaseURL)

	reportingRepository := reportingapi.NewRepository(components.gatewayClient, queryBuilder, components.logger)

	return getanalytics.NewTool(reportingRepository)
}

func initAnalyticsResource() *queryparameters.Resource {
	return queryparameters.NewResource()
}

func initAnalyticsMetricsResource() *metrics.Resource {
	return metrics.NewResource()
}

func loadServerInstructions() string {
	instructions := strings.TrimSpace(serverInstructions)
	if instructions == "" {
		log.Fatal("embedded MCP instructions file is empty: internal/app/instructions/server_instructions.md")
	}

	return instructions
}
