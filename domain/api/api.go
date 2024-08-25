package api

import (
	"encoding/json"
	"fmt"

	"github.com/scylladb/gocqlx/v3/table"
)

type ApiSerialized struct {
	Group       string   `json:"group" mapstructure:"group"`
	Name        string   `json:"name" mapstructure:"name"`
	Method      string   `json:"method" mapstructure:"method"`
	Type        string   `json:"type" mapstructure:"type"`
	Path        string   `json:"path" mapstructure:"path"`
	Description string   `json:"description" mapstructure:"description"`
	Request     string   `json:"request" mapstructure:"request"`
	Dumping     string   `json:"dumping" mapstructure:"dumping"`
	StartRules  []string `json:"rules" mapstructure:"rules"`
	Rules       string   `json:"startRules" mapstructure:"startRules"`
}

type Api struct {
	Group       string           `json:"group" mapstructure:"group"`
	Name        string           `json:"name" mapstructure:"name"`
	Method      string           `json:"method" mapstructure:"method"`
	Type        string           `json:"type" mapstructure:"type"`
	Path        string           `json:"path" mapstructure:"path"`
	Description string           `json:"description" mapstructure:"description"`
	Request     map[string]any   `json:"request" mapstructure:"request"`
	Dumping     Dumping          `json:"dumping" mapstructure:"dumping"`
	StartRules  []string         `json:"startRules" mapstructure:"startRules"`
	Rules       map[string]*Rule `json:"rules" mapstructure:"rules"`
}

var ApisMetadata = table.Metadata{
	Name:    "Apis",
	Columns: []string{"api_group", "api_name", "api_description", "api_path", "api_request", "start_rules", "rules", "queries"},
	PartKey: []string{"api_group"},
	SortKey: []string{"api_name", "api_description"},
}

func (a *Api) Serialize() (*ApiSerialized, error) {
	apiModelSerialized := ApiSerialized{
		Group:       a.Group,
		Name:        a.Name,
		Method:      a.Method,
		Type:        a.Type,
		Path:        a.Path,
		Description: a.Description,
		StartRules:  a.StartRules,
	}

	requestSerialized, err := json.Marshal(a.Request)
	if err != nil {
		return nil, fmt.Errorf("method *Api.TransformApiForSave: could not serialize request: %s", err)
	}
	apiModelSerialized.Request = string(requestSerialized)

	rulesSerialized, err := json.Marshal(a.Rules)
	if err != nil {
		return nil, fmt.Errorf("method *Api.TransformApiForSave: could not serialize rules: %s", err)
	}
	apiModelSerialized.Rules = string(rulesSerialized)

	dumpingSerialized, err := json.Marshal(a.Dumping)
	if err != nil {
		return nil, fmt.Errorf("method *Api.TransformApiForSave: could not serialize dumping: %s", err)
	}
	apiModelSerialized.Dumping = string(dumpingSerialized)

	return &apiModelSerialized, nil
}
