package repository

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	request "github.com/machtwatch/catalyst-go-skeleton/domain/common/request/api"
)

type PaginationResponse struct {
	TotalData int
	TotalPage int
}

func SetPagination(q sq.SelectBuilder, opt request.PaginationRequest) sq.SelectBuilder {
	offset := (opt.Page - 1) * opt.Size
	limit := opt.Size

	q = q.Limit(uint64(limit))
	q = q.Offset(uint64(offset))

	if opt.SortBy == nil {
		defaultSortBy := "id"
		opt.SortBy = &defaultSortBy
	}

	if opt.SortType == nil {
		defaultSortType := "asc, id" // include primary key (id) to solve duplicates in sorted column
		opt.SortType = &defaultSortType
	} else {
		sortType := fmt.Sprintf("%s, id", *opt.SortType)
		opt.SortType = &sortType
	}

	q = q.OrderBy(fmt.Sprintf("%v %v", *opt.SortBy, *opt.SortType))

	return q
}

func SetFilter(q sq.SelectBuilder, filters []interface{}) sq.SelectBuilder {
	q = q.Where("deleted_at IS NULL")

	if len(filters) == 0 {
		return q
	}

	for _, filter := range filters {
		q = q.Where(filter)
	}

	return q
}
