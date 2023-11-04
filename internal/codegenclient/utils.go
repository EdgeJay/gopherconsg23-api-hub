package codegenclient

import (
	"errors"
	"flag"
	"os"

	"github.com/pb33f/libopenapi"
	"github.com/pb33f/libopenapi/datamodel"
)

// Retrieve flags passed to app during initialisation
func NewAppFlags() AppFlags {
	appFlags := AppFlags{}
	flag.Var(&appFlags.InputFile, "input", "Path to OpenAPI 3+ specs file")
	flag.Parse()

	return appFlags
}

func NewOpenApiDocumentFromFile(file string) (*libopenapi.Document, error) {
	b, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	doc, err := libopenapi.NewDocumentWithConfiguration(b, &datamodel.DocumentConfiguration{
		AllowFileReferences:   true,
		AllowRemoteReferences: true,
	})
	if err != nil {
		return nil, err
	}

	return &doc, nil
}

func GetBaseUrlFromDocument(doc *libopenapi.Document) (string, error) {
	model, _ := (*doc).BuildV3Model()
	if model == nil {
		return "", errors.New("cannot create model")
	}
	return model.Model.Servers[0].URL, nil
}

func GetTitleFromDocument(doc *libopenapi.Document) (string, error) {
	model, _ := (*doc).BuildV3Model()
	if model == nil {
		return "", errors.New("cannot create model")
	}
	return model.Model.Info.Title, nil
}
