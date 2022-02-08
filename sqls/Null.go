package sqls

import (
	"database/sql"
	"github.com/meowalien/go-meowalien-lib/times"
	"time"
)

func NewNullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}


func NewNullTime(t *time.Time) sql.NullTime {
	if t == nil || t.Equal(times.NilTime) {
		return sql.NullTime{}
	}

	return sql.NullTime{
		Time:  *t,
		Valid: true,
	}
}

func NewNullBool(t bool) sql.NullBool {
	return sql.NullBool{
		Bool:  t,
		Valid: true,
	}
}

func NewNullFloat64(t float64) sql.NullFloat64 {
	if t == 0 {
		return sql.NullFloat64{}
	}
	return sql.NullFloat64{
		Float64: t,
		Valid:   true,
	}
}
func NewNullInt64(t int64) sql.NullInt64 {
	if t == 0 {
		return sql.NullInt64{}
	}
	return sql.NullInt64{
		Int64: t,
		Valid: true,
	}
}

func NewNullInt32(t int32) sql.NullInt32 {
	if t == 0 {
		return sql.NullInt32{}
	}
	return sql.NullInt32{
		Int32: t,
		Valid: true,
	}
}
