package httpserver

import (
	"fmt"
	"os"
)

const (
	// StartNodePath 启动链节点请求路径
	StartNodePath = "/startNode"
	// StopNodePath 停止链节点请求路径
	StopNodePath = "/stopNode"
	// CleanNodePath 清理链节点数据请求路径
	CleanNodePath = "/cleanNode"

	// StartChainPath 启动链请求路径
	StartChainPath = "/startChain"
	// StopChainPath 停止链请求路径
	StopChainPath = "/stopChain"
	// CleanChainPath 清理链数据请求路径
	CleanChainPath = "/cleanChain"

	// CheckStatusPath 查询链节点状态请求路径
	CheckStatusPath = "/checkStatus"
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
