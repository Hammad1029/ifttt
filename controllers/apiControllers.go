package controllers

import (
	"github.com/gin-gonic/gin"
)

func AddApi(c *gin.Context) {
	// reqBodyAny, _ := c.Get("Req")
	// reqBody := reqBodyAny.(*schemas.AddApi)

	// // check if client exists
	// client, found := dbutils.FindClient(reqBody.ClientId)
	// if !found {
	// 	utils.ResponseHandler(c, utils.ResponseConfig{Response: utils.Responses["ClientNotFound"]})
	// 	return
	// }

	// _, found = dbutils.FindApi(reqBody.ClientId, reqBody.ApiName)
	// if found {
	// 	utils.ResponseHandler(c, utils.ResponseConfig{Response: utils.Responses["ApiAlreadyExists"]})
	// 	return
	// }

	// // add api to mongo
	// var api models.ApisMongo
	// api.ApiName = reqBody.ApiName
	// api.ApiPath = reqBody.PathName
	// api.ClientId = reqBody.ClientId
	// apiCollection, ctx := config.GetMongoCollection("apis")
	// insertResult, err := apiCollection.InsertOne(ctx, api)
	// utils.HandleError(c, err)

	// // add api _id to client
	// client.Apis = append(client.Apis, insertResult.InsertedID.(primitive.ObjectID))
	// err = dbutils.UpdateClient(reqBody.ClientId, client)
	// utils.HandleError(c, err)

	// utils.ResponseHandler(c, utils.ResponseConfig{Data: utils.ConvertToMap("inserted", insertResult)})
}

func AddMappingToApi(c *gin.Context) {
	// reqBodyAny, _ := c.Get("Req")
	// reqBody := reqBodyAny.(*schemas.AddMappingToApi)

	// client, found := dbutils.FindClient(reqBody.ClientId)
	// if !found {
	// 	utils.ResponseHandler(c, utils.ResponseConfig{Response: utils.Responses["ClientNotFound"]})
	// 	return
	// }

	// _, found = dbutils.FindApi(reqBody.ClientId, reqBody.ApiName)
	// if !found {
	// 	utils.ResponseHandler(c, utils.ResponseConfig{Response: utils.Responses["ApiNotFound"]})
	// 	return
	// }

	// if len(client.VMap) == 0 {
	// 	client.VMap = make(map[string]map[string]string)
	// }
	// client.VMap[reqBody.ApiName] = reqBody.Mappings
	// dbutils.UpdateClient(reqBody.ClientId, client)
	// utils.ResponseHandler(c, utils.ResponseConfig{})
}

func CallApi(c *gin.Context) {
	// utils.BenchmarkFn(func() {
	// 	L := lua.NewState()
	// 	defer L.Close()

	// 	reqBody := make(map[string]string)
	// 	c.ShouldBindJSON(&reqBody)

	// 	// get from redis
	// 	ruleStr := `{"name":"createCard","type":"pre","rule":{"op1":"true","operator":"eq","op2":"true","thenActions":[{"type":"modifyReq","data":{"field":"currentBalance","value":{"get":"initialBalance","from":"req"}}}]}}`
	// 	apiMappingStr := `{"cardNumber":"String1","initialBalance":"int1","currentBalance":"int2"}`
	// 	allMappingsStr := `{"createCard":{"cardNumber":"String1","initialBalance":"int1","currentBalance":"int2"},"addFunds":{"cardNumber":"String1","amount":"int1"}}`
	// 	reqBArr, err := json.Marshal(c.Request.Body)
	// 	reqJson := lua.LString(string(reqBArr[:]))
	// 	utils.HandleError(c, err)

	// 	L.SetGlobal("ruleJson", lua.LString(ruleStr))
	// 	L.SetGlobal("apiMapping", lua.LString(apiMappingStr))
	// 	L.SetGlobal("allMappings", lua.LString(allMappingsStr))
	// 	L.SetGlobal("reqJson", lua.LString(reqJson))

	// 	// vmap, found := dbutils.GetClientVMap(c.Param("clientId"))
	// 	// apiVMap := vmap[c.Param("pathName")]
	// 	// if found && len(apiVMap) != 0 {
	// 	// 	resolvedReq := make(map[string]string)
	// 	// 	for col, value := range reqBody {
	// 	// 		internalCol, ok := apiVMap[col]
	// 	// 		if ok {
	// 	// 			resolvedReq[internalCol] = value
	// 	// 		} else {
	// 	// 			resolvedReq[col] = value
	// 	// 		}
	// 	// 	}

	// 	// 	reqBody = resolvedReq
	// 	// }

	// 	// var msg models.MsgKafka
	// 	// msg.ApiName = api.ApiName
	// 	// msg.ClientId = api.ClientId
	// 	// msg.Timestamp = time.Now()
	// 	// msg.Data = reqBody

	// 	err = L.DoFile("scripts/evaluate.lua")
	// 	utils.HandleError(c, err)
	// 	//execute pre reles
	// 	// config.CreateKafkaMsg(msg)
	// 	//execute post rules
	// })
}
