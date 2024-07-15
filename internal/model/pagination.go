package model

type PaginationParams struct {
	PageSize  int    `form:"page_size,omitempty"`
	NextToken string `form:"next_token,omitempty"`
}

type PaginationResponse[T any] struct {
	Total     int64  `json:"total"`
	Items     []T    `json:"items"`
	NextToken string `json:"next_token,omitempty"`
}

func (p *PaginationParams) Normalize() *PaginationParams {
	if p.PageSize == 0 {
		p.PageSize = 10
	}

	if p.PageSize > 40 {
		p.PageSize = 40
	}

	return p
}
