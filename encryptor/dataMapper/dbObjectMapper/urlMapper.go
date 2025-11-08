package dbObjectMapper

import "github.com/G4GANCHAUDHARY/encryption-service/encryptor/models/dbModel"

type IUrlMapper interface {
	GetUrlModel() *dbModel.Url
}

type UrlMapper struct{}

func (um *UrlMapper) GetUrlModel() *dbModel.Url {
	return &dbModel.Url{}
}
