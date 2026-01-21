// basic_syntax/10_error_handling.go
// Go 错误处理详解 - 详细注释版

package main

import (
	"errors"
	"fmt"
	"log"
	"math"
)

// ====== Go 错误处理基础 ======
/*
Go 的错误处理与其他语言不同，它不使用异常（Exception）机制，
而是使用返回值来传递错误。

设计理念：
1. 错误是值 - 错误是普通的 Go 值
2. 显式处理 - 必须显式处理错误
3. 没有异常 - 没有 try-catch-finally
4. 简单直接 - 错误就是简单的接口

错误接口：
  type error interface {
      Error() string
  }
*/

// ====== 错误定义 ======

// 自定义错误类型
// 实现 error 接口的 Error() 方法
type MyError struct {
	Code    int    // 错误码
	Message string // 错误信息
	Details string // 详细信息
}

// Error 方法实现 error 接口
func (e *MyError) Error() string {
	return fmt.Sprintf("[%d] %s: %s", e.Code, e.Message, e.Details)
}

// New 创建新的 MyError
func NewError(code int, message, details string) error {
	return &MyError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

// WrapError 带包装的错误
type WrapError struct {
	Inner error
	Msg   string
}

func (e *WrapError) Error() string {
	return fmt.Sprintf("%s: %v", e.Msg, e.Inner)
}

// Unwrap 解开包装的错误
func (e *WrapError) Unwrap() error {
	return e.Inner
}

// ====== 常见错误创建方式 ======

// 1. 使用 errors.New 创建简单错误
func divide1(a, b float64) (float64, error) {
	if b == 0 {
		// errors.New 创建简单的错误信息
		return 0, errors.New("除数不能为零")
	}
	return a / b, nil
}

// 2. 使用 fmt.Errorf 创建格式化错误
func divide2(a, b float64) (float64, error) {
	if b == 0 {
		// fmt.Errorf 支持格式化输出
		return 0, fmt.Errorf("%.2f / %.2f: 除数不能为零", a, b)
	}
	return a / b, nil
}

// 3. 使用自定义错误类型
func divide3(a, b float64) (float64, error) {
	if b == 0 {
		// 返回自定义错误
		return 0, NewError(1001, "除零错误", fmt.Sprintf("尝试将 %.2f 除以零", a))
	}
	return a / b, nil
}

// ====== 错误检查方式 ======

func checkErrorExamples() {
	// 1. 标准错误检查
	result, err := divide1(10, 0)
	if err != nil {
		log.Printf("错误: %v", err)
	} else {
		fmt.Println("结果:", result)
	}

	// 2. 使用变量捕获错误
	_, err = divide1(5, 0)
	if err != nil {
		fmt.Println("错误信息:", err.Error())
	}

	// 3. 错误比较
	if err == errors.New("除数不能为零") {
		fmt.Println("匹配到特定错误")
	}

	// 4. 错误断言（类型检查）
	if myErr, ok := err.(*MyError); ok {
		fmt.Printf("自定义错误 - 代码: %d, 消息: %s\n", myErr.Code, myErr.Message)
	}

	// 5. 使用 errors.Is 检查包装错误
	err = fmt.Errorf("上层错误: %w", NewError(1001, "底层错误", ""))
	if errors.Is(err, NewError(1001, "", "")) {
		fmt.Println("错误匹配成功")
	}

	// 6. 使用 errors.As 获取错误类型
	var myErr *MyError
	if errors.As(err, &myErr) {
		fmt.Printf("提取错误 - 代码: %d\n", myErr.Code)
	}
}

// ====== 错误包装 ======

// 包装错误保留原始错误信息
func wrapErrorExample() {
	// 1. 使用 %w 包装错误
	original := errors.New("原始错误")
	wrapped := fmt.Errorf("包装错误: %w", original)

	// 2. 解包错误
	if errors.Is(wrapped, original) {
		fmt.Println("可以追溯到原始错误")
	}

	// 3. 多层包装
	level1 := fmt.Errorf("层级1: %w", wrapped)
	level2 := fmt.Errorf("层级2: %w", level1)

	// 4. 使用 errors.Unwrap 解包
	unwrapped := errors.Unwrap(level2)
	_ = unwrapped // 继续解包...

	// 5. 使用 errors.Join 组合多个错误（Go 1.20+）
	multiErr := errors.Join(original, wrapped, level1)
	fmt.Println("组合错误:", multiErr)
}

// ====== 错误码定义 ======

// 错误码常量
const (
	ErrCodeSuccess  = 0
	ErrCodeInvalid  = 1
	ErrCodeNotFound = 2
	ErrCodeConflict = 3
	ErrCodeInternal = 4
)

// 错误码辅助函数
type AppError struct {
	Code    int
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%d] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func NewAppError(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// ====== 哨兵错误（Sentinel Errors）=====

// 预定义的错误值，用于错误比较
var (
	ErrNotFound      = errors.New("资源不存在")
	ErrAlreadyExists = errors.New("资源已存在")
	ErrInvalidInput  = errors.New("无效输入")
	ErrUnauthorized  = errors.New("未授权")
	ErrForbidden     = errors.New("禁止访问")
)

func sentinelErrorExample(data map[string]int, key string) error {
	value, ok := data[key]
	if !ok {
		// 返回哨兵错误
		return fmt.Errorf("获取 %s: %w", key, ErrNotFound)
	}

	// 检查值
	if value < 0 {
		return fmt.Errorf("值不能为负: %w", ErrInvalidInput)
	}

	return nil
}

// ====== 延迟错误处理 ======

func deferredErrorExample() {
	// 在函数结束时统一处理错误
	defer func() {
		if err := recover(); err != nil {
			// 恢复 panic，防止程序崩溃
			log.Printf("捕获 panic: %v", err)
		}
	}()

	// 可能 panic 的代码
	// panic("发生严重错误")
}

// ====== 自定义错误函数 ======

// 验证函数返回错误
type Validator struct {
	Rules []Rule
}

type Rule struct {
	Name  string
	Check func(value interface{}) error
}

func (v *Validator) Validate(value interface{}) error {
	for _, rule := range v.Rules {
		if err := rule.Check(value); err != nil {
			return fmt.Errorf("验证 %s 失败: %w", rule.Name, err)
		}
	}
	return nil
}

// 常见验证规则
var Required = func(field string) Rule {
	return Rule{
		Name: field + "_required",
		Check: func(value interface{}) error {
			if value == nil || value == "" {
				return fmt.Errorf("%s 不能为空", field)
			}
			return nil
		},
	}
}

var MinLength = func(field string, min int) Rule {
	return Rule{
		Name: field + "_min_length",
		Check: func(value interface{}) error {
			str, ok := value.(string)
			if !ok {
				return fmt.Errorf("%s 必须是字符串", field)
			}
			if len(str) < min {
				return fmt.Errorf("%s 长度必须至少 %d", field, min)
			}
			return nil
		},
	}
}

var Range = func(field string, min, max float64) Rule {
	return Rule{
		Name: field + "_range",
		Check: func(value interface{}) error {
			num, ok := value.(float64)
			if !ok {
				return fmt.Errorf("%s 必须是数字", field)
			}
			if num < min || num > max {
				return fmt.Errorf("%s 必须在 %.2f 到 %.2f 之间", field, min, max)
			}
			return nil
		},
	}
}

// ====== 错误处理最佳实践 ======

func bestPractices() {
	// 1. 不要忽略错误
	// result, _ := someFunction() // 避免这样写

	// 2. 添加上下文信息
	if err := someOperation(); err != nil {
		return fmt.Errorf("操作失败: %w", err)
	}

	// 3. 使用命名返回值处理多个错误
	result, err := safeDivide(10, 0)
	if err != nil {
		log.Printf("除法错误: %v", err)
		return
	}
	_ = result

	// 4. 避免错误字符串大写
	// 好的写法: errors.New("invalid input")
	// 不好的写法: errors.New("Invalid Input")

	// 5. 使用专门的错误处理包
	// - github.com/pkg/errors
	// - golang.org/x/xerrors
}

// safeDivide 安全的除法函数
func safeDivide(a, b float64) (result float64, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("除法发生 panic: %v", r)
		}
	}()

	if b == 0 {
		err = errors.New("除数不能为零")
		return
	}

	result = a / b

	// 检查结果
	if math.IsInf(result, 0) || math.IsNaN(result) {
		err = errors.New("无效的计算结果")
		return
	}

	return
}

// ====== 模拟业务错误 ======

type UserNotFoundError struct {
	UserID int
}

func (e *UserNotFoundError) Error() string {
	return fmt.Sprintf("用户不存在: ID=%d", e.UserID)
}

type InvalidOperationError struct {
	Operation string
	Reason    string
}

func (e *InvalidOperationError) Error() string {
	return fmt.Sprintf("无效操作: %s - %s", e.Operation, e.Reason)
}

func businessErrorExample() {
	users := map[int]string{1: "Alice", 2: "Bob"}

	// 模拟用户操作
	getUser := func(id int) (string, error) {
		name, ok := users[id]
		if !ok {
			return "", &UserNotFoundError{UserID: id}
		}
		return name, nil
	}

	// 获取用户
	name, err := getUser(1)
	if err != nil {
		// 处理特定错误
		if notFound, ok := err.(*UserNotFoundError); ok {
			log.Printf("用户不存在: %d", notFound.UserID)
		}
	} else {
		fmt.Println("用户:", name)
	}
}

// ====== 主函数 ======

func main() {
	fmt.Println("=== Go 错误处理详解 ===")

	// 1. 错误创建示例
	fmt.Println("\n--- 错误创建 ---")

	result, err := divide1(10, 2)
	fmt.Printf("divide1(10, 2) = %.2f, err = %v\n", result, err)

	result, err = divide1(10, 0)
	fmt.Printf("divide1(10, 0) = %.2f, err = %v\n", result, err)

	result, err = divide3(10, 0)
	fmt.Printf("divide3(10, 0) = %.2f, err = %v\n", result, err)

	// 2. 错误检查
	fmt.Println("\n--- 错误检查 ---")
	checkErrorExamples()

	// 3. 错误包装
	fmt.Println("\n--- 错误包装 ---")
	wrapErrorExample()

	// 4. 验证示例
	fmt.Println("\n--- 验证示例 ---")
	validator := &Validator{
		Rules: []Rule{
			Required("用户名"),
			MinLength("用户名", 3),
			Range("年龄", 0, 150),
		},
	}

	testCases := []struct {
		Username string
		Age      float64
	}{
		{"", 25},      // 用户名为空
		{"ab", 25},    // 用户名太短
		{"alice", -5}, // 年龄无效
		{"alice", 25}, // 全部正确
	}

	for i, tc := range testCases {
		err := validator.Validate(map[string]interface{}{
			"用户名": tc.Username,
			"年龄":  tc.Age,
		})
		fmt.Printf("测试 %d: 用户名='%s', 年龄=%.0f -> %v\n", i+1, tc.Username, tc.Age, err)
	}

	// 5. 业务错误
	fmt.Println("\n--- 业务错误 ---")
	businessErrorExample()

	// 6. 最佳实践
	fmt.Println("\n--- 最佳实践 ---")
	bestPractices()

	fmt.Println("\n错误处理示例完成")
}
