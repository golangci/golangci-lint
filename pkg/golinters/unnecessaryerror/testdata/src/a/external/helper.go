// Helper package.
package external

func FuncPassedToExternal(fn func() error) error {
	return fn()
}

func ClosurePassedToExternal(fn func() error) error {
	return fn()
}

func AssignedInClosure(fn func() error) error {
	return fn()
}

type Exported struct{}

func (e *Exported) FuncPassedToExternalMethod(fn func() error) error {
	return fn()
}

type AssignedToExternalStructField1 struct {
	Err error
}

type AssignedToExternalStructField2 struct {
	Err error
}

type FuncAssignedToExternalStructField1 struct {
	Func func() error
}

type FuncAssignedToExternalStructField2 struct {
	Func func() error
}

type MethodAssignedToExternalStructField struct {
	Func func() error
}
