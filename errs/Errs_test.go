package errs

import (
	"fmt"
	"testing"
)

func TestWithLine(t *testing.T) {
	rawErr := fmt.Errorf("raw err")

	weap1 := WithLine(rawErr)
	//fmt.Println("ee,ok:= e.(WithLineError)")

	fmt.Printf("%T\n",weap1)

	var err error = weap1
	weap2 := WithLine(err)
	fmt.Printf("%T\n",weap2)
}