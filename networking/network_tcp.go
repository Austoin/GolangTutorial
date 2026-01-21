// networking/network_tcp.go
// TCP 服务器与客户端示例 - 详细注释版

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

// ====== TCP 服务器基础 ======

// TCPServer 结构体表示一个 TCP 服务器
// 包含监听地址和连接管理等
type TCPServer struct {
	address  string         // 监听地址，如 ":8080"
	listener net.Listener   // 监听器，用于接受连接
	wg       sync.WaitGroup // 用于优雅关闭
}

// NewTCPServer 创建新的 TCP 服务器实例
func NewTCPServer(address string) *TCPServer {
	return &TCPServer{
		address: address,
	}
}

// Start 启动 TCP 服务器
// 这个方法会阻塞，直到服务器关闭
func (s *TCPServer) Start() error {
	// 1. 创建监听器
	// net.Listen 用于创建 TCP 监听器
	// 第一个参数是网络类型（"tcp"、"tcp4"、"tcp6"等）
	// 第二个参数是监听地址，格式为 host:port
	var err error
	s.listener, err = net.Listen("tcp", s.address)
	if err != nil {
		return fmt.Errorf("创建监听器失败: %w", err)
	}

	log.Printf("TCP 服务器启动，监听地址: %s", s.address)

	// 2. 接受连接循环
	// Accept 方法会阻塞，直到有新的连接到来
	// 返回的 net.Conn 表示一个连接，可以进行读写操作
	for {
		// 检查服务器是否已关闭
		if s.listener == nil {
			break
		}

		// 接受新连接
		conn, err := s.listener.Accept()
		if err != nil {
			// 如果是临时错误，继续接受连接
			// 如果是严重错误，可能需要停止服务器
			if netErr, ok := err.(net.Error); ok && netErr.Temporary() {
				log.Printf("临时错误: %v", err)
				continue
			}
			return fmt.Errorf("接受连接失败: %w", err)
		}

		// 3. 处理连接（使用 Goroutine 并发处理）
		// 每个连接独立处理，不会阻塞其他连接
		s.wg.Add(1)
		go s.handleConnection(conn)
	}

	return nil
}

// handleConnection 处理单个客户端连接
// conn 参数是客户端连接
func (s *TCPServer) handleConnection(conn net.Conn) {
	// 确保连接最后关闭
	defer func() {
		conn.Close()
		s.wg.Done()
		log.Printf("客户端断开: %s", conn.RemoteAddr().String())
	}()

	log.Printf("新客户端连接: %s", conn.RemoteAddr().String())

	// 4. 创建缓冲区用于读取数据
	// bufio.Scanner 提供了方便的数据读取方式
	// 默认按行分割，最大 64K
	scanner := bufio.NewScanner(conn)

	// 可以设置自定义的分割函数和缓冲区大小
	// scanner.Split(bufio.ScanLines)
	// scanner.Buffer(make([]byte, 1024), 1024*1024) // 1MB 缓冲区

	for scanner.Scan() {
		// 读取一行数据
		message := scanner.Text()
		log.Printf("收到消息: %s", message)

		// 5. 处理消息并生成响应
		response := s.processMessage(message)

		// 6. 发送响应
		// 写入数据时使用 bufio.Writer 提供缓冲
		writer := bufio.NewWriter(conn)
		fmt.Fprintf(writer, "%s\n", response)
		writer.Flush() // 确保数据发送出去
	}

	// 检查扫描错误
	if err := scanner.Err(); err != nil {
		log.Printf("读取错误: %v", err)
	}
}

// processMessage 处理客户端消息并返回响应
func (s *TCPServer) processMessage(message string) string {
	message = strings.TrimSpace(message)

	// 根据消息内容生成不同的响应
	switch message {
	case "ping":
		return "pong"
	case "time":
		return time.Now().Format("2006-01-02 15:04:05")
	case "date":
		return time.Now().Format("2006-01-02")
	case "quit":
		return "BYE"
	case "echo":
		return "请提供要回显的内容，格式: echo:<内容>"
	default:
		if strings.HasPrefix(message, "echo:") {
			return strings.TrimPrefix(message, "echo:")
		}
		return fmt.Sprintf("未知命令: %s", message)
	}
}

// Shutdown 优雅关闭服务器
// 等待所有正在处理的连接完成
func (s *TCPServer) Shutdown() error {
	log.Println("正在关闭服务器...")

	// 关闭监听器，停止接受新连接
	if s.listener != nil {
		s.listener.Close()
		s.listener = nil
	}

	// 等待所有连接处理完成
	s.wg.Wait()

	log.Println("服务器已关闭")
	return nil
}

// ====== TCP 客户端示例 ======

// TCPClient 表示 TCP 客户端
type TCPClient struct {
	address string   // 服务器地址
	conn    net.Conn // 连接
}

// NewTCPClient 创建新的 TCP 客户端
func NewTCPClient(address string) (*TCPClient, error) {
	// 1. 连接到服务器
	// net.Dial 会阻塞，直到连接建立或超时
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("连接服务器失败: %w", err)
	}

	return &TCPClient{
		address: address,
		conn:    conn,
	}, nil
}

// Send 发送消息并接收响应
func (c *TCPClient) Send(message string) (string, error) {
	// 1. 发送消息
	// 使用 bufio.Writer 写入数据
	writer := bufio.NewWriter(c.conn)
	fmt.Fprintf(writer, "%s\n", message)
	writer.Flush()

	// 2. 接收响应
	// 使用 bufio.Scanner 读取响应
	scanner := bufio.NewScanner(c.conn)
	if !scanner.Scan() {
		return "", scanner.Err()
	}

	return scanner.Text(), nil
}

// Close 关闭客户端连接
func (c *TCPClient) Close() error {
	return c.conn.Close()
}

// ====== 高级：使用 net.Pipe 进行测试 ======

// PipeServer 使用 net.Pipe 创建内存中的服务器
// 适用于单元测试，不需要真实的网络连接
func PipeServerExample() {
	// 创建一对连接的 pipe
	// 一端作为服务器，一端作为客户端
	serverPipe, clientPipe := net.Pipe()

	// 服务器端
	go func() {
		scanner := bufio.NewScanner(serverPipe)
		for scanner.Scan() {
			msg := scanner.Text()
			fmt.Printf("服务器收到: %s\n", msg)
			fmt.Fprintf(serverPipe, "服务器响应: %s\n", msg)
		}
	}()

	// 客户端
	scanner := bufio.NewScanner(clientPipe)
	fmt.Fprintf(clientPipe, "客户端消息\n")

	if scanner.Scan() {
		fmt.Printf("客户端收到: %s\n", scanner.Text())
	}
}

// ====== 高级：聊天服务器 ======

// ChatServer 实现一个简单的多人聊天服务器
type ChatServer struct {
	clients   map[net.Conn]string // 客户端连接 -> 用户名
	mu        sync.RWMutex        // 保护 clients 映射
	broadcast chan string         // 广播消息通道
}

// NewChatServer 创建新的聊天服务器
func NewChatServer() *ChatServer {
	return &ChatServer{
		clients:   make(map[net.Conn]string),
		broadcast: make(chan string, 10),
	}
}

// Start 启动聊天服务器
func (cs *ChatServer) Start(address string) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	// 启动广播处理协程
	go cs.handleBroadcast()

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		go cs.handleChatClient(conn)
	}
}

// handleBroadcast 处理广播消息
func (cs *ChatServer) handleBroadcast() {
	for msg := range cs.broadcast {
		cs.mu.RLock()
		defer cs.mu.RUnlock()

		for conn := range cs.clients {
			fmt.Fprintf(conn, "%s\n", msg)
		}
	}
}

// handleChatClient 处理聊天客户端
func (cs *ChatServer) handleChatClient(conn net.Conn) {
	defer conn.Close()

	// 读取用户名
	scanner := bufio.NewScanner(conn)
	if !scanner.Scan() {
		return
	}
	username := scanner.Text()

	// 注册客户端
	cs.mu.Lock()
	cs.clients[conn] = username
	cs.mu.Unlock()

	// 广播用户加入
	cs.broadcast <- fmt.Sprintf("[系统] %s 加入聊天", username)

	// 处理消息
	for scanner.Scan() {
		msg := scanner.Text()
		if msg == "/quit" {
			break
		}
		cs.broadcast <- fmt.Sprintf("[%s] %s", username, msg)
	}

	// 客户端离开
	cs.mu.Lock()
	delete(cs.clients, conn)
	cs.mu.Unlock()
	cs.broadcast <- fmt.Sprintf("[系统] %s 离开聊天", username)
}

// ====== 主函数 ======

func main() {
	// 示例 1: 简单回显服务器
	fmt.Println("=== 简单 TCP 回显服务器 ===")

	server := NewTCPServer(":8080")

	// 在 Goroutine 中启动服务器
	go func() {
		if err := server.Start(); err != nil {
			log.Printf("服务器错误: %v", err)
		}
	}()

	// 等待服务器启动
	time.Sleep(100 * time.Millisecond)

	// 测试客户端
	client, err := NewTCPClient("localhost:8080")
	if err != nil {
		log.Fatalf("创建客户端失败: %v", err)
	}
	defer client.Close()

	// 发送测试消息
	tests := []string{"ping", "time", "echo:Hello World", "quit"}
	for _, test := range tests {
		response, err := client.Send(test)
		if err != nil {
			log.Printf("发送失败: %v", err)
			continue
		}
		fmt.Printf("发送: %s -> 收到: %s\n", test, response)
	}

	// 关闭服务器
	server.Shutdown()

	// 示例 2: 聊天服务器（需要多个客户端测试）
	// chatServer := NewChatServer()
	// go chatServer.Start(":8081")
}
