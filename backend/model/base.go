package model

// PageInfoOutput represents information related to a particular page output
type PageInfoOutput struct {
	Page       int `json:"page"`
	PageSize   int `json:"pageSize"`
	TotalCount int `json:"totalCount"`
	PageCount  int `json:"pageCount"`
}

// PageOutput is a wrapper for any output that requires pagination information
type PageOutput struct {
	Items    interface{}    `json:"items"`
	PageInfo PageInfoOutput `json:"pageInfo"`
}
