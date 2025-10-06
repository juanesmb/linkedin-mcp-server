package campaigns

type LinkedInResponse struct {
	Elements []map[string]any `json:"elements"`
	Paging   map[string]any   `json:"paging"`
}

type SearchResult struct {
	Elements []map[string]any `json:"elements"`
	Metadata struct {
		NextPageToken string `json:"nextPageToken,omitempty"`
	} `json:"metadata,omitempty"`
}
