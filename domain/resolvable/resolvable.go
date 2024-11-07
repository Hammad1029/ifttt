package resolvable

type Resolvable struct {
	ResolveType string         `json:"resolveType" mapstructure:"resolveType"`
	ResolveData map[string]any `json:"resolveData" mapstructure:"resolveData"`
}

type apiCallResolvable struct {
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

type setCacheResolvable struct {
	Key   Resolvable `json:"key" mapstructure:"key"`
	Value Resolvable `json:"value" mapstructure:"value"`
	TTL   uint       `json:"ttl" mapstructure:"ttl"`
}

type getCacheResolvable struct {
	Key Resolvable `json:"key" mapstructure:"key"`
}

type dbDumpResolvable struct {
	Columns map[string]Resolvable `json:"columns" mapstructure:"columns"`
	Table   string                `json:"table" mapstructure:"table"`
}

type encodeResolvable struct {
	Input Resolvable `json:"input" mapstructure:"input"`
	Alg   string     `json:"alg" mapstructure:"alg"`
}

type getRequestResolvable struct{}

type getResponseResolvable struct{}

type getStoreResolvable struct{}

type getPreConfigResolvable struct{}

type getHeadersResolvable struct{}

type getConstResolvable struct {
	Value any `json:"value" mapstructure:"value"`
}

type jqResolvable struct {
	Query Resolvable `json:"query" mapstructure:"query"`
	Input any        `json:"input" mapstructure:"input"`
}

type queryResolvable struct {
	QueryString          string                `json:"queryString" mapstructure:"queryString"`
	QueryHash            string                `json:"queryHash" mapstructure:"queryHash"`
	Return               bool                  `json:"return" mapstructure:"return"`
	Named                bool                  `json:"named" mapstructure:"named"`
	NamedParameters      map[string]Resolvable `json:"namedParameters" mapstructure:"namedParameters"`
	PositionalParameters []Resolvable          `json:"positionalParameters" mapstructure:"positionalParameters"`
	Async                bool                  `json:"async" mapstructure:"async"`
	Timeout              uint                  `json:"timeout" mapstructure:"timeout"`
}

type responseResolvable struct {
	ResponseCode        string `json:"responseCode" mapstructure:"responseCode"`
	ResponseDescription string `json:"responseDescription" mapstructure:"responseDescription"`
}

type setResResolvable map[string]any

type setStoreResolvable map[string]any

type setLogResolvable struct {
	LogData any    `json:"logData" mapstructure:"logData"`
	LogType string `json:"logType" mapstructure:"logType"`
}

type stringInterpolationResolvable struct {
	Template   string       `json:"template" mapstructure:"template"`
	Parameters []Resolvable `json:"parameters" mapstructure:"parameters"`
}

type uuidResolvable struct{}
