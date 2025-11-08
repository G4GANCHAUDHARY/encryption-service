package httpDataMapper

import (
	"github.com/G4GANCHAUDHARY/encryption-service/encryptor/models/httpModel"
)

type IHttpResponseDataMapper interface {
	GetGenerateUrlCoreRes(shortUrl string) *httpModel.GenerateUrlResPayload
}

type HttpResponseDataMapper struct{}

func (dm *HttpRequestDataMapper) GetGenerateUrlCoreRes(shortUrl string) *httpModel.GenerateUrlResPayload {
	return &httpModel.GenerateUrlResPayload{ShortUrl: shortUrl}
}
