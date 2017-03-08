package pool

import (
	"fmt"
	"sync"
)

type Pool struct {
	mux     *sync.Mutex
	queue   chan int //使用chan堵塞控制协程的最大数量
	max     int64    //协程最大数量
	current int64    //当前运行协程数量
	cnt     int64    //启动协程的总数量
}

func NewPool(max int64) *Pool {
	return &Pool{
		mux:     new(sync.Mutex),
		queue:   make(chan int, max),
		max:     max,
		current: 0,
		cnt:     0,
	}
}

func (self *Pool) String() string {
	return fmt.Sprintf("<max:%d,cnt:%d,current:%d>", self.max, self.cnt, self.current)
}

func (self *Pool) add() {
	self.queue <- 1
	self.mux.Lock()
	self.current += 1
	self.cnt += 1
	self.mux.Unlock()
}

func (self *Pool) done() {
	<-self.queue
	self.mux.Lock()
	self.current -= 1
	self.mux.Unlock()
}

func (self *Pool) Wait() {
LOOP:
	switch self.current {
	case 0:
		return
	default:
		goto LOOP
	}
}

func (self *Pool) Hold() {
	select {}
}

func (self *Pool) Number() int64 {
	return self.current
}

func (self *Pool) Count() int64 {
	return self.cnt
}

func (self *Pool) Run(fn func() error) {
	self.add()
	fn()
	self.done()
}
