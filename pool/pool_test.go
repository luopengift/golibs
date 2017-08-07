package pool

import (
        "github.com/luopengift/golibs/logger"
        "testing"
        "time"
)



func Test_pool(t *testing.T) {
       factory := func() (interface{}, error) {
		var i int = 0
                return &i, nil
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
                        resp := one.(*int)
                        logger.Info("connID:%p,status:%v", one, resp)
                        err = p.Put(one)
                        if err != nil {
                                logger.Error("PUT:%#v,%#V", one, err)
                        }
                }()
        }
        select {}
}


