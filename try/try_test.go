package try

import (
	"encoding/json"
	"fmt"
	"github.com/meowalien/go-meowalien-lib/errs"
	"reflect"
	"testing"
)

type InnerStruct struct {
	TheID string `json:"the_id"`
}
type InnerStruct1 struct {
	TheID1 string `json:"the_id1"`
}
type RootStruct struct {
	InnerStruct  *InnerStruct  `json:"inner_struct"`
	InnerStruct1 *InnerStruct1 `json:"inner_struct1"`
	OtherData    any           `json:"other_data"`
}
type OtherData struct {
	DataA string `json:"dataA"`
	DataB string `json:"dataB"`
}

func TestPartUnmarshal(t *testing.T) {
	stringJson := "{\n  \"other_data\": {\n    \"dataA\": \"valueA\",\n    \"dataB\": \"valueB\"\n  },\n  \"inner_struct\": {\n    \"the_id\": \"123\"\n  }\n}"
	otd := &OtherData{}
	rs := &RootStruct{
		OtherData: otd,
	}
	err := json.Unmarshal([]byte(stringJson), rs)
	if err != nil {
		err = errs.New(err)
		panic(err)
	}
	fmt.Println(rs.InnerStruct.TheID)
	fmt.Println(rs.OtherData)

}

type testStruct struct {
}

func (t *testStruct) TestFunc(msg *OtherData) {
	fmt.Println("msg: ", msg)
}

func TestReflectFunc(t *testing.T) {
	ts := testStruct{}
	err := testSubscribe(ts.TestFunc)
	if err != nil {
		err = errs.New(err)
		panic(err)
	}
}

func testSubscribe(testFunc any) (err error) {
	valueOfFunc := reflect.ValueOf(testFunc)
	tp := reflect.TypeOf(testFunc)
	if tp.Kind() != reflect.Func {
		err = errs.New("testFunc is not a function")
		return
	}
	if tp.NumIn() != 1 {
		err = errs.New("testFunc must have only 1 input")
		return
	}
	typeOfInputParameter := tp.In(0)

	if typeOfInputParameter.Kind() != reflect.Ptr {
		err = errs.New("testFunc's input parameter must be a pointer")
		return
	}
	pointerToElementOfInputParameter := reflect.New(typeOfInputParameter.Elem())

	stringJson := "{\n  \"other_data\": {\n    \"dataA\": \"valueA\",\n    \"dataB\": \"valueB\"\n  },\n  \"inner_struct\": {\n    \"the_id\": \"123\"\n  }\n}"

	rs := &RootStruct{
		OtherData: pointerToElementOfInputParameter.Interface(),
	}

	err = json.Unmarshal([]byte(stringJson), rs)
	if err != nil {
		err = errs.New(err)
		return
	}

	fm := pointerToElementOfInputParameter.Elem().FieldByName("GaimID").String()
	fmt.Println(fm)
	//tx = apm.StartTransactionOptions(valueOfFunc.String(), g.GameID(), traceId)

	valueOfFunc.Call([]reflect.Value{pointerToElementOfInputParameter})
	return
}
