package monitor

import (
	"embed"
	"net/http"

	omi "github.com/stormi-li/omiv1"
	"github.com/stormi-li/omiv1/omihttp"
	proxy "github.com/stormi-li/omiv1/omiproxy"
	web "github.com/stormi-li/omiv1/omiweb"
)

//go:embed static/*
var embeddedSource embed.FS

func NewMonitorMux(omiClient *omi.Client) *omihttp.ServeMux {
	nodeManageHandler := NewNodeManageHandler(proxy.NewRouter(omiClient.Register.RedisClient))

	omiClient.Register.AddRegisterHandleFunc("ServerType", func() string {
		return "monitor"
	})

	mux := omiClient.NewServeMux()
	mux.ServerType = omihttp.ServerType_Monitor

	mux.HandleFunc("/GetNodes", func(w http.ResponseWriter, r *http.Request, rw *omihttp.ReadWriter) {
		nodeManageHandler.GetNodes(w, r)
	})
	mux.HandleFunc("/GetNodeInfo", func(w http.ResponseWriter, r *http.Request, rw *omihttp.ReadWriter) {
		nodeManageHandler.GetNodeInfo(w, r)
	})
	mux.HandleFunc("/SendMessage", func(w http.ResponseWriter, r *http.Request, rw *omihttp.ReadWriter) {
		nodeManageHandler.SendMessage(w, r)
	})

	omiweb := web.NewWeb(&embeddedSource)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request, rw *omihttp.ReadWriter) {
		omiweb.ServeWeb(w, r)
	})
	return mux
}
