package app

import "os"

type Configs struct {
	LinkedInConfigs LinkedInConfigs
}

type LinkedInConfigs struct {
	AccessToken string
	AccountID   string
	BaseURL     string
	Version     string
}

func readConfigs() Configs {
	linkedInConfigs := LinkedInConfigs{
		AccessToken: os.Getenv("LINKEDIN_ACCESS_TOKEN"),
		AccountID:   os.Getenv("LINKEDIN_ACCOUNT_ID"),
		BaseURL:     "https://api.linkedin.com/rest",
		Version:     "202505",
	}

	return Configs{
		LinkedInConfigs: linkedInConfigs,
	}
}
