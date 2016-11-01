package typex

////////////////////////////////////////////
// type, const, var
//
type Bundle struct {
	bools     map[string]bool
	int32s    map[string]int32
	float32s  map[string]float32
	strings   map[string]string
	userDatas map[string]interface{}
}

////////////////////////////////////////////
// func
//
func NewBundle() *Bundle {
	return &Bundle{
		bools:     make(map[string]bool),
		int32s:    make(map[string]int32),
		float32s:  make(map[string]float32),
		strings:   make(map[string]string),
		userDatas: make(map[string]interface{}),
	}
}

func (b *Bundle) SetBoolData(key string, value bool) {
	if b.bools != nil {
		b.bools[key] = value
	}
}

func (b *Bundle) BoolData(key string) bool {
	if b.bools != nil {
		return b.bools[key]
	}
	return false
}

func (b *Bundle) SetInt32Data(key string, value int32) {
	if b.int32s != nil {
		b.int32s[key] = value
	}
}

func (b *Bundle) Int32Data(key string) int32 {
	if b.int32s != nil {
		return b.int32s[key]
	}
	return 0
}

func (b *Bundle) SetFloat32Data(key string, value float32) {
	if b.float32s != nil {
		b.float32s[key] = value
	}
}

func (b *Bundle) Float32Data(key string) float32 {
	if b.float32s != nil {
		return b.float32s[key]
	}
	return 0
}

func (b *Bundle) SetStringData(key string, value string) {
	if b.strings != nil {
		b.strings[key] = value
	}
}

func (b *Bundle) StringData(key string) string {
	if b.strings != nil {
		return b.strings[key]
	}
	return ""
}

func (b *Bundle) SetUserData(key string, value interface{}) {
	if b.userDatas != nil {
		b.userDatas[key] = value
	}
}

func (b *Bundle) UserData(key string) interface{} {
	if b.userDatas != nil {
		return b.userDatas[key]
	}
	return nil
}
