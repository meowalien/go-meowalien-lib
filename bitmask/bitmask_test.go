package bitmask

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	A OffsetBitmask = iota + 1
	B
	C
)

var B_C = A.Add(B)

func TestOffsetBitmask(t *testing.T) {
	assert.True(t, A.Has(A))
	assert.True(t, B.Has(B))
	assert.True(t, C.Has(C))

	assert.False(t, A.Has(B))
	assert.False(t, B.Has(A))
	assert.False(t, C.Has(A))
	assert.False(t, C.Has(B))

	assert.True(t, A.Add(B).Has(A))
	assert.True(t, B.Add(A).Has(A))
	assert.True(t, B.Add(A).Has(B))

	assert.False(t, B.Add(A).Has(C))

}

func TestBitmask(t *testing.T) {
	a := A.Add(B).Add(C)
	b := A.Add(B)

	assert.True(t, a.Has(b))
	assert.True(t, b.Has(a))
	assert.False(t, b.Has(C))

}
