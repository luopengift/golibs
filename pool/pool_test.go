package pool

import (
	"sync"
	"testing"
	"time"

	"github.com/luopengift/log"
)

func Test_pool(t *testing.T) {
	var wg sync.WaitGroup
	factory := func() (interface{}, error) {
		var i int = 0
		return &i, nil
	}
	p := NewPool(1, 2, 2, factory)
	p.LogLevel(log.DEBUG)
	log.Info("pool init success...")
	time.Sleep(1 * time.Second)
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			one, err := p.Get()
			if err != nil {
				log.Error("pool get error:%v", err)
			}
			resp := one.(*int)
			log.Info("connID:%p,status:%v", one, resp)
			err = p.Put(one)
			if err != nil {
				log.Error("PUT:%#v,%#v", one, err)
			}
		}()
	}
	wg.Wait()
	log.Info("success")
}
