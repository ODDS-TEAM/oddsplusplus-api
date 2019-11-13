package main

import (
	"plus1/config"
	route "plus1/routes"

	"github.com/labstack/echo"
)

func main() {
	// fmt.Println("hello")
	e := echo.New()
	s := config.Spec()
	route.Init(e)
	e.Logger.Fatal(e.Start(s.APIPort))
}
