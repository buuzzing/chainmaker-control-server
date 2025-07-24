/*
 * @Author: Justin
 * @Date: 2025-07-21 05:59:00
 * @filename:
 * @version:
 * @Description:
 * @LastEditTime: 2025-07-21 08:51:59
 */

package httpserver

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"test_server/config"

	clog "github.com/kpango/glg"
)

var SingleConfig *config.SingleConfig

// 1. 接收 初始化config
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

// 2. 启动链
func HandleStartChain(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	body, _ := io.ReadAll(r.Body)

	clog.Infof("📥 [%s] Body: %s\n", "/startChain", string(body))

	resp := Response{
		Status:  "success",
		Message: fmt.Sprintf("%s received", "/startChain"),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// 3. 部署合约 初始化合约
func HandleSetupContracts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	body, _ := io.ReadAll(r.Body)

	clog.Infof("📥 [%s] Body: %s\n", "/setupContracts", string(body))

	resp := Response{
		Status:  "success",
		Message: fmt.Sprintf("%s received", "/setupContracts"),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// 4. 启动 relayer
func HandleStartRelayer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	body, _ := io.ReadAll(r.Body)

	clog.Infof("📥 [%s] Body: %s\n", "/startRelayer", string(body))

	resp := Response{
		Status:  "success",
		Message: fmt.Sprintf("%s received", "/startRelayer"),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// 5. 向 verify server 注册 relayer
func HandleSubscribeRelayer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	body, _ := io.ReadAll(r.Body)

	clog.Infof("📥 [%s] Body: %s\n", "/subscribeRelayer", string(body))

	resp := Response{
		Status:  "success",
		Message: fmt.Sprintf("%s received", "/subscribeRelayer"),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
