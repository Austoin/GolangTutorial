// testing_example.go
// Go 测试指南 - 详细注释版

/*
Go 测试框架是内置的，无需额外安装。

测试文件命名规则：
  - 以 _test.go 结尾
  - 与被测试文件在同一包

测试函数命名规则：
  - 以 Test 开头
  - 参数为 *testing.T

运行测试：
  go test ./...              # 运行所有测试
  go test -v                 # 详细输出
  go test -run TestName      # 运行指定测试
  go test -cover             # 显示覆盖率
  go test -bench=.           # 运行基准测试
*/

package main

import (
	"errors"
	"fmt"
	"math"
	"testing"
)

// ====== 被测试的代码 ======

// Calculator 计算器
type Calculator struct{}

// Add 加法
func (c *Calculator) Add(a, b float64) float64 {
	return a + b
}

// Subtract 减法
func (c *Calculator) Subtract(a, b float64) float64 {
	return a - b
}

// Multiply 乘法
func (c *Calculator) Multiply(a, b float64) float64 {
	return a * b
}

// Divide 除法
func (c *Calculator) Divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("除数不能为零")
	}
	return a / b, nil
}

// Sqrt 平方根
func (c *Calculator) Sqrt(a float64) (float64, error) {
	if a < 0 {
		return 0, errors.New("不能对负数开方")
	}
	return math.Sqrt(a), nil
}

// NewCalculator 创建计算器
func NewCalculator() *Calculator {
	return &Calculator{}
}

// ====== 单元测试 ======

// TestCalculator_Add 测试加法
func TestCalculator_Add(t *testing.T) {
	c := NewCalculator()

	tests := []struct {
		name     string
		a, b     float64
		expected float64
	}{
		{"正数相加", 1, 2, 3},
		{"负数相加", -1, -2, -3},
		{"正负相加", 1, -2, -1},
		{"小数相加", 0.1, 0.2, 0.3},
		{"零相加", 5, 0, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := c.Add(tt.a, tt.b)
			// 比较浮点数需要考虑精度
			if diff := math.Abs(result - tt.expected); diff > 1e-9 {
				t.Errorf("Add(%f, %f) = %f, 期望 %f", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

// TestCalculator_Divide 测试除法
func TestCalculator_Divide(t *testing.T) {
	c := NewCalculator()

	// 测试正常除法
	t.Run("正常除法", func(t *testing.T) {
		result, err := c.Divide(10, 2)
		if err != nil {
			t.Errorf("不应返回错误: %v", err)
		}
		if result != 5 {
			t.Errorf("Divide(10, 2) = %f, 期望 5", result)
		}
	})

	// 测试除零错误
	t.Run("除零错误", func(t *testing.T) {
		_, err := c.Divide(10, 0)
		if err == nil {
			t.Error("应该返回错误")
		}
		if err.Error() != "除数不能为零" {
			t.Errorf("错误信息不正确: %s", err.Error())
		}
	})
}

// TestCalculator_Sqrt 测试平方根
func TestCalculator_Sqrt(t *testing.T) {
	c := NewCalculator()

	tests := []struct {
		name     string
		input    float64
		expected float64
		wantErr  bool
	}{
		{"正数平方根", 4, 2, false},
		{"零的平方根", 0, 0, false},
		{"一", 1, 1, false},
		{"负数错误", -1, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := c.Sqrt(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Sqrt(%f) 错误 = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr && math.Abs(result-tt.expected) > 1e-9 {
				t.Errorf("Sqrt(%f) = %f, 期望 %f", tt.input, result, tt.expected)
			}
		})
	}
}

// ====== 表驱动测试 ======

// TestCalculator_Subtract 表驱动测试示例
func TestCalculator_Subtract(t *testing.T) {
	c := NewCalculator()

	tests := []struct {
		name     string
		a, b     float64
		expected float64
	}{
		{"基本减法", 5, 3, 2},
		{"负数结果", 3, 5, -2},
		{"零减法", 7, 0, 7},
		{"小数减法", 5.5, 2.2, 3.3},
	}

	for _, tt := range tests {
		// t.Run 创建子测试
		tt := tt // 捕获循环变量
		t.Run(tt.name, func(t *testing.T) {
			result := c.Subtract(tt.a, tt.b)
			if math.Abs(result-tt.expected) > 1e-9 {
				t.Errorf("Subtract(%f, %f) = %f, 期望 %f", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

// ====== 基准测试 ======

// BenchmarkCalculator_Add 基准测试加法
func BenchmarkCalculator_Add(b *testing.B) {
	c := NewCalculator()
	a, b := 1.0, 2.0

	b.ResetTimer() // 重置计时器
	for i := 0; i < b.N; i++ {
		_ = c.Add(a, b)
	}
}

// BenchmarkCalculator_Multiply 基准测试乘法
func BenchmarkCalculator_Multiply(b *testing.B) {
	c := NewCalculator()
	a, b := 123.456, 789.012

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = c.Multiply(a, b)
	}
}

// BenchmarkCalculator_Divide 基准测试除法
func BenchmarkCalculator_Divide(b *testing.B) {
	c := NewCalculator()
	a, b := 100.0, 3.0

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = c.Divide(a, b)
	}
}

// ====== 示例测试 ======

// ExampleAdd 加法示例测试
// 以 Example 开头的函数会被当作示例测试
// 用于展示函数用法
func ExampleAdd() {
	fmt.Println(1 + 2)
	// Output:
	// 3
}

// ExampleCalculator_Add 计算器 Add 方法示例
func ExampleCalculator_Add() {
	c := NewCalculator()
	fmt.Println(c.Add(10, 5))
	// Output:
	// 15
}

// ExampleCalculator_Divide 除法示例
func ExampleCalculator_Divide() {
	c := NewCalculator()
	result, _ := c.Divide(20, 4)
	fmt.Println(result)
	// Output:
	// 5
}

// ====== Main 测试 ======

// TestMain 主测试函数
// 包含 TestMain 的测试文件会先执行 TestMain
// 常用于设置和清理测试环境
func TestMain(m *testing.M) {
	fmt.Println("测试开始前的设置...")

	// 运行所有测试
	code := m.Run()

	fmt.Println("测试完成后的清理...")

	// 退出测试
	// os.Exit(code)
}

// ====== Mock 测试 ======

// Database 接口（用于 Mock）
type Database interface {
	GetUser(id int) (User, error)
	CreateUser(user User) error
}

// User 用户模型
type User struct {
	ID       int
	Username string
	Email    string
}

// MockDatabase Mock 数据库实现
type MockDatabase struct {
	users map[int]User
}

func NewMockDatabase() *MockDatabase {
	return &MockDatabase{
		users: make(map[int]User),
	}
}

func (m *MockDatabase) GetUser(id int) (User, error) {
	user, ok := m.users[id]
	if !ok {
		return User{}, errors.New("用户不存在")
	}
	return user, nil
}

func (m *MockDatabase) CreateUser(user User) error {
	m.users[user.ID] = user
	return nil
}

// UserService 用户服务
type UserService struct {
	db Database
}

func NewUserService(db Database) *UserService {
	return &UserService{db: db}
}

func (s *UserService) GetUserByID(id int) (User, error) {
	return s.db.GetUser(id)
}

// TestUserService_GetUserByID 使用 Mock 测试
func TestUserService_GetUserByID(t *testing.T) {
	mockDB := NewMockDatabase()

	// 添加测试数据
	mockDB.CreateUser(User{ID: 1, Username: "alice", Email: "alice@example.com"})

	service := NewUserService(mockDB)

	tests := []struct {
		name    string
		userID  int
		wantErr bool
	}{
		{"用户存在", 1, false},
		{"用户不存在", 999, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.GetUserByID(tt.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserByID() 错误 = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// ====== 子测试和跳过测试 ======

// TestWithSubTests 演示子测试
func TestWithSubTests(t *testing.T) {
	// 跳过测试示例
	t.Skip("这个测试暂时跳过")

	t.Run("子测试1", func(t *testing.T) {
		if true {
			t.Skip("跳过子测试1")
		}
	})

	t.Run("子测试2", func(t *testing.T) {
		t.Parallel() // 标记为可并行执行
		// 测试代码
	})
}

// ====== 测试覆盖率 ======

/*
测试覆盖率命令：
  go test -cover              # 显示覆盖率
  go test -coverprofile=c.out # 输出覆盖率文件
  go tool cover -func=c.out   # 查看函数覆盖率
  go tool cover -html=c.out   # 生成 HTML 报告

理想覆盖率：
  - 单元测试：80%+
  - 关键路径：100%
  - 新功能：先写测试，再写代码（TDD）
*/

// ====== 表格测试辅助函数 ======

// Sum 求和函数（用于测试）
func Sum(numbers []int) int {
	total := 0
	for _, n := range numbers {
		total += n
	}
	return total
}

// TestSum 表格驱动测试
func TestSum(t *testing.T) {
	tests := []struct {
		name     string
		numbers  []int
		expected int
	}{
		{"空切片", []int{}, 0},
		{"单个元素", []int{5}, 5},
		{"多个元素", []int{1, 2, 3, 4, 5}, 15},
		{"负数", []int{-1, 1}, 0},
		{"重复数字", []int{2, 2, 2}, 6},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Sum(tt.numbers); got != tt.expected {
				t.Errorf("Sum(%v) = %d, 期望 %d", tt.numbers, got, tt.expected)
			}
		})
	}
}

// ====== 临时文件/目录测试 ======

/*
使用 t.TempDir() 和 t.Cleanup()：
  func TestWithTempDir(t *testing.T) {
      dir := t.TempDir()
      // 使用 dir 进行测试

      t.Cleanup(func() {
          // 清理代码（可选）
      })
  }
*/

// ====== 错误处理测试 ======

// ValidateEmail 验证邮箱（用于测试）
func ValidateEmail(email string) error {
	if len(email) < 3 {
		return errors.New("邮箱太短")
	}
	if email[0] == '@' || email[len(email)-1] == '@' {
		return errors.New("邮箱格式错误")
	}
	return nil
}

// TestValidateEmail 验证邮箱测试
func TestValidateEmail(t *testing.T) {
	tests := []struct {
		email   string
		wantErr bool
		errMsg  string
	}{
		{"", true, "邮箱太短"},
		{"a@b", false, ""},
		{"@test.com", true, "邮箱格式错误"},
		{"test@", true, "邮箱格式错误"},
		{"valid@example.com", false, ""},
	}

	for _, tt := range tests {
		t.Run(tt.email, func(t *testing.T) {
			err := ValidateEmail(tt.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateEmail(%s) 错误 = %v, wantErr %v", tt.email, err, tt.wantErr)
			}
			if tt.wantErr && err.Error() != tt.errMsg {
				t.Errorf("ValidateEmail(%s) 错误信息 = %s, 期望 %s", tt.email, err.Error(), tt.errMsg)
			}
		})
	}
}

// ====== 运行测试的示例 ======

/*
在命令行运行：
  go test -v                          # 详细模式
  go test -run TestSum                # 运行指定测试
  go test -run "TestCalculator/正数"  # 运行子测试
  go test -skip "TestWithSubTests"    # 跳过测试
  go test -count=3                    # 运行多次
  go test -timeout 30s                # 设置超时
  go test -cover                      # 显示覆盖率
  go test -coverprofile=cover.out     # 输出覆盖率报告

基准测试：
  go test -bench=.                    # 运行所有基准测试
  go test -bench=BenchmarkCalculator  # 运行指定基准测试
  go test -benchmem                   # 显示内存分配
  go test -benchtime=5s               # 设置基准测试时间
*/
