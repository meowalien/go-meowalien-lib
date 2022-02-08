package math

import (
	"fmt"
	"github.com/shopspring/decimal"
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

func Add64b(left, right int64, allowOverflow bool) int64 {
	ans , over := Add64(left , right )
	if over && !allowOverflow{
		panic(fmt.Sprintf("overflow when %d + %d" , left , right))
	}
	return ans
}

func Sub64b(left, right int64, allowOverflow bool) int64 {
	return Add64b(left , right*-1 , allowOverflow )
}
// return true if overflow
func Sub(left, right int64) (int64, bool) {
	return Add64(left , -1 * right)
}

// return true if overflow
func Mul(left, right int64) (int64, bool) {
	x := left * right
	return x , x != 0 && x/right != left
}

func  Mulb(a , b int64 , allowOverflow bool) int64{
	ans , over := Mul(a , b )
	if over && !allowOverflow{
		panic(fmt.Sprintf("overflow when %d * %d" , a , b))
	}
	return ans
}


// return true if overflow
func Div(left, right int64 ,  round bool) int64 {
	ans := decimal.NewFromInt(left).Div(decimal.NewFromInt(right))
	if round {
		return ans.Round(-1).IntPart()
	}else {
		return ans.IntPart()
	}
}

