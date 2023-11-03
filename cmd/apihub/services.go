package main

import (
	"context"
	"net/http"

	"github.com/EdgeJay/gopherconsg23-api-hub/internal/codegenclient"
	"github.com/EdgeJay/gopherconsg23-api-hub/internal/codegenmulticlient"
	"github.com/EdgeJay/gopherconsg23-api-hub/internal/common"
)

type ApiService struct {
	client    *http.Client
	Config    *codegenmulticlient.ServiceConfig
	InputFile string
}

func NewApiService(file string) *ApiService {
	return &ApiService{
		client:    codegenclient.NewHttpClient(),
		InputFile: file,
	}
}

func (svc *ApiService) Load() {
	// Load OpenAPI specs file
	svc.Config = codegenmulticlient.NewServiceConfig(svc.InputFile)
}

func (svc *ApiService) GetResidentData() interface{} {
	client, err := NewClient(svc.Config.BaseUrl)
	if err != nil {
		common.LogFatalError("Unable to create client", err)
	}

	params := &GetCombinedDataParams{
		Page: common.NewInt(1),
		Size: common.NewInt(1),
	}

	res, err := client.GetCombinedData(context.Background(), params)
	if err != nil {
		common.LogFatalError("Unable to complete request", err)
	}

	data, err := ParseGetCombinedDataResponse(res)
	if err != nil {
		common.LogFatalError("Unable to parse response", err)
	}

	return data.JSON200.Data
}
