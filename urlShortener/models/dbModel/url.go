package dbModel

import (
	"gorm.io/gorm"
	"time"
)

type Url struct {
	gorm.Model
	ShortCode      string    `gorm:"short_code"`
	LongUrl        string    `gorm:"long_url"`
	LastAccessedAt time.Time `gorm:"last_accessed_at"`
	ClickCount     int       `gorm:"click_count"`
	IsCustomUrl    bool      `gorm:"is_custom_url"`
	IsActive       bool      `gorm:"is_active"`
}

// UrlAnalytics : Having pgsql as time series db
type UrlAnalytics struct {
	gorm.Model
	Date        string `gorm:"date"`
	TotalClicks int    `gorm:"total_clicks"`
}
