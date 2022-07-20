package insertable_linked_list

import (
	"fmt"
	"testing"
)

func TestList(t *testing.T) {
	theList := New()
	theList.Put(1, "T1-0")
	theList.Put(2, "T1-1")
	theList.Put(3, "T1-2")
	theList.Put(4, "T1-4")
	for checkEle := theList.Front(); checkEle != nil; checkEle = checkEle.Next() {
		fmt.Println(checkEle)
	}

	lst := theList.GetAndRemoveLowerThen(2)
	fmt.Println("after GetAndRemoveLowerThen")
	for checkEle := theList.Front(); checkEle != nil; checkEle = checkEle.Next() {
		fmt.Println(checkEle)
	}

	for _, element := range lst {
		fmt.Println("element: ", element)
	}

}
