package testdata

type GomutcheckStruct struct {
	TestField string
}

func (s *GomutcheckStruct) GomutcheckMutateWithPointer() {
	s.TestField = "new value"
}

func (s GomutcheckStruct) GomutcheckCallWithPointerMethod() {
	sPtr := &s
	sPtr.GomutcheckMutateWithPointer()
}

func (s GomutcheckStruct) GomutcheckCreateNewInstance() {
	newStruct := GomutcheckStruct{TestField: "new instance"}
	_ = newStruct
}

func (s GomutcheckStruct) GomutcheckMutateField() {
	s.TestField = "new value" // want `struct field 'TestField' is being mutated in value receiver method`
}

func (s GomutcheckStruct) GomutcheckReadOnly() {
	_ = s.TestField
}

func (s GomutcheckStruct) GomutcheckCorrectUsage() {
	println(s.TestField)
}
