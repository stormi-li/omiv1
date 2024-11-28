package monitor

import (
	"encoding/json"
	"net/http"
	"sort"

	proxy "github.com/stormi-li/omiv1/omiproxy"
)

type NodeManageHandler struct {
	nodeManager *NodeManager
}

func NewNodeManageHandler(router *proxy.Router) *NodeManageHandler {
	return &NodeManageHandler{
		nodeManager: NewManager(router),
	}
}

func (handler *NodeManageHandler) GetNodes(w http.ResponseWriter, r *http.Request) {
	var data = handler.nodeManager.GetNodes()
	type ServerInfo struct {
		ServerName string
		Address    string
		Weight     string
		Type       bool
	}
	var result []ServerInfo
	for serverName, addressMap := range data {
		for address, info := range addressMap {
			server := ServerInfo{
				ServerName: serverName,
				Address:    address,
				Weight:     info["Weight"],
			}
			result = append(result, server)
		}
	}
	sort.SliceStable(result, func(i, j int) bool {
		return result[i].Address < result[j].Address
	})
	sort.SliceStable(result, func(i, j int) bool {
		return result[i].ServerName < result[j].ServerName
	})
	lastServerName := ""
	lastType := true
	for i := 0; i < len(result); i++ {
		if result[i].ServerName != lastServerName {
			lastType = !lastType
			lastServerName = result[i].ServerName
		}
		result[i].Type = lastType
	}
	jsonStr, _ := json.Marshal(result)
	w.Write([]byte(jsonStr))
}

func (handler *NodeManageHandler) GetNodeInfo(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	address := r.URL.Query().Get("address")
	jsonStr, _ := json.Marshal(handler.nodeManager.GetNodeInfo(name, address))
	w.Write([]byte(jsonStr))
}

func (handler *NodeManageHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	address := r.URL.Query().Get("address")
	command := r.URL.Query().Get("command")
	message := r.URL.Query().Get("message")
	handler.nodeManager.SendCommand(name, address, command, message)
	jsonData, _ := json.Marshal(&Result{200})
	w.Write(jsonData)
}

type Result struct {
	Code int
}
