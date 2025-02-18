package triggerflow

import (
	"ifttt/manager/domain/resolvable"
	"ifttt/manager/domain/rule"
)

type TriggerFlow struct {
	ID          uint                 `json:"id" mapstructure:"id"`
	Name        string               `json:"name" mapstructure:"name"`
	Description string               `json:"description" mapstructure:"description"`
	StartState  uint                 `json:"startState" mapstructure:"startState"`
	Rules       map[uint]*rule.Rule  `json:"rules" mapstructure:"rules"`
	BranchFlows map[uint]*BranchFlow `json:"branchFlows" mapstructure:"branchFlows"`
}

type BranchFlow struct {
	Rule   string        `json:"rule" mapstructure:"rule"`
	States map[uint]uint `json:"states" mapstructure:"states"`
}

type TriggerCondition struct {
	If      resolvable.Condition `json:"if" mapstructure:"if"`
	Trigger TriggerFlow          `json:"trigger" mapstructure:"trigger"`
}
