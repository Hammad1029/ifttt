package orm_schema

type SchemaRepository interface {
	GetTableNames() ([]string, error)
	GetAllColumns(tables []string) (*[]Column, error)
	GetAllConstraints(tables []string) (*[]Constraint, error)
	CreateTable(schema *CreateTableRequest) error
	UpdateTable(updates *UpdateTableRequest) error
	GenerateAssociations(tables []string) (*[]ModelAssociation, error)
}

type OrmRepository interface {
	CreateModel(model *Model) error
	GetModelByIdOrName(id uint, name string) (*Model, error)
	CreateAssociation(association *ModelAssociation) error
	GetAssociationByName(name string) (*ModelAssociation, error)
	GetAssociationByNameOrTablesAndType(name string, table1 string, table2 string, relType string) (*ModelAssociation, error)
	GetAllModels() (map[string]*Model, error)
	GetAllAssociations() (map[string]*ModelAssociation, error)
}
