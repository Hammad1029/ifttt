package dbutils

import (
	"generic/config"
	"generic/models"

	"go.mongodb.org/mongo-driver/bson"
)

func FindClient(clientId string) (client models.ClientModelMongo, found bool) {
	var clientFound models.ClientModelMongo
	clientsCollection, ctx := config.GetMongoCollection("clients")
	queryClient := bson.M{"clientId": clientId}
	clientFind := clientsCollection.FindOne(ctx, queryClient)
	if clientFind.Err() != nil {
		return clientFound, false
	}
	clientFind.Decode(&clientFound)
	return clientFound, true
}

func UpdateClient(clientId string, data models.ClientModelMongo) (errReturn error) {
	queryClient := bson.M{"clientId": clientId}
	clientsCollection, ctx := config.GetMongoCollection("clients")
	_, err := clientsCollection.ReplaceOne(ctx, queryClient, data)
	return err
}
