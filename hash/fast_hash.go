package hash

import (
	"bytes"
	"encoding/gob"
	"hash/fnv"
)

func FastHash[T any](data T) (uint64, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	err := enc.Encode(data)
	if err != nil {
		return 0, err
	}

	h := fnv.New64a()
	_, err = h.Write(buf.Bytes())
	if err != nil {
		return 0, err
	}

	return h.Sum64(), nil
}
