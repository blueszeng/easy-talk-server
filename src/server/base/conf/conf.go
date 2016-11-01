package conf

import (
	"time"
)

var (
	// gate conf
	PendingWriteNum        = 2000
	MaxMsgLen       uint32 = 8192
	HTTPTimeout            = 10 * time.Second
	LenMsgLen              = 4
	LittleEndian           = false

	// skeleton conf
	GoLen              = 10000
	TimerDispatcherLen = 10000
	ChanRPCLen         = 10000

	// service conf
	ServiceTaskChLen = 10000
)
