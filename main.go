package main

import (
	"plus1/api"

	"github.com/labstack/echo"
)

func main() {
	// fmt.Println("hello")
	e := echo.New()
	e.POST("/book", api.GetBookDetail)
	e.Start(":8080")
}
