package model

import (
	"gopkg.in/mgo.v2/bson"
)

type Summary struct {
	Item    bson.ObjectId      `json:"item" bson:"item"`
	Reserve []bson.ObjectId `json:"reserve" bson:"reserve"`
}
