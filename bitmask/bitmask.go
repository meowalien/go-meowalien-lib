package bitmask

const memoryLength = uint64(64)

func newLongBitmask(index uint64) (a Bitmask) {
	if index == 0 {
		return make(longBitmask, 0)
	}
	idx := index / memoryLength
	aa := make(longBitmask, idx+1)
	aa[idx] = uint64(1) << (index%memoryLength - 1)
	a = aa
	return
}

type Bitmask interface {
	Has(b Bitmask) bool
	Add(b Bitmask) Bitmask
	blockAt(index int) uint64
}

type OffsetBitmask uint64

func (o OffsetBitmask) Has(b Bitmask) bool {
	idx := int(uint64(o) / memoryLength)

	return o.blockAt(idx)&b.blockAt(idx) != 0
}

func (o OffsetBitmask) Add(b Bitmask) Bitmask {
	return newLongBitmask(uint64(o)).Add(b)
}

func (o OffsetBitmask) blockAt(index int) uint64 {
	if o == 0 {
		return 0
	}

	if index > int(uint64(o)/memoryLength) {
		return 0
	} else {
		return uint64(1) << (uint64(o)%memoryLength - 1)
	}
}

type longBitmask []uint64

func (a longBitmask) blockLength() int {
	return len(a)
}

func (a longBitmask) blockAt(index int) uint64 {
	if len(a) <= index {
		return 0
	} else {
		return a[index]
	}
}

func (a longBitmask) Has(b Bitmask) bool {
	for i := range a {
		if a.blockAt(i)&b.blockAt(i) != 0 {
			return true
		}
	}
	return false
}

func (a longBitmask) Add(b Bitmask) Bitmask {
	for i := range a {
		a[i] |= b.blockAt(i)
	}
	return a
}
