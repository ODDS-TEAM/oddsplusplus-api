package model

import (
	"github.com/dgrijalva/jwt-go"
	"gopkg.in/mgo.v2/bson"
)

// TokenGoogle hold the jwt token of googleapi
type TokenGoogle struct {
	Token string `json:"token"`
}

// TokenRes that is generated for the requested resource from api.
type TokenRes struct {
	Token      string `json:"token"`
	FirstLogin bool   `json:"firstLogin"`
}

// JwtCustomClaims built in the payload data with ID and role user.
type JwtCustomClaims struct {
	ID bson.ObjectId `json:"id"`
	jwt.StandardClaims
}
