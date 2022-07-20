package reflect

import (
	"fmt"
	"reflect"
)

// 將 args 的指定位置解構到 targetPointer 內
func ParseArgs(args []interface{}, index int, targetPointer interface{}) error {
	if l := len(args); l < index+1 {
		return fmt.Errorf("not enough length need: %d , got:%d , args: %v", index, l, args)
	}

	valueOfTargetPointer := reflect.ValueOf(targetPointer)
	if valueOfTargetPointer.Kind() != reflect.Ptr {
		return fmt.Errorf("the targetPointer must be pointer")
	}

	targetType := valueOfTargetPointer.Elem().Type()

	typeOfGiven := reflect.ValueOf(args[index]).Type()

	if typeOfGiven != targetType {
		return fmt.Errorf("type not match expected: %v , got: %v", targetType, typeOfGiven)
	}

	reflect.ValueOf(targetPointer).Elem().Set(reflect.ValueOf(args[index]))
	return nil
}
