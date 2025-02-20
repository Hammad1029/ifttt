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

type deleteCache struct {
	Key Resolvable `json:"key" mapstructure:"key"`
}

type encode struct {
	Input Resolvable `json:"input" mapstructure:"input"`
	Alg   string     `json:"alg" mapstructure:"alg"`
}

type getErrors struct{}

type getStore struct {
	Query any `json:"query" mapstructure:"query"`
}

type getHeaders struct{}

type getConst struct {
	Value any `json:"value" mapstructure:"value"`
}

type jq struct {
	Query any `json:"query" mapstructure:"query"`
	Input any `json:"input" mapstructure:"input"`
}

type query struct {
	QueryString string       `json:"queryString" mapstructure:"queryString"`
	Scan        bool         `json:"scan" mapstructure:"scan"`
	Parameters  []Resolvable `json:"parameters" mapstructure:"parameters"`
	Async       bool         `json:"async" mapstructure:"async"`
	Timeout     uint         `json:"timeout" mapstructure:"timeout"`
}

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
	Query           *query                   `json:"query" mapstructure:"query"`
	SuccessiveQuery *query                   `json:"successiveQuery" mapstructure:"successiveQuery"`
	Operation       string                   `json:"operation" mapstructure:"operation"`
	Model           string                   `json:"model" mapstructure:"model"`
	Project         *[]orm_schema.Projection `json:"project" mapstructure:"project"`
	Columns         map[string]any           `json:"columns" mapstructure:"columns"`
	Populate        *[]orm_schema.Populate   `json:"populate" mapstructure:"populate"`
	Where           *orm_schema.Where        `json:"where" mapstructure:"where"`
	OrderBy         string                   `json:"orderBy" mapstructure:"orderBy"`
	Limit           int                      `json:"limit" mapstructure:"limit"`
	ModelsInUse     *[]string                `json:"modelsInUse" mapstructure:"modelsInUse"`
}

type filterMap struct {
	Input     any           `json:"input" mapstructure:"input"`
	Do        *[]Resolvable `json:"do" mapstructure:"do"`
	Condition Condition     `json:"condition" mapstructure:"condition"`
	Async     bool          `json:"async" mapstructure:"async"`
}

type getIter struct {
	Index bool `json:"index" mapstructure:"index"`
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

type dateIntervals struct {
	Start  dateInput `json:"start" mapstructure:"start"`
	End    dateInput `json:"end" mapstructure:"end"`
	Unit   string    `json:"unit" mapstructure:"unit"`
	Format string    `json:"format" mapstructure:"format"`
}

type response struct {
	Event uint `json:"event" mapstructure:"event"`
}

type Condition struct {
	ConditionType   string      `json:"conditionType" mapstructure:"conditionType"`
	Conditions      []Condition `json:"conditions" mapstructure:"conditions"`
	Group           bool        `json:"group" mapstructure:"group"`
	ComparisionType string      `json:"comparisionType" mapstructure:"comparisionType"`
	Operator1       *Resolvable `json:"op1" mapstructure:"op1"`
	Operand         string      `json:"opnd" mapstructure:"opnd"`
	Operator2       *Resolvable `json:"op2" mapstructure:"op2"`
}

type conditional struct {
	Condition Condition    `json:"condition" mapstructure:"condition"`
	True      []Resolvable `json:"true" mapstructure:"true"`
	False     []Resolvable `json:"false" mapstructure:"false"`
}
