package db

import (
	"server/base/db"

	"github.com/name5566/leaf/db/mongodb"
	"github.com/name5566/leaf/log"
)

////////////////////////////////////////////
// type, const, var
//
const (
	DBName = "easy_talk_db"
)

const (
	CounterCollectionName = "easy_talk_counter"
)

////////////////////////////////////////////
// func
//

func init() {
	// context
	ctxt, err := mongodb.Dial("jake:123456@localhost", 20)
	if err != nil {
		log.Fatal("db Dial failed: %s", err.Error())
		return
	}

	db.InitContext(ctxt)

	//
	// sub init
	//
	if !ensurePlayerIndex() || !ensurePlayerCounter() {
		log.Fatal("ensurePlayerIndex or ensurePlayerCounter failed")
	}

	if !ensureMsgIndex() || !ensureMsgCounter() {
		log.Fatal("ensureMsgIndex or ensureMsgCounter failed")
	}
}
