package dbutils

import (
	"generic/models"
)

func FindApi(clientId string, apiName string) (api models.ApisMongo, found bool) {
	// var apiFound models.ApisMongo
	// apisCollection, ctx := config.GetMongoCollection("apis")
	// queryClient := bson.M{"clientId": clientId, "apiName": apiName}
	// apiFind := apisCollection.FindOne(ctx, queryClient)
	// if apiFind.Err() != nil {
	// 	return apiFound, false
	// }
	// apiFind.Decode(&apiFound)
	// return apiFound, true
	return models.ApisMongo{}, true
}

func UpdateApi(clientId string, apiName string, data models.ApisMongo) (errReturn error) {
	// queryClient := bson.M{"clientId": clientId, "apiName": apiName}
	// apisCollection, ctx := config.GetMongoCollection("apis")
	// _, err := apisCollection.ReplaceOne(ctx, queryClient, data)
	// return err
	return nil
}
