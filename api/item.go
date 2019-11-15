package api

import (
	"fmt"
	"time"
	"net/http"
	"gopkg.in/mgo.v2/bson"
	"gitlab.odds.team/plus1/backend-go/model"
	"github.com/labstack/echo"
)

func (db *MongoDB) GetUserItem(c echo.Context) (error){
	var items model.Items
	userId := c.Param("userId")
	query := bson.M{
		"user": bson.ObjectIdHex(userId),
	}
	if err := db.ICol.Find(query).All(&items); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, items)
}

func (db *MongoDB) AddItem(c echo.Context) (error) {
	item := &model.Item{}
	if err := c.Bind(item); err != nil {
		fmt.Println("In c.Bind Error ", err)
		return err
	}

	Type := &model.Type{}
	query_type := bson.M{
		"type": "Book",
	}
	if err := db.TCol.Find(query_type).One(&Type); err != nil {
		fmt.Println("In Find Type Error ", err)
		return err
	}

	status := &model.Status{}
	query_status := bson.M{
		"status": "Pending",
	}
	if err := db.SCol.Find(query_status).One(&status); err != nil {
		fmt.Println("In Status Error", err)
		return err
	}
	item.ItemId = bson.NewObjectId()
	item.CreateOn = time.Now()
	item.Status = status.StatusId
	item.Type = Type.TypeId

	if err := db.ICol.Insert(item); err != nil {
		fmt.Println("In Insert Error", err)
		return err
	}
	return c.JSON(http.StatusOK, item)
}

func (db *MongoDB) DeleteItem(c echo.Context) (error) {
	userId := c.Param("userId")
	itemId := c.Param("itemId")

	query_item := bson.M{
		"_id": bson.ObjectIdHex(itemId),
	}
	item := &model.Item{}
	if err := db.ICol.Find(query_item).One(&item); err != nil {
		fmt.Println("In find Item error ", err)
		return err
	}

	query_reserve := bson.M{
		"item": item.ItemId,
	}
	reserves := []*model.Reserve{}
	if err := db.RCol.Find(query_reserve).All(&reserves); err != nil {
		fmt.Println("In find reserve Error ", err)
		return err
	}

	if len(reserves) == 1 && item.User == reserves[0].User {
		query_find_reserve := bson.M{
			"user": bson.ObjectIdHex(userId),
			"item": item.ItemId,
		}
		reserve := &model.Reserve{}
		if err := db.RCol.Find(query_find_reserve).One(reserve); err != nil {
			fmt.Println("In Find Reserve by user and item error ", err)
			return err
		}
		if err := db.RCol.RemoveId(reserve.ReserveId); err != nil {
			fmt.Println("In Remove Reserve Error ", err)
			return err
		}
		if err := db.ICol.RemoveId(item.ItemId); err != nil {
			fmt.Println("In Error Remove Item ", err)
			return err
		}
	}
	return c.JSON(http.StatusOK, "Remove Item Success!")
}

func (db *MongoDB) GetTopItem(c echo.Context) (error) {
	userId := c.Param("userId")
	user := &model.User{}
	if err := db.UCol.Find(bson.M{"_id": bson.ObjectIdHex(userId)}).One(&user); err != nil {
		fmt.Println("Error in find user ", err)
		return err
	}

	items := []model.Item{}
	sort := bson.M{
		"$sort": bson.M{
			"createOn": -1,
		},
	}
	query_item := bson.M{
		"$match": bson.M{
			"user": user.UserId,
		},
	}
	pipe_query := []bson.M{query_item, sort}
	fmt.Println("Before Pipe")
	if err := db.ICol.Pipe(pipe_query).All(&items); err != nil {
		fmt.Println("Error in find item by user id ", err)
	}
	return c.JSON(http.StatusOK, items)
}


func (db *MongoDB) GetItemReserve(c echo.Context) (error) {
	itemId := c.Param("itemId")
	item := &model.Item{}
	if err := db.ICol.Find(bson.M{"_id": bson.ObjectIdHex(itemId)}).One(&item); err != nil {
		fmt.Println("Error in find item ", err)
	}
	fmt.Println("Before find reserves")

	reserves := []model.Reserve{}
	if err := db.RCol.Find(bson.M{"item": item.ItemId}).All(&reserves); err != nil {
		fmt.Println("Error in find reserves", err)
		return err
	}
	return c.JSON(http.StatusOK, reserves)
}