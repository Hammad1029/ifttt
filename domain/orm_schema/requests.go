package orm_schema

import (
	"context"
	"ifttt/manager/common"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	RenameTableKey      = "renameTable"
	RenameColumnKey     = "renameColumn"
	AlterColumnKey      = "alterColumn"
	AddColumnKey        = "addColumn"
	RemoveColumnKey     = "removeColumn"
	AddConstraintKey    = "addConstraint"
	RemoveConstraintKey = "removeConstraint"
)

const (
	PrimaryKeyConstraintKey = "PRIMARY KEY"
	ForeignKeyConstraintKey = "FOREIGN KEY"
	UniqueConstraintKey     = "UNIQUE"
)

type CreateTableRequest struct {
	TableName   string          `mapstructure:"tableName" json:"tableName"`
	Columns     []addColumn     `mapstructure:"columns" json:"columns"`
	Constraints []addConstraint `mapstructure:"constraints" json:"constraints"`
}

func (c *CreateTableRequest) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.TableName, validation.Required),
		validation.Field(&c.Columns, validation.Required, validation.Length(1, 0),
			validation.Each(validation.By(func(value interface{}) error {
				col := value.(addColumn)
				return col.Validate()
			}))),
		validation.Field(&c.Constraints, validation.Required, validation.Length(1, 0),
			validation.Each(validation.By(func(value interface{}) error {
				constraint := value.(addConstraint)
				return constraint.Validate()
			}))),
	)
}

type UpdateTableRequest struct {
	TableName string        `mapstructure:"tableName" json:"tableName"`
	Updates   []tableUpdate `mapstructure:"updates" json:"updates"`
}

func (u *UpdateTableRequest) Validate() error {
	return validation.ValidateStruct(u,
		validation.Field(&u.TableName, validation.Required),
		validation.Field(&u.Updates, validation.Required, validation.Length(1, 0),
			validation.Each(validation.By(func(value interface{}) error {
				update := value.(tableUpdate)
				return update.Validate()
			}))),
	)
}

type tableUpdate struct {
	UpdateType       string            `mapstructure:"updateType" json:"updateType"`
	RenameTable      *renameTable      `mapstructure:"renameTable" json:"renameTable"`
	RenameColumn     *renameColumn     `mapstructure:"renameColumn" json:"renameColumn"`
	AlterColumn      *alterColumn      `mapstructure:"alterColumn" json:"alterColumn"`
	AddColumn        *addColumn        `mapstructure:"addColumn" json:"addColumn"`
	RemoveColumn     *removeColumn     `mapstructure:"removeColumn" json:"removeColumn"`
	AddConstraint    *addConstraint    `mapstructure:"addConstraint" json:"addConstraint"`
	RemoveConstraint *removeConstraint `mapstructure:"removeConstraint" json:"removeConstraint"`
}

func (t *tableUpdate) Validate() error {
	return validation.ValidateStruct(t,
		validation.Field(&t.UpdateType, validation.Required, validation.In(
			RenameTableKey, RenameColumnKey, AlterColumnKey, AddColumnKey, RemoveColumnKey, AddConstraintKey, RemoveConstraintKey)),
		validation.Field(&t.RenameTable, validation.When(t.UpdateType == RenameTableKey,
			validation.Required, validation.WithContext(
				func(ctx context.Context, value interface{}) error {
					rename := value.(*renameTable)
					return rename.Validate()
				})).Else(validation.Nil)),
		validation.Field(&t.RenameColumn, validation.When(t.UpdateType == RenameColumnKey,
			validation.Required, validation.WithContext(
				func(ctx context.Context, value interface{}) error {
					rename := value.(*renameColumn)
					return rename.Validate()
				})).Else(validation.Nil)),
		validation.Field(&t.AlterColumn, validation.When(t.UpdateType == AlterColumnKey,
			validation.Required, validation.WithContext(
				func(ctx context.Context, value interface{}) error {
					alter := value.(*alterColumn)
					return alter.Validate()
				})).Else(validation.Nil)),
		validation.Field(&t.AddColumn, validation.When(t.UpdateType == AddColumnKey,
			validation.Required, validation.WithContext(
				func(ctx context.Context, value interface{}) error {
					add := value.(*addColumn)
					return add.Validate()
				})).Else(validation.Nil)),
		validation.Field(&t.RemoveColumn, validation.When(t.UpdateType == RemoveColumnKey,
			validation.Required, validation.WithContext(
				func(ctx context.Context, value interface{}) error {
					remove := value.(*removeColumn)
					return remove.Validate()
				})).Else(validation.Nil)),
		validation.Field(&t.AddConstraint, validation.When(t.UpdateType == AddConstraintKey,
			validation.Required, validation.WithContext(
				func(ctx context.Context, value interface{}) error {
					add := value.(*addConstraint)
					return add.Validate()
				})).Else(validation.Nil)),
		validation.Field(&t.RemoveConstraint, validation.When(t.UpdateType == RemoveConstraintKey,
			validation.Required, validation.WithContext(
				func(ctx context.Context, value interface{}) error {
					remove := value.(*removeConstraint)
					return remove.Validate()
				})).Else(validation.Nil)),
	)
}

type renameTable struct {
	Name string `mapstructure:"name" json:"name"`
}

func (a *renameTable) Validate() error {
	return validation.ValidateStruct(a,
		validation.Field(&a.Name, validation.Required),
	)
}

type renameColumn struct {
	OldName string `mapstructure:"oldName" json:"oldName"`
	NewName string `mapstructure:"newName" json:"newName"`
}

func (a *renameColumn) Validate() error {
	return validation.ValidateStruct(a,
		validation.Field(&a.OldName, validation.Required),
		validation.Field(&a.NewName, validation.Required, validation.NotIn(a.OldName)),
	)
}

type alterColumn struct {
	ColumnName   string `mapstructure:"columnName" json:"columnName"`
	DataType     string `mapstructure:"dataType" json:"dataType"`
	Nullable     bool   `mapstructure:"nullable" json:"nullable"`
	DefaultValue string `mapstructure:"defaultValue" json:"defaultValue"`
}

func (a *alterColumn) Validate() error {
	return validation.ValidateStruct(a,
		validation.Field(&a.ColumnName, validation.Required),
		validation.Field(&a.DataType, validation.Required),
		validation.Field(&a.Nullable),
		validation.Field(&a.DefaultValue),
	)
}

type addColumn struct {
	ColumnName   string `mapstructure:"columnName" json:"columnName"`
	DataType     string `mapstructure:"dataType" json:"dataType"`
	Nullable     bool   `mapstructure:"nullable" json:"nullable"`
	DefaultValue string `mapstructure:"defaultValue" json:"defaultValue"`
}

func (a *addColumn) Validate() error {
	return validation.ValidateStruct(a,
		validation.Field(&a.ColumnName, validation.Required),
		validation.Field(&a.DataType, validation.Required),
		validation.Field(&a.Nullable),
		validation.Field(&a.DefaultValue),
	)
}

type removeColumn struct {
	ColumnName string `mapstructure:"columnName" json:"columnName"`
}

func (r *removeColumn) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ColumnName, validation.Required),
	)
}

type addConstraint struct {
	ConstraintType  string `mapstructure:"constraintType" json:"constraintType"`
	ColumnName      string `mapstructure:"columnName" json:"columnName"`
	ReferencesTable string `mapstructure:"referencesTable" json:"referencesTable"`
	ReferencesField string `mapstructure:"referencesField" json:"referencesField"`
}

func (a *addConstraint) Validate() error {
	return validation.ValidateStruct(a,
		validation.Field(&a.ConstraintType, validation.Required, validation.In(
			PrimaryKeyConstraintKey, ForeignKeyConstraintKey, UniqueConstraintKey)),
		validation.Field(&a.ColumnName, validation.Required),
		validation.Field(&a.ReferencesTable),
		validation.Field(&a.ReferencesField),
	)
}

type removeConstraint struct {
	ConstraintName string `mapstructure:"constraintName" json:"constraintName"`
}

func (r *removeConstraint) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ConstraintName, validation.Required),
	)
}

func (c *Constraint) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.TableName, validation.Required),
		validation.Field(&c.ConstraintName, validation.Required),
		validation.Field(&c.ConstraintType, validation.Required, validation.In(
			PrimaryKeyConstraintKey, ForeignKeyConstraintKey, UniqueConstraintKey)),
		validation.Field(&c.ColumnName, validation.Required),
		validation.Field(&c.ReferencesTable),
		validation.Field(&c.ReferencesField),
	)
}

func (m *Model) Validate() error {
	return validation.ValidateStruct(m,
		validation.Field(&m.Name, validation.Required),
		validation.Field(&m.Table, validation.Required),
		validation.Field(&m.Projections, validation.Required, validation.Each(validation.By(
			func(value interface{}) error {
				if p, ok := value.(Projection); ok {
					return p.Validate(true)
				} else {
					return validation.NewError("invalid_by_cast", "could not cast in by clause")
				}

			}))),
		validation.Field(&m.PrimaryKey),
		validation.Field(&m.OwningAssociations, validation.Nil),
		validation.Field(&m.ReferencedAssociations, validation.Nil),
	)
}

func (c *Projection) Validate(dataType bool) error {
	return validation.ValidateStruct(c,
		validation.Field(&c.As, validation.Required),
		validation.Field(&c.Column, validation.Required),
		validation.Field(&c.SchemaType,
			validation.When(dataType, validation.Required, validation.In(
				common.DatabaseTypeString, common.DatabaseTypeNumber, common.DatabaseTypeBoolean,
			)).Else(validation.Empty)),
		validation.Field(&c.ModelType,
			validation.When(dataType, validation.Required, validation.In(
				common.DatabaseTypeString, common.DatabaseTypeNumber, common.DatabaseTypeBoolean,
			)).Else(validation.Empty)),
		validation.Field(&c.NotNull),
	)
}

func (a *ModelAssociation) Validate() error {
	return validation.ValidateStruct(a,
		validation.Field(&a.Name, validation.Required),
		validation.Field(&a.Type, validation.Required, validation.In(
			common.AssociationsBelongsTo, common.AssociationsBelongsToMany,
			common.AssociationsHasMany, common.AssociationsHasOne)),
		validation.Field(&a.TableName, validation.Required),
		validation.Field(&a.ColumnName, validation.Required),
		validation.Field(&a.ReferencesTable, validation.Required),
		validation.Field(&a.ReferencesField, validation.Required),
		validation.Field(&a.JoinTable,
			validation.When(a.Type == common.AssociationsBelongsToMany, validation.Required)),
		validation.Field(&a.JoinTableSourceField,
			validation.When(a.Type == common.AssociationsBelongsToMany, validation.Required)),
		validation.Field(&a.JoinTableTargetField,
			validation.When(a.Type == common.AssociationsBelongsToMany, validation.Required)),
		validation.Field(&a.OwningModelName, validation.Required),
		validation.Field(&a.ReferencesModelName, validation.Required),
	)
}

func (w *Where) Validate(cb func(val any) error) error {
	return validation.ValidateStruct(w,
		validation.Field(&w.Template, validation.NotNil),
		validation.Field(&w.Values, validation.Each(validation.By(
			func(value any) error {
				return cb(value)
			}))),
	)
}

func (p *Populate) Validate(cb func(val any) error) error {
	return validation.ValidateStruct(p,
		validation.Field(&p.Model, validation.Required),
		validation.Field(&p.As, validation.Required),
		validation.Field(&p.Required),
		validation.Field(&p.Project, validation.Each(validation.By(
			func(value any) error {
				v := value.(Projection)
				return v.Validate(false)
			})),
		),
		validation.Field(&p.Where, validation.Each(validation.By(
			func(value any) error {
				v := value.(Where)
				return v.Validate(cb)
			}))),
		validation.Field(&p.Populate, validation.Each(validation.By(
			func(value any) error {
				v := value.(Populate)
				return v.Validate(cb)
			}))),
	)
}
