package utils

////////////////////////////////////////////
// type, const, var
//
type Action interface {
	Update(dt float32)
	IsFinished() bool
}

type ActionManager struct {
	actions    []Action
	isUpdating bool
	tobeAdded  []Action
}

type DelayRunAction struct {
	delay    float32
	curTime  float32
	callback func()
}

////////////////////////////////////////////
// func
//
func NewDelayRunAction(dt float32, callback func()) *DelayRunAction {
	return &DelayRunAction{
		delay:    dt,
		curTime:  0,
		callback: callback,
	}
}

func NewActionManager() *ActionManager {
	return &ActionManager{
		actions:   []Action{},
		tobeAdded: []Action{},
	}
}

//
// DelayRunAction method
//
func (a *DelayRunAction) Update(dt float32) {
	if a.IsFinished() {
		return
	}

	a.curTime += dt
	if a.IsFinished() {
		if a.callback != nil {
			a.callback()
		}
	}
}

func (a *DelayRunAction) IsFinished() bool {
	return a.curTime >= a.delay
}

//
// ActionManager method
//
func (m *ActionManager) Add(a Action) {
	if a == nil {
		return
	}

	if m.actions == nil {
		m.actions = []Action{}
	}

	if m.tobeAdded == nil {
		m.tobeAdded = []Action{}
	}

	if m.isUpdating {
		m.tobeAdded = append(m.tobeAdded, a)
	} else {
		m.actions = append(m.actions, a)
	}
}

func (m *ActionManager) Update(dt float32) {
	if m.actions == nil || len(m.actions) == 0 {
		return
	}

	m.isUpdating = true

	actions := []Action{}
	for _, a := range m.actions {
		if a == nil {
			continue
		}

		a.Update(dt)
		if a.IsFinished() {
			continue
		}
		actions = append(actions, a)
	}
	actions = append(actions, m.tobeAdded...)
	m.actions = actions
	m.tobeAdded = []Action{}

	m.isUpdating = false
}

func (m *ActionManager) StopAll() {
	if m.actions == nil || len(m.actions) == 0 {
		return
	}

	m.actions = []Action{}
}
