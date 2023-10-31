package mockcodegenclient

import "flag"

// Retrieve flags passed to app during initialisation
func NewAppFlags() AppFlags {
	appFlags := AppFlags{}

	flag.StringVar(&appFlags.InputFile, "input", "", "Path to OpenAPI 3+ specs file")
	flag.Parse()

	return appFlags
}
