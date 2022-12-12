package serial

import (
	"testing"
	"time"
)

func Test_Serial(t *testing.T) {
	s := New(nil, nil)
	testSerial(t, s)

	sf := NewFactory(nil, nil)
	for i := 0; i < 3; i++ {
		s = sf.Get()
		testSerial(t, s)
		sf.Put(s)
	}
}

func testSerial(t *testing.T, s *Serial) {
	loop := 5
	chValue := make(chan int)
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
			t.Fatalf("want: %v, got: %v\n", i, n)
		}
	}
}
