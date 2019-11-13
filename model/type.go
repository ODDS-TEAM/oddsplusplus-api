package model

import (
	"gopkg.in/mgo.v2/bson"
)

type Type struct {
	TypeId bson.ObjectId `json:"_id" bson:"_id"`
	Type   string        `json:"type" bson:"type"`
}
