package schedule

import (
	"time"
)

////////////////////////////////////////////
// type, const, var
//
const (
	MinutePerHour = 60
)

type TaskHandler func()
type Task struct {
	hour    int
	min     int
	handler TaskHandler

	ready bool
}

////////////////////////////////////////////
// func
//
func NewTask(hour int, minute int, f TaskHandler) *Task {
	return &Task{
		hour:    hour,
		min:     minute,
		handler: f,
		ready:   true,
	}
}

//
// Task method
//
func (t *Task) check() {
	time := time.Now()
	h, m, _ := time.Clock()
	if t.ready {
		if h == t.hour && m == t.min {
			t.active()
			t.ready = false
		}
	} else {
		if h == t.hour && m == nextMinute(t.min) {
			t.ready = true
		}
	}
}

func (t *Task) active() {
	if t.handler != nil {
		t.handler()
	}
}

//
// inner func
//
func nextMinute(curMin int) int {
	if curMin < 0 {
		return 0
	}

	if curMin < MinutePerHour-1 {
		return curMin + 1
	}
	return 0
}
