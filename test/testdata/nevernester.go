//golangcitest:args -Enevernester
package testdata

import (
	"encoding/json"
	"fmt"
	"unicode/utf8"
)

func nestedFuncThatIsNested1Level() {
	fmt.Println("Terminator and Rambo are not the same")
}

func nestedFuncThatIsNested2Levels() {
	terminator := "Arnold Schwarzenegger"
	rambo := "Sylvester Stallone"
	if terminator != rambo {
		fmt.Println("Terminator and Rambo are not the same")
	}
}

func ifNesting2Levels() {
	byt := []byte(`{"num":6.13,"strs":["a","b"]}`)
	var dat map[string]interface{}
	if err := json.Unmarshal(
		byt, &dat); err != nil {
		panic(err)
	}
}

func nestedFuncThatIsNested3Levels() {
	terminator := "Arnold Schwarzenegger"
	rambo := "Sylvester Stallone"
	if terminator != rambo {
		if len(terminator) > 0 {
			fmt.Println("Terminator has bytes")
		}
	}
}

func rowBreak() {
	type row struct {
		row *row
	}

	r1 := row{}
	r2 := row{}
	r3 := row{}
	r4 := row{}

	r1.row = &r2
	r2.row = &r3
	r3.row = &r4

	r5 := r1.
		row.
		row.
		row
	if &r5 != nil {
		r6 := r5.
			row.
			row.
			row.
			row
		fmt.Printf("%v is not nil", r6)
	}
}

func nestedFuncThatIsNested4Levels() {
	terminator := "Arnold Schwarzenegger"
	rambo := "Sylvester Stallone"
	if utf8.RuneCountInString(terminator) > 0 {
		if len(terminator) > 0 {
			if terminator != rambo {
				fmt.Println("Terminator and Rambo are not the same")
			}
			if rambo != terminator {
				fmt.Println("Rambo and Terminator are not the same")
			}
		}
	}
}

func nestedFuncThatIsNested4LevelsWithElse() {
	terminator := "Arnold Schwarzenegger"
	rambo := "Sylvester Stallone"
	if utf8.RuneCountInString(terminator) > 0 {
		if len(terminator) > 0 {
			if len(rambo) > 0 {
				fmt.Println("Terminator and Rambo are not the same")
			} else {
				fmt.Println("Rambo is bytes long")
			}
		} else {
			fmt.Println("Terminator is bytes long")
		}
	} else {
		fmt.Println("Terminator is runes long")
	}
}

type Hero struct {
	Kind string
	Name string
}

type Identity struct {
}

func (i Identity) String() (string, error) {
	return "", nil
}

type source struct {
	Id Identity
}

type pair struct {
	source  *source
	compare *source
}

func findDiffs(pair pair) ([]Hero, error) {
	diffs := []Hero{}

	if pair.compare == nil {
		id, err := pair.source.Id.String()
		if err != nil {
			return []Hero{}, nil
		}
		return []Hero{
			{
				Kind: "kind.Good",
				Name: fmt.Sprintf("no id '%s' found", id),
			},
		}, nil
	}

	return diffs, nil
}

type row struct {
	property []string
}

func (r row) withProperty(s string) row {
	r.property = append(r.property, s)
	return r
}

func assignStatement() {
	r := row{}
	if len(r.property) > 0 {
		if len(r.property) > 0 {
			if len(r.property) > 0 {
				r.property = []string{
					"one",
					"two",
				}
			}
		}
	}
}

func exprStatement() {
	r := row{}
	r.withProperty("one")
	if len(r.property) > 0 {
		if len(r.property) > 0 {
			if len(r.property) > 0 {
				r.withProperty("two").
					withProperty("three").
					withProperty("four").
					withProperty("five").
					withProperty("six")
			}
		}
	}
}

func nestedFuncThatIsNested5Levels() { // want "calculated nesting for function nestedFuncThatIsNested5Levels is 5, max is 4"
	terminator := "Arnold Schwarzenegger"
	rambo := "Sylvester Stallone"
	if utf8.RuneCountInString(terminator) > 0 {
		if len(terminator) > 0 {
			if len(rambo) > 0 {
				if terminator != rambo {
					fmt.Println("Terminator and Rambo are not the same")
				}
			}
		}
	}
}

func nestedFuncThatIsNested5LevelsWithElse() { // want "calculated nesting for function nestedFuncThatIsNested5LevelsWithElse is 5, max is 4"
	terminator := "Arnold Schwarzenegger"
	rambo := "Sylvester Stallone"
	if utf8.RuneCountInString(terminator) > 0 {
		if len(terminator) > 0 {
			if len(rambo) > 0 {
				fmt.Println("Terminator and Rambo are not the same")
			} else {
				if utf8.RuneCountInString(rambo) > 0 {
					fmt.Println("Rambo is runes long")
				}
				fmt.Println("Rambo is bytes long")
			}
		} else {
			fmt.Println("Terminator is bytes long")
		}
	} else {
		fmt.Println("Terminator is runes long")
	}
}

type Order struct {
	Header OrderHeader
	Rows   []OrderRow
}

type OrderHeader struct {
	IsValid bool
}

type OrderRow struct {
	IsValid bool
	Price   int
}

func calculate(order *Order) int { // want "calculated nesting for function calculate is 5, max is 4"
	sum := 0
	if order != nil {
		if order.Header.IsValid {
			for _, row := range order.Rows {
				if row.IsValid {
					sum = sum + row.Price
				}
			}
		}
	}
	return sum
}

func calculate2(order *Order) int {
	sum := 0
	if order == nil {
		return sum
	}

	if order.Header.IsValid {
		sum = getPrice(order.Rows)
	}

	return sum
}

func getPrice(rows []OrderRow) int {
	sum := 0
	for _, row := range rows {
		if row.IsValid {
			sum = sum + row.Price
		}
	}
	return sum
}
