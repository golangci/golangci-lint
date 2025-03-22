package simple

//nolint:nonamedreturns // testing linter
func NewOtherMyStruct() (m *MyStruct) { // want "should be placed after the struct declaration"
	m = &MyStruct{Name: "John"}
	return
}

func NewMyStruct() *MyStruct { // want "should be placed after the struct declaration"
	return &MyStruct{Name: "John"}
}

func MustMyStruct() *MyStruct { // want `function \"MustMyStruct\" for struct \"MyStruct\" should be placed after the struct declaration`
	return NewMyStruct()
}

//nolint:recvcheck // testing linter
type MyStruct struct {
	Name string
}

//nolint:unused // testing linter
func (m MyStruct) lenName() int { // want `unexported method "lenName" for struct "MyStruct" should be placed after the exported method "GetName"`
	return len(m.Name)
}

func (m MyStruct) GetName() string {
	return m.Name
}

func (m *MyStruct) SetName(name string) {
	m.Name = name
}
