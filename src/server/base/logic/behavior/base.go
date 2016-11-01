package behavior

////////////////////////////////////////////
// Result
//
type Result int

const (
	Success  Result = 0
	Continue Result = 1
	Failure  Result = 2
)

////////////////////////////////////////////
// Box
//
type Box struct {
	useData UserData
}

func NewBox() *Box {
	return &Box{}
}

func (b *Box) UserData() UserData {
	return b.useData
}

func (b *Box) SetUserData(data UserData) UserData {
	old := b.useData
	b.useData = data
	return old
}
