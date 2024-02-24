package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ClientModelMongo struct {
	ClientId   string                       `json:"clientId"`
	ClientName string                       `json:"clientName"`
	Apis       []primitive.ObjectID         `json:"apis"`
	Core       string                       `json:"clientCore"`
	VMap       map[string]map[string]string `json:"virtualMap"`
}
