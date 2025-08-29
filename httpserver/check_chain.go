package httpserver

import (
	"encoding/json"
	"io"
	"net/http"
	"os/exec"
	"strconv"
	"strings"

	clog "github.com/kpango/glg"
)

// HandleCheckNode 检查链节点状态
func HandleCheckNode(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		msg := "请求方法 " + r.Method + " 不允许"
		http.Error(w, msg, http.StatusMethodNotAllowed)
		return
	}

	// 获取节点名称
	nodeName, err := GetNodeName()
	if err != nil {
		http.Error(w, "获取节点名称失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 检查 chainmaker 进程是否存在
	pid, err := checkNodeState()
	if err != nil {
		http.Error(w, "检查节点状态失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var resp Response
	if pid != 0 {
		resp = Response{
			Status:  "success",
			Message: "节点 " + nodeName + " 正在运行, PID: " + strconv.Itoa(pid),
		}
	} else {
		resp = Response{
			Status:  "error",
			Message: "节点 " + nodeName + " 未运行",
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// HandleCheckChain 检查链状态
func HandleCheckChain(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		msg := "请求方法 " + r.Method + " 不允许"
		http.Error(w, msg, http.StatusMethodNotAllowed)
		return
	}

	var msg string
	ok := true

	// 启动四个节点
	for node, addr := range NodeAddrs {
		url := "http://" + addr + CheckNodePath
		clog.Infof("检查节点 %s, 请求地址 %s", node, url)
		resp, err := http.Get(url)
		if err != nil {
			msg += "检查节点 " + node + " 失败: " + err.Error() + "\n"
			continue
		}

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			ok = false
			msg += "检查节点 " + node + " 失败: " + string(body) + "\n"
			continue
		}

		var nodeResp Response
		err = json.NewDecoder(resp.Body).Decode(&nodeResp)
		resp.Body.Close()
		if err != nil {
			ok = false
			msg += "解析节点 " + node + " 响应失败: " + err.Error() + "\n"
			continue
		}

		if nodeResp.Status != "success" {
			ok = false
		}

		msg += nodeResp.Message + "\n"
	}

	var resp Response
	if ok {
		resp = Response{
			Status:  "success",
			Message: "长安链状态正常:\n" + msg,
		}
	} else {
		resp = Response{
			Status:  "error",
			Message: "长安链状态异常:\n" + msg,
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// 检查 chainmaker 进程是否存在
func checkNodeState() (int, error) {
	command := "ps -ef | grep '[c]hainmaker start'"

	cmd := exec.Command("bash", "-c", command)
	output, err := cmd.Output()
	if err != nil {
		return 0, nil // 没有找到进程
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		pidStr := fields[1]
		pid, err := strconv.Atoi(pidStr)
		if err == nil {
			return pid, nil
		}
	}
	return 0, nil
}
