package api

import (
	"fmt"
	"ifttt/manager/common"
	requestvalidator "ifttt/manager/domain/request_validator"
	"ifttt/manager/domain/resolvable"
	triggerflow "ifttt/manager/domain/trigger_flow"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type GetDetailsRequest struct {
	Name string `json:"name" mapstructure:"name"`
	Path string `json:"path" mapstructure:"path"`
}

func (g *GetDetailsRequest) Validate() error {
	return validation.ValidateStruct(g,
		validation.Field(&g.Name, validation.Required, validation.Length(3, 0)),
		validation.Field(&g.Path, validation.Required, validation.Length(3, 0),
			validation.Match(common.RegexEndpoint)),
	)
}

type CreateApiRequest struct {
	Name        string                                       `json:"name" mapstructure:"name"`
	Path        string                                       `json:"path" mapstructure:"path"`
	Method      string                                       `json:"method" mapstructure:"method"`
	Description string                                       `json:"description" mapstructure:"description"`
	Request     map[string]requestvalidator.RequestParameter `json:"request" mapstructure:"request"`
	Response    map[uint]ResponseDefinition                  `json:"response" mapstructure:"response"`
	PreConfig   map[string]resolvable.Resolvable             `json:"preConfig" mapstructure:"preConfig"`
	Triggers    []triggerflow.TriggerConditionRequest        `json:"triggers" mapstructure:"triggers"`
}

func (a *CreateApiRequest) Validate() error {
	return validation.ValidateStruct(a,
		validation.Field(&a.Name, validation.Required, validation.Length(3, 0)),
		validation.Field(&a.Path, validation.Required, validation.Match(common.RegexEndpoint)),
		validation.Field(&a.Method, validation.In(
			http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete)),
		validation.Field(&a.Description, validation.Required, validation.Length(3, 0)),
		validation.Field(&a.Request, validation.By(
			func(value any) error {
				pMap := value.(map[string]requestvalidator.RequestParameter)
				for _, p := range pMap {
					if err := p.Validate(); err != nil {
						return err
					}
				}
				return nil
			})),
		validation.Field(&a.Response, validation.By(
			func(value any) error {
				rdMap := value.(map[uint]ResponseDefinition)
				requiredCodes := map[uint]bool{
					common.EventSuccess:           false,
					common.EventExhaust:           false,
					common.EventBadRequest:        false,
					common.EventNotFound:          false,
					common.EventSystemMalfunction: false,
				}
				for trigger, profile := range rdMap {
					if _, ok := requiredCodes[trigger]; !ok {
						return fmt.Errorf("invalid trigger %d", trigger)
					} else if err := profile.Validate(); err != nil {
						return err
					}
					requiredCodes[trigger] = true
				}
				for trigger, validated := range requiredCodes {
					if !validated {
						return fmt.Errorf("profile for trigger %d not found", trigger)
					}
				}
				return nil
			})),
		validation.Field(&a.PreConfig, validation.By(
			func(value interface{}) error {
				rMap := value.(map[string]resolvable.Resolvable)
				for _, r := range rMap {
					if err := r.Validate(); err != nil {
						return err
					}
				}
				return nil
			})),
		validation.Field(&a.Triggers, validation.Required, validation.Each(
			validation.By(func(value interface{}) error {
				t := value.(triggerflow.TriggerConditionRequest)
				return t.Validate()
			}))),
	)
}

func (r *ResponseDefinition) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.UseProfile),
		validation.Field(&r.Definition,
			validation.When(r.UseProfile == "", validation.Required).Else(validation.Empty)),
		validation.Field(&r.HTTPStatusCode,
			validation.When(r.UseProfile == "", validation.Required).Else(validation.Empty)),
	)
}
