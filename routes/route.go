package routes

import (
	"gitlab.odds.team/plus1/backend-go/api"
	"github.com/labstack/echo"
)

// Init initialize api routes and set up a connection.
func Init(e *echo.Echo) {
	// Database connection.
	_, err := api.NewMongoDB()
	if err != nil {
		e.Logger.Fatal(err)
	}

	// a := &api.MongoDB{
	// 	Conn: db.Conn,
	// 	UCol: db.UCol,
	// }

	e.POST("/book", api.GetBookDetail)
	e.POST("/responseScrap", api.GetBookDetail)
}
