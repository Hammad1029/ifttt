package api

import (
	"ifttt/manager/domain/resolvable"
	triggerflow "ifttt/manager/domain/trigger_flow"
)

type Api struct {
	ID          uint                             `json:"id" mapstructure:"id"`
	Name        string                           `json:"name" mapstructure:"name"`
	Path        string                           `json:"path" mapstructure:"path"`
	Method      string                           `json:"method" mapstructure:"method"`
	Description string                           `json:"description" mapstructure:"description"`
	Request     map[string]any                   `json:"request" mapstructure:"request"`
	PreConfig   map[string]resolvable.Resolvable `json:"preConfig" mapstructure:"preConfig"`
	PreWare     *[]triggerflow.TriggerFlow       `json:"preWare" mapstructure:"preWare"`
	MainWare    *[]triggerflow.TriggerCondition  `json:"mainWare" mapstructure:"mainWare"`
	PostWare    *[]triggerflow.TriggerFlow       `json:"postWare" mapstructure:"postWare"`
}
