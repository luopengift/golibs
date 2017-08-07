package channel

import (
	"fmt"
	"testing"
	"time"
)

func Test_channel(t *testing.T) {
	ch := NewChannel(10)
	//启动协程查看信息
	go func() {
		for {
			fmt.Println(ch)
			time.Sleep(500 * time.Millisecond)
		}
	}()
	go func() {
		for {
			fmt.Println(ch.Get())
			time.Sleep(200 * time.Millisecond)

		}
	}()
	for i := 0; i < 10; i++ {
		ch.Put(i)
	}
	time.Sleep(10 * time.Second)
	fmt.Println("err:", ch.Close())
}

func Benchmark_Channel(b *testing.B) {
	ch := NewChannel(1000)
	/*go func() {
		for {
			if _,ok := ch.Get(); !ok {
				return
			}
		}
	}()*/
	for i := 0; i < b.N; i++ {
		ch.Put("11111111111111111111111111111asdfasdfadfl;kj11111111111111111111111")
		ch.Get()
	}
	ch.Close()
}
