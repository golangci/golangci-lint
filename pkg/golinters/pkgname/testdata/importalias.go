//golangcitest:args -Epkgname
//golangcitest:config_path importalias.yml
package hello

import (
	go_format "fmt" // want "found import 'fmt' with alias 'go_format', should not use under_score in package alias name"
	"log"
	structLog "log/slog" // want "found import 'log/slog' with alias 'structLog', should not use mixedCaps in package alias name"
)

func Hello() {
	go_format.Println("Hello, World!")
	structLog.Info("Hello, World!")
	log.Println("Hello, World!")
}
