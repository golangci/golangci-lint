//args: -Enilassign
package testdata

var num = 1

func pvar() {
	var ii *int
	*ii = 1 // ERROR "this assignment occurs invalid memory address or nil pointer dereference"

	var i *int
	i = &num // OK
	_ = i    // OK
}

func pstruct() {
	n := new(Node)

	*n.PVal = 1           // ERROR "this assignment occurs invalid memory address or nil pointer dereference"
	*n.ChildNode.PVal = 1 // ERROR "this assignment occurs invalid memory address or nil pointer dereference"

	n.ChildNode = &Node{PVal: &num} // OK
	n.PVal = &num                   // OK
}

type Node struct {
	PVal      *int
	ChildNode *Node
}
