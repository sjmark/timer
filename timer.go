package timer

import (
	"fmt"
	"sort"
	"time"
	"runtime"

	"github.com/davecgh/go-spew/spew"
)

type nTimer struct {
	timers   []*timer
	stop     chan string
	location *time.Location
	close    chan *time.Timer
	begin    chan struct{}
}

type timer struct {
	cType   clockType
	running bool
	next    time.Time
	oType   string
	sec     time.Duration
	fn      func()
}

type clockType uint8

const (
	clockOnce    clockType = 1 // 一次性定时器 只执行一次
	clockForever clockType = 2 // 循环执行定时器
)

func NewCron() *nTimer {
	return &nTimer{
		timers:   nil,
		stop:     make(chan string),
		close:    make(chan *time.Timer),
		location: time.Now().Location(),
		begin:    make(chan struct{}),
	}
}

func (c *nTimer) AddOnce(ot string, fn func(), sec time.Duration) { c.add(ot, clockOnce, fn, sec) }

func (c *nTimer) AddForever(ot string, fn func(), sec time.Duration) { c.add(ot, clockForever, fn, sec) }

func (c *nTimer) Start() { c.run() }

func (c *nTimer) Stop(ot string) {

	var size = len(c.timers)

	for i := 0; i < size; i++ {
		if i > size {
			break
		}

		if c.timers[i].oType == ot {
			c.timers[i].running = false
			continue
		}
	}
}

//PrintPanicStack 产生panic时的调用栈打印
func printPanicStack(extras ...interface{}) {
	if x := recover(); x != nil {
		i := 0
		funcName, file, line, ok := runtime.Caller(i)
		for ok {
			fmt.Println("frame %v:[func:%v,file:%v,line:%v]\n", i, runtime.FuncForPC(funcName).Name(), file, line)
			i++
			funcName, file, line, ok = runtime.Caller(i)
		}

		for k := range extras {
			fmt.Println("EXRAS:%d DATA:%s\n", k, spew.Sdump(extras[k]))
		}
	}
}

func (c *nTimer) runWorking(j func()) {
	defer printPanicStack(j)
	j()
}

func (c *nTimer) run() {
	t := time.NewTicker(time.Millisecond * 100)
	n := 0

	for {
		select {
		case <-c.begin:
			n++
		case <-func() <-chan time.Time {
			if n > 0 {
				return t.C
			}
			return nil
		}():
			if len(c.timers) > 0 {
				cro := c.timers[0]

				if !cro.running {
					c.timers = append(c.timers[:0], c.timers[1:]...)
					n--
					continue
				}

				if time.Now().After(cro.next) {

					if cro.running {
						go c.runWorking(cro.fn)

						if cro.cType == clockOnce {
							c.timers = append(c.timers[:0], c.timers[1:]...)
							n--
						}

						if cro.cType == clockForever {
							c.timers[0].next = time.Now().Add(cro.sec)
							sort.Sort(clockTime(c.timers))
						}
					} else {
						c.timers = append(c.timers[:0], c.timers[1:]...)
						continue
					}
				}
			}
		}
	}
}

func (c *nTimer) add(ot string, cType clockType, fn func(), sec time.Duration) {

	c.timers = append(
		c.timers,
		&timer{
			oType:   ot,
			cType:   cType,
			fn:      fn,
			sec:     sec,
			next:    time.Now().Add(sec),
			running: true,
		},
	)

	sort.Sort(clockTime(c.timers))

	if !c.timers[0].running {
		c.timers = append(c.timers[:0], c.timers[1:]...)
	}

	c.begin <- struct{}{}
}

func (c *nTimer) now() time.Time { return time.Now().In(c.location) }

type clockTime []*timer

func (s clockTime) Len() int { return len(s) }

func (s clockTime) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s clockTime) Less(i, j int) bool {
	if s[i].next.IsZero() {
		return false
	}
	if s[j].next.IsZero() {
		return true
	}
	return s[i].next.Before(s[j].next)
}
