package errs

import (
	"errors"
	"fmt"
	"github.com/meowalien/go-meowalien-lib/bitmask"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWithLine_errors(t *testing.T) {
	err1 := errors.New("Error 1")
	err2 := errors.New("Error 2")
	err4 := errors.New("Error 4")
	err5 := errors.New("Error 5")
	err3 := New(err1, err2, err4, err5)
	fmt.Println(err3)
}
func TestWithLine_new(t *testing.T) {
	err2 := New("Error 2")
	err4 := New("Error 4")
	err5 := New("Error 5")

	err1 := New("Error 1")
	err3 := New(err1, err2, err4, err5)
	fmt.Println(err3)
}

func TestWithLine_one_with_line_one_plain_error(t *testing.T) {
	err1 := errors.New("Error 1")
	err2 := New("Error 2")
	err3 := New(err1, err2)
	fmt.Println(err3)
	fmt.Println(err1)
	fmt.Println(err2)

	//assert.EqualError(t, err3, "errs/Errs_test.go:21: { \n\tError 1\n\t=> errs/Errs_test.go:20: Error 2 \n}")
}

func TestWithLine_with_string(t *testing.T) {
	err3 := New("Error 1")
	fmt.Println(err3)
	//assert.EqualError(t, err3, "errs/Errs_test.go:28: Error 1")
}

func TestWithLine_with_strings(t *testing.T) {
	err3 := New("Error 1 { %s }", "some string")
	fmt.Println(err3)
	//assert.EqualError(t, err3, "errs/Errs_test.go:34: Error 1 some string")
}

func TestWithLineError_Wrap(t *testing.T) {
	err1 := New("Error 1")
	err2 := New("Error 2")
	err3 := New(err1, err2)
	fmt.Println("err3: ", err3)

	errx := New("Error 4")
	err4 := New(err3, errx)

	fmt.Println("err4: ", err4)
	err := errors.Unwrap(err4)
	fmt.Println("err1: ", err)
	err = errors.Unwrap(err)
	fmt.Println("err2: ", err)

	//fmt.Println(err)
	//if err != nil {
	//	assert.EqualError(t, err, "errs/Errs_test.go:46: Error 1")
	//}
}

func TestWithLine_with_line_code(t *testing.T) {
	err3 := New(errors.New("Error 1"))
	fmt.Println(err3)
	//assert.EqualError(t, err3, "errs/Errs_test.go:62: Error 1")
}

func TestWithLine_with_line_code_and_wrap(t *testing.T) {
	err3 := New(errors.New("Error 1"), "Error 2")
	fmt.Println(err3)
	//assert.EqualError(t, err3, "errs/Errs_test.go:68: { \n\tError 1\n\t=> errs/Errs_test.go:68: Error 2 \n}")
}

func TestWithLine_already_with_line(t *testing.T) {
	err1 := New("Error 1")
	err := New(err1)
	fmt.Println(err)
	//assert.EqualError(t, err, "errs/Errs_test.go:75: \n\terrs/Errs_test.go:74: Error 1")
}

func TestWithLine_nil_parent_case(t *testing.T) {
	var err1 error
	err2 := New("Error 2")
	err3 := New(err1, err2)
	fmt.Println(err3)
	//assert.EqualError(t, err3, "errs/Errs_test.go:83: Error 2")
}

func TestWithLine_defer(t *testing.T) {
	err := testFunc()
	fmt.Println(err)
	//assert.EqualError(t, err3, "errs/Errs_test.go:97: { \n\terrs/Errs_test.go:99: Error 1\n\t=> errs/Errs_test.go:96: Error 2 \n}")

}

func testFunc() (err error) {
	defer func() {
		err1 := errors.New("Error 2")
		err = New(err, err1)
	}()
	erra := New("Error 1")
	return New(erra)
}
func TestWithLine_defernil(t *testing.T) {
	err3 := testFuncnil()
	fmt.Println(err3)
	//assert.EqualError(t, err3, "errs/Errs_test.go:97: { \n\terrs/Errs_test.go:99: Error 1\n\t=> errs/Errs_test.go:96: Error 2 \n}")

}

func testFuncnil() (err error) {
	defer func() {
		var err1 error = nil
		if err1 != nil {
			err = New(err, err1)
		}
	}()
	erra := New("Error 1")
	return New(erra)
}

func TestErrorAndNil(t *testing.T) {
	var err2 error
	err1 := New("Error 1")
	err3 := New(err1, err2)
	fmt.Println(err3)
	//assert.EqualError(t, err3, "errs/Errs_test.go:105: Error 1")
}

// error from std pkg
func F1() error {
	return errors.New("Error 1")
}

func F2() (err error) {
	err = F1()
	if err != nil {
		err = New(err)
		return
	}
	return
}

func F3() (err error) {
	err = F2()
	if err != nil {
		err = New(err)
		return
	}
	return
}

func TestFF(t *testing.T) {
	err := F3()
	if err != nil {
		err = New(err)
		fmt.Println(err)
		return
	}
}

const (
	_ bitmask.OffsetBitmask = iota
	A
	B
	D
	E
)

func TestErrCode(t *testing.T) {
	err1 := New("Error 1").WithCode(A)
	err2 := New("Error 2").WithCode(B)
	fmt.Println(err1.HasCode(A))
	fmt.Println(err2.HasCode(A))
}

func TestErrCode1(t *testing.T) {
	err1 := New("Error 1").WithCode(A)
	err2 := New("Error 2").WithCode(B).WithCode(A)
	fmt.Println(err1.HasCode(A))
	fmt.Println(err2.HasCode(A))
}

func TestErrCode2(t *testing.T) {
	err0 := New("Error 0")
	err1 := New("Error 1").WithCode(A).WithCode(B)
	err2 := New("Error 2").WithCode(A)

	assert.True(t, errors.Is(err1, err2))
	assert.True(t, err1.HasCode(A))
	assert.True(t, err1.Is(err2))

	assert.False(t, errors.Is(err0, err2))
	assert.True(t, errors.Is(err0, err0))
	assert.False(t, err0.HasCode(A))
	assert.False(t, err0.HasCode(B))
}

func TestErrWrap(t *testing.T) {
	a := New("Error 1")
	b := New(a, "Error { %s } 2", "Error 2-1")
	c := New(b, "Error 3")
	fmt.Println(c)
}

func TestPassOnly(t *testing.T) {
	a := New("Error 1")
	b := New(a)
	c := New(b)
	fmt.Println(c)
}
