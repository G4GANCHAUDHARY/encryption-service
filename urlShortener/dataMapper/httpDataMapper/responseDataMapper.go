package httpDataMapper

import (
	"github.com/G4GANCHAUDHARY/encryption-service/urlShortener/models/dbModel"
	"github.com/G4GANCHAUDHARY/encryption-service/urlShortener/models/httpModel"
)

type IHttpResponseDataMapper interface {
	GetGenerateUrlCoreRes(shortUrl string) *httpModel.GenerateUrlResPayload
	GetUrlCoreRes(longUrl string) *httpModel.GetUrlResPayload
	GetUrlListRes(urls *[]dbModel.Url) *httpModel.GetUrlListPayload
}

type HttpResponseDataMapper struct{}

func (dm *HttpResponseDataMapper) GetGenerateUrlCoreRes(shortUrl string) *httpModel.GenerateUrlResPayload {
	return &httpModel.GenerateUrlResPayload{ShortUrl: shortUrl}
}

func (dm *HttpResponseDataMapper) GetUrlCoreRes(longUrl string) *httpModel.GetUrlResPayload {
	return &httpModel.GetUrlResPayload{LongUrl: longUrl}
}

func (dm *HttpResponseDataMapper) GetUrlListRes(urls *[]dbModel.Url) *httpModel.GetUrlListPayload {
	res := httpModel.GetUrlListPayload{}
	for _, url := range *urls {
		urlListObj := httpModel.Url{
			Id:             int(url.ID),
			ShortCode:      url.ShortCode,
			LongUrl:        url.LongUrl,
			LastAccessedAt: url.LastAccessedAt,
			ClickCount:     url.ClickCount,
			IsCustomUrl:    url.IsCustomUrl,
		}
		res.UrlList = append(res.UrlList, urlListObj)
	}
	return &res
}
