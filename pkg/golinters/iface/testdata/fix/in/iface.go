//golangcitest:args -Eiface
//golangcitest:config_path testdata/iface_fix.yml
//golangcitest:expected_exitcode 0
package testdata

import "fmt"

// identical

type Pinger interface {
	Ping() error
}

type Healthcheck interface {
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

func NewServer(addr string) Server {
	return &server{addr: addr}
}

// unused

type User struct {
	ID   string
	Name string
}

type UserRepository interface {
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
