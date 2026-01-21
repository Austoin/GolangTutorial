# Go 测试指南

## 目录
- [1. 测试概述](#1-测试概述)
- [2. 单元测试](#2-单元测试)
- [3. 基准测试](#3-基准测试)
- [4. Mock 技术](#4-mock-技术)
- [5. 集成测试](#5-集成测试)
- [6. 测试最佳实践](#6-测试最佳实践)

---

## 1. 测试概述

### 1.1 Go 测试框架

Go 内置了强大的测试框架，位于 `testing` 包中。

```go
import "testing"

// 测试函数格式
func TestXxx(*testing.T) {}

// 基准测试函数格式
func BenchmarkXxx(*testing.B) {}

// 示例函数格式
func ExampleXxx() {}
```

### 1.2 测试命令

```bash
# 运行所有测试
go test -v ./...

# 运行指定测试
go test -v -run TestFunctionName

# 跳过测试
go test -v -skip TestFunctionName

# 运行基准测试
go test -bench=. -benchmem

# 运行特定基准测试
go test -bench=BenchmarkFunctionName

# 生成测试覆盖率
go test -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html

# 运行测试并显示覆盖率
go test -cover

# 测试超时
go test -timeout 30s

# 并行测试
go test -parallel 4

# 竞赛检测
go test -race

# 内存泄漏检测
go test -msan

# 详细输出
go test -v

# 重试失败的测试
go test -count=3
```

---

## 2. 单元测试

### 2.1 基础测试

```go
// math.go
package math

// Add 加法
func Add(a, b int) int {
    return a + b
}

// Subtract 减法
func Subtract(a, b int) int {
    return a - b
}

// Multiply 乘法
func Multiply(a, b int) int {
    return a * b
}

// Divide 除法
func Divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, fmt.Errorf("除数不能为零")
    }
    return a / b, nil
}
```

```go
// math_test.go
package math

import (
    "testing"
)

func TestAdd(t *testing.T) {
    tests := []struct {
        name     string
        a, b     int
        expected int
    }{
        {"正数相加", 1, 2, 3},
        {"负数相加", -1, -2, -3},
        {"正负相加", 5, -3, 2},
        {"零相加", 0, 5, 5},
        {"大数相加", 1000000, 2000000, 3000000},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := Add(tt.a, tt.b)
            if result != tt.expected {
                t.Errorf("Add(%d, %d) = %d, want %d", 
                    tt.a, tt.b, result, tt.expected)
            }
        })
    }
}

func TestSubtract(t *testing.T) {
    tests := []struct {
        name     string
        a, b     int
        expected int
    }{
        {"基本减法", 10, 5, 5},
        {"负数减法", -5, -3, -2},
        {"结果为负", 3, 8, -5},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := Subtract(tt.a, tt.b)
            if result != tt.expected {
                t.Errorf("Subtract(%d, %d) = %d, want %d", 
                    tt.a, tt.b, result, tt.expected)
            }
        })
    }
}

func TestMultiply(t *testing.T) {
    tests := []struct {
        name     string
        a, b     int
        expected int
    }{
        {"基本乘法", 3, 4, 12},
        {"零乘法", 0, 100, 0},
        {"负数乘法", -2, 3, -6},
        {"负数相乘", -2, -3, 6},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := Multiply(tt.a, tt.b)
            if result != tt.expected {
                t.Errorf("Multiply(%d, %d) = %d, want %d", 
                    tt.a, tt.b, result, tt.expected)
            }
        })
    }
}

func TestDivide(t *testing.T) {
    tests := []struct {
        name        string
        a, b        float64
        expected    float64
        expectError bool
    }{
        {"基本除法", 10, 2, 5, false},
        {"小数除法", 7, 2, 3.5, false},
        {"零除", 10, 0, 0, true},
        {"负数除法", -10, 2, -5, false},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := Divide(tt.a, tt.b)
            
            if tt.expectError {
                if err == nil {
                    t.Errorf("Divide(%f, %f) expected error, got %f", 
                        tt.a, tt.b, result)
                }
            } else {
                if err != nil {
                    t.Errorf("Divide(%f, %f) unexpected error: %v", 
                        tt.a, tt.b, err)
                }
                if result != tt.expected {
                    t.Errorf("Divide(%f, %f) = %f, want %f", 
                        tt.a, tt.b, result, tt.expected)
                }
            }
        })
    }
}
```

### 2.2 表驱动测试

```go
// string_test.go
package stringutils

import (
    "testing"
)

func TestToUpper(t *testing.T) {
    tests := []struct {
        input    string
        expected string
    }{
        {"hello", "HELLO"},
        {"world", "WORLD"},
        {"", ""},
        {"Go", "GO"},
        {"go123", "GO123"},
    }

    for _, tt := range tests {
        result := ToUpper(tt.input)
        if result != tt.expected {
            t.Errorf("ToUpper(%q) = %q, want %q", 
                tt.input, result, tt.expected)
        }
    }
}

func TestTruncate(t *testing.T) {
    tests := []struct {
        input    string
        length   int
        expected string
        ellipsis string
    }{
        {"hello", 10, "hello", "..."},
        {"hello world", 5, "hello", "..."},
        {"hi", 5, "hi", "..."},
        {"", 5, "", "..."},
    }

    for _, tt := range tests {
        result := Truncate(tt.input, tt.length)
        if result != tt.expected {
            t.Errorf("Truncate(%q, %d) = %q, want %q", 
                tt.input, tt.length, result, tt.expected)
        }
    }
}
```

### 2.3 错误处理测试

```go
// validator_test.go
package validator

import (
    "testing"
)

func TestValidateEmail(t *testing.T) {
    tests := []struct {
        name    string
        email   string
        valid   bool
    }{
        {"有效邮箱", "user@example.com", true},
        {"有效邮箱2", "user.name@example.co.uk", true},
        {"无效邮箱", "invalid-email", false},
        {"缺少@", "userexample.com", false},
        {"缺少域名", "user@", false},
        {"空字符串", "", false},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateEmail(tt.email)
            isValid := err == nil
            
            if isValid != tt.valid {
                if tt.valid {
                    t.Errorf("ValidateEmail(%q) expected valid, got error: %v", 
                        tt.email, err)
                } else {
                    t.Errorf("ValidateEmail(%q) expected invalid, got valid", 
                        tt.email)
                }
            }
        })
    }
}

func TestValidatePassword(t *testing.T) {
    tests := []struct {
        name        string
        password    string
        minLength   int
        requireDigit bool
        expectError bool
    }{
        {"有效密码", "Password123", 8, true, false},
        {"太短", "Pass1", 8, true, true},
        {"无数字", "Password", 8, true, true},
        {"最小长度", "Pass1234", 8, true, false},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidatePassword(tt.password, tt.minLength, tt.requireDigit)
            hasError := err != nil
            
            if hasError != tt.expectError {
                t.Errorf("ValidatePassword() error = %v, expectError = %v", 
                    err, tt.expectError)
            }
        })
    }
}
```

### 2.4 表格测试（Table-Driven Tests）扩展

```go
// slice_test.go
package sliceutils

import (
    "reflect"
    "testing"
)

func TestFilter(t *testing.T) {
    tests := []struct {
        name     string
        input    []int
        predicate func(int) bool
        expected []int
    }{
        {
            name:     "过滤偶数",
            input:    []int{1, 2, 3, 4, 5, 6},
            predicate: func(n int) bool { return n%2 == 1 },
            expected: []int{1, 3, 5},
        },
        {
            name:     "过滤大于10",
            input:    []int{5, 10, 15, 20},
            predicate: func(n int) bool { return n > 10 },
            expected: []int{15, 20},
        },
        {
            name:     "空切片",
            input:    []int{},
            predicate: func(n int) bool { return n > 0 },
            expected: []int{},
        },
        {
            name:     "无匹配",
            input:    []int{1, 2, 3},
            predicate: func(n int) bool { return n > 100 },
            expected: []int{},
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := Filter(tt.input, tt.predicate)
            if !reflect.DeepEqual(result, tt.expected) {
                t.Errorf("Filter() = %v, want %v", result, tt.expected)
            }
        })
    }
}

func TestMap(t *testing.T) {
    tests := []struct {
        name     string
        input    []int
        mapper   func(int) string
        expected []string
    }{
        {
            name:     "数字转字符串",
            input:    []int{1, 2, 3},
            mapper:   func(i int) string { return string(rune('a' + i - 1)) },
            expected: []string{"a", "b", "c"},
        },
        {
            name:     "乘2转字符串",
            input:    []int{2, 4, 6},
            mapper:   func(i int) string { return fmt.Sprintf("%d", i*2) },
            expected: []string{"4", "8", "12"},
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := Map(tt.input, tt.mapper)
            if !reflect.DeepEqual(result, tt.expected) {
                t.Errorf("Map() = %v, want %v", result, tt.expected)
            }
        })
    }
}
```

### 2.5 子测试

```go
// calculator_test.go
package calculator

import (
    "testing"
)

func TestCalculator(t *testing.T) {
    calc := NewCalculator()

    t.Run("Add", func(t *testing.T) {
        tests := []struct{ a, b, expected int }{
            {1, 2, 3},
            {0, 0, 0},
            {-1, 1, 0},
        }
        for _, tt := range tests {
            result := calc.Add(tt.a, tt.b)
            if result != tt.expected {
                t.Errorf("Add(%d, %d) = %d, want %d", 
                    tt.a, tt.b, result, tt.expected)
            }
        }
    })

    t.Run("Multiply", func(t *testing.T) {
        tests := []struct{ a, b, expected int }{
            {3, 4, 12},
            {0, 100, 0},
            {-2, 3, -6},
        }
        for _, tt := range tests {
            result := calc.Multiply(tt.a, tt.b)
            if result != tt.expected {
                t.Errorf("Multiply(%d, %d) = %d, want %d", 
                    tt.a, tt.b, result, tt.expected)
            }
        }
    })
}
```

---

## 3. 基准测试

### 3.1 基础基准测试

```go
// fib_test.go
package fib

// Fib 递归实现
func Fib(n int) int {
    if n <= 1 {
        return n
    }
    return Fib(n-1) + Fib(n-2)
}

// FibDP 动态规划实现
func FibDP(n int) int {
    if n <= 1 {
        return n
    }
    a, b := 0, 1
    for i := 2; i <= n; i++ {
        a, b = b, a+b
    }
    return b
}
```

```go
// fib_benchmark_test.go
package fib

import (
    "testing"
)

func BenchmarkFibRecursive(b *testing.B) {
    // b.N 会根据函数运行时间自动调整
    for i := 0; i < b.N; i++ {
        Fib(20) // 计算第 20 个斐波那契数
    }
}

func BenchmarkFibDP(b *testing.B) {
    for i := 0; i < b.N; i++ {
        FibDP(20)
    }
}

func BenchmarkFibRecursive10(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Fib(10)
    }
}

func BenchmarkFibDP10(b *testing.B) {
    for i := 0; i < b.N; i++ {
        FibDP(10)
    }
}
```

```bash
# 运行基准测试
go test -bench=. -benchmem

# 输出示例
# BenchmarkFibRecursive-8    1000   1,234,567 ns/op   0 B/op   0 allocs/op
# BenchmarkFibDP-8           1000000    35 ns/op        0 B/op   0 allocs/op
# BenchmarkFibRecursive10-8  10000    15,678 ns/op     0 B/op   0 allocs/op
# BenchmarkFibDP10-8         1000000    28 ns/op        0 B/op   0 allocs/op
```

### 3.2 内存分配测试

```go
// strings_test.go
package strings

import (
    "testing"
)

func BenchmarkConcatPlus(b *testing.B) {
    b.ReportAllocs()
    for i := 0; i < b.N; i++ {
        result := ""
        for j := 0; j < 100; j++ {
            result += "a"
        }
        _ = result
    }
}

func BenchmarkConcatBuilder(b *testing.B) {
    b.ReportAllocs()
    for i := 0; i < b.N; i++ {
        var builder strings.Builder
        for j := 0; j < 100; j++ {
            builder.WriteString("a")
        }
        _ = builder.String()
    }
}

func BenchmarkConcatSlice(b *testing.B) {
    b.ReportAllocs()
    for i := 0; i < b.N; i++ {
        slices := make([]string, 100)
        for j := 0; j < 100; j++ {
            slices[j] = "a"
        }
        result := strings.Join(slices, "")
        _ = result
    }
}
```

### 3.3 并发基准测试

```go
// channel_test.go
package channel

import (
    "sync"
    "testing"
)

func BenchmarkChannelUnbuffered(b *testing.B) {
    b.ReportAllocs()
    b.ResetTimer()
    
    for i := 0; i < b.N; i++ {
        ch := make(chan int)
        var wg sync.WaitGroup
        wg.Add(2)
        
        go func() {
            ch <- 1
            wg.Done()
        }()
        
        go func() {
            <-ch
            wg.Done()
        }()
        
        wg.Wait()
    }
}

func BenchmarkChannelBuffered(b *testing.B) {
    b.ReportAllocs()
    b.ResetTimer()
    
    for i := 0; i < b.N; i++ {
        ch := make(chan int, 1)
        var wg sync.WaitGroup
        wg.Add(2)
        
        go func() {
            ch <- 1
            wg.Done()
        }()
        
        go func() {
            <-ch
            wg.Done()
        }()
        
        wg.Wait()
    }
}

func BenchmarkMutex(b *testing.B) {
    var mu sync.Mutex
    b.ReportAllocs()
    b.ResetTimer()
    
    for i := 0; i < b.N; i++ {
        mu.Lock()
        _ = i
        mu.Unlock()
    }
}

func BenchmarkRWMutex(b *testing.B) {
    var mu sync.RWMutex
    b.ReportAllocs()
    b.ResetTimer()
    
    for i := 0; i < b.N; i++ {
        mu.RLock()
        _ = i
        mu.RUnlock()
    }
}
```

### 3.4 子基准测试

```go
// sort_benchmark_test.go
package sort

import (
    "math/rand"
    "testing"
)

func generateSlice(n int) []int {
    slice := make([]int, n)
    for i := range slice {
        slice[i] = rand.Intn(n)
    }
    return slice
}

func BenchmarkSort(b *testing.B) {
    sizes := []int{100, 1000, 10000}
    
    for _, size := range sizes {
        b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
            b.ReportAllocs()
            b.ResetTimer()
            
            for i := 0; i < b.N; i++ {
                slice := generateSlice(size)
                BubbleSort(slice)
            }
        })
    }
}

func BenchmarkSortParallel(b *testing.B) {
    sizes := []int{100, 1000, 10000}
    
    for _, size := range sizes {
        b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
            b.ReportAllocs()
            b.ResetTimer()
            b.SetParallelism(4)
            
            for i := 0; i < b.N; i++ {
                slice := generateSlice(size)
                ParallelSort(slice)
            }
        })
    }
}
```

---

## 4. Mock 技术

### 4.1 接口 Mock

```go
// repository.go
package repository

// UserRepository 用户仓储接口
type UserRepository interface {
    GetByID(id int64) (*User, error)
    GetByEmail(email string) (*User, error)
    Create(user *User) error
    Update(user *User) error
    Delete(id int64) error
}

// User 用户模型
type User struct {
    ID    int64
    Email string
    Name  string
}
```

```go
// mock_repository.go
package repository

import (
    "sync"
)

// MockUserRepository Mock 实现
type MockUserRepository struct {
    users map[int64]*User
    mu    sync.RWMutex
    nextID int64
}

func NewMockUserRepository() *MockUserRepository {
    return &MockUserRepository{
        users:  make(map[int64]*User),
        nextID: 1,
    }
}

func (m *MockUserRepository) GetByID(id int64) (*User, error) {
    m.mu.RLock()
    defer m.mu.RUnlock()
    
    user, ok := m.users[id]
    if !ok {
        return nil, ErrNotFound
    }
    return user, nil
}

func (m *MockUserRepository) GetByEmail(email string) (*User, error) {
    m.mu.RLock()
    defer m.mu.RUnlock()
    
    for _, user := range m.users {
        if user.Email == email {
            return user, nil
        }
    }
    return nil, ErrNotFound
}

func (m *MockUserRepository) Create(user *User) error {
    m.mu.Lock()
    defer m.mu.Unlock()
    
    user.ID = m.nextID
    m.nextID++
    m.users[user.ID] = user
    return nil
}

func (m *MockUserRepository) Update(user *User) error {
    m.mu.Lock()
    defer m.mu.Unlock()
    
    if _, ok := m.users[user.ID]; !ok {
        return ErrNotFound
    }
    m.users[user.ID] = user
    return nil
}

func (m *MockUserRepository) Delete(id int64) error {
    m.mu.Lock()
    defer m.mu.Unlock()
    
    if _, ok := m.users[id]; !ok {
        return ErrNotFound
    }
    delete(m.users, id)
    return nil
}

// 错误定义
var ErrNotFound = fmt.Errorf("user not found")
```

### 4.2 使用 mockery 生成 Mock

```bash
# 安装 mockery
go install github.com/vektra/mockery/v2@latest

# 生成配置
cat > .mockery.yaml >
name: "UserRepository"
interface: "UserRepository"
outpkg: "repository"
dir: "repository"
```

```go
// user_service_test.go
package service

import (
    "testing"

    "github.com/example/user-service/repository"
)

func TestUserService_GetUser(t *testing.T) {
    // 创建 Mock
    mockRepo := repository.NewMockUserRepository()
    
    // 准备测试数据
    testUser := &repository.User{
        ID:    1,
        Email: "test@example.com",
        Name:  "Test User",
    }
    mockRepo.Create(testUser)
    
    // 创建服务（注入 Mock）
    service := NewUserService(mockRepo)
    
    // 测试
    user, err := service.GetUser(1)
    if err != nil {
        t.Fatalf("GetUser() error = %v", err)
    }
    
    if user.ID != testUser.ID {
        t.Errorf("GetUser() ID = %d, want %d", user.ID, testUser.ID)
    }
}

func TestUserService_GetUser_NotFound(t *testing.T) {
    mockRepo := repository.NewMockUserRepository()
    service := NewUserService(mockRepo)
    
    _, err := service.GetUser(999)
    if err == nil {
        t.Error("GetUser() expected error, got nil")
    }
}
```

### 4.3 使用 testify/mock

```go
// 安装
go get github.com/stretchr/testify
```

```go
// service_test.go
package service

import (
    "testing"
    "time"

    "github.com/stretchr/testify/mock"
)

// MockUserRepository 使用 testify/mock
type MockUserRepository struct {
    mock.Mock
}

func (m *MockUserRepository) GetByID(id int64) (*User, error) {
    args := m.Called(id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*User), args.Error(1)
}

func (m *MockUserRepository) Create(user *User) error {
    args := m.Called(user)
    return args.Error(0)
}

func TestUserService_GetUser_WithTestify(t *testing.T) {
    mockRepo := new(MockUserRepository)
    
    expectedUser := &User{
        ID:        1,
        Email:     "test@example.com",
        Name:      "Test User",
        CreatedAt: time.Now(),
    }
    
    // 设置预期行为
    mockRepo.On("GetByID", int64(1)).Return(expectedUser, nil)
    
    service := NewUserService(mockRepo)
    user, err := service.GetUser(1)
    
    // 断言
    if err != nil {
        t.Fatalf("GetUser() error = %v", err)
    }
    
    mockRepo.AssertExpectations(t)
    mockRepo.AssertNumberOfCalls(t, "GetByID", 1)
}
```

---

## 5. 集成测试

### 5.1 数据库集成测试

```go
// database_test.go
package database

import (
    "database/sql"
    "testing"

    _ "github.com/go-sql-driver/mysql"
    "github.com/testcontainers/testcontainers-go"
    "github.com/testcontainers/testcontainers-go/modules/mysql"
)

func TestUserRepository_Integration(t *testing.T) {
    // 使用 testcontainers 启动 MySQL
    ctx := context.Background()
    
    mysqlContainer, err := mysql.RunContainer(ctx,
        mysql.WithDatabase("testdb"),
        mysql.WithUsername("test"),
        mysql.WithPassword("test"),
    )
    if err != nil {
        t.Fatal(err)
    }
    defer mysqlContainer.Terminate(ctx)
    
    // 获取连接字符串
    connStr, err := mysqlContainer.ConnectionString(ctx)
    if err != nil {
        t.Fatal(err)
    }
    
    // 连接数据库
    db, err := sql.Open("mysql", connStr)
    if err != nil {
        t.Fatal(err)
    }
    defer db.Close()
    
    // 运行测试
    repository := NewUserRepository(db)
    
    t.Run("CreateAndGet", func(t *testing.T) {
        user := &User{
            Email: "test@example.com",
            Name:  "Test User",
        }
        
        // 创建
        err := repository.Create(user)
        if err != nil {
            t.Fatalf("Create() error = %v", err)
        }
        
        if user.ID == 0 {
            t.Error("Create() user.ID should be set")
        }
        
        // 获取
        fetched, err := repository.GetByID(user.ID)
        if err != nil {
            t.Fatalf("GetByID() error = %v", err)
        }
        
        if fetched.Email != user.Email {
            t.Errorf("GetByID() Email = %s, want %s", 
                fetched.Email, user.Email)
        }
    })
}
```

### 5.2 HTTP 集成测试

```go
// handler_test.go
package handler

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gin-gonic/gin"
)

func TestHealthCheck(t *testing.T) {
    // 设置路由
    gin.SetMode(gin.TestMode)
    router := gin.New()
    router.GET("/health", HealthCheck)
    
    // 创建测试请求
    req, _ := http.NewRequest("GET", "/health", nil)
    w := httptest.NewRecorder()
    
    // 执行请求
    router.ServeHTTP(w, req)
    
    // 断言
    if w.Code != http.StatusOK {
        t.Errorf("HealthCheck() status = %d, want %d", 
            w.Code, http.StatusOK)
    }
    
    var response map[string]string
    json.Unmarshal(w.Body.Bytes(), &response)
    
    if response["status"] != "ok" {
        t.Errorf("HealthCheck() status = %s, want ok", 
            response["status"])
    }
}

func TestCreateUser(t *testing.T) {
    gin.SetMode(gin.TestMode)
    router := gin.New()
    
    // 注入 Mock 服务
    mockService := NewMockUserService()
    handler := NewUserHandler(mockService)
    router.POST("/users", handler.CreateUser)
    
    t.Run("ValidRequest", func(t *testing.T) {
        user := CreateUserRequest{
            Name:  "Test User",
            Email: "test@example.com",
        }
        body, _ := json.Marshal(user)
        
        req, _ := http.NewRequest("POST", "/users", 
            bytes.NewBuffer(body))
        req.Header.Set("Content-Type", "application/json")
        w := httptest.NewRecorder()
        
        router.ServeHTTP(w, req)
        
        if w.Code != http.StatusCreated {
            t.Errorf("CreateUser() status = %d, want %d", 
                w.Code, http.StatusCreated)
        }
    })
    
    t.Run("InvalidRequest", func(t *testing.T) {
        invalid := map[string]string{
            "name": "", // 缺少必填字段
        }
        body, _ := json.Marshal(invalid)
        
        req, _ := http.NewRequest("POST", "/users", 
            bytes.NewBuffer(body))
        req.Header.Set("Content-Type", "application/json")
        w := httptest.NewRecorder()
        
        router.ServeHTTP(w, req)
        
        if w.Code != http.StatusBadRequest {
            t.Errorf("CreateUser() status = %d, want %d", 
                w.Code, http.StatusBadRequest)
        }
    })
}
```

### 5.3 gRPC 集成测试

```go
// grpc_test.go
package grpc

import (
    "context"
    "testing"

    "google.golang.org/grpc"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
    "google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func bufDialer(context.Context, string) (net.Conn, error) {
    return lis.Dial()
}

func TestUserService(t *testing.T) {
    // 创建内存服务器
    lis = bufconn.Listen(bufSize)
    server := grpc.NewServer()
    RegisterUserServiceServer(server, &UserServer{})
    
    go func() {
        if err := server.Serve(lis); err != nil {
            t.Fatalf("Server exited with error: %v", err)
        }
    }()
    
    // 创建客户端
    ctx := context.Background()
    conn, err := grpc.DialContext(ctx, "bufnet", 
        grpc.WithBlock(),
        grpc.WithInsecure(),
        grpc.WithContextDialer(bufDialer),
    )
    if err != nil {
        t.Fatalf("Failed to dial: %v", err)
    }
    defer conn.Close()
    
    client := NewUserServiceClient(conn)
    
    t.Run("GetUser", func(t *testing.T) {
        // 创建测试用户
        createResp, err := client.CreateUser(ctx, &CreateUserRequest{
            Username: "test",
            Email:    "test@example.com",
        })
        if err != nil {
            t.Fatalf("CreateUser() error = %v", err)
        }
        
        // 获取用户
        getResp, err := client.GetUser(ctx, &GetUserRequest{
            Id: createResp.User.Id,
        })
        if err != nil {
            t.Fatalf("GetUser() error = %v", err)
        }
        
        if getResp.User.Username != "test" {
            t.Errorf("GetUser() username = %s, want test", 
                getResp.User.Username)
        }
    })
    
    t.Run("GetUser_NotFound", func(t *testing.T) {
        _, err := client.GetUser(ctx, &GetUserRequest{
            Id: 999999,
        })
        
        st, ok := status.FromError(err)
        if !ok {
            t.Fatal("Expected gRPC error")
        }
        
        if st.Code() != codes.NotFound {
            t.Errorf("GetUser() code = %v, want NotFound", 
                st.Code())
        }
    })
    
    server.Stop()
}
```

---

## 6. 测试最佳实践

### 6.1 测试命名规范

```go
// 好的测试命名
func TestUserService_CreateUser_Success(t *testing.T) {}
func TestUserService_CreateUser_ValidationError(t *testing.T) {}
func TestUserRepository_GetByID_NotFound(t *testing.T) {}

// 避免的命名
func Test1(t *testing.T) {}
func TestUser(t *testing.T) {}  // 太笼统
func TestCreate(t *testing.T) {} // 缺少上下文
```

### 6.2 测试组织结构

```
project/
├── cmd/
│   └── main.go
├── internal/
│   ├── service/
│   │   ├── user_service.go
│   │   └── user_service_test.go  # 服务测试
│   ├── handler/
│   │   ├── user_handler.go
│   │   └── user_handler_test.go  # HTTP 测试
│   └── repository/
│       ├── user_repository.go
│       └── user_repository_test.go # 数据库测试
├── pkg/
│   ├── utils/
│   │   ├── string_utils.go
│   │   └── string_utils_test.go
│   └── mock/
│       └── mock.go
├── integration/
│   └── user_flow_test.go  # 端到端测试
└── go.mod
```

### 6.3 测试覆盖率

```bash
# 生成覆盖率报告
go test -coverprofile=coverage.out ./...

# 查看覆盖率
go tool cover -func=coverage.out

# 生成 HTML 报告
go tool cover -html=coverage.out -o coverage.html

# 检查覆盖率变化
go test -coverprofile=new_coverage.out ./...
go diff coverage.out new_coverage.out
```

### 6.4 持续集成中的测试

```yaml
# .github/workflows/test.yml
name: Tests

on:
  push:
    branches: [main, develop]
  pull_request:

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.21'
      
      - name: Cache Go modules
        uses: actions/cache@v3
        with:
          path: |
            ~/.cache/go-build
            ~/go/pkg/mod
          key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
      
      - name: Run unit tests
        run: |
          go test -v -race -cover ./...
      
      - name: Run integration tests
        run: |
          go test -v -tags=integration ./integration/...
      
      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          files: ./coverage.out
```

### 6.5 常见问题

| 问题 | 解决方案 |
|------|---------|
| 测试太慢 | 使用 t.Parallel() 并行测试 |
| 依赖外部服务 | 使用 Mock 或 testcontainers |
| 测试不稳定 | 使用重试机制或增加超时 |
| 代码覆盖率低 | 逐步增加测试用例 |
| 测试数据复杂 | 使用测试工厂或 fixtures |

### 6.6 测试金字塔

```
                    /\
                   /  \
                  /    \
                 / 单元  \
                / 测试    \
               /──────────\
              /            \
             /   集成测试    \
            /                \
           /──────────────────\
          /                    \
         /      端到端测试       \
        /                        \
       /────────────────────────--\
```

- **单元测试**（70%）：快速、独立
- **集成测试**（20%）：验证组件交互
- **端到端测试**（10%）：验证完整流程
