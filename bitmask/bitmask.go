package bitmask

type Bitmask uint64

func (a Bitmask) Has(b Bitmask) bool {
	return a&b != 0
}

func (a Bitmask) Add(b Bitmask) Bitmask {
	return a | b
}
