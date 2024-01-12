package models

import (
	"generic/config"
	"generic/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ApisMongo struct {
	ApiName   string               `json:"apiName"`
	ClientId  string               `json:"clientId"`
	ApiPath   string               `json:"pathName"`
	PreRules  []primitive.ObjectID `json:"preRules"`
	PostRules []primitive.ObjectID `json:"postRules"`
	VMap      map[string]string    `json:"virtualMap"`
}

func AddApi(c *gin.Context) *mongo.InsertOneResult {
	rulesCollection, ctx := config.GetMongoCollection("apis")
	var api ApisMongo
	err := c.BindJSON(api)
	utils.HandleError(c, err)
	insertResult, err := rulesCollection.InsertOne(ctx, api)
	utils.HandleError(c, err)
	return insertResult
}
