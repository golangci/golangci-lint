//golangcitest:args -Egobreakselectinfor
package testdata

func bad() {
	var ch chan string
	for {
		select {
		case <-ch:
			break // want "break statement inside select statement inside for loop"
		}
	}
}
