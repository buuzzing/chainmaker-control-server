package httpserver

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
)

// HandleReceiveConfig 处理接收配置的请求
func HandleReceiveConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
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

	// 限制最大上传大小（如10MB）
	r.ParseMultipartForm(10 << 20)
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "无法获取上传文件: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// 获取文件名参数（可选，优先使用form中的filename字段，否则用上传文件名）
	filename := r.FormValue("filename")
	if filename == "" {
		filename = handler.Filename
	}

	// 保存路径
	saveDir := "/home/chainmaker/relayer"
	savePath := saveDir + "/" + filename

	out, err := os.Create(savePath)
	if err != nil {
		http.Error(w, "保存文件失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		http.Error(w, "写入文件失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	resp := Response{
		Status:  "success",
		Message: "文件已保存在 " + savePath,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// HandleCleanConfig 清理配置文件
func HandleCleanConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
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

	// 只删除 .toml 文件
	saveDir := "/home/chainmaker/relayer"
	files, err := os.ReadDir(saveDir)
	if err == nil {
		for _, f := range files {
			if !f.IsDir() && len(f.Name()) > 5 && f.Name()[len(f.Name())-5:] == ".toml" {
				_ = os.Remove(saveDir + "/" + f.Name())
			}
		}
	}

	resp := Response{
		Status:  "success",
		Message: "配置文件已清理",
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
