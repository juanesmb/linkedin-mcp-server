package adaccounts

type LinkedInResponse struct {
	Elements []map[string]any `json:"elements"`
	Paging   map[string]any   `json:"paging"`
}

type SearchResult struct {
	Elements []map[string]any `json:"elements"`
	Paging   map[string]any   `json:"paging,omitempty"`
}
