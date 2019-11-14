package model

type Summary struct {
	Item    Item      `json:"item" bson:"item"`
	Reserve []Reserve `json:"reserve" bson:"reserve"`
}
