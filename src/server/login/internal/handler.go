package internal

import (
	"net/http"
	"server/login/internal/handler"
	"time"

	"github.com/name5566/leaf/log"
)

////////////////////////////////////////////
// type, const, var
//
type Handler func(http.ResponseWriter, *http.Request)

var (
	urlToHandler = map[string]Handler{
		"/push": handler.HandlePush,
		"/pull": handler.HandlePull,
	}
)

////////////////////////////////////////////
// func
//
func init() {
	for url, handler := range urlToHandler {
		http.HandleFunc(url, handler)
	}

	go startHttpServer()
}

func startHttpServer() {
	server := &http.Server{
		Addr:         ":9123",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("login ListenAndServe failed, %s", err)
	}
}
