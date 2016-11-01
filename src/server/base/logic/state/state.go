package state

import (
	"server/base/event"
)

////////////////////////////////////////////
// type, const, var
//
type State interface {
	Name() string
	IsReplaceable() bool
	OnEnter()
	OnUpdate(dt float32)
	OnEvent(evt event.Event)
	OnExit()
}
