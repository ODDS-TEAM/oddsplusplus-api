package api

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	model "gitlab.odds.team/plus1/backend-go/model"
	"gopkg.in/mgo.v2/bson"
)

func (db *MongoDB) GetUserReserveItem(c echo.Context) error {
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

func (db *MongoDB) GetItemOrder(c echo.Context) error {
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

func (db *MongoDB) GetOrderCount(c echo.Context) error {
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

func (db *MongoDB) GetSummary(c echo.Context) error {
	var data model.Summary
	var itemData model.Item
	var reserveData []model.Reserve
	itemID := c.Param("itemId")
	query := bson.M{
		"item": bson.ObjectIdHex(itemID),
	}
	if err := db.ICol.Find(query).One(&itemData); err != nil {
		return err
	}
	if err := db.RCol.Find(query).All(&reserveData); err != nil {
		return err
	}
	data.Item = itemData
	data.Reserve = reserveData
	if err := db.SumCol.Insert(data); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, data)
}

func (db *MongoDB) DeleteReserve(c echo.Context) error {
	reserveID := c.Param("reserveId")
	var data model.Reserve
	query := bson.M{
		"_id": bson.ObjectIdHex(reserveID),
	}
	if err := db.RCol.Find(query).One(&data); err != nil {
		fmt.Println("In find Reserve error ", err)
		return err
	}
	var itemData model.Item
	queryItem := bson.M{
		"_id": data.Item,
	}
	if err := db.ICol.Find(queryItem).One(&itemData); err != nil {
		fmt.Println("In find Item error ", err)
		return err
	}
	q := bson.M{
		"_id": itemData.ItemId,
	}
	ob := bson.M{
		"$set": bson.M{
			"count": itemData.Count - data.Count,
		},
	}
	// itemData.Count = itemData.Count - data.Count
	if err := db.ICol.Update(q, &ob); err != nil {
		fmt.Println("In update Item error ", err)
		return err
	}
	if err := db.RCol.RemoveId(reserveID); err != nil {
		fmt.Println("In remove Reserve Error ", err)
		return err
	}
	return c.JSON(http.StatusOK, "Remove Reserve Success!")
}