package behavior

////////////////////////////////////////////
// type, const, var
//
type UserData interface{}

type Driver interface {
	UserData() UserData
	Drive(ctx *Context) UserData
}

type Drivers []Driver
