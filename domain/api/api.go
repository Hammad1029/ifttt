package api

import (
	requestvalidator "ifttt/manager/domain/request_validator"
	"ifttt/manager/domain/resolvable"
	triggerflow "ifttt/manager/domain/trigger_flow"
)

type Api struct {
	ID          uint                                         `json:"id" mapstructure:"id"`
	Name        string                                       `json:"name" mapstructure:"name"`
	Path        string                                       `json:"path" mapstructure:"path"`
	Method      string                                       `json:"method" mapstructure:"method"`
	Description string                                       `json:"description" mapstructure:"description"`
	Request     map[string]requestvalidator.RequestParameter `json:"request" mapstructure:"request"`
	PreConfig   map[string]resolvable.Resolvable             `json:"preConfig" mapstructure:"preConfig"`
	Triggers    *[]triggerflow.TriggerCondition              `json:"triggers" mapstructure:"triggers"`
}
