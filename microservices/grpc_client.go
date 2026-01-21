// microservices/grpc_client.go
// gRPC 客户端示例 - 详细注释版

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	pb "GolangTutorial/microservices/proto"
)

// ====== gRPC 客户端基础 ======

// UserClient gRPC 客户端封装
type UserClient struct {
	client pb.UserServiceClient // 生成的客户端接口
	conn   *grpc.ClientConn     // 连接实例
}

// NewUserClient 创建新的客户端
func NewUserClient(address string) (*UserClient, error) {
	// 1. 创建连接
	// grpc.Dial 连接到 gRPC 服务器
	// WithTransportCredentials 设置传输凭证
	// insecure.NewCredentials() 表示不使用 TLS（仅用于开发）
	conn, err := grpc.Dial(address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(), // 阻塞直到连接成功或超时
		grpc.WithTimeout(5*time.Second),
	)
	if err != nil {
		return nil, fmt.Errorf("连接服务器失败: %w", err)
	}

	// 2. 创建客户端
	client := pb.NewUserServiceClient(conn)

	log.Printf("连接到 gRPC 服务器: %s", address)

	return &UserClient{
		client: client,
		conn:   conn,
	}, nil
}

// Close 关闭连接
func (c *UserClient) Close() error {
	return c.conn.Close()
}

// ====== 客户端方法 ======

// CreateUser 创建用户
func (c *UserClient) CreateUser(username, email, password string) (*pb.User, error) {
	// 1. 创建请求
	req := &pb.CreateUserRequest{
		Username: username,
		Email:    email,
		Password: password,
	}

	// 2. 调用远程方法
	// context 用于设置超时和取消
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := c.client.CreateUser(ctx, req)
	if err != nil {
		// 3. 处理错误
		// 从错误中提取 gRPC 状态码
		if st, ok := status.FromError(err); ok {
			return nil, fmt.Errorf("gRPC 错误 [%d]: %s", st.Code(), st.Message())
		}
		return nil, fmt.Errorf("调用失败: %w", err)
	}

	return resp.User, nil
}

// GetUser 获取用户
func (c *UserClient) GetUser(id int64) (*pb.User, error) {
	req := &pb.GetUserRequest{Id: id}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := c.client.GetUser(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.User, nil
}

// ListUsers 列出所有用户
func (c *UserClient) ListUsers() ([]*pb.User, error) {
	req := &pb.ListUsersRequest{}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := c.client.ListUsers(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.Users, nil
}

// UpdateUser 更新用户
func (c *UserClient) UpdateUser(id int64, username, email string) (*pb.User, error) {
	req := &pb.UpdateUserRequest{
		Id:       id,
		Username: username,
		Email:    email,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := c.client.UpdateUser(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.User, nil
}

// DeleteUser 删除用户
func (c *UserClient) DeleteUser(id int64) error {
	req := &pb.DeleteUserRequest{Id: id}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := c.client.DeleteUser(ctx, req)
	return err
}

// SearchUsers 搜索用户（客户端流式）
func (c *UserClient) SearchUsers(usernamePrefix string, minAge int32) error {
	// 1. 创建搜索请求
	req := &pb.SearchUsersRequest{
		UsernamePrefix: usernamePrefix,
		MinAge:         minAge,
	}

	// 2. 发起流式调用
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	stream, err := c.client.SearchUsers(ctx, req)
	if err != nil {
		return err
	}

	// 3. 接收流式响应
	log.Println("搜索结果:")
	for {
		resp, err := stream.Recv()
		if err != nil {
			// 流结束时退出
			break
		}

		user := resp.User
		log.Printf("  - %s (%s), Age: %d", user.Username, user.Email, user.Age)
	}

	return nil
}

// Chat 聊天（双向流式）
func (c *UserClient) Chat(userID int64, messages []string) error {
	// 1. 发起双向流式调用
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	stream, err := c.client.Chat(ctx)
	if err != nil {
		return err
	}

	// 2. 发送消息
	for _, msg := range messages {
		if err := stream.Send(&pb.ChatRequest{
			UserId:  userID,
			Message: msg,
		}); err != nil {
			return err
		}
	}

	// 3. 关闭发送流
	if err := stream.CloseSend(); err != nil {
		return err
	}

	// 4. 接收响应
	for {
		resp, err := stream.Recv()
		if err != nil {
			break
		}
		log.Printf("收到响应: %s", resp.Message)
	}

	return nil
}

// ====== 主函数 ======

func main() {
	fmt.Println("=== gRPC 客户端示例 ===")

	// 1. 解析命令行参数
	serverAddr := flag.String("server", "localhost:50051", "gRPC 服务器地址")
	flag.Parse()

	// 2. 创建客户端
	client, err := NewUserClient(*serverAddr)
	if err != nil {
		log.Fatalf("创建客户端失败: %v", err)
	}
	defer client.Close()

	// 3. 测试创建用户
	fmt.Println("\n--- 创建用户 ---")
	user1, err := client.CreateUser("alice", "alice@example.com", "password123")
	if err != nil {
		log.Printf("创建用户 alice 失败: %v", err)
	} else {
		fmt.Printf("创建用户成功: %s (ID: %d)\n", user1.Username, user1.Id)
	}

	user2, _ := client.CreateUser("bob", "bob@example.com", "password456")
	if user2 != nil {
		fmt.Printf("创建用户成功: %s (ID: %d)\n", user2.Username, user2.Id)
	}

	user3, _ := client.CreateUser("charlie", "charlie@example.com", "password789")
	if user3 != nil {
		fmt.Printf("创建用户成功: %s (ID: %d)\n", user3.Username, user3.Id)
	}

	// 4. 测试获取用户
	fmt.Println("\n--- 获取用户 ---")
	user, err := client.GetUser(1)
	if err != nil {
		log.Printf("获取用户失败: %v", err)
	} else {
		fmt.Printf("用户信息: %s, %s, 年龄: %d\n", user.Username, user.Email, user.Age)
	}

	// 5. 测试列出用户
	fmt.Println("\n--- 列出所有用户 ---")
	users, err := client.ListUsers()
	if err != nil {
		log.Printf("列出用户失败: %v", err)
	} else {
		fmt.Printf("共有 %d 个用户:\n", len(users))
		for _, u := range users {
			fmt.Printf("  - %s (%s)\n", u.Username, u.Email)
		}
	}

	// 6. 测试更新用户
	fmt.Println("\n--- 更新用户 ---")
	updatedUser, err := client.UpdateUser(1, "alice_updated", "alice.new@example.com")
	if err != nil {
		log.Printf("更新用户失败: %v", err)
	} else {
		fmt.Printf("更新成功: %s, %s\n", updatedUser.Username, updatedUser.Email)
	}

	// 7. 测试搜索用户
	fmt.Println("\n--- 搜索用户 ---")
	_ = client.SearchUsers("a", 0)

	// 8. 测试聊天
	fmt.Println("\n--- 聊天测试 ---")
	_ = client.Chat(1, []string{"Hello!", "How are you?", "Bye!"})

	// 9. 测试删除用户
	fmt.Println("\n--- 删除用户 ---")
	err = client.DeleteUser(3)
	if err != nil {
		log.Printf("删除用户失败: %v", err)
	} else {
		fmt.Println("删除用户成功")
	}

	// 10. 验证删除
	fmt.Println("\n--- 验证删除 ---")
	users, _ = client.ListUsers()
	fmt.Printf("剩余 %d 个用户\n", len(users))

	fmt.Println("\n客户端测试完成")
}
