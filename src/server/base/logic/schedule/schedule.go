package schedule

import (
	"server/base/service"
	"server/base/utils"
)

////////////////////////////////////////////
// type, const, var
//
const (
	CheckInterval = 30
)

var (
	srv              = service.NewService1(update)
	curTime  float32 = 0
	taskList         = []*Task{}
)

////////////////////////////////////////////
// func
//
func init() {
	srv.Start()
}

func AddTask(task *Task) bool {
	ret, _ := srv.AddTask(addTask, service.Args{task}).Result().(bool)
	return ret
}

//
// implement
//
func addTask(args service.Args) service.Result {
	if len(args) < 1 {
		utils.ArgsNumberErr("schedule.addTask", 1)
		return false
	}

	task, _ := args[0].(*Task)
	if task == nil {
		utils.InvalidValueErr("schedule.addTask", "task == nil")
		return false
	}

	taskList = append(taskList, task)
	return true
}

//
// inner func
//
func update(dt float32) {
	curTime += dt
	if curTime > CheckInterval {
		curTime = 0
		check()
	}
}

func check() {
	for _, t := range taskList {
		t.check()
	}
}
