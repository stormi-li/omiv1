package monitor

import (
	"embed"
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/stormi-li/omiv1/omihttp"
	proxy "github.com/stormi-li/omiv1/omiproxy"
	web "github.com/stormi-li/omiv1/omiweb"
)

//go:embed static/*
var embeddedSource embed.FS

func NewMonitorMux(mux *omihttp.ServeMux, redisClient *redis.Client) *omihttp.ServeMux {
	nodeManageHandler := NewNodeManageHandler(proxy.NewRouter(redisClient))

	mux.ServeMux.HandleFunc("/GetNodes", func(w http.ResponseWriter, r *http.Request) {
		nodeManageHandler.GetNodes(w, r)
	})
	mux.ServeMux.HandleFunc("/GetNodeInfo", func(w http.ResponseWriter, r *http.Request) {
		nodeManageHandler.GetNodeInfo(w, r)
	})
	mux.ServeMux.HandleFunc("/SendMessage", func(w http.ResponseWriter, r *http.Request) {
		nodeManageHandler.SendMessage(w, r)
	})

	omiweb := web.NewWebServer(&embeddedSource)
	mux.ServeMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		omiweb.ServeWeb(w, r)
	})
	return mux
}
