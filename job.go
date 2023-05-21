package serial

import "sync"

var jobEmpty job

type job struct {
	f  func()
	fv func(v interface{})
	v  interface{}
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
	*jo = jobEmpty
	jobPool.Put(jo)
}
