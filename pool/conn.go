package pool

import (
	//"github.com/luopengift/golibs/logger"
	"time"
)

type Ctx struct {
	conn interface{}
	time time.Time
}

func NewCtx(conn interface{}) *Ctx {
	ctx := new(Ctx)
	ctx.conn = conn
	ctx.time = time.Now()
	return ctx
}
