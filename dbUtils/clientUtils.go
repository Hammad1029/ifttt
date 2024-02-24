package dbutils

import (
	"generic/models"
)

func FindClient(clientId string) (client models.ClientModelMongo, found bool) {
	// var clientFound models.ClientModelMongo
	// clientsCollection, ctx := config.GetMongoCollection("clients")
	// queryClient := bson.M{"clientId": clientId}
	// clientFind := clientsCollection.FindOne(ctx, queryClient)
	// if clientFind.Err() != nil {
	// 	return clientFound, false
	// }
	// clientFind.Decode(&clientFound)
	// return clientFound, true
	return models.ClientModelMongo{}, false
}

func UpdateClient(clientId string, data models.ClientModelMongo) (errReturn error) {
	// queryClient := bson.M{"clientId": clientId}
	// clientsCollection, ctx := config.GetMongoCollection("clients")
	// _, err := clientsCollection.ReplaceOne(ctx, queryClient, data)
	// return err
	return nil
}

func GetClientVMap(clientId string) (vmapReturn map[string]map[string]string, found bool) {
	// var vmap map[string]map[string]string
	// clientsCollection, ctx := config.GetMongoCollection("clients")
	// filter := bson.M{"clientId": clientId}
	// opts := options.FindOne().SetProjection(bson.D{{Key: "virtualMap", Value: 1}})
	// vmapFound := clientsCollection.FindOne(ctx, filter, opts)
	// if vmapFound.Err() != nil {
	// 	return vmap, false
	// }
	// vmapFound.Decode(&vmap)
	// return vmap, true
	return make(map[string]map[string]string), false
}
