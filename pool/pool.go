package pool

import (
	"fmt"
	"sync"
)

type Pool struct {
	mux   *sync.Mutex
	queue chan interface{} //使用chan堵塞控制协程的最大数量
	max   int64            //协程最大数量
	cur   int64            //当前运行协程数量
	cnt   int64            //启动协程的总数量
}

func NewPool(max int64, v interface{}) *Pool {
	pool := &Pool{
		mux:   new(sync.Mutex),
		queue: make(chan interface{}, max),
		max:   max,
		cur:   0,
		cnt:   0,
	}
	var i int64 = 0
	for ; i < max; i++ {
		pool.Put(v)
	}
	return pool
}

func (self *Pool) Init() {}

func (self *Pool) String() string {
	return fmt.Sprintf("<max:%d,cnt:%d,current:%d>", self.max, self.cnt, self.cur)
}

func (self *Pool) Put(v interface{}) {
	self.queue <- v
	self.mux.Lock()
	self.cur += 1
	self.cnt += 1
	self.mux.Unlock()
}

func (self *Pool) Get() (v interface{}) {
	v = <-self.queue
	self.mux.Lock()
	self.cur -= 1
	self.mux.Unlock()
	return
}

func (self *Pool) Wait() {
	for {
		switch self.cur {
		case 0:
			return
		default:
			continue
		}
	}
}

func (self *Pool) Hold() {
	select {}
}

func (self *Pool) Number() int64 {
	return self.cur
}

func (self *Pool) Count() int64 {
	return self.cnt
}

func (self *Pool) Run(fn func() error) {
	v := self.Get()
	fn()
	self.Put(v)
}
