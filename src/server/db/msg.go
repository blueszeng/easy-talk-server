package db

import (
	"server/base/db"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/name5566/leaf/log"
)

type MsgInfo struct {
	Mid            int64
	Pid            int64
	Date           int64
	LocationX      float32
	LocationY      float32
	LocationZ      float32
	LocationDetail string
	Channel        int32
	MsgType        string
	Msg            string
}

const (
	MsgCollectionName = "msg"

	MsgCollectionMidKey = "mid"
	MsgCollectionPidKey = "pid"

	MsgMidCounterKey = "mid"
)

func GetLastMsgId() int64 {
	s := db.Session()
	defer db.Disuse(s)

	c := s.DB(DBName).C(CounterCollectionName)

	var st struct {
		Seq uint64
	}
	err := c.FindId(MsgMidCounterKey).One(&st)
	if err != nil {
		log.Error("GetLastMsgId failed, err: %s", err)
		return 0
	}
	return int64(st.Seq)
}

func LoadMsgInfoByMid(mid int64) *MsgInfo {
	s := db.Session()
	defer db.Disuse(s)

	c := s.DB(DBName).C(MsgCollectionName)

	info := &MsgInfo{}
	err := c.Find(bson.M{
		MsgCollectionMidKey: mid,
	}).One(info)
	if err != nil {
		return nil
	}
	return info
}

func AddMsgInfo(info *MsgInfo) bool {
	if info == nil {
		return false
	}

	mid := nextMid()
	if mid <= 0 {
		return false
	}

	info.Mid = mid

	s := db.Session()
	defer db.Disuse(s)

	c := s.DB(DBName).C(MsgCollectionName)
	if !addMsgInfo(c, info) {
		return false
	}
	return true
}

//
// inner func
//
func ensureMsgIndex() bool {
	err := db.EnsureUniqueIndex(DBName, MsgCollectionName, []string{MsgCollectionMidKey})
	if err != nil {
		log.Error("ensureMsgIndex mid failed, err: %s", err)
		return false
	}

	err = db.EnsureIndex(DBName, PlayerCollectionName, []string{MsgCollectionPidKey})
	if err != nil {
		log.Error("ensureMsgIndex pid failed, err: %s", err)
		return false
	}

	return true
}

func ensureMsgCounter() bool {
	err := db.EnsureCounter(DBName, CounterCollectionName, MsgMidCounterKey)
	if err != nil {
		log.Error("ensureMsgCounter failed, err: %s", err)
		return false
	}
	return true
}

func nextMid() int64 {
	id, err := db.NextSeq(DBName, CounterCollectionName, MsgMidCounterKey)
	if err != nil {
		log.Error("nextMid failed, err: %s", err)
		return 0
	}
	return int64(id)
}

func addMsgInfo(c *mgo.Collection, info *MsgInfo) bool {
	err := c.Insert(info)
	if err != nil {
		log.Error("addMsgInfo failed, err: %s", err)
		return false
	}
	return true
}
