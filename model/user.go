package model

import (
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	UserId bson.ObjectId `json:"_id," bson:"_id"`
	Name   string        `json:"name" bson:"name"`
	Email  string        `json:"email" bson:"email"`
	ImgUrl string        `json:"imgUrl" bson:"imgUrl"`
	Token  string        `json:"token" bson:"token"`
}
