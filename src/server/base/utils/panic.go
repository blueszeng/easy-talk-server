package utils

import (
	"github.com/name5566/leaf/conf"
	"github.com/name5566/leaf/log"
)

func Recover() {
	if r := recover(); r != nil {
		if conf.LenStackBuf > 0 {
			Traceback(func(buff []byte) {
				log.Error("%v, %s", r, buff)
			})
		} else {
			log.Error("%v", r)
		}
	}
}
