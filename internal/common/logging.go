package common

import "log"

func LogFatal(message string) {
	log.Fatalln(message)
}

func LogFatalError(message string, err interface{}) {
	log.Fatalf("%s: %v\n", message, err)
}
