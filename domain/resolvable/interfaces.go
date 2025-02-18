package resolvable

import (
	"ifttt/manager/common"
	"ifttt/manager/domain/orm_schema"
)

type resolvableInterface interface {
	common.Manipulatable
	common.Validatable
}

type OrmQueryGenerator interface {
	GenerateSelect(r *Orm, rootModel *orm_schema.Model, models map[string]*orm_schema.Model) (string, error)
	GenerateInsert(tableName string, colSq []string) (string, error)
	GenerateUpdate(tableName string, where string, colSq []string) (string, error)
	GenerateDelete(tableName string, where string) (string, error)
	GenerateSuccessive(r *Orm, rootModel *orm_schema.Model) (string, *[]Resolvable, error)
}
