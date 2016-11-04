/*
cache 用于缓存游戏动态数据。
*/
package cache

import (
	"server/base/service"

	"github.com/name5566/leaf/log"
)

////////////////////////////////////////////
// type, const, var
//
var (
	srv             = service.NewService1(update)
	curTime float32 = 0
)

////////////////////////////////////////////
// func
//
func init() {
	srv.Start()
}

func update(dt float32) {
	curTime += dt
	if curTime > CheckOnlineNumInterval {
		curTime = 0
		checkOnlineNum()
	}
}

func checkOnlineNum() {
	log.Release("online player num: %d", getInnerPlayerNum())
}
