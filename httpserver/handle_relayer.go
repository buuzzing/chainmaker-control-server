package httpserver

import (
	"encoding/json"
	"net/http"
	"os"
	"os/exec"

	clog "github.com/kpango/glg"
)

// HandleStartRelayer 启动 relayer
func HandleStartRelayer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		http.Error(w, "请求方法不允许", http.StatusMethodNotAllowed)
		return
	}

	// 获取节点名称
	nodeName, err := GetNodeName()
	if err != nil {
		http.Error(w, "获取节点名称失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 判断是否允许操作 relayer
	if !AllowHandleRelayer() {
		http.Error(w, "当前节点 "+nodeName+" 不允许操作 relayer", http.StatusForbidden)
		return
	}

	// 判断 cmrelayer 可执行文件是否存在
	if _, err := os.Stat("/home/chainmaker/relayer/cmrelayer"); err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "cmrelayer 不存在", http.StatusInternalServerError)
			return
		}
		http.Error(w, "检查 cmrelayer 文件失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 构造启动命令
	command := "cd /home/chainmaker/relayer"
	// 判断是否有 conf_name 参数
	confName := r.URL.Query().Get("conf_name")
	if confName == "" {
		command += " && ./start_cmrelayer.sh"
	} else {
		command += " && ./start_cmrelayer.sh " + confName
	}
	clog.Debugf("启动 relayer 命令: %s", command)

	// 执行启动命令
	cmd := exec.Command("bash", "-c", command)
	err = cmd.Run()
	if err != nil {
		http.Error(w, "启动 relayer 失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	clog.Infof("节点 %s 启动 relayer 成功", nodeName)

	resp := Response{
		Status:  "success",
		Message: "节点 " + nodeName + " 启动 relayer 成功",
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// HandleStopRelayer 停止 relayer
func HandleStopRelayer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		http.Error(w, "请求方法不允许", http.StatusMethodNotAllowed)
		return
	}

	// 获取节点名称
	nodeName, err := GetNodeName()
	if err != nil {
		http.Error(w, "获取节点名称失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 判断是否允许操作 relayer
	if !AllowHandleRelayer() {
		http.Error(w, "当前节点 "+nodeName+" 不允许操作 relayer", http.StatusForbidden)
		return
	}

	// 构造停止命令
	command := "cd /home/chainmaker/relayer && ./stop_cmrelayer.sh"
	// 执行停止命令
	cmd := exec.Command("bash", "-c", command)
	err = cmd.Run()
	if err != nil {
		http.Error(w, "停止 relayer 失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	clog.Infof("节点 %s 停止 relayer 成功", nodeName)

	resp := Response{
		Status:  "success",
		Message: "节点 " + nodeName + " 停止 relayer 成功",
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
