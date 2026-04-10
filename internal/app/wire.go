package app

import (
	"log"
	"os"
	"strings"

	"linkedin-mcp/internal/infrastructure/api"
	adaccountsapi "linkedin-mcp/internal/infrastructure/api/adaccounts"
	"linkedin-mcp/internal/infrastructure/api/campaigns"
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

const (
	serverInstructionsFilePath = "internal/app/instructions/server_instructions.md"
)

type Components struct {
	httpClient api.Client
	logger     infrastructurelog.Logger
}

func initServer(configs Configs) *mcp.Server {
	components := initCommonComponents()

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
	}, initSearchAdAccountsTool(configs, *components).SearchAdAccounts)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "search_campaigns",
		Description: "Search for LinkedIn ad campaigns. Requires the accountID argument.",
	}, initSearchCampaignsTool(configs, *components).SearchCampaigns)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_analytics",
		Description: "Get LinkedIn ad analytics data. Requires accountID and should be used after reading analytics resources.",
	}, initReportingTool(configs, *components).GetAnalytics)

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

func initCommonComponents() *Components {
	return &Components{
		httpClient: http.NewClient(nil),
		logger:     locallogger.NewLogger(),
	}
}

func initSearchCampaignsTool(configs Configs, components Components) *searchcampaigns.Tool {
	queryBuilder := campaigns.NewQueryBuilder(configs.LinkedInConfigs.BaseURL,
		configs.LinkedInConfigs.Version,
		configs.LinkedInConfigs.AccessToken,
	)

	campaignsRepository := campaigns.NewRepository(components.httpClient, queryBuilder, components.logger)

	return searchcampaigns.NewTool(campaignsRepository)
}

func initSearchAdAccountsTool(configs Configs, components Components) *searchadaccounts.Tool {
	queryBuilder := adaccountsapi.NewQueryBuilder(configs.LinkedInConfigs.BaseURL,
		configs.LinkedInConfigs.Version,
		configs.LinkedInConfigs.AccessToken,
	)

	repository := adaccountsapi.NewRepository(components.httpClient, queryBuilder, components.logger)

	return searchadaccounts.NewTool(repository)
}

func initReportingTool(configs Configs, components Components) *getanalytics.Tool {
	queryBuilder := reportingapi.NewQueryBuilder(configs.LinkedInConfigs.BaseURL,
		configs.LinkedInConfigs.Version,
		configs.LinkedInConfigs.AccessToken,
	)

	reportingRepository := reportingapi.NewRepository(components.httpClient, queryBuilder, components.logger)

	return getanalytics.NewTool(reportingRepository)
}

func initAnalyticsResource() *queryparameters.Resource {
	return queryparameters.NewResource()
}

func initAnalyticsMetricsResource() *metrics.Resource {
	return metrics.NewResource()
}

func loadServerInstructions() string {
	content, err := os.ReadFile(serverInstructionsFilePath)
	if err != nil {
		log.Fatalf("failed to read MCP instructions file %q: %v", serverInstructionsFilePath, err)
	}

	instructions := strings.TrimSpace(string(content))
	if instructions == "" {
		log.Fatalf("MCP instructions file is empty: %q", serverInstructionsFilePath)
	}

	return instructions
}
