package httpDataMapper

import (
	"github.com/G4GANCHAUDHARY/encryption-service/encryptor/models/coreModel"
	"github.com/G4GANCHAUDHARY/encryption-service/encryptor/models/httpModel"
)

type IHttpRequestDataMapper interface {
	GetGenerateUrlCoreReq(reqPayload httpModel.GenerateUrlReqPayload) coreModel.IGenerateUrlReq
}

type HttpRequestDataMapper struct{}

func (dm *HttpRequestDataMapper) GetGenerateUrlCoreReq(reqPayload httpModel.GenerateUrlReqPayload) coreModel.IGenerateUrlReq {
	return &coreModel.GenerateUrlReq{
		LongUrl:     reqPayload.LongUrl,
		CustomUrl:   reqPayload.CustomAlias,
		IsCustomUrl: reqPayload.IsCustomUrl,
	}
}
