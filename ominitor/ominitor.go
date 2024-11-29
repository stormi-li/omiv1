package monitor

import (
	"embed"
	"net/http"

	proxy "github.com/stormi-li/omiv1/omiproxy"
	register "github.com/stormi-li/omiv1/omiregister"
	web "github.com/stormi-li/omiv1/omiweb"
)

type Monitor struct {
	Register   *register.Register
	EmbedModel bool
}

func NewMonitor(register *register.Register) *Monitor {
	return &Monitor{
		Register:   register,
		EmbedModel: true,
	}
}

//go:embed dev/static/*
var embeddedSource embed.FS
var sourcePath = "dev/static"
var indexPath = "/index.html"

func (monitor *Monitor) Start(address string) {
	monitor.Register.AddRegisterHandleFunc("ServerType", func() string {
		return "monitor"
	})
	nodeManageHandler := NewNodeManageHandler(proxy.NewRouter(monitor.Register.RedisClient))
	http.HandleFunc("/GetNodes", func(w http.ResponseWriter, r *http.Request) {
		nodeManageHandler.GetNodes(w, r)
	})
	http.HandleFunc("/GetNodeInfo", func(w http.ResponseWriter, r *http.Request) {
		nodeManageHandler.GetNodeInfo(w, r)
	})
	http.HandleFunc("/SendMessage", func(w http.ResponseWriter, r *http.Request) {
		nodeManageHandler.SendMessage(w, r)
	})
	omiweb := web.NewWeb("static", "/index.html", nil)
	if monitor.EmbedModel {
		omiweb = web.NewWeb(sourcePath, indexPath, &embeddedSource)
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		omiweb.ServeFile(w, r)
	})
	monitor.Register.Register("monitor", address)
	http.ListenAndServe(monitor.Register.Port, nil)
}
