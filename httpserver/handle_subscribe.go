package httpserver

import (
	"net/http"
	"os"
	"time"
)

// HandleSubscribeRelayer 监听 relayer 输出日志
func HandleSubscribeRelayer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Transfer-Encoding", "chunked")

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

	file, err := os.Open("/home/chainmaker/relayer/cmrelayer.log")
	if err != nil {
		http.Error(w, "无法打开日志文件: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// 定位到文件末尾
	stat, _ := file.Stat()
	file.Seek(stat.Size(), 0)

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	notify := r.Context().Done()
	buf := make([]byte, 4096)
	for {
		select {
		case <-notify:
			return // 客户端断开
		default:
			n, err := file.Read(buf)
			if n > 0 {
				w.Write(buf[:n])
				flusher.Flush()
			}
			if err != nil {
				if err.Error() == "EOF" {
					// 没有新内容，稍等再读
					time.Sleep(500 * time.Millisecond)
					continue
				}
				return
			}
		}
	}
}
