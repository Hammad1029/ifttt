package orm_schema

type SchemaRepository interface {
	GetTableNames() ([]string, error)
	GetAllColumns(tables []string) (*[]Column, error)
	GetAllConstraints(tables []string) (*[]Constraint, error)
	CreateTable(schema *CreateTableRequest) error
	UpdateTable(updates *UpdateTableRequest) error
}

type OrmRepository interface {
	CreateModel(model *Model) error
	GetModelByName(name string) (*Model, error)
	CreateAssociation(association *ModelAssociation) error
	GetAssociationByName(name string) (*ModelAssociation, error)
}
