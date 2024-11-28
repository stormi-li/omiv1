package monitor

import (
	"embed"
	"net/http"

	"github.com/go-redis/redis/v8"
	proxy "github.com/stormi-li/omiv1/omiproxy"
	server "github.com/stormi-li/omiv1/omiserver"
	web "github.com/stormi-li/omiv1/omiweb"
)

type Monitor struct {
	Server     *server.Server
	Proxy      *proxy.Proxy
	EmbedModel bool
}

func NewMonitor(redisClient *redis.Client, serverName, address string) *Monitor {
	return &Monitor{
		Server:     server.NewServer(redisClient, serverName, address),
		Proxy:      proxy.NewProxy(redisClient, &http.Transport{}),
		EmbedModel: true,
	}
}

//go:embed dev/static/*
var embeddedSource embed.FS
var sourcePath = "dev/static"
var indexPath = "/index.html"

func (monitor *Monitor) Start() {
	monitor.Server.Register.AddRegisterHandleFunc("ServerType", func() string {
		return "monitor"
	})
	nodeManageHandler := NewNodeManageHandler(monitor.Proxy.Reslover.Router)
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
	monitor.Server.Start(nil)
}
