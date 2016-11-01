package behavior

////////////////////////////////////////////
// type, const, var
//
type Continuation struct {
	subContinuation *Continuation
	curStep         int
	userData        UserData
}

////////////////////////////////////////////
// func
//
func NewContinuation() *Continuation {
	return &Continuation{}
}

func (c *Continuation) SubContinuation() *Continuation {
	return c.subContinuation
}

func (c *Continuation) SetSubContinuation(sub *Continuation) {
	c.subContinuation = sub
}

func (c *Continuation) CurStep() int {
	return c.curStep
}

func (c *Continuation) SetCurStep(step int) {
	c.curStep = step
}

func (c *Continuation) UserData() UserData {
	return c.userData
}

func (c *Continuation) SetUserData(data UserData) {
	c.userData = data
}
