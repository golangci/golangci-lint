//args: -Eexhaustive
//config_path: testdata/configs/exhaustive.yml
package testdata

type Direction int

const (
	North Direction = iota
	East
	South
	West
)

func processDirection(d Direction) {
	switch d { // ERROR "missing cases in switch of type Direction: East, West"
	case North, South:
	}
}

func processDirectionDefault(d Direction) {
	switch d {
	case North, South:
	default:
	}
}
