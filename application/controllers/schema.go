package controllers

import (
	"ifttt/manager/application/core"
	"ifttt/manager/application/middlewares"
	"ifttt/manager/common"
	"ifttt/manager/domain/schema"

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
		common.HandleErrorResponse(c, err)
		return
	}

	columns, err := s.serverCore.DataStore.SchemaRepo.GetAllColumns(tableNames)
	if err != nil {
		common.HandleErrorResponse(c, err)
		return
	}

	constraints, err := s.serverCore.DataStore.SchemaRepo.GetAllConstraints(tableNames)
	if err != nil {
		common.HandleErrorResponse(c, err)
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

	common.ResponseHandler(c, common.ResponseConfig{Data: schemas})

}

func (s *schemaController) CreateTable(c *gin.Context) {
	err, reqBodyAny := middlewares.Validator(c, schema.CreateTableRequest{})
	if err != nil {
		common.HandleErrorResponse(c, err)
		return
	}
	reqBody := reqBodyAny.(*schema.CreateTableRequest)

	if existingTables, err := s.serverCore.DataStore.SchemaRepo.GetTableNames(); err != nil {
		common.HandleErrorResponse(c, err)
		return
	} else {
		if _, exists := lo.Find(existingTables, func(tName string) bool {
			return tName == reqBody.TableName
		}); exists {
			common.ResponseHandler(c, common.ResponseConfig{Response: common.Responses["TableAlreadyExists"]})
			return
		}
	}

	if err := s.serverCore.DataStore.SchemaRepo.CreateTable(reqBody); err != nil {
		common.HandleErrorResponse(c, err)
		return
	}

	common.ResponseHandler(c, common.ResponseConfig{})
}

func (s *schemaController) UpdateTable(c *gin.Context) {
	err, reqBodyAny := middlewares.Validator(c, schema.UpdateTableRequest{})
	if err != nil {
		common.HandleErrorResponse(c, err)
		return
	}
	reqBody := reqBodyAny.(*schema.UpdateTableRequest)

	if existingTables, err := s.serverCore.DataStore.SchemaRepo.GetTableNames(); err != nil {
		common.HandleErrorResponse(c, err)
		return
	} else {
		if _, exists := lo.Find(existingTables, func(tName string) bool {
			return tName == reqBody.TableName
		}); !exists {
			common.ResponseHandler(c, common.ResponseConfig{Response: common.Responses["TableNotFound"]})
			return
		}
	}

	if err := s.serverCore.DataStore.SchemaRepo.UpdateTable(reqBody); err != nil {
		common.HandleErrorResponse(c, err)
		return
	}

	common.ResponseHandler(c, common.ResponseConfig{})
}
