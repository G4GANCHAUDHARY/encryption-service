package httpDataMapper

import (
	"github.com/G4GANCHAUDHARY/encryption-service/encryptor/models/httpModel"
)

type IHttpResponseDataMapper interface {
	GetGenerateUrlCoreRes() *httpModel.GenerateUrlResPayload
}

type HttpResponseDataMapper struct{}

func (dm *HttpRequestDataMapper) GetGenerateUrlCoreRes() *httpModel.GenerateUrlResPayload {
	return &httpModel.GenerateUrlResPayload{}
}
