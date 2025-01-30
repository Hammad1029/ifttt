package api

import (
	requestvalidator "ifttt/manager/domain/request_validator"
	triggerflow "ifttt/manager/domain/trigger_flow"
)

type Api struct {
	ID          uint                                         `json:"id" mapstructure:"id"`
	Name        string                                       `json:"name" mapstructure:"name"`
	Path        string                                       `json:"path" mapstructure:"path"`
	Method      string                                       `json:"method" mapstructure:"method"`
	Description string                                       `json:"description" mapstructure:"description"`
	Request     map[string]requestvalidator.RequestParameter `json:"request" mapstructure:"request"`
	Response    map[uint]ResponseDefinition                  `json:"response" mapstructure:"response"`
	Triggers    *[]triggerflow.TriggerCondition              `json:"triggers" mapstructure:"triggers"`
}

type ResponseDefinition struct {
	UseProfile     string         `json:"useProfile" mapstructure:"useProfile"`
	Definition     map[string]any `json:"definition" mapstructure:"Definition"`
	HTTPStatusCode int            `json:"httpStatusCode" mapstructure:"httpStatusCode"`
}
