//golangcitest:args -Ethelper
package testdata

import "testing"

func thelperWithHelperAfterAssignment(t *testing.T) { // want "test helper function should start from t.Helper()"
	_ = 0
	t.Helper()
}

func thelperWithNotFirst(s string, t *testing.T, i int) { // want `parameter \*testing.T should be the first`
	t.Helper()
}

func thelperWithIncorrectName(o *testing.T) { // want `parameter \*testing.T should have name t`
	o.Helper()
}

func bhelperWithHelperAfterAssignment(b *testing.B) { // want "test helper function should start from b.Helper()"
	_ = 0
	b.Helper()
}

func bhelperWithNotFirst(s string, b *testing.B, i int) { // want `parameter \*testing.B should be the first`
	b.Helper()
}

func bhelperWithIncorrectName(o *testing.B) { // want `parameter \*testing.B should have name b`
	o.Helper()
}

func tbhelperWithHelperAfterAssignment(tb testing.TB) { // want "test helper function should start from tb.Helper()"
	_ = 0
	tb.Helper()
}

func tbhelperWithNotFirst(s string, tb testing.TB, i int) { // want `parameter testing.TB should be the first`
	tb.Helper()
}

func tbhelperWithIncorrectName(o testing.TB) { // want `parameter testing.TB should have name tb`
	o.Helper()
}

func TestSubtestShouldNotBeChecked(t *testing.T) {
	testCases := []struct {
		desc string
	}{
		{
			desc: "example",
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			t.Error("test")
		})
	}
}
