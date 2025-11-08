package core

import (
	"context"
	"github.com/G4GANCHAUDHARY/encryption-service/encryptor/models/coreModel"
	"github.com/G4GANCHAUDHARY/encryption-service/encryptor/models/httpModel"
	"github.com/G4GANCHAUDHARY/encryption-service/encryptor/repo"
	"github.com/G4GANCHAUDHARY/encryption-service/providers"
	"gorm.io/gorm"
)

type IUrlEncryptionCore interface {
	EncryptUrl(ctx context.Context, req coreModel.IGenerateUrlReq) *httpModel.GenerateUrlResPayload
}

type UrlEncryption struct {
	Db            *gorm.DB
	UrlRepository repo.IUrlRepository
	Redis         *providers.RedisLib
}

func (ue *UrlEncryption) EncryptUrl(ctx context.Context, req coreModel.IGenerateUrlReq) *httpModel.GenerateUrlResPayload {
	// check if already exist, if yes return short code

	// check if custom aliases, if yes save mapping without unique handling

	// get uniq counter from redis

	// generate short url using some algo

	// save url mapping

	// return
	return nil
}
