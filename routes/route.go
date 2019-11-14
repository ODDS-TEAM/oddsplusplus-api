package routes

import (
	"gitlab.odds.team/plus1/backend-go/api"
	"github.com/labstack/echo"
)

// Init initialize api routes and set up a connection.
func Init(e *echo.Echo) {
	// Database connection.
	db, err := api.NewMongoDB()
	if err != nil {
		e.Logger.Fatal(err)
	}

	a := &api.MongoDB{
		Conn: db.Conn,
		UCol: db.UCol,
		ICol: db.ICol,
		SCol: db.SCol,
		SumCol: db.SumCol,
		RCol: db.RCol,
		TCol: db.TCol,
	}

	e.POST("/book", a.GetBookDetail)
	e.POST("/responseScrap", a.GetBookDetail)

	//Api item
	e.GET("/items/users/:userId", a.GetUserItem)
	e.POST("/additem", a.AddItem)
}
