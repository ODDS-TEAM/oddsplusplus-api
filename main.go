package main

import (
	// "plus1/controller"

	"github.com/labstack/echo"
)

func main() {
	// fmt.Println("hello")
	e := echo.New()
	// e.POST("/book", controller.GetBookDetail)
	e.Start(":8080")
}
