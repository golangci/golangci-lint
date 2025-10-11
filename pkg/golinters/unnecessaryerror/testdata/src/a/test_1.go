package a

import (
	"errors"
)

// Not sure how this test case is different from [callchain1], [callchain2] or [callchain3],
// but it caused false positives while the above three passed.

func causeErr() error {
	return errors.New("")
}

func passErr() error {
	err := causeErr()
	if err != nil {
		return err
	}
	return nil
}

type Rec struct{}

func (r *Rec) Method() error {
	err := passErr()
	if err != nil {
		return err
	}
	return nil
}
