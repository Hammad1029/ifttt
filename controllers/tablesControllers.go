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
	"github.com/mitchellh/mapstructure"
	"github.com/samber/lo"
	"github.com/scylladb/gocqlx/v2/qb"
	"github.com/scylladb/gocqlx/v2/table"
)

var Tables = struct {
	AddTable    func(*gin.Context)
	GetTables   func(*gin.Context)
	UpdateTable func(*gin.Context)
}{
	AddTable: func(c *gin.Context) {
		err, reqBodyAny := middlewares.Validator(c, schemas.AddTableRequest{})

		if err != nil {
			return
		}
		reqBody := reqBodyAny.(*schemas.AddTableRequest)

		// check if table of this name already exists
		stmt, names := qb.Select("tables").Where(qb.Eq("name")).ToCql()
		q := config.GetScylla().Query(stmt, names).BindStruct(models.TablesModel{Name: reqBody.Name})
		var tables []models.TablesModel
		if err := q.SelectRelease(&tables); err != nil {
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

		newTableQuery := "created_at timestamp,"
		allowedDataTypes := config.GetSchemas().GetStringSlice("allowedDataTypes")

		var partitionKeys []string
		var clusteringKeys []string

		// validate columns
		for _, col := range reqBody.Columns {
			if col.ClusteringKey && col.PartitionKey {
				utils.ResponseHandler(c, utils.ResponseConfig{
					Response: utils.Responses["WrongTableFormat"],
					Data:     utils.ConvertToMap("Error", "Column cannot be both partition and clustering key"),
				})
				return
			} else if !lo.Contains(allowedDataTypes, col.DataType) {
				utils.ResponseHandler(c, utils.ResponseConfig{
					Response: utils.Responses["WrongTableFormat"],
					Data:     utils.ConvertToMap("Error", fmt.Sprintf("Wrong column dataType: %s", col.DataType)),
				})
				return
			} else {
				newTableQuery += fmt.Sprintf("%s %s,", col.Name, col.DataType)
				if col.PartitionKey {
					partitionKeys = append(partitionKeys, col.Name)
				} else if col.ClusteringKey {
					clusteringKeys = append(clusteringKeys, col.Name)
				}
			}
		}

		// record table creation
		TablesTable := table.New(models.TablesMetadata)
		newTableName := fmt.Sprintf("table_%s_%s", reqBody.Name, strings.ReplaceAll(gocql.TimeUUID().String(), "-", "_"))[:40]
		newTable := models.TablesModel{
			InternalName: newTableName,
			Name:         reqBody.Name,
			Description:  reqBody.Description,
			Columns:      reqBody.Columns,
		}
		entry := config.GetScylla().Query(TablesTable.Insert()).BindStruct(&newTable)
		if err = entry.ExecRelease(); err != nil {
			utils.HandleErrorResponse(c, err)
			return
		}

		// create table
		clusteringKeysQueryChunk := ""
		if len(clusteringKeys) != 0 {
			clusteringKeysQueryChunk = fmt.Sprintf(", %s", strings.Join(clusteringKeys, ","))
		}
		newTableQuery = fmt.Sprintf("%sPRIMARY KEY ((%s)%s)", newTableQuery,
			strings.Join(partitionKeys, ","), clusteringKeysQueryChunk)
		newTableQuery = fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s.%s (%s);",
			config.GetConfigProp("scylla.keyspace"), newTableName, newTableQuery)
		if err := config.GetScylla().ExecStmt(newTableQuery); err != nil {
			utils.HandleErrorResponse(c, err)
			return
		}

		utils.ResponseHandler(c, utils.ResponseConfig{})
	},
	GetTables: func(c *gin.Context) {
		err, reqBodyAny := middlewares.Validator(c, schemas.GetTablesRequest{})
		if err != nil {
			return
		}
		reqBody := reqBodyAny.(*schemas.GetTablesRequest)

		tableModel := models.TablesModel{}
		err = mapstructure.Decode(reqBody, &tableModel)
		if err != nil {
			utils.HandleErrorResponse(c, err)
			return
		}

		var queryParams []qb.Cmp
		if tableModel.InternalName != "" {
			queryParams = append(queryParams, qb.Eq("internal_name"))
		}

		stmt, names := qb.Select("tables").Where(queryParams...).ToCql()

		q := config.GetScylla().Query(stmt, names).BindStruct(tableModel)
		var tables []models.TablesModel
		if err := q.SelectRelease(&tables); err != nil {
			utils.HandleErrorResponse(c, err)
			return
		}

		filteredTables := lo.Filter(tables, func(item models.TablesModel, _ int) bool {
			return strings.Contains(item.Name, tableModel.Name) || strings.Contains(item.Description, tableModel.Description)
		})

		utils.ResponseHandler(c, utils.ResponseConfig{Data: filteredTables})
	},
	UpdateTable: func(c *gin.Context) {
		err, reqBodyAny := middlewares.Validator(c, schemas.UpdateTableRequest{})
		if err != nil {
			return
		}
		reqBody := reqBodyAny.(*schemas.UpdateTableRequest)

		// check if table of this name already exists
		stmt, names := qb.Select("tables").Where(qb.Eq("internal_name")).ToCql()
		q := config.GetScylla().Query(stmt, names).BindStruct(models.TablesModel{InternalName: reqBody.InternalName})
		var tables []models.TablesModel
		if err := q.SelectRelease(&tables); err != nil {
			utils.HandleErrorResponse(c, err)
			return
		}
		if len(tables) == 0 {
			utils.ResponseHandler(c, utils.ResponseConfig{
				Response: utils.Responses["WrongTableFormat"],
				Data:     utils.ConvertToMap("Error", "Table with this name does not exists"),
			})
			return
		}

		selectedTable := tables[0]
		selectedTable.Description = reqBody.Description
		allowedDataTypes := config.GetSchemas().GetStringSlice("allowedDataTypes")
		keyspaceName := config.GetConfigProp("scylla.keyspace")
		tableName := selectedTable.Name
		alterQueries := map[string][]string{
			"add":    {fmt.Sprintf("ALTER TABLE %s.%s ADD (", keyspaceName, tableName)},
			"alter":  {},
			"drop":   {fmt.Sprintf("ALTER TABLE %s.%s DROP ", keyspaceName, tableName)},
			"rename": {},
		}

		for _, col := range reqBody.Columns {
			if col.Add && lo.Contains(allowedDataTypes, col.DataType) {
				alterQueries["add"][0] += fmt.Sprintf("%s %s ", col.Name, col.DataType)
			} else if col.Alter && lo.Contains(allowedDataTypes, col.DataType) {
				alterQueries["alter"] = append(alterQueries["alter"],
					fmt.Sprintf("ALTER TABLE %s.%s ALTER %s TYPE %s;", keyspaceName, tableName, col.Name, col.DataType))
			} else if col.Drop {
				alterQueries["drop"][0] += fmt.Sprintf("%s AND ", col.Name)
			} else if col.Rename != "" {
				alterQueries["rename"] = append(alterQueries["rename"],
					fmt.Sprintf("ALTER TABLE %s.%s RENAME %s TO %s;", keyspaceName, tableName, col.Name, col.Rename))
			} else {
				utils.ResponseHandler(c, utils.ResponseConfig{
					Response: utils.Responses["WrongTableFormat"],
					Data:     utils.ConvertToMap("Error", fmt.Sprintf("Invalid column data: %v", col)),
				})
				return
			}
		}

		alterQueries["add"][0] += ");"
		alterQueries["drop"][0], _ = strings.CutSuffix(alterQueries["drop"][0], " AND ")
		alterQueries["drop"][0] += ";"

		batchQuery := fmt.Sprintf("BEGIN BATCH \n %s \n %s \n %s \n %s \n APPLY BATCH;",
			strings.Join(alterQueries["add"], "\n"),
			strings.Join(alterQueries["alter"], "\n"),
			strings.Join(alterQueries["drop"], "\n"),
			strings.Join(alterQueries["rename"], "\n"))

		if err := config.GetScylla().ExecStmt(batchQuery); err != nil {
			utils.HandleErrorResponse(c, err)
			return
		}

	},
}
