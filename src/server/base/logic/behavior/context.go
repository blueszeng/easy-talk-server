package behavior

////////////////////////////////////////////
// type, const, var
//
type Context struct {
	continuation *Continuation
	self         UserData
}

////////////////////////////////////////////
// func
//
func NewContext(self UserData) *Context {
	return &Context{
		self: self,
	}
}

func (c *Context) Continuation() *Continuation {
	return c.continuation
}

func (c *Context) SetContinuation(continuation *Continuation) {
	c.continuation = continuation
}

func (c *Context) Self() UserData {
	return c.self
}
