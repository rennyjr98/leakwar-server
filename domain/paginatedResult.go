package domain

// PaginatedResult : Returns a collection of entities by pages
type PaginatedResult struct {
	PageSize int            `json:"pageSize"`
	Current  int            `json:"current"`
	Items    []*interface{} `json:"items"`
}
