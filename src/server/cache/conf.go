package cache

//
// checker
//
const (
	CheckOnlineNumInterval   = 60 * 5
	CheckPlayerAliveInterval = 60
)

const (
	MaxPlayerAliveInterval = 60 * 3
)

//
// channel
//
const (
	MinChannel = 0

	Nearby = 0

	MaxChannel = Nearby
)

func ValidChannel(ch int32) bool {
	return (MinChannel <= ch) && (ch <= MaxChannel)
}

//
// msg
//
const (
	MaxCacheMsgNum = 500
)
