package shell

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

type lineGetter struct {
	r *bufio.Scanner
}

func newLineGetter(a io.Reader) lineGetter {
	return lineGetter{
		r: bufio.NewScanner(a),
	}
}

func (lg lineGetter) Line() (string, error) {
	for lg.r.Scan() {
		line := strings.TrimSpace(lg.r.Text())
		if line == "" || line[0] == '#' {
			continue
		}
		return line, nil
	}
	err := lg.r.Err()
	if err == nil {
		err = io.EOF
	}
	return "", err
}

// Equal determines if two imputs are functionally the same shell scripts.

func Equal(a, b io.Reader) bool {
	s1 := newLineGetter(a)
	s2 := newLineGetter(b)
	for {
		line1, err1 := s1.Line()
		line2, err2 := s2.Line()
		if err1 == io.EOF && err2 == io.EOF {
			return true
		}
		if line1 != line2 || err1 != err2 {
			return false
		}
		// no errors, and lines are the same, continue
	}
}

func EqualBytes(a, b []byte) bool {
	if bytes.Equal(a, b) {
		return true
	}
	return Equal(bytes.NewReader(a), bytes.NewReader(b))
}

func EqualString(a, b string) bool {
	if a == b {
		return true
	}
	return Equal(strings.NewReader(a), strings.NewReader(b))
}
