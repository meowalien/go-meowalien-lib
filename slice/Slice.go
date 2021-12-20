package slice

import (
	"fmt"
	"reflect"
)

func Remove(s interface{}, idx int) interface{} {
	if reflect.TypeOf(s).Kind() != reflect.Slice {
		panic(fmt.Sprintf("the input is not a slice: %v", s))
	}

	//fmt.Println(v.Len() > idx)
	//fmt.Println(v.Len() )
	//fmt.Println( idx)
	if v := reflect.ValueOf(s); v.Len() > idx {
		return reflect.AppendSlice(v.Slice(0, idx), v.Slice(idx+1, v.Len())).Interface()
	}

	return s
}

func ToInterfaceSlice(sl interface{}) (ans []interface{}) {
	switch aa := sl.(type) {
	case []int8:
		for i := range aa {
			ans = append(ans, aa[i])
		}
		return ans
	case []int16:
		for i := range aa {
			ans = append(ans, aa[i])
		}
		return ans
	case []int32:
		for i := range aa {
			ans = append(ans, aa[i])
		}
		return ans
	case []int64:
		for i := range aa {
			ans = append(ans, aa[i])
		}
		return ans
	case []uint8:
		for i := range aa {
			ans = append(ans, aa[i])
		}
		return ans
	case []uint16:
		for i := range aa {
			ans = append(ans, aa[i])
		}
		return ans
	case []uint32:
		for i := range aa {
			ans = append(ans, aa[i])
		}
		return ans
	case []uint64:
		for i := range aa {
			ans = append(ans, aa[i])
		}
		return ans
	case []int:
		for i := range aa {
			ans = append(ans, aa[i])
		}
		return ans
	case []uint:
		for i := range aa {
			ans = append(ans, aa[i])
		}
		return ans
	case []uintptr:
		for i := range aa {
			ans = append(ans, aa[i])
		}
		return ans
	case []string:
		for i := range aa {
			ans = append(ans, aa[i])
		}
		return ans
	default:
		return reflectInterfaceSlice(sl)
	}
}


func reflectInterfaceSlice(slice interface{}) []interface{} {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		panic("ToInterfaceSlice() given a non-slice type")
	}

	if s.IsNil() {
		return nil
	}

	ret := make([]interface{}, s.Len())

	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}
	return ret
}

//
//func ToInterfaceSlice(slice ...int{}) []interface{} {
//
//	s := reflect.ValueOf(slice)
//	if s.Kind() != reflect.Slice {
//		panic("ToInterfaceSlice() given a non-slice type")
//	}
//
//	if s.IsNil() {
//		return nil
//	}
//
//	ret := make([]interface{}, s.Len())
//
//	for i:=0; i<s.Len(); i++ {
//		ret[i] = s.Index(i).Interface()
//	}
//
//	return ret
//}

//func abc(x ...interface{}) {
//	for _, _ = range x {
//		//	do something
//	}
//}
//
//func XXX() {
//	a := []string{"one", "two"}
//	//b := make([]interface{} , len(a))
//	for _, s := range a {
//		//b[i] = s
//		abc(s)
//	}
//	abc(b...)
//}
