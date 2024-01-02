package presentation

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
)

type Paging struct {
	OffSet  int64  `form:"offset"`
	Limit   int64  `form:"limit"`
	OrderBy string `form:"order_by"`
	SortBy  string `form:"sort_by"`
}

func (ref Paging) ToDomain() values.Paging {
	return values.Paging{
		OffSet:  ref.OffSet,
		Limit:   ref.Limit,
		OrderBy: ref.OrderBy,
		SortBy:  ref.SortBy,
	}
}

type PagingResponse struct {
	NumberOfPages int64 `json:"number_of_pages"`
	TotalCount    int64 `json:"total_count"`
}
