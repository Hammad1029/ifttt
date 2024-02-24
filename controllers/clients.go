package controllers

import (
	"github.com/gin-gonic/gin"
)

func AddClient(c *gin.Context) {
	// reqBodyAny, _ := c.Get("Req")
	// reqBody := reqBodyAny.(*schemas.AddClient)

	// var client models.ClientModelMongo
	// client.ClientName = reqBody.Name
	// client.ClientId = uuid.New().String()

	// clientsCollection, ctx := config.GetMongoCollection("clients")
	// insertResult, err := clientsCollection.InsertOne(ctx, client)
	// utils.HandleError(c, err)

	// utils.ResponseHandler(c, utils.ResponseConfig{Data: utils.ConvertToMap("inserted", insertResult)})
}
