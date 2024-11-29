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
	ProductEnv bool
}

func NewMonitor(register *register.Register) *Monitor {
	return &Monitor{
		Register:   register,
		ProductEnv: true,
	}
}

//go:embed static/*
var embeddedSource embed.FS

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
	omiweb := web.NewWeb(nil)
	if monitor.ProductEnv {
		omiweb = web.NewWeb(&embeddedSource)
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		omiweb.ServeWeb(w, r)
	})
	monitor.Register.RegisterAndServe("monitor", address, func(port string) {
		http.ListenAndServe(port, nil)
	})
}
