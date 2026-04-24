package gateway

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"linkedin-mcp/internal/infrastructure/api"
)

const (
	headerGatewaySecret = "x-gateway-secret"
)

type Client struct {
	httpClient     api.Client
	baseURL        string
	internalSecret string
}

func NewClient(httpClient api.Client, baseURL, internalSecret string) *Client {
	return &Client{
		httpClient:     httpClient,
		baseURL:        strings.TrimRight(baseURL, "/"),
		internalSecret: internalSecret,
	}
}

func (c *Client) GetLinkedInConnection(ctx context.Context, userID string) (*api.Response, error) {
	path := fmt.Sprintf("%s/api/internal/connections/linkedin/current?userId=%s", c.baseURL, url.QueryEscape(userID))
	return c.httpClient.Get(ctx, path, c.authHeaders())
}

func (c *Client) ProxyLinkedIn(ctx context.Context, userID, resourcePath string, query map[string]string, headers map[string]string) (*api.Response, error) {
	path := fmt.Sprintf("%s/api/internal/providers/linkedin/proxy", c.baseURL)
	body := map[string]interface{}{
		"userId": userID,
		"method": "GET",
		"path":   strings.TrimLeft(resourcePath, "/"),
		"query":  query,
	}
	if len(headers) > 0 {
		body["headers"] = headers
	}
	return c.httpClient.Post(ctx, path, body, c.authHeaders())
}

func (c *Client) RefreshLinkedIn(ctx context.Context, userID string) (*api.Response, error) {
	path := fmt.Sprintf("%s/api/internal/providers/linkedin/refresh", c.baseURL)
	body := map[string]string{
		"userId": userID,
	}
	return c.httpClient.Post(ctx, path, body, c.authHeaders())
}

// ProxyLinkedInOrRefresh proxies the REST call through Jumon and retries once after a token
// refresh if the gateway responds with 401. Call [Client.GetLinkedInConnection] first if you
// need to verify the user has a LinkedIn connection before proxying.
func (c *Client) ProxyLinkedInOrRefresh(ctx context.Context, userID, resourcePath string, query map[string]string, headers map[string]string) (*api.Response, error) {
	resp, err := c.ProxyLinkedIn(ctx, userID, resourcePath, query, headers)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 401 {
		return resp, nil
	}
	if _, err := c.RefreshLinkedIn(ctx, userID); err != nil {
		return resp, nil
	}
	return c.ProxyLinkedIn(ctx, userID, resourcePath, query, headers)
}

func (c *Client) authHeaders() map[string]string {
	return map[string]string{
		headerGatewaySecret: c.internalSecret,
		"Content-Type":      "application/json",
		"Accept":            "application/json",
	}
}
