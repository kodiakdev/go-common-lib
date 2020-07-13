package commonreq

import (
	"strconv"

	"github.com/emicklei/go-restful"
)

const (
	PathParamTenant     = "tenant"
	PathParamUserID     = "userId"
	QueryParamLimit     = "limit"
	QueryParamSort      = "sort"
	QueryParamPage      = "page"
	DefaultLimitPerPage = 30
	DefaultSortOnPage   = 0
	DefaultPage         = 0
)

//ExtractPagingSortingRequest extract limit and sort query param and convert to int
func ExtractPagingSortingRequest(req *restful.Request) PagingSortingRequest {

	limitRaw := req.QueryParameter(QueryParamLimit)
	sortRaw := req.QueryParameter(QueryParamSort)
	pageRaw := req.QueryParameter(QueryParamPage)

	limit, err := strconv.ParseInt(limitRaw, 10, 64)
	if err != nil {
		limit = DefaultLimitPerPage
	}
	sort, err := strconv.Atoi(sortRaw)
	if err != nil {
		sort = DefaultSortOnPage
	}
	page, err := strconv.ParseInt(pageRaw, 10, 64)
	if err != nil {
		page = DefaultPage
	}
	return PagingSortingRequest{
		Limit: limit,
		Sort:  sort,
		Page:  page,
	}
}

//PagingSortingRequest dto for paging-sorting request
type PagingSortingRequest struct {
	Limit int64
	Sort  int
	Page  int64
}
