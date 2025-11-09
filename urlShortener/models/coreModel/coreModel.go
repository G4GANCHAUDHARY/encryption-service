package coreModel

import "github.com/lib/pq"

type GenerateUrlReq struct {
	LongUrl     string
	IsCustomUrl bool
	CustomUrl   string
}

func (r *GenerateUrlReq) GetLongUrl() string {
	return r.LongUrl
}

func (r *GenerateUrlReq) GetIsCustomUrl() bool {
	return r.IsCustomUrl
}

func (r *GenerateUrlReq) GetCustomUrl() string {
	return r.CustomUrl
}

func (r *GenerateUrlReq) SetLongUrl(url string) {
	r.LongUrl = url
}

func (r *GenerateUrlReq) SetIsCustomUrl(isCustom bool) {
	r.IsCustomUrl = isCustom
}

func (r *GenerateUrlReq) SetCustomUrl(url string) {
	r.CustomUrl = url
}

type GetUrlReq struct {
	ShortCode string
}

func (r *GetUrlReq) GetShortCode() string {
	return r.ShortCode
}

type GetListRequestModel struct {
	Filters      pq.StringArray
	PageSize     int
	Cursor       int
	SortingOrder string
	OrderBy      string
}

func (gl GetListRequestModel) GetPageSize() int {
	return gl.PageSize
}

func (gl GetListRequestModel) GetFilters() pq.StringArray {
	return gl.Filters
}

func (gl GetListRequestModel) GetCursor() int {
	return gl.Cursor
}

func (gl GetListRequestModel) GetSortingOrder() string {
	return gl.SortingOrder
}

func (gl GetListRequestModel) GetOrderBy() string {
	return gl.OrderBy
}
