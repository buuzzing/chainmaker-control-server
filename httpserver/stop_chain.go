package httpserver

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strings"

	clog "github.com/kpango/glg"
)

// HandleStopNode 停止链节点
func HandleStopNode(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Transfer-Encoding", "chunked")

	if r.Method != http.MethodGet {
		msg := fmt.Sprintf("请求方法 %s 不允许", r.Method)
		http.Error(w, msg, http.StatusMethodNotAllowed)
		return
	}

	// 停止脚本绝对路径
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

	command := fmt.Sprintf("cd %s && ./stop.sh full", binPath)
	clog.Debugf("停止节点 %s 命令: %s", nodeName, command)

	// 执行启动命令
	cmd := exec.Command("bash", "-c", command)

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
	var builder strings.Builder
	tee := io.TeeReader(stdout, &builder)
	io.Copy(w, tee)
	output := builder.String()

	// 等待命令执行完成
	if err := cmd.Wait(); err != nil {
		http.Error(w, "命令执行失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 输出日志信息
	clog.Infof("节点 %s 停止完成", nodeName)
	clog.Debugf("节点 %s 输出: %s", nodeName, output)

	// 写入完成响应
	msg := fmt.Sprintf("长安链节点 %s 停止完成\n", nodeName)
	w.Write([]byte(msg))
}

// HandleStopChain 停止长安链
func HandleStopChain(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		msg := fmt.Sprintf("请求方法 %s 不允许", r.Method)
		http.Error(w, msg, http.StatusMethodNotAllowed)
		return
	}

	var msg string

	// 停止四个节点
	for node, addr := range NodeAddrs {
		resp, err := http.Get(fmt.Sprintf("http://%s%s", addr, StopNodePath))
		if err != nil {
			msg += fmt.Sprintf("停止节点 %s 失败: %s\n", node, err.Error())
			continue
		}
		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			msg += fmt.Sprintf("停止节点 %s 失败: %s\n", node, body)
			continue
		}
	}
	// 停止节点阶段出现错误
	if msg != "" {
		resp := Response{
			Status:  "error",
			Message: msg,
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp := Response{
		Status:  "success",
		Message: "长安链已停止",
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
