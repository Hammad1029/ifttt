package triggerflow

import validation "github.com/go-ozzo/ozzo-validation/v4"

type GetDetailsRequest struct {
	Name string `json:"name" mapstructure:"name"`
}

func (g *GetDetailsRequest) Validate() error {
	return validation.Validate(g.Name, validation.Required, validation.Length(3, 0))
}

type CreateTriggerFlowRequest struct {
	Name        string `json:"name" mapstructure:"name"`
	Description string `json:"description" mapstructure:"description"`
	Class       uint   `json:"class" mapstructure:"class"`
	StartRules  []uint `json:"startRules" mapstructure:"startRules"`
	AllRules    []uint `json:"allRules" mapstructure:"allRules"`
}

func (c *CreateTriggerFlowRequest) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.Name, validation.Required, validation.Length(3, 0)),
		validation.Field(&c.Description, validation.Required, validation.Length(3, 0)),
		validation.Field(&c.Class, validation.Required),
		validation.Field(&c.StartRules, validation.Required, validation.Each(validation.In(c.AllRules))),
		validation.Field(&c.AllRules, validation.Required),
	)
}
