//golangcitest:args -Enilnil
//golangcitest:config_path testdata/nilnil_multiple_nils.yml
package testdata

func withoutArgs()                                {}
func withoutError1() *User                        { return nil }
func withoutError2() (*User, *User)               { return nil, nil }
func withoutError3() (*User, *User, *User)        { return nil, nil, nil }
func withoutError4() (*User, *User, *User, *User) { return nil, nil, nil, nil }

func invalidOrder() (error, *User)               { return nil, nil }
func withError3rd() (*User, bool, error)         { return nil, false, nil }    // want "return both a `nil` error and an invalid value: use a sentinel error instead"
func withError4th() (*User, *User, *User, error) { return nil, nil, nil, nil } // want "return both a `nil` error and an invalid value: use a sentinel error instead"

type User struct{}
