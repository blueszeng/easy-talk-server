package static_data

import (
	"reflect"
	"server/base/utils"
	"strconv"

	"github.com/name5566/leaf/log"
	"github.com/name5566/leaf/recordfile"
)

////////////////////////////////////////////
// type, const, var
//
type DataManager map[string]*recordfile.RecordFile
type Key interface{}
type DataUnit interface{}

var (
	dataManager = make(DataManager)
)

////////////////////////////////////////////
// func
//
func Load(st interface{}) {
	rf := readRf(st)
	stName := reflect.TypeOf(st).Name()
	dataManager[stName] = rf
}

func GetDataUnitByKey(stName string, key Key) DataUnit {
	rf := getData(stName)
	if rf == nil {
		return nil
	}

	unit := rf.Index(key)
	if unit == nil {
		utils.InvalidValueErr("static_data.GetDataUnitByKey", "key, stName:"+stName)
		return nil
	}
	return unit
}

func GetDataCount(stName string) int {
	rf := getData(stName)
	if rf == nil {
		return 0
	}
	return rf.NumRecord()
}

func GetDataUnitByIndex(stName string, index int) DataUnit {
	rf := getData(stName)
	if rf == nil {
		return nil
	}

	unit := rf.Record(index)
	if unit == nil {
		utils.InvalidValueErr("static_data.GetDataUnitByIndex", "index: "+strconv.Itoa(index)+", stName:"+stName)
		return nil
	}
	return unit
}

//
// implement
//
func readRf(st interface{}) *recordfile.RecordFile {
	rf, err := recordfile.New(st)
	if err != nil {
		log.Fatal("%v", err)
	}
	fn := reflect.TypeOf(st).Name() + ".txt"
	err = rf.Read("gamedata/" + fn)
	if err != nil {
		log.Fatal("%v: %v", fn, err)
	}

	return rf
}

func getData(stName string) *recordfile.RecordFile {
	rf := dataManager[stName]
	if rf == nil {
		utils.InvalidValueErr("static_data.GetData", "stName: "+stName)
		return nil
	}
	return rf
}
