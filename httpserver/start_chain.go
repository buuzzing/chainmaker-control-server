package httpserver

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strings"
	"time"

	clog "github.com/kpango/glg"
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

	command := fmt.Sprintf("cd %s && ./start.sh full -y", binPath)
	clog.Debugf("启动节点 %s 命令: %s", nodeName, command)

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
	go io.Copy(w, tee)

	// 等待命令执行完成
	if err := cmd.Wait(); err != nil {
		http.Error(w, "命令执行失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 输出日志信息
	output := builder.String()
	clog.Infof("节点 %s 启动完成", nodeName)
	clog.Debugf("节点 %s 输出: %s", nodeName, output)

	// 写入完成响应
	msg := fmt.Sprintf("长安链节点 %s 启动完成\n", nodeName)
	w.Write([]byte(msg))
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
		url := fmt.Sprintf("http://%s%s", addr, StartNodePath)
		clog.Infof("启动节点 %s, 请求地址 %s", node, url)
		resp, err := http.Get(url)
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

	time.Sleep(20 * time.Second)
	// 检查链状态
	statusResp, err := http.Get(fmt.Sprintf("http://%s%s", NodeAddrs["Server5"], CheckChainPath))
	if err != nil {
		http.Error(w, "检查链状态失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if statusResp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(statusResp.Body)
		statusResp.Body.Close()
		http.Error(w, "检查链状态失败: "+string(body), http.StatusInternalServerError)
		return
	}
	var status Response
	err = json.NewDecoder(statusResp.Body).Decode(&status)
	statusResp.Body.Close()
	if err != nil {
		http.Error(w, "解析链状态响应失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var resp Response
	if status.Status == "success" {
		resp = Response{
			Status:  "success",
			Message: "长安链启动成功，\n" + status.Message,
		}
	} else {
		resp = Response{
			Status:  "error",
			Message: "长安链启动失败，\n" + status.Message,
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)

}
