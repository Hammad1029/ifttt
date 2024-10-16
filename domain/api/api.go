package api

import (
	"ifttt/manager/domain/condition"
	"ifttt/manager/domain/resolvable"
	triggerflow "ifttt/manager/domain/trigger_flow"
)

type Api struct {
	Name         string                           `json:"name" mapstructure:"name"`
	Path         string                           `json:"path" mapstructure:"path"`
	Method       string                           `json:"method" mapstructure:"method"`
	Description  string                           `json:"description" mapstructure:"description"`
	Request      map[string]any                   `json:"request" mapstructure:"request"`
	PreConfig    map[string]resolvable.Resolvable `json:"preConfig" mapstructure:"preConfig"`
	TriggerFlows *[]TriggerCondition              `json:"triggerFlows" mapstructure:"triggerFlows"`
}

type TriggerCondition struct {
	If      condition.Condition     `json:"if" mapstructure:"if"`
	Trigger triggerflow.TriggerFlow `json:"trigger" mapstructure:"trigger"`
}
