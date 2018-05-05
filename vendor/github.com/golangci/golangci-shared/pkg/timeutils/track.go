package timeutils

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

func Track(now time.Time, format string, args ...interface{}) {
	logrus.Infof("[timing] %s took %s", fmt.Sprintf(format, args...), time.Since(now))
}
