//golangcitest:args -Etagalign
//golangcitest:config_path testdata/configs/tagalign_strict.yml
package testdata

import "time"

type TagAlignExampleStrictKO struct {
	Foo    time.Time `json:"foo,omitempty" validate:"required" zip:"foo"`                                                       // want `                     json:"foo,omitempty"    validate:"required"                            zip:"foo"`
	FooBar struct{}  `gorm:"column:fooBar" validate:"required" zip:"fooBar" xml:"fooBar" json:"fooBar,omitempty" yaml:"fooBar"` // want `gorm:"column:fooBar" json:"fooBar,omitempty" validate:"required" xml:"fooBar" yaml:"fooBar" zip:"fooBar"`
}
