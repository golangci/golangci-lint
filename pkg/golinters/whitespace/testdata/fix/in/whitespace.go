//golangcitest:args -Ewhitespace
//golangcitest:config_path testdata/whitespace-fix.yml
//golangcitest:expected_exitcode 0
package p

import "fmt"

func oneLeadingNewline() {

	fmt.Println("Hello world")
}

func oneNewlineAtBothEnds() {

	fmt.Println("Hello world")

}

func noNewlineFunc() {
}

func oneNewlineFunc() {

}

func twoNewlinesFunc() {


}

func noNewlineWithCommentFunc() {
	// some comment
}

func oneTrailingNewlineWithCommentFunc() {
	// some comment

}

func oneLeadingNewlineWithCommentFunc() {

	// some comment
}

func twoLeadingNewlines() {


	fmt.Println("Hello world")
}

func multiFuncFunc(a int,
	b int) {
	fmt.Println("Hello world")
}

func multiIfFunc() {
	if 1 == 1 &&
		2 == 2 {
		fmt.Println("Hello multi-line world")
	}

	if true {
		if true {
			if true {
				if 1 == 1 &&
					2 == 2 {
						fmt.Println("Hello nested multi-line world")
				}
			}
		}
	}
}

func notGoFmted() {
         
        


         fmt.Println("Hello world")

			
               
}
