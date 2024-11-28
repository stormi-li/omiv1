package main

var redisAddr = "118.25.196.166:3934"
var password = "12982397StrongPassw0rd"

type User struct {
	ID    int
	Name  string
	Email string
}

// func main() {
// 	c := omi.NewClient(&redis.Options{Addr: redisAddr, Password: password})
// 	omiserver := c.NewServer("hello_server", "118.25.196.166:7860")
// 	omiserver.HandleFunc("/hello", func(w omidto.ResponseWriter, r *omidto.Request) {
// 		fmt.Fprintf(w, "hello afasdfasdfii")
// 	})
// 	omiserver.HandleFunc("/socket_hello", func(w omidto.ResponseWriter, r *omidto.Request) {
// 		upgrader := websocket.Upgrader{}
// 		c, err := upgrader.Upgrade(w.ResponseWriter, r.Request, nil)
// 		if err != nil {
// 			fmt.Println(err)
// 			return
// 		}
// 		c.WriteMessage(1, []byte("nihaonihoanihao"))
// 	})
// 	omiserver.HandleFunc("/user", func(w omidto.ResponseWriter, r *omidto.Request) {
// 		var user User
// 		r.OmiRead(&user)
// 		user.Name = "server"
// 		w.OmiWrite(&user)
// 	})
// 	omiserver.Start(nil)
// }
