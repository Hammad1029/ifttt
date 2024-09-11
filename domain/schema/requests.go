package schema

type CreateTableRequest struct {
	TableName   string          `mapstructure:"tableName" json:"tableName"`
	Columns     []addColumn     `mapstructure:"columns" json:"columns"`
	Constraints []addConstraint `mapstructure:"constraints" json:"constraints"`
}

type UpdateTableRequest struct {
	TableName string         `mapstructure:"tableName" json:"tableName"`
	Updates   *[]tableUpdate `mapstructure:"updates" json:"updates"`
}

type tableUpdate struct {
	UpdateType       string            `mapstructure:"updateType" json:"updateType"`
	RenameTable      *renameTable      `mapstructure:"renameTable" json:"renameTable"`
	RenameColumn     *renameColumn     `mapstructure:"renameColumn" json:"renameColumn"`
	AlterColumn      *alterColumn      `mapstructure:"alterColumn" json:"alterColumn"`
	AddColumn        *addColumn        `mapstructure:"addColumn" json:"addColumn"`
	RemoveColumn     *removeColumn     `mapstructure:"removeColumn" json:"removeColumn"`
	AddConstraint    *addConstraint    `mapstructure:"constraint" json:"constraint"`
	RemoveConstraint *removeConstraint `mapstructure:"removeConstraint" json:"removeConstraint"`
}

type renameTable struct {
	Name string `mapstructure:"name" json:"name"`
}

type renameColumn struct {
	OldName string `mapstructure:"oldName" json:"oldName"`
	NewName string `mapstructure:"newName" json:"newName"`
}

type alterColumn struct {
	ColumnName   string `mapstructure:"columnName" json:"columnName"`
	DataType     string `mapstructure:"dataType" json:"dataType"`
	Nullable     bool   `mapstructure:"nullable" json:"nullable"`
	DefaultValue string `mapstructure:"columnDefault" json:"columnDefault"`
}

type addColumn struct {
	ColumnName   string `mapstructure:"columnName" json:"columnName"`
	DataType     string `mapstructure:"dataType" json:"dataType"`
	Nullable     bool   `mapstructure:"nullable" json:"nullable"`
	DefaultValue string `mapstructure:"defaultValue" json:"defaultValue"`
}

type removeColumn struct {
	ColumnName string `mapstructure:"columnName" json:"columnName"`
}

type addConstraint struct {
	ConstraintType  string `mapstructure:"constraintType" json:"constraintType"`
	ColumnName      string `mapstructure:"columnName" json:"columnName"`
	ReferencesTable string `mapstructure:"referencesTable" json:"referencesTable"`
	ReferencesField string `mapstructure:"referencesField" json:"referencesField"`
}

type removeConstraint struct {
	ConstraintName string `mapstructure:"constraintName" json:"constraintName"`
}
