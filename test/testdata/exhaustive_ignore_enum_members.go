//golangcitest:args -Eexhaustive
//golangcitest:config_path testdata/configs/exhaustive_ignore_enum_members.yml
package testdata

type Direction int

const (
	North Direction = iota
	East
	South
	West
)

// Should only report East as missing because the enum member West is ignored
// using the ignore-enum-members setting.

func processDirectionIgnoreEnumMembers(d Direction) {
	switch d { // want "missing cases in switch of type Direction: East"
	case North, South:
	}
}
