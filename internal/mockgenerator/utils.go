package mockgenerator

import (
	"flag"
	"os"
)

// Retrieve flags passed to app during initialisation
func NewAppFlags() AppFlags {
	appFlags := AppFlags{}

	flag.StringVar(&appFlags.InputFile, "input", "", "Path to OpenAPI 3+ specs file")
	flag.StringVar(&appFlags.Method, "method", "", "HTTP method")
	flag.IntVar(&appFlags.HttpStatus, "status", 0, "HTTP status")
	flag.StringVar(&appFlags.Path, "path", "", "Relative path of url")
	flag.Parse()

	return appFlags
}

func GetInputFileFromAppFlags(appFlags AppFlags) ([]byte, error) {
	return os.ReadFile(appFlags.InputFile)
}
