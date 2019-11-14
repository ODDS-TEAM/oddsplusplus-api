package api

import (
	"fmt"
	"net/http"
	"strconv"
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

func (db *MongoDB) AddReserve(c echo.Context) (error) {
	userId := c.Param("userId")
	itemId := c.Param("itemId")
	count := c.Param("count")

	item := &model.Item{}
	query_item := bson.M{
		"_id": bson.ObjectIdHex(itemId),
	}
	if err := db.ICol.Find(query_item).One(&item); err != nil {
		fmt.Println("Error in find Item ", err)
		return err
	}

	user := &model.User{}
	query_user := bson.M{
		"_id": bson.ObjectIdHex(userId),
	}
	if err := db.UCol.Find(query_user).One(&user); err != nil {
		fmt.Println("Error in find user ", err)
		return err
	}

	reserve := &model.Reserve{}
	query_reserve := bson.M{
		"user": user.UserId,
		"item": item.ItemId,
	}

	if err := db.RCol.Find(query_reserve).One(&reserve); err != nil && err.Error() != "not found" {
		fmt.Println("Error in find reserve ", err)
		return err
	}
	temp, err := strconv.Atoi(count)
	if err != nil {
		fmt.Println("Error in pharse string to int ", err)
		return err
	}
	isDefualt := &model.Reserve{}
	if reserve.ReserveId !=  isDefualt.ReserveId {
		query_update_reserve_q := bson.M{
			"_id": reserve.ReserveId,
		}
		query_update_reserve_ob := bson.M{
			"$set": bson.M{
				"count": temp + reserve.Count,
			},
		}
		if err := db.RCol.Update(query_update_reserve_q, &query_update_reserve_ob); err != nil {
			fmt.Println("Error in update reserve ", err)
			return err
		}

		query_update_item_q := bson.M{
			"_id": item.ItemId,
		}
		query_update_item_ob := bson.M{
			"$set": bson.M{
				"count": item.Count + temp,
			},
		}
		if err := db.ICol.Update(query_update_item_q, &query_update_item_ob); err != nil {
			fmt.Println("Error in update item ", err)
			return err
		}
	} else {
		reserve.ReserveId = bson.NewObjectId()
		reserve.Count = temp
		reserve.Item = item.ItemId
		reserve.User = user.UserId

		if err := db.RCol.Insert(reserve); err != nil {
			fmt.Println("Error in insert reserve ", err)
			return err
		}
	}
	return nil
}