package main

import (
	"flag"
	"test_server/httpserver"
)

var port *string

func init() {
	port = flag.String("p", "60082", "长安链控制服务端口号")
}

func main() {
	flag.Parse()

	// 实例化自定义 HttpEngine
	engine := httpserver.NewHttpEngine(*port)

	// 注册路由与对应处理函数
	engine.Handle("/config", httpserver.HandleConfig)
	engine.Handle("/startChain", httpserver.HandleStartChain)
	engine.Handle("/setupContracts", httpserver.HandleSetupContracts)
	engine.Handle("/startRelayer", httpserver.HandleStartRelayer)
	engine.Handle("/subscribeRelayer", httpserver.HandleSubscribeRelayer)

	engine.Handle(httpserver.StartNodePath, httpserver.HandleStartNode)
	engine.Handle(httpserver.StopNodePath, httpserver.HandleStopNode)
	engine.Handle(httpserver.CleanNodePath, httpserver.HandleCleanNode)
	engine.Handle(httpserver.StartChainPath, httpserver.HandleStartChain)
	engine.Handle(httpserver.StopChainPath, httpserver.HandleStopChain)
	engine.Handle(httpserver.CleanChainPath, httpserver.HandleCleanChain)
	engine.Handle(httpserver.CheckStatusPath, httpserver.HandleCheckChainStatus)

	// 启动 HTTP 服务
	if err := engine.Start(); err != nil {
		panic(err)
	}
}
