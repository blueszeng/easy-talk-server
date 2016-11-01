package utils

import (
	"fmt"

	"github.com/name5566/leaf/log"
)

////////////////////////////////////////////
// func
//
func LogErr(tag string, format string, a ...interface{}) {
	log.Error("%s, %s", tag, fmt.Sprintf(format, a...))
}

func ArgsNumberErr(tag string, needNum int) {
	LogErr(tag, "args wrong number, need: %d", needNum)
}

func ArgsTypeCastErr(tag string, index int) {
	LogErr(tag, "type cast failed, args index: %d", index)
}

func TypeCastErr(tag string, targetType string) {
	LogErr(tag, "type cast failed, target type: %s", targetType)
}

func InvalidValueErr(tag string, detail string) {
	LogErr(tag, "invalid value: %s", detail)
}
