package gateway

import (
	"fmt"
	"net/url"
	"strings"
)

// ParseLinkedInRESTProxyTarget splits a full LinkedIn REST URL into a gateway-relative
// resource path (under /rest/) and a flat query map for the Jumon proxy body.
func ParseLinkedInRESTProxyTarget(requestURL string) (resourcePath string, query map[string]string, err error) {
	parsed, err := url.Parse(requestURL)
	if err != nil {
		return "", nil, fmt.Errorf("invalid request URL: %w", err)
	}

	resourcePath = strings.TrimPrefix(parsed.Path, "/rest/")
	if resourcePath == "" {
		return "", nil, fmt.Errorf("empty LinkedIn REST resource path")
	}

	query = make(map[string]string)
	for key, values := range parsed.Query() {
		if len(values) == 0 {
			continue
		}
		query[key] = values[0]
	}

	return resourcePath, query, nil
}
