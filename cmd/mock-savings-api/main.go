package main

import (
	"fmt"
	"os"

	"github.com/EdgeJay/gopherconsg23-api-hub/cmd/mock-savings-api/server"
	"github.com/pb33f/libopenapi"
	"github.com/pb33f/libopenapi/renderer"

	"github.com/labstack/echo/v4"
)

func main() {

	// create a new JSON mock generator
	mg := renderer.NewMockGenerator(renderer.JSON)

	// tell the mock generator to pretty print the output
	mg.SetPretty()

	savingsApiDoc, _ := os.ReadFile("./docs/savings-api/savings-api.yaml")

	// create a new document from specification and build a v3 model.
	document, err := libopenapi.NewDocument(savingsApiDoc)
	fmt.Println("Create doc error:", err)
	v3Model, _ := document.BuildV3Model()

	// get the index
	index := v3Model.Index
	fmt.Println(len(index.GetAllPaths()))

	// create model
	savingsRecordsModel := v3Model.Model.Components.Schemas["GetSavingsRecordsSuccessResponse"]

	// build the schema
	savingsRecords := savingsRecordsModel.Schema()

	// generate a mock of the schema
	mock, err := mg.GenerateMock(savingsRecords, "")

	if err != nil {
		panic(err)
	}

	// print the mock to stdout
	fmt.Println(string(mock))

	e := echo.New()
	server.RegisterHandlersWithMockData(e)
	e.Logger.Fatal(e.Start(":1337"))
}
