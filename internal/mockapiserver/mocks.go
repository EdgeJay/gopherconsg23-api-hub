package mockapiserver

import (
	"fmt"
	"net/http"
	"os"

	"github.com/EdgeJay/gopherconsg23-api-hub/internal/common"
	"github.com/labstack/gommon/log"
	"github.com/pb33f/libopenapi"
	"github.com/pb33f/libopenapi/datamodel"
	v3high "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/pb33f/libopenapi/renderer"
	"gopkg.in/yaml.v3"
)

func NewJSONMockGenerator(prettyPrint bool) *renderer.MockGenerator {
	// create a new JSON mock generator
	generator := renderer.NewMockGenerator(renderer.JSON)

	// tell the mock generator to pretty print the output
	if prettyPrint {
		generator.SetPretty()
	}

	return generator
}

func NewOpenAPIV3ModelFromFile(appFlags AppFlags) *libopenapi.DocumentModel[v3high.Document] {
	inputFileBytes, err := GetInputFileFromAppFlags(appFlags)
	if err != nil {
		common.LogFatalError("Input failed retrieval failed", err)
	}

	// create a new document from specification and build a v3 model.
	basePath, err := os.Getwd()
	if err != nil {
		common.LogFatalError("Document creation failed, cannot get working directory", err)
	}

	document, err := libopenapi.NewDocumentWithConfiguration(inputFileBytes, &datamodel.DocumentConfiguration{
		AllowFileReferences: true,
		BasePath:            basePath,
	})
	if err != nil {
		common.LogFatalError("Document creation failed", err)
	}

	model, errs := document.BuildV3Model()
	if model == nil {
		common.LogFatalError("Model creation failed", errs)
	}

	return model
}

func NewMockDataMapping(
	v3Model *libopenapi.DocumentModel[v3high.Document],
	prettyPrintMock bool,
) *MockDataMapping {
	// create a new JSON mock generator
	generator := NewJSONMockGenerator(prettyPrintMock)

	// iterate through all paths
	mockDataMapping := &MockDataMapping{}
	for path, pathItem := range v3Model.Model.Paths.PathItems {
		addMockDataForPath(mockDataMapping, generator, v3Model, path, pathItem)
	}

	return mockDataMapping
}

func addMockDataForMethod(
	mapping *MockDataMapping,
	generator *renderer.MockGenerator,
	model *libopenapi.DocumentModel[v3high.Document],
	path, method string,
	op *v3high.Operation,
) {
	if op != nil {

		// iterate thru HTTP status codes
		for code, res := range op.Responses.Codes {

			key := method + ":" + code + ":" + path

			// iterate thru Content-Type
			for _, mediaType := range res.Content {

				schema := getSchemaForMockGeneration(mediaType)

				// generate a mock of the schema
				if b, err := generator.GenerateMock(schema, ""); err == nil {
					(*mapping)[key] = string(b)
				} else {
					log.Printf("Cannot generator mock, %v\n", err)
				}
			}
		}
	}
}

func addMockDataForPath(
	mapping *MockDataMapping,
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

// This function helps to choose, resolve and return best node
// for mock generation.
// If mediaType has examples referencing external sources,
// this function will retrieve the file as well.
func getSchemaForMockGeneration(mediaType *v3high.MediaType) any {
	if mediaType != nil {

		if mappingFilePath, ok := mediaType.Extensions["x-examples-mapping"].(string); ok {
			if b, err := os.ReadFile(mappingFilePath); err == nil {
				var examplesMapping MockDataExamplesMapping
				yaml.Unmarshal(b, &examplesMapping)
				fmt.Println(examplesMapping)
			}
		}

		if mediaType.Examples != nil && len(mediaType.Examples) > 0 {
			for _, ex := range mediaType.Examples {

				if ex.ExternalValue != "" {
					// extract examples from file, then overwrite value in mediaType
					if b, err := os.ReadFile(ex.ExternalValue); err == nil {
						ex.Value = string(b)
					}
				}
			}
			return mediaType
		} else {
			return mediaType.Schema.Schema()
		}
	}
	return nil
}
