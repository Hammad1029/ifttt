package orm_schema

type Model struct {
	Name        string       `mapstructure:"name" json:"name"`
	Table       string       `mapstructure:"table" json:"table"`
	PrimaryKey  string       `mapstructure:"primaryKey" json:"primaryKey"`
	Projections []Projection `mapstructure:"projections" json:"projections"`
}

type Projection struct {
	Column   string `mapstructure:"column" json:"column"`
	As       string `mapstructure:"as" json:"as"`
	DataType string `mapstructure:"dataType" json:"dataType"`
}

type ModelAssociation struct {
	Name                 string `mapstructure:"name" json:"name"`
	Type                 string `mapstructure:"type" json:"type"`
	TableName            string `mapstructure:"tableName" json:"tableName"`
	ColumnName           string `mapstructure:"columnName" json:"columnName"`
	ReferencesTable      string `mapstructure:"referencesTable" json:"referencesTable"`
	ReferencesField      string `mapstructure:"referencesField" json:"referencesField"`
	JoinTable            string `mapstructure:"joinTable" json:"joinTable"`
	JoinTableSourceField string `mapstructure:"joinTableSourceField" json:"joinTableSourceField"`
	JoinTableTargetField string `mapstructure:"joinTableTargetField" json:"joinTableTargetField"`
}

type Populate struct {
	Association string     `mapstructure:"association" json:"association"`
	Populate    []Populate `mapstructure:"populate" json:"populate"`
}
