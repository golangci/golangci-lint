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
	switch d { // want "missing cases in switch of type Direction: East, West"
	case North, South:
	}
}
