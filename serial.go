package serial

import (
	"sync"
)

var DefaultExecSync = func(f func()) {
	f()
}

var DefaultExecAsync = func(f func()) {
	go func() {
		f()
	}()
}

type job struct {
	f    func()
	next *job
}

var jobPool = sync.Pool{
	New: func() interface{} {
		return &job{}
	},
}

func getJob() *job {
	return jobPool.Get().(*job)
}

func putJob(jo *job) {
	jo.f = nil
	jo.next = nil
	jobPool.Put(jo)
}

var serialPool = sync.Pool{
	New: func() interface{} {
		return &Serial{}
	},
}

type Serial struct {
	mux       sync.Mutex
	head      *job
	tail      *job
	execSync  func(f func())
	execAsync func(f func())
	// closed    bool
}

func (s *Serial) Go(f func()) {
	if f == nil {
		return
	}
	jo := getJob()
	jo.f = f

	s.mux.Lock()
	if s.tail != nil {
		s.tail.next = jo
		s.tail = jo
		s.mux.Unlock()
		return
	}

	s.head = jo
	s.tail = jo
	s.mux.Unlock()

	s.execAsync(func() {
		var next *job
		for {
			s.execSync(jo.f)
			s.mux.Lock()
			next = jo.next
			putJob(jo)
			if next == nil {
				s.head = nil
				s.tail = nil
				s.mux.Unlock()
				return
			}
			jo = next
			s.head = jo
			s.mux.Unlock()
		}
	})
}

func New(execSync, execAsync func(f func())) *Serial {
	if execSync == nil {
		execSync = DefaultExecSync
	}
	if execAsync == nil {
		execAsync = DefaultExecAsync
	}
	return &Serial{
		execSync:  execSync,
		execAsync: execAsync,
	}
}

type SerialFactory struct {
	execSync  func(f func())
	execAsync func(f func())
}

func (sf *SerialFactory) Get() *Serial {
	s := serialPool.Get().(*Serial)
	s.execSync = sf.execSync
	s.execAsync = sf.execAsync
	return s
}

func (sf *SerialFactory) Put(s *Serial) {
	s.head = nil
	s.tail = nil
	serialPool.Put(s)
}

func NewFactory(execSync, execAsync func(f func())) *SerialFactory {
	if execSync == nil {
		execSync = DefaultExecSync
	}
	if execAsync == nil {
		execAsync = DefaultExecAsync
	}
	return &SerialFactory{
		execSync:  execSync,
		execAsync: execAsync,
	}
}
