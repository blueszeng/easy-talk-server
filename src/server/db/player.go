package db

import (
	"server/base/db"

	"github.com/name5566/leaf/log"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type PlayerInfo struct {
	Did  string
	Pid  int64
	Name string
}

const (
	PlayerCollectionName = "player"

	PlayerCollectionDidKey  = "did"
	PlayerCollectionPidKey  = "pid"
	PlayerCollectionNameKey = "name"

	PlayerPidCounterKey = "pid"
)

func LoadPlayerInfoByDid(did string) *PlayerInfo {
	s := db.Session()
	defer db.Disuse(s)

	c := s.DB(DBName).C(PlayerCollectionName)

	info := &PlayerInfo{}
	err := c.Find(bson.M{
		PlayerCollectionDidKey: did,
	}).One(info)
	if err != nil {
		return nil
	}
	return info
}

func AddPlayerInfo(info *PlayerInfo) bool {
	if info == nil {
		return false
	}

	pid := nextPid()
	if pid <= 0 {
		return false
	}

	info.Pid = pid

	s := db.Session()
	defer db.Disuse(s)

	c := s.DB(DBName).C(PlayerCollectionName)
	if !addPlayerInfo(c, info) {
		return false
	}
	return true
}

func UpdatePlayerInfo(info *PlayerInfo) bool {
	if info == nil {
		return false
	}

	s := db.Session()
	defer db.Disuse(s)

	c := s.DB(DBName).C(PlayerCollectionName)
	if !updatePlayerInfo(c, info) {
		return false
	}

	return true
}

//
// inner func
//
func ensurePlayerIndex() bool {
	err := db.EnsureUniqueIndex(DBName, PlayerCollectionName, []string{PlayerCollectionPidKey})
	if err != nil {
		log.Error("ensurePlayerIndex pid failed, err: %s", err)
		return false
	}

	err = db.EnsureIndex(DBName, PlayerCollectionName, []string{PlayerCollectionDidKey})
	if err != nil {
		log.Error("ensurePlayerIndex did failed, err: %s", err)
		return false
	}

	return true
}

func ensurePlayerCounter() bool {
	err := db.EnsureCounter(DBName, CounterCollectionName, PlayerPidCounterKey)
	if err != nil {
		log.Error("ensurePlayerCounter failed, err: %s", err)
		return false
	}
	return true
}

func nextPid() int64 {
	id, err := db.NextSeq(DBName, CounterCollectionName, PlayerPidCounterKey)
	if err != nil {
		log.Error("nextPid failed, err: %s", err)
		return 0
	}
	return int64(id)
}

func addPlayerInfo(c *mgo.Collection, info *PlayerInfo) bool {
	err := c.Insert(info)
	if err != nil {
		log.Error("addPlayerInfo failed, err: %s", err)
		return false
	}
	return true
}

func updatePlayerInfo(c *mgo.Collection, info *PlayerInfo) bool {
	err := c.Update(bson.M{
		PlayerCollectionPidKey: info.Pid,
	}, bson.M{
		db.SetCmd: bson.M{
			PlayerCollectionNameKey: info.Name,
		},
	})
	if err != nil {
		log.Error("updatePlayerInfo failed, err: %s", err)
		return false
	}
	return true
}
