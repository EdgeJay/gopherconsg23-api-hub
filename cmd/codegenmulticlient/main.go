package main

import (
	"fmt"

	"github.com/EdgeJay/gopherconsg23-api-hub/cmd/codegenmulticlient/services"
	"github.com/EdgeJay/gopherconsg23-api-hub/internal/codegenclient"
)

const (
	SAVINGS_API_TITLE = "Savings API"
	HOUSING_API_TITLE = "Housing API"
)

func generateServicesFromFiles(files []string) *services.ServicesMapping {
	mapping := make(services.ServicesMapping)

	for _, file := range files {
		config := services.NewServiceConfig(file)
		config.Load()

		switch config.OpenApiDocTitle() {
		case SAVINGS_API_TITLE:
			mapping[SAVINGS_API_TITLE] = services.NewSavingsService(config)
		case HOUSING_API_TITLE:
			mapping[HOUSING_API_TITLE] = services.NewHousingService(config)
		}
	}

	return &mapping
}

func main() {
	appFlags := codegenclient.NewAppFlags()

	// Load OpenAPI specs files
	servicesMapping := generateServicesFromFiles(appFlags.InputFile)

	savingsService := (*servicesMapping)[SAVINGS_API_TITLE].(*services.SavingsService)
	savings := *savingsService.GetResidentSavingsRecord(1, 10)
	fmt.Println(*savings[0].Amount)

	housingService := (*servicesMapping)[HOUSING_API_TITLE].(*services.HousingService)
	purchasedApt := housingService.GetPurchasedApartment()
	fmt.Println(*purchasedApt)
}
