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
		SumCol: db.SCol,
		ICol:   db.ICol,
		TCol:   db.TCol,
	}

	e.POST("/book", api.GetBookDetail)
	e.POST("/responseScrap", api.GetBookDetail)
	e.POST("/login", a.Login)
	e.POST("/update", a.Register)
}
