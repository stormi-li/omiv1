package main

import (
	"fmt"

	"github.com/gorilla/websocket"
)

func main() {
	d := websocket.Dialer{}
	c, _, err := d.Dial("ws://localhost:9015/websocket", nil)
	fmt.Println(err)
	_, data, err := c.ReadMessage()
	fmt.Println(string(data))
	fmt.Println(err)
}
