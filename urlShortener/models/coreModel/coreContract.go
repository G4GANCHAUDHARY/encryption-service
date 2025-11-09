package coreModel

import "github.com/lib/pq"

type IGenerateUrlReq interface {
	GetLongUrl() string
	GetIsCustomUrl() bool
	GetCustomUrl() string
}

type IGetUrlReq interface {
	GetShortCode() string
}

type IGetList interface {
	GetPageSize() int
	GetFilters() pq.StringArray
	GetCursor() int
	GetSortingOrder() string
	GetOrderBy() string
}
