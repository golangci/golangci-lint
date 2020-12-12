package testifysubpkg

type StructA struct {
	Field1 int
}

func (s *StructA) Reset() {}

func (s *StructA) String() string { return "" }

func (s *StructA) ProtoMessage() {}
