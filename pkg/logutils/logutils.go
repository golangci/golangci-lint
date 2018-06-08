package logutils

import (
	"os"

	"github.com/sirupsen/logrus"
)

var isGolangCIRun = os.Getenv("GOLANGCI_COM_RUN") == "1"

func HiddenWarnf(format string, args ...interface{}) {
	if isGolangCIRun {
		logrus.Warnf(format, args...)
	} else {
		logrus.Infof(format, args...)
	}
}
