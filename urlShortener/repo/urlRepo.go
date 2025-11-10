package repo

import (
	"context"
	"github.com/G4GANCHAUDHARY/encryption-service/urlShortener/models/dbModel"
	"github.com/G4GANCHAUDHARY/encryption-service/utils"
	"gorm.io/gorm"
)

type IUrlRepository interface {
	Save(ctx context.Context, tx *gorm.DB, url *dbModel.Url) (*dbModel.Url, error)
	Get(ctx context.Context, tx *gorm.DB, filters string) (*dbModel.Url, error)
	GetList(ctx context.Context, tx *gorm.DB, filtersString []string, sortObject utils.ISortObject, pagination utils.IPagination) (*[]dbModel.Url, error)
}

type UrlRepository struct{}

func (u *UrlRepository) Save(ctx context.Context, tx *gorm.DB, url *dbModel.Url) (*dbModel.Url, error) {
	if err := tx.Save(&url).Error; err != nil {
		return nil, err
	}
	return url, nil
}

func (u *UrlRepository) Get(ctx context.Context, tx *gorm.DB, filters string) (*dbModel.Url, error) {
	var urlEntity dbModel.Url
	if err := tx.Where(filters).Order("created_at DESC").First(&urlEntity, "is_active = ?", true).Error; err != nil {
		return nil, err
	}
	return &urlEntity, nil
}

func (u *UrlRepository) GetList(ctx context.Context, tx *gorm.DB, filtersString []string, sortObject utils.ISortObject, pagination utils.IPagination) (*[]dbModel.Url, error) {
	var urls []dbModel.Url
	tx = tx.Table("url_shortener.url")
	tx = tx.Select("url.id as id, url.short_code as short_code, url.long_url as long_url, url.click_count as click_count, url.last_accessed_at as last_accessed_at, url.is_custom_url as is_custom_url")
	for _, filter := range filtersString {
		tx = tx.Where(filter)
	}
	tx.Limit(pagination.GetLimit()).Offset(pagination.GetOffSet()).Order(sortObject.GetQueryString()).Find(&urls)
	if tx.Error != nil {
		return &urls, tx.Error
	}

	return &urls, nil
}
