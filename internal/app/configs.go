package app

import (
	"os"
)

type Configs struct {
	LinkedInConfigs LinkedInConfigs
	ServerConfig    ServerConfig
}

type LinkedInConfigs struct {
	AccessToken string
	BaseURL     string
	Version     string
}

type ServerConfig struct {
	BindAddress string
	Path        string
}

func readConfigs() Configs {
	return Configs{
		LinkedInConfigs: LinkedInConfigs{
			AccessToken: os.Getenv("LINKEDIN_ACCESS_TOKEN"),
			BaseURL:     "https://api.linkedin.com/rest",
			// Default matches LinkedIn Marketing API "Latest Version" (YYYYMM header value).
			// See: https://learn.microsoft.com/en-us/linkedin/marketing/versioning
			Version: "202603",
		},
		ServerConfig: ServerConfig{
			BindAddress: ":8080",
			Path:        "/mcp",
		},
	}
}
