package pool

import (
	"github.com/luopengift/golibs/logger"
	"sync"
	"testing"
	"time"
)

func Test_pool(t *testing.T) {
	var wg sync.WaitGroup
	factory := func() (interface{}, error) {
		var i int = 0
		return &i, nil
	}
	p := NewPool(1, 2, 2, factory)
	p.LogLevel(logger.DEBUG)
	logger.Info("pool init success...")
	time.Sleep(1 * time.Second)
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			one, err := p.Get()
			if err != nil {
				logger.Error("pool get error:%v", err)
			}
			resp := one.(*int)
			logger.Info("connID:%p,status:%v", one, resp)
			err = p.Put(one)
			if err != nil {
				logger.Error("PUT:%#v,%#V", one, err)
			}
		}()
	}
	wg.Wait()
	logger.Info("success")
}
