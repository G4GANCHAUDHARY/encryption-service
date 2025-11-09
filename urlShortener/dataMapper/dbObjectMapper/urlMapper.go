package dbObjectMapper

import (
	"github.com/G4GANCHAUDHARY/encryption-service/urlShortener/models/coreModel"
	"github.com/G4GANCHAUDHARY/encryption-service/urlShortener/models/dbModel"
	"time"
)

type IUrlMapper interface {
	GetUrlModel(req coreModel.IGenerateUrlReq, shortCode string) *dbModel.Url
	GetUpdateUrlModel(urlEntity *dbModel.Url) *dbModel.Url
}

type UrlMapper struct{}

func (um *UrlMapper) GetUrlModel(req coreModel.IGenerateUrlReq, shortCode string) *dbModel.Url {
	urlEntity := &dbModel.Url{
		LongUrl:     req.GetLongUrl(),
		IsCustomUrl: req.GetIsCustomUrl(),
		IsActive:    true,
	}
	if req.GetIsCustomUrl() {
		urlEntity.ShortCode = req.GetCustomUrl()
	} else {
		urlEntity.ShortCode = shortCode
	}
	return urlEntity
}

func (um *UrlMapper) GetUpdateUrlModel(urlEntity *dbModel.Url) *dbModel.Url {
	urlEntity.ClickCount++
	urlEntity.LastAccessedAt = time.Now()
	return urlEntity
}
