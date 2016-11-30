package cache

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"image"
	"image/png"
	"os"
	"server/base/conf"
	"server/base/service"
	"server/base/utils"
	"server/db"
	"strconv"
	"time"
)

type Msg struct {
	Mid     int64   `json:"mid"`
	Date    int64   `json:"date"`
	Channel int32   `json:"channel"`
	MsgType string  `json:"msgType"`
	Msg     string  `json:"msg"`
	Player  *Player `json:"player"`
}

var (
	msgs     []*Msg
	startMid int64
)

func init() {
	msgs = make([]*Msg, 0, MaxCacheMsgNum)
	last := db.GetLastMsgId()
	start := last - MinCacheMsgNum + 1
	if start < 1 {
		start = 1
	}

	cachedPlayerInfos := map[int64]*db.PlayerInfo{}
	for id := start; id <= last; id++ {
		msgInfo := db.LoadMsgInfoByMid(id)
		if msgInfo != nil {
			playerInfo := cachedPlayerInfos[msgInfo.Pid]
			if playerInfo == nil {
				playerInfo = db.LoadPlayerInfoByPid(msgInfo.Pid)
				if playerInfo != nil {
					cachedPlayerInfos[playerInfo.Pid] = playerInfo
				}
			}

			if playerInfo != nil {
				msgs = append(msgs, &Msg{
					Mid:     msgInfo.Mid,
					Date:    msgInfo.Date,
					Channel: msgInfo.Channel,
					MsgType: msgInfo.MsgType,
					Msg:     msgInfo.Msg,
					Player: &Player{
						Pid:            playerInfo.Pid,
						Name:           playerInfo.Name,
						Color:          getRandomColor(),
						LocationX:      msgInfo.LocationX,
						LocationY:      msgInfo.LocationY,
						LocationZ:      msgInfo.LocationZ,
						LocationDetail: msgInfo.LocationDetail,
					},
				})
			}
		}
	}

	if len(msgs) > 0 {
		startMid = msgs[0].Mid
	}
}

func GetMsgsFromTheMid(channel int32, mid int64) []*Msg {
	ret, _ := srv.AddTask(getMsgsFromTheMid, service.Args{channel, mid}).Result().([]*Msg)
	return ret
}

func AddMsgByPid(pid int64, channel int32, msgType string, msg string) *Msg {
	ret, _ := srv.AddTask(addMsgByPid, service.Args{pid, channel, msgType, msg}).Result().(*Msg)
	return ret
}

//
// implement
//
func getMsgsFromTheMid(args service.Args) service.Result {
	if len(msgs) == 0 {
		return nil
	}

	if len(args) < 2 {
		utils.ArgsNumberErr("getMsgsFromTheMid", 2)
		return nil
	}

	channel, ok := args[0].(int32)
	if !ok {
		utils.ArgsTypeCastErr("getMsgsFromTheMid", 0)
		return nil
	}

	if !ValidChannel(channel) {
		utils.InvalidValueErr("getMsgsFromTheMid", "invalid channel")
		return nil
	}

	mid, ok := args[1].(int64)
	if !ok {
		utils.ArgsTypeCastErr("getMsgsFromTheMid", 1)
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
	if len(args) < 4 {
		utils.ArgsNumberErr("addMsgByPid", 4)
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

	channel, ok := args[1].(int32)
	if !ok {
		utils.ArgsTypeCastErr("addMsgByPid", 1)
		return nil
	}

	if !ValidChannel(channel) {
		utils.InvalidValueErr("addMsgByPid", "invalid channel")
		return nil
	}

	msgType, ok := args[2].(string)
	if !ok {
		utils.ArgsTypeCastErr("addMsgByPid", 2)
		return nil
	}

	if !ValidMsgType(msgType) {
		utils.InvalidValueErr("addMsgByPid", "invalid msg type")
	}

	msg, ok := args[3].(string)
	if !ok {
		utils.ArgsTypeCastErr("addMsgByPid", 3)
		return nil
	}

	if len(msg) == 0 {
		utils.InvalidValueErr("addMsgByPid", "len(msg) == 0")
		return nil
	}

	switch msgType {
	case TextMsg:
		return addInnerTextMsg(player, channel, msg)
	case ImageMsg:
		return addInnerImageMsg(player, channel, msg)
	default:
		utils.InvalidValueErr("addMsgByPid", "invalid msg type: "+msgType)
		return nil
	}
}

//
// inner func
//
func addInnerTextMsg(player *Player, channel int32, msg string) *Msg {
	msgInfo := &db.MsgInfo{
		Pid:            player.Pid,
		Date:           time.Now().Unix(),
		LocationX:      player.LocationX,
		LocationY:      player.LocationY,
		LocationZ:      player.LocationZ,
		LocationDetail: player.LocationDetail,
		Channel:        channel,
		MsgType:        TextMsg,
		Msg:            msg,
	}
	return addInnerMsgByMsgInfo(msgInfo, player)
}

func addInnerImageMsg(player *Player, channel int32, msg string) *Msg {
	buff, err := base64.StdEncoding.DecodeString(msg)
	if err != nil {
		utils.InvalidValueErr("addInnerImageMsg", "base64 decode failed: "+err.Error())
		return nil
	}

	buffer := bytes.NewBuffer(buff)
	gzipReader, err := gzip.NewReader(buffer)
	if err != nil {
		utils.InvalidValueErr("addInnerImageMsg", "new gzip reader failed: "+err.Error())
		return nil
	}
	defer gzipReader.Close()

	img, _, err := image.Decode(gzipReader)
	if err != nil {
		utils.InvalidValueErr("addInnerImageMsg", "image decode failed: "+err.Error())
		return nil
	}

	imgFileName := strconv.Itoa(int(player.Pid)) + "_" + utils.GetUUID() + ".png"
	imgFile, err := os.Create("./" + ImgDir + "/" + imgFileName)
	if err != nil {
		utils.InvalidValueErr("addInnerImageMsg", "create file failed: "+err.Error())
		return nil
	}
	defer imgFile.Close()

	png.Encode(imgFile, img)

	imgUrl := "http://" + conf.Server.TCPAddr + ":" + conf.Server.TCPPort + ImgPrefix + imgFileName
	msgInfo := &db.MsgInfo{
		Pid:            player.Pid,
		Date:           time.Now().Unix(),
		LocationX:      player.LocationX,
		LocationY:      player.LocationY,
		LocationZ:      player.LocationZ,
		LocationDetail: player.LocationDetail,
		Channel:        channel,
		MsgType:        ImageMsg,
		Msg:            imgUrl,
	}
	return addInnerMsgByMsgInfo(msgInfo, player)
}

func addInnerMsgByMsgInfo(msgInfo *db.MsgInfo, player *Player) *Msg {
	if msgInfo == nil || player == nil {
		utils.InvalidValueErr("addInnerMsgByMsgInfo", "msgInfo == nil || player == nil")
		return nil
	}

	if !db.AddMsgInfo(msgInfo) {
		return nil
	}

	playerCopy := *player
	msgSt := &Msg{
		Mid:     msgInfo.Mid,
		Date:    msgInfo.Date,
		Channel: msgInfo.Channel,
		MsgType: msgInfo.MsgType,
		Msg:     msgInfo.Msg,
		Player:  &playerCopy,
	}
	msgs = append(msgs, msgSt)
	updateInnerPlayerLastAliveTime(player.Pid)

	if len(msgs) >= MaxCacheMsgNum {
		resetInnerMsgs()
	}

	return msgSt
}

func resetInnerMsgs() {
	cnt := len(msgs)
	if cnt >= MaxCacheMsgNum {
		newStartIndex := cnt - MinCacheMsgNum
		for i := 0; i < MinCacheMsgNum; i++ {
			msgs[i] = msgs[newStartIndex+i]
		}
		msgs = msgs[:MinCacheMsgNum]
		startMid = msgs[0].Mid
	}
}
