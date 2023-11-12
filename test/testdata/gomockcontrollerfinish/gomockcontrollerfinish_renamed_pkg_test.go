//golangcitest:args -Egomockcontrollerfinish
package gomockcontrollerfinish

import (
	"testing"

	gomockRenamed "github.com/golang/mock/gomock"
)

func TestRenamedFinishCall(t *testing.T) {
	ctrl := gomockRenamed.NewController(t)
	ctrl.Finish() // want "calling Finish on gomock.Controller is no longer needed"
}

func TestRenamedFinishCallDefer(t *testing.T) {
	ctrl := gomockRenamed.NewController(t)
	defer ctrl.Finish() // want "calling Finish on gomock.Controller is no longer needed"
}

func TestRenamedFinishCallWithoutT(t *testing.T) {
	ctrl := gomockRenamed.NewController(nil)
	ctrl.Finish() // want "calling Finish on gomock.Controller is no longer needed"
}

func TestRenamedFinsihCallInAnotherFunction(t *testing.T) {
	ctrl := gomockRenamed.NewController(t)
	callFinish(ctrl)
}

func callFinish(ctrl *gomockRenamed.Controller) {
	ctrl.Finish() // want "calling Finish on gomock.Controller is no longer needed"
}

func TestNoFinishCall(t *testing.T) {
	gomockRenamed.NewController(t)
}
