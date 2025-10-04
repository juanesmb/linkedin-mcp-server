package app

import (
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
	return searchcampaigns.NewTool(configs.LinkedInConfigs.AccessToken, configs.LinkedInConfigs.AccountID,
		configs.LinkedInConfigs.BaseURL, configs.LinkedInConfigs.Version)
}
