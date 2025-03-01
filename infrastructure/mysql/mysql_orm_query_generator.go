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

func (m *MySqlOrmQueryGeneratorRepository) GenerateSuccessive(r *resolvable.Orm, rootModel *orm_schema.Model) (string, *[]resolvable.Resolvable, error) {
	projections := m.buildProjections(rootModel, nil, rootModel.Name)
	query := fmt.Sprintf("SELECT %s FROM `%s` AS `%s`",
		strings.Join(projections, ","), rootModel.Table, rootModel.Name)
	var whereParams []resolvable.Resolvable

	switch r.Operation {
	case common.OrmUpdate:
		if r.Where.Template != "" {
			query = fmt.Sprintf("%s WHERE %s", query, r.Where.Template)
			if len(r.Where.Values) > 0 && len(r.Query.Parameters) >= len(r.Where.Values) {
				whereParams = r.Query.Parameters[len(r.Query.Parameters)-len(r.Where.Values):]
			}
		}
	case common.OrmInsert:
		query = fmt.Sprintf("%s WHERE %s = LAST_INSERT_ID()", query, rootModel.PrimaryKey)
	default:
		return "", nil, fmt.Errorf("no successive generator for %s", r.Operation)
	}

	return query, &whereParams, nil
}

func (m *MySqlOrmQueryGeneratorRepository) GenerateUpdate(
	tableName string, where string, colSq []string,
) (string, error) {
	if len(colSq) == 0 {
		return "", fmt.Errorf("nil columns in update")
	}

	setString := make([]string, len(colSq))
	for idx, col := range colSq {
		setString[idx] = fmt.Sprintf("%s = ?", col)
	}

	queryString := fmt.Sprintf("UPDATE %s SET %s", tableName, strings.Join(setString, ", "))

	if where != "" {
		queryString += fmt.Sprintf(" WHERE %s", where)
	}

	return queryString, nil
}

func (m *MySqlOrmQueryGeneratorRepository) GenerateDelete(
	tableName string, where string,
) (string, error) {
	queryString := fmt.Sprintf("DELETE FROM %s", tableName)
	if where != "" {
		queryString += fmt.Sprintf(" WHERE %s", where)
	}
	return queryString, nil
}

func (m *MySqlOrmQueryGeneratorRepository) GenerateInsert(
	tableName string, colSq []string,
) (string, error) {
	if len(colSq) == 0 {
		return "", fmt.Errorf("nil columns in insert")
	}

	insertCols := make([]string, len(colSq))
	insertValues := make([]string, len(colSq))

	for idx, col := range colSq {
		insertCols[idx] = col
		insertValues[idx] = "?"
	}

	queryString := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		tableName, strings.Join(insertCols, ","), strings.Join(insertValues, ","))

	return queryString, nil
}

func (m *MySqlOrmQueryGeneratorRepository) GenerateSelect(
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
	if customProjections != nil && len(*customProjections) > 0 {
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
		if association.ReferencesTable == joinAlias {
			return fmt.Sprintf("%s `%s` AS `%s` ON `%s`.`%s` = `%s`.`%s`",
				joinType,
				association.TableName,
				alias,
				joinAlias,
				association.ReferencesField,
				alias,
				association.ColumnName,
			), nil
		} else {
			return fmt.Sprintf("%s `%s` AS `%s` ON `%s`.`%s` = `%s`.`%s`",
				joinType,
				association.ReferencesTable,
				alias,
				joinAlias,
				association.ColumnName,
				alias,
				association.ReferencesField,
			), nil
		}
	case common.AssociationsBelongsToMany:
		m2mAlias := fmt.Sprintf("%s_m2m_%s", joinAlias, alias)
		return fmt.Sprintf("%s `%s` AS `%s` ON `%s`.`%s` = `%s`.`%s` %s `%s` AS `%s` ON `%s`.`%s` = `%s`.`%s`",
			joinType,
			association.JoinTable,
			m2mAlias,
			m2mAlias,
			association.JoinTableSourceField,
			joinAlias,
			association.ColumnName,
			joinType,
			association.ReferencesTable,
			alias,
			alias,
			association.ReferencesField,
			m2mAlias,
			association.JoinTableTargetField,
		), nil

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
