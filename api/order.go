package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	model "gitlab.odds.team/plus1/backend-go/model"
	"gopkg.in/mgo.v2/bson"
)

func (db *MongoDB) GetAllReserves(c echo.Context) error {
	var data []model.Item
	if err := db.RCol.Find(bson.M{}).All(&data); err != nil {
		fmt.Println("In fine All Reserves error", err)
		return err
	}
	return c.JSON(http.StatusOK, data)
}

func (db *MongoDB) GetUserReserveItem(c echo.Context) error {
	// fmt.Println("In Get reserve by userid")
	data := []bson.M{}
	userID := c.Param("userId")

	ItemLookup := bson.M{
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

	SLoopup := bson.M{
		"$lookup": bson.M{
			"from":         "Status",
			"localField":   "Item.status",
			"foreignField": "_id",
			"as":           "Item.Status",
		},
	}

	SUnwind := bson.M{
		"$unwind": bson.M{
			"path":                       "$Item.Status",
			"preserveNullAndEmptyArrays": true,
		},
	}

	Ulookup := bson.M{
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

	MReserve := bson.M{
		"$match": bson.M{
			"user": bson.ObjectIdHex(userID),
		},
	}
	query := []bson.M{MReserve, ItemLookup, IUnwind, SLoopup, SUnwind, Ulookup, UUnwind}
	if err := db.RCol.Pipe(query).All(&data); err != nil {
		fmt.Println("Error in find reserve ", err)
		return err
	}
	// fmt.Println("reserve = ", data)
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
	fmt.Println("Before find user id = ", userID, " item id = ", itemID)

	query := bson.M{
		"user": bson.ObjectIdHex(userID),
		"item": bson.ObjectIdHex(itemID),
	}
	if err := db.RCol.Find(query).One(&data); err != nil && err.Error() != "not found" {
		fmt.Println("Error in find reserve ", err)
		return err
	} 
	if data.Item != bson.ObjectIdHex(itemID) {
		c.JSON(http.StatusOK, 0)
	}
	return c.JSON(http.StatusOK, data.Count)
}

func (db *MongoDB) AddReserve(c echo.Context) error {
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
		fmt.Println("Error in parse string to int ", err)
		return err
	}
	isDefualt := &model.Reserve{}
	if reserve.ReserveId != isDefualt.ReserveId {
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

		item.Count = item.Count + temp
		if err := db.ICol.Update(bson.M{"_id": item.ItemId}, &item); err != nil {
			fmt.Println("Error in update item ", err)
			return err
		}
	}
	return c.JSON(http.StatusOK, "Reserve Successed!")
}

func (db *MongoDB) Order(c echo.Context) error {
	userId := c.Param("userId")
	itemId := c.Param("itemId")
	count, err := strconv.Atoi(c.Param("count"))
	if err != nil {
		fmt.Println("Error in parse string to int ", err)
		return err
	}
	item := &model.Item{}
	query_item := bson.M{
		"_id": bson.ObjectIdHex(itemId),
	}
	if err := db.ICol.Find(query_item).One(&item); err != nil {
		fmt.Println("Error in find item ", err)
		return err
	}

	query_user := bson.M{
		"_id": bson.ObjectIdHex(userId),
	}
	user := &model.User{}
	if err := db.UCol.Find(query_user).One(&user); err != nil {
		fmt.Println("Error in find user ", err)
		return err
	}

	reserve := &model.Reserve{}
	query_reserve := bson.M{
		"user": user.UserId,
		"item": item.ItemId,
	}
	if err := db.RCol.Find(query_reserve).One(&reserve); err != nil {
		fmt.Println("Error in find reserve ", err)
		return err
	}
	item.Count = item.Count + count - reserve.Count
	if err := db.ICol.Update(bson.M{"_id": item.ItemId}, &item); err != nil {
		fmt.Println("Error in update item ", err)
		return err
	}
	reserve.Count = count
	if err := db.RCol.Update(bson.M{"_id": reserve.ReserveId}, &reserve); err != nil {
		fmt.Println("Error in update reserve ", err)
		return err
	}
	return c.JSON(http.StatusOK, "Oder Successed!")
}

func (db *MongoDB) UpdateOrder(c echo.Context) error {
	itemId := c.Param("itemId")
	totaoPrice, err := strconv.ParseFloat(c.Param("totalPrice"), 64)
	if err != nil {
		fmt.Println("Error in parsefloat ", err)
		return err
	}
	charge, err := strconv.ParseFloat(c.Param("charge"), 64)
	if err != nil {
		fmt.Println("Error in parse float ", err)
	}

	item := &model.Item{}
	if err := db.ICol.Find(bson.M{"_id": bson.ObjectIdHex(itemId)}).One(&item); err != nil {
		fmt.Println("Error in find item ", err)
		return err
	}
	item.Cost = totaoPrice
	item.ShippingCharge = charge
	status := &model.Status{}
	if err := db.SCol.Find(bson.M{"status": "Shipping"}).One(&status); err != nil {
		fmt.Println("Error in find status ", err)
		return err
	}
	item.Status = status.StatusId
	if err := db.ICol.Update(bson.M{"_id": item.ItemId}, &item); err != nil {
		fmt.Println("Error in update item ", err)
		return err
	}

	costPerItem := totaoPrice / float64(item.Count)
	chargePerItem := charge / float64(item.Count)

	reserves := []model.Reserve{}
	if err := db.RCol.Find(bson.M{"item": item.ItemId}).All(&reserves); err != nil {
		fmt.Print("Error in find reserves ", err)
		return err
	}

	for i, reserve := range reserves {
		reserve.Cost = float64(reserve.Count) * costPerItem
		reserve.ShippingCharge = float64(reserve.Count) * chargePerItem
		if err := db.RCol.Update(bson.M{"_id": reserve.ReserveId}, &reserve); err != nil {
			fmt.Println("Error in update Reserve ", err, " index = ", i)
			return err
		}
	}
	return c.JSON(http.StatusOK, "Update Order Successed!")
}

func (db *MongoDB) DeleteOrderByUserAndItem(c echo.Context) error {
	userId := c.Param("userId")
	itemId := c.Param("itemId")

	item := &model.Item{}
	defaultItem := &model.Item{}
	if err := db.ICol.Find(bson.M{"_id": bson.ObjectIdHex(itemId), "user": bson.ObjectIdHex(userId)}).One(&item); err != nil && err.Error() != "not found" {
		fmt.Println("Error in find item ", err)
		return err
	} else if item.ItemId != defaultItem.ItemId {
		reserves := []model.Reserve{}
		if err := db.RCol.Find(bson.M{"item": item.ItemId}).All(&reserves); err != nil {
			fmt.Println("Error in find reserves all by itemid ", err)
			return err
		}

		for index, reserveEle := range reserves {
			if err := db.RCol.RemoveId(reserveEle.ReserveId); err != nil {
				fmt.Println("Error in remove reserve index  = ", index, " error = ", err)
				return err
			}
		}

		if err := db.ICol.RemoveId(item.ItemId); err != nil {
			fmt.Println("Error in remove item when count = 0")
			return err
		}
	} else {
		reserve := &model.Reserve{}
		if err := db.RCol.Find(bson.M{"item": bson.ObjectIdHex(itemId), "user": bson.ObjectIdHex(userId)}).One(&reserve); err != nil {
			fmt.Println("Error in find reserve ", err)
			return err
		}

		if err := db.RCol.RemoveId(reserve.ReserveId); err != nil {
			fmt.Println("Error in remove reserve ", err)
			return err
		}
	}

	// user := &model.User{}
	// if err := db.UCol.Find(bson.M{"_id": bson.ObjectIdHex(userId)}).One(&user); err != nil {
	// 	fmt.Println("Error in find user ", err)
	// 	return err
	// }

	return c.JSON(http.StatusOK, "Delete Reserve Successed!")
}

func (db *MongoDB) GetSummary(c echo.Context) error {
	var data model.Summary
	var itemData model.Item
	var reserveData []model.Reserve
	itemID := c.Param("itemId")
	query := bson.M{
		"_id": bson.ObjectIdHex(itemID),
	}
	if err := db.ICol.Find(query).One(&itemData); err != nil {
		return err
	}
	if err := db.RCol.Find(bson.M{"item": itemData.ItemId}).All(&reserveData); err != nil {
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
	userId := c.Param("userId")
	// fmt.Println("User = ", userId, " Item = ", )
	user := &model.User{}
	if err := db.UCol.Find(bson.M{"_id": bson.ObjectIdHex(userId)}).One(&user); err != nil {
		fmt.Println("Error in find user ", err)
		return err
	}

	var data model.Reserve
	query := bson.M{
		"_id": bson.ObjectIdHex(reserveID),
	}
	if err := db.RCol.Find(query).One(&data); err != nil {
		fmt.Println("In find Reserve error ", err)
		return err
	}
	itemData := &model.Item{}
	queryItem := bson.M{
		"_id":  data.Item,
		"user": user.UserId,
	}
	fmt.Println("Before find item")
	if err := db.ICol.Find(queryItem).One(&itemData); err != nil && err.Error() != "not found" {
		fmt.Println("In find Item error ", err)
		return err
	}
	if itemData.User == user.UserId { // find this book is your
		fmt.Println("Delete all reserve")

		reserves := []model.Reserve{}
		if err := db.RCol.Find(bson.M{"item": itemData.ItemId}).All(&reserves); err != nil {
			fmt.Println("Error in find reserves ", err)
			return err
		}

		for index, reserveEle := range reserves {
			if err := db.RCol.RemoveId(reserveEle.ReserveId); err != nil {
				fmt.Println("Error in remove reserve index  = ", index, " error = ", err)
				return err
			}
		}

		if err := db.ICol.Remove(itemData); err != nil {
			fmt.Println("Error in remove item " , err)
			return err
		}

	}else {
		fmt.Println("Not book owner")
		item := &model.Item{}
		if err := db.ICol.Find(bson.M{"_id": data.Item}).One(&item); err != nil {
			fmt.Println("Error in find item ", err)
			return err
		}
		q := bson.M{
			"_id": item.ItemId,
		}
		ob := bson.M{
			"$set": bson.M{
				"count": item.Count - data.Count,
			},
		}

 		if err := db.ICol.Update(q, &ob); err != nil {
			fmt.Println("In update Item error ", err)
			return err
		}
		if err := db.RCol.RemoveId(data.ReserveId); err != nil {
			fmt.Println("In remove Reserve Error ", err)
			return err
		}
	}

	return c.JSON(http.StatusOK, "Remove Reserve Success!")
}
