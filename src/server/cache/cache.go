/*
cache 用于缓存游戏动态数据。
*/
package cache

import (
	"server/base/service"

	"github.com/name5566/leaf/log"
)

// checker
type Checker struct {
	curTime  float32
	interval float32
	check    func()
}

func (c *Checker) update(dt float32) {
	c.curTime += dt
	if c.curTime > c.interval {
		c.curTime = 0
		if c.check != nil {
			c.check()
		}
	}
}

//
var (
	srv      = service.NewService1(update)
	checkers = []*Checker{
		&Checker{
			interval: CheckOnlineNumInterval,
			check:    checkOnlineNum,
		},
		&Checker{
			interval: CheckPlayerAliveInterval,
			check:    checkInnerPlayerAliveTime,
		},
	}
)

func init() {
	srv.Start()
}

func update(dt float32) {
	for _, c := range checkers {
		c.update(dt)
	}
}

func checkOnlineNum() {
	log.Release("online player num: %d", getInnerPlayerNum())
}
