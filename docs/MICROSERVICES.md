# Go 微服务详解

## 目录
- [1. 微服务概述](#1-微服务概述)
- [2. ProtoBuf 语言](#2-protobuf-语言)
- [3. gRPC 基础](#3-grpc-基础)
- [4. gRPC 高级特性](#4-grpc-高级特性)
- [5. 服务发现与注册](#5-服务发现与注册)
- [6. 负载均衡](#6-负载均衡)
- [7. 链路追踪](#7-链路追踪)

---

## 1. 微服务概述

### 1.1 微服务架构

```
┌─────────────────────────────────────────────────────────────┐
│                      微服务架构                               │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ┌─────────┐   ┌─────────┐   ┌─────────┐                   │
│  │ API     │   │ User    │   │ Order   │                   │
│  │ Gateway │──▶│ Service │──▶│ Service │                   │
│  └─────────┘   └─────────┘   └─────────┘                   │
│       │              │              │                       │
│       └──────────────┼──────────────┘                       │
│                      │                                      │
│              ┌───────┴───────┐                              │
│              │   Service     │                              │
│              │   Discovery   │                              │
│              └───────────────┘                              │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

### 1.2 Go 微服务技术栈

| 层级 | 技术 | 用途 |
|------|------|------|
| 通信 | gRPC, HTTP/2 | 高效服务间通信 |
| 序列化 | ProtoBuf, JSON | 数据序列化 |
| 服务发现 | Consul, Etcd, Kubernetes | 服务注册与发现 |
| 网关 | Kong, Nginx, Traefik | API 网关 |
| 追踪 | Jaeger, Zipkin | 分布式追踪 |
| 监控 | Prometheus, Grafana | 指标监控 |
| 配置 | Consul, ConfigMap | 配置管理 |

---

## 2. ProtoBuf 语言

### 2.1 ProtoBuf 简介

Protocol Buffers（ProtoBuf）是 Google 开发的一种语言无关、平台无关的可扩展序列化机制。

```protobuf
// 安装 protoc
// macOS: brew install protoc
// Linux: apt install protobuf-compiler
// Windows: 下载 protoc-xxx-win64.zip

// 安装 Go 插件
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

### 2.2 ProtoBuf 语法

```protobuf
// 语法声明（必须）
syntax = "proto3";

// 包声明（避免命名冲突）
package user;

// Go 包名（生成的 Go 文件包名）
option go_package = "github.com/example/proto/user";

// 导入其他 proto 文件
import "common/address.proto";

// ========== 消息定义 ==========

// 简单消息
message User {
    // 字段规则（proto3 不支持 required/optional）
    string name = 1;        // 字符串
    int32 age = 2;          // 32 位整数
    int64 id = 3;           // 64 位整数
    float height = 4;       // 32 位浮点
    double weight = 5;      // 64 位浮点
    bool is_active = 6;     // 布尔值
    bytes avatar = 7;       // 字节数组
    Status status = 8;      // 枚举类型
    Address address = 9;    // 嵌套消息
    repeated string tags = 10; // 数组（repeated）
}

// 枚举类型
enum Status {
    STATUS_UNSPECIFIED = 0; // 必须有默认值 0
    STATUS_ACTIVE = 1;
    STATUS_INACTIVE = 2;
    STATUS_DELETED = 3;
}

// 嵌套消息
message Address {
    string street = 1;
    string city = 2;
    string country = 3;
    int32 zip_code = 4;
}

// ========== 服务定义 ==========

// 用户服务
service UserService {
    // 简单 RPC
    rpc GetUser(GetUserRequest) returns (User);
    
    // 服务端流式 RPC
    rpc ListUsers(ListUsersRequest) returns (stream User);
    
    // 客户端流式 RPC
    rpc CreateUsers(stream CreateUserRequest) returns (CreateUsersResponse);
    
    // 双向流式 RPC
    rpc Chat(stream ChatRequest) returns (stream ChatResponse);
}

// 请求/响应消息
message GetUserRequest {
    int64 id = 1;
}

message ListUsersRequest {
    int32 page = 1;
    int32 page_size = 2;
    string keyword = 3;
}

message CreateUserRequest {
    User user = 1;
}

message CreateUsersResponse {
    int32 created_count = 1;
    repeated int64 ids = 2;
}

// ========== Map 类型 ==========

message UserMap {
    map<int64, User> users = 1;           // map<key_type, value_type>
    map<string, string> metadata = 2;
}

// ========== Oneof 类型 ==========

message Result {
    oneof result {
        User user = 1;
        string error = 2;
    }
}

// ========== 高级特性 ==========

// 保留字段（避免使用已删除的字段编号）
message Legacy {
    reserved 3, 15 to 20;  // 保留字段编号
    reserved "deprecated_field"; // 保留字段名
}

// 包别名
import "google/protobuf/timestamp.proto" as timestamp;

// 时间戳（使用标准类型）
google.protobuf.Timestamp created_at = 1;

// 任意类型
google.protobuf.Any metadata = 2;

// JSON 映射
message Data {
    google.protobuf.Struct json_data = 1;
}
```

### 2.3 ProtoBuf 最佳实践

```protobuf
// 1. 使用 proto3 语法（更简洁，无 required/optional）

// 2. 字段编号从 1 开始，保留小编号给常用字段
message User {
    int64 id = 1;           // 最重要
    string name = 2;        // 常用
    // ... 其他字段
}

// 3. 为字段添加注释
message User {
    // 用户唯一标识
    int64 id = 1;
    
    // 用户名，长度 2-50 字符
    string name = 2;
}

// 4. 使用枚举代替布尔值（更易扩展）
enum UserStatus {
    USER_STATUS_UNSPECIFIED = 0;
    USER_STATUS_ACTIVE = 1;
    USER_STATUS_INACTIVE = 2;
}

message User {
    int64 id = 1;
    UserStatus status = 2;
}

// 5. 使用 repeated 表示数组
message User {
    repeated string email = 1; // 多个邮箱
}

// 6. 使用 map 表示键值对
message UserPreferences {
    map<string, string> settings = 1;
}

// 7. 合理使用嵌套和组合
message Order {
    message Item {
        string product_id = 1;
        int32 quantity = 2;
    }
    
    repeated Item items = 1;
    User user = 2;
}

// 8. 使用包名避免冲突
package payment;
message Payment {
    int64 id = 1;
}
```

---

## 3. gRPC 基础

### 3.1 安装和配置

```bash
# 安装 gRPC
go get google.golang.org/grpc@v1.59.0

# 安装 ProtoBuf 工具
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.31.0
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0
```

### 3.2 定义服务

```protobuf
// proto/user.proto
syntax = "proto3";

package user;

option go_package = "github.com/example/user-service/proto/user";

import "google/protobuf/timestamp.proto";

// 用户服务
service UserService {
    // 获取用户
    rpc GetUser(GetUserRequest) returns (User);
    
    // 创建用户
    rpc CreateUser(CreateUserRequest) returns (CreateUserResponse);
    
    // 更新用户
    rpc UpdateUser(UpdateUserRequest) returns (UpdateUserResponse);
    
    // 删除用户
    rpc DeleteUser(DeleteUserRequest) returns (DeleteUserResponse);
    
    // 列出用户
    rpc ListUsers(ListUsersRequest) returns (ListUsersResponse);
    
    // 批量获取用户
    rpc BatchGetUsers(BatchGetUsersRequest) returns (BatchGetUsersResponse);
}

// 消息定义
message User {
    int64 id = 1;
    string username = 2;
    string email = 3;
    string phone = 4;
    UserStatus status = 5;
    google.protobuf.Timestamp created_at = 6;
    google.protobuf.Timestamp updated_at = 7;
}

enum UserStatus {
    USER_STATUS_UNSPECIFIED = 0;
    USER_STATUS_ACTIVE = 1;
    USER_STATUS_INACTIVE = 2;
    USER_STATUS_DELETED = 3;
}

message GetUserRequest {
    int64 id = 1;
}

message CreateUserRequest {
    string username = 1;
    string email = 2;
    string phone = 3;
    string password = 4;
}

message CreateUserResponse {
    User user = 1;
}

message UpdateUserRequest {
    int64 id = 1;
    string username = 2;
    string email = 3;
    string phone = 4;
}

message UpdateUserResponse {
    User user = 1;
}

message DeleteUserRequest {
    int64 id = 1;
}

message DeleteUserResponse {
    bool success = 1;
}

message ListUsersRequest {
    int32 page = 1;
    int32 page_size = 2;
    string keyword = 3;
}

message ListUsersResponse {
    repeated User users = 1;
    int32 total = 2;
}

message BatchGetUsersRequest {
    repeated int64 ids = 1;
}

message BatchGetUsersResponse {
    repeated User users = 1;
}
```

### 3.3 生成代码

```bash
# 生成 Go 代码
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       proto/user.proto
```

### 3.4 实现服务端

```go
// server/main.go
package main

import (
    "context"
    "log"
    "net"

    "google.golang.org/grpc"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
    pb "github.com/example/user-service/proto/user"
)

// UserServer 实现 UserServiceServer 接口
type UserServer struct {
    pb.UnimplementedUserServiceServer
    
    // 存储层（实际项目中应该是数据库）
    users map[int64]*pb.User
    nextID int64
}

func NewUserServer() *UserServer {
    return &UserServer{
        users:  make(map[int64]*pb.User),
        nextID: 1,
    }
}

func (s *UserServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
    user, ok := s.users[req.Id]
    if !ok {
        return nil, status.Errorf(codes.NotFound, "user not found: %d", req.Id)
    }
    return user, nil
}

func (s *UserServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
    // 验证请求
    if req.Username == "" {
        return nil, status.Error(codes.InvalidArgument, "username is required")
    }
    if req.Email == "" {
        return nil, status.Error(codes.InvalidArgument, "email is required")
    }

    // 创建用户
    user := &pb.User{
        Id:       s.nextID,
        Username: req.Username,
        Email:    req.Email,
        Phone:    req.Phone,
        Status:   pb.UserStatus_USER_STATUS_ACTIVE,
    }
    s.users[user.Id] = user
    s.nextID++

    return &pb.CreateUserResponse{User: user}, nil
}

func (s *UserServer) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
    user, ok := s.users[req.Id]
    if !ok {
        return nil, status.Errorf(codes.NotFound, "user not found: %d", req.Id)
    }

    // 更新字段
    if req.Username != "" {
        user.Username = req.Username
    }
    if req.Email != "" {
        user.Email = req.Email
    }
    if req.Phone != "" {
        user.Phone = req.Phone
    }

    return &pb.UpdateUserResponse{User: user}, nil
}

func (s *UserServer) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
    if _, ok := s.users[req.Id]; !ok {
        return nil, status.Errorf(codes.NotFound, "user not found: %d", req.Id)
    }
    
    delete(s.users, req.Id)
    return &pb.DeleteUserResponse{Success: true}, nil
}

func (s *UserServer) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
    var users []*pb.User
    for _, user := range s.users {
        users = append(users, user)
    }

    return &pb.ListUsersResponse{
        Users: users,
        Total: int32(len(users)),
    }, nil
}

func (s *UserServer) BatchGetUsers(ctx context.Context, req *pb.BatchGetUsersRequest) (*pb.BatchGetUsersResponse, error) {
    var users []*pb.User
    for _, id := range req.Ids {
        if user, ok := s.users[id]; ok {
            users = append(users, user)
        }
    }
    return &pb.BatchGetUsersResponse{Users: users}, nil
}

func main() {
    // 创建 gRPC 服务器
    server := grpc.NewServer(
        grpc.UnaryInterceptor(loggingInterceptor),
    )

    // 注册服务
    userServer := NewUserServer()
    pb.RegisterUserServiceServer(server, userServer)

    // 监听端口
    listener, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }

    log.Println("gRPC server listening on :50051")
    if err := server.Serve(listener); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}

// 日志拦截器
func loggingInterceptor(ctx context.Context, req interface{}, 
    info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
    log.Printf("gRPC call: %s", info.FullMethod)
    resp, err := handler(ctx, req)
    if err != nil {
        log.Printf("gRPC error: %v", err)
    }
    return resp, err
}
```

### 3.5 实现客户端

```go
// client/main.go
package main

import (
    "context"
    "log"
    "time"

    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
    pb "github.com/example/user-service/proto/user"
)

func main() {
    // 创建连接
    conn, err := grpc.Dial("localhost:50051",
        grpc.WithTransportCredentials(insecure.NewCredentials()),
        grpc.WithBlock(),
        grpc.WithTimeout(5*time.Second),
    )
    if err != nil {
        log.Fatalf("failed to connect: %v", err)
    }
    defer conn.Close()

    // 创建客户端
    client := pb.NewUserServiceClient(conn)

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // 1. 创建用户
    createResp, err := client.CreateUser(ctx, &pb.CreateUserRequest{
        Username: "张三",
        Email:    "zhangsan@example.com",
        Phone:    "13800138000",
    })
    if err != nil {
        log.Fatalf("failed to create user: %v", err)
    }
    log.Printf("Created user: %v", createResp.User)

    userID := createResp.User.Id

    // 2. 获取用户
    getResp, err := client.GetUser(ctx, &pb.GetUserRequest{
        Id: userID,
    })
    if err != nil {
        log.Fatalf("failed to get user: %v", err)
    }
    log.Printf("Got user: %v", getResp.User)

    // 3. 更新用户
    updateResp, err := client.UpdateUser(ctx, &pb.UpdateUserRequest{
        Id:       userID,
        Username: "张三三",
    })
    if err != nil {
        log.Fatalf("failed to update user: %v", err)
    }
    log.Printf("Updated user: %v", updateResp.User)

    // 4. 列出用户
    listResp, err := client.ListUsers(ctx, &pb.ListUsersRequest{
        Page:     1,
        PageSize: 10,
    })
    if err != nil {
        log.Fatalf("failed to list users: %v", err)
    }
    log.Printf("Total users: %d", listResp.Total)

    // 5. 删除用户
    _, err = client.DeleteUser(ctx, &pb.DeleteUserRequest{
        Id: userID,
    })
    if err != nil {
        log.Fatalf("failed to delete user: %v", err)
    }
    log.Printf("Deleted user: %d", userID)
}
```

---

## 4. gRPC 高级特性

### 4.1 流式 RPC

```protobuf
// stream.proto
syntax = "proto3";

package chat;

option go_package = "github.com/example/chat-service/proto/chat";

service ChatService {
    // 服务端流式
    rpc GetMessages(GetMessagesRequest) returns (stream Message);
    
    // 客户端流式
    rpc SendMessages(stream MessageRequest) returns (SendMessagesResponse);
    
    // 双向流式
    rpc Chat(stream ChatRequest) returns (stream ChatResponse);
}

message Message {
    int64 id = 1;
    string content = 2;
    string sender = 3;
    int64 timestamp = 4;
}

message GetMessagesRequest {
    string room_id = 1;
    int64 since = 2;
}

message MessageRequest {
    string room_id = 1;
    string content = 2;
    string sender = 3;
}

message SendMessagesResponse {
    int32 sent_count = 1;
}

message ChatRequest {
    string room_id = 1;
    string content = 2;
}

message ChatResponse {
    string content = 1;
    string sender = 2;
    int64 timestamp = 3;
}
```

```go
// 流式服务端实现
type ChatServer struct {
    pb.UnimplementedChatServiceServer
    messages chan *pb.Message
}

func (s *ChatServer) GetMessages(req *pb.GetMessagesRequest, 
    stream pb.ChatService_GetMessagesServer) error {
    
    for msg := range s.messages {
        if msg.Timestamp > req.Since {
            if err := stream.Send(msg); err != nil {
                return err
            }
        }
    }
    return nil
}

func (s *ChatServer) SendMessages(stream pb.ChatService_SendMessagesServer) error {
    var count int32
    for {
        req, err := stream.Recv()
        if err == io.EOF {
            break
        }
        if err != nil {
            return err
        }
        
        // 处理消息
        s.messages <- &pb.Message{
            Content:   req.Content,
            Sender:    req.Sender,
            Timestamp: time.Now().Unix(),
        }
        count++
    }
    return stream.SendAndClose(&pb.SendMessagesResponse{
        SentCount: count,
    })
}

func (s *ChatServer) Chat(stream pb.ChatService_ChatServer) error {
    for {
        req, err := stream.Recv()
        if err == io.EOF {
            return nil
        }
        if err != nil {
            return err
        }

        // 广播消息
        resp := &pb.ChatResponse{
            Content:   req.Content,
            Sender:    "Server",
            Timestamp: time.Now().Unix(),
        }
        if err := stream.Send(resp); err != nil {
            return err
        }
    }
}
```

```go
// 流式客户端调用
func main() {
    conn, _ := grpc.Dial("localhost:50052")
    client := pb.NewChatServiceClient(conn)
    ctx := context.Background()

    // 服务端流式
    stream, err := client.GetMessages(ctx, &pb.GetMessagesRequest{
        RoomId: "room1",
    })
    for {
        msg, err := stream.Recv()
        if err == io.EOF {
            break
        }
        log.Println(msg)
    }

    // 客户端流式
    stream, err = client.SendMessages(ctx)
    for i := 0; i < 5; i++ {
        stream.Send(&pb.MessageRequest{
            RoomId:  "room1",
            Content: fmt.Sprintf("Message %d", i),
            Sender:  "client",
        })
    }
    resp, _ := stream.CloseAndRecv()
    log.Printf("Sent: %d messages", resp.SentCount)

    // 双向流式
    stream, err = client.Chat(ctx)
    go func() {
        for {
            stream.Send(&pb.ChatRequest{
                RoomId:  "room1",
                Content: "Hello",
            })
            time.Sleep(time.Second)
        }
    }()
    for {
        resp, err := stream.Recv()
        if err != nil {
            break
        }
        log.Println(resp)
    }
}
```

### 4.2 拦截器

```go
// 拦截器类型
type UnaryServerInterceptor func(ctx context.Context, req interface{}, 
    info *UnaryServerInfo, handler UnaryHandler) (interface{}, error)

// 创建拦截器
func chainInterceptor(interceptors ...grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
    n := len(interceptors)
    return func(ctx context.Context, req interface{}, 
        info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
        chain := func(n int) grpc.UnaryHandler {
            if n == 0 {
                return handler
            }
            return func(ctx context.Context, req interface{}) (interface{}, error) {
                return interceptors[n-1](ctx, req, info, chain(n-1))
            }
        }
        return chain(n)(ctx, req)
    }
}

// 日志拦截器
func loggingInterceptor(ctx context.Context, req interface{}, 
    info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
    
    start := time.Now()
    
    resp, err := handler(ctx, req)
    
    log.Printf("gRPC %s - %v - %v", 
        info.FullMethod, time.Since(start), err)
    
    return resp, err
}

// 认证拦截器
func authInterceptor(ctx context.Context, req interface{}, 
    info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
    
    // 从 metadata 获取 token
    md, ok := metadata.FromIncomingContext(ctx)
    if !ok {
        return nil, status.Error(codes.Unauthenticated, "missing metadata")
    }
    
    token := md.Get("authorization")
    if len(token) == 0 {
        return nil, status.Error(codes.Unauthenticated, "missing token")
    }
    
    // 验证 token...
    
    return handler(ctx, req)
}

// 限流拦截器
func rateLimitInterceptor(limiter *RateLimiter) grpc.UnaryServerInterceptor {
    return func(ctx context.Context, req interface{}, 
        info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
        
        if !limiter.Allow() {
            return nil, status.Error(codes.ResourceExhausted, "rate limited")
        }
        
        return handler(ctx, req)
    }
}

// 使用拦截器
server := grpc.NewServer(
    grpc.ChainUnaryInterceptor(
        loggingInterceptor,
        authInterceptor,
        rateLimitInterceptor(limiter),
    ),
)
```

### 4.3 错误处理

```go
// 错误处理工具
import (
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

// 创建错误
func validateRequest(req *pb.User) error {
    if req.Username == "" {
        return status.Error(codes.InvalidArgument, "username is required")
    }
    if req.Email == "" {
        return status.Error(codes.InvalidArgument, "email is required")
    }
    return nil
}

// 包装错误
func getUser(id int64) (*User, error) {
    user, ok := users[id]
    if !ok {
        return nil, status.Errorf(codes.NotFound, "user not found: %d", id)
    }
    return user, nil
}

// 客户端处理错误
_, err := client.GetUser(ctx, &pb.GetUserRequest{Id: id})
if err != nil {
    st, ok := status.FromError(err)
    if !ok {
        log.Printf("unknown error: %v", err)
        return
    }
    
    switch st.Code() {
    case codes.NotFound:
        log.Printf("user not found")
    case codes.InvalidArgument:
        log.Printf("invalid argument: %v", st.Details())
    case codes.Internal:
        log.Printf("internal error")
    default:
        log.Printf("error: %v", st.Message())
    }
}
```

### 4.4 负载均衡

```go
// 客户端负载均衡
import (
    "google.golang.org/grpc/balancer/roundrobin"
    "google.golang.org/grpc/resolver"
)

// 服务发现
resolver.Register(&customResolverBuilder{})

// 创建连接（使用服务名）
conn, err := grpc.Dial(
    "service:///user-service", // service:// 是自定义 scheme
    grpc.WithBalancerName(roundrobin.Name),
    grpc.WithInsecure(),
)
```

```go
// 自定义 Resolver
type customResolverBuilder struct{}

func (b *customResolverBuilder) Build(target resolver.Target, 
    cc resolver.ClientConn) (resolver.Resolver, error) {
    
    r := &customResolver{
        target: target,
        cc:     cc,
    }
    
    r.start()
    return r, nil
}

func (b *customResolverBuilder) Scheme() string {
    return "service"
}

type customResolver struct {
    target resolver.Target
    cc     resolver.ClientConn
}

func (r *customResolver) start() {
    // 从服务发现获取地址
    addresses := []resolver.Address{
        {Addr: "localhost:50051"},
        {Addr: "localhost:50052"},
    }
    
    r.cc.UpdateState(resolver.State{
        Addresses: addresses,
    })
}

func (r *customResolver) ResolveNow(options resolver.ResolveNowOptions) {
    r.start()
}

func (r *customResolver) Close() {}
```

---

## 5. 服务发现与注册

### 5.1 Consul 服务注册

```go
package main

import (
    "log"
    "time"

    "github.com/hashicorp/consul/api"
)

type ServiceConfig struct {
    ID      string
    Name    string
    Address string
    Port    int
    Tags    []string
}

func RegisterToConsul(config ServiceConfig) error {
    // 创建 Consul 客户端
    client, err := api.NewClient(&api.Config{
        Address: "localhost:8500",
    })
    if err != nil {
        return err
    }

    // 创建注册信息
    registration := &api.AgentServiceRegistration{
        ID:      config.ID,
        Name:    config.Name,
        Address: config.Address,
        Port:    config.Port,
        Tags:    config.Tags,
        Check: &api.AgentServiceCheck{
            HTTP:                           "http://localhost:8080/health",
            Interval:                       "10s",
            Timeout:                        "5s",
            DeregisterCriticalServiceAfter: "30s",
        },
    }

    // 注册服务
    err = client.Agent().ServiceRegister(registration)
    if err != nil {
        return err
    }

    log.Printf("Service %s registered to Consul", config.Name)
    return nil
}

func DeregisterFromConsul(client *api.Client, serviceID string) error {
    return client.Agent().ServiceDeregister(serviceID)
}

func main() {
    config := ServiceConfig{
        ID:      "user-service-1",
        Name:    "user-service",
        Address: "localhost",
        Port:    50051,
        Tags:    []string{"grpc", "v1"},
    }

    if err := RegisterToConsul(config); err != nil {
        log.Fatal(err)
    }
}
```

### 5.2 Consul 服务发现

```go
package main

import (
    "log"

    "github.com/hashicorp/consul/api"
)

func DiscoverServices(client *api.Client, serviceName string) ([]string, error) {
    services, _, err := client.Health().Service(serviceName, "", true, nil)
    if err != nil {
        return nil, err
    }

    var addresses []string
    for _, service := range services {
        addr := service.Service.Address + ":" + 
            string(rune(service.Service.Port))
        addresses = append(addresses, addr)
    }

    return addresses, nil
}

func main() {
    client, _ := api.NewClient(&api.Config{
        Address: "localhost:8500",
    })

    addresses, err := DiscoverServices(client, "user-service")
    if err != nil {
        log.Fatal(err)
    }

    for _, addr := range addresses {
        log.Printf("Found service at: %s", addr)
    }
}
```

---

## 6. 负载均衡

### 6.1 服务端负载均衡

```go
// 使用 gRPC 的加权轮询负载均衡
import (
    "google.golang.org/grpc"
    "google.golang.org/grpc/balancer/weightedroundrobin"
)

func main() {
    // 注册负载均衡策略
    weightedroundrobin.Register()
    
    // 创建连接
    conn, _ := grpc.Dial(
        "consul:///user-service",
        grpc.WithBalancerName(weightedroundrobin.Name),
        grpc.WithInsecure(),
    )
}
```

### 6.2 客户端负载均衡器

```go
// 自定义负载均衡器
type CustomBalancer struct {
    picker    *WeightedPicker
    balancer  *grpc.Balancer
}

func NewCustomBalancer() grpc.Balancer {
    return &CustomBalancer{
        picker: NewWeightedPicker(),
    }
}

func (b *CustomBalancer) HandleSubConnStateChange(
    sc grpc.SubConn, state grpc.ConnectivityState) {
    // 处理连接状态变化
}

func (b *CustomBalancer) HandleResolvedAddrs(
    addrs []resolver.Address) {
    // 处理解析的地址
}

func (b *CustomBalancer) Pick(ctx context.Context, 
    opts grpc.PickOptions) (grpc.SubConn, func(balancer.DoneInfo), error) {
    return b.picker.Pick()
}
```

---

## 7. 链路追踪

### 7.1 OpenTelemetry 集成

```go
package main

import (
    "context"
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/exporters/jaeger"
    "go.opentelemetry.io/otel/propagation"
    "go.opentelemetry.io/otel/sdk/resource"
    sdktrace "go.opentelemetry.io/otel/sdk/trace"
    semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

func InitTracer(serviceName string) (func(context.Context) error, error) {
    // 创建 Jaeger exporter
    exporter, err := jaeger.New(
        jaeger.WithCollectorEndpoint(
            jaeger.WithEndpoint("http://localhost:14268/api/traces"),
        ),
    )
    if err != nil {
        return nil, err
    }

    // 创建资源
    res, err := resource.Merge(
        resource.Default(),
        resource.NewWithAttributes(
            semconv.SchemaURL,
            semconv.ServiceName(serviceName),
        ),
    )
    if err != nil {
        return nil, err
    }

    // 创建 tracer provider
    tp := sdktrace.NewTracerProvider(
        sdktrace.WithBatcher(exporter),
        sdktrace.WithResource(res),
    )

    // 设置全局 tracer provider
    otel.SetTracerProvider(tp)
    otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
        propagation.TraceContext{},
        propagation.Baggage{},
    ))

    return tp.Shutdown, nil
}
```

```go
// 在 gRPC 中使用追踪
import (
    "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
)

func main() {
    // 初始化追踪
    shutdown, _ := InitTracer("user-service")
    defer shutdown(context.Background())

    // 创建带追踪的 gRPC 服务器
    server := grpc.NewServer(
        grpc.StatsHandler(otelgrpc.NewServerHandler()),
    )
}
```

---

## 8. 配置管理

### 8.1 环境变量配置

```go
package config

import (
    "os"
    "strconv"
)

type Config struct {
    Port    int
    DBHost  string
    DBPort  int
    DBUser  string
    DBPass  string
    DBName  string
}

func Load() *Config {
    return &Config{
        Port:    getEnvInt("PORT", 50051),
        DBHost:  os.Getenv("DB_HOST"),
        DBPort:  getEnvInt("DB_PORT", 3306),
        DBUser:  os.Getenv("DB_USER"),
        DBPass:  os.Getenv("DB_PASS"),
        DBName:  os.Getenv("DB_NAME"),
    }
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
    if value := os.Getenv(key); value != "" {
        if intValue, err := strconv.Atoi(value); err == nil {
            return intValue
        }
    }
    return defaultValue
}
```

---

## 最佳实践

### 1. 项目结构

```
user-service/
├── proto/
│   └── user.proto
├── internal/
│   ├── server/
│   │   └── grpc_server.go
│   ├── client/
│   │   └── user_client.go
│   ├── service/
│   │   └── user_service.go
│   └── config/
│       └── config.go
├── cmd/
│   └── main.go
├── go.mod
└── Dockerfile
```

### 2. 错误处理

- 使用 `status.Errorf` 返回明确错误码
- 在拦截器中统一处理错误
- 使用错误详情（Details）

### 3. 性能优化

- 使用流式 RPC 处理大数据
- 合理设置消息大小限制
- 使用连接池复用连接

### 4. 安全性

- 使用 TLS 加密通信
- 实现认证拦截器
- 验证所有输入
