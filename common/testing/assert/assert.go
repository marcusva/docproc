// Package assert provides minimalistic testing enhancements
package assert

import (
	"fmt"
	"reflect"
	"runtime"
	"testing"
)

func printErr(t *testing.T, defmsg string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(2)
	msg := defmsg
	if len(args) > 0 {
		// If args are passed in, ignore the default message
		msg = fmt.Sprintf(args[0].(string), args[1:]...)
	}
	if ok {
		// call stack available
		t.Errorf("%v:%v:\n%v", file, line, msg)
	} else {
		t.Errorf("[no stacktrace]:\n%v", msg)
	}
}

// Equal checks two values on equality. This performs a deep object equality
// check via reflect.DeepEqual. The additional args are used for a customized
// error output, if the values are not equal. If no args are provided, a simple
// standard message will be printed via t.Errorf().
func Equal(t *testing.T, v1, v2 interface{}, args ...interface{}) {
	if v1 == nil && v2 == nil {
		return
	}
	if !reflect.DeepEqual(v1, v2) {
		printErr(t, fmt.Sprintf("Assertion failure (equal): '%v' == '%v'",
			v1, v2), args...)
	}
}

// NotEqual Checks two values on inequality. This performs a deep object
// equality check via reflect.DeepEqual. The additional args are used for a
// customized error output, if the values are not equal. If no args are
// provided, a simple standard message will be printed via t.Errorf().
func NotEqual(t *testing.T, v1, v2 interface{}, args ...interface{}) {
	if reflect.DeepEqual(v1, v2) {
		printErr(t, fmt.Sprintf("Assertion failure (not equal): '%v' != '%v'",
			v1, v2), args...)
	}
}

// Nil checks, if the passed value is nil. The additional args are used for a
// customized error output, if the value is not nil. If no args are provided,
// a simple standard message will be printed via t.Errorf().
func Nil(t *testing.T, v interface{}, args ...interface{}) {
	if v != nil {
		printErr(t, fmt.Sprintf("Assertion failure (nil): nil == '%v'", v),
			args...)
	}
}

// NotNil checks, if the passed value is not nil. The additional args are used
// for a customized error output, if the value is nil. If no args are provided,
// a simple standard message will be printed via t.Errorf().
func NotNil(t *testing.T, v interface{}, args ...interface{}) {
	if v == nil {
		printErr(t, fmt.Sprintf("Assertion failure (not nil): nil != '%v' ",
			v), args...)
	}
}

// Err checks, if the passed error is not nil.
func Err(t *testing.T, err error) {
	if err == nil {
		printErr(t, "", fmt.Sprintf("%v", err))
	}
}

// NoErr checks, if the passed error is nil.
func NoErr(t *testing.T, err error) {
	if err != nil {
		printErr(t, "", fmt.Sprintf("%v", err))
	}
}

// FailOnErr fails the test via t.FailNow(), if the passed error is not nil.
func FailOnErr(t *testing.T, err error) {
	if err != nil {
		printErr(t, "", fmt.Sprintf("%v", err))
		t.FailNow()
	}
}

// FailIf fails the test via t.FailNow(), if condition evaluates to true. The
// additional args are used for a customized error output, if the condition
// is met. If no args are  provided, a simple standard message will be printed
// via t.Errorf().
func FailIf(t *testing.T, condition bool, args ...interface{}) {
	if condition {
		printErr(t, "Condition not met", args...)
		t.FailNow()
	}
}

// FailIfNot fails the test, if condition evaluates to false. The additional
// args are used for a customized error output, if the condition is not met.
// If no args are provided, a simple standard message will be printed via
// t.Errorf().
func FailIfNot(t *testing.T, condition bool, args ...interface{}) {
	if !condition {
		printErr(t, "Condition not met", args...)
		t.FailNow()
	}
}

// Panics checks, if the passed function fn panics.
func Panics(t *testing.T, fn func()) {
	panic := false
	var msg interface{}
	exec := func() {
		defer func() {
			if msg = recover(); msg != nil {
				panic = true
			}
		}()
		fn()
	}
	exec()
	FailIfNot(t, panic, "function did not panic")
}

// ContainsS checks, if array contains the passed value. The additional args are
// used for a customized error output, if the values are not equal. If no args
// are provided, a simple standard message will be printed via t.Errorf().
func ContainsS(t *testing.T, array []string, val string, args ...interface{}) {
	for _, v := range array {
		if v == val {
			return
		}
	}
	printErr(t, fmt.Sprintf("Array failure: '%v' does not contain '%v' ", array, val), args...)
	t.FailNow()
}
