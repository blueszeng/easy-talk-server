package handler

import (
	"net/http"
	"server/base/utils"
	"server/cache"
)

//
// handle push
//
var actionHandler = map[Action]Handler{
	Login:          handleLogin,
	UpdateLocation: handleUpdateLocation,
	SendMsg:        handleSendMsg,
	ChangeName:     handleChangeName,
}

func HandlePush(w http.ResponseWriter, req *http.Request) {
	if !checkMethodAndParseForm(req, "HandlePush") {
		return
	}

	action, ok := parseString(req, "action")
	if !ok {
		sayErr(w, Failed)
		return
	}

	if len(action) == 0 {
		sayErr(w, Failed)
		return
	}

	handler := actionHandler[Action(action)]
	if handler == nil {
		utils.LogErr("HandlePush", "invalid action: %s", action)
		sayErr(w, Failed)
		return
	}

	handler(w, req)
}

func handleLogin(w http.ResponseWriter, req *http.Request) {
	did, ok := parseString(req, "did")
	if !ok {
		sayErr(w, Failed)
		return
	}

	if len(did) == 0 {
		sayErr(w, Failed)
		return
	}

	player := cache.GetPlayerByDid(did)
	if player == nil {
		utils.InvalidValueErr("handleLogin", "player == nil")
		sayErr(w, Failed)
		return
	}

	resp := &PushResponse{
		Errcode: Success,
		Pid:     player.Pid,
		Name:    player.Name,
	}
	w.Write(jsonMarshal(resp))
}

func handleUpdateLocation(w http.ResponseWriter, req *http.Request) {
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
		sayErr(w, Offline)
		return
	}

	locationX, ok := parseFloat32(req, "locationX")
	if !ok {
		sayErr(w, Failed)
		return
	}

	locationY, ok := parseFloat32(req, "locationY")
	if !ok {
		sayErr(w, Failed)
		return
	}

	locationZ, ok := parseFloat32(req, "locationZ")
	if !ok {
		sayErr(w, Failed)
		return
	}

	locationDetail, ok := parseString(req, "locationDetail")
	if !ok {
		sayErr(w, Failed)
		return
	}

	player = cache.UpdatePlayerLocationByPid(pid, locationX, locationY, locationZ, locationDetail)
	if player == nil {
		sayErr(w, Offline)
		return
	}

	sayErr(w, Success)
}

func handleSendMsg(w http.ResponseWriter, req *http.Request) {
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
		sayErr(w, Offline)
		return
	}

	channel, ok := parseInt32(req, "channel")
	if !ok {
		sayErr(w, Failed)
		return
	}

	if !cache.ValidChannel(channel) {
		sayErr(w, Failed)
		return
	}

	msg, ok := parseString(req, "msg")
	if !ok {
		sayErr(w, Failed)
		return
	}

	if len(msg) == 0 {
		sayErr(w, Failed)
		return
	}

	if cache.AddMsgByPid(pid, channel, msg) == nil {
		sayErr(w, Failed)
		return
	}

	sayErr(w, Success)
}

func handleChangeName(w http.ResponseWriter, req *http.Request) {
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
		sayErr(w, Offline)
		return
	}

	name, ok := parseString(req, "name")
	if !ok {
		sayErr(w, Failed)
		return
	}

	if len(name) == 0 {
		sayErr(w, Failed)
		return
	}

	nameArr := []rune(name)
	if len(nameArr) > MaxNameLen {
		sayErr(w, NameTooLong)
		return
	}

	player = cache.ChangePlayerNameByPid(pid, name)
	if player == nil {
		sayErr(w, Failed)
		return
	}

	resp := &PushResponse{
		Errcode: Success,
		Name:    player.Name,
	}
	w.Write(jsonMarshal(resp))
}

//
// handle pull
//
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
		sayErr(w, Offline)
		return
	}

	channel, ok := parseInt32(req, "channel")
	if !ok {
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

	msgs := cache.GetMsgsFromTheMid(channel, mid)
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
