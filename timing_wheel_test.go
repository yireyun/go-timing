// timing_sheel_test
package timing

import (
	"testing"
	"time"
)

func Test(t *testing.T) {
	TimingWheel := NewTimingWheel(time.Millisecond*10, 600)
	c := TimingWheel.After(time.Millisecond)
	if _, ok := <-c; ok {
		t.Errorf("TimingWheel After")
	}
	start := time.Now()
	c1 := TimingWheel.After(time.Millisecond * 10)
	c2 := TimingWheel.After(time.Millisecond * 10 * 2)
	c3 := TimingWheel.After(time.Millisecond * 10 * 3)
	for i := 0; i < 3; i++ {
		select {
		case <-c1:
			end := time.Now()
			t.Logf("After 10 Millisecond, Use:%v", end.Sub(start))
			c1 = nil
		case <-c2:
			end := time.Now()
			t.Logf("After 20 Millisecond, Use:%v", end.Sub(start))
			c2 = nil
		case <-c3:
			end := time.Now()
			t.Logf("After 30 Millisecond, Use:%v", end.Sub(start))
			c3 = nil
		}
	}
}

func BenchmarkAfterAsyncDo(b *testing.B) {
	tw := NewTimingWheel(time.Second, 100)
	for i := 0; i < b.N; i++ {
		tw.After(time.Second)
	}
}

func BenchmarkAfterAsyncGo(b *testing.B) {
	tw := NewTimingWheel(time.Second, 100)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			tw.After(time.Second)
		}
	})
}
