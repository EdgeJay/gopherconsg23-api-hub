package services

import (
	"context"

	"github.com/EdgeJay/gopherconsg23-api-hub/cmd/codegenmulticlient/housing"
	"github.com/EdgeJay/gopherconsg23-api-hub/internal/codegenmulticlient"
	"github.com/EdgeJay/gopherconsg23-api-hub/internal/common"
)

type HousingService struct {
	Config *codegenmulticlient.ServiceConfig
}

type PurchasedApartment struct {
	AptType      string
	CreatedOn    string
	PurchasedOn  string
	SalePeriodId string
}

func NewHousingService(config *codegenmulticlient.ServiceConfig) *HousingService {
	return &HousingService{
		Config: config,
	}
}

func (svc *HousingService) GetPurchasedApartment() *PurchasedApartment {

	client, err := housing.NewClient(svc.Config.BaseUrl)
	if err != nil {
		common.LogFatalError("Unable to create client", err)
	}

	res, err := client.GetPurchasedApartment(context.Background())
	if err != nil {
		common.LogFatalError("Unable to complete request", err)
	}

	record, err := housing.ParseGetPurchasedApartmentResponse(res)
	if err != nil {
		common.LogFatalError("Unable to parse response", err)
	}

	return &PurchasedApartment{
		AptType:      string(*record.JSON200.Data.AptType),
		CreatedOn:    *record.JSON200.Data.CreatedOn,
		PurchasedOn:  *record.JSON200.Data.PurchasedOn,
		SalePeriodId: *record.JSON200.Data.SalePeriodId,
	}
}
