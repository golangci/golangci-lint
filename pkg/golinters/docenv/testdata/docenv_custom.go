//golangcitest:args -Edocenv
//golangcitest:config_path testdata/docenv_custom.yml
package testdata

type Config struct {
	// Host name.
	Host string `foo:"HOST"`

	// Port number.
	Port int `foo:"PORT"`

	Undocumented string `foo:"UNDOCUMENTED"` // want "field `Undocumented` with `foo` tag should have a documentation comment"

	NotEnv string
}
