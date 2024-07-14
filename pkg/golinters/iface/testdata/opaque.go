//golangcitest:args -Eiface
//golangcitest:config_path testdata/opaque.yml
package testdata

type Server interface {
	Serve() error
}

type server struct {
	addr string
}

func (s server) Serve() error {
	return nil
}

func NewServer(addr string) Server { // want "NewServer function return Server interface at the 1st result, abstract a single concrete implementation of \\*server"
	return &server{addr: addr}
}
