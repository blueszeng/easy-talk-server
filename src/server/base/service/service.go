package service

import (
	"server/base/conf"
	"server/base/utils"
	"time"

	"github.com/name5566/leaf/log"
)

////////////////////////////////////////////
// type, const, var
//
type EnterHandler func()
type UpdateHanlder func(float32)

type Service struct {
	isRunning     bool
	hasStoped     bool
	tasksCh       chan *Task
	stopCh        chan bool
	enterHandler  EnterHandler
	updateHandler UpdateHanlder
}

const (
	TaskResultChLen   = 1
	GetResultWaitTime = 3 * time.Second
)

const (
	UpdateIntervalInSecond = 0.1
	UpdateInterval         = time.Duration(UpdateIntervalInSecond * float64(time.Second))
)

////////////////////////////////////////////
// func
//
func NewService0() *Service {
	return NewService2(nil, nil)
}

func NewService1(updateHandler UpdateHanlder) *Service {
	return NewService2(nil, updateHandler)
}

func NewService2(enterHandler EnterHandler, updateHandler UpdateHanlder) *Service {
	return &Service{
		isRunning:     false,
		hasStoped:     false,
		tasksCh:       make(chan *Task, conf.ServiceTaskChLen),
		stopCh:        make(chan bool, 1),
		enterHandler:  enterHandler,
		updateHandler: updateHandler,
	}
}

//
// service method
//
func (s *Service) Start() {
	if s.hasStoped {
		log.Error("Service has stoped, can't start")
		return
	}

	if s.isRunning {
		log.Error("Service is running, can't start again")
		return
	}
	s.isRunning = true
	go s.run()
}

func (s *Service) AddTask(handler TaskHandler, args Args) *Task {
	if s.hasStoped {
		log.Error("Service has stoped, can't add task")
		task := &Task{resultCh: make(ResultCh, TaskResultChLen)}
		task.SendResult(nil)
		return task
	}

	task := &Task{
		handler:  handler,
		args:     args,
		resultCh: make(ResultCh, TaskResultChLen),
	}
	s.tasksCh <- task
	return task
}

func (s *Service) Stop() {
	if s.hasStoped {
		log.Error("Service has stoped, can't stop again")
		return
	}
	s.hasStoped = true
	s.stopCh <- true
}

func (s *Service) run() {
	ticker := time.NewTicker(UpdateInterval)
	defer ticker.Stop()

	if s.enterHandler != nil {
		s.enterHandler()
	}

	for {
		select {
		case task := <-s.tasksCh:
			s.handleTask(task)
		case <-ticker.C:
			s.update()
		case <-s.stopCh:
			s.destroy()
			return
		}
	}
	return
}

func (s *Service) handleTask(t *Task) {
	if t == nil {
		return
	}

	defer utils.Recover()

	t.Call()
}

func (s *Service) update() {
	if s.updateHandler == nil {
		return
	}

	defer utils.Recover()

	s.updateHandler(UpdateIntervalInSecond)
}

func (s *Service) destroy() {
	s.hasStoped = true
	close(s.stopCh)
	close(s.tasksCh)
	s.enterHandler = nil
	s.updateHandler = nil
	s.isRunning = false
}
