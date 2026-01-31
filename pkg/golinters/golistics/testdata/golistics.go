//golangcitest:args -Egolistics
package ast

import "image/color"

type Dimensional map[string]uint

type (
	Display struct {
		Outside string
		Inside  string
	}
	Border struct {
		Color any
		Style any
		Width any
	}
	BorderRadiuses struct {
		TopLeft, TopRight, BottomRight, BottomLeft any
	}
	Borders struct {
		Top, Right, Bottom, Left Border
	}
	Margin struct {
		Top, Right, Bottom, Left any
	}
	Padding struct {
		Top, Right, Bottom, Left any
	}
	Font struct {
		Family any
		Size   any
		Weight any
	}
	Text struct {
		Color         any
		LineHeight    any
		TextAlignment any
	}
	Dimensions struct {
		Height     any
		Width      any
		unexported color.RGBA
	}

	Styles struct {
		Dimensions      Dimensions
		Margin          Margin
		Padding         Padding
		Display         Display
		Text            Text
		Font            Font
		Border          Borders
		BorderRadiuses  BorderRadiuses
		BackgroundColor any
	}
	Rule struct {
		Selector string
		Styles   *Styles
	}
)

// implementation is not needed for test purposes
func collect(_ map[string]any) []string { return nil }
func tree(_ ...any) string              { return "" }
func safeEq(a, b any) bool              { return a == b }

//golistics:exported
func (s Text) Strings() []string {
	return collect(map[string]any{
		"Color":         s.Color,
		"LineHeight":    s.LineHeight,
		"TextAlignment": s.TextAlignment,
	})
}

//golistics:exported // want `missing fields: Height, Width`
func (Dimensions) Strings() []string {
	return nil
}

//golistics:exported
func (s *Styles) Strings() []string {
	return collect(map[string]any{
		"Dimensions":      s.Dimensions,
		"Margin":          s.Margin,
		"Padding":         s.Padding,
		"Display":         s.Display,
		"Text":            s.Text,
		"Font":            s.Font,
		"Border":          s.Border,
		"BorderRadiuses":  s.BorderRadiuses,
		"BackgroundColor": s.BackgroundColor,
	})
}

//golistics:exported
func (r Rule) String() string {
	return tree(r.Selector, r.Styles.Strings())
}

//golistics:all // want `missing fields: Bottom, Left, Right, Top`
func (Borders) IsEqual(Borders) bool {
	return false
}

//golistics:all // want `missing field: Top`
func (s Margin) IsEqual(y Margin) bool {
	return safeEq(s.Right, y.Right) &&
		safeEq(s.Bottom, y.Bottom) &&
		safeEq(s.Left, y.Left)
}

//golistics:exported
func (s Dimensions) IsEqual(y Dimensions) bool {
	return safeEq(s.Height, y.Height) &&
		safeEq(s.Width, y.Width)
}

// Suppress is exist to suppress linter error on unused & unexported
// field that is needed in [Dimensions.Strings] for testing linter behavior
func (d Dimensions) Suppress() {
	_ = d.unexported
}
