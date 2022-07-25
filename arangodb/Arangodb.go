package arangodb

import (
	"context"
	"github.com/arangodb/go-driver"
)

type ReadDocumentFunc interface {
	ReadDocument(ctx context.Context, result interface{}) (driver.DocumentMeta, error)
	HasMore() bool
}

func ReadDocument[T any](ctx context.Context, f ReadDocumentFunc) (result []T, err error) {
	for f.HasMore() {
		var raw T
		_, err = f.ReadDocument(ctx, &raw)
		if err != nil {
			return
		}
		result = append(result, raw)
	}
	return
}
