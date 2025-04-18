package swaggo

import "github.com/swaggo/swag"

const Name = "swaggo"

type Formatter struct {
	formatter *swag.Formatter
}

func New() *Formatter {
	return &Formatter{
		formatter: swag.NewFormatter(),
	}
}

func (*Formatter) Name() string {
	return Name
}

func (f *Formatter) Format(path string, src []byte) ([]byte, error) {
	return f.formatter.Format(path, src)
}
