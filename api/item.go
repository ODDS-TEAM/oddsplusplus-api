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
	fmt.Println("In Add Item")
	item := &model.Item{}
	if err := c.Bind(item); err != nil {
		fmt.Println("In c.Bind Error")
		return err
	}

	Type := &model.Type{}
	query_type := bson.M{
		"type": "Book",
	}
	if err := db.TCol.Find(query_type).One(&Type); err != nil {
		fmt.Println("In Find Type Error ", Type, err)
		return err
	}

	status := &model.Status{}
	query_status := bson.M{
		"status": "Pending",
	}
	if err := db.SCol.Find(query_status).One(&status); err != nil {
		fmt.Println("In Status Error")
		return err
	}
	
	item.User = bson.ObjectId(item.User)
	item.CreateOn = time.Now()
	item.Status = status.StatusId
	item.Type = Type.TypeId

	if err := db.ICol.Insert(item); err != nil {
		fmt.Println("In Insert Error")
		return err
	}
	return c.JSON(http.StatusOK, item)
}