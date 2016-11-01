package cache

import (
	"math/rand"
	"server/base/service"
	"server/base/utils"
	"server/db"
	"strconv"
	"time"
)

type Color struct {
	R int32 `json:"r"`
	G int32 `json:"g"`
	B int32 `json:"b"`
}

type Player struct {
	Pid            int64   `json:"pid"`
	Name           string  `json:"name"`
	Color          Color   `json:"color"`
	LocationX      float32 `json:"locationX"`
	LocationY      float32 `json:"locationY"`
	LocationDetail string  `json:"locationDetail"`
}

var (
	didToPid    = map[string]int64{}
	pidToPlayer = map[int64]*Player{}
)

func GetPlayerByDid(did string) *Player {
	ret, _ := srv.AddTask(getPlayerByDid, service.Args{did}).Result().(*Player)
	return ret
}

func GetPlayerByPid(pid int64) *Player {
	ret, _ := srv.AddTask(getPlayerByPid, service.Args{pid}).Result().(*Player)
	return ret
}

//
// implement
//
func getPlayerByDid(args service.Args) service.Result {
	if len(args) < 1 {
		utils.ArgsNumberErr("getPlayerByDid", 1)
		return nil
	}

	did, ok := args[0].(string)
	if !ok {
		utils.ArgsTypeCastErr("getPlayerByDid", 0)
		return nil
	}

	if len(did) == 0 {
		utils.InvalidValueErr("getPlayerByDid", "len(did) == 0")
		return nil
	}

	player := getInnerPlayerByDid(did)
	if player == nil {
		playerInfo := db.LoadPlayerInfoByDid(did)
		if playerInfo == nil {
			playerInfo = &db.PlayerInfo{
				Did:  did,
				Name: "Talker" + strconv.FormatInt(time.Now().Unix(), 10),
			}
			if !db.AddPlayerInfo(playerInfo) {
				return nil
			}
		}

		player = &Player{
			Pid:   playerInfo.Pid,
			Name:  playerInfo.Name,
			Color: getRandomColor(),
		}
		pidToPlayer[player.Pid] = player
	}

	return player
}

func getPlayerByPid(args service.Args) service.Result {
	if len(args) < 1 {
		utils.ArgsNumberErr("getPlayerByPid", 1)
		return nil
	}

	pid, ok := args[0].(int64)
	if !ok {
		utils.ArgsTypeCastErr("getPlayerByPid", 0)
		return nil
	}

	if pid <= 0 {
		utils.InvalidValueErr("getPlayerByPid", "pid <= 0")
		return nil
	}

	return getInnerPlayerByPid(pid)
}

//
// inner func
//
func getRandomColor() Color {
	return Color{
		R: int32(rand.Intn(170)),
		G: int32(rand.Intn(170)),
		B: int32(rand.Intn(170)),
	}
}

func getInnerPlayerByDid(did string) *Player {
	if len(did) == 0 {
		utils.InvalidValueErr("getInnerPlayerByDid", "len(did) == 0")
		return nil
	}

	pid := didToPid[did]
	if pid <= 0 {
		return nil
	}
	return getInnerPlayerByPid(pid)
}

func getInnerPlayerByPid(pid int64) *Player {
	if pid <= 0 {
		utils.InvalidValueErr("getInnerPlayerByPid", "pid <= 0")
		return nil
	}

	return pidToPlayer[pid]
}
