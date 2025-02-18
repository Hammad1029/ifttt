package cron

import (
	"ifttt/manager/common"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type GetCronRequest struct {
	Name string `json:"name" mapstructure:"name"`
}

func (c *GetCronRequest) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.Name, validation.Required),
	)
}

func (c *Cron) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.Name, validation.Required, validation.Length(3, 0)),
		validation.Field(&c.Description, validation.Required, validation.Length(3, 0)),
		validation.Field(&c.CronExpr, validation.Required, validation.Match(common.RegexCron)),
		validation.Field(&c.ApiName, validation.Required),
		validation.Field(&c.API, validation.Nil),
	)
}
