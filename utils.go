package gotq

import (
	"reflect"
	"runtime"
)

// Get func name , see detail
// http://stackoverflow.com/questions/7052693/how-to-get-the-name-of-a-function-in-go
func GetFuncName(f interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}
