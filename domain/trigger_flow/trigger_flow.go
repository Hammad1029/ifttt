package triggerflow

import (
	"ifttt/manager/domain/rule"
)

type Class struct {
	Name string `mapstructure:"name"`
}

type TriggerFlow struct {
	Name        string      `mapstructure:"name"`
	Description string      `mapstructure:"description"`
	Class       Class       `mapstructure:"class"`
	StartRules  []rule.Rule `mapstructure:"startRules"`
	AllRules    []rule.Rule `mapstructure:"allRules"`
}
