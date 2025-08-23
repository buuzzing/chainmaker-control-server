package httpserver

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

// HandleStartNode 启动链节点
func HandleStartNode(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Transfer-Encoding", "chunked")

	if r.Method != http.MethodGet {
		msg := fmt.Sprintf("请求方法 %s 不允许", r.Method)
		http.Error(w, msg, http.StatusMethodNotAllowed)
		return
	}

	// 启动脚本绝对路径
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
	startBinPath := binPath + "/start.sh"

	// 执行启动命令
	cmd := exec.Command(startBinPath, "full", "-y")

	// 获取 stdout 管道
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		http.Error(w, "无法获取 stdout: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 启动命令
	if err := cmd.Start(); err != nil {
		http.Error(w, "无法启动命令: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 使用 http.ResponseWriter 写入响应
	io.Copy(w, stdout)

	// 等待命令执行完成
	if err := cmd.Wait(); err != nil {
		http.Error(w, "命令执行失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 写入完成响应
	w.WriteHeader(http.StatusOK)
	msg := fmt.Sprintf("长安链节点 %s 启动完成\n", nodeName)
	io.WriteString(w, msg)
}

// HandleCheckChainStatus 检查链状态
func HandleCheckChainStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		msg := fmt.Sprintf("请求方法 %s 不允许", r.Method)
		http.Error(w, msg, http.StatusMethodNotAllowed)
		return
	}

	// 获取节点名称
	nodeName, err := GetNodeName()
	if err != nil {
		http.Error(w, "获取节点名称失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 检查日志文件中是否包含 "all necessary"
	found, err := checkAllNecessary()
	if err != nil {
		http.Error(w, "检查日志文件失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var resp Response
	if found {
		resp = Response{
			Status:  "success",
			Message: fmt.Sprintf("节点 %s 已启动并运行正常", nodeName),
		}
	} else {
		resp = Response{
			Status:  "error",
			Message: fmt.Sprintf("节点 %s 未启动或运行异常", nodeName),
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// HandleStartChain 启动链
func HandleStartChain(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		msg := fmt.Sprintf("请求方法 %s 不允许", r.Method)
		http.Error(w, msg, http.StatusMethodNotAllowed)
		return
	}

	var msg string

	// 启动四个节点
	for node, addr := range NodeAddrs {
		resp, err := http.Get(fmt.Sprintf("http://%s%s", StartNodePath, addr))
		if err != nil {
			msg += fmt.Sprintf("启动节点 %s 失败: %v\n", node, err)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			msg += fmt.Sprintf("启动节点 %s 失败: %v\n", node, body)
			continue
		}
	}
	// 启动节点阶段出现错误
	if msg != "" {
		resp := Response{
			Status:  "error",
			Message: msg,
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	time.Sleep(5 * time.Second)
	// 检查链状态
	found, err := checkAllNecessary()
	if err != nil {
		http.Error(w, "检查日志文件失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var resp Response
	if found {
		resp = Response{
			Status:  "success",
			Message: "长安链启动成功",
		}
	} else {
		resp = Response{
			Status:  "error",
			Message: "长安链启动失败",
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)

}

// 检查日志文件中是否包含 "all necessary"
func checkAllNecessary() (bool, error) {
	// 日志路径
	binPath, err := GetChainmakerBinPath()
	if err != nil {
		return false, err
	}
	logFilePath := binPath + "../log/system.log"

	// 打开日志文件
	file, err := os.Open(logFilePath)
	if err != nil {
		return false, err
	}
	defer file.Close()

	// 逐行扫描日志文件
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "all necessary") {
			return true, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return false, err
	}

	return false, nil
}
