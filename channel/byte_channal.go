package channel

import (
	"fmt"
	"sync"
)

type ByteChannel struct {
	mux     *sync.Mutex
	channel chan []byte //使用chan堵塞控制队列的最大数量
	max     int64       //队列最大数量
	total   int64       //通过队列总数量
}

func NewByteChannel(max int64) *ByteChannel {
	channel := &ByteChannel{
		mux:     new(sync.Mutex),
		channel: make(chan []byte, max),
		max:     max,
		total:   0,
	}
	return channel
}

func (self *ByteChannel) String() string {
	return fmt.Sprintf("<max:%d,total:%d,idle:%d>", self.max, self.total, self.Idle())
}

func (self *ByteChannel) Close() {
	close(self.channel)
}

//计数
func (self *ByteChannel) Count() {
	self.mux.Lock()
	self.total = self.total + 1
	self.mux.Unlock()
}

//往管道中写数据
func (self *ByteChannel) Recv() chan []byte {
	self.Count()
	return self.channel
}

//从管道中读数据
func (self *ByteChannel) Send() <-chan []byte {
	return self.channel
}

//往管道中写数据
func (self *ByteChannel) Put(v []byte) {
	self.channel <- v
	self.Count()
}

//从管道中读数据
func (self *ByteChannel) Get() []byte {
	v := <-self.channel
	return v
}

func (self *ByteChannel) Idle() int64 {
	return self.max - int64(len(self.channel))
}

func (self *ByteChannel) Total() int64 {
	return self.total
}
