package adaccounts

// SearchInput represents the supported filters for querying LinkedIn Ad Accounts.
// It mirrors the REST parameters described in LinkedIn's documentation:
// https://learn.microsoft.com/es-mx/linkedin/marketing/integrations/ads/account-structure/create-and-manage-accounts?view=li-lms-2025-09&tabs=http#search-for-accounts
type SearchInput struct {
	// Status filters by ad account status (e.g., DRAFT, ACTIVE).
	Status []string

	// Test determines whether to fetch test accounts (true), non-test accounts (false), or both when nil.
	Test *bool

	// AccountIDs filters by specific ad account IDs.
	AccountIDs []string

	// References filters by associated entity URNs (organization or person).
	References []string

	// Names filters by exact ad account names.
	Names []string

	// SortField orders the results using LinkedIn-supported sort fields (e.g., id, name, createdTime).
	SortField string

	// SortOrder complements SortField (ASCENDING or DESCENDING).
	SortOrder string

	// Start controls pagination starting index.
	Start int

	// Count controls number of results per page (max 1000).
	Count int
}
