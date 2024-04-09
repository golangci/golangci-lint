//golangcitest:args -Etagalign
//golangcitest:config_path testdata/tagalign_align_only.yml
package testdata

import "time"

type TagAlignExampleAlignOnlyKO struct {
	Foo    time.Time `gorm:"column:foo" json:"foo,omitempty" xml:"foo" yaml:"foo" zip:"foo"`                // want `gorm:"column:foo"    json:"foo,omitempty" xml:"foo"               yaml:"foo"   zip:"foo"`
	FooBar struct{}  `gorm:"column:fooBar" zip:"fooBar" json:"fooBar,omitempty" xml:"fooBar" yaml:"fooBar"` // want  `gorm:"column:fooBar" zip:"fooBar"         json:"fooBar,omitempty" xml:"fooBar" yaml:"fooBar"`
	FooFoo struct {
		Foo    int    `json:"foo" yaml:"foo"`   // want `json:"foo"    yaml:"foo"`
		Bar    int    `yaml:"bar"   json:"bar"` // want `yaml:"bar"    json:"bar"`
		BarBar string `json:"barBar" yaml:"barBar"`
	} `xml:"fooFoo" json:"fooFoo"`
	NoTag  struct{}
	BarBar struct{} `json:"barBar,omitempty" gorm:"column:barBar" yaml:"barBar" xml:"barBar" zip:"barBar"`
	Boo    struct{} `gorm:"column:boo" json:"boo,omitempty" xml:"boo" yaml:"boo" zip:"boo"` // want `gorm:"column:boo"       json:"boo,omitempty" xml:"boo"     yaml:"boo"   zip:"boo"`
}

type TagAlignExampleAlignOnlyOK struct {
	Foo    time.Time `gorm:"column:foo"    json:"foo,omitempty" xml:"foo"               yaml:"foo"   zip:"foo"`
	FooBar struct{}  `gorm:"column:fooBar" zip:"fooBar"         json:"fooBar,omitempty" xml:"fooBar" yaml:"fooBar"`
	FooFoo struct {
		Foo    int    `json:"foo"    yaml:"foo"`
		Bar    int    `yaml:"bar"    json:"bar"`
		BarBar string `json:"barBar" yaml:"barBar"`
	} `xml:"fooFoo" json:"fooFoo"`
	NoTag  struct{}
	BarBar struct{} `json:"barBar,omitempty" gorm:"column:barBar" yaml:"barBar" xml:"barBar" zip:"barBar"`
	Boo    struct{} `gorm:"column:boo"       json:"boo,omitempty" xml:"boo"     yaml:"boo"   zip:"boo"`
}
