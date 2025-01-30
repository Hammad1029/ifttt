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
	OrmSelect = "select"
	OrmUpdate = "update"
	OrmInsert = "insert"
	OrmDelete = "delete"
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
