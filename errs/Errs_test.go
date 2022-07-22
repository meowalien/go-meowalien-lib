package errs

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWithLine_two_errors(t *testing.T) {
	err1 := errors.New("Error 1")
	err2 := errors.New("Error 2")
	err3 := WithLine(err1, err2)
	fmt.Println(err3)
	assert.EqualError(t, err3, "errs/Errs_test.go:13: { \n\tError 1\n\t=> Error 2 \n}")
}

func TestWithLine_one_with_line_one_plain_error(t *testing.T) {
	err1 := errors.New("Error 1")
	err2 := WithLine("Error 2")
	err3 := WithLine(err1, err2)
	fmt.Println(err3)

	assert.EqualError(t, err3, "errs/Errs_test.go:21: { \n\tError 1\n\t=> errs/Errs_test.go:20: Error 2 \n}")
}

func TestWithLine_with_string(t *testing.T) {
	err3 := WithLine("Error 1")
	fmt.Println(err3)
	assert.EqualError(t, err3, "errs/Errs_test.go:28: Error 1")
}

func TestWithLine_with_strings(t *testing.T) {
	err3 := WithLine("Error 1", "some string")
	fmt.Println(err3)
	assert.EqualError(t, err3, "errs/Errs_test.go:34: Error 1 some string")
}

func TestWithLine_with_stringf(t *testing.T) {
	err3 := WithLine("Error {%s} 1", "some string")
	fmt.Println(err3)
	assert.EqualError(t, err3, "errs/Errs_test.go:40: Error {some string} 1")
}

func TestWithLineError_Wrap(t *testing.T) {
	err1 := WithLine("Error 1")
	err2 := WithLine("Error 2")
	err3 := WithLine(err1, err2)
	fmt.Println(err3)
	err := errors.Unwrap(err3)
	if err != nil {
		assert.EqualError(t, err, "{ \n\terrs/Errs_test.go:46: Error 1\n\t=> errs/Errs_test.go:47: Error 2 \n}")
	}
	err = errors.Unwrap(err)
	fmt.Println(err)
	if err != nil {
		assert.EqualError(t, err, "errs/Errs_test.go:46: Error 1")
	}
}

func TestWithLine_with_line_code(t *testing.T) {
	err3 := WithLine(errors.New("Error 1"))
	fmt.Println(err3)
	assert.EqualError(t, err3, "errs/Errs_test.go:62: Error 1")
}

func TestWithLine_with_line_code_and_wrap(t *testing.T) {
	err3 := WithLine(errors.New("Error 1"), "Error 2")
	fmt.Println(err3)
	assert.EqualError(t, err3, "errs/Errs_test.go:68: { \n\tError 1\n\t=> errs/Errs_test.go:68: Error 2 \n}")
}

func TestWithLine_already_with_line(t *testing.T) {
	err1 := WithLine("Error 1")
	err := WithLine(err1)
	fmt.Println(err)
	assert.EqualError(t, err, "errs/Errs_test.go:75: \n\terrs/Errs_test.go:74: Error 1")
}

func TestWithLine_nil_parent_case(t *testing.T) {
	var err1 error
	err2 := errors.New("Error 2")
	err3 := WithLine(err1, err2)
	fmt.Println(err3)
	assert.EqualError(t, err3, "errs/Errs_test.go:83: Error 2")
}

func TestWithLine_defer(t *testing.T) {
	err3 := testFunc()
	assert.EqualError(t, err3, "errs/Errs_test.go:97: { \n\terrs/Errs_test.go:99: Error 1\n\t=> errs/Errs_test.go:96: Error 2 \n}")

}

func testFunc() (err error) {
	defer func() {
		err1 := WithLine("Error 2")
		err = WithLine(err, err1)
	}()
	return WithLine("Error 1")

}
