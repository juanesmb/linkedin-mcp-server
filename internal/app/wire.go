package app

import (
	adaccountsapi "linkedin-mcp/internal/infrastructure/api/adaccounts"
	"linkedin-mcp/internal/infrastructure/api/campaigns"
	reportingapi "linkedin-mcp/internal/infrastructure/api/reporting"
	"linkedin-mcp/internal/infrastructure/http"
	"linkedin-mcp/internal/infrastructure/log"
	locallogger "linkedin-mcp/internal/infrastructure/log/local"
	"linkedin-mcp/internal/infrastructure/prompts/accountid"
	"linkedin-mcp/internal/infrastructure/prompts/systemguidelines"
	"linkedin-mcp/internal/infrastructure/resources/analytics/metrics"
	"linkedin-mcp/internal/infrastructure/resources/analytics/queryparameters"
	"linkedin-mcp/internal/infrastructure/tools/getanalytics"
	"linkedin-mcp/internal/infrastructure/tools/searchadaccounts"
	"linkedin-mcp/internal/infrastructure/tools/searchcampaigns"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func initServer(configs Configs) *mcp.Server {
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "LinkedIn",
		Version: "v1.0.0",
		Title:   "LinkedIn Advertising MCP server. Use 'system_guidelines' prompt first.",
	}, nil)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "search_ad_accounts",
		Description: "Search for LinkedIn ad accounts without needing an ID. REQUIRES: Use 'system_guidelines' prompt first.",
	}, initSearchAdAccountsTool(configs).SearchAdAccounts)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "search_campaigns",
		Description: "Search for LinkedIn ad campaigns. REQUIRES: Use 'system_guidelines' prompt first.",
	}, initSearchCampaignsTool(configs).SearchCampaigns)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_analytics",
		Description: "Get LinkedIn ad analytics data. REQUIRES: Use 'system_guidelines' prompt first.",
	}, initReportingTool(configs).GetAnalytics)

	analyticsResource := initAnalyticsResource()
	server.AddResource(&mcp.Resource{
		URI:         "linkedin://analytics/parameters",
		Name:        "LinkedIn Analytics Query Parameters",
		Description: "JSON schema containing LinkedIn analytics API query parameters and their descriptions",
		MIMEType:    "application/json",
	}, analyticsResource.ReadResource)

	analyticsMetricsResource := initAnalyticsMetricsResource()
	server.AddResource(&mcp.Resource{
		URI:         "linkedin://analytics/metrics",
		Name:        "LinkedIn Analytics Metrics",
		Description: "JSON schema containing LinkedIn analytics API metrics and their descriptions",
		MIMEType:    "application/json",
	}, analyticsMetricsResource.ReadResource)

	systemGuidelinesPrompt := initSystemGuidelinesPrompt()
	server.AddPrompt(&mcp.Prompt{
		Name:        "system_guidelines",
		Description: "System guidelines for using the LinkedIn MCP server. READ THIS FIRST before using any tools or prompts.",
		Arguments:   []*mcp.PromptArgument{},
	}, systemGuidelinesPrompt.GetPrompt)

	accountIDPrompt := initAccountIDPrompt()
	server.AddPrompt(&mcp.Prompt{
		Name:        "linkedin_account_id_required",
		Description: "Instructions to request LinkedIn Account ID before using tools",
		Arguments: []*mcp.PromptArgument{
			{
				Name:        "accountID",
				Description: "The LinkedIn Ad Account ID provided by the user",
				Required:    true,
			},
		},
	}, accountIDPrompt.GetPrompt)

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

func resolveLogger(configs Configs) log.Logger {
	return locallogger.NewLogger()
}

func initAnalyticsResource() *queryparameters.Resource {
	return queryparameters.NewResource()
}

func initAnalyticsMetricsResource() *metrics.Resource {
	return metrics.NewResource()
}

func initSystemGuidelinesPrompt() *systemguidelines.Prompt {
	return systemguidelines.NewPrompt()
}

func initAccountIDPrompt() *accountid.Prompt {
	return accountid.NewPrompt()
}
