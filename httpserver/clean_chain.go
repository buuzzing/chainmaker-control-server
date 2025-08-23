package httpserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"time"
)

// HandleCleanNode 清理链节点数据
func HandleCleanNode(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Transfer-Encoding", "chunked")

	if r.Method != http.MethodGet {
		msg := fmt.Sprintf("请求方法 %s 不允许", r.Method)
		http.Error(w, msg, http.StatusMethodNotAllowed)
		return
	}

	// 节点数据路径
	nodeName, err := GetNodeName()
	if err != nil {
		http.Error(w, "获取节点名称失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	binPath, err := GetChainmakerBinPath()
	if err != nil {
		http.Error(w, "获取 Chainmaker bin 路径失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	cleanPath := binPath + "../"

	// 数据目录
	cmd := exec.Command("rm", "-rf", cleanPath+"data")
	_ = cmd.Start()
	_ = cmd.Wait()
	// 无法删除的部分移动到 /tmp 目录
	tmpName := fmt.Sprintf("/tmp/%s-%s", time.Now().Format("20060102150405"), nodeName)
	cmd = exec.Command("mv", cleanPath+"data/go", tmpName)
	_ = cmd.Start()
	_ = cmd.Wait()

	// 日志
	cmd = exec.Command("rm", "-rf", cleanPath+"log", "bin/panic.log")
	_ = cmd.Start()
	_ = cmd.Wait()

	w.WriteHeader(http.StatusOK)
}

// HandleCleanChain 清理链数据
func HandleCleanChain(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		msg := fmt.Sprintf("请求方法 %s 不允许", r.Method)
		http.Error(w, msg, http.StatusMethodNotAllowed)
		return
	}

	// 清理四个节点
	for node, addr := range NodeAddrs {
		resp, err := http.Get(fmt.Sprintf("http://%s%s", addr, CleanNodePath))
		if err != nil || resp.StatusCode != http.StatusOK {
			http.Error(w, fmt.Sprintf("清理节点 %s 失败: %v", node, err), http.StatusInternalServerError)
			return
		}
	}

	resp := Response{
		Status:  "success",
		Message: "链数据清理成功",
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
