package controllers

import (
	"github.com/gin-gonic/gin"
)

func AddRuleToApi(c *gin.Context) {
	// reqBodyAny, _ := c.Get("Req")
	// reqBody := reqBodyAny.(*schemas.AddRuleToApi)

	// _, found := dbutils.FindClient(reqBody.ClientId)
	// if !found {
	// 	utils.ResponseHandler(c, utils.ResponseConfig{Response: utils.Responses["ClientNotFound"]})
	// 	return
	// }

	// rulesCollection, ctx := config.GetMongoCollection("rules")
	// insertResult, err := rulesCollection.InsertOne(ctx, reqBody.Rule)
	// utils.HandleError(nil, err)

	// api, found := dbutils.FindApi(reqBody.ClientId, reqBody.ApiName)
	// if !found {
	// 	utils.ResponseHandler(c, utils.ResponseConfig{Response: utils.Responses["ApiNotFound"]})
	// 	return
	// }

	// if reqBody.Rule.Type == "pre" {
	// 	api.PreRules = append(api.PreRules, insertResult.InsertedID.(primitive.ObjectID))
	// } else if reqBody.Rule.Type == "post" {
	// 	api.PostRules = append(api.PostRules, insertResult.InsertedID.(primitive.ObjectID))
	// }
	// err = dbutils.UpdateApi(reqBody.ClientId, reqBody.ApiName, api)
	// utils.HandleError(c, err)

	// utils.ResponseHandler(c, utils.ResponseConfig{Data: utils.ConvertToMap("inserted", insertResult.InsertedID)})
}

func RulesCall(c *gin.Context) {
	// L := lua.NewState()
	// defer L.Close()

	// rule := lua.LString(`{"name":"createCard","type":"pre","rule":{"op1":"true","operator":"eq","op2":"true","thenActions":[{"type":"modifyReq","data":{"field":"currentBalance","value":{"get":"initialBalance","from":"req"}}}]}}`)
	// apiMapping := lua.LString(`{"cardNumber":"String1","initialBalance":"int1","currentBalance":"int2"}`)
	// allMappings := lua.LString(`{"createCard":{"cardNumber":"String1","initialBalance":"int1","currentBalance":"int2"},"addFunds":{"cardNumber":"String1","amount":"int1"}}`)
	// reqBArr, err := json.Marshal(c.Request.Body)
	// reqJson := lua.LString(string(reqBArr[:]))
	// utils.HandleError(c, err)

	// L.SetGlobal("ruleJson", rule)
	// L.SetGlobal("apiMapping", apiMapping)
	// L.SetGlobal("allMappings", allMappings)
	// L.SetGlobal("reqJson", reqJson)
	// err = L.DoFile("scripts/evaluate.lua")
	// utils.HandleError(c, err)

	// _, ruleStr := models.GetRule(c.Param("id"))
	// if ruleStr == "" {
	// 	utils.ResponseHandler(c, utils.Config{Response: utils.Responses["NotFound"]})
	// 	return
	// }
	// strReq, err := json.Marshal(c.Request.Body)
	// utils.HandleError(c, err)
	// scriptStr := utils.GetScriptString("./scripts/rulesProto.js")
	// var result string

	// utils.BenchmarkFn(func() {
	// 	ctx.RunScript(scriptStr, "rulesProto.js")
	// 	value, err := ctx.RunScript(fmt.Sprintf("parseAndEvaluate('%s','%s')", ruleStr, strReq), "rulesProto.js")
	// 	utils.HandleError(c, err)
	// 	result = value.String()
	// })

	// utils.ResponseHandler(c, utils.Config{Data: utils.ConvertToMap("Returned", result)})
}
