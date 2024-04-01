//golangcitest:args -Etagalign
//golangcitest:config_path testdata/tagalign_order_only.yml
package testdata

import "time"

type TagAlignExampleOrderOnlyKO struct {
	Foo    time.Time `xml:"foo" json:"foo,omitempty" yaml:"foo" zip:"foo" gorm:"column:foo" validate:"required"`                // want `xml:"foo" json:"foo,omitempty" yaml:"foo" gorm:"column:foo" validate:"required" zip:"foo"`
	FooBar struct{}  `gorm:"column:fooBar" validate:"required" zip:"fooBar" xml:"fooBar" json:"fooBar,omitempty" yaml:"fooBar"` // want `xml:"fooBar" json:"fooBar,omitempty" yaml:"fooBar" gorm:"column:fooBar" validate:"required" zip:"fooBar"`
}

type TagAlignExampleOrderOnlyOK struct {
	Foo    time.Time `xml:"foo" json:"foo,omitempty" yaml:"foo" gorm:"column:foo" validate:"required" zip:"foo"`
	FooBar struct{}  `xml:"fooBar" json:"fooBar,omitempty" yaml:"fooBar" gorm:"column:fooBar" validate:"required" zip:"fooBar"`
}
