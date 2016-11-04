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
	LocationZ      float32 `json:"locationZ"`
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

func UpdatePlayerLocationByPid(pid int64, locationX float32, locationY float32, locationZ float32, locationDetail string) *Player {
	ret, _ := srv.AddTask(updatePlayerLocationByPid, service.Args{pid, locationX, locationY, locationZ, locationDetail}).Result().(*Player)
	return ret
}

func ChangePlayerNameByPid(pid int64, name string) *Player {
	ret, _ := srv.AddTask(changePlayerNameByPid, service.Args{pid, name}).Result().(*Player)
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
				Name: "ET" + strconv.FormatInt(time.Now().Unix(), 10),
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

func updatePlayerLocationByPid(args service.Args) service.Result {
	if len(args) < 5 {
		utils.ArgsNumberErr("updatePlayerLocationByPid", 5)
		return nil
	}

	pid, ok := args[0].(int64)
	if !ok {
		utils.ArgsTypeCastErr("updatePlayerLocationByPid", 0)
		return nil
	}

	if pid <= 0 {
		utils.InvalidValueErr("updatePlayerLocationByPid", "pid <= 0")
		return nil
	}

	locationX, ok := args[1].(float32)
	if !ok {
		utils.ArgsTypeCastErr("updatePlayerLocationByPid", 1)
		return nil
	}

	locationY, ok := args[2].(float32)
	if !ok {
		utils.ArgsTypeCastErr("updatePlayerLocationByPid", 2)
		return nil
	}

	locationZ, ok := args[3].(float32)
	if !ok {
		utils.ArgsTypeCastErr("updatePlayerLocationByPid", 3)
		return nil
	}

	locationDetail, ok := args[4].(string)
	if !ok {
		utils.ArgsTypeCastErr("updatePlayerLocationByPid", 4)
		return nil
	}

	player := getInnerPlayerByPid(pid)
	if player == nil {
		utils.InvalidValueErr("updatePlayerLocationByPid", "player == nil")
		return nil
	}

	player.LocationX = locationX
	player.LocationY = locationY
	player.LocationZ = locationZ
	player.LocationDetail = locationDetail

	return player
}

func changePlayerNameByPid(args service.Args) service.Result {
	if len(args) < 2 {
		utils.ArgsNumberErr("changePlayerNameByPid", 2)
		return nil
	}

	pid, ok := args[0].(int64)
	if !ok {
		utils.ArgsTypeCastErr("changePlayerNameByPid", 0)
		return nil
	}

	if pid <= 0 {
		utils.InvalidValueErr("changePlayerNameByPid", "pid <= 0")
		return nil
	}

	name, ok := args[1].(string)
	if !ok {
		utils.ArgsTypeCastErr("changePlayerNameByPid", 1)
		return nil
	}

	if len(name) == 0 {
		utils.InvalidValueErr("changePlayerNameByPid", "len(name) == 0")
		return nil
	}

	player := getInnerPlayerByPid(pid)
	if player == nil {
		utils.InvalidValueErr("changePlayerNameByPid", "player == nil")
		return nil
	}

	if name == player.Name {
		return player
	}

	playerInfo := &db.PlayerInfo{
		Pid:  player.Pid,
		Name: name,
	}
	if !db.UpdatePlayerInfo(playerInfo) {
		utils.InvalidValueErr("changePlayerNameByPid", "db.UpdatePlayerInfo failed")
		return nil
	}

	player.Name = name
	return player
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

func getInnerPlayerNum() int {
	return len(pidToPlayer)
}
