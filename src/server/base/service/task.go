package service

import (
	"server/base/utils"
	"time"

	"github.com/name5566/leaf/log"
)

////////////////////////////////////////////
// type, const, var
//
type Args []interface{}
type Result interface{}
type TaskHandler func(Args) Result
type ResultCh chan Result

type Task struct {
	handler  TaskHandler
	args     Args
	resultCh ResultCh
}

////////////////////////////////////////////
// func
//

//
// task method
//
func (t *Task) Call() {
	if t.handler != nil {
		t.SendResult(t.handler(t.args))
	} else {
		log.Error("Task.Call, task miss handler")
	}
}

func (t *Task) SendResult(r Result) {
	if len(t.resultCh) != 0 {
		log.Error("Task.SendResult, result channel is not empty, len: %d", len(t.resultCh))
		return
	}

	t.resultCh <- r
}

func (t *Task) Result() Result {
	select {
	case result := <-t.resultCh:
		return result
	case <-time.After(GetResultWaitTime):
		log.Error("Task.Result, time out")
		utils.Traceback(nil)
		return nil
	}
}
