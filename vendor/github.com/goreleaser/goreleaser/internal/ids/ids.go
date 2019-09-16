// Package ids provides id validation code used my multiple pipes.
package ids

import "fmt"

// IDs is the IDs type
type IDs struct {
	ids  map[string]int
	kind string
}

// New IDs
func New(kind string) IDs {
	return IDs{
		ids:  map[string]int{},
		kind: kind,
	}
}

// Inc increment the counter of the given id
func (i IDs) Inc(id string) {
	i.ids[id]++
}

// Validate errors if there are any ids with counter > 1
func (i IDs) Validate() error {
	for id, count := range i.ids {
		if count > 1 {
			return fmt.Errorf(
				"found %d %s with the ID '%s', please fix your config",
				count, i.kind, id,
			)
		}
	}
	return nil
}
