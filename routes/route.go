package routes

import (
	"github.com/labstack/echo"
	"gitlab.odds.team/plus1/backend-go/api"
)

// Init initialize api routes and set up a connection.
func Init(e *echo.Echo) {
	// Database connection.
	db, err := api.NewMongoDB()
	if err != nil {
		e.Logger.Fatal(err)
	}

	a := &api.MongoDB{
		Conn:   db.Conn,
		UCol:   db.UCol,
		ICol:   db.ICol,
		SCol:   db.SCol,
		SumCol: db.SumCol,
		RCol:   db.RCol,
		TCol:   db.TCol,
	}

	e.POST("/book", a.GetBookDetail)
	e.POST("/responseScrap", a.GetBookDetail)

	//Api item
	e.GET("/items/users/:userId", a.GetUserItem)
	e.POST("/additem", a.AddItem)
	e.DELETE("/items/users/:itemId/:userId", a.DeleteItem)
	
	e.GET("/items/:itemId", a.GetItemData)
	e.GET("/items", a.GetAllItem)

	//Api order
	e.GET("/reserves/users/:userId", a.GetUserReserveItem)
	e.GET("/reserves/items/:itemId", a.GetItemOrder)
	e.GET("/reserves/users/:userId/:itemId", a.GetOrderCount)

	e.DELETE("/reserves/:userId/:itemId", a.DeleteOrderByUserAndItem)
	e.PATCH("/updateOrder/:itemId/:totalPrice/:charge", a.UpdateOrder)
	e.POST("/order/:userId/:itemId/:count", a.Order)
	e.POST("/reserves/:userId/:itemId/:count", a.AddReserve)

	e.GET("/reserves/sum/:itemId", a.GetSummary)
	e.DELETE("reserves/:reserveId", a.DeleteReserve)
	e.GET("/getreserves/:itemId", a.GetItemReserve)

	e.GET("/reserves", a.GetAllReserves)
	e.GET("/reserves/:itemId", a.GetItemReserve)
	e.GET("/findTopItem/:userId", a.GetTopItem)
}
