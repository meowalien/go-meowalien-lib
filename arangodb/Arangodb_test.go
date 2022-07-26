package arangodb

import (
	"context"
	"github.com/arangodb/go-driver"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

type ReadDocumentFuncMock struct {
	*MockCursor
	ReadTimes int
}

func (r *ReadDocumentFuncMock) Close() error {
	return nil
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
	assert.NotPanics(t,
		func() {
			cursor := &ReadDocumentFuncMock{ReadTimes: 1}
			ss, err := ReadDocument[string](context.TODO(), cursor)
			assert.NoError(t, err)
			assert.Equal(t, ss, []string{"test"})
		})
}
func TestReadDocuments(t *testing.T) {
	assert.NotPanics(t,
		func() {
			cursor := &ReadDocumentFuncMock{ReadTimes: 3}
			ss, err := ReadDocument[string](context.TODO(), cursor)
			assert.NoError(t, err)
			assert.Equal(t, ss, []string{"test", "test", "test"})
		})
}

func TestQueryAndRead(t *testing.T) {
	assert.NotPanics(t,
		func() {
			testController := gomock.NewController(t)
			defer testController.Finish()
			mockCursor := NewMockCursor(testController)

			mockQuery := NewMockQuery(testController)
			cursor := &ReadDocumentFuncMock{ReadTimes: 1, MockCursor: mockCursor}
			mockQuery.EXPECT().Query(context.TODO(), "", map[string]interface{}{}).Return(cursor, nil)

			result, err := QueryAndRead[string](context.TODO(), mockQuery, "", map[string]interface{}{})
			assert.NoError(t, err)
			assert.Equal(t, result, []string{"test"})
		})
}
