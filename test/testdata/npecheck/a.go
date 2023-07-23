//golangcitest:args -Enpecheck
package a

import (
	"fmt"
)

type Node struct {
	A int
	B *ChildNode
	C *float32
	D ChildNode
}

type ChildNode struct {
	Score int32
	GNode *GrandsonNode
}

type ChildBrotherNode struct {
	Score int32
	GNode *GrandsonNode
}

type GrandsonNode struct {
	Age  int32
	High *int32
}

type DataInfo struct {
	A *Node
}

func (d *DataInfo) printDataInfo() {
	fmt.Println(d)
}

func (d *DataInfo) GetChildNodePtr() *ChildNode {
	return nil
}

func (d *DataInfo) GetChildNodeNonPtr() ChildBrotherNode {
	return ChildBrotherNode{}
}

func (c *ChildNode) PrintScore() int {
	fmt.Println(c.Score)
	return 0
}

func (b ChildBrotherNode) PrintScore() int {
	fmt.Println(b.Score)
	return 0
}

func (b ChildBrotherNode) GetGrandsonNodePtr() *GrandsonNode {
	return nil
}

// Function parameter is pointer, and its variable is directly referenced without validation
func np1Example(d *DataInfo) {
	fmt.Println(d.A) // want "potential nil pointer reference"

	// d may be is a nil pointer, you had better to check it before reference.
	// such as :
	// if d != nil {
	//	fmt.Println(d.A)
	// }

	// Or:
	// if d == nil {
	//	return
	// }
	// fmt.Println(d.A)

	// Otherwise it is potential nil pointer reference, sometimes it's unexpected disaster
}

// Function parameter is pointer, and its variables are directly referenced in a chain without validation
func np2Example(d *DataInfo) {
	fmt.Println(d.A.B) // want "potential nil pointer reference"

	// d is a potential nil pointer
	// d.A is also a potential nil pointer
	// You can follow the writing below, and will be more safe:

	// if d != nil && d.A != nil {
	//		fmt.Println(d.A.B)
	//}

	// Or:
	// if d == nil {
	//	 return
	// }
	//
	// if d.A == nil {
	//	 return
	// }
	//
	// fmt.Println(d.A.B)
}

// A pointer variable obtained by calling an external function, unverified, directly referenced
func np3Example() {
	d := GetDataInfo() // d is a pointer obtained by calling an external function
	fmt.Println(d.A)   // want "potential nil pointer reference"

	// d may be is a nil pointer, you had better to check it before reference.
	// such as :
	// if d != nil {
	//	fmt.Println(d.A)
	// }

	// Or:
	// if d == nil {
	//	return
	// }
	// fmt.Println(d.A)

	// Otherwise it is potential nil pointer reference, sometimes it's unexpected disaster
}

// A pointer variable obtained by calling an external function, unverified, directly chain referenced
func np4Example() {
	d := GetDataInfo()
	fmt.Println(d.A.B) // want "potential nil pointer reference"

	// d is a potential nil pointer
	// d.A is also a potential nil pointer
	// You can follow the writing below, and will be more safe:

	// if d != nil && d.A != nil {
	//		fmt.Println(d.A.B)
	//}

	// Or:
	// if d == nil {
	//	 return
	// }
	//
	// if d.A == nil {
	//	 return
	// }
	//
	// fmt.Println(d.A.B)
}

// Function input parameter is slice including pointer, and their elements are directly referenced without validation
func np5Example(infoList []*Node) {
	for _, info := range infoList {
		fmt.Println(info.A) // want "potential nil pointer reference"
		// info is a potential nil pointer
		// It can be written as follows, and will be more safe.

		// if info != nil {
		// 	fmt.Println(info.A)
		// }

		// Or:
		// if info == nil {
		// 	  continue
		// }
		// fmt.Println(info.A)
	}
}

// An slice including pointers obtained by a function, whose pointer elements are not checked and are directly referenced
func np6Example() {
	infoList := GetDataInfoList()
	for _, info := range infoList {
		fmt.Println(info.A) // want "potential nil pointer reference"

		// info is a potential nil pointer
		// It can be written as follows, and will be more safe.

		// if info != nil {
		// 	fmt.Println(info.A)
		// }

		// Or:
		// if info == nil {
		// 	  continue
		// }
		// 	fmt.Println(info.A)
	}
}

// Function parameter is pointer, and its method is directly referenced without validation
func np7Example1(d *DataInfo) {
	d.printDataInfo() // want "potential nil pointer reference"

	// d is a potential nil pointer
	// It can be written as follows, and will be more safe.
	// if d != nil {
	// 	d.printDataInfo()
	// }

	// Or:
	// if d == nil {
	// 	 return
	// }
	// d.printDataInfo()
}

// A pointer variable obtained by calling an external function, unverified, directly reference its method
func np7Example2() {
	d := GetDataInfo()
	d.printDataInfo() // want "potential nil pointer reference"

	// d is a potential nil pointer
	// It can be written as follows, and will be more safe.
	// if d != nil {
	// 	d.printDataInfo()
	// }

	// Or:
	// if d == nil {
	// 	 return
	// }
	// d.printDataInfo()
}

// Function parameter is pointer, and its method is directly referenced in chain without validation
func np8Example(d *DataInfo) {
	_ = d.GetChildNodePtr().PrintScore() // want "potential nil pointer reference"

	// d is a potential nil pointer reference
	// d.GetChildNodePtr() is also a potential nil pointer

	// It can be written as follows, and will be more safe.
	// if d != nil && d.GetChildNodePtr() != nil {
	// 	_ = d.GetChildNodePtr().PrintScore()
	// }

	// Or:
	// if d == nil {
	// 	 return
	// }
	//
	// if d.GetChildNodePtr() == nil {
	//	 return
	// }
	//
	// _ = d.GetChildNodePtr().PrintScore()
}

// A pointer variable obtained by calling an external function, and its method is directly referenced in chain without validation
func np9Example() {
	d := GetDataInfo()
	_ = d.GetChildNodePtr().PrintScore() // want "potential nil pointer reference"

	// d is a potential nil pointer reference
	// d.GetChildNodePtr() is also a potential nil pointer

	// It can be written as follows, and will be more safe.
	// if d != nil && d.GetChildNodePtr() != nil {
	// 	_ = d.GetChildNodePtr().PrintScore()
	// }

	// Or:
	// if d == nil {
	// 	 return
	// }
	//
	// if d.GetChildNodePtr() == nil {
	//	 return
	// }
	//
	// _ = d.GetChildNodePtr().PrintScore()
}

// Function parameter is pointer, and its child-node method is directly referenced in chain without validation
func np10Example(d *DataInfo) {
	age := d.GetChildNodeNonPtr().GetGrandsonNodePtr().Age // want "potential nil pointer reference"
	fmt.Println(age)
	// d is a potential nil pointer
	// d.GetChildNodeNonPtr() is not a pointer, just a struct variable
	// d.GetChildNodeNonPtr().GetGrandsonNodePtr() is a potential nil pointer

	// It can be written as follows, and will be more safe.
	// if d == nil {
	// 	 return
	// }
	//
	// if d.GetChildNodeNonPtr().GetGrandsonNodePtr() != nil {
	//	 age := d.GetChildNodeNonPtr().GetGrandsonNodePtr().Age
	//	 fmt.Println(age)
	// }

	// Or:
	// if d != nil && d.GetChildNodeNonPtr().GetGrandsonNodePtr() != nil {
	//	 age := d.GetChildNodeNonPtr().GetGrandsonNodePtr().Age
	//	 fmt.Println(age)
	// }
}

// A pointer variable obtained by calling an external function, and its child-node method is directly referenced in chain without validation
func np11Example() {
	d := GetDataInfo()
	age := d.GetChildNodeNonPtr().GetGrandsonNodePtr().Age // want "potential nil pointer reference"
	fmt.Println(age)

	// d is a potential nil pointer
	// d.GetChildNodeNonPtr() is not a pointer, just a struct variable
	// d.GetChildNodeNonPtr().GetGrandsonNodePtr() is a potential nil pointer

	// It can be written as follows, and will be more safe.
	// if d == nil {
	// 	 return
	// }
	//
	// if d.GetChildNodeNonPtr().GetGrandsonNodePtr() != nil {
	//	 age := d.GetChildNodeNonPtr().GetGrandsonNodePtr().Age
	//	 fmt.Println(age)
	// }

	// Or:
	// if d != nil && d.GetChildNodeNonPtr().GetGrandsonNodePtr() != nil {
	//	 age := d.GetChildNodeNonPtr().GetGrandsonNodePtr().Age
	//	 fmt.Println(age)
	// }
}

// Skip the parent node pointer check and directly verify the child node.
func np12Example(d *DataInfo) {
	if d.A != nil { // want "potential nil pointer reference"
		fmt.Println(d.A.B) // want "potential nil pointer reference"
	}

	// d is a potential nil pointer. It should valid d first.
	// It can be written as follows, and will be more safe.

	// if d != nil && d.A != nil {
	//	 fmt.Println(d.A.B)
	// }

	// Or:
	// if d == nil {
	//	 return
	// }
	//
	// if d.A != nil {
	//	 fmt.Println(d.A.B)
	//}

	// Or:
	// if d == nil {
	//	 return
	// }
	// if d.A == nil {
	//	 return
	// }
	// fmt.Println(d.A.B)
}

func GetDataInfo() *DataInfo {
	return nil
}

func GetDataInfoList() []*DataInfo {
	return []*DataInfo{
		nil, nil, nil,
	}
}
