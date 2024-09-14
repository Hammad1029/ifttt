package triggerflow

import (
	"ifttt/manager/domain/rule"
)

type classes struct {
	Name string `mapstructure:"name"`
}

type TriggerFlow struct {
	Name        string      `mapstructure:"name"`
	Description string      `mapstructure:"description"`
	Class       classes     `mapstructure:"class"`
	StartRules  []rule.Rule `mapstructure:"startRules"`
	AllRules    []rule.Rule `mapstructure:"allRules"`
}
