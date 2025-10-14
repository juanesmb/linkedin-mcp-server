package dto

type Output struct {
	Elements []map[string]any `json:"elements" jsonschema:"Ad account results"`
	Paging   map[string]any   `json:"paging,omitempty" jsonschema:"LinkedIn paging metadata"`
}
