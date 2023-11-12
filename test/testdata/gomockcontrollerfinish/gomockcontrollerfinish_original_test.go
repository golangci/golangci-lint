//golangcitest:args -Egomockcontrollerfinish
package gomockcontrollerfinish

import (
	"testing"

	"github.com/golang/mock/gomock"
)

func TestFinishCall(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctrl.Finish() // want "calling Finish on gomock.Controller is no longer needed"
}

func TestFinishCallDefer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish() // want "calling Finish on gomock.Controller is no longer needed"
}

func TestFinishCallWithoutT(t *testing.T) {
	ctrl := gomock.NewController(nil)
	ctrl.Finish() // want "calling Finish on gomock.Controller is no longer needed"
}

func TestFinsihCallInAnotherFunction(t *testing.T) {
	ctrl := gomock.NewController(t)
	callFinish(ctrl)
}

func callFinish(ctrl *gomock.Controller) {
	ctrl.Finish() // want "calling Finish on gomock.Controller is no longer needed"
}

func TestNoFinishCall(t *testing.T) {
	gomock.NewController(t)
}
