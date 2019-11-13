package model

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Item struct {
	ItemId    bson.ObjectId `json:"_id" bson:"_id"`
	Url       string        `json:"url" bson:"url"`
	Title     string        `json:"title" bson:"title"`
	Author    string        `json:"author" bson:"author"`
	Format    string        `json:"format" bson:"format"`
	ImgUrl    string        `json:"imgUrl" bson:"imgUrl"`
	Price     float64       `json:"price" bson:"price"`
	OrderDate time.Time     `json:"orderDate" bson:"orderDate"`
	Count     int           `json:"count" bson:"count"`
	Cost      float64       `json:"cost" bson:"cost"`
	CreateOn  time.Time     `json:"createOn" bson:"createOn"`
	Status    Status        `json:"status" bson:"status"`
	User      User          `json:"user" bson:"user"`
	Type      Type          `json:"type" bson:"type"`
}
