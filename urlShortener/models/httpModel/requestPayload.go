package httpModel

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"strings"
)

type GenerateUrlReqPayload struct {
	LongUrl     string `json:"long_url"`
	CustomAlias string `json:"custom_alias"`
	IsCustomUrl bool   `json:"is_custom_url"`
}

func (p *GenerateUrlReqPayload) Validate() error {
	p.LongUrl = strings.TrimSpace(p.LongUrl)
	p.CustomAlias = strings.TrimSpace(p.CustomAlias)
	return validation.ValidateStruct(p,
		validation.Field(&p.LongUrl, validation.Required.Error("long url is required"), is.URL.Error("enter valid long url")),
		validation.Field(&p.CustomAlias,
			validation.By(func(value interface{}) error {
				
				alias := value.(string)
				if p.IsCustomUrl {
					if alias == "" {
						return validation.NewError("validation_custom_alias", "custom_alias is required when is_custom_url is true")
					}
				}

				if alias != "" {
					if len(alias) < 3 || len(alias) > 10 {
						return validation.NewError("validation_custom_alias_length", "custom_alias must be between 3 and 10 characters")
					}
				}
				return nil
			}),
		),
	)
}
