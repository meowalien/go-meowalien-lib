package math

import (
	"errors"
	"math"
)

// 絕對值
func AbsInt64(x int64) int64 {
	if x >= 0{
		return x
	}else{
		return x *-1
	}
}


// 四捨五入
func Round(x float64) int {
	return int(math.Floor(x + 0.5))
}

// 10進制數字位數
func IntLength(a int64) int {
	count := 0
	for a != 0 {
		a /= 10
		count++
	}
	return count
}

var ErrOverflow = errors.New("integer overflow")

func Add32(left, right int32) (int32, error) {
	if right > 0 {
		if left > math.MaxInt32-right {
			return 0, ErrOverflow
		}
	} else {
		if left < math.MinInt32-right {
			return 0, ErrOverflow
		}
	}
	return left + right, nil
}

// return true if overflow
func Add64(left, right int64) (int64, bool) {
	if !(right == 0 || left == 0){
		if right > 0 {
			if left > math.MaxInt64-right {
				return 0, true
			}
		} else {
			if left < math.MinInt64-right {
				return 0, true
			}
		}
	}

	return left + right, false
}
