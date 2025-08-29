package httpserver

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"

	clog "github.com/kpango/glg"
)

func HandleSetupContracts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Transfer-Encoding", "chunked")

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

	// 判断是否允许操作 relayer
	if !AllowHandleRelayer() {
		http.Error(w, "当前节点 "+nodeName+" 不允许操作 relayer", http.StatusForbidden)
		return
	}

	// 判断 deploytools 可执行文件是否存在
	if _, err := os.Stat("/home/chainmaker/relayer/deploytools"); err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "deploytools 不存在", http.StatusInternalServerError)
			return
		}
		http.Error(w, "检查 deploytools 文件失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 构造部署命令
	command := "cd /home/chainmaker/relayer"
	// 判断是否有 conf_name 参数
	confName := r.URL.Query().Get("conf_name")
	if confName == "" {
		command += " && ./deploytools"
	} else {
		command += " && ./deploytools -c " + confName
	}
	clog.Debugf("部署合约命令: %s", command)

	// 执行部署命令
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
	clog.Infof("节点 %s 部署合约完成", nodeName)
	clog.Debugf("节点 %s 输出: %s", nodeName, output)

	// 写入完成响应
	msg := fmt.Sprintf("节点 %s 部署合约完成\n", nodeName)
	w.Write([]byte(msg))
}
