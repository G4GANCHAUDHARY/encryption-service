package dbModel

import "time"

type Url struct {
	Id             int       `gorm:"id"`
	ShortCode      string    `gorm:"short_code"`
	LongUrl        string    `gorm:"long_url"`
	CreatedAt      time.Time `gorm:"created_at"`
	LastAccessedAt time.Time `gorm:"last_accessed_at"`
	ClickCount     int       `gorm:"click_count"`
	IsCustomUrl    bool      `gorm:"is_custom_url"`
	IsActive       bool      `gorm:"is_active"`
}
