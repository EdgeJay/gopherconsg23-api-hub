package main

import (
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.Static("/docs", "./docs")
	e.Logger.Fatal(e.Start(":1336"))
}
