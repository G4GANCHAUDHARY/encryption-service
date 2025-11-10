package httpModel

import validation "github.com/go-ozzo/ozzo-validation/v4"

type GenerateUrlReqPayload struct {
	LongUrl     string `json:"long_url"`
	CustomAlias string `json:"custom_alias"`
	IsCustomUrl bool   `json:"is_custom_url"`
}

func (p *GenerateUrlReqPayload) Validate() error {
	return validation.ValidateStruct(p)
}
