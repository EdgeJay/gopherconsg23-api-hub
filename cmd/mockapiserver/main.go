package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/EdgeJay/gopherconsg23-api-hub/internal/mockapiserver"
	"github.com/labstack/echo/v4"
	"github.com/pb33f/libopenapi"
	v3high "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/pb33f/libopenapi/renderer"
)

type AppFlags struct {
	InputFile string
}

func retrieveAppFlags() AppFlags {
	appFlags := AppFlags{}

	flag.StringVar(&appFlags.InputFile, "input", "", "Path to OpenAPI 3+ specs file")
	flag.Parse()

	return appFlags
}

func addMockDataForMethod(
	mapping *mockapiserver.MockDataMapping,
	generator *renderer.MockGenerator,
	model *libopenapi.DocumentModel[v3high.Document],
	path, method string,
	op *v3high.Operation,
) {
	if op != nil {

		fmt.Println(op.Extensions["x-mock-mapping-file"])

		// iterate thru HTTP status codes
		for code, res := range op.Responses.Codes {

			// iterate thru Content-Type
			for _, mediaType := range res.Content {

				key := method + ":" + code + ":" + path

				// build schema
				schema := mediaType.Schema.Schema()
				// generate a mock of the schema
				if b, err := generator.GenerateMock(schema, ""); err == nil {
					(*mapping)[key] = string(b)
				}
			}
		}
	}
}

func addMockDataForPath(
	mapping *mockapiserver.MockDataMapping,
	generator *renderer.MockGenerator,
	model *libopenapi.DocumentModel[v3high.Document],
	path string,
	pathItem *v3high.PathItem,
) {
	addMockDataForMethod(mapping, generator, model, path, http.MethodGet, pathItem.Get)
	addMockDataForMethod(mapping, generator, model, path, http.MethodPut, pathItem.Put)
	addMockDataForMethod(mapping, generator, model, path, http.MethodPost, pathItem.Post)
	addMockDataForMethod(mapping, generator, model, path, http.MethodDelete, pathItem.Delete)
	addMockDataForMethod(mapping, generator, model, path, http.MethodOptions, pathItem.Options)
	addMockDataForMethod(mapping, generator, model, path, http.MethodHead, pathItem.Head)
	addMockDataForMethod(mapping, generator, model, path, http.MethodPatch, pathItem.Patch)
	addMockDataForMethod(mapping, generator, model, path, http.MethodTrace, pathItem.Trace)
}

func main() {

	appFlags := retrieveAppFlags()

	// create a new JSON mock generator
	mg := renderer.NewMockGenerator(renderer.JSON)

	// tell the mock generator to pretty print the output
	mg.SetPretty()

	savingsApiDoc, _ := os.ReadFile(appFlags.InputFile)

	// create a new document from specification and build a v3 model.
	document, err := libopenapi.NewDocument(savingsApiDoc)
	fmt.Println("Create doc error:", err)
	v3Model, _ := document.BuildV3Model()

	// iterate through all paths
	mockDataMapping := &mockapiserver.MockDataMapping{}
	for path, pathItem := range v3Model.Model.Paths.PathItems {
		addMockDataForPath(mockDataMapping, mg, v3Model, path, pathItem)
	}

	e := echo.New()
	RegisterHandlersWithMockData(e, mockDataMapping)
	e.Logger.Fatal(e.Start(":1337"))
}
