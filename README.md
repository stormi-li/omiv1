# 简介
## 什么是 OMI？
OMI 是一款用于搭建全栈微服务的框架，它基于 Go 和 Redis 开发，并提供了一套搭建前后端微服务和反向代理微服务的方法，帮助你高效地搭建微服务项目。
## 为什么选择 OMI？
### 基于 Redis 搭建注册中心
基于几乎每个项目都会使用的 Redis 搭建注册中心，无需专门的注册中心中间件，降低了架构复杂度，同时降低了学习成本和使用成本。下面是一个最基本的搭建反向代理的示例：
```go
package main

import (
	"net/http"

	omi "github.com/stormi-li/omiv1"
)

var RedisAddr = "localhost:6379"

func main() {
	options := &omi.Options{Addr: RedisAddr} // 设置连接选项

	proxy := omi.NewProxy(options) // 连接 Redis 创建 Proxy

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeProxy(w, r) // 启动代理服务
	})

	http.ListenAndServe(":80", nil)
}
```
在浏览器输入
```
http:localhost
```
终端输出结果
```
解析失败: /
解析失败: /favicon.ico
```
由于我们还未启动并注册任何一个微服务，所以会解析失败，但上面的示例已经展示了 OMI 的三个核心功能中的两个功能：
* 连接 Redis 客户端：由 OMI 创建的所有服务均需要连接 Redis 客户端以注册自身服务或发现其它服务，所有连接在同一个 Redis 上同一个 DB 的服务均可相互发现和调用。
* 启动 Proxy 服务：Proxy 服务会自动解析请求中的域名和路径，如果在 Redis 上发现了对应的服务会自动回去对应的地址并将请求转发到对应的节点。

你可能会有些疑问，比如代理是如何发现和解析服务的，怎么使用 Redis 注册和发现服务———先别急，在后续的文档中我会详细介绍每一个细节。
### 使用 RESTful 风格调用远程服务
前后端均采用 RESTful 风格调用远程服务，这种风格的统一使得远程调用使用起来异常简单。远程调用的更高级用法会在后续文档进行详细演示，这里展示一个最基本的示例：
注册并启动一个最基本的微服务：
```go
package main

import (
	"fmt"
	"net/http"

	omi "github.com/stormi-li/omiv1"
)

var RedisAddr = "localhost:6379"

func main() {
	options := &omi.Options{Addr: RedisAddr}

	register := omi.NewRegister(options) 
	
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello world")
	})

    // 注册并启动微服务
	register.RegisterAndServe("service-demo", "localhost:9014", func(port string) {
		http.ListenAndServe(port, nil)
	})
}
```
在浏览器输入
```
http://localhost/service-demo
```
返回结果
```
hello world
```
上述示例展示了最后一个核心功能———服务注册，并使用 URL 的方式调用了远程服务
* 服务注册：注册服务以 [服务名+地址] 的形式注册在 Redis 上，其它服务使用 http 协议通过将 Path 的第一个参数设置为服务名调用该服务。

你可能会觉得很神奇，每个服务仅仅通过连接 Redis 就能调用其他服务，整个过程没有配置任何地址信息，仿佛地址不存在了一般，别着急，我会在后续文档里详细展示其原理。
# 快速上手

## 安装
注意：go 版本 >= 1.18
```go
go get github.com/stormi-li/omiv1
```
## 启动 Redis
在使用后续功能前你需要先启动一个 Redis 实例，并且当前只支持 Redis 单例模式。建议新建一个 Redis 实例用于搭建注册中心，防止 Key 冲突。如果要进行服务隔离可以用 DB 区分，或者使用其它 Redis 实例。在这里我们假设你已经启动了一个 Redis 实例，并且默认 Addr:"localhost:6379"，Password:""， DB:0
## 启动配置中心监视器
### 代码
```go
package main

import (
	omi "github.com/stormi-li/omiv1"
)

var RedisAddr = "localhost:6379"

func main() {
	options := &omi.Options{Addr: RedisAddr}

	m := omi.NewMonitor(options)

	m.Start("localhost:9013")
}
```
### 在浏览器输入
```
http://localhost:9013
```
### 如果看到这个画面表示监视服务启动成功，并且 monitor 服务已经注册在了 Redis 上面


![60f3962373f0d72c311ef533b4916fb6](media/17329824885630/60f3962373f0d72c311ef533b4916fb6.png)
## 注册并启动反向代理服务
### 代码
```go
package main

import (
	"net/http"

	omi "github.com/stormi-li/omiv1"
)

var RedisAddr = "localhost:6379"

func main() {
	options := &omi.Options{Addr: RedisAddr}

	proxy := omi.NewProxy(options)
	register := omi.NewRegister(options)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeProxy(w, r)
	})

	register.RegisterAndServe("http-80代理", "localhost:80", func(port string) {
		http.ListenAndServe(port, nil)
	})
}
```
### 如果 monitor 页面显示如下表示代理注册成功。

![9c8c11e0930057e1d87f57edf81fef94](media/17329824885630/9c8c11e0930057e1d87f57edf81fef94.png)
### 在浏览器输入，如果页面跳转成功表示代理启动成功。
```
http://localhost/monitor/
```
## 启动 web 服务
```go
package main

import (
	"net/http"

	omi "github.com/stormi-li/omiv1"
)

var RedisAddr = "localhost:6379"

func main() {
	web := omi.NewWeb(nil)
	web.GenerateTemplate()

	options := &omi.Options{Addr: RedisAddr}

	proxy := omi.NewProxy(options)
	register := omi.NewRegister(options)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if web.ServeWeb(w, r) {
			return
		}
		proxy.ServeProxy(w, r)
	})

	register.RegisterAndServe("localhost", "localhost:9014", func(address string) {
		http.ListenAndServe(address, nil)
	})
}
```
### 在浏览器输入
```
http://localhost
```
### 看到如下页面表示注册并启动成功。
![f0ae39b72e9a3bf55ce67cf12f93d29e](media/17329824885630/f0ae39b72e9a3bf55ce67cf12f93d29e.png)

你还可以在 monitor 页面查看注册详情。
## 启动后端服务
### 代码
```go

package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	omi "github.com/stormi-li/omiv1"
)

var RedisAddr = "localhost:6379"

func main() {
	options := &omi.Options{Addr: RedisAddr}

	register := omi.NewRegister(options)

	http.HandleFunc("/http", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello, send by http")
	})

	http.HandleFunc("/websocket", func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{}
		c, _ := upgrader.Upgrade(w, r, nil)
		c.WriteMessage(1, []byte("hello, send by websocket"))
		time.Sleep(100 * time.Millisecond)
		c.Close()
	})

	register.RegisterAndServe("hello", "localhost:9015", func(port string) {
		http.ListenAndServe(port, nil)
	})
}
```
### 测试请求
在刚刚打开的 web 页面中输入 http 请求地址
```
hello/http
```
返回如下
```
请求地址: hello/http
返回数据:hello, send by http
```
在刚刚打开的 web 页面中输入 websocket 请求地址
```
hello/websocket
```
返回如下
```
WebSocket 连接已建立: hello/websocket
接收到数据:hello, send by websocket
WebSocket 连接已关闭。
```
## 下一步
上述示例已经演示了 OMI 框架的全部模块，包括监控服务、反向代理服务、前端服务和后端服务，你已经可以使用这些功能进行微服务全栈开发了，如果你想了解其原理，或注意事项和更高级的功能，可以继续访问官网[https://stormili.site](https://stormili.site)。
