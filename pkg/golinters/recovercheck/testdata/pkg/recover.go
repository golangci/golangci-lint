package pkg

import "log"

func PanicRecover() func() {
	return func() {
		if r := recover(); r != nil {
			log.Println("Recovered from panic:", r)
		}
	}
}
