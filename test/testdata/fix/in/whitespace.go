//args: -Ewhitespace
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
