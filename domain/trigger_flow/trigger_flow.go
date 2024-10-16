package triggerflow

import (
	"ifttt/manager/domain/resolvable"
	"ifttt/manager/domain/rule"
)

type Class struct {
	Name string `mapstructure:"name"`
}

type TriggerFlow struct {
	Name        string                `mapstructure:"name"`
	Description string                `mapstructure:"description"`
	Class       Class                 `mapstructure:"class"`
	StartRules  []rule.Rule           `mapstructure:"startRules"`
	AllRules    []rule.Rule           `mapstructure:"allRules"`
	BranchFlows map[uint][]BranchFlow `mapstructure:"branchFlows"`
}

type BranchFlow struct {
	IfReturn resolvable.Resolvable `json:"ifReturn" mapstructure:"ifReturn"`
	Jump     uint                  `json:"jump" mapstructure:"jump"`
}
