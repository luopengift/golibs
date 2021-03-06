package pool

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/luopengift/golibs/channel"
	"github.com/luopengift/log"
)

// Factory factory
type Factory func() (interface{}, error)

// Pool pool
type Pool struct {
	mutex   *sync.Mutex
	maxIdle int              //最大空闲数
	maxOpen int              //最大连接数
	timeout time.Duration    //链接最大存活时间
	factory Factory          //连接生成方式
	pool    chan *Ctx        //连接存放的channel
	channel *channel.Channel //并发最大连接控制
	*log.Log
}

// NewPool new pool
func NewPool(maxIdle, maxOpen, timeout int, factory Factory) *Pool {
	if maxIdle < 0 || maxOpen <= 0 || maxIdle > maxOpen {
		log.Error("maxIdle or maxOpen args error!")
		return nil
	}
	p := new(Pool)
	p.mutex = new(sync.Mutex)
	p.maxIdle = maxIdle
	p.maxOpen = maxOpen
	p.timeout = time.Duration(timeout) * time.Second
	p.factory = factory
	p.pool = make(chan *Ctx, p.maxOpen)
	p.channel = channel.NewChannel(p.maxOpen)
	p.Log = log.NewLog("2006/01/02 15:04:05.000", os.Stdout)
	for i := 0; i < p.maxIdle; i++ {
		if err := p.create(); err != nil {
			p.Log.Error("%v", err)
		}
	}
	return p
}

// LogLevel LogLevel
func (p *Pool) LogLevel(lv uint8) {
	p.Log.SetLevel(lv)
}

//生成一个新的连接,放入连接池中
func (p *Pool) create() error {
	if p.channel.Len() >= p.maxOpen {
		return fmt.Errorf("conn pool is full,can not create new conn!")
	}
	p.channel.Add()
	conn, err := p.factory()
	if err != nil {
		return fmt.Errorf("create a new conn error!%v", err)
	}
	if err = p.Put(conn); err != nil {
		return fmt.Errorf("conn can not put into pool!%v", err)
	}
	return nil
}

// Get 从Pool中取出一个连接
func (p *Pool) Get() (interface{}, error) {
	for {
		select {
		case ctx := <-p.pool:
			if ctx.time.Add(p.timeout).Before(time.Now()) {
				p.Log.Debug("Get one is timeout,release:%p", ctx.conn)
				p.release(ctx)
				continue
			}
			p.Log.Debug("GET one from full pool:%#v", ctx)
			return ctx.conn, nil

		default:
			if p.channel.Len() < p.maxOpen {
				p.Log.Debug("pool is null,create one")
				p.create()
				continue
			} else {
				p.Log.Debug("all conn is used,please wait...")
				ctx := <-p.pool
				if ctx.time.Add(p.timeout).Before(time.Now()) {
					p.release(ctx)
					continue
				}
				p.Log.Debug("GET one from full pool:%#v", ctx)
				return ctx.conn, nil
			}
		}
	}
}

// 将连接放入Pool中
func (p *Pool) Put(conn interface{}) error {
	if conn == nil {
		return fmt.Errorf("conn is nil")
	}
	select {
	case p.pool <- NewCtx(conn):
		p.Log.Debug("PUT conn into pool:%p", conn)
		return nil
	default:
		//连接池已满，直接关闭该链接
		p.Log.Warn("pool if full,release:%p", conn)
		return p.release(conn)
	}
	return nil
}

// 释放指定连接
func (p *Pool) release(ctx interface{}) error {
	p.channel.Done()
	return nil
}

func (p *Pool) Close() error {
	return nil
}
