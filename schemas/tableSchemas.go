package schemas

import (
	"generic/config"
	"generic/utils"

	lo "github.com/samber/lo"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type AllowedColumnsType struct {
	Name     string `mapstructure:"name"`
	DataType string `mapstructure:"type"`
}

type AddTableRequest struct {
	Name           string            `json:"name"`
	Description    string            `json:"description"`
	PartitionKeys  []string          `json:"partitionKeys"`
	ClusteringKeys []string          `json:"clusteringKeys"`
	AllColumns     []string          `json:"allColumns"`
	Mappings       map[string]string `json:"mappings"`
}

func (r AddTableRequest) Validate() error {
	schemas := config.GetSchemas()
	allowedPartitionKeys := schemas.GetStringSlice("allowedPartitionKeys")
	allowedClusteringKeys := schemas.GetStringSlice("allowedClusteringKeys")
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
			&r.Name,
			validation.Required,
			validation.Length(3, 0),
		),
		validation.Field(
			&r.Description,
			validation.Required,
			validation.Length(3, 0),
		),
		validation.Field(
			&r.PartitionKeys,
			validation.Required,
			validation.Length(1, 3),
			validation.Each(validation.In(utils.ConvertStringToInterfaceArray(allowedPartitionKeys)...)),
		),
		validation.Field(
			&r.ClusteringKeys,
			validation.Required,
			validation.Length(0, 3),
			validation.Each(validation.In(utils.ConvertStringToInterfaceArray(allowedClusteringKeys)...)),
		),
		validation.Field(
			&r.AllColumns,
			validation.Required,
			validation.Length(1, len(allowedColumns)),
			validation.Each(validation.In(utils.ConvertStringToInterfaceArray(allowedColumns)...)),
		),
		validation.Field(
			&r.Mappings,
			validation.Required,
			validation.Each(validation.NilOrNotEmpty),
			validation.Each(validation.In(utils.ConvertStringToInterfaceArray(allowedColumns)...)),
		),
	)
}
