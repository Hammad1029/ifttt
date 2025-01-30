package requestvalidator

const (
	dataTypeText    = "text"
	dataTypeNumber  = "number"
	dataTypeBoolean = "boolean"
	dataTypeArray   = "array"
	dataTypeMap     = "map"
)

type RequestParameter struct {
	Regex       string         `json:"regex" mapstructure:"regex"`
	DataType    string         `json:"dataType" mapstructure:"dataType"`
	Required    bool           `json:"required" mapstructure:"required"`
	InternalTag string         `json:"internalTag" mapstructure:"internalTag"`
	Config      map[string]any `json:"config" mapstructure:"config"`
}

type textValue struct {
	Alpha   bool  `json:"alpha" mapstructure:"alpha"`
	Numeric bool  `json:"alphanumeric" mapstructure:"numeric"`
	Special bool  `json:"special" mapstructure:"special"`
	Minimum int   `json:"minimum" mapstructure:"minimum"`
	Maximum int   `json:"maximum" mapstructure:"maximum"`
	In      []any `json:"in" mapstructure:"in"`
}

type numberValue struct {
	Minimum int   `json:"minimum" mapstructure:"minimum"`
	Maximum int   `json:"maximum" mapstructure:"maximum"`
	In      []any `json:"in" mapstructure:"in"`
}

type booleanValue struct {
}

type arrayValue struct {
	OfType  *RequestParameter `json:"ofType" mapstructure:"ofType"`
	Minimum int               `json:"minimum" mapstructure:"minimum"`
	Maximum int               `json:"maximum" mapstructure:"maximum"`
}

type mapValue map[string]RequestParameter
