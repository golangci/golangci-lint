//golangcitest:args -Efuncorder
//golangcitest:config_path testdata/funcorder_disable_constructor_alphabetical.yml
package testdata

import "time"

type MyStruct struct {
	Name string
}

func (m MyStruct) lenName() int { // want `unexported method "lenName" for struct "MyStruct" should be placed after the exported method "GetName"`
	return len(m.Name)
}

func (m *MyStruct) SetName(name string) {
	m.Name = name
}

func (m MyStruct) GetName() string { // want `"GetName" for struct "MyStruct" should be placed before method "SetName"`
	return m.Name
}

type MyStruct2 struct {
	Name string
}

func (m *MyStruct2) SetName(name string) {
	m.Name = name
}

func (m MyStruct2) GetName() string { // want `"GetName" for struct "MyStruct2" should be placed before method "SetName"`
	return m.Name
}

func NewMyStruct2() *MyStruct2 {
	return &MyStruct2{Name: "John"}
}

func NewTime() time.Time {
	return time.Now()
}
