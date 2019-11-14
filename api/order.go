package api

import (
	"net/http"

	"github.com/labstack/echo"
	model "gitlab.odds.team/plus1/backend-go/model"
	"gopkg.in/mgo.v2/bson"
)

func (db *MongoDB) GetUserReserveItem(c echo.Context) (err error) {
	var data []model.Reserve
	userID := c.Param("userId")
	query := bson.M{
		"user": bson.ObjectIdHex(userID),
	}
	if err := db.RCol.Find(query).All(&data); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, data)
}

func (db *MongoDB) GetItemOrder(c echo.Context) (err error) {
	var data []model.Reserve
	itemID := c.Param("itemId")
	query := bson.M{
		"item": bson.ObjectIdHex(itemID),
	}
	if err := db.RCol.Find(query).All(&data); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, data)
}

func (db *MongoDB) GetOrderCount(c echo.Context) (err error) {
	var data model.Reserve
	userID := c.Param("userId")
	itemID := c.Param("itemId")
	query := bson.M{
		"user": bson.ObjectIdHex(userID),
		"item": bson.ObjectIdHex(itemID),
	}
	if err := db.RCol.Find(query).One(&data); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, data.Count)
}

