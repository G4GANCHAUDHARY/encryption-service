package dbObjectMapper

import (
	"github.com/G4GANCHAUDHARY/encryption-service/urlShortener/models/dbModel"
	"time"
)

type IUrlAnalyticsMapper interface {
	GetUrlAnalyticsModel() *dbModel.UrlAnalytics
}

type UrlAnalyticsMapper struct{}

func (um *UrlAnalyticsMapper) GetUrlAnalyticsModel() *dbModel.UrlAnalytics {
	return &dbModel.UrlAnalytics{
		Date: time.Now().Format("2006-01-02"),
	}
}
