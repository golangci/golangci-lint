//golangcitest:args -Eiface
//golangcitest:config_path testdata/unused.yml
package testdata

import "context"

type User struct {
	ID   string
	Name string
}

type UserRepository interface { // want "interface UserRepository is declared but not used within the package"
	UserOf(ctx context.Context, id string) (*User, error)
}

type UserRepositorySQL struct {
}

func (r *UserRepositorySQL) UserOf(ctx context.Context, id string) (*User, error) {
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

func Allow(x interface{}) {
	_ = x.(Allower)
}
