// args: -Ethelper
// config_path: testdata/configs/thelper.yml
package testdata

import "testing"

func thelperWithHelperAfterAssignmentWO(t *testing.T) { // ERROR "test helper function should start from t.Helper()"
	_ = 0
	t.Helper()
}

func thelperWithNotFirstWO(s string, t *testing.T, i int) { // ERROR `parameter \*testing.T should be the first`
	t.Helper()
}

func thelperWithIncorrectNameWO(o *testing.T) {
	o.Helper()
}
