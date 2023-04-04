//golangcitest:args -Egovet
//golangcitest:expected_exitcode 0
package main

import (
	"bytes"
	"fmt"
	"os"
)

const (
	escape = "\x1b"
	reset  = escape + "[0m"
	green  = escape + "[32m"

	minValue = 0.0
	maxValue = 1.0
)

// Bar is a progress bar.
type Bar float64

var _ fmt.Formatter = Bar(maxValue)

// Format the progress bar as output
func (h Bar) Format(state fmt.State, r rune) {
	switch r {
	case 'r':
	default:
		panic(fmt.Sprintf("%v: unexpected format character", float64(h)))
	}

	if h > maxValue {
		h = maxValue
	}

	if h < minValue {
		h = minValue
	}

	if state.Flag('-') {
		h = maxValue - h
	}

	width, ok := state.Width()
	if !ok {
		// default width of 40
		width = 40
	}

	var pad int

	extra := len([]byte(green)) + len([]byte(reset))

	p := make([]byte, width+extra)
	p[0], p[len(p)-1] = '|', '|'
	pad += 2

	positive := int(Bar(width-pad) * h)
	negative := width - pad - positive

	n := 1
	n += copy(p[n:], green)
	n += copy(p[n:], bytes.Repeat([]byte("+"), positive))
	n += copy(p[n:], reset)

	if negative > 0 {
		copy(p[n:len(p)-1], bytes.Repeat([]byte("-"), negative))
	}

	_, _ = state.Write(p)
}

func main() {
	var b Bar = 0.9
	_, _ = fmt.Fprintf(os.Stdout, "%r", b)
}
