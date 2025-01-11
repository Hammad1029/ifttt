package postgres

import (
	"fmt"
	"ifttt/manager/domain/orm_schema"
	"strings"

	"gorm.io/gorm"
)

type PostgresSchemaRepository struct {
	*PostgresBaseRepository
}

func NewPostgresSchemaRepository(base *PostgresBaseRepository) *PostgresSchemaRepository {
	return &PostgresSchemaRepository{PostgresBaseRepository: base}
}

func (p *PostgresSchemaRepository) GetTableNames() ([]string, error) {
	var names []string
	if err := p.client.Table("information_schema.tables").Where(
		"table_type = ? AND table_schema NOT IN ?",
		"BASE TABLE", []string{"pg_catalog", "information_schema"},
	).Pluck("table_name", &names).Error; err != nil {
		return nil,
			fmt.Errorf("method *PostgresSchemaRepository.GetAllTables: could not get table names: %s", err)
	}
	return names, nil
}

func (p *PostgresSchemaRepository) GetAllColumns(tables []string) (*[]orm_schema.Column, error) {
	var columns []orm_schema.Column
	if err := p.client.Table("information_schema.columns").
		Select("table_name,ordinal_position,column_name,data_type,column_default,is_nullable,character_maximum_length,numeric_precision").
		Where("table_name IN ?", tables).
		Order("ordinal_position").
		Scan(&columns).Error; err != nil {
		return nil,
			fmt.Errorf("method *PostgresSchemaRepository.GetAllColumns: could not get columns: %s", err)
	}
	return &columns, nil
}

func (p *PostgresSchemaRepository) GetAllConstraints(tables []string) (*[]orm_schema.Constraint, error) {
	var constraints []orm_schema.Constraint
	if err := p.client.Table("information_schema.table_constraints AS tc").
		Select("tc.constraint_name, tc.constraint_type, tc.table_name, kcu.column_name, ccu.table_name AS references_table, ccu.column_name AS references_field").
		Joins("LEFT JOIN information_schema.key_column_usage AS kcu ON tc.constraint_catalog = kcu.constraint_catalog AND tc.constraint_schema = kcu.constraint_schema AND tc.constraint_name = kcu.constraint_name").
		Joins("LEFT JOIN information_schema.referential_constraints AS rc ON tc.constraint_catalog = rc.constraint_catalog AND tc.constraint_schema = rc.constraint_schema AND tc.constraint_name = rc.constraint_name").
		Joins("LEFT JOIN information_schema.constraint_column_usage AS ccu ON rc.unique_constraint_catalog = ccu.constraint_catalog AND rc.unique_constraint_schema = ccu.constraint_schema AND rc.unique_constraint_name = ccu.constraint_name").
		Where("tc.table_name in ?", tables).
		Scan(&constraints).Error; err != nil {
		return nil,
			fmt.Errorf("method *PostgresSchemaRepository.GetAllColumns: could not get columns: %s", err)
	}
	return &constraints, nil
}

func (p *PostgresSchemaRepository) CreateTable(newSchema *orm_schema.CreateTableRequest) error {
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

func (p *PostgresSchemaRepository) UpdateTable(updates *orm_schema.UpdateTableRequest) error {
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

func (p *PostgresSchemaRepository) GenerateAssociations(tables []string) (
	*[]orm_schema.ModelAssociation, error) {
	return nil, nil
}
