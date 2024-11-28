package monitor

import (
	"encoding/json"
	"net/http"

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
	jsonStr, _ := json.Marshal(data)
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
}
