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
	r *resolvable.Orm, rootModel *orm_schema.Model, models map[string]*orm_schema.Model,
) (string, error) {
	switch r.Operation {
	case common.OrmSelect:
		return m.generateSelect(r, rootModel, models)
	case common.OrmInsert:
		return m.generateInsert(r, rootModel)
	default:
		return "", fmt.Errorf("generator for %s not available", r.Operation)
	}
}

func (m *MySqlOrmQueryGeneratorRepository) generateInsert(
	r *resolvable.Orm, rootModel *orm_schema.Model,
) (string, error) {
	if r.Columns == nil || len(*r.Columns) == 0 {
		return "", fmt.Errorf("nil columns in insert")
	}

	insertCols := make([]string, len(*r.Columns))
	insertValues := make([]string, len(*r.Columns))

	idx := 0
	for col := range *r.Columns {
		insertCols[idx] = col
		insertValues[idx] = fmt.Sprintf("@%s", col)
		idx++
	}

	queryString := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		rootModel.Table, strings.Join(insertCols, ","), strings.Join(insertValues, ","))

	return queryString, nil
}

func (m *MySqlOrmQueryGeneratorRepository) generateSelect(
	r *resolvable.Orm, rootModel *orm_schema.Model, models map[string]*orm_schema.Model,
) (string, error) {
	var joins []string
	projections := m.buildProjections(rootModel, r.Project, rootModel.Name)
	if joinProjections, joinClauses, err :=
		m.buildAssociations(r.Populate, rootModel, models, rootModel.Name); err != nil {
		return "", err
	} else {
		projections = append(projections, joinProjections...)
		joins = joinClauses
	}

	query := fmt.Sprintf("SELECT %s FROM `%s` AS `%s` %s",
		strings.Join(projections, ","), rootModel.Table, rootModel.Name,
		strings.Join(joins, " "))

	if r.Where.Template != "" {
		query += fmt.Sprintf(" WHERE %s", r.Where.Template)
	}
	if r.OrderBy != "" {
		query += fmt.Sprintf(" ORDER BY %s", r.OrderBy)
	} else {
		query += fmt.Sprintf(" ORDER BY %s.%s ASC", rootModel.Name, rootModel.PrimaryKey)
	}
	if r.Limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", r.Limit)
	}

	return query, nil
}

func (m *MySqlOrmQueryGeneratorRepository) buildProjections(
	model *orm_schema.Model, customProjections *[]orm_schema.Projection, alias string) []string {
	var projs []string
	var pKey bool
	if len(*customProjections) > 0 {
		for _, p := range *customProjections {
			if p.Column == model.PrimaryKey {
				pKey = true
			}
			projs = append(projs, fmt.Sprintf("`%s`.`%s` AS `%s.%s`", alias, p.Column, alias, p.As))
		}
	} else {
		for _, p := range model.Projections {
			if p.Column == model.PrimaryKey {
				pKey = true
			}
			projs = append(projs, fmt.Sprintf("`%s`.`%s` AS `%s.%s`", alias, p.Column, alias, p.As))
		}
	}
	if !pKey {
		projs = append(projs, fmt.Sprintf("`%s`.`%s` AS `%s.%s`", alias, model.PrimaryKey, alias, model.PrimaryKey))
	}
	return projs
}

func (m *MySqlOrmQueryGeneratorRepository) buildAssociations(
	populate *[]orm_schema.Populate, parent *orm_schema.Model, models map[string]*orm_schema.Model, alias string,
) ([]string, []string, error) {
	joinProjections := []string{}
	joinClauses := []string{}

	for _, p := range *populate {
		joinModel, ok := models[p.Model]
		if !ok {
			return nil, nil, fmt.Errorf("model %s not found", p.Model)
		}

		association := m.findAssociation(parent, joinModel)
		if association == nil {
			return nil, nil, fmt.Errorf("association between %s and %s not found", parent.Name, joinModel.Name)
		}

		concatAlias := fmt.Sprintf("%s_%s", alias, p.As)

		hasRequiredNestedWhere := false
		for _, nested := range p.Populate {
			if nested.Where.Template != "" && m.isJoinRequired(&nested) {
				hasRequiredNestedWhere = true
				break
			}
		}

		if hasRequiredNestedWhere && !m.isJoinRequired(&p) {

			joinClauses = append(joinClauses, "LEFT OUTER JOIN (")
			joinClauses = append(joinClauses, fmt.Sprintf("`%s` AS `%s`", association.ReferencesTable, concatAlias))

			if childProjections, childJoins, err := m.buildAssociations(&p.Populate, joinModel, models, concatAlias); err != nil {
				return nil, nil, err
			} else {
				joinProjections = append(joinProjections, childProjections...)
				joinClauses = append(joinClauses, childJoins...)
			}

			joinClauses = append(joinClauses, fmt.Sprintf(") ON `%s`.`%s` = `%s`.`%s`",
				concatAlias,
				association.ReferencesField,
				alias,
				association.ColumnName))
		} else {

			if joinClause, err := m.buildJoin(association, concatAlias, alias, m.isJoinRequired(&p)); err != nil {
				return nil, nil, err
			} else {
				joinClauses = append(joinClauses, joinClause)
				if p.Where.Template != "" {
					joinClauses = append(joinClauses, fmt.Sprintf("AND %s", p.Where.Template))
				}
			}

			if childProjections, childJoins, err := m.buildAssociations(&p.Populate, joinModel, models, concatAlias); err != nil {
				return nil, nil, err
			} else {
				joinProjections = append(joinProjections, childProjections...)
				joinClauses = append(joinClauses, childJoins...)
			}
		}

		joinProjections = append(joinProjections, m.buildProjections(joinModel, &p.Project, concatAlias)...)
	}

	return joinProjections, joinClauses, nil
}

func (m *MySqlOrmQueryGeneratorRepository) buildJoin(
	association *orm_schema.ModelAssociation,
	alias string,
	joinAlias string,
	required bool,
) (string, error) {
	joinType := "LEFT OUTER JOIN"
	if required {
		joinType = "INNER JOIN"
	}

	switch association.Type {
	case common.AssociationsHasOne, common.AssociationsHasMany, common.AssociationsBelongsTo:
		return fmt.Sprintf("%s `%s` AS `%s` ON `%s`.`%s` = `%s`.`%s`",
			joinType,
			association.ReferencesTable,
			alias,
			joinAlias,
			association.ColumnName,
			alias,
			association.ReferencesField,
		), nil

	case common.AssociationsBelongsToMany:
		return fmt.Sprintf("%s `%s` AS `%s` ON `%s`.`%s` = `%s`.`%s` %s `%s` ON `%s`.`%s` = `%s`.`%s`",
			joinType,
			association.JoinTable,
			alias,
			alias,
			association.JoinTableSourceField,
			joinAlias,
			association.ColumnName,
			joinType,
			association.ReferencesTable,
			association.ReferencesTable,
			association.ReferencesField,
			association.JoinTable,
			association.JoinTableTargetField), nil

	default:
		return "", fmt.Errorf("unsupported association type: %s", association.Type)
	}
}

func (m *MySqlOrmQueryGeneratorRepository) isJoinRequired(p *orm_schema.Populate) bool {
	if p.Required != nil {
		return *p.Required
	}

	return p.Where.Template != ""
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
