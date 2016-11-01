/*
cache 用于缓存游戏动态数据。
*/

package cache

import (
	"server/base/service"
)

////////////////////////////////////////////
// type, const, var
//
var (
	srv = service.NewService0()
)

////////////////////////////////////////////
// func
//
func init() {
	srv.Start()
}

func Service() *service.Service {
	return srv
}
