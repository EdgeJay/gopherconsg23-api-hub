package mockapiserver

import (
	"encoding/json"
	"errors"
	"flag"
	"net/http"
	"os"
	"strconv"
)

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

func GetMockDataForRequest(mapping *MockDataMapping, httpStatus int, req *http.Request) (interface{}, error) {
	if str, ok := (*mapping)[req.Method+":"+strconv.Itoa(httpStatus)+":"+req.URL.Path]; ok {
		var v interface{}
		if err := json.Unmarshal([]byte(str), &v); err == nil {
			return v, nil
		}
	}
	return nil, errors.New("missing mock data")
}
