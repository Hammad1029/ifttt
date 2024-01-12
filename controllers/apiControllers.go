package controllers

import (
	"generic/config"
	dbutils "generic/dbUtils"
	"generic/models"
	"generic/schemas"
	"generic/utils"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddApi(c *gin.Context) {
	reqBodyAny, _ := c.Get("Req")
	reqBody := reqBodyAny.(*schemas.AddApi)

	// check if client exists
	client, found := dbutils.FindClient(reqBody.ClientId)
	if !found {
		utils.ResponseHandler(c, utils.Config{Response: utils.Responses["ClientNotFound"]})
		return
	}

	_, found = dbutils.FindApi(reqBody.ClientId, reqBody.ApiName)
	if found {
		utils.ResponseHandler(c, utils.Config{Response: utils.Responses["ApiAlreadyExists"]})
		return
	}

	// add api to mongo
	var api models.ApisMongo
	api.ApiName = reqBody.ApiName
	api.ApiPath = reqBody.PathName
	api.ClientId = reqBody.ClientId
	apiCollection, ctx := config.GetMongoCollection("apis")
	insertResult, err := apiCollection.InsertOne(ctx, api)
	utils.HandleError(c, err)

	// add api _id to client
	client.Apis = append(client.Apis, insertResult.InsertedID.(primitive.ObjectID))
	err = dbutils.UpdateClient(reqBody.ClientId, client)
	utils.HandleError(c, err)

	utils.ResponseHandler(c, utils.Config{Data: utils.ConvertToMap("inserted", insertResult)})
}

func AddMappingToApi(c *gin.Context) {
	reqBodyAny, _ := c.Get("Req")
	reqBody := reqBodyAny.(*schemas.AddMappingToApi)

	api, found := dbutils.FindApi(reqBody.ClientId, reqBody.ApiName)
	if !found {
		utils.ResponseHandler(c, utils.Config{Response: utils.Responses["ApiAlreadyExists"]})
		return
	}

	api.VMap = reqBody.Mappings
	dbutils.UpdateApi(reqBody.ClientId, reqBody.ApiName, api)
	utils.ResponseHandler(c, utils.Config{})
}

func CallApi(c *gin.Context) {
	utils.BenchmarkFn(func() {

		reqBody := make(map[string]string)
		c.ShouldBindJSON(&reqBody)
		var api models.ApisMongo
		apisCollection, ctx := config.GetMongoCollection("apis")
		queryClient := bson.M{"clientId": c.Param("clientId"), "apiName": c.Param("pathName")}
		apiFind := apisCollection.FindOne(ctx, queryClient)
		if apiFind.Err() != nil {
			utils.ResponseHandler(c, utils.Config{Response: utils.Responses["NotFound"]})
			return
		}
		apiFind.Decode(&api)

		if len(api.VMap) != 0 {
			resolvedReq := make(map[string]string)
			for col, value := range reqBody {
				internalCol, ok := api.VMap[col]
				if ok {
					resolvedReq[internalCol] = value
				} else {
					resolvedReq[col] = value
				}
			}

			reqBody = resolvedReq
		}

		var msg models.MsgKafka
		msg.ApiName = api.ApiName
		msg.ClientId = api.ClientId
		msg.Timestamp = time.Now()
		msg.Data = reqBody

		//execute pre reles
		config.CreateKafkaMsg(msg)
		//execute post rules
	})
}
