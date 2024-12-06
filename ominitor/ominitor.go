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

	mux.ServeMux.HandleFunc("/GetNodes", func(w http.ResponseWriter, r *http.Request) {
		nodeManageHandler.GetNodes(w, r)
	})
	mux.ServeMux.HandleFunc("/GetNodeInfo", func(w http.ResponseWriter, r *http.Request) {
		nodeManageHandler.GetNodeInfo(w, r)
	})
	mux.ServeMux.HandleFunc("/SendMessage", func(w http.ResponseWriter, r *http.Request) {
		nodeManageHandler.SendMessage(w, r)
	})

	omiweb := web.NewWeb(&embeddedSource)
	mux.ServeMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		omiweb.ServeWeb(w, r)
	})
	return mux
}
