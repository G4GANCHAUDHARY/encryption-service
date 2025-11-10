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

// UrlAnalytics : Taking pgsql as time series db
type UrlAnalytics struct {
	ID          uint      `gorm:"primarykey"`
	CreatedAt   time.Time `gorm:"created_at"`
	UpdatedAt   time.Time `gorm:"updated_at"`
	Date        string    `gorm:"date"`
	TotalClicks int       `gorm:"total_clicks"`
}
