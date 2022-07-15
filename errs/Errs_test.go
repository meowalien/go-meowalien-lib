package errs

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWithLine_two_errors(t *testing.T) {
	err1 := errors.New("Error 1")
	err2 := errors.New("Error 2")
	err3 := WithLine(err1, err2)
	assert.EqualError(t, err3, "{ errs/Errors.go:16: Error 1 } -> { errs/Errors.go:16: Error 2 }")
}
func TestWithLine_one_with_line_one_plain_error(t *testing.T) {
	err1 := errors.New("Error 1")
	err2 := WithLine("Error 2")
	err3 := WithLine(err1, err2)
	assert.EqualError(t, err3, "{ errs/Errors.go:16: Error 1 } -> { errs/Errs_test.go:17: Error 2 }")
}
func TestWithLine_with_string(t *testing.T) {
	err3 := WithLine("Error 1")
	assert.EqualError(t, err3, "errs/Errs_test.go:22: Error 1")
}
func TestWithLine_with_strings(t *testing.T) {
	err3 := WithLine("Error 1", "some string")
	assert.EqualError(t, err3, "errs/Errs_test.go:26: Error 1 some string")
}

func TestWithLine_with_stringf(t *testing.T) {
	err3 := WithLine("Error {%s} 1", "some string")
	assert.EqualError(t, err3, "errs/Errs_test.go:31: Error {some string} 1")
}
