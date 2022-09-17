//golangcitest:args -Epairedbrackets
package testdata

import (
	"fmt"
	"net/http"

	"github.com/stretchr/testify/assert"
	alias "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func _() {

	// good - empty
	fmt.Println()

	// good - one line, one element
	fmt.Println("xxx")

	// good - one line, several elements
	fmt.Printf("%s %d", "xxx", 10)

	// good - multiline
	fmt.Printf(
		"%s %d",
		"xxx",
		10,
	)

	// good - multiline, arguments are not validated (it should be different linter)
	fmt.Printf(
		"%s %d",
		"xxx", 10,
	)

	// good - last item exception
	http.HandleFunc("/api/v1", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// bad left - good right
	fmt.Printf("%s %d", // want `^left parenthesis should either be the last character on a line or be on the same line with the last argument$`
		"xxx", 10,
	)

	// bad left - right is ignored
	fmt.Printf("%s %d", // want `^left parenthesis should either be the last character on a line or be on the same line with the last argument$`
		"xxx", 10)

	// bad right - next
	fmt.Printf(
		"%s %d",
		"xxx",
		10) // want `^right parenthesis should be on the next line$`

	// bad right - previous, multiline
	http.HandleFunc("/api/v1", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	},
	) // want `^right parenthesis should be on the previous line$`

}

// default ignore-func-calls
func _() {

	assert.Equal(nil, []int{
		1,
	}, nil)

	alias.Equal(nil, []int{
		1,
	}, nil)

	require.Equalf(nil, []int{
		1,
	}, nil, "")

	assert.New(nil).JSONEq("", "", []int{
		1,
	}, "")

	require.New(nil).Eventually(func() bool {
		return false
	}, 0, 0)

}

func _() {

	// good - empty
	_ = []int{}

	// good - one line, one element
	_ = []int{1}

	// good - one line, several elements
	_ = []int{1, 2, 3}

	// good - multiline
	_ = []int{
		1,
		2,
		3,
	}

	// good - multiline, elements are not validated (it should be different linter)
	_ = []int{
		1, 2,
		3,
	}

	// good - last item exception
	_ = []any{1, 2, 3, `x
	y`}

	// bad left - good right
	_ = []int{1, 2, // want `^left brace should either be the last character on a line or be on the same line with the last composite element$`
		3,
	}

	// bad left - right is ignored
	_ = []int{1, 2, // want `^left brace should either be the last character on a line or be on the same line with the last composite element$`
		3}

	// bad right - next
	_ = []int{
		1, 2, 3} // want `^right brace should be on the next line$`

	// bad right - previous, one line
	// ../../testdata/no_go_fmt/ast_composite_lit.go

	// bad right - previous, multiline
	_ = []any{1, 2, 3, `x
	y`,
	} // want `^right brace should be on the previous line$`
}

func _() {

	type goodEmpty func()

	type goodOneLineOneParam func(int)

	type goodOneLineTwoParams func(int, string)

	type goodMultiline func(
		int,
		string,
	)

	type goodMultilineParamsNotValidated func(
		int, bool,
		string,
	)

	type goodLastItemException func(int, struct {
	})

	type badLeftGoodRight func(int, // want `^left parenthesis should either be the last character on a line or be on the same line with the last parameter$`
		string,
	)

	type badLeftRightIsIgnored func(int, // want `^left parenthesis should either be the last character on a line or be on the same line with the last parameter$`
		string)

	type badRightNext func(
		int,
		string) // want `^right parenthesis should be on the next line$`

	type badRightPreviousOneLine func(int, string,
	) // want `^right parenthesis should be on the previous line$`

	type badRightPreviousMultiline func(int, struct {
	},
	) // want `^right parenthesis should be on the previous line$`

}

func _() {

	type x[T, V, R any] struct{}

	// good - one line
	type _ x[int, string, bool]

	// good - multiline
	type _ x[
		int,
		string,
		bool,
	]

	// good - multiline, types are not validated (it should be different linter)
	type _ x[
		int, string,
		bool,
	]

	// good - last item exception
	type _ x[int, string, struct {
	}]

	// bad left - good right
	type _ x[int, string, // want `^left bracket should either be the last character on a line or be on the same line with the last element$`
		bool,
	]

	// bad left - right is ignored
	type _ x[int, string, // want `^left bracket should either be the last character on a line or be on the same line with the last element$`
		bool]

	// bad right - next
	type _ x[
		int, string,
		bool] // want `^right bracket should be on the next line$`

	// bad right - previous, one line
	// ../../testdata/no_go_fmt/ast_index_list_expr.go

	// bad right - previous, multiline
	type _ x[int, string, struct {
	},
	] // want `^right bracket should be on the previous line$`

}

func _() {

	// good - [ast.TypeSpec].TypeParams == nil
	type _ int

	// good - one line, one type parameter
	type _[T int] int

	// good - one line, several elements
	type _[T int, V string] int

	// good - multiline
	type _[
		T int,
		V string,
	] int

	// good - multiline, type parameters are not validated (it should be different linter)
	type _[
		T int, V string,
		R bool,
	] int

	// good - last item exception
	type _[T int, V struct {
	}] int

	// bad left - good right
	type _[T int, V string, // want `^left bracket should either be the last character on a line or be on the same line with the last type parameter$`
		R bool,
	] int

	// bad left - right is ignored
	type _[T int, V string, // want `^left bracket should either be the last character on a line or be on the same line with the last type parameter$`
		R bool] int

	// bad right - next
	type _[
		T int,
		V string] int // want `^right bracket should be on the next line$`

	// bad right - previous, one line
	type _[T int, V string,
	] int // want `^right bracket should be on the previous line$`

	// bad right - previous, multiline
	type _[T int, V struct {
	},
	] int // want `^right bracket should be on the previous line$`

}
