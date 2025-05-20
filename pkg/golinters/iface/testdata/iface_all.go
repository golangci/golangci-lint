//golangcitest:args -Eiface
//golangcitest:config_path testdata/iface_all.yml
package testdata

import "fmt"

// identical

type Pinger interface { // want "identical: interface 'Pinger' contains identical methods or type constraints with another interface, causing redundancy"
	Ping() error
}

type Healthcheck interface { // want "identical: interface 'Healthcheck' contains identical methods or type constraints with another interface, causing redundancy"
	Ping() error
}

// opaque

type Server interface {
	Serve() error
}

type server struct {
	addr string
}

func (s server) Serve() error {
	return nil
}

func NewServer(addr string) Server { // want "opaque: 'NewServer' function return 'Server' interface at the 1st result, abstract a single concrete implementation of '\\*server'"
	return &server{addr: addr}
}

// unused

type User struct {
	ID   string
	Name string
}

type UserRepository interface { // want "unused: interface 'UserRepository' is declared but not used within the package"
	UserOf(id string) (*User, error)
}

type UserRepositorySQL struct {
}

func (r *UserRepositorySQL) UserOf(id string) (*User, error) {
	return nil, nil
}

type Granter interface {
	Grant(permission string) error
}

func AllowAll(g Granter) error {
	return g.Grant("all")
}

type Allower interface {
	Allow(permission string) error
}

func Allow(x any) {
	_ = x.(Allower)
	fmt.Println("allow")
}

// unexported

type unexportedReader interface {
	Read([]byte) (int, error)
}

func ReadAll(r unexportedReader) ([]byte, error) { // want "unexported interface 'unexportedReader' used as parameter in exported function 'ReadAll'"
	buf := make([]byte, 1024)
	_, err := r.Read(buf)
	return buf, err
}

func NewUnexportedReader() unexportedReader { // want "unexported interface 'unexportedReader' used as return value in exported function 'NewUnexportedReader'"
	return nil // stub
}
