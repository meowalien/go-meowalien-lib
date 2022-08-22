package trace

import (
	"context"
	"testing"
)

func TestNewBlock(t *testing.T) {
	end := NewBlock(context.TODO(), "wayne")
	defer end()
}
