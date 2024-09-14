package resolvable

type Resolvable struct {
	ResolveType string         `json:"resolveType" mapstructure:"resolveType"`
	ResolveData map[string]any `json:"resolveData" mapstructure:"resolveData"`
}

type apiCallResolvable struct {
	Method  string         `json:"method" mapstructure:"method"`
	Url     string         `json:"url" mapstructure:"url"`
	Headers map[string]any `json:"headers" mapstructure:"headers"`
	Body    map[string]any `json:"body" mapstructure:"body"`
}

type arithmetic struct {
	Group     bool         `json:"group" mapstructure:"group"`
	Operation string       `json:"operation" mapstructure:"operation"`
	Operators []arithmetic `json:"operators" mapstructure:"operators"`
	Value     Resolvable   `json:"value" mapstructure:"value"`
}

// type getRequestResolvable map[string]any

// type getResponseResolvable map[string]any

// type getStoreResolvable map[string]any

// type getApiResultsResolvable map[string]map[string]any

// type getQueryResultsResolvable map[string][]map[string]any

// type getConstResolvable struct {
// 	Value any `json:"value" mapstructure:"value"`
// }

type preConfigResolvable map[string]any

type jqResolvable struct {
	Query Resolvable `json:"query" mapstructure:"query"`
	Input any        `json:"input" mapstructure:"input"`
}

type callRuleResolvable struct {
	RuleId uint `json:"ruleId" mapstructure:"ruleId"`
}

type stringInterpolationResolvable struct {
	Template   string       `json:"template" mapstructure:"template"`
	Parameters []Resolvable `json:"parameters" mapstructure:"parameters"`
}

type queryResolvable struct {
	QueryString          string                `json:"queryString" mapstructure:"queryString"`
	QueryHash            string                `json:"queryHash" mapstructure:"queryHash"`
	Return               bool                  `json:"return" mapstructure:"return"`
	Named                bool                  `json:"named" mapstructure:"named"`
	NamedParameters      map[string]Resolvable `json:"namedParameters" mapstructure:"namedParameters"`
	PositionalParameters []Resolvable          `json:"positionalParameters" mapstructure:"positionalParameters"`
}

type responseResolvable struct {
	ResponseCode        string `json:"responseCode" mapstructure:"responseCode"`
	ResponseDescription string `json:"responseDescription" mapstructure:"responseDescription"`
}

type setResResolvable map[string]any

type setStoreResolvable map[string]any
