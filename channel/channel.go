package channel

import (
	"fmt"
	"sync"
)

type Channel struct {
	mux   *sync.Mutex
	channel chan interface{} //使用chan堵塞控制队列的最大数量
	max   int64            //队列最大数量
	idle  int64            //当前队列空闲数量
	total int64            //通过队列总数量
}

func NewChannel(max int64) *Channel {
	channel := &Channel{
		mux:   new(sync.Mutex),
		channel: make(chan interface{}, max),
		max:   max,
		idle:  max,
		total: 0,
	}
	return channel
}

func (self *Channel) String() string {
	return fmt.Sprintf("<max:%d,total:%d,idle:%d>", self.max, self.total, self.idle)
}

func (self *Channel) Close() {
    close(self.channel)
}
//接收数据
func (self *Channel) recv() {
	self.mux.Lock()
	self.idle = self.idle - 1
	self.total = self.total + 1
	self.mux.Unlock()
}
//从管道中读取数据
func (self *Channel) send() {
	self.mux.Lock()
	self.idle = self.idle + 1
	self.mux.Unlock()
}

func (self *Channel) Send() (chan interface{}) {
    self.send()
    return self.channel
}


func (self *Channel) Recv() (chan interface{}) {
    self.recv()
    return self.channel
}


func (self *Channel) Put(v interface{}) {
	self.channel <- v
    self.recv()
}

func (self *Channel) Add() {
    self.Put(true)
}

func (self *Channel) Get() interface{} {
	v := <-self.channel
    self.send()
    return v
}

func (self *Channel) Done() {
    self.Get()
}


func (self *Channel) Idle() int64 {
	return self.idle
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
