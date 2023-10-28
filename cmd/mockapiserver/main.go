package main

import (
	"fmt"
	"log"

	"github.com/EdgeJay/gopherconsg23-api-hub/internal/mockapiserver"
	"github.com/labstack/echo/v4"
)

func main() {

	appFlags := mockapiserver.NewAppFlags()

	v3Model := mockapiserver.NewOpenAPIV3ModelFromFile(appFlags)

	mockDataMapping := mockapiserver.NewMockDataMapping(v3Model, false)

	e := echo.New()
	RegisterHandlersWithMockData(e, mockDataMapping)
	log.Printf("Mock API server starting up\nListening to port %d", appFlags.Port)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", appFlags.Port)))
}
