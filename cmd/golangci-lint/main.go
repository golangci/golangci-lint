package main

import (
	"log"

	"github.com/golangci/golangci-lint/internal/commands"
	"github.com/golangci/golangci-shared/pkg/analytics"
	"github.com/sirupsen/logrus"
)

func main() {
	log.SetFlags(0) // don't print time
	analytics.SetLogLevel(logrus.WarnLevel)

	e := commands.NewExecutor()
	if err := e.Execute(); err != nil {
		panic(err)
	}
}
