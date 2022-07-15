//golangcitest:args -Eexhaustive
//golangcitest:config_path testdata/configs/exhaustive_default.yml
package testdata

type Direction int

const (
	North Direction = iota
	East
	South
	West
)

// Should not report missing cases in the switch statement below even though
// some enum members (East, West) are not listed, because the switch statement
// has a 'default' case and the default-signifies-exhaustive setting is true.

func processDirectionDefault(d Direction) {
	switch d {
	case North, South:
	default:
	}
}
