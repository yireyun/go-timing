package timing

import (
	"sync/atomic"
	"time"
	"unsafe"
)

type waitSignal chan struct{}

func newWaitSignal() *waitSignal {
	var ws waitSignal = make(chan struct{})
	return &ws
}

//时间轮
type TimingWheel struct {
	bucketCnt  uint32        //时间轮总刻度数量
	bucketMod  uint32        //时间轮总刻度取模
	interval   time.Duration //时间轮间隔时间
	maxTimeout time.Duration //最大超时
	ticker     *time.Ticker  //时间轮定时器
	quitChan   chan struct{} //关闭时间轮
	position   uint32        //时间轮当前位置
	timeScales []*waitSignal //时间刻度(数组管道),struct{}的长度为0，减少内存
	preSignal  *waitSignal   //预备等待信号量

}

//创建时间轮
//interval	是输入精度
//buckets	是输入时间槽数量
func NewTimingWheel(interval time.Duration, buckets uint32) *TimingWheel {
	w := new(TimingWheel)

	w.bucketCnt = minQuantity(buckets)
	w.bucketMod = w.bucketCnt - 1
	w.interval = interval

	w.maxTimeout = interval * time.Duration(w.bucketCnt)

	w.ticker = time.NewTicker(interval)
	w.quitChan = make(chan struct{})

	w.position = 0
	w.timeScales = make([]*waitSignal, w.bucketCnt)
	for i := range w.timeScales {
		w.timeScales[i] = newWaitSignal()
	}
	w.preSignal = newWaitSignal()

	go w.run()
	return w
}

func (w *TimingWheel) Stop() {
	close(w.quitChan)
}

//定时返回
//timeout	输入超时时间，最大时间=精度*时间槽数量
func (w *TimingWheel) After(timeout time.Duration) <-chan struct{} {
	if timeout >= w.maxTimeout {
		panic("timeout greater than max size")
	}

	return *w.timeScales[(w.position+uint32(timeout/w.interval))&w.bucketMod]
}

func (w *TimingWheel) run() {
	for {
		select {
		case <-w.ticker.C:
			w.onTicker()
		case <-w.quitChan:
			w.ticker.Stop()
			return
		}
	}
}

func (w *TimingWheel) onTicker() {
	pos := w.position

	wsp := &w.timeScales[pos]

	ws := atomic.SwapPointer(
		(*unsafe.Pointer)(unsafe.Pointer(wsp)), unsafe.Pointer(w.preSignal))
	atomic.SwapUint32(&w.position, (w.position+1)&w.bucketMod)
	close(*((*waitSignal)(ws)))
	var prepare waitSignal = make(chan struct{})
	w.preSignal = &prepare
}

// round 到最近的2的倍数
func minQuantity(v uint32) uint32 {
	v--
	v |= v >> 1
	v |= v >> 2
	v |= v >> 4
	v |= v >> 8
	v |= v >> 16
	v++
	return v
}
