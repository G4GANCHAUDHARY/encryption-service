package repo

import (
	"context"
	"github.com/G4GANCHAUDHARY/encryption-service/urlShortener/models/dbModel"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type IUrlAnalyticsRepository interface {
	SaveAnalytics(ctx context.Context, tx *gorm.DB) error
}

type UrlAnalyticsRepository struct{}

func (u *UrlAnalyticsRepository) SaveAnalytics(ctx context.Context, tx *gorm.DB) error {
	// upsert ops
	return tx.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "date"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"total_clicks": gorm.Expr("url_analytics.total_clicks + EXCLUDED.total_clicks"),
		}),
	}).Create(&dbModel.UrlAnalytics{
		Date:        time.Now().Format("2006-01-02"),
		TotalClicks: 1,
	}).Error
}
