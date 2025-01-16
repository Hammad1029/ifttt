package orm_schema

type Model struct {
	ID                     uint               `mapstructure:"ID" json:"ID"`
	Name                   string             `mapstructure:"name" json:"name"`
	Table                  string             `mapstructure:"table" json:"table"`
	Projections            []Projection       `mapstructure:"projections" json:"projections"`
	PrimaryKey             string             `mapstructure:"primaryKey" json:"primaryKey"`
	OwningAssociations     []ModelAssociation `mapstructure:"owningAssociations" json:"owningAssociations"`
	ReferencedAssociations []ModelAssociation `mapstructure:"referencedAssociations" json:"referencedAssociations"`
}

type Projection struct {
	Column     string `mapstructure:"column" json:"column"`
	As         string `mapstructure:"as" json:"as"`
	SchemaType string `mapstructure:"schemaType" json:"schemaType"`
	ModelType  string `mapstructure:"modelType" json:"modelType"`
	NotNull    bool   `mapstructure:"notNull" json:"notNull"`
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
	OwningModelID        uint   `mapstructure:"owningModelID" json:"owningModelID"`
	ReferencesModelID    uint   `mapstructure:"referencesModelID" json:"referencesModelID"`
	OwningModelName      string `mapstructure:"owningModelName" json:"owningModelName"`
	ReferencesModelName  string `mapstructure:"referencesModelName" json:"referencesModelName"`
	OwningModel          Model  `mapstructure:"owningModel" json:"owningModel"`
	ReferencesModel      Model  `mapstructure:"referencesModel" json:"referencesModel"`
}

type Populate struct {
	Model    string       `mapstructure:"model" json:"model"`
	As       string       `mapstructure:"as" json:"as"`
	Project  []Projection `mapstructure:"project" json:"project"`
	Where    Where        `mapstructure:"where" json:"where"`
	Populate []Populate   `mapstructure:"populate" json:"populate"`
	Required *bool        `mapstructure:"required" json:"required"`
}

type Where struct {
	Template string `json:"template" mapstructure:"template"`
	Values   []any  `json:"values" mapstructure:"values"`
}
