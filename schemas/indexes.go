package schemas

import (
	"generic/config"
	"generic/utils"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/samber/lo"
)

type AddIndex struct {
	Table   string   `json:"table"`
	Columns []string `json:"columns"`
	Local   bool     `json:"local"`
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
			&r.Table,
			validation.Required,
			validation.Length(3, 0),
		),
		validation.Field(
			&r.Columns,
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
