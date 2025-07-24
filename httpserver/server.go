package httpserver

import (
	"net/http"

	clog "github.com/kpango/glg"
)

// Response 自定义响应结构体
type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// 定义路由处理函数类型
type HandlerFunc func(w http.ResponseWriter, r *http.Request)

// HttpEngine = 自定义 HTTP Server
type HttpEngine struct {
	Port   string
	router map[string]HandlerFunc
}

// 创建新实例
func NewHttpEngine(port string) *HttpEngine {
	return &HttpEngine{
		Port:   port,
		router: make(map[string]HandlerFunc),
	}
}

// 注册路由
func (e *HttpEngine) Handle(path string, handler HandlerFunc) {
	e.router[path] = handler
	clog.Infof("注册路由，%s", path)
}

// 实现 ServeHTTP 方法，作为 http.Handler 使用
func (e *HttpEngine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler, ok := e.router[r.URL.Path]
	if !ok {
		http.NotFound(w, r)
		return
	}
	handler(w, r)
}

// 启动 HTTP 服务
func (e *HttpEngine) Start() error {
	clog.Infof("🚀 Listening on :%s\n", e.Port)
	return http.ListenAndServe(":"+e.Port, e)
}
