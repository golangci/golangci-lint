//golangcitest:args -Etagalign
//golangcitest:config_path testdata/configs/tagalign_sort_only.yml
package testdata

import "time"

type TagAlignExampleSortOnlyKO struct {
	Foo    time.Time `xml:"foo" json:"foo,omitempty" yaml:"foo" gorm:"column:foo" validate:"required" zip:"foo"`                // want `gorm:"column:foo" json:"foo,omitempty" validate:"required" xml:"foo" yaml:"foo" zip:"foo"`
	FooBar struct{}  `gorm:"column:fooBar" validate:"required" zip:"fooBar" xml:"fooBar" json:"fooBar,omitempty" yaml:"fooBar"` // want `gorm:"column:fooBar" json:"fooBar,omitempty" validate:"required" xml:"fooBar" yaml:"fooBar" zip:"fooBar"`
}

type TagAlignExampleSortOnlyOK struct {
	Foo    time.Time `gorm:"column:foo" json:"foo,omitempty" validate:"required" xml:"foo" yaml:"foo" zip:"foo"`
	FooBar struct{}  `gorm:"column:fooBar" json:"fooBar,omitempty" validate:"required" xml:"fooBar" yaml:"fooBar" zip:"fooBar"`
}
