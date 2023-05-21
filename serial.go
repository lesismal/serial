package serial

import (
	"sync"
)

var serialPool = sync.Pool{
	New: func() interface{} {
		return &Serial{}
	},
}

type Serial struct {
	mux sync.Mutex

	jobs      []*job
	maxSize   int
	execSync  func(f func())
	execAsync func(f func())
}

func (s *Serial) Go(f func()) {
	if f == nil {
		panic("nil value func")
	}

	jo := getJob()
	jo.f = f

	s.mux.Lock()
	isHead := (len(s.jobs) == 0)
	s.jobs = append(s.jobs, jo)
	s.mux.Unlock()

	if isHead {
		s.execAsync(func() {
			s.doAll(jo)
		})
	}
}

func (s *Serial) GoWithValue(fv func(interface{}), v interface{}) {
	if fv == nil {
		panic("nil func")
	}

	jo := getJob()
	jo.fv = fv
	jo.v = v

	s.mux.Lock()
	isHead := (len(s.jobs) == 0)
	s.jobs = append(s.jobs, jo)
	s.mux.Unlock()

	if isHead {
		s.execAsync(func() {
			s.doAll(jo)
		})
	}
}

func (s *Serial) doAll(jo *job) {
	i := 0
	for {
		jo2 := jo
		f := func() {
			s.execSync(func() {
				defer putJob(jo2)
				if jo2.f != nil {
					jo2.f()
				}
				if jo2.fv != nil {
					jo2.fv(jo2.v)
				}
			})
		}
		s.mux.Lock()
		s.jobs[i] = nil
		i++
		if len(s.jobs) == i {
			s.jobs = s.jobs[0:0]
			s.mux.Unlock()
			f()
			return
		}
		jo = s.jobs[i]
		s.mux.Unlock()
		f()
	}
}
