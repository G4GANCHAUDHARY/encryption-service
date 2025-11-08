package coreModel

type IGenerateUrlReq interface {
	GetLongUrl() string
	GetIsCustomUrl() bool
	GetCustomUrl() string
}
