package middleware

import (
	"errors"
	"time"
)

// 实现context.Context接口
type ExitContext struct {
	Chan         chan struct{}
	DeadLineTime time.Time
}

func NewExitContext(duration time.Duration) *ExitContext {
	cxt := ExitContext{}
	cxt.DeadLineTime = time.Now().Add(duration)
	cxt.Chan = make(chan struct{}, 1)
	return &cxt
}

func (cxt *ExitContext) Done() <-chan struct{} {
	if time.Now().After(cxt.DeadLineTime) {
		cxt.Chan <- struct{}{}
	}
	return cxt.Chan
}

func (cxt *ExitContext) Err() error {
	return errors.New("can't exit before Specified time")
}

// 无意义的空函数
func (cxt *ExitContext) Value(key interface{}) interface{} {
	return nil
}

func (ctx *ExitContext) Deadline() (deadline time.Time, ok bool) {
	deadline = ctx.DeadLineTime
	ok = true
	return
}
