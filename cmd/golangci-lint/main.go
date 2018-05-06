package main

import (
	"log"
	"runtime"

	"github.com/golangci/golangci-lint/internal/commands"
	"github.com/golangci/golangci-shared/pkg/analytics"
	"github.com/sirupsen/logrus"
)

func main() {
	log.SetFlags(0) // don't print time
	analytics.SetLogLevel(logrus.WarnLevel)
	runtime.GOMAXPROCS(runtime.NumCPU())

	e := commands.NewExecutor()
	if err := e.Execute(); err != nil {
		panic(err)
	}
}
