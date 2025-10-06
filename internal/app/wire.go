package app

import (
	"linkedin-mcp/internal/infrastructure/api/campaigns"
	"linkedin-mcp/internal/infrastructure/http"
	"linkedin-mcp/internal/infrastructure/tools/searchcampaigns"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func initServer(configs Configs) *mcp.Server {
	server := mcp.NewServer(&mcp.Implementation{Name: "LinkedIn", Version: "v1.0.0"}, nil)

	searchCampaignsTool := initSearchCampaignsTool(configs)

	mcp.AddTool(server, &mcp.Tool{Name: "search_campaigns", Description: "Search for LinkedIn ad campaigns"}, searchCampaignsTool.SearchCampaigns)

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
