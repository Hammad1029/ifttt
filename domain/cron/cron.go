package cron

import (
	"ifttt/manager/domain/resolvable"
	triggerflow "ifttt/manager/domain/trigger_flow"
)

type Cron struct {
	ID           uint                             `json:"id" mapstructure:"id"`
	Name         string                           `json:"name" mapstructure:"name"`
	Description  string                           `json:"description" mapstructure:"description"`
	Cron         string                           `json:"cron" mapstructure:"cron"`
	PreConfig    map[string]resolvable.Resolvable `json:"preConfig" mapstructure:"preConfig"`
	TriggerFlows *[]triggerflow.TriggerCondition  `json:"triggerFlows" mapstructure:"triggerFlows"`
}
