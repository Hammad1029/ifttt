package triggerflow

import (
	"ifttt/manager/domain/resolvable"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type CreateTriggerFlowRequest struct {
	Name        string               `json:"name" mapstructure:"name"`
	Description string               `json:"description" mapstructure:"description"`
	StartState  uint                 `json:"startState" mapstructure:"startState"`
	Rules       []string             `json:"rules" mapstructure:"rules"`
	BranchFlows map[uint]*BranchFlow `json:"branchFlows" mapstructure:"branchFlows"`
}

type GetDetailsRequest struct {
	Name string `json:"name" mapstructure:"name"`
}

type TriggerConditionRequest struct {
	If      resolvable.Condition `json:"if" mapstructure:"if"`
	Trigger string               `json:"trigger" mapstructure:"trigger"`
}

func (g *GetDetailsRequest) Validate() error {
	return validation.Validate(g.Name, validation.Required, validation.Length(3, 0))
}

func (c *CreateTriggerFlowRequest) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.Name, validation.Required, validation.Length(3, 0)),
		validation.Field(&c.Description, validation.Required, validation.Length(3, 0)),
		validation.Field(&c.StartState, validation.Required),
		validation.Field(&c.Rules, validation.Required),
		validation.Field(&c.BranchFlows, validation.Each(
			validation.By(func(value interface{}) error {
				b := value.(BranchFlow)
				return b.Validate()
			}))),
	)
}

func (b *BranchFlow) Validate() error {
	return validation.ValidateStruct(b,
		validation.Field(&b.Rule, validation.Required),
		validation.Field(&b.States),
	)
}

func (t *TriggerConditionRequest) Validate() error {
	return validation.ValidateStruct(t,
		validation.Field(&t.If, validation.By(func(value interface{}) error {
			c := value.(resolvable.Condition)
			return c.Validate()
		})),
		validation.Field(&t.Trigger, validation.Required),
	)
}
