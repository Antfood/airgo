package testutils

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

func Assert(tb testing.TB, conditional bool, msg string, v ...any) {
	if !conditional {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: "+msg+"\033[39m\n\n", append([]any{filepath.Base(file), line}, v...)...)
		tb.FailNow()
	}
}

func Ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

func Equals(tb testing.TB, exp, act any) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}
