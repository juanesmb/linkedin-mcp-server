package dto

type Output struct {
	Elements []map[string]any `json:"elements" jsonschema:"Campaign results"`
	Metadata struct {
		NextPageToken string `json:"nextPageToken,omitempty" jsonschema:"Cursor for next page if available"`
	} `json:"metadata,omitempty" jsonschema:"Metadata containing pagination info"`
}
