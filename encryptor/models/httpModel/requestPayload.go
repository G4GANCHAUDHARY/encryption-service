package httpModel

type GenerateUrlReqPayload struct {
	LongUrl     string `json:"long_url"`
	CustomAlias string `json:"custom_alias"`
	IsCustomUrl bool   `json:"is_custom_url"`
}

// add validations here
