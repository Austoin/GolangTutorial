// networking/network_http_server.go
// HTTP 服务器示例 - 详细注释版

package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// ====== HTTP 服务器基础 ======

// http.Handler 接口是 Go 语言中处理 HTTP 请求的核心接口
// 任何实现了 ServeHTTP(ResponseWriter, *Request) 方法的类型都可以作为 HTTP 处理器
type HelloHandler struct {
	// 可以在这里添加结构体字段，如配置信息等
	name string
}

// ServeHTTP 方法实现了 http.Handler 接口
// 当有请求到达时，Go 的 HTTP 服务器会自动调用这个方法
// 参数说明：
//   - w: 用于写入响应内容的 ResponseWriter
//   - r: 包含请求信息的 *Request 对象
func (h *HelloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 1. 设置响应头 - Content-Type 告诉客户端响应内容的 MIME 类型
	//    浏览器会根据这个值来正确渲染内容
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	// 2. 写入状态码 - 200 OK 表示请求成功
	//    常见状态码：200(成功), 404(未找到), 500(服务器错误)等
	w.WriteHeader(http.StatusOK)

	// 3. 写入响应内容 - 使用 fmt.Fprintf 格式化输出到响应流
	//    也可以使用 w.Write([]byte("内容")) 直接写入字节
	fmt.Fprintf(w, "Hello, %s!\n", h.name)
}

// ====== 使用 http.HandleFunc 简化路由 ======

// 处理根路径 "/" 的请求
// http.HandleFunc 是最常用的路由注册方式，它接受一个函数作为处理器
// 函数签名必须为：func(http.ResponseWriter, *Request)
func homeHandler(w http.ResponseWriter, r *http.Request) {
	// 检查请求路径 - 确保只处理正确的路径
	// 可以使用 r.URL.Path 来获取请求的路径
	if r.URL.Path != "/" {
		// 返回 404 状态码，表示资源未找到
		http.NotFound(w, r)
		return
	}

	// 使用 html/template 进行 HTML 渲染（更安全，防止 XSS 攻击）
	// 简单示例中可以直接写入字符串
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, `
		<!DOCTYPE html>
		<html>
		<head><title>Go HTTP Server</title></head>
		<body>
			<h1>欢迎使用 Go 语言 HTTP 服务器</h1>
			<p>当前时间：%s</p>
			<ul>
				<li><a href="/hello?name=World">/hello?name=World</a></li>
				<li><a href="/time">/time</a></li>
			</ul>
		</body>
		</html>
	`, time.Now().Format("2006-01-02 15:04:05"))
}

// 处理 /hello 路径的请求
func helloHandler(w http.ResponseWriter, r *http.Request) {
	// 获取查询参数 - r.URL.Query() 返回 URL 查询参数的键值对
	// 使用 Get() 方法可以安全地获取参数值，如果参数不存在则返回空字符串
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "Guest"
	}

	// 设置响应头
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// 返回 JSON 格式的响应
	// fmt.Fprintf 可以直接写入 JSON 字符串
	// 更推荐使用 encoding/json 包来序列化结构体
	fmt.Fprintf(w, `{"message": "Hello, %s!", "time": "%s"}`,
		name, time.Now().Format(time.RFC3339))
}

// 处理 /time 路径的请求
func timeHandler(w http.ResponseWriter, r *http.Request) {
	// 获取当前时间并格式化
	// Go 的时间格式化使用特定的参考时间：Mon Jan 2 15:04:05 MST 2006
	// 这个日期的特殊之处在于它按顺序展示了时间的各个部分
	now := time.Now().Format("2006年01月02日 15:04:05")

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "当前服务器时间：%s\n", now)
}

// ====== 中间件示例 ======

// LoggerMiddleware 创建日志中间件
// 中间件是一个接收 http.Handler 返回另一个 http.Handler 的函数
// 中间件可以在请求处理前后执行额外的逻辑，如日志记录、认证等
func LoggerMiddleware(next http.Handler) http.Handler {
	// 返回的 HandlerFunc 是一个适配器，将函数转换为 http.Handler
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 请求到达前的处理
		start := time.Now()
		log.Printf("收到请求: %s %s", r.Method, r.URL.Path)

		// 创建自定义的 ResponseWriter 来记录响应状态码
		// 因为 http.ResponseWriter 的 WriteHeader 方法是延迟调用的
		lrw := &loggingResponseWriter{ResponseWriter: w, statusCode: 200}

		// 调用下一个处理器（请求处理）
		next.ServeHTTP(lrw, r)

		后的 // 请求处理处理
		duration := time.Since(start)
		log.Printf("请求完成: %s %s - %d - %v",
			r.Method, r.URL.Path, lrw.statusCode, duration)
	})
}

// loggingResponseWriter 包装 http.ResponseWriter 以记录响应状态码
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader 覆盖原始方法，记录实际的状态码
// 如果多次调用 WriteHeader，只有第一次会生效
func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

// ====== 静态文件服务 ======

// 使用 http.FileServer 提供静态文件服务
// FileServer 接受一个 FileSystem 参数，返回一个 http.Handler
// 常用的文件系统实现：
//   - http.Dir: 将字符串路径转换为 FileSystem
//   - http.Dir("."): 当前目录
//   - http.FS: 将 embed.FS 转换为 FileSystem（Go 1.16+）
func staticFileHandler() http.Handler {
	// http.StripPrefix 用于去除请求路径的前缀
	// 这样 /static/files/logo.png 会被映射到 ./files/logo.png
	return http.StripPrefix("/static/", http.FileServer(http.Dir("./static")))
}

// ====== 主函数 - 服务器入口 ======

func main() {
	// 1. 注册路由处理器
	// http.Handle 和 http.HandleFunc 用于注册路由
	// 路由规则：
	//   - 精确匹配：如 "/hello" 只匹配 /hello
	//   - 路径前缀：如 "/static/" 匹配所有以 /static/ 开头的路径

	// 使用自定义处理器
	helloHandler := &HelloHandler{name: "Guest"}
	http.Handle("/custom", helloHandler)

	// 使用函数处理器
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/time", timeHandler)

	// 注册静态文件服务
	// 所有 /static/* 的请求都会从 ./static 目录提供文件
	http.Handle("/static/", staticFileHandler())

	// 2. 应用中间件
	// 使用 http.TimeoutHandler 添加超时控制
	// 这可以防止慢请求占用过多服务器资源
	// 超时处理会自动返回 503 Service Unavailable 响应
	// wrappedHandler := LoggerMiddleware(http.DefaultServeMux)

	// 3. 配置服务器
	// http.Server 结构体用于配置 HTTP 服务器
	server := &http.Server{
		Addr:         ":8080",           // 监听地址和端口，格式为 host:port
		ReadTimeout:  10 * time.Second,  // 读取请求的超时时间
		WriteTimeout: 10 * time.Second,  // 写入响应的超时时间
		IdleTimeout:  120 * time.Second, // 空闲连接的最大存活时间
	}

	// 4. 启动服务器
	// ListenAndServe 会阻塞当前 Goroutine
	// 直到服务器因错误或调用 Shutdown 而停止
	log.Println("服务器启动在 http://localhost:8080")
	log.Println("按 Ctrl+C 停止服务器")

	// 使用 TLS（HTTPS）启动服务器
	// server.ListenAndServeTLS("cert.pem", "key.pem")

	// 普通 HTTP 启动
	if err := server.ListenAndServe(); err != nil {
		// 处理服务器启动错误
		// 常见错误：端口已被占用、权限不足等
		log.Fatalf("服务器启动失败: %v", err)
	}
}

// ====== 进阶：使用 ServeMux ======

// http.ServeMux 是 Go 标准库中的多路复用器（路由）
// 它提供了更灵活的路由管理
func createCustomMux() *http.ServeMux {
	mux := http.NewServeMux()

	// 注册处理器
	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/hello", helloHandler)
	mux.HandleFunc("/time", timeHandler)

	// 添加自定义处理器
	mux.Handle("/custom", &HelloHandler{name: "Custom"})

	// 添加静态文件服务
	mux.Handle("/static/", http.StripPrefix("/static/",
		http.FileServer(http.Dir("./static"))))

	return mux
}
