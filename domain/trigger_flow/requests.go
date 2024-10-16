package triggerflow

import (
	"ifttt/manager/domain/resolvable"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type GetDetailsRequest struct {
	Name string `json:"name" mapstructure:"name"`
}

func (g *GetDetailsRequest) Validate() error {
	return validation.Validate(g.Name, validation.Required, validation.Length(3, 0))
}

type CreateTriggerFlowRequest struct {
	Name        string                `json:"name" mapstructure:"name"`
	Description string                `json:"description" mapstructure:"description"`
	Class       uint                  `json:"class" mapstructure:"class"`
	StartRules  []uint                `json:"startRules" mapstructure:"startRules"`
	AllRules    []uint                `json:"allRules" mapstructure:"allRules"`
	BranchFlows map[uint][]BranchFlow `json:"branchFlows" mapstructure:"branchFlows"`
}

func (c *CreateTriggerFlowRequest) Validate() error {
	allRules := make([]any, len(c.AllRules))
	for _, v := range c.AllRules {
		allRules = append(allRules, v)
	}

	return validation.ValidateStruct(c,
		validation.Field(&c.Name, validation.Required, validation.Length(3, 0)),
		validation.Field(&c.Description, validation.Required, validation.Length(3, 0)),
		validation.Field(&c.Class, validation.Required),
		validation.Field(&c.StartRules, validation.Required, validation.Each(validation.In(allRules...))),
		validation.Field(&c.AllRules, validation.Required),
		validation.Field(&c.BranchFlows, validation.Each(validation.Each(
			validation.By(func(value interface{}) error {
				b := value.(BranchFlow)
				return b.Validate()
			})))),
	)
}

func (b *BranchFlow) Validate() error {
	return validation.ValidateStruct(b,
		validation.Field(&b.IfReturn, validation.By(func(value interface{}) error {
			r := value.(resolvable.Resolvable)
			return r.Validate()
		})),
		validation.Field(&b.Jump, validation.Required),
	)
}
