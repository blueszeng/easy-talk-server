package handler

import (
	"encoding/json"
	"net/http"
	"server/base/utils"

	"github.com/name5566/leaf/log"
)

func checkMethodAndParseForm(req *http.Request, tag string) bool {
	if req.Method != "POST" {
		log.Error("%s, invalid req Method: %s", tag, req.Method)
		return false
	}

	err := req.ParseForm()
	if err != nil {
		log.Error("%s, req ParseForm failed", tag)
		return false
	}

	return true
}

func parseString(req *http.Request, key string) (string, bool) {
	datas := req.PostForm[key]
	if len(datas) == 0 {
		return "", false
	}

	return datas[0], true
}

func parseInt64(req *http.Request, key string) (int64, bool) {
	datas := req.PostForm[key]
	if len(datas) == 0 {
		return 0, false
	}

	data, err := utils.StrToInt64(datas[0])
	if err != nil {
		return 0, false
	}
	return data, true
}

func parseFloat32(req *http.Request, key string) (float32, bool) {
	datas := req.PostForm[key]
	if len(datas) == 0 {
		return 0, false
	}

	data, err := utils.StrToFloat32(datas[0])
	if err != nil {
		return 0, false
	}
	return data, true
}

func sayErr(w http.ResponseWriter, errcode int32) {
	w.Write(jsonMarshal(Error{errcode}))
}

func jsonMarshal(info interface{}) []byte {
	json, err := json.Marshal(info)
	if err != nil {
		log.Error("jsonMarshal failed, info: %v, err: %v", info, err)
		return []byte("")
	}
	return json
}
