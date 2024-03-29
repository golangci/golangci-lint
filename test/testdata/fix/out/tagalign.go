//golangcitest:args -Etagalign
//golangcitest:config_path testdata/configs/tagalign_strict.yml
//golangcitest:expected_exitcode 0
package p

import "time"

type TagAlignExampleStrictKO struct {
	Foo    time.Time `                     json:"foo,omitempty"    validate:"required"                            zip:"foo"`
	FooBar struct{}  `gorm:"column:fooBar" json:"fooBar,omitempty" validate:"required" xml:"fooBar" yaml:"fooBar" zip:"fooBar"`
}
