package serial

import "sync/atomic"

type SerialFactory struct {
	inited    int32
	initSize  int
	maxSize   int
	execSync  func(f func())
	execAsync func(f func())
}

func (sf *SerialFactory) Get() *Serial {
	if atomic.CompareAndSwapInt32(&sf.inited, 0, 1) {
		if sf.execSync == nil {
			sf.execSync = DefaultExecSync
		}
		if sf.execAsync == nil {
			sf.execAsync = DefaultExecAsync
		}
	}
	s := serialPool.Get().(*Serial)
	s.execSync = sf.execSync
	s.execAsync = sf.execAsync
	s.maxSize = sf.maxSize
	s.jobs = make([]*job, sf.initSize)[:0]
	return s
}

func (sf *SerialFactory) Put(s *Serial) {
	s.jobs = nil
	serialPool.Put(s)
}

func NewFactory(execSync, execAsync func(f func()), initQueueSize, maxQueueSize int) *SerialFactory {
	if execSync == nil {
		execSync = DefaultExecSync
	}
	if execAsync == nil {
		execAsync = DefaultExecAsync
	}
	if initQueueSize <= 0 {
		initQueueSize = DefaultInitSize
	}
	if maxQueueSize > 0 && maxQueueSize < initQueueSize {
		maxQueueSize = initQueueSize * 2
	}

	return &SerialFactory{
		execSync:  execSync,
		execAsync: execAsync,
		initSize:  initQueueSize,
		maxSize:   maxQueueSize,
	}
}

var DefaultFactory = NewFactory(nil, nil, DefaultInitSize, DefaultMaxSize)

func Get() *Serial {
	return DefaultFactory.Get()
}

func Put(s *Serial) {
	DefaultFactory.Put(s)
}
