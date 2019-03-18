package unused

import "go/types"

// lookupMethod returns the index of and method with matching package and name, or (-1, nil).
func lookupMethod(T *types.Interface, pkg *types.Package, name string) (int, *types.Func) {
	if name != "_" {
		for i := 0; i < T.NumMethods(); i++ {
			m := T.Method(i)
			if sameId(m, pkg, name) {
				return i, m
			}
		}
	}
	return -1, nil
}

func sameId(obj types.Object, pkg *types.Package, name string) bool {
	// spec:
	// "Two identifiers are different if they are spelled differently,
	// or if they appear in different packages and are not exported.
	// Otherwise, they are the same."
	if name != obj.Name() {
		return false
	}
	// obj.Name == name
	if obj.Exported() {
		return true
	}
	// not exported, so packages must be the same (pkg == nil for
	// fields in Universe scope; this can only happen for types
	// introduced via Eval)
	if pkg == nil || obj.Pkg() == nil {
		return pkg == obj.Pkg()
	}
	// pkg != nil && obj.pkg != nil
	return pkg.Path() == obj.Pkg().Path()
}

func (c *Checker) implements(V types.Type, T *types.Interface) bool {
	// fast path for common case
	if T.Empty() {
		return true
	}

	if ityp, _ := V.Underlying().(*types.Interface); ityp != nil {
		for i := 0; i < T.NumMethods(); i++ {
			m := T.Method(i)
			_, obj := lookupMethod(ityp, m.Pkg(), m.Name())
			switch {
			case obj == nil:
				return false
			case !types.Identical(obj.Type(), m.Type()):
				return false
			}
		}
		return true
	}

	// A concrete type implements T if it implements all methods of T.
	ms := c.msCache.MethodSet(V)
	for i := 0; i < T.NumMethods(); i++ {
		m := T.Method(i)
		sel := ms.Lookup(m.Pkg(), m.Name())
		if sel == nil {
			return false
		}

		f, _ := sel.Obj().(*types.Func)
		if f == nil {
			return false
		}

		if !types.Identical(f.Type(), m.Type()) {
			return false
		}
	}
	return true
}
