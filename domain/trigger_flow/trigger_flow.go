package triggerflow

import (
	"ifttt/manager/domain/condition"
	"ifttt/manager/domain/rule"
)

type Class struct {
	Name string `json:"name" mapstructure:"name"`
}

type TriggerFlow struct {
	ID          uint                 `json:"id" mapstructure:"id"`
	Name        string               `json:"name" mapstructure:"name"`
	Description string               `json:"description" mapstructure:"description"`
	Class       Class                `json:"class" mapstructure:"class"`
	StartState  uint                 `json:"startState" mapstructure:"startState"`
	Rules       map[uint]*rule.Rule  `json:"rules" mapstructure:"rules"`
	BranchFlows map[uint]*BranchFlow `json:"branchFlows" mapstructure:"branchFlows"`
}

type BranchFlow struct {
	Rule   uint          `json:"rule" mapstructure:"rule"`
	States map[uint]uint `json:"states" mapstructure:"states"`
}

type TriggerCondition struct {
	If      condition.Condition `json:"if" mapstructure:"if"`
	Trigger TriggerFlow         `json:"trigger" mapstructure:"trigger"`
}
