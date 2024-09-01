package api

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
