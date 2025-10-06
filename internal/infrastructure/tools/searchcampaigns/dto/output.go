package dto

type Output struct {
	Elements []map[string]any `json:"elements" jsonschema:"Campaign results"`
	Metadata Metadata         `json:"metadata,omitempty" jsonschema:"Metadata containing pagination info"`
}

type Metadata struct {
	NextPageToken string `json:"nextPageToken,omitempty" jsonschema:"Cursor for next page if available"`
}
