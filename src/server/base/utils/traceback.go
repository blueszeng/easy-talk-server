package utils

import (
	"runtime"

	"github.com/name5566/leaf/conf"
	"github.com/name5566/leaf/log"
)

func Traceback(printer func([]byte)) {
	if conf.LenStackBuf > 0 {
		buf := make([]byte, conf.LenStackBuf)
		l := runtime.Stack(buf, false)
		if printer != nil {
			printer(buf[:l])
		} else {
			log.Debug("%s", buf[:l])
		}
	} else {
		log.Error("Traceback, conf.LenStackBuf <= 0")
	}
}
