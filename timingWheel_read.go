// +build 386 amd64 amd64p32 s390x

package timing

import (
	"sync/atomic"
	"time"
	"unsafe"
)

//定时返回
//timeout	输入超时时间，最大时间=精度*时间槽数量
func (w *TimingWheel) After(timeout time.Duration) <-chan struct{} {
	if timeout > w.maxTimeout {
		timeout = w.maxTimeout
	}
	return *w.timeScales[(w.position+uint32(timeout/w.interval))&w.bucketMod]
}

func (w *TimingWheel) onTicker() {
	pos := w.position

	wsp := &w.timeScales[pos]

	ws := atomic.SwapPointer(
		(*unsafe.Pointer)(unsafe.Pointer(wsp)), unsafe.Pointer(w.preSignal))
	atomic.SwapUint32(&w.position, (pos+1)&w.bucketMod)

	close(*((*waitSignal)(ws)))
	var prepare waitSignal = make(chan struct{})
	w.preSignal = &prepare
}
