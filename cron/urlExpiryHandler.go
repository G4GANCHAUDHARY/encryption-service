package cron

import (
	"context"
	"github.com/G4GANCHAUDHARY/encryption-service/urlShortener/core"
	"github.com/G4GANCHAUDHARY/encryption-service/urlShortener/models/dbModel"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
	"log"
	"time"
)

type UrlExpiryCron struct {
	Core core.IUrlShortenerCore
	Db   *gorm.DB
	cron *cron.Cron
}

func (c *UrlExpiryCron) ExpireUrl() {
	var err error
	ctx := context.Background()
	ctx = context.WithValue(ctx, "source", "cron")

	// five-year old data
	fiveYearsAgo := time.Now().UTC().AddDate(-5, 0, 0)

	var oldActiveURLs []dbModel.Url
	err = c.Db.Where("is_active = ? AND created_at < ?", true, fiveYearsAgo).Find(&oldActiveURLs).Error
	if err != nil {
		log.Println("err while fetching old url records")
	}

	for _, oldUrl := range oldActiveURLs {
		err = c.Core.ExpireUrls(ctx, oldUrl.ShortCode)
		if err != nil {
			log.Printf("err while expiring old url %v", oldUrl.ShortCode)
		}
	}
}

func (c *UrlExpiryCron) StartDailyCron() {
	c.cron = cron.New(cron.WithLocation(time.UTC))

	// every 3AM :}
	_, err := c.cron.AddFunc("0 3 * * *", func() {
		c.ExpireUrl()
	})
	if err != nil {
		log.Printf("failed to schedule cron: %v", err)
	}

	c.cron.Start()
}

func (c *UrlExpiryCron) Stop() {
	if c.cron != nil {
		c.cron.Stop()
	}
}
