package dto

import "linkedin-mcp/internal/infrastructure/api/creatives"

type Output struct {
	Elements []creatives.NormalizedCreative `json:"elements" jsonschema:"Normalized creative metadata rows"`
	Paging   creatives.PagingSummary        `json:"paging,omitempty" jsonschema:"Pagination summary"`
}
