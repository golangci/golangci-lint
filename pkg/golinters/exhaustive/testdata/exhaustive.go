//golangcitest:args -Eexhaustive
package testdata

type Direction int

const (
	North Direction = iota
	East
	South
	West
)

func processDirection(d Direction) {
	switch d { // want "switch not exhaustive: missing cases: East, West"
	case North, South:
	}
}
