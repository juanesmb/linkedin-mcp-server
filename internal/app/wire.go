package app

import (
	"log"
	"os"
	"strings"

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

func initServer(configs Configs) *mcp.Server {
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
	}, initSearchAdAccountsTool(configs).SearchAdAccounts)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "search_campaigns",
		Description: "Search for LinkedIn ad campaigns. Requires the accountID argument.",
	}, initSearchCampaignsTool(configs).SearchCampaigns)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_analytics",
		Description: "Get LinkedIn ad analytics data. Requires accountID and should be used after reading analytics resources.",
	}, initReportingTool(configs).GetAnalytics)

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

func initSearchCampaignsTool(configs Configs) *searchcampaigns.Tool {
	httpClient := http.NewClient(nil)
	logger := resolveLogger(configs)

	queryBuilder := campaigns.NewQueryBuilder(configs.LinkedInConfigs.BaseURL,
		configs.LinkedInConfigs.Version,
		configs.LinkedInConfigs.AccessToken,
	)

	campaignsRepository := campaigns.NewRepository(httpClient, queryBuilder, logger)

	return searchcampaigns.NewTool(campaignsRepository)
}

func initSearchAdAccountsTool(configs Configs) *searchadaccounts.Tool {
	httpClient := http.NewClient(nil)
	logger := resolveLogger(configs)

	queryBuilder := adaccountsapi.NewQueryBuilder(configs.LinkedInConfigs.BaseURL,
		configs.LinkedInConfigs.Version,
		configs.LinkedInConfigs.AccessToken,
	)

	repository := adaccountsapi.NewRepository(httpClient, queryBuilder, logger)

	return searchadaccounts.NewTool(repository)
}

func initReportingTool(configs Configs) *getanalytics.Tool {
	httpClient := http.NewClient(nil)
	logger := resolveLogger(configs)

	queryBuilder := reportingapi.NewQueryBuilder(configs.LinkedInConfigs.BaseURL,
		configs.LinkedInConfigs.Version,
		configs.LinkedInConfigs.AccessToken,
	)

	reportingRepository := reportingapi.NewRepository(httpClient, queryBuilder, logger)

	return getanalytics.NewTool(reportingRepository)
}

func resolveLogger(configs Configs) infrastructurelog.Logger {
	return locallogger.NewLogger()
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
