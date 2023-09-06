package proto

func (x *Embedded) CustomMethod() interface{} {
	return nil
}

type Other struct {
}

func (x *Other) MyMethod(certs *Test) *Embedded {
	return nil
}
