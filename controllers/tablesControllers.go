package controllers

import (
	"fmt"
	"generic/config"
	"generic/middlewares"
	"generic/models"
	"generic/schemas"
	"generic/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"github.com/samber/lo"
	"github.com/scylladb/gocqlx"
	"github.com/scylladb/gocqlx/qb"
	"github.com/scylladb/gocqlx/v2/table"
)

var Tables = struct {
	AddTable func(*gin.Context)
}{
	AddTable: func(c *gin.Context) {
		err, reqBodyAny := middlewares.Validator(c, schemas.AddTableRequest{})
		if err != nil {
			return
		}
		reqBody := reqBodyAny.(*schemas.AddTableRequest)

		partitionKeys := reqBody.PartitionKeys
		clusteringKeys := reqBody.ClusteringKeys
		allColumns := reqBody.AllColumns
		mappings := reqBody.Mappings

		// check that partition keys and clustering keys are independent
		if intersection := lo.Intersect(partitionKeys, clusteringKeys); len(intersection) != 0 {
			utils.ResponseHandler(c, utils.ResponseConfig{
				Response: utils.Responses["WrongTableFormat"],
				Data: utils.ConvertToMap("Error",
					fmt.Sprintf("Partition keys and clustering keys overlapping: %+v", intersection)),
			})
			return
		}

		// verify that all partition key are part of all columns
		if len(lo.Intersect(partitionKeys, allColumns)) != len(partitionKeys) {
			utils.ResponseHandler(c, utils.ResponseConfig{
				Response: utils.Responses["WrongTableFormat"],
				Data:     utils.ConvertToMap("Error", "All partition keys not included in field allColumns"),
			})
			return
		}

		// verify that all clustering columns are part of all columns
		if len(lo.Intersect(clusteringKeys, allColumns)) != len(clusteringKeys) {
			utils.ResponseHandler(c, utils.ResponseConfig{
				Response: utils.Responses["WrongTableFormat"],
				Data:     utils.ConvertToMap("Error", "All clustering keys not included in field allColumns"),
			})
			return
		}

		// verify that mapped columns are selected in all columns
		if mappedCols := lo.MapToSlice(mappings, func(_, value string) string {
			return value
		}); len(lo.Intersect(mappedCols, allColumns)) != len(mappedCols) {
			utils.ResponseHandler(c, utils.ResponseConfig{
				Response: utils.Responses["WrongTableFormat"],
				Data:     utils.ConvertToMap("Error", "All mapped keys not included in field allColumns"),
			})
			return
		}

		// check if table of this name already exists
		stmt, names := qb.Select("tables").Where(qb.Eq("name")).ToCql()
		q := config.GetScylla().Query(stmt, names).BindStruct(models.TablesModel{Name: reqBody.Name})
		var tables []models.TablesModel
		if err := gocqlx.Select(&tables, q.Query); err != nil {
			utils.HandleErrorResponse(c, err)
			return
		}
		if len(tables) > 0 {
			utils.ResponseHandler(c, utils.ResponseConfig{
				Response: utils.Responses["WrongTableFormat"],
				Data:     utils.ConvertToMap("Error", "Table with this name already exists"),
			})
			return
		}

		allowedSchemas := config.GetSchemas()
		allowedColumnsMap := []schemas.AllowedColumnsType{}
		if err := allowedSchemas.UnmarshalKey("allowedColumns", &allowedColumnsMap); err != nil {
			utils.HandleErrorResponse(c, err)
			return
		}

		// record table creation
		TablesTable := table.New(models.TablesMetadata)
		newTableName := fmt.Sprintf("table_%s_%s", reqBody.Name, strings.ReplaceAll(gocql.TimeUUID().String(), "-", "_"))[:40]
		newTable := models.TablesModel{
			InternalName:   newTableName,
			Name:           reqBody.Name,
			Description:    reqBody.Description,
			PartitionKeys:  reqBody.PartitionKeys,
			ClusteringKeys: reqBody.ClusteringKeys,
			AllColumns:     reqBody.AllColumns,
			Mappings:       reqBody.Mappings,
		}
		entry := config.GetScylla().Query(TablesTable.Insert()).BindStruct(&newTable)
		if err = entry.ExecRelease(); err != nil {
			utils.HandleErrorResponse(c, err)
			return
		}

		// create table
		newTableQuery := "created_at timestamp,"
		lo.ForEach(allowedColumnsMap, func(i schemas.AllowedColumnsType, idx int) {
			if lo.Contains(allColumns, i.Name) {
				newTableQuery += fmt.Sprintf("%s %s,", i.Name, i.DataType)
			}
		})
		newTableQuery = fmt.Sprintf("%sPRIMARY KEY ((%s),%s)", newTableQuery,
			strings.Join(reqBody.PartitionKeys, ","), strings.Join(reqBody.ClusteringKeys, ","))
		newTableQuery = fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s.%s (%s);",
			config.GetConfigProp("scylla.keyspace"), newTableName, newTableQuery)
		if err := config.GetScylla().ExecStmt(newTableQuery); err != nil {
			utils.HandleErrorResponse(c, err)
			return
		}

		utils.ResponseHandler(c, utils.ResponseConfig{})
	},
}
