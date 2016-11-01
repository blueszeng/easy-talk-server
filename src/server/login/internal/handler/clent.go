package handler

import (
	"net/http"
	"server/cache"
)

type Error struct {
	Errcode int32 `json:"errcode"`
}

type PushResponse struct {
	Errcode int32 `json:"errcode"`
	Pid     int64 `json:"pid"`
}

type PullResponse struct {
	Errcode int32        `json:"errcode"`
	Msgs    []*cache.Msg `json:"msgs"`
}

// errcode
const (
	Success int32 = 0
	Failed  int32 = 1
)

////////////////////////////////////////////
// func
//
func HandlePush(w http.ResponseWriter, req *http.Request) {
	if !checkMethodAndParseForm(req, "HandlePush") {
		return
	}
}

func HandlePull(w http.ResponseWriter, req *http.Request) {
	if !checkMethodAndParseForm(req, "HandlePull") {
		return
	}

	pid, ok := parseInt64(req, "pid")
	if !ok {
		sayErr(w, Failed)
		return
	}

	if pid <= 0 {
		sayErr(w, Failed)
		return
	}

	player := cache.GetPlayerByPid(pid)
	if player == nil {
		sayErr(w, Failed)
		return
	}

	mid, ok := parseInt64(req, "mid")
	if !ok {
		sayErr(w, Failed)
		return
	}

	if mid <= 0 {
		sayErr(w, Failed)
		return
	}

	msgs := cache.GetMsgsFromTheMid(mid)
	if len(msgs) == 0 {
		sayErr(w, Success)
		return
	}

	resp := &PullResponse{
		Errcode: Success,
		Msgs:    msgs,
	}
	w.Write(jsonMarshal(resp))
}
