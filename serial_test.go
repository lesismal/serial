package serial

import (
	"testing"
	"time"
)

func Test_Serial(t *testing.T) {
	factory := NewFactory(nil, nil, 8, 0)
	for i := 0; i < 3; i++ {
		s := factory.Get()
		testGo(t, s)
		testGoWithValue(t, s)
		factory.Put(s)
	}
}

func testGo(t *testing.T, s *Serial) {
	loop := 5
	chValue := make(chan int, loop)
	for i := 0; i < loop; i++ {
		n := i
		s.Go(func() {
			chValue <- n
		})
	}

	time.Sleep(time.Second / 100)
	for i := 0; i < loop; i++ {
		n := <-chValue
		if n != i {
			t.Fatalf("want: %v, got: %v", i, n)
		}
	}
}

func testGoWithValue(t *testing.T, s *Serial) {
	loop := 5
	chValue := make(chan interface{}, loop)
	for i := 0; i < loop; i++ {
		n := i
		s.GoWithValue(func(v interface{}) {
			chValue <- v
		}, n)
	}

	time.Sleep(time.Second / 100)
	for i := 0; i < loop; i++ {
		v := <-chValue
		n, ok := v.(int)
		if !ok {
			t.Fatalf("invalid type v: %v", v)
		}
		if n != i {
			t.Fatalf("want: %v, got: %v", i, n)
		}
	}
}
