package models

import (
	"generic/config"
	"generic/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RuleModelMongo struct {
	Name string                `json:"name" binding:"required"`
	Type string                `json:"type" binding:"required"`
	Rule ruleDetailsModelMongo `json:"rule" binding:"required"`
}

type ruleDetailsModelMongo struct {
	Operator1   interface{}   `json:"op1" binding:"required"`
	Operator2   interface{}   `json:"op2" binding:"required"`
	Operator    string        `json:"operator" binding:"required"`
	ThenActions []interface{} `json:"thenActions" binding:"required"`
	ElseActions []interface{} `json:"elseActions" binding:"required"`
	Then        interface{}   `json:"then" binding:"required"`
	Else        interface{}   `json:"else" binding:"required"`
}

func GetRule(id string) (RuleModelMongo, string) {
	var rule RuleModelMongo
	rulesCollection, ctx := config.GetMongoCollection("rules")
	objectId, err := primitive.ObjectIDFromHex(id)
	utils.HandleError(nil, err)
	result := rulesCollection.FindOne(ctx, bson.M{"_id": objectId})
	err = result.Decode(&rule)
	if err == mongo.ErrNoDocuments {
		return RuleModelMongo{}, ""
	}
	raw, err := result.Raw()
	utils.HandleError(nil, err)
	return rule, raw.String()
}
