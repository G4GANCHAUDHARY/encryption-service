package httpDataMapper

import (
	"github.com/G4GANCHAUDHARY/encryption-service/urlShortener/models/coreModel"
	"github.com/G4GANCHAUDHARY/encryption-service/urlShortener/models/httpModel"
	"github.com/lib/pq"
)

type IHttpRequestDataMapper interface {
	GetGenerateUrlCoreReq(reqPayload *httpModel.GenerateUrlReqPayload) coreModel.IGenerateUrlReq
	GetUrlListCoreReq(pageSize int, cursor int, sortingOrder string, filters pq.StringArray, orderBy string) coreModel.IGetList
	GetDecryptUrlCoreReq(shortCode string) coreModel.IGetUrlReq
}

type HttpRequestDataMapper struct{}

func (dm *HttpRequestDataMapper) GetGenerateUrlCoreReq(reqPayload *httpModel.GenerateUrlReqPayload) coreModel.IGenerateUrlReq {
	return &coreModel.GenerateUrlReq{
		LongUrl:     reqPayload.LongUrl,
		CustomUrl:   reqPayload.CustomAlias,
		IsCustomUrl: reqPayload.IsCustomUrl,
	}
}

func (dm *HttpRequestDataMapper) GetUrlListCoreReq(pageSize int, cursor int, sortingOrder string, filters pq.StringArray, orderBy string) coreModel.IGetList {
	return &coreModel.GetListRequestModel{
		Filters:      filters,
		PageSize:     pageSize,
		Cursor:       cursor,
		SortingOrder: sortingOrder,
		OrderBy:      orderBy,
	}
}

func (dm *HttpRequestDataMapper) GetDecryptUrlCoreReq(shortCode string) coreModel.IGetUrlReq {
	return &coreModel.GetUrlReq{ShortCode: shortCode}
}
