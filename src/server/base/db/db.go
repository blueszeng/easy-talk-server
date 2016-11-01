package db

import (
	"server/base/utils"

	"github.com/name5566/leaf/db/mongodb"
)

////////////////////////////////////////////
// type, const, var
//

const (
	SetCmd = "$set"
)

var (
	context *mongodb.DialContext
)

////////////////////////////////////////////
// func
//

func InitContext(ctxt *mongodb.DialContext) bool {
	if ctxt == nil {
		utils.InvalidValueErr("db.InitContext", "ctxt == nil")
		return false
	}

	context = ctxt
	return true
}

func Session() *mongodb.Session {
	if context == nil {
		utils.InvalidValueErr("db.Session not init", "context == nil")
		return nil
	}
	return context.Ref()
}

func Disuse(s *mongodb.Session) {
	if context == nil {
		utils.InvalidValueErr("db.Disuse not init", "context == nil")
		return
	}
	context.UnRef(s)
}

func EnsureUniqueIndex(db string, collection string, key []string) error {
	if context == nil {
		utils.InvalidValueErr("db.EnsureUniqueIndex not init", "context == nil")
		return nil
	}
	return context.EnsureUniqueIndex(db, collection, key)
}

func EnsureIndex(db string, collection string, key []string) error {
	if context == nil {
		utils.InvalidValueErr("db.EnsureIndex not init", "context == nil")
		return nil
	}
	return context.EnsureIndex(db, collection, key)
}

func EnsureCounter(db string, collection string, id string) error {
	if context == nil {
		utils.InvalidValueErr("db.EnsureCounter not init", "context == nil")
		return nil
	}
	return context.EnsureCounter(db, collection, id)
}

func NextSeq(db string, collection string, id string) (uint64, error) {
	if context == nil {
		utils.InvalidValueErr("db.NextSeq not init", "context == nil")
		return 0, nil
	}
	return context.NextSeq(db, collection, id)
}
