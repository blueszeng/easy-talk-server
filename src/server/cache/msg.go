package cache

import (
	"server/base/service"
	"server/base/utils"
	"server/db"
	"time"
)

type Msg struct {
	Mid    int64   `json:"mid"`
	Date   int64   `json:"date"`
	Msg    string  `json:"msg"`
	Player *Player `json:"player"`
}

const (
	MaxCacheMsgNum = 1000
)

var (
	msgs     []*Msg
	startMid int64
)

func init() {
	msgs = []*Msg{}
	startMid = db.GetLastMsgId()
	if startMid <= 0 {
		startMid = 1
	}
}

func GetMsgsFromTheMid(mid int64) []*Msg {
	ret, _ := srv.AddTask(getMsgsFromTheMid, service.Args{mid}).Result().([]*Msg)
	return ret
}

func AddMsgByPid(pid int64, msg string) *Msg {
	ret, _ := srv.AddTask(addMsgByPid, service.Args{pid, msg}).Result().(*Msg)
	return ret
}

//
// implement
//
func getMsgsFromTheMid(args service.Args) service.Result {
	if len(msgs) == 0 {
		return nil
	}

	if len(args) < 1 {
		utils.ArgsNumberErr("getMsgsFromTheMid", 1)
		return nil
	}

	mid, ok := args[0].(int64)
	if !ok {
		utils.ArgsTypeCastErr("getMsgsFromTheMid", 0)
		return nil
	}

	if mid < startMid {
		mid = startMid
	}

	index := int(mid - startMid)
	if index >= len(msgs) {
		return nil
	}

	return msgs[index:]
}

func addMsgByPid(args service.Args) service.Result {
	if len(args) < 2 {
		utils.ArgsNumberErr("addMsgByPid", 2)
		return nil
	}

	pid, ok := args[0].(int64)
	if !ok {
		utils.ArgsTypeCastErr("addMsgByPid", 0)
		return nil
	}

	if pid <= 0 {
		utils.InvalidValueErr("addMsgByPid", "pid <= 0")
		return nil
	}

	player := getInnerPlayerByPid(pid)
	if player == nil {
		utils.InvalidValueErr("addMsgByPid", "player == nil")
		return nil
	}

	msg, ok := args[1].(string)
	if !ok {
		utils.ArgsTypeCastErr("addMsgByPid", 1)
		return nil
	}

	if len(msg) == 0 {
		utils.InvalidValueErr("addMsgByPid", "len(msg) == 0")
		return nil
	}

	msgInfo := &db.MsgInfo{
		Pid:  pid,
		Date: time.Now().Unix(),
		Msg:  msg,
	}
	if !db.AddMsgInfo(msgInfo) {
		return nil
	}

	if len(msgs) >= MaxCacheMsgNum {
		resetInnerMsgs()
	}

	msgSt := &Msg{
		Mid:    msgInfo.Mid,
		Date:   msgInfo.Date,
		Msg:    msgInfo.Msg,
		Player: player,
	}
	msgs = append(msgs, msgSt)

	return msgSt
}

func resetInnerMsgs() {
	cnt := len(msgs)
	if cnt > 0 {
		lastMsg := msgs[cnt-1]
		startMid = lastMsg.Mid + 1
		msgs = msgs[:0]
	}
}
