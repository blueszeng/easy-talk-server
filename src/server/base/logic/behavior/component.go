package behavior

////////////////////////////////////////////
// SequenceDriver
//
type SequenceDriver struct {
	subDrivers Drivers
}

func NewSequenceDriver(subs Drivers) *SequenceDriver {
	return &SequenceDriver{
		subDrivers: subs,
	}
}

func (d *SequenceDriver) UserData() UserData {
	return d
}

func (d *SequenceDriver) Drive(ctx *Context) UserData {
	if d.subDrivers == nil || len(d.subDrivers) == 0 {
		return Failure
	}

	continuation := ctx.Continuation()
	if continuation == nil {
		continuation = NewContinuation()
		continuation.SetCurStep(0)
	}

	ctx.SetContinuation(continuation.SubContinuation())
	subDriver := d.subDrivers[continuation.CurStep()]
	result := subDriver.Drive(ctx)
	if result, ok := result.(Result); ok {
		if result == Failure {
			ctx.SetContinuation(nil)
			return Failure
		}

		if result == Continue {
			continuation.SetSubContinuation(ctx.Continuation())
			ctx.SetContinuation(continuation)
			return Continue
		}

		curStep := continuation.CurStep() + 1
		if curStep >= len(d.subDrivers) {
			ctx.SetContinuation(nil)
			return Success
		}

		continuation.SetCurStep(curStep)
		continuation.SetSubContinuation(nil)
		ctx.SetContinuation(continuation)
		return Continue
	} else {
		ctx.SetContinuation(nil)
		return Failure
	}
}

////////////////////////////////////////////
// SequenceLoopDriver
//
type SequenceLoopDriver struct {
	subDrivers Drivers
}

func NewSequenceLoopDriver(subs Drivers) *SequenceLoopDriver {
	return &SequenceLoopDriver{
		subDrivers: subs,
	}
}

func (d *SequenceLoopDriver) UserData() UserData {
	return d
}

func (d *SequenceLoopDriver) Drive(ctx *Context) UserData {
	if d.subDrivers == nil || len(d.subDrivers) == 0 {
		return Failure
	}

	continuation := ctx.Continuation()
	if continuation == nil {
		continuation = NewContinuation()
		continuation.SetCurStep(0)
	}

	ctx.SetContinuation(continuation.SubContinuation())
	subDriver := d.subDrivers[continuation.CurStep()]
	result := subDriver.Drive(ctx)
	if result, ok := result.(Result); ok {
		if result == Failure {
			ctx.SetContinuation(nil)
			return Failure
		}

		if result == Continue {
			continuation.SetSubContinuation(ctx.Continuation())
			ctx.SetContinuation(continuation)
			return Continue
		}

		curStep := continuation.CurStep() + 1
		if curStep >= len(d.subDrivers) {
			curStep = 0
		}

		continuation.SetCurStep(curStep)
		continuation.SetSubContinuation(nil)
		ctx.SetContinuation(continuation)
		return Continue
	} else {
		ctx.SetContinuation(nil)
		return Failure
	}
}

////////////////////////////////////////////
// SelectDriver
//
type SelectDriver struct {
	subDrivers Drivers
}

func NewSelectDriver(subs Drivers) *SelectDriver {
	return &SelectDriver{
		subDrivers: subs,
	}
}

func (d *SelectDriver) UserData() UserData {
	return d
}

func (d *SelectDriver) Drive(ctx *Context) UserData {
	if d.subDrivers == nil || len(d.subDrivers) == 0 {
		return Failure
	}

	continuation := ctx.Continuation()
	if continuation == nil {
		continuation = NewContinuation()
		continuation.SetCurStep(0)
	}

	ctx.SetContinuation(continuation.SubContinuation())
	subDriver := d.subDrivers[continuation.CurStep()]
	result := subDriver.Drive(ctx)
	if result, ok := result.(Result); ok {
		if result == Success {
			ctx.SetContinuation(nil)
			return Success
		}

		if result == Continue {
			continuation.SetSubContinuation(ctx.Continuation())
			ctx.SetContinuation(continuation)
			return Continue
		}

		curStep := continuation.CurStep() + 1
		if curStep >= len(d.subDrivers) {
			ctx.SetContinuation(nil)
			return Failure
		}

		continuation.SetCurStep(curStep)
		continuation.SetSubContinuation(nil)
		ctx.SetContinuation(continuation)
		return Continue
	} else {
		ctx.SetContinuation(nil)
		return Failure
	}
}

////////////////////////////////////////////
// NegateDriver
//
type NegateDriver struct {
	subDriver Driver
}

func NewNegateDriver(sub Driver) *NegateDriver {
	return &NegateDriver{
		subDriver: sub,
	}
}

func (d *NegateDriver) UserData() UserData {
	return d
}

func (d *NegateDriver) Drive(ctx *Context) UserData {
	if d.subDriver == nil {
		return Success
	}

	result := d.subDriver.Drive(ctx)
	if result, ok := result.(Result); ok {
		if result == Success {
			return Failure
		}

		if result == Failure {
			return Success
		}

		return Continue
	}

	return Success
}

////////////////////////////////////////////
// CheckDriver
//
type CheckDriver struct {
	getter Driver
}

func NewCheckDriver(getter Driver) *CheckDriver {
	return &CheckDriver{
		getter: getter,
	}
}

func (d *CheckDriver) UserData() UserData {
	return d
}

func (d *CheckDriver) Drive(ctx *Context) UserData {
	if d.getter == nil {
		return Failure
	}

	continuation := ctx.continuation
	ctx.SetContinuation(nil)

	val := d.getter.Drive(ctx)
	ctx.SetContinuation(continuation)

	if val, ok := val.(bool); ok {
		if val {
			return Success
		}
	}
	return Failure
}

////////////////////////////////////////////
// WithDriver
//
type WithDriver struct {
	box       *Box
	getter    Driver
	subDriver Driver
}

func NewWithDriver(box *Box, getter Driver, sub Driver) *WithDriver {
	return &WithDriver{
		box:       box,
		getter:    getter,
		subDriver: sub,
	}
}

func (d *WithDriver) UserData() UserData {
	return d
}

func (d *WithDriver) Drive(ctx *Context) UserData {
	if d.box == nil || d.getter == nil || d.subDriver == nil {
		return Failure
	}

	continuation := ctx.Continuation()

	ctx.SetContinuation(nil)
	val := d.getter.Drive(ctx)
	oldVal := d.box.SetUserData(val)

	if continuation != nil {
		ctx.SetContinuation(continuation.SubContinuation())
	}
	result := d.subDriver.Drive(ctx)
	d.box.SetUserData(oldVal)

	if ctx.Continuation() != nil {
		if continuation == nil {
			continuation = NewContinuation()
		}
		continuation.SetSubContinuation(ctx.Continuation())
		ctx.SetContinuation(continuation)
	}

	return result
}

////////////////////////////////////////////
// IfElseDriver
//
const (
	UnTested  = 0
	TestTrue  = 1
	TestFalse = 2
)

type IfElseDriver struct {
	tester   Driver
	leftSub  Driver
	rightSub Driver
}

func NewIfElseDriver(tester Driver, leftSub Driver, rightSub Driver) *IfElseDriver {
	return &IfElseDriver{
		tester:   tester,
		leftSub:  leftSub,
		rightSub: rightSub,
	}
}

func (d *IfElseDriver) UserData() UserData {
	return d
}

func (d *IfElseDriver) Drive(ctx *Context) UserData {
	if d.tester == nil || d.leftSub == nil {
		return Failure
	}

	continuation := ctx.Continuation()

	needTest := true
	is := false
	if continuation != nil {
		// curStep: 0 (needTest), 1 (Test==true), 2 (Test==false)
		if continuation.CurStep() != UnTested {
			needTest = false
			if continuation.CurStep() == TestTrue {
				is = true
			}
		}
	}

	if needTest {
		ctx.SetContinuation(nil)
		val := d.tester.Drive(ctx)
		if val, ok := val.(bool); ok {
			is = val
		} else {
			return Failure
		}
	}

	if continuation != nil {
		ctx.SetContinuation(continuation.SubContinuation())
	}

	var result UserData
	if is {
		result = d.leftSub.Drive(ctx)
	} else {
		if d.rightSub == nil {
			ctx.SetContinuation(nil)
			return Failure
		}
		result = d.rightSub.Drive(ctx)
	}

	if result, ok := result.(Result); ok {
		if result == Continue {
			if ctx.Continuation() != nil {
				if continuation == nil {
					continuation = NewContinuation()
				}
				continuation.SetSubContinuation(ctx.Continuation())
				if continuation.CurStep() == UnTested {
					if is {
						continuation.SetCurStep(TestTrue)
					} else {
						continuation.SetCurStep(TestFalse)
					}
				}
				ctx.SetContinuation(continuation)
			}
		}
		return result
	} else {
		ctx.SetContinuation(nil)
		return Failure
	}
}

////////////////////////////////////////////
// WaitDriver
//
type WaitDriver struct {
	subDriver Driver
	maxTimes  int
}

func NewWaitDriver(sub Driver, maxTimes int) *WaitDriver {
	return &WaitDriver{
		subDriver: sub,
		maxTimes:  maxTimes,
	}
}

func (d *WaitDriver) UserData() UserData {
	return d
}

func (d *WaitDriver) Drive(ctx *Context) UserData {
	if d.subDriver == nil {
		return Failure
	}

	continuation := ctx.Continuation()
	if continuation != nil {
		ctx.SetContinuation(continuation.subContinuation)
	}

	result := d.subDriver.Drive(ctx)
	if result, ok := result.(Result); ok {
		if result == Success {
			ctx.SetContinuation(nil)
			return Success
		}

		if d.maxTimes > 0 {
			if continuation == nil {
				continuation = NewContinuation()
			}

			curTimes := continuation.CurStep() + 1
			if curTimes >= d.maxTimes {
				ctx.SetContinuation(nil)
				return Failure
			}

			continuation.SetCurStep(curTimes)
		}

		if ctx.Continuation() != nil {
			if continuation == nil {
				continuation = NewContinuation()
			}
			continuation.SetSubContinuation(ctx.Continuation())
		}
		ctx.SetContinuation(continuation)

		return Continue
	}

	ctx.SetContinuation(nil)
	return Failure
}
