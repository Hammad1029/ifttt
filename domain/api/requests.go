package api

import (
	"ifttt/manager/common"
	"ifttt/manager/domain/resolvable"
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
	Name         string                           `json:"name" mapstructure:"name"`
	Path         string                           `json:"path" mapstructure:"path"`
	Method       string                           `json:"method" mapstructure:"method"`
	Description  string                           `json:"description" mapstructure:"description"`
	Request      map[string]any                   `json:"request" mapstructure:"request"`
	PreConfig    map[string]resolvable.Resolvable `json:"preConfig" mapstructure:"preConfig"`
	TriggerFlows []uint                           `json:"triggerFlows" mapstructure:"triggerFlows"`
}

func (a *CreateApiRequest) Validate() error {
	return validation.ValidateStruct(a,
		validation.Field(&a.Name, validation.Required, validation.Length(3, 0)),
		validation.Field(&a.Path, validation.Required, validation.Match(common.RegexEndpoint)),
		validation.Field(&a.Method, validation.In(
			http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete)),
		validation.Field(&a.Description, validation.Required, validation.Length(3, 0)),
		validation.Field(&a.Request),
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
		validation.Field(&a.TriggerFlows, validation.Required),
	)
}
