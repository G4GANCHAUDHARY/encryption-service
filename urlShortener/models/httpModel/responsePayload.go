package httpModel

import (
	"time"
)

type GenerateUrlResPayload struct {
	ShortUrl string `json:"short_url"`
}

type GetUrlResPayload struct {
	LongUrl string `json:"long_url"`
}

type GetUrlListPayload struct {
	UrlList []Url `json:"url_list"`
}

type Url struct {
	Id             int       `json:"id"`
	ShortCode      string    `json:"short_code"`
	LongUrl        string    `json:"long_url"`
	LastAccessedAt time.Time `json:"last_accessed_at"`
	ClickCount     int       `json:"click_count"`
	IsCustomUrl    bool      `json:"is_custom_url"`
}
