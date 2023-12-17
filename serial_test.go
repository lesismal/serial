package serial

import (
	"testing"
)

func Test_Serial(t *testing.T) {
	s := New(nil)
	for i := 0; i < 1; i++ {
		loop := 10000
		chValue := make(chan int, loop)
		for j := 0; j < loop; j++ {
			func(n int) {
				s.Go(func() {
					chValue <- n
				})
			}(j)
		}

		for j := 0; j < loop; j++ {
			n := <-chValue
			if n != j {
				t.Fatalf("want: %v, got: %v", j, n)
			}
		}
	}
}
