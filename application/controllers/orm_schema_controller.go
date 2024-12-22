package controllers

import (
	"ifttt/manager/application/core"
	"ifttt/manager/common"
	"ifttt/manager/domain/orm_schema"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
)

type ormSchemaController struct {
	serverCore *core.ServerCore
}

func newOrmSchemaController(serverCore *core.ServerCore) *ormSchemaController {
	return &ormSchemaController{
		serverCore: serverCore,
	}
}

func (s *ormSchemaController) GetSchema(c *gin.Context) {
	var schemas []orm_schema.Schema

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

	groupedColumns := lo.GroupBy(*columns, func(col orm_schema.Column) string {
		return col.TableName
	})
	groupedConstraints := lo.GroupBy(*constraints, func(constraint orm_schema.Constraint) string {
		return constraint.TableName
	})

	var newSchema orm_schema.Schema
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

func (s *ormSchemaController) CreateTable(c *gin.Context) {
	var reqBody orm_schema.CreateTableRequest
	if ok := validateAndBind(c, &reqBody); !ok {
		return
	}

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

	if err := s.serverCore.DataStore.SchemaRepo.CreateTable(&reqBody); err != nil {
		common.HandleErrorResponse(c, err)
		return
	}

	common.ResponseHandler(c, common.ResponseConfig{})
}

func (s *ormSchemaController) UpdateTable(c *gin.Context) {
	var reqBody orm_schema.UpdateTableRequest
	if ok := validateAndBind(c, &reqBody); !ok {
		return
	}

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

	if err := s.serverCore.DataStore.SchemaRepo.UpdateTable(&reqBody); err != nil {
		common.HandleErrorResponse(c, err)
		return
	}

	common.ResponseHandler(c, common.ResponseConfig{})
}

func (s *ormSchemaController) CreateModel(c *gin.Context) {
	var reqBody orm_schema.Model
	if ok := validateAndBind(c, &reqBody); !ok {
		return
	}

	if model, err := s.serverCore.ConfigStore.OrmRepo.GetModelByName(reqBody.Name); err != nil {
		common.HandleErrorResponse(c, err)
		return
	} else if model != nil {
		common.ResponseHandler(c, common.ResponseConfig{Response: common.Responses["AlreadyExists"]})
		return
	}

	tableNames, err := s.serverCore.DataStore.SchemaRepo.GetTableNames()
	if err != nil {
		common.HandleErrorResponse(c, err)
		return
	} else if !lo.Contains(tableNames, reqBody.Name) {
		common.ResponseHandler(c, common.ResponseConfig{Response: common.Responses["NotFound"]})
	}

	tableColumns, err := s.serverCore.DataStore.SchemaRepo.GetAllColumns([]string{reqBody.Name})
	if err != nil {
		common.HandleErrorResponse(c, err)
		return
	} else {
		colNames := lo.Map(*tableColumns, func(col orm_schema.Column, _ int) string {
			return col.ColumnName
		})
		reqCols := lo.Map(reqBody.Projections, func(proj orm_schema.Projection, _ int) string {
			return proj.Column
		})
		if len(reqBody.Projections) != len(lo.Intersect(colNames, reqCols)) {
			common.ResponseHandler(c, common.ResponseConfig{Response: common.Responses["InvalidORM"]})
			return
		}
	}

	tableConstraints, err := s.serverCore.DataStore.SchemaRepo.GetAllConstraints([]string{reqBody.Name})
	if err != nil {
		common.HandleErrorResponse(c, err)
		return
	} else {
		var same bool
		for _, c := range *tableConstraints {
			if c.ColumnName == reqBody.PrimaryKey && c.ConstraintType == "PRIMARY KEY" {
				same = true
			}
		}
		if !same {
			common.ResponseHandler(c, common.ResponseConfig{Response: common.Responses["InvalidORM"]})
			return
		}
	}

	if err := s.serverCore.ConfigStore.OrmRepo.CreateModel(&reqBody); err != nil {
		common.HandleErrorResponse(c, err)
		return
	}

	common.ResponseHandler(c, common.ResponseConfig{})
}

func (s *ormSchemaController) CreateAssociation(c *gin.Context) {
	var reqBody orm_schema.ModelAssociation
	if ok := validateAndBind(c, &reqBody); !ok {
		return
	}

	if model, err := s.serverCore.ConfigStore.OrmRepo.GetAssociationByName(reqBody.Name); err != nil {
		common.HandleErrorResponse(c, err)
		return
	} else if model != nil {
		common.ResponseHandler(c, common.ResponseConfig{Response: common.Responses["AlreadyExists"]})
		return
	}

	if err := s.serverCore.ConfigStore.OrmRepo.CreateAssociation(&reqBody); err != nil {
		common.HandleErrorResponse(c, err)
		return
	}

	common.ResponseHandler(c, common.ResponseConfig{})

}
