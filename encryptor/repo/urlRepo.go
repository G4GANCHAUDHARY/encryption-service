package repo

import (
	"context"
	"github.com/G4GANCHAUDHARY/encryption-service/encryptor/models/dbModel"
	"gorm.io/gorm"
)

type IUrlRepository interface {
	Save(ctx context.Context, tx *gorm.DB, url *dbModel.Url) (*dbModel.Url, error)
	Get(ctx context.Context, tx *gorm.DB, filters string) (*dbModel.Url, error)
}

type UrlRepository struct {
	//ToDo logger and tracer DI
}

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
