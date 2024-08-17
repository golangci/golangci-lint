//golangcitest:args -Egobreakselectinfor
package testdata

func bad(ch <-chan string) {
	for {
		select {
		case <-ch:
			break // want "break statement inside select statement inside for loop"
		}
	}
}

func good(ch <-chan string) {
OUTER:
	for {
		select {
		case <-ch:
			break OUTER
		}
	}
}
