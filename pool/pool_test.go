package pool

import (
	"github.com/luopengift/gohttp"
	"github.com/luopengift/golibs/logger"
	"testing"
	"time"
)

func Test_pool(t *testing.T) {
	factory := func() (interface{}, error) {
		client := gohttp.NewClient().Url("http://www.baidu.com")
		logger.Info("create conn:%p",client)
		return client, nil
	}
	p := NewPool(1, 2, 2, factory)
	p.LogLevel(logger.DEBUG)
	logger.Info("pool init success...")
	time.Sleep(4*time.Second)
	for i := 0; i < 4; i++ {
		go func() {
			one, err := p.Get()
			if err != nil {
				logger.Error("pool get error:%v", err)
			}
			resp, err := one.(*gohttp.Client).Get()
			if err != nil {
				logger.Error("http get error:%v", err)
			}
			logger.Info("connID:%p,status:%v", one, resp.StatusCode)
			err = p.Put(one)
			if err != nil {
				logger.Error("PUT:%#v,%#V", one, err)
			}
		}()
	}
	select {}
}
