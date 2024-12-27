package infrastructure

import (
	"fmt"
	"ifttt/manager/common"
	"ifttt/manager/domain/orm_schema"
	"ifttt/manager/domain/resolvable"
	"strings"
)

type MySqlOrmQueryGeneratorRepository struct {
	*MySqlBaseRepository
}

func NewMySqlOrmQueryGeneratorRepository(base *MySqlBaseRepository) *MySqlOrmQueryGeneratorRepository {
	return &MySqlOrmQueryGeneratorRepository{MySqlBaseRepository: base}
}

func (m *MySqlOrmQueryGeneratorRepository) Generate(
	r *resolvable.OrmResolvable, models map[string]*orm_schema.Model) (string, error) {
	rootModel, ok := models[r.Model]
	if !ok {
		return "", fmt.Errorf("root model %s not found", r.Model)
	}

	var joins []string
	projections := m.buildProjections(rootModel, rootModel.Name)
	if joinProjections, joinClauses, err :=
		m.buildAssociations(r.Populate, rootModel, models, rootModel.Name); err != nil {
		return "", err
	} else {
		projections = append(projections, joinProjections...)
		joins = joinClauses
	}

	query := fmt.Sprintf("SELECT %s FROM `%s` AS `%s` %s WHERE %s",
		strings.Join(projections, ","), rootModel.Table, rootModel.Name,
		strings.Join(joins, " "), r.ConditionsTemplate)

	return query, nil
}

func (m *MySqlOrmQueryGeneratorRepository) buildProjections(model *orm_schema.Model, alias string) []string {
	var projs []string
	for _, p := range model.Projections {
		projs = append(projs, fmt.Sprintf("`%s`.`%s` AS `%s.%s`", alias, p.Column, alias, p.As))
	}
	return projs
}

func (m *MySqlOrmQueryGeneratorRepository) buildAssociations(
	populate []orm_schema.Populate, parent *orm_schema.Model, models map[string]*orm_schema.Model, alias string,
) ([]string, []string, error) {
	joinProjections := []string{}
	joinClauses := []string{}

	for _, p := range populate {
		joinModel, ok := models[p.Model]
		if !ok {
			return nil, nil, fmt.Errorf("model %s not found", p.Model)
		}

		association := m.findAssociation(parent, joinModel)
		if association == nil {
			return nil, nil, fmt.Errorf("association between %s and %s not found", parent.Name, joinModel.Name)
		}

		concatAlias := fmt.Sprintf("%s_%s", alias, p.As)

		if joinClause, err := m.buildJoin(association, concatAlias, alias); err != nil {
			return nil, nil, err
		} else {
			joinClauses = append(joinClauses, joinClause)
		}
		joinProjections = append(joinProjections, m.buildProjections(joinModel, concatAlias)...)

		if childProjections, childJoins, err := m.buildAssociations(p.Populate, joinModel, models, concatAlias); err != nil {
			return nil, nil, err
		} else {
			joinProjections = append(joinProjections, childProjections...)
			joinClauses = append(joinClauses, childJoins...)
		}
	}

	return joinProjections, joinClauses, nil
}

func (m *MySqlOrmQueryGeneratorRepository) buildJoin(association *orm_schema.ModelAssociation, alias string, joinAlias string) (string, error) {
	switch association.Type {
	case common.AssociationsHasOne:
		return fmt.Sprintf("LEFT OUTER JOIN `%s` AS `%s` ON `%s`.`%s` = `%s`.`%s`",
			association.ReferencesTable,
			alias,
			alias,
			association.ReferencesField,
			joinAlias,
			association.ColumnName), nil

	case common.AssociationsHasMany:
		return fmt.Sprintf("LEFT OUTER JOIN `%s` AS `%s` ON `%s`.`%s` = `%s`.`%s`",
			association.ReferencesTable,
			alias,
			alias,
			association.ReferencesField,
			joinAlias,
			association.ColumnName), nil

	case common.AssociationsBelongsTo:
		return fmt.Sprintf("LEFT OUTER JOIN `%s` AS `%s` ON `%s`.`%s` = `%s`.`%s`",
			association.ReferencesTable,
			alias,
			alias,
			association.ReferencesField,
			joinAlias,
			association.ColumnName), nil

	case common.AssociationsBelongsToMany:
		return fmt.Sprintf("LEFT OUTER JOIN `%s` AS `%s` ON `%s`.`%s` = `%s`.`%s` LEFT JOIN `%s` ON `%s`.`%s` = `%s`.`%s`",
			association.JoinTable,
			alias,
			alias,
			association.JoinTableSourceField,
			joinAlias,
			association.ColumnName,
			association.ReferencesTable,
			association.ReferencesTable,
			association.ReferencesField,
			association.JoinTable,
			association.JoinTableTargetField), nil

	default:
		return "", fmt.Errorf("unsupported association type: %s", association.Type)
	}
}

func (m *MySqlOrmQueryGeneratorRepository) findAssociation(
	parent *orm_schema.Model, join *orm_schema.Model) *orm_schema.ModelAssociation {
	for _, a := range parent.OwningAssociations {
		if a.ReferencesModel.Name == join.Name {
			return &a
		}
	}
	for _, a := range parent.ReferencedAssociations {
		if a.OwningModel.Name == join.Name {
			return &a
		}
	}
	return nil
}
