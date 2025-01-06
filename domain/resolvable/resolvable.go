package resolvable

import "ifttt/manager/domain/orm_schema"

type Resolvable struct {
	ResolveType string         `json:"resolveType" mapstructure:"resolveType"`
	ResolveData map[string]any `json:"resolveData" mapstructure:"resolveData"`
}

type apiCall struct {
	Method  string         `json:"method" mapstructure:"method"`
	URL     Resolvable     `json:"url" mapstructure:"url"`
	Headers map[string]any `json:"headers" mapstructure:"headers"`
	Body    map[string]any `json:"body" mapstructure:"body"`
	Aysnc   bool           `json:"async" mapstructure:"async"`
	Timeout uint           `json:"timeout" mapstructure:"timeout"`
}

type arithmetic struct {
	Group     bool         `json:"group" mapstructure:"group"`
	Operation string       `json:"operation" mapstructure:"operation"`
	Operators []arithmetic `json:"operators" mapstructure:"operators"`
	Value     *Resolvable  `json:"value" mapstructure:"value"`
}

type setCache struct {
	Key   Resolvable `json:"key" mapstructure:"key"`
	Value Resolvable `json:"value" mapstructure:"value"`
	TTL   uint       `json:"ttl" mapstructure:"ttl"`
}

type getCache struct {
	Key Resolvable `json:"key" mapstructure:"key"`
}

type encode struct {
	Input Resolvable `json:"input" mapstructure:"input"`
	Alg   string     `json:"alg" mapstructure:"alg"`
}

type getRequest struct{}

type getResponse struct{}

type getStore struct{}

type getPreConfig struct{}

type getHeaders struct{}

type getConst struct {
	Value any `json:"value" mapstructure:"value"`
}

type jq struct {
	Query any `json:"query" mapstructure:"query"`
	Input any `json:"input" mapstructure:"input"`
}

type query struct {
	QueryString          string                `json:"queryString" mapstructure:"queryString"`
	Named                bool                  `json:"named" mapstructure:"named"`
	NamedParameters      map[string]Resolvable `json:"namedParameters" mapstructure:"namedParameters"`
	PositionalParameters []Resolvable          `json:"positionalParameters" mapstructure:"positionalParameters"`
	Async                bool                  `json:"async" mapstructure:"async"`
	Timeout              uint                  `json:"timeout" mapstructure:"timeout"`
}

type response struct {
	ResponseCode        string `json:"responseCode" mapstructure:"responseCode"`
	ResponseDescription string `json:"responseDescription" mapstructure:"responseDescription"`
}

type setRes map[string]any

type setStore map[string]any

type setLog struct {
	LogData any    `json:"logData" mapstructure:"logData"`
	LogType string `json:"logType" mapstructure:"logType"`
}

type stringInterpolation struct {
	Template   string       `json:"template" mapstructure:"template"`
	Parameters []Resolvable `json:"parameters" mapstructure:"parameters"`
}

type uuid struct{}

type cast struct {
	Input any    `json:"input" mapstructure:"input"`
	To    string `json:"to" mapstructure:"to"`
}

type Orm struct {
	Query     *query                   `json:"query" mapstructure:"query"`
	Operation string                   `json:"operation" mapstructure:"operation"`
	Model     string                   `json:"model" mapstructure:"model"`
	Project   *[]orm_schema.Projection `json:"project" mapstructure:"project"`
	Columns   *map[string]any          `json:"columns" mapstructure:"columns"`
	Populate  *[]orm_schema.Populate   `json:"populate" mapstructure:"populate"`
	Where     *orm_schema.Where        `json:"where" mapstructure:"where"`
	OrderBy   string                   `json:"orderBy" mapstructure:"orderBy"`
	Limit     int                      `json:"limit" mapstructure:"limit"`
}

type forEach struct {
	Input any           `json:"input" mapstructure:"input"`
	Do    *[]Resolvable `json:"do" mapstructure:"do"`
}

type getIter struct {
}

type dateFunc struct {
	Input        dateInput         `json:"input" mapstructure:"input"`
	Manipulators []dateManipulator `json:"manipulators" mapstructure:"manipulators"`
	Format       string            `json:"format" mapstructure:"format"`
	UTC          bool              `json:"utc" mapstructure:"utc"`
}

type dateManipulator struct {
	Operator string     `json:"operator" mapstructure:"operator"`
	Operand  Resolvable `json:"operand" mapstructure:"operand"`
	Unit     string     `json:"unit" mapstructure:"unit"`
}

type dateInput struct {
	Input    *Resolvable `json:"input" mapstructure:"input"`
	Parse    string      `json:"parse" mapstructure:"parse"`
	Timezone string      `json:"timezone" mapstructure:"timezone"`
}
