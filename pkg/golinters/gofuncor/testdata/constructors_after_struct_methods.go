package simple

//nolint:recvcheck // testing linter
type MyStruct2 struct {
	Name string
}

func (m MyStruct2) GetName() string {
	return m.Name
}

func (m *MyStruct2) SetName(name string) {
	m.Name = name
}

//nolint:nonamedreturns // testing linter
func NewOtherMyStruct2() (m *MyStruct2) { // want `constructor 'NewOtherMyStruct2' for struct 'MyStruct2' should be placed before struct method 'GetName'`
	m = &MyStruct2{Name: "John"}
	return
}

func NewMyStruct2() *MyStruct2 { // want `constructor \"NewMyStruct2\" for struct \"MyStruct2\" should be placed before struct method \"GetName\"`
	return &MyStruct2{Name: "John"}
}
