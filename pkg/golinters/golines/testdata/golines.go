//golangcitest:config_path testdata/golines.yml
package testdata

import "fmt"

var (
	// want +1 "File is not properly formatted"
	abc = []string{"a really long string", "another really long string", "a third really long string", "a fourth really long string", fmt.Sprintf("%s %s %s %s >>>>> %s %s", "first argument", "second argument", "third argument", "fourth argument", "fifth argument", "sixth argument")}
)

type MyStruct struct {
	aLongProperty       int `help:"This is a really long string for this property"`
	anotherLongProperty int `help:"This is a really long string for this property, part 2"`
	athirdLongProperty  int `help:"This is a really long string for this property, part 3...."`
}

type MyInterface interface {
	aReallyLongFunctionName(argument1 string, argument2 string, argument3 string, argument4 string, argument5 string, argument6 string) (string, error)
}

// Something here

// Another comment
// A third comment
// This is a really really long comment that needs to be split up into multiple lines. I don't know how easy it will be to do, but I think we can do it!
func longLine(aReallyLongName string, anotherLongName string, aThirdLongName string) (string, error) {
	argument1 := "argument1"
	argument2 := "argument2"
	argument3 := "argument3"
	argument4 := "argument4"

	fmt.Printf("This is a really long string with a bunch of arguments: %s %s %s %s >>>>>>>>>>>>>>>>>>>>>>", argument1, argument2, argument3, argument4)
	fmt.Printf("This is a short statement: %d %d %d", 1, 2, 3)

	z := argument1 + argument2 + fmt.Sprintf("This is a really long statement that should be broken up %s %s %s", argument1, argument2, argument3)

	fmt.Printf("This is a really long line that can be broken up twice %s %s", fmt.Sprintf("This is a really long sub-line that should be broken up more because %s %s", argument1, argument2), fmt.Sprintf("A short one %d", 3))

	fmt.Print("This is a function with a really long single argument. We want to see if it's properly split")

	fmt.Println(z)

	// This is a really long comment on an indented line. Do you think we can split it up or should we just leave it as is?
	if argument4 == "5" {
		return "", fmt.Errorf("a very long query with ID %d failed. Check Query History in AWS UI", 12341251)
	}

	go func() {
		fmt.Printf("This is a really long line inside of a go routine call. It should be split if at all possible.")
	}()

	if "hello this is a big string" == "this is a small string" && "this is another big string" == "this is an even bigger string >>>" {
		fmt.Print("inside if statement")
	}

	fmt.Println(map[string]string{"key1": "a very long value", "key2": "a very long value", "key3": "another very long value"})

	return "", nil
}

func shortFunc(a int, b int) error {
	c := make(chan int)

	for {
		select {
		case <-c:
			switch a {
			case 1:
				return fmt.Errorf("This is a really long line that can be broken up twice %s %s", fmt.Sprintf("This is a really long sub-line that should be broken up more because %s %s", "xxxx", "yyyy"), fmt.Sprintf("A short one %d", 3))
			case 2:
			}
		}

		break
	}

	if a > 5 {
		panic(fmt.Sprintf(">>>>>>>>>>>>>>>>>>> %s %s %s %s", "really long argument", "another really long argument", "a third really long arguement", abc[1:2]))
	}

	return nil
	// This is an end decoration
}
