//golangcitest:args -Emodernize
//golangcitest:expected_exitcode 0
package stditerators

import "go/types"

func _(tuple *types.Tuple) {
	for i := 0; i < tuple.Len(); i++ { // want "Len/At loop can simplified using Tuple.Variables iteration"
		print(tuple.At(i))
	}
}

func _(scope *types.Scope) {
	for i := 0; i < scope.NumChildren(); i++ { // want "NumChildren/Child loop can simplified using Scope.Children iteration"
		print(scope.Child(i))
	}
	{
		const child = 0                            // shadowing of preferred name at def
		for i := 0; i < scope.NumChildren(); i++ { // want "NumChildren/Child loop can simplified using Scope.Children iteration"
			print(scope.Child(i))
		}
	}
	{
		for i := 0; i < scope.NumChildren(); i++ {
			const child = 0 // nope: shadowing of fresh name at use
			print(scope.Child(i))
		}
	}
	{
		for i := 0; i < scope.NumChildren(); i++ { // want "NumChildren/Child loop can simplified using Scope.Children iteration"
			elem := scope.Child(i) // => preferred name = "elem"
			print(elem)
		}
	}
	{
		for i := 0; i < scope.NumChildren(); i++ { // want "NumChildren/Child loop can simplified using Scope.Children iteration"
			first := scope.Child(0) // the name heuristic should not be fooled by this
			print(first, scope.Child(i))
		}
	}
}

func _(union, union2 *types.Union) {
	for i := 0; i < union.Len(); i++ { // want "Len/Term loop can simplified using Union.Terms iteration"
		print(union.Term(i))
		print(union.Term(i))
	}
	for i := union.Len() - 1; i >= 0; i-- { // nope: wrong loop form
		print(union.Term(i))
	}
	for i := 0; i <= union.Len(); i++ { // nope: wrong loop form
		print(union.Term(i))
	}
	for i := 0; i <= union.Len(); i++ { // nope: use of i not in x.At(i)
		print(i, union.Term(i))
	}
	for i := 0; i <= union.Len(); i++ { // nope: x.At and x.Len have different receivers
		print(i, union2.Term(i))
	}
}
