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

type Value interface {
	sql.Scanner
	driver.Valuer
	json.Unmarshaler
	MultiplyBy(dc Value) Value
	DivideBy(dc Value) Value
	Add(dc Value) Value
	RAdd(dc Value) Value
	Sub(dc Value) Value
	Negative() Value
	toModify(modify modify) Value
	ActualEqual(i int64) bool
	Actual() int64
	As(ratio int64) int64
	Reserved() int64
}

type Option struct {
	Round         bool
	AllowOverflow bool
}

// base = 相對於其他比例的倍數
// reserved = 實際存儲的倍數
// number = 當前倍率下數字
func New(number int64, base int64, reserved int64, opt ...Option) Value {
	var op Option
	if opt == nil || len(opt) == 0 {
		op = DefaultOption()
	}
	return &value{
		actual: math.Div(math.Mulb(number, reserved, op.AllowOverflow), base, op.Round),
		modify: modify{
			base:     base,
			reserved: reserved,
		},
		Option: op,
	}
}

func DefaultOption() Option {
	return Option{
		Round:         false,
		AllowOverflow: false,
	}
}

type modify struct {
	base     int64
	reserved int64
}

func (m modify) Reserved() int64 {
	return m.reserved
}

type value struct {
	Option
	modify
	actual int64
}

func (d value) toModify(modify modify) Value {
	if d.modify.reserved == modify.reserved && d.modify.base == modify.base {
		return &d
	}
	d.actual = math.Div(math.Mulb(d.actual, modify.reserved, d.AllowOverflow), d.reserved, d.Round)
	d.modify = modify
	return &d
}

func (d value) MultiplyBy(dc Value) Value {
	d.actual = math.Div(math.Mulb(d.actual, dc.toModify(d.modify).Actual(), d.AllowOverflow), d.reserved, d.Round)
	return &d
}

func (d value) DivideBy(dc Value) Value {
	d.actual = math.Div(math.Mulb(d.actual, d.reserved, d.AllowOverflow), dc.toModify(d.modify).Actual(), d.Round)
	return &d
}

func (d value) Add(dc Value) Value {
	d.actual = math.Add64b(d.actual, dc.toModify(d.modify).Actual(), d.AllowOverflow)
	return &d
}
func (d *value) RAdd(dc Value) Value {
	d.actual = math.Add64b(d.actual, dc.toModify(d.modify).Actual(), d.AllowOverflow)
	return d
}
func (d *value) As(base int64) int64 {
	return math.Div(math.Mulb(d.actual, base, d.AllowOverflow), d.reserved, d.Round)
}

func (d *value) ActualEqual(i int64) bool {
	return d.actual == i
}

func (d *value) String() string {
	return fmt.Sprintf("%d", d.Actual())
}

func (d *value) Actual() int64 {
	return d.actual
}

func (d *value) Sub(dc Value) Value {
	return d.Add(dc.Negative())
}

func (d value) Negative() Value {
	d.actual *= -1
	return &d
}

func (d *value) Scan(v interface{}) error {
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
	d.actual = vl
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
	return nil
}
