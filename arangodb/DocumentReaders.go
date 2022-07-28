package arangodb

import (
	"context"
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

type decoder[R any] interface {
	func(ctx context.Context, f ReadDocumentFunc) (result R, err error)
}

func withCursor[R any, D decoder[R]](ctx context.Context, q Query, aqlQuery string, keys map[string]interface{}, documentReader D) (res R, err error) {
	cursor, err := q.Query(ctx, aqlQuery, keys)
	if err != nil {
		err = errs.New(err)
		return
	}
	defer func(cursor io.Closer) {
		err1 := cursor.Close()
		if err1 != nil {
			err = errs.New(err, err1)
		}
	}(cursor)
	return documentReader(ctx, cursor)
}

func ReadDocument[T any](ctx context.Context, f ReadDocumentFunc) (result T, err error) {
	r, err := ReadDocumentPtr[T](ctx, f)
	if r != nil {
		result = *r
	}
	return
}

func ReadDocumentPtr[T any](ctx context.Context, f ReadDocumentFunc) (result *T, err error) {
	var r T
	_, err = f.ReadDocument(ctx, &r)
	if err != nil {
		if driver.IsNotFound(err) || driver.IsNoMoreDocuments(err) {
			err = nil
			return
		}
		err = errs.New(err)
		return
	}
	result = &r
	return
}

func ReadDocuments[T any, R []T](ctx context.Context, f ReadDocumentFunc) (result []T, err error) {
	for f.HasMore() {
		var raw *T
		raw, err = ReadDocumentPtr[T](ctx, f)
		if err != nil {
			err = errs.New(err)
			return
		}
		result = append(result, *raw)
	}
	return
}

func ReadDocumentsPtr[T any](ctx context.Context, f ReadDocumentFunc) (result []*T, err error) {
	for f.HasMore() {
		var raw T
		_, err = f.ReadDocument(ctx, &raw)
		if err != nil {
			return
		}
		result = append(result, &raw)
	}
	return
}

func QueryAndReadDocumentPtr[T any](ctx context.Context, q Query, aqlQuery string, keys map[string]interface{}) (result *T, err error) {
	return withCursor(ctx, q, aqlQuery, keys, ReadDocumentPtr[T])
}

func QueryAndReadDocument[T any](ctx context.Context, q Query, aqlQuery string, keys map[string]interface{}) (result T, err error) {
	return withCursor(ctx, q, aqlQuery, keys, ReadDocument[T])
}

func QueryAndReadDocuments[T any](ctx context.Context, q Query, aqlQuery string, keys map[string]interface{}) (result []T, err error) {
	return withCursor(ctx, q, aqlQuery, keys, ReadDocuments[T])
}

func QueryAndReadDocumentsPtr[T any](ctx context.Context, q Query, aqlQuery string, keys map[string]interface{}) (result []*T, err error) {
	return withCursor(ctx, q, aqlQuery, keys, ReadDocumentsPtr[T])
}
