package pool

import (
	"fmt"
	"sync"
)

type Pool struct {
	mux   *sync.Mutex
	queue chan interface{} //使用chan堵塞控制协程的最大数量
	max   int64            //协程池最大数量
	idle  int64            //当前协程池空闲数量
	total int64            //启动协程的总数量
}

func NewPool(max int64, v interface{}) *Pool {
	pool := &Pool{
		mux:   new(sync.Mutex),
		queue: make(chan interface{}, max),
		max:   max,
		idle:  0,
		total: 0,
	}
	pool.Init(v)
	return pool
}

func (self *Pool) Init(v interface{}) *Pool {
	var i int64 = 0
	for ; i < self.max; i++ {
		self.Put(v)
	}
	return self
}

func (self *Pool) String() string {
	return fmt.Sprintf("<max:%d,total:%d,idle:%d>", self.max, self.total, self.idle)
}

func (self *Pool) Put(v interface{}) {
	self.queue <- v
	self.mux.Lock()
	self.idle = self.idle + 1
	self.mux.Unlock()
}

func (self *Pool) Get() (v interface{}) {
	v = <-self.queue
	self.mux.Lock()
	self.total = self.total + 1
	self.idle = self.idle - 1
	self.mux.Unlock()
	return
}

func (self *Pool) Wait() {
	for {
		switch self.idle {
		case self.max:
			fmt.Println(self)
			return
		default:
			continue
		}
	}
}

func (self *Pool) Hold() {
	select {}
}

func (self *Pool) Idle() int64 {
	return self.idle
}

func (self *Pool) Total() int64 {
	return self.total
}

func (self *Pool) Run(fn func() error) {
	v := self.Get()
	fn()
	self.Put(v)
}
