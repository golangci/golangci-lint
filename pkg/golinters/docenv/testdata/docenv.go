//golangcitest:args -Edocenv
package testdata

type Config struct {
	// Host name.
	Host string `env:"HOST"`

	// Port number.
	Port int `env:"PORT"`

	Undocumented string `env:"UNDOCUMENTED"` // want "field `Undocumented` with `env` tag should have a documentation comment"

	NotEnv string
}
