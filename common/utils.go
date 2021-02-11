package common

import (
	"fmt"
	"reflect"
	"runtime"
)

func GetFileInfo(skip int) string {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return " [get file info failed] "
	}
	fn := runtime.FuncForPC(pc)
	return fmt.Sprintf(" [%s %s:%d] ", fn.Name(), file, line)
}


func InArray(v, list interface{}) bool {
	listValue := reflect.ValueOf(list)
	if listValue.Kind() == reflect.Slice || listValue.Kind() == reflect.Array {
		for i := 0; i < listValue.Len(); i++ {
			if reflect.DeepEqual(v, listValue.Index(i).Interface()) {
				return true
			}
		}
	}
	return false
}