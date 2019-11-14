package api

import (
	"gitlab.odds.team/plus1/backend-go/model"
	"google.golang.org/api/oauth2/v1"
	"gopkg.in/mgo.v2/bson"
)

// FindUser search and validation of user data.
func (db *MongoDB) FindUser(tokeninfo *oauth2.Tokeninfo) (*model.User, string) {
	u := &model.User{}
	status := "new"
	q := bson.M{
		"email": tokeninfo.Email,
	}
	if err := db.UCol.Find(q).One(&u); err != nil {
		return nil, status
	}

	if u.Name == "" || u.Email == "" || u.ImgUrl == "" {
		status = "notComplete"
	} else {
		status = "exist"
	}

	return u, status
}
