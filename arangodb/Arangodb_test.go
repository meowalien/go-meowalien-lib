package arangodb

import (
	"context"
	"fmt"
	"github.com/arangodb/go-driver"
	"reflect"
	"testing"
)

type ReadDocumentFuncMock struct {
	ReadTimes int
}

func (r *ReadDocumentFuncMock) ReadDocument(ctx context.Context, result interface{}) (driver.DocumentMeta, error) {
	reflect.ValueOf(result).Elem().SetString("test")
	return driver.DocumentMeta{}, nil
}

func (r *ReadDocumentFuncMock) HasMore() bool {
	r.ReadTimes--
	return r.ReadTimes >= 0
}

func TestReadDocument(t *testing.T) {
	cursor := &ReadDocumentFuncMock{ReadTimes: 1}
	ss, err := ReadDocument[string](context.TODO(), cursor)
	fmt.Println(err)
	fmt.Println(ss)
}
func TestReadDocuments(t *testing.T) {
	cursor := &ReadDocumentFuncMock{ReadTimes: 3}
	ss, err := ReadDocument[string](context.TODO(), cursor)
	fmt.Println(err)
	fmt.Println(ss)
}
