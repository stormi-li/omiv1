package register

import (
	"net/http"
	"strconv"

	"github.com/stormi-li/omiv1/omirpc"
)

type ServeMux struct {
	*http.ServeMux
	Register  *Register
	RouterMap map[string]bool
}

func NewServeMux(register *Register) *ServeMux {
	return &ServeMux{ServeMux: http.NewServeMux(), Register: register, RouterMap: map[string]bool{}}
}

func (mux *ServeMux) Handle(pattern string, handler omirpc.Handler) {
	mux.HandleFunc(pattern, handler.ServeHTTP)
}

func (mux *ServeMux) HandleFunc(pattern string, handler func(w http.ResponseWriter, r *http.Request, rw *omirpc.ReadWriter)) {
	mux.RouterMap[pattern] = true
	mux.Register.AddRegisterHandleFunc("HandleFunc["+pattern+"]State", func() string {
		if mux.RouterMap[pattern] {
			return "open"
		} else {
			return "closed"
		}
	})
	mux.Register.AddMessageHandleFunc("SwitchHandleFunc["+pattern+"]State", func(message string) {
		state, err := strconv.Atoi(message)
		if err == nil {
			if state == 1 {
				mux.RouterMap[pattern] = true
			}
			if state == 0 {
				mux.RouterMap[pattern] = false
			}
		}
	})
	mux.ServeMux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		if mux.RouterMap[pattern] {
			handler(w, r, omirpc.NewReadWriter(w, r))
		}
	})
}
