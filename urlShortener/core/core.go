package core

import (
	"context"
	"errors"
	"github.com/G4GANCHAUDHARY/encryption-service/global"
	"github.com/G4GANCHAUDHARY/encryption-service/providers"
	"github.com/G4GANCHAUDHARY/encryption-service/urlShortener/dataMapper/dbObjectMapper"
	"github.com/G4GANCHAUDHARY/encryption-service/urlShortener/dataMapper/httpDataMapper"
	"github.com/G4GANCHAUDHARY/encryption-service/urlShortener/domainModel"
	"github.com/G4GANCHAUDHARY/encryption-service/urlShortener/models/coreModel"
	"github.com/G4GANCHAUDHARY/encryption-service/urlShortener/models/dbModel"
	"github.com/G4GANCHAUDHARY/encryption-service/urlShortener/models/httpModel"
	"github.com/G4GANCHAUDHARY/encryption-service/urlShortener/repo"
	"github.com/G4GANCHAUDHARY/encryption-service/utils"
	"gorm.io/gorm"
)

type IUrlShortenerCore interface {
	EncryptUrl(ctx context.Context, req coreModel.IGenerateUrlReq) (*httpModel.GenerateUrlResPayload, error)
	DecryptUrl(ctx context.Context, req coreModel.IGetUrlReq) (*httpModel.GetUrlResPayload, error)
	ExpireUrls(ctx context.Context, shortCode string) error
	GetUrls(ctx context.Context, req coreModel.IGetList) (*httpModel.GetUrlListPayload, error)
}

type UrlShortener struct {
	Db               *gorm.DB
	UrlRepository    repo.IUrlRepository
	RedisCache       *providers.RedisLib
	RedisCounter     *providers.RedisLib
	HttpResMapper    httpDataMapper.IHttpResponseDataMapper
	DbMapper         dbObjectMapper.IUrlMapper
	UrlAnalyticsRepo repo.IUrlAnalyticsRepository
}

// EncryptUrl : encrypts long url to short code
func (ue *UrlShortener) EncryptUrl(ctx context.Context, req coreModel.IGenerateUrlReq) (res *httpModel.GenerateUrlResPayload, err error) {
	// check if already exist, if yes return short code
	existingUrl, err := ue.UrlRepository.Get(ctx, ue.Db, domainModel.GetLongUrlFilterString(req.GetLongUrl()))
	if err == nil {
		return ue.HttpResMapper.GetGenerateUrlCoreRes(existingUrl.ShortCode), nil
	}

	// begin tx
	tx := ue.Db.Begin()
	defer utils.HandlePanic(utils.HandlePanicRequest{
		Tx:       tx,
		FuncName: "EncryptUrl",
		Err:      &err,
	})

	// handle custom url
	if req.GetIsCustomUrl() {
		return ue.handleCustomUrl(ctx, tx, req)
	}

	// generate and save url with collision handling
	url, err := ue.generateUrlAndSave(ctx, tx, req)
	if err != nil {
		return nil, err
	}

	// commit tx
	if err = tx.Commit().Error; err != nil {
		return nil, err
	}

	return ue.HttpResMapper.GetGenerateUrlCoreRes(url.ShortCode), nil
}

// DecryptUrl : decrypts short code to long url
func (ue *UrlShortener) DecryptUrl(ctx context.Context, req coreModel.IGetUrlReq) (res *httpModel.GetUrlResPayload, err error) {
	// defer saving url mapping analytics data
	var urlEntity *dbModel.Url
	defer func() {
		if err == nil && urlEntity != nil {
			go func(entity *dbModel.Url) {
				ue.saveUrlMappingAnalytics(context.Background(), req, entity)
			}(urlEntity)
		}
	}()

	// check if exist in redis -> if yes return
	var longUrl string
	longUrl, err = ue.RedisCache.Get(ctx, global.Url+req.GetShortCode())
	if err == nil {
		return ue.HttpResMapper.GetUrlCoreRes(longUrl), nil
	}

	// check if exist in db -> if yes return
	urlEntity, err = ue.UrlRepository.Get(ctx, ue.Db, domainModel.GetShortUrlFilterString(req.GetShortCode()))
	if err == nil {
		return ue.HttpResMapper.GetUrlCoreRes(urlEntity.LongUrl), nil
	}

	return nil, errors.New("redirect url not found")
}

// ExpireUrls : cron would call this func for each expired url
func (ue *UrlShortener) ExpireUrls(ctx context.Context, shortCode string) error {
	urlEntity, err := ue.UrlRepository.Get(ctx, ue.Db, domainModel.GetShortUrlFilterString(shortCode))
	if err != nil {
		return err
	}

	tx := ue.Db.Begin()
	defer utils.HandlePanic(utils.HandlePanicRequest{
		Tx:       tx,
		FuncName: "ExpireUrls",
		Err:      &err,
	})

	// set is active false
	urlEntity.IsActive = false
	urlEntity, err = ue.UrlRepository.Save(ctx, tx, urlEntity)
	if err != nil {
		return err
	}

	if err = tx.Commit().Error; err != nil {
		return nil
	}

	// delete from redis
	if _, err = ue.RedisCache.Delete(ctx, global.Url+shortCode); err != nil {
		return err
	}

	return nil
}

// GetUrls : returns list of active url mappings to admin
func (ue *UrlShortener) GetUrls(ctx context.Context, req coreModel.IGetList) (*httpModel.GetUrlListPayload, error) {

	urlList, err := ue.UrlRepository.GetList(ctx, ue.Db, req.GetFilters(), &utils.SortObject{SortBy: req.GetOrderBy(), SortOrder: req.GetSortingOrder()}, &utils.Pagination{OffSet: req.GetCursor(), Limit: req.GetPageSize()})
	if err != nil {
		return nil, err
	}

	return ue.HttpResMapper.GetUrlListRes(urlList), nil
}

func (ue *UrlShortener) generateUrlAndSave(ctx context.Context, tx *gorm.DB, req coreModel.IGenerateUrlReq) (*dbModel.Url, error) {
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

func (ue *UrlShortener) getShortCodeForUrl(ctx context.Context) (string, error) {
	// get unique key to create short code hash using redis-counter
	uniqKey, err := ue.RedisCounter.Increment(ctx, global.Counter)
	if err != nil {
		return "", errors.New("redis-err while generating short url")
	}

	return domainModel.GetShortUrlFromUniqKey(uniqKey), nil
}

func (ue *UrlShortener) handleCustomUrl(ctx context.Context, tx *gorm.DB, req coreModel.IGenerateUrlReq) (*httpModel.GenerateUrlResPayload, error) {
	// save url mapping with short code as custom url itself
	if _, err := ue.UrlRepository.Save(ctx, tx, ue.DbMapper.GetUrlModel(req, req.GetCustomUrl())); err != nil {
		return nil, err
	}

	// commit tx
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return ue.HttpResMapper.GetGenerateUrlCoreRes(req.GetCustomUrl()), nil
}

func (ue *UrlShortener) saveUrlMappingAnalytics(ctx context.Context, req coreModel.IGetUrlReq, urlEntity *dbModel.Url) {
	var err error
	tx := ue.Db.Begin()
	defer utils.HandlePanic(utils.HandlePanicRequest{
		Tx:       tx,
		FuncName: "saveUrlMappingAnalytics",
		Err:      &err,
	})

	// update url entity with click count and last accessed at
	updatedUrlEntity := ue.DbMapper.GetUpdateUrlModel(urlEntity)

	// save url entity
	updatedUrlEntity, err = ue.UrlRepository.Save(ctx, tx, updatedUrlEntity)
	if err != nil {
		return
	}

	// save url analytics
	err = ue.UrlAnalyticsRepo.SaveAnalytics(ctx, tx)
	if err != nil {
		return
	}

	if err = tx.Commit().Error; err != nil {
		return
	}

	// save in cache
	_ = ue.RedisCache.Set(ctx, global.Url+req.GetShortCode(), urlEntity.LongUrl)
}
