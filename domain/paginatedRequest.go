package domain

// PaginatedRequest : Struct for request a page of data
type PaginatedRequest struct {
	PageSize int `json:"pageSize"`
	Current  int `json:"current"`
}
