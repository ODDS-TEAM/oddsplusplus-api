package main

import (
	"gitlab.odds.team/plus1/backend-go/routes"
	"gitlab.odds.team/plus1/backend-go/config"
	"github.com/labstack/echo"
)

func main() {
	// fmt.Println("hello")
	e := echo.New()
	s := config.Spec()
	routes.Init(e)
	e.Logger.Fatal(e.Start(s.APIPort))
}
