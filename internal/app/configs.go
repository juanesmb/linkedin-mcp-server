package app

import "os"

type Configs struct {
	LinkedInConfigs LinkedInConfigs
}

type LinkedInConfigs struct {
	AccessToken string
	BaseURL     string
	Version     string
}

func readConfigs() Configs {
	linkedInConfigs := LinkedInConfigs{
		AccessToken: os.Getenv("LINKEDIN_ACCESS_TOKEN"),
		BaseURL:     "https://api.linkedin.com/rest",
		Version:     "202505",
	}

	return Configs{
		LinkedInConfigs: linkedInConfigs,
	}
}
