package campaigns

type SearchInput struct {
	CampaignGroupURNs      []string
	AssociatedEntityValues []string
	CampaignURNs           []string
	Status                 []string
	Type                   []string
	Name                   []string
	Test                   *bool
	SortOrder              string
	PageSize               int
	PageToken              string
}
