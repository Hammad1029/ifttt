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
	Generate(r *Orm, rootModel *orm_schema.Model, models map[string]*orm_schema.Model) (string, error)
}
