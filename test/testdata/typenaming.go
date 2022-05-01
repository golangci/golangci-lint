// args: -Etypenaming
package testdata

type UserType struct{} // ERROR "trim suffix `Type` from type name"

type userType string // ERROR "trim suffix `Type` from type name"
