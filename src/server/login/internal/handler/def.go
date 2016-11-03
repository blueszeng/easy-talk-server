package handler

import (
	"net/http"
	"server/cache"
)

//
// handler
//
type Handler func(http.ResponseWriter, *http.Request)

//
// errcode
//
type Errcode int32
type Error struct {
	Errcode Errcode `json:"errcode"`
}

const (
	Success Errcode = 0
	Failed  Errcode = 1
	Offline Errcode = 2

	NameTooLong Errcode = 10
)

//
// push
//
type Action string

const (
	Login          Action = "login"
	UpdateLocation Action = "update-locaction"
	SendMsg        Action = "send-msg"
	ChangeName     Action = "change-name"
)

const (
	MaxNameLen = 12
)

type PushResponse struct {
	Errcode Errcode `json:"errcode"`
	Pid     int64   `json:"pid"`
	Name    string  `json:"name"`
}

//
// pull
//
type PullResponse struct {
	Errcode Errcode      `json:"errcode"`
	Msgs    []*cache.Msg `json:"msgs"`
}
