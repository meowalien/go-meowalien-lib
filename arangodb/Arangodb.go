package arangodb

import (
	"context"
	"fmt"
	"github.com/arangodb/go-driver"
	"github.com/meowalien/go-meowalien-lib/errs"
	"io"
)

type Query interface {
	Query(ctx context.Context, query string, bindVars map[string]interface{}) (driver.Cursor, error)
}

type ReadDocumentFunc interface {
	ReadDocument(ctx context.Context, result interface{}) (driver.DocumentMeta, error)
	HasMore() bool
	io.Closer
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

func QueryAndRead[T any](ctx context.Context, q Query, aqlQuery string, keys map[string]interface{}) (result []T, err error) {
	cursor, err := q.Query(ctx, aqlQuery, keys)
	if err != nil {
		return nil, fmt.Errorf("Repo GetGameIDsInThemes failed: %w", err)
	}
	defer func(cursor io.Closer) {
		err1 := cursor.Close()
		if err1 != nil {
			err = errs.New(err, err1)
		}
	}(cursor)
	return ReadDocument[T](ctx, cursor)

}
