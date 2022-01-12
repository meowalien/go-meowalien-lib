package value

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/meowalien/go-meowalien-lib/errs"
	"github.com/meowalien/go-meowalien-lib/math"
	"strconv"
)

// 會以此比例為預設基準輸出與讀取
//var BaseRate int64 = 1

type Value interface {
	sql.Scanner
	driver.Valuer
	json.Unmarshaler
	MultiplyByActual(dc int64) Value
	MultiplyBy(dc Value) Value
	DivideBy(dc Value) Value
	DivideByActual(dc int64) Value
	Add(dc Value) Value
	AddActual(dc int64) Value
	ActualEqual(i int64) bool
	Actual() int64
	To(ratio int64) int64
	//BaseValue() int64
	// 取得基準倍率
	BaseRate() int64
	Sub(dc Value) Value
	SubActual(dc int64) Value
	//ExpansionRate() int64
	Negative() Value
	ChangeBase(base int64)Value

	Replace(dc Value)
	RMultiplyBy(dc Value)
	RDivideBy(dc Value)
	RAdd(dc Value)
	RSub(dc Value)
}

func New(actual int64, base int64) Value {
	return &value{
		actual:        actual,
		//expansionRate: rate,
		base:          base,
	}
}

type value struct {
	actual        int64
	//expansionRate int64
	base          int64
}

func (d *value) Replace(dc Value) {
	d.actual = dc.Actual()
	//d.expansionRate = dc.ExpansionRate()
	d.base = dc.BaseRate()
}

func (d *value) BaseRate() int64 {
	return d.base
}

func (d *value) RMultiplyBy(dc Value) {
	actual := dc.Actual()
	//if d.expansionRate != dc.ExpansionRate() {
	//	actual = changeRate(dc.Actual(), dc.ExpansionRate(), d.expansionRate)
	//}
	d.actual *= actual
}

func (d *value) RDivideBy(dc Value) {
	actual := dc.Actual()
	//if d.expansionRate != dc.ExpansionRate() {
	//	actual = changeRate(dc.Actual(), dc.ExpansionRate(), d.expansionRate)
	//}
	d.actual /= actual
}

func (d *value) RAdd(dc Value) {
	actual := dc.Actual()
	//if d.expansionRate != dc.ExpansionRate() {
	//	actual = changeRate(dc.Actual(), dc.ExpansionRate(), d.expansionRate)
	//}
	d.actual += actual
	return
}

func (d *value) RSub(dc Value) {
	actual := dc.Actual()
	//if d.expansionRate != dc.ExpansionRate() {
	//	actual = changeRate(dc.Actual(), dc.ExpansionRate(), d.expansionRate)
	//}
	d.actual -= actual
}

func (d *value) ChangeBase(base int64) Value {
	d.base = base
	return d
}

func (d value) MultiplyByActual(dc int64) Value {
	return &value{
		actual:        d.actual * dc,
		//expansionRate: d.expansionRate,
		base:          d.base,
	}
}

func (d value) MultiplyBy(dc Value) Value {
	actual := dc.Actual()
	//if d.expansionRate != dc.ExpansionRate() {
	//	actual = changeRate(dc.Actual(), dc.ExpansionRate(), d.expansionRate)
	//}
	return d.MultiplyByActual(actual)
}

func (d value) DivideBy(dc Value) Value {
	actual := dc.Actual()
	//if d.expansionRate != dc.ExpansionRate() {
	//	actual = changeRate(dc.Actual(), dc.ExpansionRate(), d.expansionRate)
	//}
	return d.DivideByActual(actual)
}

func (d value) DivideByActual(dc int64) Value {
	if d.actual%dc != 0 {
		panic(fmt.Sprintf("not divisible, actual: %d , base: %d", d.actual, d.base))
	}
	return &value{
		actual:        d.actual / dc,
		//expansionRate: d.expansionRate,
		base:          d.base,
	}
}

func (d value) Add(dc Value) Value {
	actual := dc.Actual()
	//if d.expansionRate != dc.ExpansionRate() {
	//	actual = changeRate(dc.Actual(), dc.ExpansionRate(), d.expansionRate)
	//}
	return d.AddActual(actual)
}

func (d value) AddActual(dc int64) Value {
	a, over := math.Add64(d.actual, dc)
	if over {
		panic(fmt.Sprintf("overflow if ( %d + %d )", d.actual, dc))
	}
	return &value{
		actual:        a,
		//expansionRate: d.expansionRate,
		base:          d.base,
	}
}

func (d value) ActualEqual(i int64) bool {
	return d.actual == i
}

func (d value) To(ratio int64) int64 {
	return d.actual * ratio / d.base
}

func (d value) String() string {
	return fmt.Sprintf("%d", d.Actual())
}

//func (d value) BaseValue() int64 {
//	x := d.actual * d.base
//	//if x%d.expansionRate != 0 {
//	//	panic(fmt.Sprintf("not divisible, x: %d , expansionRate: %d", x, d.expansionRate))
//	//}
//	return x / d.expansionRate
//}

func (d value) Actual() int64 {
	return d.actual
}
//
//func (d value) ExpansionRate() int64 {
//	return d.expansionRate
//}

func (d value) Sub(dc Value) Value {
	actual := dc.Actual()
	//if d.expansionRate != dc.ExpansionRate() {
	//	actual = changeRate(dc.Actual(), dc.ExpansionRate(), d.expansionRate)
	//}
	return d.SubActual(actual)
}

func (d value) SubActual(dc int64) Value {
	return d.AddActual(-1 * dc)
}

func (d value) Negative() Value {
	return d.MultiplyByActual(-1)
}

func (d *value) Scan(v interface{}) error {
	//fmt.Printf("Scan:%v \n",v)

	if v == nil {
		return errs.WithLine("the DecimalValue should not be nil")
	}
	var vl int64
	switch vv := v.(type) {
	case int64:
		vl = vv
	case []byte:
		ac, err := strconv.ParseInt(string(vv), 10, 64)
		if err != nil {
			return errs.WithLine(err)
		}
		vl = ac
	default:
		return errs.WithLine("not supported type : %T", v)
	}
	//fmt.Println("changeRate(d.actual, d.base, d.expansionRate):", changeRate(d.actual, d.base, d.expansionRate))
	//fmt.Println(" d.expansionRate != d.base: ", d.expansionRate != d.base)
	//if d.expansionRate != d.base {
	//	d.actual = changeRate(vl, d.base, d.expansionRate)
	//} else {
	//	d.actual = vl
	//}

	d.actual = vl
	//fmt.Printf("d.expansionRate:%v \n",d.expansionRate)
	//fmt.Printf("d.base:%v \n",d.base)
	return nil
}

func (d value) Value() (driver.Value, error) {
	return d.Actual(), nil
}

func (d value) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%d\"", d.Actual())), nil
}

func (d *value) UnmarshalJSON(bytes []byte) error {
	status := string(bytes)
	if n := len(status); n > 1 && status[0] == '"' && status[n-1] == '"' {
		status = status[1 : n-1]
	}
	parseInt, err := strconv.ParseInt(status, 10, 64)
	if err != nil {
		return err
	}

	d.actual = parseInt
	//if d.expansionRate != d.base {
	//	d.actual = changeRate(parseInt, d.base, d.expansionRate)
	//} else {
	//	d.actual = parseInt
	//}

	return nil
}

//func changeRate(i int64, from int64, to int64) int64 {
//	x := i * to
//	if x%from != 0 {
//		panic(fmt.Sprintf("not divisible, i: %d , from: %d , to: %d" , i, from, to))
//	}
//	return x / from
//}
