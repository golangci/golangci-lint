//golangcitest:args -Etagalign
//golangcitest:config_path testdata/configs/tagalign.yml
package testdata

type TagAlignSortExample struct {
	Foo    int `json:"foo,omitempty" yaml:"bar" xml:"baz" binding:"required" gorm:"column:foo" zip:"foo" validate:"required"`    // want `json:"foo,omitempty" yaml:"bar" xml:"baz" binding:"required" gorm:"column:foo" validate:"required" zip:"foo"`
	Bar    int `validate:"required"  yaml:"foo" xml:"bar" binding:"required" json:"bar,omitempty" gorm:"column:bar" zip:"bar" `  // want `json:"bar,omitempty" yaml:"foo" xml:"bar" binding:"required" gorm:"column:bar" validate:"required" zip:"bar"`
	FooBar int `gorm:"column:bar" validate:"required"   xml:"bar" binding:"required" json:"bar,omitempty"  zip:"bar" yaml:"foo"` // want `json:"bar,omitempty" yaml:"foo" xml:"bar" binding:"required" gorm:"column:bar" validate:"required" zip:"bar"`
}

type TagAlignSortExample2 struct {
	Foo int ` xml:"baz"  json:"foo,omitempty" yaml:"bar"  zip:"foo"  binding:"required" gorm:"column:foo"  validate:"required"` // want `tag is not aligned , should be: json:"foo,omitempty" yaml:"bar" xml:"baz" binding:"required" gorm:"column:foo" validate:"required" zip:"foo"`

	Bar int `validate:"required" gorm:"column:bar"  yaml:"foo" xml:"bar" binding:"required" json:"bar,omitempty" zip:"bar" ` // want `tag is not aligned , should be: json:"bar,omitempty" yaml:"foo" xml:"bar" binding:"required" gorm:"column:bar" validate:"required" zip:"bar"`
}
