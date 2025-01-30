package configuration

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (p *ResponseProfile) Validate() error {
	return validation.ValidateStruct(p,
		validation.Field(&p.Name, validation.Required),
		validation.Field(&p.ResponseHTTPStatus, validation.Required),
		validation.Field(&p.BodyFormat),
	)
}

func (g *InternalTagGroup) Validate() error {
	return validation.ValidateStruct(g,
		validation.Field(&g.Name, validation.Required),
		validation.Field(&g.Tags, validation.Nil),
	)
}

type InternalTagRequest struct {
	Name     string `json:"name" mapstructure:"name"`
	Groups   []uint `json:"groups" mapstructure:"groups"`
	All      bool   `json:"all" mapstructure:"all"`
	Reserved bool   `json:"reserved" mapstructure:"reserved"`
}

func (p *InternalTagRequest) Validate() error {
	return validation.ValidateStruct(p,
		validation.Field(&p.Name, validation.Required),
		validation.Field(&p.Groups, validation.Required),
		validation.Field(&p.All),
		validation.Field(&p.Reserved),
	)
}
