package channel

import (
	"fmt"
	"testing"
	"time"
)

func Test_channel(t *testing.T) {
	q := NewChannel(10)
	//启动协程查看信息
	go func() {
		for {
			fmt.Println(q)
			time.Sleep(500 * time.Millisecond)
		}
	}()
	for i := 0; i < 20; i++ {
			q.Run(func() error {
				fmt.Println(fmt.Sprintf("groutine no.%d start,time %v", i, time.Now().Format("15:04:05")))
				time.Sleep(2 * time.Second)
				fmt.Println(fmt.Sprintf("groutine no.%d end,time %v", i, time.Now().Format("15:04:05")))
				return nil
			})
	}
    select{}
}
