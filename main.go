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

	//engine.Handle("/config", httpserver.HandleConfig)

	// 接收配置文件
	engine.Handle(httpserver.ReceiveConfPath, httpserver.HandleReceiveConfig)

	// 部署与注册合约
	engine.Handle(httpserver.SetupContractsPath, httpserver.HandleSetupContracts)

	// 启动 relayer
	engine.Handle(httpserver.StartRelayerPath, httpserver.HandleStartRelayer)
	// 停止 relayer
	engine.Handle(httpserver.StopRelayerPath, httpserver.HandleStopRelayer)

	// 监控 relayer 日志
	engine.Handle("/subscribeRelayer", httpserver.HandleSubscribeRelayer)

	// 启动单个链节点
	engine.Handle(httpserver.StartNodePath, httpserver.HandleStartNode)
	// 停止单个链节点
	engine.Handle(httpserver.StopNodePath, httpserver.HandleStopNode)
	// 清理单个链节点数据
	engine.Handle(httpserver.CleanNodePath, httpserver.HandleCleanNode)
	// 检查单个链节点状态
	engine.Handle(httpserver.CheckNodePath, httpserver.HandleCheckNode)

	// 启动长安链
	engine.Handle(httpserver.StartChainPath, httpserver.HandleStartChain)
	// 停止长安链
	engine.Handle(httpserver.StopChainPath, httpserver.HandleStopChain)
	// 清理长安链数据
	engine.Handle(httpserver.CleanChainPath, httpserver.HandleCleanChain)
	// 检查长安链状态
	engine.Handle(httpserver.CheckChainPath, httpserver.HandleCheckChain)

	// 启动 HTTP 服务
	if err := engine.Start(); err != nil {
		panic(err)
	}
}
