package global_pool

import "sync"

const defaultByteCap = 1024*32

var byteArrayPool = sync.Pool{New: func() interface{}{
		return make([]byte, 0,defaultByteCap )
}}
func GetByteArray() []byte {
return byteArrayPool.Get().([]byte)
}

func PutByteArray(b []byte ) {
	byteArrayPool.Put(b)
}
