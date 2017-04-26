package queue

import (
	"fmt"
	"testing"
	"time"
)

func Test_queue(t *testing.T) {
	q := NewQueue(10)
	//启动协程查看信息
	go func() {
		for {
			fmt.Println(q)
			time.Sleep(500 * time.Millisecond)
		}
	}()
	for i := 0; i < 20; i++ {
		go func(i int) {
			q.Run(func() error {
				fmt.Println(fmt.Sprintf("groutine no.%d start,time %v", i, time.Now().Format("15:04:05")))
				time.Sleep(2 * time.Second)
				fmt.Println(fmt.Sprintf("groutine no.%d end,time %v", i, time.Now().Format("15:04:05")))
				return nil
			})
		}(i)
	}
    select{}
}
