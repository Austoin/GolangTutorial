# Go 网络编程详解

## 目录
- [1. 网络编程概述](#1-网络编程概述)
- [2. HTTP 服务器](#2-http-服务器)
- [3. TCP 编程](#3-tcp-编程)
- [4. UDP 编程](#4-udp-编程)
- [5. WebSocket](#5-websocket)
- [6. 高级主题](#6-高级主题)

---

## 1. 网络编程概述

### 1.1 Go 网络模型

Go 的网络编程基于 **CSP（Communicating Sequential Processes）** 模型，提供了简洁而强大的网络抽象。

```go
// net 包核心接口
type Conn interface {
    Read(b []byte) (n int, err error)
    Write(b []byte) (n int, err error)
    Close() error
    LocalAddr() Addr
    RemoteAddr() Addr
    SetDeadline(t time.Time) error
    // ...
}

type Listener interface {
    Accept() (Conn, error)
    Close() error
    Addr() Addr
}
```

### 1.2 网络分层模型

```
应用层    │ HTTP, FTP, SMTP, WebSocket
传输层    │ TCP, UDP
网络层    │ IP, ICMP
链路层    │ Ethernet, WiFi
```

---

## 2. HTTP 服务器

### 2.1 标准库 HTTP 服务器

Go 的 `net/http` 包提供了功能完整的 HTTP 服务器实现。

```go
package main

import (
    "fmt"
    "log"
    "net/http"
    "time"
)

// Handler 接口：所有 HTTP 处理器都需要实现这个接口
type Handler interface {
    ServeHTTP(w ResponseWriter, r *Request)
}

// http.HandlerFunc 类型：函数也可以作为处理器
func helloHandler(w http.ResponseWriter, r *http.Request) {
    // 1. 获取请求方法
    method := r.Method
    fmt.Printf("收到 %s 请求\n", method)

    // 2. 获取 URL 参数
    query := r.URL.Query()
    name := query.Get("name")
    if name == "" {
        name = "World"
    }

    // 3. 获取请求头
    userAgent := r.Header.Get("User-Agent")

    // 4. 获取请求体
    if method == "POST" {
        body := make([]byte, r.ContentLength)
        r.Body.Read(body)
        fmt.Printf("POST 数据: %s\n", string(body))
    }

    // 5. 设置响应头
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("X-Custom-Header", "Golang")

    // 6. 写入响应
    fmt.Fprintf(w, `{"message": "Hello, %s!", "method": "%s"}`, name, method)
    fmt.Fprintf(w, "\nUser-Agent: %s\n", userAgent)
}

// 使用 http.HandlerFunc 适配器
func timeHandler(w http.ResponseWriter, r *http.Request) {
    now := time.Now().Format("2006-01-02 15:04:05")
    fmt.Fprintf(w, "当前时间: %s\n", now)
}

// 自定义 Handler 结构体
type StatsHandler struct {
    requestCount int
}

func (s *StatsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    s.requestCount++
    w.Header().Set("Content-Type", "text/plain")
    fmt.Fprintf(w, "请求次数: %d\n", s.requestCount)
}

func main() {
    // 方法1：使用 http.HandleFunc 注册处理器
    http.HandleFunc("/hello", helloHandler)
    http.HandleFunc("/time", timeHandler)

    // 方法2：使用 http.Handle 注册自定义 Handler
    statsHandler := &StatsHandler{}
    http.Handle("/stats", statsHandler)

    // 方法3：使用 http.FileServer 提供静态文件服务
    http.Handle("/static/", http.StripPrefix("/static/", 
        http.FileServer(http.Dir("./public"))))

    // 方法4：使用 ServeMux 进行路由分组
    mux := http.NewServeMux()
    mux.HandleFunc("/api/hello", helloHandler)
    mux.HandleFunc("/api/time", timeHandler)

    // 复杂路由示例
    mux.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
        id := r.URL.Path[len("/users/"):]
        fmt.Fprintf(w, "用户 ID: %s\n", id)
    })

    // 中间件示例
    loggingMux := loggingMiddleware(mux)

    // 服务器配置
    server := &http.Server{
        Addr:         ":8080",
        Handler:      loggingMux,
        ReadTimeout:  10 * time.Second,
        WriteTimeout: 10 * time.Second,
        IdleTimeout:  60 * time.Second,
    }

    log.Println("服务器启动在 http://localhost:8080")
    log.Fatal(server.ListenAndServe())
}

// 中间件：日志记录
func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next.ServeHTTP(w, r)
        duration := time.Since(start)
        log.Printf("%s %s %s %v", r.Method, r.URL.Path, r.RemoteAddr, duration)
    })
}
```

### 2.2 HTTP 客户端

```go
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "time"
)

type User struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

func main() {
    // 1. 基础 GET 请求
    resp, err := http.Get("https://jsonplaceholder.typicode.com/users/1")
    if err != nil {
        fmt.Printf("请求失败: %v\n", err)
        return
    }
    defer resp.Body.Close()

    body, _ := io.ReadAll(resp.Body)
    fmt.Printf("GET 响应: %s\n", string(body))

    // 2. 带参数的 GET 请求
    url := "https://jsonplaceholder.typicode.com/posts"
    req, _ := http.NewRequest("GET", url, nil)
    q := req.URL.Query()
    q.Add("userId", "1")
    req.URL.RawQuery = q.Encode()

    client := &http.Client{Timeout: 5 * time.Second}
    resp, err = client.Do(req)
    if err != nil {
        fmt.Printf("请求失败: %v\n", err)
        return
    }
    defer resp.Body.Close()

    // 3. POST 请求（JSON）
    user := User{ID: 1, Name: "张三"}
    jsonData, _ := json.Marshal(user)

    resp, err = http.Post(
        "https://jsonplaceholder.typicode.com/users",
        "application/json",
        bytes.NewBuffer(jsonData),
    )
    if err != nil {
        fmt.Printf("POST 失败: %v\n", err)
        return
    }
    defer resp.Body.Close()

    // 4. 自定义请求（PUT, DELETE 等）
    req, _ = http.NewRequest("PUT", "https://jsonplaceholder.typicode.com/users/1", 
        bytes.NewBuffer(jsonData))
    req.Header.Set("Content-Type", "application/json")

    resp, err = client.Do(req)
    if err != nil {
        fmt.Printf("PUT 失败: %v\n", err)
        return
    }
    defer resp.Body.Close()

    // 5. 管理 Cookie
    jar, _ := cookieJar()
    cookieClient := &http.Client{Jar: jar, Timeout: 5 * time.Second}

    resp, err = cookieClient.Get("https://example.com")
    if err != nil {
        fmt.Printf("Cookie 请求失败: %v\n", err)
    }
}

func cookieJar() (*http.CookieJar, error) {
    return nil, nil // 使用默认 CookieJar
}
```

### 2.3 HTTP 中间件

```go
package main

import (
    "fmt"
    "log"
    "net/http"
    "time"
)

// Middleware 是中间件函数类型
type Middleware func(http.Handler) http.Handler

// 1. 日志中间件
func LoggerMiddleware() Middleware {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()
            next.ServeHTTP(w, r)
            log.Printf("%s %s %s %v", 
                r.Method, r.URL.Path, r.RemoteAddr, time.Since(start))
        })
    }
}

// 2. 认证中间件
func AuthMiddleware(token string) Middleware {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            auth := r.Header.Get("Authorization")
            if auth != "Bearer "+token {
                http.Error(w, "Unauthorized", http.StatusUnauthorized)
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}

// 3. CORS 中间件
func CORSMiddleware() Middleware {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            w.Header().Set("Access-Control-Allow-Origin", "*")
            w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
            w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
            
            if r.Method == "OPTIONS" {
                w.WriteHeader(http.StatusOK)
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}

// 4. 限流中间件（简单实现）
func RateLimitMiddleware(requestsPerSecond int) Middleware {
    limiter := time.NewTicker(time.Second / time.Duration(requestsPerSecond))
    
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            <-limiter.C // 等待令牌
            next.ServeHTTP(w, r)
        })
    }
}

// 5. 恢复中间件（Panic 恢复）
func RecoveryMiddleware() Middleware {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            defer func() {
                if err := recover(); err != nil {
                    log.Printf("Panic recovered: %v", err)
                    http.Error(w, "Internal Server Error", 
                        http.StatusInternalServerError)
                }
            }()
            next.ServeHTTP(w, r)
        })
    }
}

// 链式应用中间件
func ChainMiddleware(h http.Handler, mids ...Middleware) http.Handler {
    for i := len(mids) - 1; i >= 0; i-- {
        h = mids[i](h)
    }
    return h
}

func main() {
    handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello, World!\n")
    })

    // 链式应用中间件
    wrappedHandler := ChainMiddleware(handler,
        LoggerMiddleware(),
        AuthMiddleware("secret-token"),
        CORSMiddleware(),
        RecoveryMiddleware(),
    )

    http.Handle("/", wrappedHandler)
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

---

## 3. TCP 编程

### 3.1 TCP 服务器

```go
package main

import (
    "bufio"
    "fmt"
    "log"
    "net"
    "strings"
    "sync"
    "time"
)

type TCPServer struct {
    address string
    clients map[string]*Client
    mu      sync.RWMutex
}

type Client struct {
    conn     net.Conn
    name     string
    sendChan chan string
}

func NewTCPServer(address string) *TCPServer {
    return &TCPServer{
        address: address,
        clients: make(map[string]*Client),
    }
}

func (s *TCPServer) Start() error {
    listener, err := net.Listen("tcp", s.address)
    if err != nil {
        return err
    }
    defer listener.Close()
    
    log.Printf("TCP 服务器启动在 %s", s.address)

    for {
        conn, err := listener.Accept()
        if err != nil {
            log.Printf("接受连接失败: %v", err)
            continue
        }
        go s.handleConnection(conn)
    }
}

func (s *TCPServer) handleConnection(conn net.Conn) {
    defer conn.Close()
    
    client := &Client{
        conn:     conn,
        sendChan: make(chan string, 100),
    }

    // 读取客户端名称
    reader := bufio.NewReader(conn)
    fmt.Fprint(conn, "请输入您的名称: ")
    name, _ := reader.ReadString('\n')
    client.name = strings.TrimSpace(name)
    
    s.mu.Lock()
    s.clients[client.name] = client
    s.mu.Unlock()

    log.Printf("客户端 %s 已连接", client.name)

    // 广播用户加入消息
    s.broadcast(fmt.Sprintf("%s 加入聊天室\n", client.name))

    // 启动读取协程
    go s.readMessages(client)
    
    // 启动写入协程
    s.writeMessages(client)
}

func (s *TCPServer) readMessages(client *Client) {
    reader := bufio.NewReader(client.conn)
    for {
        message, err := reader.ReadString('\n')
        if err != nil {
            s.mu.Lock()
            delete(s.clients, client.name)
            s.mu.Unlock()
            s.broadcast(fmt.Sprintf("%s 离开聊天室\n", client.name))
            return
        }
        message = strings.TrimSpace(message)
        if message == "/quit" {
            client.conn.Close()
            return
        }
        s.broadcast(fmt.Sprintf("[%s]: %s\n", client.name, message))
    }
}

func (s *TCPServer) writeMessages(client *Client) {
    writer := bufio.NewWriter(client.conn)
    for message := range client.sendChan {
        _, err := writer.WriteString(message)
        if err != nil {
            client.conn.Close()
            return
        }
        writer.Flush()
    }
}

func (s *TCPServer) broadcast(message string) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    for _, client := range s.clients {
        select {
        case client.sendChan <- message:
        default:
            // 发送队列满，跳过
        }
    }
}

func main() {
    server := NewTCPServer(":8080")
    if err := server.Start(); err != nil {
        log.Fatal(err)
    }
}
```

### 3.2 TCP 客户端

```go
package main

import (
    "bufio"
    "fmt"
    "log"
    "net"
    "os"
    "time"
)

type TCPClient struct {
    conn    net.Conn
    address string
}

func NewTCPClient(address string) (*TCPClient, error) {
    conn, err := net.DialTimeout("tcp", address, 5*time.Second)
    if err != nil {
        return nil, err
    }
    return &TCPClient{conn: conn, address: address}, nil
}

func (c *TCPClient) Send(message string) error {
    _, err := c.conn.Write([]byte(message))
    return err
}

func (c *TCPClient) Receive() (string, error) {
    reader := bufio.NewReader(c.conn)
    return reader.ReadString('\n')
}

func (c *TCPClient) Close() error {
    return c.conn.Close()
}

func (c *TCPClient) Run() {
    // 启动读取协程
    done := make(chan bool)
    go func() {
        for {
            message, err := c.Receive()
            if err != nil {
                fmt.Println("连接已关闭")
                done <- true
                return
            }
            fmt.Print(message)
        }
    }()

    // 主线程处理用户输入
    reader := bufio.NewReader(os.Stdin)
    for {
        fmt.Print("请输入消息: ")
        input, _ := reader.ReadString('\n')
        if strings.TrimSpace(input) == "/quit" {
            c.Close()
            break
        }
        c.Send(input)
    }
    
    <-done
}

func main() {
    client, err := NewTCPClient("localhost:8080")
    if err != nil {
        log.Fatal("连接服务器失败:", err)
    }
    defer client.Close()

    fmt.Println("成功连接到服务器")
    client.Run()
}
```

### 3.3 TCP 心跳机制

```go
package main

import (
    "fmt"
    "log"
    "net"
    "sync"
    "time"
)

const (
    heartbeatInterval = 30 * time.Second
    heartbeatTimeout  = 90 * time.Second
)

type HeartbeatServer struct {
    listeners map[net.Conn]*ClientInfo
    mu        sync.RWMutex
}

type ClientInfo struct {
    lastSeen   time.Time
    isActive   bool
    heartbeat  chan bool
}

func NewHeartbeatServer() *HeartbeatServer {
    return &HeartbeatServer{
        listeners: make(map[net.Conn]*ClientInfo),
    }
}

func (s *HeartbeatServer) Start() error {
    listener, err := net.Listen("tcp", ":8080")
    if err != nil {
        return err
    }
    defer listener.Close()

    // 启动心跳检测协程
    go s.heartbeatChecker()

    for {
        conn, err := listener.Accept()
        if err != nil {
            log.Printf("接受连接失败: %v", err)
            continue
        }
        go s.handleConnection(conn)
    }
}

func (s *HeartbeatServer) handleConnection(conn net.Conn) {
    client := &ClientInfo{
        lastSeen:  time.Now(),
        isActive:  true,
        heartbeat: make(chan bool, 1),
    }

    s.mu.Lock()
    s.listeners[conn] = client
    s.mu.Unlock()

    buf := make([]byte, 1024)
    for {
        conn.SetReadDeadline(time.Now().Add(heartbeatTimeout))
        n, err := conn.Read(buf)
        if err != nil {
            s.mu.Lock()
            delete(s.listeners, conn)
            s.mu.Unlock()
            conn.Close()
            return
        }
        
        // 更新最后活跃时间
        client.lastSeen = time.Now()
        client.isActive = true
        
        // 处理心跳消息
        if string(buf[:n]) == "PING" {
            conn.Write([]byte("PONG"))
        }
    }
}

func (s *HeartbeatServer) heartbeatChecker() {
    ticker := time.NewTicker(heartbeatInterval)
    for range ticker.C {
        s.mu.RLock()
        for conn, client := range s.listeners {
            if time.Since(client.lastSeen) > heartbeatTimeout {
                log.Printf("客户端 %v 超时断开", conn.RemoteAddr())
                s.mu.RUnlock()
                s.mu.Lock()
                delete(s.listeners, conn)
                s.mu.Unlock()
                s.mu.RLock()
                conn.Close()
            }
        }
        s.mu.RUnlock()
    }
}

func main() {
    server := NewHeartbeatServer()
    if err := server.Start(); err != nil {
        log.Fatal(err)
    }
}
```

---

## 4. UDP 编程

### 4.1 UDP 服务器

```go
package main

import (
    "fmt"
    "log"
    "net"
)

func main() {
    // 解析 UDP 地址
    udpAddr, err := net.ResolveUDPAddr("udp", ":8080")
    if err != nil {
        log.Fatal(err)
    }

    // 创建 UDP 连接
    conn, err := net.ListenUDP("udp", udpAddr)
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    log.Println("UDP 服务器启动在 :8080")

    buffer := make([]byte, 1024)

    for {
        // 读取数据（同时获取客户端地址）
        n, addr, err := conn.ReadFromUDP(buffer)
        if err != nil {
            log.Printf("读取数据失败: %v", err)
            continue
        }

        fmt.Printf("收到来自 %v 的消息: %s\n", addr, string(buffer[:n]))

        // 发送响应
        response := fmt.Sprintf("服务器收到: %s", string(buffer[:n]))
        _, err = conn.WriteToUDP([]byte(response), addr)
        if err != nil {
            log.Printf("发送响应失败: %v", err)
        }
    }
}
```

### 4.2 UDP 客户端

```go
package main

import (
    "fmt"
    "log"
    "net"
    "time"
)

func main() {
    // 解析服务器地址
    serverAddr, err := net.ResolveUDPAddr("udp", "localhost:8080")
    if err != nil {
        log.Fatal(err)
    }

    // 创建 UDP 连接
    conn, err := net.DialUDP("udp", nil, serverAddr)
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    // 设置超时
    conn.SetReadDeadline(time.Now().Add(5 * time.Second))

    // 发送数据
    message := "Hello, UDP Server!"
    _, err = conn.Write([]byte(message))
    if err != nil {
        log.Fatal("发送数据失败:", err)
    }

    fmt.Println("发送消息:", message)

    // 接收响应
    buffer := make([]byte, 1024)
    n, _, err := conn.ReadFromUDP(buffer)
    if err != nil {
        log.Fatal("接收数据失败:", err)
    }

    fmt.Printf("收到服务器响应: %s\n", string(buffer[:n]))
}
```

### 4.3 UDP 广播

```go
package main

import (
    "fmt"
    "log"
    "net"
    "time"
)

func main() {
    // 广播地址
    broadcastAddr, err := net.ResolveUDPAddr("udp", "255.255.255.255:8080")
    if err != nil {
        log.Fatal(err)
    }

    // 创建 UDP 连接
    conn, err := net.DialUDP("udp", nil, broadcastAddr)
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    // 启用广播
    conn.SetWriteBuffer(1024 * 1024) // 1MB 缓冲区

    // 定时广播
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()

    for range ticker.C {
        message := fmt.Sprintf("服务发现广播 %s", time.Now().Format("15:04:05"))
        _, err := conn.Write([]byte(message))
        if err != nil {
            log.Printf("广播失败: %v", err)
        }
        fmt.Println("广播:", message)
    }
}
```

---

## 5. WebSocket

### 5.1 WebSocket 服务器

```go
package main

import (
    "fmt"
    "log"
    "net/http"
    "sync"

    "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    // 允许跨域
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

type Client struct {
    conn *websocket.Conn
    send chan []byte
}

type Hub struct {
    clients    map[*Client]bool
    broadcast  chan []byte
    register   chan *Client
    unregister chan *Client
    mu         sync.RWMutex
}

func NewHub() *Hub {
    return &Hub{
        clients:    make(map[*Client]bool),
        broadcast:  make(chan []byte, 256),
        register:   make(chan *Client),
        unregister: make(chan *Client),
    }
}

func (h *Hub) Run() {
    for {
        select {
        case client := <-h.register:
            h.mu.Lock()
            h.clients[client] = true
            h.mu.Unlock()

        case client := <-h.unregister:
            h.mu.Lock()
            if _, ok := h.clients[client]; ok {
                delete(h.clients, client)
                close(client.send)
            }
            h.mu.Unlock()

        case message := <-h.broadcast:
            h.mu.RLock()
            for client := range h.clients {
                select {
                case client.send <- message:
                default:
                    close(client.send)
                    delete(h.clients, client)
                }
            }
            h.mu.RUnlock()
        }
    }
}

func (h *Hub) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Printf("升级连接失败: %v", err)
        return
    }

    client := &Client{
        conn: conn,
        send: make(chan []byte, 256),
    }

    h.register <- client

    // 启动读取协程
    go func() {
        for {
            _, message, err := conn.ReadMessage()
            if err != nil {
                h.unregister <- client
                conn.Close()
                return
            }
            fmt.Printf("收到消息: %s\n", string(message))
            // 广播消息
            h.broadcast <- message
        }
    }()

    // 启动写入协程
    go func() {
        for message := range client.send {
            if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
                h.unregister <- client
                return
            }
        }
    }()
}

func main() {
    hub := NewHub()
    go hub.Run()

    http.HandleFunc("/ws", hub.HandleWebSocket)

    log.Println("WebSocket 服务器启动在 ws://localhost:8080/ws")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

### 5.2 WebSocket 客户端

```go
package main

import (
    "fmt"
    "log"
    "os"
    "os/signal"
    "syscall"

    "github.com/gorilla/websocket"
)

func main() {
    interrupt := make(chan os.Signal, 1)
    signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

    conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/ws", nil)
    if err != nil {
        log.Fatal("连接失败:", err)
    }
    defer conn.Close()

    done := make(chan struct{})

    // 接收消息
    go func() {
        defer close(done)
        for {
            _, message, err := conn.ReadMessage()
            if err != nil {
                log.Printf("读取消息失败: %v", err)
                return
            }
            fmt.Printf("收到消息: %s\n", string(message))
        }
    }()

    // 发送消息
    ticker := NewTicker(3 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case t := <-ticker.C:
            err := conn.WriteMessage(websocket.TextMessage, 
                []byte(fmt.Sprintf("消息 %s", t.Format("15:04:05"))))
            if err != nil {
                log.Printf("发送消息失败: %v", err)
                return
            }

        case <-interrupt:
            log.Println("中断连接")
            conn.WriteMessage(websocket.CloseMessage, 
                websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
            return

        case <-done:
            return
        }
    }
}

func NewTicker(d time.Duration) *time.Ticker {
    return time.NewTicker(d)
}
```

---

## 6. 高级主题

### 6.1 连接池

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

type Pool struct {
    mu       sync.Mutex
    conns    chan net.Conn
    factory  func() (net.Conn, error)
    maxIdle  int
    maxOpen  int
    openConn int
}

func NewPool(factory func() (net.Conn, error), maxIdle, maxOpen int) *Pool {
    return &Pool{
        conns:   make(chan net.Conn, maxIdle),
        factory: factory,
        maxIdle: maxIdle,
        maxOpen: maxOpen,
    }
}

func (p *Pool) Get() (net.Conn, error) {
    select {
    case conn := <-p.conns:
        return conn, nil
    default:
        p.mu.Lock()
        if p.openConn >= p.maxOpen {
            p.mu.Unlock()
            return nil, fmt达到最大连接数")
        }
        p.openConn++
        p.mu.Unlock()
        return p.factory()
    }
}

func (p *Pool) Put(conn net.Conn) error {
    select {
    case p.conns <- conn:
        return nil
    default:
        // 连接池已满，关闭连接
        p.mu.Lock()
        p.openConn--
        p.mu.Unlock()
        return conn.Close()
    }
}

func (p *Pool) Close() error {
    p.mu.Lock()
    defer p.mu.Unlock()
    
    for i := 0; i < len(p.conns); i++ {
        <-p.conns
    }
    return nil
}
```

### 6.2 TLS/SSL 加密

```go
package main

import (
    "crypto/rand"
    "crypto/rsa"
    "crypto/tls"
    "crypto/x509"
    "crypto/x509/pkix"
    "encoding/pem"
    "fmt"
    "math/big"
    "net"
    "time"
)

// 生成自签名证书
func generateCert() (tls.Certificate, error) {
    // 生成私钥
    privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
    if err != nil {
        return tls.Certificate{}, err
    }

    // 创建证书模板
    template := x509.Certificate{
        SerialNumber: big.NewInt(1),
        Subject: pkix.Name{
            Organization: []string{"Example Corp"},
        },
        NotBefore:             time.Now(),
        NotAfter:              time.Now().Add(365 * 24 * time.Hour),
        KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
        ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
        BasicConstraintsValid: true,
        IPAddresses:           []net.IP{net.ParseIP("127.0.0.1")},
    }

    // 创建证书
    certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, 
        &privateKey.PublicKey, privateKey)
    if err != nil {
        return tls.Certificate{}, err
    }

    // PEM 编码
    certPEM := pem.EncodeToMemory(&pem.Block{
        Type:  "CERTIFICATE",
        Bytes: certDER,
    })
    keyPEM := pem.EncodeToMemory(&pem.Block{
        Type:  "RSA PRIVATE KEY",
        Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
    })

    return tls.X509KeyPair(certPEM, keyPEM)
}

// TLS 服务器
func startTLSServer() error {
    cert, err := generateCert()
    if err != nil {
        return err
    }

    config := &tls.Config{
        Certificates: []tls.Certificate{cert},
    }

    listener, err := tls.Listen("tcp", ":8443", config)
    if err != nil {
        return err
    }
    defer listener.Close()

    fmt.Println("TLS 服务器启动在 https://localhost:8443")
    for {
        conn, err := listener.Accept()
        if err != nil {
            fmt.Printf("接受连接失败: %v\n", err)
            continue
        }
        // 处理连接...
        conn.Close()
    }
}

func main() {
    startTLSServer()
}
```

---

## 最佳实践

1. **资源管理**：始终使用 `defer` 关闭连接
2. **超时设置**：为所有网络操作设置合理的超时
3. **并发安全**：使用适当的锁或 channel 保护共享资源
4. **错误处理**：区分可恢复和不可恢复的错误
5. **日志记录**：记录关键操作和错误信息
6. **监控指标**：监控连接数、请求延迟、错误率等

## 性能优化建议

1. **连接复用**：使用 HTTP/2 或连接池
2. **缓冲区调优**：根据场景调整读写缓冲区大小
3. **零拷贝**：使用 `io.Copy` 进行高效数据传输
4. **批量处理**：合并小请求减少网络往返
5. **压缩**：对大响应启用 gzip 压缩
