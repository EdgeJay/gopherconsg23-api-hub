package common

import "log"

func LogFatalError(message string, err interface{}) {
	log.Fatalf("%s: %v\n", message, err)
}
