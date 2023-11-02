package main

import (
	"fmt"
	"log"
	"os"

	"github.com/EdgeJay/gopherconsg23-api-hub/internal/codegenclient"
	"github.com/EdgeJay/gopherconsg23-api-hub/internal/common"
	"github.com/pb33f/libopenapi"
	validator "github.com/pb33f/libopenapi-validator"
)

func main() {
	appFlags := codegenclient.NewAppFlags()

	// Load OpenAPI specs files
	b, err := os.ReadFile(appFlags.InputFile[0])
	if err != nil {
		common.LogFatalError("Unable to load OpenAPI specs file", err)
	}

	doc, err := libopenapi.NewDocument(b)
	if err != nil {
		common.LogFatalError("Unable to parse OpenAPI specs file", err)
	}
	apiValidator, errs := validator.NewValidator(doc)
	if len(errs) > 0 {
		common.LogFatalError("Unable to create OpenAPI validator", err)
	}

	// Validate document
	docValid, validationErrs := apiValidator.ValidateDocument()
	if !docValid {
		for _, e := range validationErrs {
			log.Printf("Type: %s, Error: %s\n", e.ValidationType, e.Message)
			log.Printf("Reason: %s\n", e.Reason)

			for _, err := range e.SchemaValidationErrors {
				log.Printf("SchemaValidationError: %s, Location: %s\n", err.Reason, err.Location)
			}

			log.Printf("How to fix: %s\n\n", e.HowToFix)
		}
		common.LogFatal("OpenAPI specs file validation failed")
	}

	params := GetResidentSavingsRecordsParams{
		Page: common.NewInt(1),
		Size: common.NewInt(10),
	}

	// Build model and get server url
	model, errs := doc.BuildV3Model()
	if model == nil {
		common.LogFatalError("Model creation failed", errs)
	}
	baseUrl := model.Model.Servers[0].URL
	log.Printf("Server url found from file: %s\n", baseUrl)
	req, err := NewGetResidentSavingsRecordsRequest(baseUrl, &params)
	if err != nil {
		common.LogFatalError("Unable to create http request", err)
	}

	// Validate request
	requestValid, validationErrs := apiValidator.ValidateHttpRequest(req)
	if !requestValid {
		for _, e := range validationErrs {
			log.Println(e.Message)
		}
		common.LogFatal("http request validation failed")
	} else {
		log.Println("http request validation passed")
	}

	client := codegenclient.NewHttpClient()
	res, err := client.Do(req)
	if err != nil {
		common.LogFatalError("Unable to make http request", err)
	}

	// Validate response
	responseValid, validationErrs := apiValidator.ValidateHttpResponse(req, res)
	if !responseValid {
		for _, e := range validationErrs {
			log.Println(e.Message)

			for _, err := range e.SchemaValidationErrors {
				log.Printf("SchemaValidationError: %s, Location: %s\n", err.Reason, err.Location)
			}
		}
		common.LogFatal("http response validation failed")
	} else {
		log.Println("http response validation passed")
	}

	r, err := ParseGetResidentSavingsRecordsResponse(res)
	if err != nil {
		common.LogFatalError("Unable to parse http request", err)
	}

	rec := (*r.JSON200.Data)[0]
	fmt.Println(*rec.Amount, *rec.CreatedOn, *rec.RecordType)
}
