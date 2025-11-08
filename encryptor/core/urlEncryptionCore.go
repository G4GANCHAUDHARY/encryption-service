package core

import (
	"context"
	"errors"
	"github.com/G4GANCHAUDHARY/encryption-service/encryptor/dataMapper/dbObjectMapper"
	"github.com/G4GANCHAUDHARY/encryption-service/encryptor/dataMapper/httpDataMapper"
	"github.com/G4GANCHAUDHARY/encryption-service/encryptor/domainModel"
	"github.com/G4GANCHAUDHARY/encryption-service/encryptor/models/coreModel"
	"github.com/G4GANCHAUDHARY/encryption-service/encryptor/models/dbModel"
	"github.com/G4GANCHAUDHARY/encryption-service/encryptor/models/httpModel"
	"github.com/G4GANCHAUDHARY/encryption-service/encryptor/repo"
	"github.com/G4GANCHAUDHARY/encryption-service/global"
	"github.com/G4GANCHAUDHARY/encryption-service/providers"
	"gorm.io/gorm"
)

type IUrlEncryptionCore interface {
	EncryptUrl(ctx context.Context, req coreModel.IGenerateUrlReq) (*httpModel.GenerateUrlResPayload, error)
}

type UrlEncryption struct {
	Db            *gorm.DB
	UrlRepository repo.IUrlRepository
	Redis         *providers.RedisLib
	HttpResMapper httpDataMapper.IHttpResponseDataMapper
	DbMapper      dbObjectMapper.IUrlMapper
}

func (ue *UrlEncryption) EncryptUrl(ctx context.Context, req coreModel.IGenerateUrlReq) (*httpModel.GenerateUrlResPayload, error) {
	// check if already exist, if yes return short code
	if existingUrl, err := ue.UrlRepository.Get(ctx, ue.Db, domainModel.GetUrlFilterString(req.GetLongUrl())); err != nil {
		return ue.HttpResMapper.GetGenerateUrlCoreRes(existingUrl.ShortCode), nil
	}

	// begin tx
	tx := ue.Db.Begin()

	// handle custom url
	if req.GetIsCustomUrl() {
		return ue.handleCustomUrl(ctx, tx, req)
	}

	// generate and save url with collision handling
	url, err := ue.GenerateUrlAndSave(ctx, tx, req)
	if err != nil {
		return nil, err
	}

	// commit tx
	if err = tx.Commit().Error; err != nil {
		return nil, err
	}

	return ue.HttpResMapper.GetGenerateUrlCoreRes(url.ShortCode), nil
}

func (ue *UrlEncryption) GenerateUrlAndSave(ctx context.Context, tx *gorm.DB, req coreModel.IGenerateUrlReq) (*dbModel.Url, error) {
	// max retry count
	for retry := 0; retry < global.RetryCount; retry++ {
		shortCode, err := ue.getShortCodeForUrl(ctx)
		if err != nil {
			return nil, err
		}

		urlEntity, err := ue.UrlRepository.Save(ctx, tx, ue.DbMapper.GetUrlModel(req, shortCode))
		if err == nil {
			// no collision found
			return urlEntity, nil
		}

		if domainModel.IsUniqueConstraintError(err) {
			// Collision: retry with new short code
			continue
		}
	}

	return nil, errors.New("collisions found while url generation")
}

func (ue *UrlEncryption) getShortCodeForUrl(ctx context.Context) (string, error) {
	// get unique key to create short code hash using redis-counter
	uniqKey, err := ue.Redis.Increment(ctx, global.Counter)
	if err != nil {
		return "", errors.New("redis-err while generating short url")
	}

	return domainModel.GetShortUrlFromUniqKey(uniqKey), nil
}

func (ue *UrlEncryption) handleCustomUrl(ctx context.Context, tx *gorm.DB, req coreModel.IGenerateUrlReq) (*httpModel.GenerateUrlResPayload, error) {
	if _, err := ue.UrlRepository.Save(ctx, tx, ue.DbMapper.GetUrlModel(req, "")); err != nil {
		return nil, err
	}
	return ue.HttpResMapper.GetGenerateUrlCoreRes(req.GetCustomUrl()), nil
}
