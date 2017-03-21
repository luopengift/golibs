package pool

import (
	"fmt"
	"testing"
	"time"
)

func Test_pool(t *testing.T) {
	pool := NewPool(10, nil)
	//启动协程查看信息
	go func() {
		for {
			fmt.Println(pool)
			time.Sleep(500 * time.Millisecond)
		}
	}()
	for i := 0; i < 20; i++ {
		go func(i int) {
			pool.Run(func() error {
				fmt.Println(fmt.Sprintf("groutine no.%d start,time %v", i, time.Now().Format("15:04:05")))
				time.Sleep(2 * time.Second)
				fmt.Println(fmt.Sprintf("groutine no.%d end,time %v", i, time.Now().Format("15:04:05")))
				return nil
			})
		}(i)
	}
	pool.Wait()
}
