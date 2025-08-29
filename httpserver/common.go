package httpserver

import (
	"fmt"
	"os"
)

const (
	// StartNodePath 启动链节点请求路径
	StartNodePath = "/start-node"
	// StopNodePath 停止链节点请求路径
	StopNodePath = "/stop-node"
	// CleanNodePath 清理链节点数据请求路径
	CleanNodePath = "/clean-node"
	// CheckNodePath 查询链节点状态请求路径
	CheckNodePath = "/check-node"

	// StartChainPath 启动链请求路径
	StartChainPath = "/start-chain"
	// StopChainPath 停止链请求路径
	StopChainPath = "/stop-chain"
	// CleanChainPath 清理链数据请求路径
	CleanChainPath = "/clean-chain"
	// CheckChainPath 查询链状态请求路径
	CheckChainPath = "/check-chain"

	// ReceiveConfPath 接收配置文件请求路径
	ReceiveConfPath = "/receive-conf"

	// StartRelayerPath 启动 relayer 请求路径
	StartRelayerPath = "/start-relayer"
	// StopRelayerPath 停止 relayer 请求路径
	StopRelayerPath = "/stop-relayer"

	// SetupContractsPath 部署合约请求路径
	SetupContractsPath = "/setup-contracts"

	// SubscribeRelayerPath 监听 relayer 日志请求路径
	SubscribeRelayerPath = "/subscribe-relayer"
)

// NodeAddrs 长安链节点列表
var NodeAddrs = map[string]string{
	"Server5": "39.99.41.202:60082",
	"Server6": "47.111.78.240:60082",
	"Server7": "115.29.212.238:60082",
	"Server8": "112.74.104.251:60082",
}

// GetNodeName 获取长安链节点名称
func GetNodeName() (string, error) {
	nodeName, ok := os.LookupEnv("cm_node_name")
	if !ok {
		return "", fmt.Errorf("环境变量 cm_node_name 未设置")
	}
	return nodeName, nil
}

// GetChainmakerBinPath 获取 Chainmaker 控制脚本文件夹路径
func GetChainmakerBinPath() (string, error) {
	nodeName, err := GetNodeName()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("/home/chainmaker/chainmaker/chainmaker-v2.3.5-%s/bin", nodeName), nil
}

// AllowHandleRelayer 判断当前节点是否允许操作 relayer
// 只有 node1 和 node2 所在服务器允许操作 relayer
func AllowHandleRelayer() bool {
	nodeName, err := GetNodeName()
	if err != nil {
		return false
	}
	if nodeName == "node1" || nodeName == "node2" {
		return true
	}
	return false
}
