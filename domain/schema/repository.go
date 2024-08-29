package schema

type Repository interface {
	GetTableNames() ([]string, error)
	GetAllColumns(tables []string) (*[]Column, error)
	GetAllConstraints(tables []string) (*[]Constraint, error)
	CreateTable(schema *CreateTableRequest) error
	UpdateTable(updates *UpdateTableRequest) error
}
