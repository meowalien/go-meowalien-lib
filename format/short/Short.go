package short

import (
	"go.uber.org/zap/buffer"
	"log"
	"math"
)

var base58ShortMap []byte
var deBase58ShortMap = map[uint8]int64{}

func init() {
	var i byte = 48
	var key int64

	doappend := func(r byte) {
		base58ShortMap = append(base58ShortMap, r)
		deBase58ShortMap[r] = key
		key++
	}

	//0-9
	for ; i <= 57; i++ {
		if i == '0' {
			continue
		}

		doappend(i)
	}
	i = 65
	//A-Z
	for ; i <= 90; i++ {
		if i == 'I' || i == 'O' {
			continue
		}
		doappend(i)
	}

	i = 97
	// a-z
	for ; i <= 122; i++ {
		if i == 'l' || i == 'O' {
			continue
		}
		doappend(i)
	}
}
var maxShortMap []byte
var demaxShortMap = map[uint8]int64{}


func init() {
	var i byte = 48
	var key int64

	add := func(r byte) {
		maxShortMap = append(maxShortMap, r)
		demaxShortMap[r] = key
		key++
	}

	//0-9
	for ; i <= 57; i++ {
		add(i)
	}
	i = 65
	//A-Z
	for ; i <= 90; i++ {
		add(i)
	}

	i = 97
	// a-z
	for ; i <= 122; i++ {
		add(i)
	}

	add('_')
	add('-')
	add('=')
	add('+')
	add('#')
	add('@')
}

func DeBase58Short(n string) int64 {
	return DeShortWithMap(n, deBase58ShortMap)
}

func DeShortWithMap(n string, mp map[uint8]int64) int64 {
	max := float64(len(mp))
	if n == "" {
		panic("n == \"\"")
	}
	l := len(n)

	d := int64(0)
	for i := l - 1; i >= 0; i-- {
		d += mp[n[i]] *int64( math.Pow(max, float64(i)))
	}
	return d
}

func Base58Short(n int64) string {
	return ShortWithMap(n, base58ShortMap)
}

func MaxShort(n int64) string {
	return ShortWithMap(n, maxShortMap)
}

func DeMaxShort(n string) int64 {
	return DeShortWithMap(n, demaxShortMap)
}
var bytePool = buffer.NewPool()

func ShortWithMap(n int64, mp []byte) string {
	max := int64(len(mp))
	if n < 0 {
		panic("n < 0")
	}

	bt := bytePool.Get()
	defer bt.Free()
	defer bt.Reset()

	for q := n; q > 0; q = q / max {
		px := q%max
		b := mp[px]

		err := bt.WriteByte(b)
		if err != nil {
			log.Println("error when Write buffer at ShortWithMap")
			return ""
		}
	}
	return bt.String()
}
