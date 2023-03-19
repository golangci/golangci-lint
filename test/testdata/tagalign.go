//golangcitest:args -Estylecheck
package testdata

type SortExample struct {
	Foo    int `json:"foo,omitempty" yaml:"bar" xml:"baz" binding:"required" gorm:"column:foo" zip:"foo" validate:"required"`    // want `tag is not aligned, should be: json:"foo,omitempty" yaml:"bar" xml:"baz" binding:"required" gorm:"column:foo" validate:"required" zip:"foo"`
	Bar    int `validate:"required"  yaml:"foo" xml:"bar" binding:"required" json:"bar,omitempty" gorm:"column:bar" zip:"bar" `  // want `tag is not aligned, should be: json:"bar,omitempty" yaml:"foo" xml:"bar" binding:"required" gorm:"column:bar" validate:"required" zip:"bar"`
	FooBar int `gorm:"column:bar" validate:"required"   xml:"bar" binding:"required" json:"bar,omitempty"  zip:"bar" yaml:"foo"` // want `tag is not aligned, should be: json:"bar,omitempty" yaml:"foo" xml:"bar" binding:"required" gorm:"column:bar" validate:"required" zip:"bar"`
}

type FooBar struct {
	Foo    int    `json:"foo"        validate:"required"`                    // want `tag is not aligned, should be: json:"foo"`
	Bar    string `json:"___bar___,omitempty"  validate:"required"`          // want `json:"___bar___,omitempty" validate:"required"`
	FooFoo int8   `json:"foo_foo"    validate:"required"      yaml:"fooFoo"` // want `tag is not aligned, should be: json:"foo_foo"`
	BarBar int    `json:"bar_bar"         validate:"required"`               // want `tag is not aligned, should be: json:"bar_bar"`
	FooBar struct {
		Foo    int    `json:"foo"    yaml:"foo"     validate:"required"` // want `tag is not aligned, should be: json:"foo"    yaml:"foo"          validate:"required"`
		Bar222 string `json:"bar222"   validate:"required"  yaml:"bar"`  // want `tag is not aligned, should be: json:"bar222" validate:"required" yaml:"bar"`
	} `json:"foo_bar" validate:"required"`
	FooFooFoo struct {
		BarBarBar struct {
			BarBarBarBar    string `json:"bar_bar_bar_bar" validate:"required"`                // want `json:"bar_bar_bar_bar"     validate:"required"`
			BarBarBarFooBar string `json:"bar_bar_bar_foo_bar" yaml:"bar" validate:"required"` // want `tag is not aligned, should be: json:"bar_bar_bar_foo_bar" yaml:"bar"          validate:"required"`
		} `json:"bar_bar_bar" validate:"required"`
	}
	BarFooBarFoo struct{}
	// test comment
	// test commnet 2
	BarFoo    string `json:"bar_foo" validate:"required"` // want `tag is not aligned, should be: json:"bar_foo"     validate:"required"`
	BarFooBar string `json:"bar_foo_bar" validate:"required"`
}
