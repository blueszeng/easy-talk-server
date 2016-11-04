package cache

// channel
const (
	MinChannel = 0

	Nearby = 0

	MaxChannel = Nearby
)

func ValidChannel(ch int32) bool {
	return (MinChannel <= ch) && (ch <= MaxChannel)
}
