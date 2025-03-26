//golangcitest:args -Egofuncor
package testdata

import "time"

type MyStruct struct {
	Name string
}

func (m MyStruct) lenName() int { // want `unexported method "lenName" for struct "MyStruct" should be placed after the exported method "SetName"`
	return len(m.Name)
}

func (m MyStruct) GetName() string {
	return m.Name
}

func (m *MyStruct) SetName(name string) {
	m.Name = name
}

type MyStruct2 struct {
	Name string
}

func (m MyStruct2) GetName() string {
	return m.Name
}

func (m *MyStruct2) SetName(name string) {
	m.Name = name
}

func NewMyStruct2() *MyStruct2 { // want `constructor "NewMyStruct2" for struct "MyStruct2" should be placed before struct method "GetName"`
	return &MyStruct2{Name: "John"}
}

func NewTime() time.Time {
	return time.Now()
}
