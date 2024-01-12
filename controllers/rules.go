package controllers

import (
	"encoding/json"
	"fmt"
	"generic/config"
	dbutils "generic/dbUtils"
	"generic/models"
	"generic/schemas"
	"generic/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"rogchap.com/v8go"
)

func AddRuleToApi(c *gin.Context) {
	reqBodyAny, _ := c.Get("Req")
	reqBody := reqBodyAny.(*schemas.AddRuleToApi)

	_, found := dbutils.FindClient(reqBody.ClientId)
	if !found {
		utils.ResponseHandler(c, utils.Config{Response: utils.Responses["ClientNotFound"]})
		return
	}

	rulesCollection, ctx := config.GetMongoCollection("rules")
	insertResult, err := rulesCollection.InsertOne(ctx, reqBody.Rule)
	utils.HandleError(nil, err)

	api, found := dbutils.FindApi(reqBody.ClientId, reqBody.ApiName)
	if !found {
		utils.ResponseHandler(c, utils.Config{Response: utils.Responses["ApiNotFound"]})
		return
	}

	if reqBody.Rule.Type == "pre" {
		api.PreRules = append(api.PreRules, insertResult.InsertedID.(primitive.ObjectID))
	} else if reqBody.Rule.Type == "post" {
		api.PostRules = append(api.PostRules, insertResult.InsertedID.(primitive.ObjectID))
	}
	err = dbutils.UpdateApi(reqBody.ClientId, reqBody.ApiName, api)
	utils.HandleError(c, err)

	utils.ResponseHandler(c, utils.Config{Data: utils.ConvertToMap("inserted", insertResult.InsertedID)})
}

func RulesCall(c *gin.Context) {
	iso := v8go.NewIsolate()
	defer iso.Dispose()
	ctx := v8go.NewContext(iso)
	defer ctx.Close()

	_, ruleStr := models.GetRule(c.Param("id"))
	if ruleStr == "" {
		utils.ResponseHandler(c, utils.Config{Response: utils.Responses["NotFound"]})
		return
	}
	strReq, err := json.Marshal(c.Request.Body)
	utils.HandleError(c, err)
	scriptStr := utils.GetScriptString("./scripts/rulesProto.js")
	var result string

	utils.BenchmarkFn(func() {
		ctx.RunScript(scriptStr, "rulesProto.js")
		value, err := ctx.RunScript(fmt.Sprintf("parseAndEvaluate('%s','%s')", ruleStr, strReq), "rulesProto.js")
		utils.HandleError(c, err)
		result = value.String()
	})

	utils.ResponseHandler(c, utils.Config{Data: utils.ConvertToMap("Returned", result)})
}
