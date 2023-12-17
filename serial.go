package serial

import (
	"log"
	"runtime"
	"sync"
	"unsafe"
)

type Serial struct {
	idx     int
	jobs    []func()
	mux     sync.Mutex
	running bool
	caller  func(f func())
}

func (s *Serial) Go(f func()) {
	s.mux.Lock()
	running := s.running
	s.running = true
	s.jobs = append(s.jobs, f)
	s.mux.Unlock()

	if !running {
		go s.doAll()
	}
}

func (s *Serial) doAll() {
	caller := s.caller
	if caller == nil {
		caller = defaultCaller
	}
	var f func()
	for {
		s.mux.Lock()
		if s.idx < len(s.jobs) {
			f = s.jobs[s.idx]
			s.jobs[s.idx] = nil
			s.idx++
		} else {
			s.idx = 0
			s.jobs = s.jobs[:0]
			s.running = false
			s.mux.Unlock()
			return
		}
		s.mux.Unlock()
		caller(f)
	}
}

func New(caller func(f func())) *Serial {
	if caller == nil {
		caller = defaultCaller
	}
	return &Serial{jobs: make([]func(), 0, 10000), caller: caller}
}

func defaultCaller(f func()) {
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
