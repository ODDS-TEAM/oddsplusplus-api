package model

import (
	"gopkg.in/mgo.v2/bson"
)

type Status struct {
	StatusId	bson.ObjectId	`json:"_id," bson:"_id"`
	Status		string 			`json:"status" bson:"status"`
}