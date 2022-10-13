package parse

import (
	"encoding/json"
	"github.com/meowalien/go-meowalien-lib/errs"
	"strings"
)

func ParseNestedMap(m map[string]interface{}, keys ...string) (value any, exist bool, err error) {
	if m == nil {
		err = errs.New("map is nil")
		return
	}
	if len(keys) == 0 {
		err = errs.New("keys is empty")
		return
	}
	v, exist := m[keys[0]]
	if !exist {
		return
	}
	if len(keys) == 1 {
		switch vv := v.(type) {
		case json.Number:
			// if the map was created with decoder.UseNumber(),
			// then all numerical value will be parsed as json.Number
			if strings.Contains(vv.String(), ".") {
				value, err = vv.Float64()
			} else {
				value, err = vv.Int64()
			}
		default:
			value = vv
		}
		return
	}

	var ok bool
	m, ok = v.(map[string]interface{})
	if !ok {
		err = errs.New("value %v is not map", v)
		return
	}
	return ParseNestedMap(m, keys[1:]...)
}
