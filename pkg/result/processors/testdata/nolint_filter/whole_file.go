//nolint:errcheck
package nolint_filter

func RetError() error {
	return nil
}

func MissedErrorCheck1() {
	RetErr()
}

func jo(v, w bool) bool {
	return v == w || v == w
}
