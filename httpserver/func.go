package httpserver

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"test_server/config"

	clog "github.com/kpango/glg"
)

// SingleConfig 单个配置
var SingleConfig *config.SingleConfig

// HandleConfig 接收并初始化 config
func HandleConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, _ := io.ReadAll(r.Body)
	SingleConfig = config.DecodeSingleConfig(body)

	clog.Infof("📥 [%s] Body: %s\n", "/Config", string(body))
	clog.Infof("📥 [%s] Config: %+v\n", "/Config", SingleConfig)

	resp := Response{
		Status:  "success",
		Message: fmt.Sprintf("%s received", "/config"),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
