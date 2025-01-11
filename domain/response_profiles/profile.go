package responseprofiles

type Profile struct {
	ID             uint       `json:"id" mapstructure:"id"`
	MappedCode     string     `json:"mappedCode" mapstructure:"mappedCode"`
	HttpStatus     int        `json:"httpStatus" mapstructure:"httpStatus"`
	Internal       bool       `json:"internal" mapstructure:"internal"`
	Code           Field      `json:"code" mapstructure:"code"`
	Description    Field      `json:"description" mapstructure:"description"`
	Data           Field      `json:"data" mapstructure:"data"`
	Errors         Field      `json:"errors" mapstructure:"errors"`
	MappedProfiles *[]Profile `json:"mappedProfiles" mapstructure:"mappedProfiles"`
}

type Field struct {
	Key      string `json:"key" mapstructure:"key"`
	Default  any    `json:"default" mapstructure:"default"`
	Disabled bool   `json:"disabled" mapstructure:"disabled"`
}
