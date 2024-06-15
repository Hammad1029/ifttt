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
	"github.com/samber/lo"
	"github.com/scylladb/gocqlx/v2/qb"
)

var Indexes = struct {
	AddIndex  func(*gin.Context)
	FindIndex func(*gin.Context)
	DropIndex func(*gin.Context)
}{
	AddIndex: func(c *gin.Context) {
		err, reqBodyAny := middlewares.Validator(c, schemas.AddIndexRequest{})
		if err != nil {
			return
		}
		reqBody := reqBodyAny.(*schemas.AddIndexRequest)

		// get table
		stmt, names := qb.Select("tables").Where(qb.Eq("name")).ToCql()
		q := config.GetScylla().Query(stmt, names).BindStruct(models.TablesModel{Name: reqBody.TableName})
		var tables []models.TablesModel
		if err := q.SelectRelease(&tables); err != nil {
			utils.HandleErrorResponse(c, err)
			return
		}
		if len(tables) == 0 {
			utils.ResponseHandler(c, utils.ResponseConfig{Response: utils.Responses["TableNotFound"]})
			return
		}

		selectedTable := tables[0]

		for _, col := range reqBody.IndexedColumns {
			if colModel, found := lo.Find(selectedTable.Columns, func(item models.TableColumn) bool {
				return item.Name == col
			}); found {
				if colModel.PartitionKey || colModel.ClusteringKey {
					utils.ResponseHandler(c, utils.ResponseConfig{
						Response: utils.Responses["IndexNotPossible"],
						Data:     utils.ConvertToMap("Error", "Can not create index on partition and clustering keys"),
					})
					return
				}
			} else {
				utils.ResponseHandler(c, utils.ResponseConfig{
					Response: utils.Responses["IndexNotPossible"],
					Data:     utils.ConvertToMap("Error", "Column not found in table"),
				})
				return
			}

		}

		indexName := utils.GenerateIndexName(selectedTable.InternalName, reqBody.IndexedColumns)

		if selectedTable.Indexes == nil {
			selectedTable.Indexes = make(map[string]models.IndexModel)
		}

		// check if index already exists
		if len(lo.PickByKeys(selectedTable.Indexes, []string{indexName})) != 0 {
			utils.ResponseHandler(c, utils.ResponseConfig{
				Response: utils.Responses["IndexNotPossible"],
				Data:     utils.ConvertToMap("Error", "Index already exists"),
			})
			return
		}

		selectedTable.Indexes[indexName] = models.IndexModel{
			Local:     reqBody.Local,
			IndexName: indexName,
			TableName: selectedTable.InternalName,
			Columns:   reqBody.IndexedColumns,
		}

		stmt, names = qb.Update("tables").Set("indexes").
			Where(qb.Eq("internal_name"), qb.Eq("name"), qb.Eq("description")).ToCql()
		if err := config.GetScylla().Query(stmt, names).BindStruct(&selectedTable).ExecRelease(); err != nil {
			utils.HandleErrorResponse(c, err)
			return
		}

		newIndexQuery := ""
		if reqBody.Local {
			// create local index
			newIndexQuery = fmt.Sprintf(
				"CREATE INDEX %s ON %s((%s),%s)",
				indexName,
				selectedTable.InternalName,
				strings.Join(
					lo.Map(
						lo.Filter(selectedTable.Columns,
							func(col models.TableColumn, _ int) bool {
								return col.PartitionKey == true
							},
						),
						func(col models.TableColumn, _ int) string {
							return col.Name
						}), ", "),
				strings.Join(reqBody.IndexedColumns, ", "),
			)
		} else {
			// create global index
			newIndexQuery = fmt.Sprintf(
				"CREATE INDEX %s ON %s(%s)",
				indexName,
				selectedTable.InternalName,
				strings.Join(reqBody.IndexedColumns, ", "),
			)
		}
		if err := config.GetScylla().ExecStmt(newIndexQuery); err != nil {
			utils.HandleErrorResponse(c, err)
			return
		}

		utils.ResponseHandler(c, utils.ResponseConfig{})
	},
	FindIndex: func(c *gin.Context) {
		err, reqBodyAny := middlewares.Validator(c, schemas.FindIndexRequest{})
		if err != nil {
			return
		}
		reqBody := reqBodyAny.(*schemas.FindIndexRequest)

		// check if table exists
		stmt, names := qb.Select("tables").Where(qb.Eq("name")).ToCql()
		q := config.GetScylla().Query(stmt, names).BindStruct(models.TablesModel{Name: reqBody.TableName})
		var tables []models.TablesModel
		if err := q.SelectRelease(&tables); err != nil {
			utils.HandleErrorResponse(c, err)
			return
		}
		if len(tables) == 0 {
			utils.ResponseHandler(c, utils.ResponseConfig{Response: utils.Responses["TableNotFound"]})
			return
		}

		// filter out all indexes containing columns
		filteredIndexes := lo.PickBy(tables[0].Indexes, func(key string, _ models.IndexModel) bool {
			return lo.EveryBy(reqBody.IndexedColumns, func(col string) bool {
				return strings.Contains(key, col)
			})
		})

		utils.ResponseHandler(c, utils.ResponseConfig{Data: utils.ConvertToMap("filteredIndexes", filteredIndexes)})
	},
	DropIndex: func(c *gin.Context) {
		err, reqBodyAny := middlewares.Validator(c, schemas.DropIndexRequest{})
		if err != nil {
			return
		}
		reqBody := reqBodyAny.(*schemas.DropIndexRequest)

		// check if table exists
		stmt, names := qb.Select("tables").Where(qb.Eq("name")).ToCql()
		q := config.GetScylla().Query(stmt, names).BindStruct(models.TablesModel{Name: reqBody.TableName})
		var tables []models.TablesModel
		if err := q.SelectRelease(&tables); err != nil {
			utils.HandleErrorResponse(c, err)
			return
		}
		if len(tables) == 0 {
			utils.ResponseHandler(c, utils.ResponseConfig{Response: utils.Responses["TableNotFound"]})
			return
		}

		selectedTable := tables[0]

		// check if index exists
		if _, ok := selectedTable.Indexes[reqBody.IndexName]; !ok {
			utils.ResponseHandler(c, utils.ResponseConfig{
				Response: utils.Responses["IndexNotFound"],
				Data:     utils.ConvertToMap("Error", "Index does not exist"),
			})
			return
		}

		// remove index from table details
		delete(selectedTable.Indexes, reqBody.IndexName)
		stmt, names = qb.Update("tables").Set("indexes").
			Where(qb.Eq("internal_name"), qb.Eq("name"), qb.Eq("description")).ToCql()
		if err := config.GetScylla().Query(stmt, names).BindStruct(&selectedTable).ExecRelease(); err != nil {
			utils.HandleErrorResponse(c, err)
			return
		}

		// drop index
		stmt = fmt.Sprintf("DROP INDEX %s", reqBody.IndexName)
		if err := config.GetScylla().ExecStmt(stmt); err != nil {
			utils.HandleErrorResponse(c, err)
			return
		}

		utils.ResponseHandler(c, utils.ResponseConfig{})
	},
}
