package controllers

import (
	"ifttt/manager/application/core"
	"ifttt/manager/application/middlewares"
	"ifttt/manager/domain/schema"
	"ifttt/manager/utils"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
)

type schemaController struct {
	serverCore *core.ServerCore
}

func newSchemaController(serverCore *core.ServerCore) *schemaController {
	return &schemaController{
		serverCore: serverCore,
	}
}

func (s *schemaController) GetSchema(c *gin.Context) {
	var schemas []schema.Schema

	tableNames, err := s.serverCore.DataStore.SchemaRepo.GetTableNames()
	if err != nil {
		utils.HandleErrorResponse(c, err)
		return
	}

	columns, err := s.serverCore.DataStore.SchemaRepo.GetAllColumns(tableNames)
	if err != nil {
		utils.HandleErrorResponse(c, err)
		return
	}

	constraints, err := s.serverCore.DataStore.SchemaRepo.GetAllConstraints(tableNames)
	if err != nil {
		utils.HandleErrorResponse(c, err)
		return
	}

	groupedColumns := lo.GroupBy(*columns, func(col schema.Column) string {
		return col.TableName
	})
	groupedConstraints := lo.GroupBy(*constraints, func(constraint schema.Constraint) string {
		return constraint.TableName
	})

	var newSchema schema.Schema
	for _, tableName := range tableNames {
		newSchema.TableName = tableName
		if columns, ok := groupedColumns[newSchema.TableName]; ok {
			newSchema.Columns = columns
		}
		if constraints, ok := groupedConstraints[newSchema.TableName]; ok {
			newSchema.Constraints = constraints
		}
		schemas = append(schemas, newSchema)
	}

	utils.ResponseHandler(c, utils.ResponseConfig{Data: schemas})

}

func (s *schemaController) CreateTable(c *gin.Context) {
	err, reqBodyAny := middlewares.Validator(c, schema.CreateTableRequest{})
	if err != nil {
		utils.HandleErrorResponse(c, err)
		return
	}
	reqBody := reqBodyAny.(*schema.CreateTableRequest)

	if err := s.serverCore.DataStore.SchemaRepo.CreateTable(reqBody); err != nil {
		utils.HandleErrorResponse(c, err)
		return
	}

	utils.ResponseHandler(c, utils.ResponseConfig{})
}

func (s *schemaController) UpdateTable(c *gin.Context) {
	err, reqBodyAny := middlewares.Validator(c, schema.UpdateTableRequest{})
	if err != nil {
		utils.HandleErrorResponse(c, err)
		return
	}
	reqBody := reqBodyAny.(*schema.UpdateTableRequest)

	if err := s.serverCore.DataStore.SchemaRepo.UpdateTable(reqBody); err != nil {
		utils.HandleErrorResponse(c, err)
		return
	}

	utils.ResponseHandler(c, utils.ResponseConfig{})
}
