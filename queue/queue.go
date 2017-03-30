package crawler

import (
	"fmt"
	"sync"
)

type Queue struct {
	mux   *sync.Mutex
	queue chan interface{} //使用chan堵塞控制队列的最大数量
	max   int64            //队列最大数量
	idle  int64            //当前队列空闲数量
	total int64            //启动队列总数量
}

func NewQueue(max int64, v interface{}) *Queue {
	queue := &Queue{
		mux:   new(sync.Mutex),
		queue: make(chan interface{}, max),
		max:   max,
		idle:  max,
		total: 0,
	}
	return queue
}

func (self *Queue) String() string {
	return fmt.Sprintf("<max:%d,total:%d,idle:%d>", self.max, self.total, self.idle)
}

func (self *Queue) Close() {
    close(self.queue)
}

func (self *Queue) Put(v interface{}) {
	self.queue <- v
	self.mux.Lock()
	self.idle = self.idle - 1
	self.mux.Unlock()
}

func (self *Queue) Get() (v interface{}) {
	v = <-self.queue
	self.mux.Lock()
	self.total = self.total + 1
	self.idle = self.idle + 1
	self.mux.Unlock()
	return
}

func (self *Queue) Idle() int64 {
	return self.idle
}

func (self *Queue) Total() int64 {
	return self.total
}
