package main

import (
	"net/http"

	"github.com/EdgeJay/gopherconsg23-api-hub/internal/codegenclient"
	"github.com/labstack/echo/v4"
)

func main() {
	appFlags := codegenclient.NewAppFlags()

	service := NewApiService(appFlags.InputFile[0])
	service.Load()

	e := echo.New()
	e.GET("/v1/resident-data", func(c echo.Context) error {
		data := service.GetResidentData()
		return c.JSON(http.StatusOK, data)
	})
	e.Logger.Fatal(e.Start(":1323"))
}
