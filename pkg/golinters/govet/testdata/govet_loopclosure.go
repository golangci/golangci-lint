//golangcitest:args -Egovet
//golangcitest:config_path testdata/govet_loopclosure.yml
//golangcitest:expected_exitcode 0
package testdata

func InGo22(l []int) {
	for i, v := range l {
		go func() {
			print(i) // Not reported.
			print(v) // Not reported.
		}()
	}
}
