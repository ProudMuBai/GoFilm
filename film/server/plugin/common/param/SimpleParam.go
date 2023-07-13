package param

import (
	"reflect"
)

/*
	简单参数处理
*/

// IsEmpty 判断各种基本类型是否为空
func IsEmpty(target any) bool {
	switch target.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return target == reflect.Zero(reflect.TypeOf(target)).Interface()
	case string:
		return target == ""
	case bool:
		return target.(bool)
	default:
		return false
	}
}

func IsEmptyRe[T ~int | ~uint | float32 | float64 | string | bool](target T) bool {
	v := reflect.ValueOf(target)
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Float32,
		reflect.Float64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Interface() == 0
	case reflect.String:
		return v.String() == ""
	case reflect.Bool:
		return v.Bool()
	default:
		return false
	}
}
