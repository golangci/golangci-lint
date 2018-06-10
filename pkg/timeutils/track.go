package timeutils

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

func Track(from time.Time, format string, args ...interface{}) {
	logrus.Infof("%s took %s", fmt.Sprintf(format, args...), time.Since(from))
}
