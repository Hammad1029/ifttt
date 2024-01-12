package config

import (
	"context"
	"fmt"

	"github.com/grafadruid/go-druid"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var druidClient *druid.Client
var mongoClient *mongo.Client
var mongoCtx context.Context

func startDruid() {
	var druidOpts []druid.ClientOption
	druidOpts = append(druidOpts, druid.WithSkipTLSVerify())
	d, err := druid.NewClient("http://localhost:8082", druidOpts...)
	handleError(err)
	_, _, err = d.Common().Status()
	handleError(err)
	druidClient = d
}

func GetDruid() *druid.Client {
	return druidClient
}

func startMongo() {
	uri := fmt.Sprintf("mongodb://%s:%s/%s",
		GetConfigProp("database.mongo.host"), GetConfigProp("database.mongo.port"), GetConfigProp("database.mongo.name"))
	mongoCtx = context.TODO()
	bsonoptions := options.BSONOptions{UseJSONStructTags: true}
	newClient, err := mongo.Connect(mongoCtx, options.Client().ApplyURI(uri).SetBSONOptions(&bsonoptions))
	mongoClient = newClient
	handleError(err)
	err = mongoClient.Ping(mongoCtx, nil)
	handleError(err)
	fmt.Println("MongoDB connected")
}

func GetMongoCollection(name string) (*mongo.Collection, context.Context) {
	collection := mongoClient.Database(GetConfigProp("database.mongo.name")).Collection(name)
	return collection, mongoCtx
}
