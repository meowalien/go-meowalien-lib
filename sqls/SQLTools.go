package sqls

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// JsonTypeStruct provide a simple solution
// to make struct use as Json in SQL query
type JsonTypeStruct struct {
	Thing interface{}
}

func (m JsonTypeStruct) Value() (driver.Value, error) {
	j, err := json.Marshal(m.Thing)
	if err != nil {
		return nil, err
	}
	return driver.Value(j), nil
}

func (m *JsonTypeStruct) Scan(src interface{}) error {
	var source []byte
	switch s := src.(type) {
	case []uint8:
		source = s
	case nil:
		return nil
	default:
		return errors.New("incompatible type")
	}
	err := json.Unmarshal(source, src)
	if err != nil {
		return err
	}
	return nil
}
