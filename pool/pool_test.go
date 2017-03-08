package pool

import (
	"fmt"
	"testing"
	"time"
)

func Test_pool(t *testing.T) {
	pool := NewPool(100)
	for i := 0; i < 1000; i++ {
		go pool.Run(func() error {
			fmt.Println(pool)
			time.Sleep(1 * time.Second)
			return nil
		})
	}
	pool.Wait()
}
