package model

import (
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	UserId  bson.ObjectId `json:"_id," bson:"_id"`
	Name 	string `json:"name"`
	Email 	string `json:"email"`
	ImgUrl 	string `json:"imgUrl"`
	Token	string `json:"token"`
}
