package controllers

import (
	"generic/config"
	"generic/middlewares"
	"generic/models"
	"generic/schemas"
	"generic/utils"

	"github.com/gin-gonic/gin"
	"github.com/scylladb/gocqlx/v2/qb"
	"github.com/scylladb/gocqlx/v2/table"
)

var Apis = struct {
	AddApi  func(*gin.Context)
	GetApis func(*gin.Context)
}{
	AddApi: func(c *gin.Context) {
		err, reqBodyAny := middlewares.Validator(c, schemas.AddApiRequest{})
		if err != nil {
			return
		}
		reqBody := reqBodyAny.(*schemas.AddApiRequest)

		var apis []models.ApiModel

		// check if api of this name already exists
		stmt, names := qb.Select("apis").Where(qb.Eq("api_group"), qb.Eq("api_name")).ToCql()
		q := config.GetScylla().Query(stmt, names).BindStruct(models.ApiModelSerialized{ApiName: reqBody.ApiName, ApiGroup: reqBody.ApiGroup})
		if err := q.SelectRelease(&apis); err != nil {
			utils.HandleErrorResponse(c, err)
			return
		}
		if len(apis) > 0 {
			utils.ResponseHandler(c, utils.ResponseConfig{Response: utils.Responses["ApiAlreadyExists"]})
			return
		}

		// create api struct
		api := models.ApiModel{
			ApiGroup:       reqBody.ApiGroup,
			ApiName:        reqBody.ApiName,
			ApiPath:        reqBody.ApiPath,
			ApiDescription: reqBody.ApiDescription,
			ApiRequest:     reqBody.ApiRequest,
			StartRules:     reqBody.StartRules,
			Rules:          reqBody.Rules,
		}

		// generate parameterized queries & serialize data
		apiSerialized, err := api.TransformApiForSave()
		if err != nil {
			utils.HandleErrorResponse(c, err)
			return
		}

		// insert api
		ApisTable := table.New(models.ApisMetadata)
		entry := config.GetScylla().Query(ApisTable.Insert()).BindStruct(&apiSerialized)
		if err = entry.ExecRelease(); err != nil {
			utils.HandleErrorResponse(c, err)
			return
		}

		utils.ResponseHandler(c, utils.ResponseConfig{})
	},
	GetApis: func(c *gin.Context) {
		err, reqBodyAny := middlewares.Validator(c, schemas.GetApisRequest{})
		if err != nil {
			return
		}
		reqBody := reqBodyAny.(*schemas.GetApisRequest)

		var conditions []qb.Cmp
		switch {
		case reqBody.ApiGroup != "":
			conditions = append(conditions, qb.Eq("api_group"))
		case reqBody.ApiName != "":
			conditions = append(conditions, qb.Eq("api_name"))
		case reqBody.ApiDescription != "":
			conditions = append(conditions, qb.Eq("api_description"))
		}

		var apis []models.ApiModel
		stmt, names := qb.Select("apis").Where(conditions...).ToCql()
		q := config.GetScylla().Query(stmt, names).BindStruct(reqBody)
		if err := q.SelectRelease(&apis); err != nil {
			utils.HandleErrorResponse(c, err)
			return
		}

		utils.ResponseHandler(c, utils.ResponseConfig{Data: apis})
	},
}
