package codegenmulticlient

import (
	"github.com/EdgeJay/gopherconsg23-api-hub/internal/codegenclient"
	"github.com/pb33f/libopenapi"
)

type ServicesMapping map[string]interface{}

type ServiceConfig struct {
	OpenApiFilePath string
	OpenApiDoc      *libopenapi.Document
	BaseUrl         string
	title           string
}

func NewServiceConfig(file string) *ServiceConfig {
	return &ServiceConfig{
		OpenApiFilePath: file,
	}
}

func (svc *ServiceConfig) Load() error {
	doc, err := codegenclient.NewOpenApiDocumentFromFile(svc.OpenApiFilePath)
	if err != nil {
		return err
	}

	svc.OpenApiDoc = doc

	url, err := codegenclient.GetBaseUrlFromDocument(svc.OpenApiDoc)
	if err != nil {
		return err
	}

	svc.BaseUrl = url

	title, err := codegenclient.GetTitleFromDocument(svc.OpenApiDoc)
	if err != nil {
		return err
	}

	svc.title = title

	return nil
}

func (svc *ServiceConfig) OpenApiDocTitle() string {
	return svc.title
}
