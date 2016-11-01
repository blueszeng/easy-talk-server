package state

import (
	"server/base/event"
	"server/base/utils"
)

////////////////////////////////////////////
// type, const, var
//
type StateMachine struct {
	globalStates map[string]State
	curState     State
}

////////////////////////////////////////////
// func
//
func NewStateMachine() *StateMachine {
	return &StateMachine{
		globalStates: make(map[string]State),
	}
}

//
// StateMachine method
//
func (sm *StateMachine) AddGlobal(s State) {
	if s == nil {
		utils.InvalidValueErr("StateMachine.AddGlobal", "s == nil")
		return
	}

	name := s.Name()
	if len(name) == 0 {
		utils.InvalidValueErr("StateMachine.AddGlobal", "len(name) == 0")
		return
	}

	if sm.globalStates[name] != nil {
		utils.InvalidValueErr("StateMachine.AddGlobal", "state already exist")
		return
	}

	sm.globalStates[name] = s
	s.OnEnter()
}

func (sm *StateMachine) Change(s State) {
	if s == nil {
		utils.InvalidValueErr("StateMachine.Change", "s == nil")
		return
	}

	if sm.curState != nil {
		if !sm.curState.IsReplaceable() {
			return
		}

		sm.curState.OnExit()
	}

	sm.curState = s
	s.OnEnter()
}

func (sm *StateMachine) Update(dt float32) {
	for _, s := range sm.globalStates {
		s.OnUpdate(dt)
	}

	if sm.curState != nil {
		sm.curState.OnUpdate(dt)
	}
}

func (sm *StateMachine) Event(evt event.Event) {
	for _, s := range sm.globalStates {
		s.OnEvent(evt)
	}

	if sm.curState != nil {
		sm.curState.OnEvent(evt)
	}
}
