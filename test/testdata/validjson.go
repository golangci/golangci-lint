//golangcitest:args -Evalidjson
package testdata

type Key struct {
	A, B string
}

type StrAlias string

type ChanAlias chan<- int

type Textable [2]byte

func (Textable) MarshalText() ([]byte, error) {
	return nil, nil
}

func (*Textable) UnmarshalText(_ []byte) error {
	return nil
}

type ValidJsonTest struct {
	A int                           `json:"A"`
	B int                           `json:"B,omitempty"`
	C int                           `json:",omitempty"`
	D int                           `json:"-"`
	E chan<- int                    `json:"E"` // want "struct field has json tag but non-serializable type chan<- int"
	F chan<- int                    `json:"-"`
	G chan<- int                    `json:"-,"` // want "struct field has json tag but non-serializable type chan<- int"
	H ChanAlias                     `json:"H"`  // want "struct field has json tag but non-serializable type command-line-arguments.ChanAlias"
	I map[int]int                   `json:"I"`
	J map[int64]int                 `json:"J"`
	K map[Key]int                   `json:"K"` // want "struct field has json tag but non-serializable type"
	L map[struct{ A, B string }]int `json:"L"` // want "struct field has json tag but non-serializable type"
	M map[string]int                `json:"M"`
	N map[StrAlias]int              `json:"N"`
	O func()                        `json:"O"` // want "struct field has json tag but non-serializable type func()"
	P complex64                     `json:"P"` // want "struct field has json tag but non-serializable type complex64"
	Q complex128                    `json:"Q"` // want "struct field has json tag but non-serializable type complex128"
	R StrAlias                      `json:"R"`
	S map[Textable]int              `json:"S"`
}
