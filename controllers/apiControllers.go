package controllers

import (
	"generic/config"
	"generic/middlewares"
	"generic/models"
	"generic/schemas"
	"generic/utils"

	"github.com/gin-gonic/gin"
	"github.com/scylladb/gocqlx"
	"github.com/scylladb/gocqlx/qb"
	"github.com/scylladb/gocqlx/v2/table"
)

var Apis = struct {
	AddApi func(*gin.Context)
}{
	AddApi: func(c *gin.Context) {
		err, reqBodyAny := middlewares.Validator(c, schemas.AddApiRequest{})
		if err != nil {
			return
		}
		reqBody := reqBodyAny.(*schemas.AddApiRequest)

		var apis []models.ApiModel

		// check if api of this name already exists
		stmt, names := qb.Select("apis").Where(qb.Eq("api_name")).ToCql()
		q := config.GetScylla().Query(stmt, names).BindStruct(models.ApiModel{ApiName: reqBody.ApiName})
		if err := gocqlx.Select(&apis, q.Query); err != nil {
			utils.HandleErrorResponse(c, err)
			return
		}
		if len(apis) > 0 {
			utils.ResponseHandler(c, utils.ResponseConfig{Response: utils.Responses["ApiAlreadyExists"]})
			return
		}

		// check if api of this path already exists
		stmt, names = qb.Select("apis").Where(qb.Eq("api_path")).ToCql()
		q = config.GetScylla().Query(stmt, names).BindStruct(models.ApiModel{ApiPath: reqBody.ApiPath})
		if err := gocqlx.Select(&apis, q.Query); err != nil {
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
			StartRules:     reqBody.StartRules,
			Rules:          reqBody.Rules,
		}

		// generate parameterized queries & stringify resolvable data
		api.TransformApiForSave()

		// insert api
		ApisTable := table.New(models.ApisMetadata)
		entry := config.GetScylla().Query(ApisTable.Insert()).BindStruct(&api)
		if err = entry.ExecRelease(); err != nil {
			utils.HandleErrorResponse(c, err)
			return
		}

		utils.ResponseHandler(c, utils.ResponseConfig{})
	},
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
