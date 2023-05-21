package serial

import (
	"log"
	"runtime"
	"unsafe"
)

var (
	DefaultInitSize = 8

	DefaultMaxSize = 0

	DefaultExecSync = func(f func()) {
		defer func() {
			if err := recover(); err != nil {
				const size = 64 << 10
				buf := make([]byte, size)
				buf = buf[:runtime.Stack(buf, false)]
				log.Printf("execute failed: %v\n%v\n", err, *(*string)(unsafe.Pointer(&buf)))
			}
		}()
		f()
	}

	DefaultExecAsync = func(f func()) {
		go func() {
			f()
		}()
	}
)
