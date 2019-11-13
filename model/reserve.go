package model

import (
	"gopkg.in/mgo.v2/bson"
)

type Reserve struct {
	ReserveId      bson.ObjectId `json:"_id" bson:"_id"`
	Count          int           `json:"count" bson:"count"`
	Cost           float64       `json:"cost" json:"cost"`
	ShippingCharge float64       `json:"shippingCharge" bson:"shippingCharge"`
	Item           bson.ObjectId `json:"item" bson:"item"`
	User           bson.ObjectId `json:"user" bson:"user"`
}
