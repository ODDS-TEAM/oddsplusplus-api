package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"gitlab.odds.team/plus1/backend-go/model"
	"gopkg.in/mgo.v2/bson"
)

func (db *MongoDB) GetUserItem(c echo.Context) error {
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

func (db *MongoDB) AddItem(c echo.Context) error {
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

	reserve := &model.Reserve{}
	reserve.ReserveId = bson.NewObjectId()
	reserve.Count = 1
	reserve.Item = item.ItemId
	reserve.User = item.User

	if err := db.RCol.Insert(reserve); err != nil {
		fmt.Println("Error in insert reserve ", err)
		return err
	}

	return c.JSON(http.StatusOK, item)
}

func (db *MongoDB) DeleteItem(c echo.Context) error {
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

func (db *MongoDB) GetItemData(c echo.Context) error {
	var data model.Item
	itemId := c.Param("itemId")
	query := bson.M{
		"_id": bson.ObjectIdHex(itemId),
	}
	if err := db.ICol.Find(query).One(&data); err != nil {
		fmt.Println("In fine Item error", err)
		return err
	}
	return c.JSON(http.StatusOK, data)
}

func (db *MongoDB) GetAllItem(c echo.Context) error {
	data := []bson.M{}
	var itemData []model.Item
	if err := db.ICol.Find(bson.M{}).All(&itemData); err != nil {
		fmt.Println("In fine All Items error", err)
		return err
	}
	SLookup := bson.M{
		"$lookup": bson.M{
			"from":         "Status",
			"localField":   "status",
			"foreignField": "_id",
			"as":           "Status",
		},
	}
	SUnwind := bson.M{
		"$unwind": bson.M{
			"path":                       "$Status",
			"preserveNullAndEmptyArrays": true,
		},
	}
	ULookup := bson.M{
		"$lookup": bson.M{
			"from":         "Users",
			"localField":   "user",
			"foreignField": "_id",
			"as":           "User",
		},
	}
	UUnwind := bson.M{
		"$unwind": bson.M{
			"path":                       "$User",
			"preserveNullAndEmptyArrays": true,
		},
	}
	TLookup := bson.M{
		"$lookup": bson.M{
			"from":         "Type",
			"localField":   "type",
			"foreignField": "_id",
			"as":           "Type",
		},
	}
	TUnwind := bson.M{
		"$unwind": bson.M{
			"path":                       "$Type",
			"preserveNullAndEmptyArrays": true,
		},
	}
	sort := bson.M{
		"$sort": bson.M{
			"createOn": -1,
		},
	}

	pipe_query := []bson.M{SLookup, SUnwind, ULookup, UUnwind, TLookup, TUnwind, sort}
	if err := db.ICol.Pipe(pipe_query).All(&data); err != nil {
		fmt.Println("Error in find item ", err)
		return err
	}
	fmt.Println(data)
	return c.JSON(http.StatusOK, data)
}

func (db *MongoDB) GetTopItem(c echo.Context) error {
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

func (db *MongoDB) GetItemReserve(c echo.Context) error {
	fmt.Println("In Get Reserve by item Id")
	itemId := c.Param("itemId")
	item := &model.Item{}
	if err := db.ICol.Find(bson.M{"_id": bson.ObjectIdHex(itemId)}).One(&item); err != nil {
		fmt.Println("Error in find item ", err)
	}
	fmt.Println("Before find reserves")
	reserves := []bson.M{}

	ULookup := bson.M{
		"$lookup": bson.M{
			"from":         "Users",
			"localField":   "user",
			"foreignField": "_id",
			"as":           "User",
		},
	}
	UUnwind := bson.M{
		"$unwind": bson.M{
			"path":                       "$User",
			"preserveNullAndEmptyArrays": true,
		},
	}

	ILookup := bson.M{
		"$lookup": bson.M{
			"from":         "Item",
			"localField":   "item",
			"foreignField": "_id",
			"as":           "Item",
		},
	}
	IUnwind := bson.M{
		"$unwind": bson.M{
			"path":                       "$Item",
			"preserveNullAndEmptyArrays": true,
		},
	}
	Slookup := bson.M{
		"$lookup": bson.M{
			"from":         "Status",
			"localField":   "Item.status",
			"foreignField": "_id",
			"as":           "Item.Status",
		},
	}
	MReserve := bson.M{
		"$match": bson.M{
			"item": bson.ObjectIdHex(itemId),
		},
	}
	query := []bson.M{MReserve, ILookup, IUnwind, Slookup, ULookup, UUnwind}
	if err := db.RCol.Pipe(query).All(&reserves); err != nil {
		fmt.Println("Error in find reserves", err)
		return err
	}
	return c.JSON(http.StatusOK, reserves)
}
