package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"gitlab.odds.team/plus1/backend-go/config"
	"gitlab.odds.team/plus1/backend-go/routes"
)

func main() {
	// fmt.Println("hello")

	// Use labstack/echo for rich routing.
	// See https://echo.labstack.com/
	e := echo.New()
	s := config.Spec()

	// Middleware
	e.Logger.SetLevel(log.ERROR)
	e.Use(
		middleware.CORS(),
		middleware.Recover(),
		middleware.Logger(),
		// middleware.JWTWithConfig(middleware.JWTConfig{
		// 	SigningKey: []byte("sMJuczqQPYzocl1s6SLj"),
		// 	Skipper: func(c echo.Context) bool {
		// 		// Skip authentication for and login requests
		// 		if c.Path() == "/login" || c.Path() == "/_ah/health" {
		// 			return true
		// 		}
		// 		return false
		// 	},
		// }),
	)
	routes.Init(e)
	e.Logger.Fatal(e.Start(s.APIPort))
}
