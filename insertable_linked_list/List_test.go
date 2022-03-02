package insertable_linked_list

import (
	"fmt"
	"testing"
)

func TestList(t *testing.T) {
	theList := New()
	theList.Put(0, "T1-0")
	//fmt.Println("After T1-0")
	theList.Put(3, "T1-1")
	//fmt.Println("After T1-1")
	theList.Put(2, "T1-2")
	//fmt.Println("After T1-2")
	for  i2 := range theList.Iterator() {
		fmt.Println("ans: ",i2)
	}
}