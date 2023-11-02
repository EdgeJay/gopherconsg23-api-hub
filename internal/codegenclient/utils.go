package codegenclient

import "flag"

// Retrieve flags passed to app during initialisation
func NewAppFlags() AppFlags {
	appFlags := AppFlags{}
	flag.Var(&appFlags.InputFile, "input", "Path to OpenAPI 3+ specs file")
	flag.Parse()

	return appFlags
}
