package hash

import (
	"encoding/base64"
	"encoding/gob"
	"hash/fnv"
)

func FastHash[T any](data T) (string, error) {
	h := fnv.New64a()
	enc := gob.NewEncoder(h)
	err := enc.Encode(data)
	if err != nil {
		return "", err
	}
	hashBytes := h.Sum(nil)
	return base64.URLEncoding.EncodeToString(hashBytes), nil
}
