/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package main

import (
	"os"
	"go_aggspotter/cmd"
	log "github.com/sirupsen/logrus"
)

func main() {
	loglevel := os.Getenv("LOG_LEVEL")

	switch loglevel {
	case "TRACE":
		log.SetLevel(log.TraceLevel | log.DebugLevel | log.InfoLevel | log.WarnLevel | log.ErrorLevel | log.FatalLevel | log.PanicLevel)
	case "DEBUG":
		log.SetLevel(log.DebugLevel | log.InfoLevel | log.WarnLevel | log.ErrorLevel | log.FatalLevel | log.PanicLevel)
	case "INFO":
		log.SetLevel(log.InfoLevel | log.WarnLevel | log.ErrorLevel | log.FatalLevel | log.PanicLevel)
	case "WARN":
		log.SetLevel(log.WarnLevel | log.ErrorLevel | log.FatalLevel | log.PanicLevel)
	case "ERROR":
		log.SetLevel(log.ErrorLevel | log.FatalLevel | log.PanicLevel)
	case "FATAL":
		log.SetLevel(log.FatalLevel | log.PanicLevel)
	case "PANIC":
		log.SetLevel(log.PanicLevel)
	default:
		log.SetLevel(log.InfoLevel | log.WarnLevel | log.ErrorLevel | log.FatalLevel | log.PanicLevel)
	}

	spotter.Execute()
}
