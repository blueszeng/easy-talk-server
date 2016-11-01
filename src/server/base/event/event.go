package event

import "server/base/typex"

////////////////////////////////////////////
// type, const, var
//
type Event interface {
	Name() string
	Data() *typex.Bundle
}

type BaseEvent struct {
	name string
	data *typex.Bundle
}

////////////////////////////////////////////
// func
//
func NewBaseEvent(name string) *BaseEvent {
	return &BaseEvent{
		name: name,
		data: typex.NewBundle(),
	}
}

func NewBaseEventWithData(name string, data *typex.Bundle) *BaseEvent {
	return &BaseEvent{
		name: name,
		data: data,
	}
}

//
// BaseEvent method
//
func (evt *BaseEvent) Name() string {
	return evt.name
}

func (evt *BaseEvent) SetData(data *typex.Bundle) {
	evt.data = data
}

func (evt *BaseEvent) Data() *typex.Bundle {
	return evt.data
}
