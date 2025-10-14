package dto

type Input struct {
	AccountID  string   `json:"accountID,omitempty" jsonschema:"Optional LinkedIn Ad Account ID (numeric value, e.g., 512345678)"`
	AccountIDs []string `json:"accountIDs" jsonschema:"Optional additional LinkedIn Ad Account IDs to include (numeric values)"`
	Status     []string `json:"status" jsonschema:"Filter by account status (e.g., DRAFT, ACTIVE, PAUSED, ARCHIVED)"`
	Test       *bool    `json:"test,omitempty" jsonschema:"Filter by test accounts: true, false, or omit for both"`
	References []string `json:"references" jsonschema:"Filter by associated entity URNs (urn:li:organization:{id} or urn:li:person:{id})"`
	Names      []string `json:"names" jsonschema:"Filter by ad account names (exact match)"`
	Start      *int     `json:"start,omitempty" jsonschema:"Pagination start offset (>=0)"`
	Count      *int     `json:"count,omitempty" jsonschema:"Results per page (1-1000). Default 100"`
}
