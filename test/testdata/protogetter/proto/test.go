package proto

func (x *Embedded) CustomMethod() interface{} {
	return nil
}

type Other struct {
}

func (x *Other) MyMethod(certs *Test) *Embedded {
	return nil
}

func (x *Test) Equal(v *Test) bool {
	return false
}

func (x *Embedded) SetS(s string) {
	x.S = s
}

func (x *Embedded) SetMap(_ map[string]string) {
}
