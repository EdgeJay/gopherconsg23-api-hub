package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"

	"github.com/EdgeJay/gopherconsg23-api-hub/internal/common"
	"github.com/EdgeJay/gopherconsg23-api-hub/internal/mockapiserver"
	"github.com/EdgeJay/gopherconsg23-api-hub/internal/mockgenerator"
)

func main() {

	appFlags := mockgenerator.NewAppFlags()

	v3Model := mockapiserver.NewOpenAPIV3ModelFromFile(mockapiserver.AppFlags{
		InputFile: appFlags.InputFile,
	})

	mockDataMapping := mockapiserver.NewMockDataMapping(v3Model, false)

	// Get mock data
	req := &http.Request{
		Method: appFlags.Method,
		URL: &url.URL{
			Path: appFlags.Path,
		},
	}

	if mockData, err := mockapiserver.GetMockDataForRequest(mockDataMapping, http.StatusOK, req); err == nil {
		if b, err := json.MarshalIndent(mockData, "", "  "); err == nil {
			log.Printf("%v\n", string(b))
		} else {
			common.LogFatalError("Unable to generate mock data", err)
		}
	} else {
		common.LogFatalError("Unable to retrieve mock data", err)
	}
}
