package mockapiserver

import (
	"flag"
	"log"
	"os"
)

func LogFatalError(message string, err interface{}) {
	log.Fatalf("%s: %v\n", message, err)
}

// Retrieve flags passed to app during initialisation
func NewAppFlags() AppFlags {
	appFlags := AppFlags{}

	flag.StringVar(&appFlags.InputFile, "input", "", "Path to OpenAPI 3+ specs file")
	flag.IntVar(&appFlags.Port, "port", 1337, "Port number for mock API server to listen to")
	flag.Parse()

	return appFlags
}

func GetInputFileFromAppFlags(appFlags AppFlags) ([]byte, error) {
	return os.ReadFile(appFlags.InputFile)
}
