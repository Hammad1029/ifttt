package cron

import (
	"ifttt/manager/common"
	"ifttt/manager/domain/resolvable"
	triggerflow "ifttt/manager/domain/trigger_flow"

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

type CreateCronRequest struct {
	Name        string                                `json:"name" mapstructure:"name"`
	Description string                                `json:"description" mapstructure:"description"`
	Cron        string                                `json:"cron" mapstructure:"cron"`
	PreConfig   map[string]resolvable.Resolvable      `json:"preConfig" mapstructure:"preConfig"`
	Triggers    []triggerflow.TriggerConditionRequest `json:"triggers" mapstructure:"triggers"`
}

func (c *CreateCronRequest) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.Name, validation.Required, validation.Length(3, 0)),
		validation.Field(&c.Description, validation.Required, validation.Length(3, 0)),
		validation.Field(&c.Cron, validation.Required, validation.Match(common.RegexCron)),
		validation.Field(&c.PreConfig, validation.By(
			func(value interface{}) error {
				rMap := value.(map[string]resolvable.Resolvable)
				for _, r := range rMap {
					if err := r.Validate(); err != nil {
						return err
					}
				}
				return nil
			})),
		validation.Field(&c.Triggers, validation.Required, validation.Each(
			validation.By(func(value interface{}) error {
				t := value.(triggerflow.TriggerConditionRequest)
				return t.Validate()
			}))),
	)
}
