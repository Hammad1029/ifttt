package controllers

import (
	"generic/config"
	"generic/middlewares"
	"generic/models"
	"generic/schemas"
	"generic/utils"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	"github.com/scylladb/gocqlx"
	"github.com/scylladb/gocqlx/qb"
)

var Indexes = struct {
	AddIndex func(*gin.Context)
}{
	AddIndex: func(c *gin.Context) {
		err, reqBodyAny := middlewares.Validator(c, schemas.AddIndex{})
		if err != nil {
			return
		}
		reqBody := reqBodyAny.(*schemas.AddIndex)

		// check if table of this name already exists
		stmt, names := qb.Select("tables").Where(qb.Eq("name")).ToCql()
		q := config.GetScylla().Query(stmt, names).BindStruct(models.TablesModel{Name: reqBody.Table})
		var tables []models.TablesModel
		if err := gocqlx.Select(&tables, q.Query); err != nil {
			utils.HandleErrorResponse(c, err)
			return
		}
		if len(tables) == 0 {
			utils.ResponseHandler(c, utils.ResponseConfig{Response: utils.Responses["TableNotFound"]})
			return
		}

		selectedTable := tables[0]

		// check if column to be indexed is a partition key
		if len(lo.Intersect(selectedTable.PartitionKeys, reqBody.Columns)) != 0 {
			utils.ResponseHandler(c, utils.ResponseConfig{
				Response: utils.Responses["IndexNotPossible"],
				Data:     utils.ConvertToMap("Error", "Can not create index on partition keys"),
			})
			return
		}

		// check if column to be indexed is a clustering key
		if len(lo.Intersect(selectedTable.ClusteringKeys, reqBody.Columns)) != 0 {
			utils.ResponseHandler(c, utils.ResponseConfig{
				Response: utils.Responses["IndexNotPossible"],
				Data:     utils.ConvertToMap("Error", "Can not create index on clustering keys"),
			})
			return
		}

		// check if table has these columns
		if len(lo.Intersect(selectedTable.AllColumns, reqBody.Columns)) != len(reqBody.Columns) {
			utils.ResponseHandler(c, utils.ResponseConfig{
				Response: utils.Responses["IndexNotPossible"],
				Data:     utils.ConvertToMap("Error", "All of the specified columns are not activated for this table"),
			})
			return
		}

		indexName := utils.GenerateIndexName(selectedTable.InternalName, reqBody.Columns)

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
			Name:      indexName,
			TableName: selectedTable.InternalName,
			Columns:   reqBody.Columns,
		}

		stmt, names = qb.Update("tables").Set("indexes").
			Where(qb.Eq("internal_name"), qb.Eq("name"), qb.Eq("description")).ToCql()
		if err := config.GetScylla().Query(stmt, names).BindStruct(&selectedTable).ExecRelease(); err != nil {
			utils.HandleErrorResponse(c, err)
			return
		}

		// newIndexQuery := ""
		// if reqBody.Local {
		// 	// create local index
		// 	newIndexQuery = fmt.Sprintf(
		// 		"CREATE INDEX %s ON %s((%s),%s)",
		// 		indexName,
		// 		selectedTable.InternalName,
		// 		strings.Join(selectedTable.PartitionKeys, ", "),
		// 		strings.Join(reqBody.Columns, ", "),
		// 	)
		// } else {
		// 	// create global index
		// 	newIndexQuery = fmt.Sprintf(
		// 		"CREATE INDEX %s ON %s(%s)",
		// 		indexName,
		// 		selectedTable.InternalName,
		// 		strings.Join(reqBody.Columns, ", "),
		// 	)
		// }
		// if err := config.GetScylla().ExecStmt(newIndexQuery); err != nil {
		// 	utils.HandleErrorResponse(c, err)
		// 	return
		// }

		utils.ResponseHandler(c, utils.ResponseConfig{})
	},
}
