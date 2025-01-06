package infrastructure

import (
	"fmt"
	"ifttt/manager/domain/orm_schema"
	"strings"

	"gorm.io/gorm"
)

type MySqlSchemaRepository struct {
	*MySqlBaseRepository
}

func NewMySqlSchemaRepository(base *MySqlBaseRepository) *MySqlSchemaRepository {
	return &MySqlSchemaRepository{MySqlBaseRepository: base}
}

func (p *MySqlSchemaRepository) GetTableNames() ([]string, error) {
	var names []string
	if err := p.client.Table("information_schema.tables").Where(
		"table_schema = DATABASE()",
	).Pluck("table_name", &names).Error; err != nil {
		return nil,
			fmt.Errorf("method *PostgresSchemaRepository.GetAllTables: could not get table names: %s", err)
	}
	return names, nil
}

func (p *MySqlSchemaRepository) GetAllColumns(tables []string) (*[]orm_schema.Column, error) {
	var columns []orm_schema.Column
	if err := p.client.Table("information_schema.columns").
		Select(
			"TABLE_NAME as TableName, "+
				"ORDINAL_POSITION as OrdinalPosition, "+
				"COLUMN_NAME as ColumnName, "+
				"DATA_TYPE as DataType, "+
				"COLUMN_DEFAULT as ColumnDefault, "+
				"IS_NULLABLE as IsNullable, "+
				"CHARACTER_MAXIMUM_LENGTH as CharacterMaximumLength, "+
				"NUMERIC_PRECISION as NumericPrecision",
		).
		Where("table_name IN ?", tables).
		Order("ordinal_position").
		Scan(&columns).Error; err != nil {
		return nil,
			fmt.Errorf("could not get columns: %s", err)
	}
	return &columns, nil
}

func (p *MySqlSchemaRepository) GetAllConstraints(tables []string) (*[]orm_schema.Constraint, error) {
	var constraints []orm_schema.Constraint
	if err := p.client.Table("information_schema.table_constraints AS tc").
		Select(
			"tc.table_name as TableName, "+
				"tc.constraint_name as ConstraintName, "+
				"tc.constraint_type as ConstraintType, "+
				"kcu.column_name as ColumnName, "+
				"kcu.referenced_table_name as ReferencesTable, "+
				"kcu.referenced_column_name as ReferencesField",
		).
		Joins("LEFT JOIN information_schema.key_column_usage AS kcu ON "+
			"kcu.table_schema = tc.table_schema AND "+
			"kcu.table_name = tc.table_name AND "+
			"kcu.constraint_name = tc.constraint_name").
		Where("tc.table_schema = DATABASE() AND tc.table_name IN ?", tables).
		Scan(&constraints).Error; err != nil {
		return nil,
			fmt.Errorf("method *PostgresSchemaRepository.GetAllColumns: could not get columns: %s", err)
	}
	return &constraints, nil
}

func (p *MySqlSchemaRepository) CreateTable(newSchema *orm_schema.CreateTableRequest) error {
	var columnDefs []string
	for _, col := range newSchema.Columns {
		colDef := fmt.Sprintf("%s %s", col.ColumnName, col.DataType)
		if col.DefaultValue != "" {
			colDef += fmt.Sprintf(" DEFAULT %s", col.DefaultValue)
		}
		if !col.Nullable {
			colDef += " NOT NULL"
		}
		columnDefs = append(columnDefs, colDef)
	}

	var constraintDefs []string
	for _, constr := range newSchema.Constraints {
		switch constr.ConstraintType {
		case orm_schema.PrimaryKeyConstraintKey:
			constraintDefs = append(constraintDefs, fmt.Sprintf("PRIMARY KEY (%s)", constr.ColumnName))
		case orm_schema.UniqueConstraintKey:
			constraintDefs = append(constraintDefs, fmt.Sprintf("UNIQUE (%s)", constr.ColumnName))
		case orm_schema.ForeignKeyConstraintKey:
			constraintDefs = append(constraintDefs, fmt.Sprintf("FOREIGN KEY (%s) REFERENCES %s(%s)",
				constr.ColumnName, constr.ReferencesTable, constr.ReferencesField))
		default:
			return fmt.Errorf("method *PostgresSchemaRepository.CreateTable: constraint type %s not supported", constr.ConstraintType)
		}
	}

	query := fmt.Sprintf("CREATE TABLE %s (\n%s", newSchema.TableName, strings.Join(columnDefs, ",\n"))
	if len(constraintDefs) > 0 {
		query += fmt.Sprintf(",\n%s", strings.Join(constraintDefs, ",\n"))
	}
	query += "\n);"

	if err := p.client.Exec(query).Error; err != nil {
		return fmt.Errorf("method *PostgresSchemaRepository.CreateTable: could not create table: %s", err)
	}
	return nil
}

func (p *MySqlSchemaRepository) UpdateTable(updates *orm_schema.UpdateTableRequest) error {
	var transactionQueries []string
	baseAlterQuery := fmt.Sprintf("ALTER TABLE %s ", updates.TableName)

	for _, update := range updates.Updates {
		switch update.UpdateType {
		case orm_schema.RenameTableKey:
			{
				transactionQueries = append(transactionQueries,
					baseAlterQuery+fmt.Sprintf("RENAME TO %s;", update.RenameTable.Name))
			}
		case orm_schema.RenameColumnKey:
			{
				transactionQueries = append(transactionQueries,
					baseAlterQuery+fmt.Sprintf("RENAME COLUMN %s TO %s;", update.RenameColumn.OldName, update.RenameColumn.NewName))
			}
		case orm_schema.AlterColumnKey:
			{
				transactionQueries = append(transactionQueries,
					baseAlterQuery+fmt.Sprintf(
						"ALTER COLUMN %s SET DATA TYPE %s;",
						update.AlterColumn.ColumnName, update.AlterColumn.DataType))
				if update.AlterColumn.Nullable {
					transactionQueries = append(transactionQueries,
						baseAlterQuery+fmt.Sprintf(
							"ALTER COLUMN %s DROP NOT NULL;", update.AlterColumn.ColumnName))
				} else {
					transactionQueries = append(transactionQueries,
						baseAlterQuery+fmt.Sprintf(
							"ALTER COLUMN %s SET NOT NULL;", update.AlterColumn.ColumnName))
				}
				if update.AlterColumn.DefaultValue == "" {
					transactionQueries = append(transactionQueries,
						baseAlterQuery+fmt.Sprintf(
							"ALTER COLUMN %s DROP DEFAULT;", update.AlterColumn.ColumnName))
				} else {
					transactionQueries = append(transactionQueries,
						baseAlterQuery+fmt.Sprintf("ALTER COLUMN %s SET DEFAULT %s;",
							update.AlterColumn.ColumnName, update.AlterColumn.DefaultValue))
				}
			}
		case orm_schema.AddColumnKey:
			{
				var defaultValueQuery string
				if update.AddColumn.DefaultValue != "" {
					defaultValueQuery = fmt.Sprintf("DEFAULT %s", update.AddColumn.DefaultValue)
				}
				transactionQueries = append(transactionQueries,
					baseAlterQuery+fmt.Sprintf(
						"ADD COLUMN %s %s %s;", update.AddColumn.ColumnName,
						update.AddColumn.DataType, defaultValueQuery))
			}
		case orm_schema.RemoveColumnKey:
			{
				transactionQueries = append(transactionQueries,
					baseAlterQuery+fmt.Sprintf("DROP COLUMN %s;", update.RemoveColumn.ColumnName))
			}
		case orm_schema.AddConstraintKey:
			{
				constraintQuery := baseAlterQuery + "ADD "
				switch update.AddConstraint.ConstraintType {
				case orm_schema.ForeignKeyConstraintKey:
					constraintQuery += fmt.Sprintf("FOREIGN KEY (%s) REFERENCES %s(%s)",
						update.AddConstraint.ColumnName, update.AddConstraint.ReferencesTable,
						update.AddConstraint.ReferencesField)
				case orm_schema.UniqueConstraintKey:
					constraintQuery += fmt.Sprintf("UNIQUE (%s)", update.AddConstraint.ColumnName)
				}
				transactionQueries = append(transactionQueries, fmt.Sprintf("%s;", constraintQuery))
			}
		case orm_schema.RemoveConstraintKey:
			{
				transactionQueries = append(transactionQueries,
					baseAlterQuery+fmt.Sprintf("DROP CONSTRAINT %s;", update.RemoveConstraint.ConstraintName))
			}
		}
	}

	if err := p.client.Transaction(func(tx *gorm.DB) error {
		for _, query := range transactionQueries {
			if err := tx.Exec(query).Error; err != nil {
				return fmt.Errorf("query failed: %s", err)
			}
		}
		return nil
	}); err != nil {
		return fmt.Errorf(
			"method *PostgresSchemaRepository.UpdateTable: update table transaction failed. rolling back: %s", err)
	}

	return nil
}

func (m *MySqlSchemaRepository) GenerateAssociations(tables []string) (
	*[]orm_schema.ModelAssociation, error) {
	return nil, nil
}
