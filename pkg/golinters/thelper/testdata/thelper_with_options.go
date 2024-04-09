//golangcitest:args -Ethelper
//golangcitest:config_path testdata/thelper.yml
package testdata

import "testing"

func thelperWithHelperAfterAssignmentWO(t *testing.T) { // want "test helper function should start from t.Helper()"
	_ = 0
	t.Helper()
}

func thelperWithNotFirstWO(s string, t *testing.T, i int) { // want `parameter \*testing.T should be the first`
	t.Helper()
}

func thelperWithIncorrectNameWO(o *testing.T) {
	o.Helper()
}
