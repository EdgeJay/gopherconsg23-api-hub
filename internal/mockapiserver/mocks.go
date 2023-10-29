package mockapiserver

import (
	"net/http"

	"github.com/pb33f/libopenapi"
	v3high "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/pb33f/libopenapi/renderer"
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
		LogFatalError("Input failed retrieval failed", err)
	}

	// create a new document from specification and build a v3 model.
	document, err := libopenapi.NewDocument(inputFileBytes)
	if err != nil {
		LogFatalError("Document creation failed", err)
	}

	model, errs := document.BuildV3Model()
	if model == nil {
		LogFatalError("Model creation failed", errs)
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
