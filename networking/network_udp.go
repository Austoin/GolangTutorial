// networking/network_udp.go
// UDP 服务器与客户端示例 - 详细注释版

package main

import (
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

// ====== UDP 协议特点 ======
/*
UDP (User Datagram Protocol) 是一种无连接的传输层协议。

主要特点：
1. 无连接 - 发送数据前不需要建立连接
2. 不可靠 - 不保证数据包的顺序和到达
3. 速度快 - 没有连接建立和确认的开销
4. 数据报 - 数据以独立的数据报形式发送
5. 有大小限制 - 单个数据报最大约 64KB（通常限制在 512 字节以内）

适用场景：
- DNS 查询
- 视频流、语音通话
- 实时游戏
- IoT 设备通信
- 简单请求-响应场景
*/

// ====== UDP 服务器基础 ======

// UDPServer 结构体表示一个 UDP 服务器
type UDPServer struct {
	address string
	conn    *net.UDPConn
	wg      sync.WaitGroup
}

// NewUDPServer 创建新的 UDP 服务器
func NewUDPServer(address string) *UDPServer {
	return &UDPServer{
		address: address,
	}
}

// Start 启动 UDP 服务器
func (s *UDPServer) Start() error {
	// 1. 解析 UDP 地址
	// net.ResolveUDPAddr 用于解析 UDP 地址
	// 参数：网络类型("udp")、地址字符串
	udpAddr, err := net.ResolveUDPAddr("udp", s.address)
	if err != nil {
		return fmt.Errorf("解析地址失败: %w", err)
	}

	// 2. 创建 UDP 监听器
	// 与 TCP 不同，UDP 使用 ListenUDP 而不是 Listen
	s.conn, err = net.ListenUDP("udp", udpAddr)
	if err != nil {
		return fmt.Errorf("创建监听器失败: %w", err)
	}

	log.Printf("UDP 服务器启动，监听地址: %s", s.address)

	// 设置读写超时（可选）
	// 这可以防止服务器在无数据时永久阻塞
	s.conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	s.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

	// 3. 处理数据报循环
	buf := make([]byte, 1024) // 数据缓冲区

	for {
		// 4. 读取数据报
		// ReadFromUDP 返回：读取的字节数、发送方地址、错误
		// 这个方法会阻塞，直到收到数据或超时
		n, addr, err := s.conn.ReadFromUDP(buf)
		if err != nil {
			// 检查是否是超时错误
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				// 继续循环，检查服务器是否关闭
				continue
			}
			if netErr, ok := err.(net.Error); ok && netErr.Temporary() {
				log.Printf("临时错误: %v", err)
				continue
			}
			return fmt.Errorf("读取数据失败: %w", err)
		}

		// 5. 处理数据
		data := string(buf[:n])
		log.Printf("收到来自 %s 的消息: %s", addr.String(), data)

		// 6. 生成响应
		response := s.processMessage(data)

		// 7. 发送响应
		// WriteToUDP 将数据发送到指定地址
		_, err = s.conn.WriteToUDP([]byte(response), addr)
		if err != nil {
			log.Printf("发送响应失败: %v", err)
		}
	}
}

// processMessage 处理消息并返回响应
func (s *UDPServer) processMessage(message string) string {
	message = strings.TrimSpace(message)

	switch message {
	case "ping":
		return "pong"
	case "time":
		return time.Now().Format("2006-01-02 15:04:05")
	case "date":
		return time.Now().Format("2006-01-02")
	default:
		return fmt.Sprintf("Echo: %s", message)
	}
}

// Close 关闭 UDP 连接
func (s *UDPServer) Close() error {
	if s.conn != nil {
		return s.conn.Close()
	}
	return nil
}

// ====== UDP 客户端基础 ======

// UDPClient 表示 UDP 客户端
type UDPClient struct {
	address string
	conn    *net.UDPConn
}

// NewUDPClient 创建新的 UDP 客户端
func NewUDPClient(address string) (*UDPClient, error) {
	// 1. 解析服务器地址
	serverAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return nil, fmt.Errorf("解析服务器地址失败: %w", err)
	}

	// 2. 建立 UDP 连接
	// DialUDP 建立一个 "连接" 到服务器
	// 这样可以避免每次发送都指定地址
	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		return nil, fmt.Errorf("创建 UDP 连接失败: %w", err)
	}

	// 设置超时
	conn.SetDeadline(time.Now().Add(5 * time.Second))
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	conn.SetWriteDeadline(time.Now().Add(5 * time.Second))

	return &UDPClient{
		address: address,
		conn:    conn,
	}, nil
}

// Send 发送消息并接收响应
func (c *UDPClient) Send(message string) (string, error) {
	// 1. 发送数据
	_, err := c.conn.Write([]byte(message))
	if err != nil {
		return "", fmt.Errorf("发送失败: %w", err)
	}

	// 2. 接收响应
	buf := make([]byte, 1024)
	n, err := c.conn.Read(buf)
	if err != nil {
		return "", fmt.Errorf("接收失败: %w", err)
	}

	return string(buf[:n]), nil
}

// Close 关闭连接
func (c *UDPClient) Close() error {
	return c.conn.Close()
}

// ====== 高级：无连接 UDP 通信 ======

// UnconnectedUDP 演示无连接的 UDP 通信
// 每次发送都需要指定目标地址
func UnconnectedUDPExample() {
	// 解析地址
	addr, err := net.ResolveUDPAddr("udp", "localhost:8080")
	if err != nil {
		log.Fatal(err)
	}

	// 创建 UDP 连接（本地地址为 nil，系统自动选择）
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// 发送数据（每次 Write 都是独立的数据报）
	_, err = conn.Write([]byte("Hello UDP"))
	if err != nil {
		log.Fatal(err)
	}

	// 接收响应
	buf := make([]byte, 1024)
	n, _, err := conn.ReadFromUDP(buf)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("收到响应: %s\n", string(buf[:n]))
}

// ====== 高级：广播与组播 ======

// BroadcastServer 广播服务器示例
func BroadcastServerExample() {
	// 解析广播地址
	broadcastAddr, err := net.ResolveUDPAddr("udp", "255.255.255.255:9999")
	if err != nil {
		log.Fatal(err)
	}

	// 创建连接
	conn, err := net.DialUDP("udp", nil, broadcastAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// 发送广播消息
	for {
		_, err := conn.Write([]byte("Broadcast message"))
		if err != nil {
			log.Printf("广播失败: %v", err)
		}
		time.Sleep(5 * time.Second)
	}
}

// MulticastServer 组播服务器示例
func MulticastServerExample() {
	// 组播地址通常是 224.0.0.0 到 239.255.255.255
	multicastAddr, err := net.ResolveUDPAddr("udp", "239.0.0.1:9999")
	if err != nil {
		log.Fatal(err)
	}

	// 创建 UDP 连接
	conn, err := net.DialUDP("udp", nil, multicastAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// 发送组播消息
	for {
		_, err := conn.Write([]byte("Multicast message"))
		if err != nil {
			log.Printf("组播失败: %v", err)
		}
		time.Sleep(5 * time.Second)
	}
}

// MulticastReceiver 组播接收者示例
func MulticastReceiverExample() {
	// 解析组播地址和本地接口
	multicastAddr, err := net.ResolveUDPAddr("udp", "239.0.0.1:9999")
	if err != nil {
		log.Fatal(err)
	}

	// 创建 UDP 连接
	conn, err := net.ListenMulticastUDP("udp", nil, multicastAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// 设置读写缓冲区大小（可选）
	conn.SetReadBuffer(1024 * 1024) // 1MB

	buf := make([]byte, 1024)
	for {
		n, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Printf("读取错误: %v", err)
			continue
		}
		fmt.Printf("收到组播消息 from %s: %s\n", addr.String(), string(buf[:n]))
	}
}

// ====== 高级：心跳检测 ======

// HeartbeatServer 带心跳检测的 UDP 服务器
type HeartbeatServer struct {
	clients    map[string]time.Time // 客户端地址 -> 最后活跃时间
	mu         sync.RWMutex
	timeout    time.Duration
	cleanupInt time.Duration
}

// NewHeartbeatServer 创建心跳检测服务器
func NewHeartbeatServer(timeout time.Duration) *HeartbeatServer {
	return &HeartbeatServer{
		clients:    make(map[string]time.Time),
		timeout:    timeout,
		cleanupInt: timeout / 2, // 每超时时间的一半清理一次
	}
}

func (hs *HeartbeatServer) Start(address string) error {
	udpAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return err
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return err
	}
	defer conn.Close()

	buf := make([]byte, 1024)

	// 启动清理协程
	go hs.cleanupLoop()

	for {
		conn.SetReadDeadline(time.Now().Add(hs.cleanupInt))
		n, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			continue
		}

		message := string(buf[:n])

		// 更新客户端活跃时间
		hs.mu.Lock()
		hs.clients[addr.String()] = time.Now()
		hs.mu.Unlock()

		// 处理消息
		if message == "ping" {
			conn.WriteToUDP([]byte("pong"), addr)
		}
	}
}

func (hs *HeartbeatServer) cleanupLoop() {
	ticker := time.NewTicker(hs.cleanupInt)
	defer ticker.Stop()

	for range ticker.C {
		hs.mu.Lock()
		now := time.Now()
		for addr, lastActive := range hs.clients {
			if now.Sub(lastActive) > hs.timeout {
				delete(hs.clients, addr)
				log.Printf("客户端超时移除: %s", addr)
			}
		}
		hs.mu.Unlock()
	}
}

// ====== 主函数 ======

func main() {
	fmt.Println("=== UDP 服务器与客户端示例 ===")

	// 启动服务器
	server := NewUDPServer(":8080")

	go func() {
		if err := server.Start(); err != nil {
			log.Printf("服务器错误: %v", err)
		}
	}()

	// 等待服务器启动
	time.Sleep(100 * time.Millisecond)

	// 创建客户端
	client, err := NewUDPClient("localhost:8080")
	if err != nil {
		log.Fatalf("创建客户端失败: %v", err)
	}
	defer client.Close()

	// 测试各种消息
	tests := []string{"ping", "time", "Hello UDP", "date"}

	for _, test := range tests {
		response, err := client.Send(test)
		if err != nil {
			log.Printf("测试 %s 失败: %v", test, err)
			continue
		}
		fmt.Printf("发送: %s -> 收到: %s\n", test, response)
	}

	// 关闭服务器
	server.Close()

	fmt.Println("UDP 示例完成")
}
