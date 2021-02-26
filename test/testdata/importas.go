//args: -Eimportas
//config: linters-settings.importas.fff=fmt
package testdata

import (
	wrong_alias "fmt" // ERROR "import \"fmt\" imported as \"wrong_alias\" but must be \"fff\" according to config"
)

func ImportAsWrongAlias() {
	wrong_alias.Println("foo")
}
