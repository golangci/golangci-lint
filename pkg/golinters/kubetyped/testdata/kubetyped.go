//golangcitest:args -Ekubetyped
package testdata

func example() {
	m := map[string]any{ // want `use \*corev1\.Pod`
		"apiVersion": "v1",
		"kind":       "Pod",
		"metadata":   map[string]any{"name": "test"},
	}
	_ = m
}
