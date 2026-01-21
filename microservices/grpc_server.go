// microservices/grpc_server.go
// gRPC 服务端示例 - 详细注释版

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	pb "GolangTutorial/microservices/proto"
)

// ====== gRPC 服务端基础 ======
/*
gRPC 是 Google 开发的高性能远程过程调用（RPC）框架。

主要特点：
1. 高性能 - 基于 HTTP/2 和 Protocol Buffers
2. 跨语言 - 支持多种编程语言
3. 接口定义 - 使用 Protocol Buffers 定义服务
4. 流式支持 - 支持单向、双向流式传输
5. 认证支持 - 内置 TLS 和 Token 认证

安装：
  go get -u google.golang.org/grpc
  go get -u google.golang.org/protobuf/proto
*/

// server 结构体实现 UserServiceServer 接口
type server struct {
	pb.UnimplementedUserServiceServer
	users map[int64]*pb.User // 内存存储用户数据
}

// NewServer 创建新的服务器实例
func NewServer() *server {
	return &server{
		users: make(map[int64]*pb.User),
	}
}

// CreateUser 创建用户
// 实现 UserServiceServer 接口的 CreateUser 方法
func (s *server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	// 1. 验证请求
	if req.Username == "" {
		return nil, status.Error(codes.InvalidArgument, "Username is required")
	}
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "Email is required")
	}

	// 2. 创建用户
	user := &pb.User{
		Id:       generateID(), // 生成唯一 ID
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	// 3. 存储用户
	s.users[user.Id] = user

	log.Printf("创建用户: %s (ID: %d)", user.Username, user.Id)

	// 4. 返回响应
	return &pb.CreateUserResponse{
		User: user,
	}, nil
}

// GetUser 获取用户
func (s *server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	// 1. 验证请求
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "User ID is required")
	}

	// 2. 查找用户
	user, exists := s.users[req.Id]
	if !exists {
		return nil, status.Error(codes.NotFound, "User not found")
	}

	return &pb.GetUserResponse{
		User: user,
	}, nil
}

// ListUsers 列出所有用户
func (s *server) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	// 1. 收集所有用户
	users := make([]*pb.User, 0, len(s.users))
	for _, user := range s.users {
		users = append(users, user)
	}

	// 2. 返回响应
	return &pb.ListUsersResponse{
		Users: users,
		Count: int32(len(users)),
	}, nil
}

// UpdateUser 更新用户
func (s *server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	// 1. 验证请求
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "User ID is required")
	}

	// 2. 查找用户
	user, exists := s.users[req.Id]
	if !exists {
		return nil, status.Error(codes.NotFound, "User not found")
	}

	// 3. 更新字段
	if req.Username != "" {
		user.Username = req.Username
	}
	if req.Email != "" {
		user.Email = req.Email
	}

	// 4. 返回响应
	return &pb.UpdateUserResponse{
		User: user,
	}, nil
}

// DeleteUser 删除用户
func (s *server) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	// 1. 验证请求
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "User ID is required")
	}

	// 2. 删除用户
	if _, exists := s.users[req.Id]; !exists {
		return nil, status.Error(codes.NotFound, "User not found")
	}

	delete(s.users, req.Id)

	log.Printf("删除用户: ID=%d", req.Id)

	return &pb.DeleteUserResponse{
		Success: true,
	}, nil
}

// SearchUsers 搜索用户（服务端流式）
func (s *server) SearchUsers(req *pb.SearchUsersRequest, stream pb.UserService_SearchUsersServer) error {
	// 1. 遍历所有用户
	for _, user := range s.users {
		// 2. 检查是否匹配搜索条件
		matched := true

		if req.UsernamePrefix != "" {
			// 检查用户名是否以指定前缀开头
			matched = matched && len(user.Username) >= len(req.UsernamePrefix) &&
				user.Username[:len(req.UsernamePrefix)] == req.UsernamePrefix
		}

		if req.MinAge > 0 {
			// 检查年龄是否大于最小年龄
			matched = matched && user.Age >= req.MinAge
		}

		if matched {
			// 3. 发送匹配的用户到流
			if err := stream.Send(&pb.SearchUsersResponse{User: user}); err != nil {
				return err
			}
		}
	}

	return nil
}

// Chat stream 用户聊天（双向流式）
func (s *server) Chat(stream pb.UserService_ChatServer) error {
	for {
		// 1. 接收消息
		req, err := stream.Recv()
		if err != nil {
			// 流结束时返回 nil
			return nil
		}

		// 2. 处理消息
		message := req.Message
		userID := req.UserId

		// 3. 生成响应
		response := fmt.Sprintf("收到来自用户 %d 的消息: %s", userID, message)

		// 4. 发送响应
		if err := stream.Send(&pb.ChatResponse{
			Message: response,
			UserId:  userID,
		}); err != nil {
			return err
		}
	}
}

// ====== 辅助函数 ======

// generateID 生成唯一 ID
// 在实际应用中，应该使用数据库自增 ID 或 UUID
var idCounter int64 = 0

func generateID() int64 {
	idCounter++
	return idCounter
}

// ====== 主函数 ======

func main() {
	// 1. 解析命令行参数
	port := flag.Int("port", 50051, "gRPC 服务器端口")
	flag.Parse()

	// 2. 创建监听器
	addr := fmt.Sprintf(":%d", *port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("监听端口失败: %v", err)
	}

	log.Printf("gRPC 服务器启动，监听地址: %s", addr)

	// 3. 创建 gRPC 服务器
	// grpc.NewServer 创建新的 gRPC 服务器实例
	s := grpc.NewServer(
		// 4. 配置服务器选项（可选）
		grpc.MaxRecvMsgSize(10*1024*1024), // 最大接收消息大小 10MB
		grpc.MaxSendMsgSize(10*1024*1024), // 最大发送消息大小 10MB
	)

	// 5. 注册服务
	// 将服务实现注册到 gRPC 服务器
	pb.RegisterUserServiceServer(s, NewServer())

	// 6. 启用反射（用于调试工具如 grpcurl）
	reflection.Register(s)

	// 7. 启动服务器
	// Serve 开始接受连接并处理请求
	if err := s.Serve(lis); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
}
