package common

var RuleAllowedReturns = []uint{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
var RuleDefaultReturn uint = 0

const (
	AccessTokenKey  = "user-session-at"
	RefreshTokenKey = "user-session-rt"
)

const (
	LogError = "error"
	LogInfo  = "info"
)

const (
	EncodeMD5          = "md5"
	EncodeSHA1         = "sha1"
	EncodeSHA2         = "sha2"
	EncodeBcrypt       = "bcrypt"
	EncodeBase64Decode = "base64-de"
	EncodeBase64Encode = "base64-en"
)

const (
	CastToString  = "string"
	CastToNumber  = "number"
	CastToBoolean = "boolean"
)

const (
	DatabaseTypeString  = "string"
	DatabaseTypeNumber  = "number"
	DatabaseTypeBoolean = "boolean"
)

const (
	AssociationsHasOne        = "hasOne"
	AssociationsHasMany       = "hasMany"
	AssociationsBelongsTo     = "belongsTo"
	AssociationsBelongsToMany = "belongsToMany"
)

const (
	DbNameMySql    = "mysql"
	DbNamePostgres = "postgres"
	DbNameRedis    = "redis"
)

const (
	OrmSelect = "SELECT"
	OrmUpdate = "UPDATE"
	OrmInsert = "INSERT"
	OrmDelete = "DELETE"
)

const (
	DependencyOrmQueryRepo = iota
	DependencyOrmSchemaRepo
	DependencyInternalTagRepo
)

const (
	DateOperatorAdd      = "+"
	DateOperatorSubtract = "-"
)

const (
	EventSuccess           = 0
	EventExhaust           = 10
	EventBadRequest        = 400
	EventNotFound          = 404
	EventSystemMalfunction = 500
)

const (
	ComparatorEquals            = "eq"
	ComparatorNotEquals         = "ne"
	ComparatorIn                = "in"
	ComparatorNotIn             = "ni"
	ComparatorLessThan          = "lt"
	ComparatorLessThanEquals    = "lte"
	ComparatorGreaterThan       = "gt"
	ComparatorGreaterThanEquals = "gte"
)

const (
	ComparisionTypeString  = "string"
	ComparisionTypeNumber  = "number"
	ComparisionTypeBoolean = "boolean"
	ComparisionTypeDate    = "date"
	ComparisionTypeBcrypt  = "bcrypt"
)

const (
	CalculatorAdd      = "+"
	CalculatorSubtract = "-"
	CalculatorMultiply = "*"
	CalculatorDivide   = "/"
	CalculatorModulus  = "%"
)

var DateManipulatorUnits = []string{
	"y", "year", "years",
	"Q", "quarter", "quarters",
	"M", "month", "months",
	"w", "week", "weeks",
	"d", "day", "days",
	"h", "hour", "hours",
	"m", "minute", "minutes",
	"s", "second", "seconds",
	"ms", "millisecond", "milliseconds",
	"ns", "nanosecond", "nanoseconds",
}

var (
	ConditionTypeAnd = "AND"
	ConditionTypeOr  = "OR"
	OperandTypes     = []string{
		ComparatorEquals,
		ComparatorNotEquals,
		ComparatorIn,
		ComparatorNotIn,
		ComparatorLessThan,
		ComparatorLessThanEquals,
		ComparatorGreaterThan,
		ComparatorGreaterThanEquals,
	}
)
