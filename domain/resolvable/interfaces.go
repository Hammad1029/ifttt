package resolvable

import (
	"ifttt/manager/common"
	"ifttt/manager/domain/orm_schema"
)

type resolvableInterface interface {
	common.Manipulatable
	common.ValidatorInterface
}

type OrmQueryGenerator interface {
	Generate(r *OrmResolvable, models map[string]*orm_schema.Model) (string, error)
}
