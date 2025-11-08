package coreModel

type GenerateUrlReq struct {
	LongUrl     string
	IsCustomUrl bool
	CustomUrl   string
}

func (r *GenerateUrlReq) GetLongUrl() string {
	return r.LongUrl
}

func (r *GenerateUrlReq) GetIsCustomUrl() bool {
	return r.IsCustomUrl
}

func (r *GenerateUrlReq) GetCustomUrl() string {
	return r.CustomUrl
}

func (r *GenerateUrlReq) SetLongUrl(url string) {
	r.LongUrl = url
}

func (r *GenerateUrlReq) SetIsCustomUrl(isCustom bool) {
	r.IsCustomUrl = isCustom
}

func (r *GenerateUrlReq) SetCustomUrl(url string) {
	r.CustomUrl = url
}
