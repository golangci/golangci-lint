//args: -Ewsl
//config: linters-settings.wsl.tests=1
package testdata

func main() {
	var (
		y = 0
	)
	if y < 1 { // ERROR "if statements should only be cuddled with assignments"

	}
}
