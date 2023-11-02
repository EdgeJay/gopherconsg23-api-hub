package services

import (
	"net/http"

	"github.com/EdgeJay/gopherconsg23-api-hub/cmd/codegenmulticlient/savings"
	"github.com/EdgeJay/gopherconsg23-api-hub/internal/codegenclient"
	"github.com/EdgeJay/gopherconsg23-api-hub/internal/codegenmulticlient"
	"github.com/EdgeJay/gopherconsg23-api-hub/internal/common"
)

type SavingsService struct {
	client *http.Client
	Config *codegenmulticlient.ServiceConfig
}

func NewSavingsService(config *codegenmulticlient.ServiceConfig) *SavingsService {
	return &SavingsService{
		client: codegenclient.NewHttpClient(),
		Config: config,
	}
}

func (svc *SavingsService) GetResidentSavingsRecord(page, size int) *[]savings.Savings {
	params := savings.GetResidentSavingsRecordsParams{
		Page: common.NewInt(page),
		Size: common.NewInt(size),
	}

	req, err := savings.NewGetResidentSavingsRecordsRequest(svc.Config.BaseUrl, &params)
	if err != nil {
		common.LogFatalError("Unable to create http request", err)
	}

	res, err := svc.client.Do(req)
	if err != nil {
		common.LogFatalError("Unable to make http request", err)
	}

	records, err := savings.ParseGetResidentSavingsRecordsResponse(res)
	if err != nil {
		common.LogFatalError("Unable to parse http request", err)
	}

	return records.JSON200.Data
}
