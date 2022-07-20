package deap_copy

import "reflect"

func Copy(v interface{}) interface{} {
	x := reflect.ValueOf(v).Elem()
	vp2 := reflect.New(x.Type())
	vp2.Elem().Set(x)
	return vp2.Interface()
}
