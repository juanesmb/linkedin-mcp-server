package app

import (
	"os"
	"strings"
)

type Configs struct {
	LinkedInConfigs LinkedInConfigs
	AuthConfig      AuthConfig
	GatewayConfig   GatewayConfig
	ServerConfig    ServerConfig
}

type LinkedInConfigs struct {
	BaseURL string
}

type AuthConfig struct {
	ClerkIssuer      string
	ClerkJWKSURL     string
	ClerkAudience    string
	RequiredScope    string
	AuthorizationURL string
}

type GatewayConfig struct {
	BaseURL        string
	InternalSecret string
	ConnectURL     string
}

type ServerConfig struct {
	BindAddress string
	Path        string
	PublicURL   string
}

func readConfigs() Configs {
	host := strings.TrimSpace(envOrDefault("MCP_SERVER_HOST", "0.0.0.0"))
	port := strings.TrimSpace(envOrDefault("PORT", "8080"))
	// Go 1.22+ ServeMux treats "METHOD /path" patterns; a trailing space in "/mcp "
	// makes "/mcp" parse as an invalid HTTP method and panics at registration.
	path := strings.TrimSpace(envOrDefault("MCP_SERVER_PATH", "/mcp"))
	if path == "" {
		path = "/mcp"
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	clerkIssuer := strings.TrimSpace(os.Getenv("CLERK_ISSUER"))
	authorizationURL := strings.TrimSpace(os.Getenv("AUTHORIZATION_SERVER_URL"))
	if authorizationURL == "" {
		authorizationURL = clerkIssuer
	}

	return Configs{
		LinkedInConfigs: LinkedInConfigs{
			BaseURL: "https://api.linkedin.com/rest",
		},
		AuthConfig: AuthConfig{
			ClerkIssuer:      clerkIssuer,
			ClerkJWKSURL:     strings.TrimSpace(os.Getenv("CLERK_JWKS_URL")),
			ClerkAudience:    strings.TrimSpace(os.Getenv("CLERK_AUDIENCE")),
			RequiredScope:    strings.TrimSpace(os.Getenv("MCP_REQUIRED_SCOPE")),
			AuthorizationURL: authorizationURL,
		},
		GatewayConfig: GatewayConfig{
			BaseURL:        gatewayBaseURL(),
			InternalSecret: gatewayInternalSecret(),
			ConnectURL:     deriveConnectURL(gatewayBaseURL(), connectURLExplicit()),
		},
		ServerConfig: ServerConfig{
			BindAddress: host + ":" + port,
			Path:        path,
			PublicURL:   strings.TrimSpace(os.Getenv("PUBLIC_BASE_URL")),
		},
	}
}

func gatewayBaseURL() string {
	if v := strings.TrimSpace(os.Getenv("GATEWAY_BASE_URL")); v != "" {
		return v
	}
	return strings.TrimSpace(os.Getenv("JUMON_GATEWAY_BASE_URL"))
}

func gatewayInternalSecret() string {
	if v := strings.TrimSpace(os.Getenv("GATEWAY_INTERNAL_SECRET")); v != "" {
		return v
	}
	return strings.TrimSpace(os.Getenv("JUMON_GATEWAY_INTERNAL_SECRET"))
}

func connectURLExplicit() string {
	if v := strings.TrimSpace(os.Getenv("GATEWAY_CONNECT_URL")); v != "" {
		return v
	}
	return strings.TrimSpace(os.Getenv("JUMON_CONNECT_URL"))
}

func deriveConnectURL(gatewayBaseURL, explicitConnectURL string) string {
	if explicitConnectURL != "" {
		return explicitConnectURL
	}
	if v := strings.TrimSpace(os.Getenv("JUMON_CONNECT_URL")); v != "" {
		return v
	}
	if gatewayBaseURL == "" {
		return "/connections"
	}
	return strings.TrimRight(gatewayBaseURL, "/") + "/connections"
}

func envOrDefault(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
