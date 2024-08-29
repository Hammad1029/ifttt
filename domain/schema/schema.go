package schema

type Schema struct {
	TableName   string       `mapstructure:"tableName" json:"tableName"`
	Columns     []Column     `mapstructure:"columns" json:"columns"`
	Constraints []Constraint `mapstructure:"constraints" json:"constraints"`
}

type Column struct {
	TableName              string `mapstructure:"tableName" json:"tableName"`
	OrdinalPosition        int    `mapstructure:"ordinalPosition" json:"ordinalPosition"`
	ColumnName             string `mapstructure:"columnName" json:"columnName"`
	DataType               string `mapstructure:"dataType" json:"dataType"`
	ColumnDefault          string `mapstructure:"columnDefault" json:"columnDefault"`
	IsNullable             string `mapstructure:"isNullable" json:"isNullable"`
	CharacterMaximumLength int    `mapstructure:"characterMaximumLength" json:"characterMaximumLength"`
	NumericPrecision       int    `mapstructure:"numericPrecision" json:"numericPrecision"`
}

type Constraint struct {
	TableName       string `mapstructure:"tableName" json:"tableName"`
	ConstraintName  string `mapstructure:"constraintName" json:"constraintName"`
	ConstraintType  string `mapstructure:"constraintType" json:"constraintType"`
	ColumnName      string `mapstructure:"columnName" json:"columnName"`
	ReferencesTable string `mapstructure:"referencesTable" json:"referencesTable"`
	ReferencesField string `mapstructure:"referencesField" json:"referencesField"`
}
