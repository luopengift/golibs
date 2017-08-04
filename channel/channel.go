package channel

import (
	"fmt"
	"sync"
)

type Channel struct {
	mux     *sync.Mutex
	channel chan interface{} //使用chan堵塞控制队列的最大数量
	max     int64            //队列最大数量
	total   int64            //通过队列总数量
}

func NewChannel(max int64) *Channel {
	channel := &Channel{
		mux:     new(sync.Mutex),
		channel: make(chan interface{}, max),
		max:     max,
		total:   0,
	}
	return channel
}

func (self *Channel) String() string {
	return fmt.Sprintf("<max:%d,total:%d,idle:%d>", self.max, self.total, self.Idle())
}

func (self *Channel) Close() {
	close(self.channel)
}

//计数
func (self *Channel) Count() {
	self.mux.Lock()
	self.total = self.total + 1
	self.mux.Unlock()
}

//往管道中写数据
func (self *Channel) Put(v interface{}) {
	self.channel <- v
	self.Count()
}

//从管道中读数据
func (self *Channel) Get() interface{} {
	v := <-self.channel
	return v
}

//往管道中写数据
func (self *Channel) Add() {
	self.Put(struct{}{})
}

//从管道中读数据
func (self *Channel) Done() {
	self.Get()
}

func (self *Channel) Idle() int64 {
	return self.max - int64(len(self.channel))
}

func (self *Channel) Total() int64 {
	return self.total
}

func (self *Channel) Run(fun func() error) {
	self.Add()
	go func() {
		fun()
		self.Done()
	}()
}
