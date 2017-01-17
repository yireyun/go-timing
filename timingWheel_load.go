// +build arm arm64 mips64 mips64le ppc64 ppc64le

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
	pos := atomic.LoadUint32(&w.position)

	wsp := &w.timeScales[(pos+uint32(timeout/w.interval))&w.bucketMod]

	ws := atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(wsp)))

	return *((*waitSignal)(ws))
}

func (w *TimingWheel) onTicker() {
	pos := atomic.LoadUint32(&w.position)

	wsp := &w.timeScales[pos]

	ws := atomic.SwapPointer((*unsafe.Pointer)(unsafe.Pointer(wsp)),
		unsafe.Pointer(w.preSignal))

	atomic.SwapUint32(&w.position, (pos+1)&w.bucketMod)

	close(*((*waitSignal)(ws)))
	var prepare waitSignal = make(chan struct{})
	w.preSignal = &prepare
}
