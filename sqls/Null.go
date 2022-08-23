package sqls

import (
	"database/sql"
	"database/sql/driver"
	"time"
)

type ScanAndValuer interface {
	sql.Scanner
	driver.Valuer
}

type NullSqlValue interface {
	*string | *time.Time | *bool | *float64 | *int64 | *int32 | *int16 | *byte
}

func Null[T NullSqlValue](t T) ScanAndValuer {
	switch tp := (any)(t).(type) {
	case *string:
		if t == nil {
			return &sql.NullString{}
		}
		return &sql.NullString{
			String: *tp,
			Valid:  true,
		}
	case *time.Time:
		if t == nil {
			return &sql.NullTime{}
		}
		return &sql.NullTime{
			Time:  *tp,
			Valid: true,
		}
	case *bool:
		if t == nil {
			return &sql.NullBool{}
		}
		return &sql.NullBool{
			Bool:  *tp,
			Valid: true,
		}
	case *float64:
		if t == nil {
			return &sql.NullFloat64{}
		}
		return &sql.NullFloat64{
			Float64: *tp,
			Valid:   true,
		}
	case *int64:
		if t == nil {
			return &sql.NullInt64{}
		}
		return &sql.NullInt64{
			Int64: *tp,
			Valid: true,
		}
	case *int32:
		if t == nil {
			return &sql.NullInt32{}
		}
		return &sql.NullInt32{
			Int32: *tp,
			Valid: true,
		}
	case *int16:
		if t == nil {
			return &sql.NullInt16{}
		}
		return &sql.NullInt16{
			Int16: *tp,
			Valid: true,
		}
	case *byte:
		if t == nil {
			return &sql.NullByte{}
		}
		return &sql.NullByte{
			Byte:  *tp,
			Valid: true,
		}
	default:
		panic("unsupported type")
	}
}
