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

// HandleStartChain 启动链
func HandleStartChain(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	body, _ := io.ReadAll(r.Body)

	clog.Infof("📥 [%s] Body: %s\n", "/startChain", string(body))

	// TODO: 启动链逻辑

	resp := Response{
		Status:  "success",
		Message: fmt.Sprintf("%s received", "/startChain"),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// HandleSetupContracts 部署合约与初始化合约
func HandleSetupContracts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	body, _ := io.ReadAll(r.Body)

	clog.Infof("📥 [%s] Body: %s\n", "/setupContracts", string(body))

	// TODO: 部署合约与初始化合约逻辑

	resp := Response{
		Status:  "success",
		Message: fmt.Sprintf("%s received", "/setupContracts"),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// HandleStartRelayer 启动 relayer
func HandleStartRelayer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	body, _ := io.ReadAll(r.Body)

	clog.Infof("📥 [%s] Body: %s\n", "/startRelayer", string(body))

	// TODO: 启动 relayer 逻辑

	resp := Response{
		Status:  "success",
		Message: fmt.Sprintf("%s received", "/startRelayer"),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// HandleSubscribeRelayer 向 verify server 注册 relayer
func HandleSubscribeRelayer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	body, _ := io.ReadAll(r.Body)

	clog.Infof("📥 [%s] Body: %s\n", "/subscribeRelayer", string(body))

	// TODO: 向 verify server 注册 relayer 逻辑

	resp := Response{
		Status:  "success",
		Message: fmt.Sprintf("%s received", "/subscribeRelayer"),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
