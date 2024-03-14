package schemas

import (
	"generic/config"
	"generic/utils"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/samber/lo"
)

type AddIndex struct {
	TableName      string   `json:"tableName"`
	IndexedColumns []string `json:"indexedColumns"`
	Local          bool     `json:"local"`
}

func (r AddIndex) Validate() error {
	schemas := config.GetSchemas()
	allowedColumnsMap := []AllowedColumnsType{}
	err := schemas.UnmarshalKey("allowedColumns", &allowedColumnsMap)
	if err != nil {
		utils.HandleError(err)
		return err
	}
	allowedColumns := lo.Map(allowedColumnsMap, func(col AllowedColumnsType, _ int) string {
		return col.Name
	})
	return validation.ValidateStruct(&r,
		validation.Field(
			&r.TableName,
			validation.Required,
			validation.Length(3, 0),
		),
		validation.Field(
			&r.IndexedColumns,
			validation.Required,
			validation.Length(1, len(allowedColumns)),
			validation.Each(validation.In(utils.ConvertStringToInterfaceArray(allowedColumns)...)),
		),
		validation.Field(
			&r.Local,
			validation.NotNil,
			validation.In(true, false),
		),
	)
}

type FindIndex struct {
	IndexedColumns []string `json:"indexedColumns"`
	TableName      string   `json:"tableName"`
}

func (r FindIndex) Validate() error {
	schemas := config.GetSchemas()
	allowedColumnsMap := []AllowedColumnsType{}
	err := schemas.UnmarshalKey("allowedColumns", &allowedColumnsMap)
	if err != nil {
		utils.HandleError(err)
		return err
	}
	allowedColumns := lo.Map(allowedColumnsMap, func(col AllowedColumnsType, _ int) string {
		return col.Name
	})
	return validation.ValidateStruct(&r,
		validation.Field(
			&r.TableName,
			validation.Required,
			validation.Length(3, 0),
		),
		validation.Field(
			&r.IndexedColumns,
			validation.Length(0, 3),
			validation.Each(validation.In(utils.ConvertStringToInterfaceArray(allowedColumns)...)),
		),
	)
}

type DropIndex struct {
	IndexName string `json:"indexName"`
	TableName string `json:"tableName"`
}

func (r DropIndex) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(
			&r.IndexName,
			validation.Required,
			validation.Length(3, 0),
		),
		validation.Field(
			&r.TableName,
			validation.Required,
			validation.Length(3, 0),
		),
	)
}
