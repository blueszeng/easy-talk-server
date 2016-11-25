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
	ImgDir    = "img"
	ImgPrefix = "/msg-img/"
)

// msg type
const (
	TextMsg  = "text"
	ImageMsg = "image"
)

func ValidMsgType(mtype string) bool {
	return mtype == TextMsg ||
		mtype == ImageMsg
}

const (
	MaxCacheMsgNum = 500
)
