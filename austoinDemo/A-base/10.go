package mybase

import (
	"errors"
	"fmt"

	// "os"
	"time"
	// "honnef.co/go/tools/config"
)

// 错误处理
func Lesson10() {
	errorDemo()   // error 接口
	panicDemo()   // panic 异常
	recoverDemo() // recover 恢复

	// defer + panic + recover 完整流程
	completeDemo()

	//  errors.Is 错误类型判断
	checkErrorType()

	viliConfig := map[string]string{
		"port": "",
		"host": "",
	}

	validateConfig(viliConfig)
}

func errorDemo() {
	// 创建错误的方式
	// 1. error.New() 创建简单错误
	err1 := errors.New("文件不存在")
	fmt.Printf("err1:%v\n", err1)

	// 2. fmt.Errorf() 创建格式化错误
	filename := "config.json"
	err2 := fmt.Errorf("文件 %s 读取失败", filename)
	fmt.Println("err2:", err2)

	// 函数返回错误
	// 除法函数，返回结果和错误
	result, err := divide(10, 0)
	if err != nil {
		fmt.Println("除法错误:", err)
	} else {
		fmt.Printf("除法结果: %.2f\n", result)
	}

	// 成功的例子
	result2, err := divide(10, 2)
	if err != nil {
		fmt.Printf("错误: %v\n", err)
	} else {
		fmt.Printf("10 / 2 = %.2f\n", result2)
	}

	// 判断错误
	value := -1
	if _, err := validateAge(value); err != nil {
		if errors.Is(err, ErrInvalidAge) {
			fmt.Printf("年龄验证失败: %v\n", err)
		}
	} else {
		fmt.Printf("年龄 %d 验证通过\n", value)
	}
	fmt.Println()
}

// 定义自定义错误变量（包内可见）
var ErrInvalidAge = errors.New("年龄无效")

// 验证年龄
func validateAge(age int) (bool, error) {
	if age < 0 || age > 150 {
		return false, ErrInvalidAge
	}
	return true, nil
}

// // 除法函数
// func divide(a, b float64) (float64, error) {
// 	// 检查除数是否为0
// 	if b == 0 {
// 		// 返回0和错误
// 		return 0, errors.New("除数不能为零")
// 	}
// 	// 返回结果和nil（无错误）
// 	return a / b, nil
// }

// panic 演示
func panicDemo() {

	defer func() {
		// recover() 返回panic的值
		if r := recover(); r != nil {
			fmt.Printf("recover 捕获到 panic: %v\n", r)
		}
		fmt.Println()
	}()

	fmt.Println("panicDemo 函数开始执行...")

	// 这个函数会触发 panic
	triggerPanic()

	// 有derer func recover这行也不会执行，因为panic 之后的代码不会执行
	fmt.Println("这行不会执行, 因为上面panic了")
}

// 出发 panic 的函数
func triggerPanic() {
	fmt.Println("准备触发panic...")

	// 手动触发panic
	panic("这是一个panic消息")

	// panic 之后的代码不会执行，只是为了后续函数正常执行所以在paincDemo加上recover()获取异常
	fmt.Println("这行不会执行", "\n")
}

// 注意：如果直接运行上面的代码，会导致整个程序崩溃,所以我加了recover捕获异常
// 下面的recover会演示如何捕获panic
func recoverDemo() {
	fmt.Println("recoverDemo函数开始...")
	// 运行被保护的函数
	result := safeFunction()

	fmt.Printf("safeFunction 返回值: %d\n", result)
	fmt.Println("程序继续正常运行，没有崩溃", "\n")
}

// 被recover保护的函数
func safeFunction() int {
	// defer 必须在panic之前定义
	// defer 函数会在panic时（或其他退出时）执行
	defer func() {
		// recover() 返回panic的值
		if r := recover(); r != nil {
			fmt.Printf("recover 捕获到 panic: %v\n", r)
		}
	}()

	fmt.Println("正常执行...")
	time.Sleep(10 * time.Millisecond)

	// 触发 panic
	panic("模拟严重错误")

	// 这行不会执行
	fmt.Println("这行不会执行")
	return 0
}

// defer 执行顺序：先定义的后执行（栈的后进先出）
func completeDemo() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("捕获到panic 异常: %v\n\n", r)
		}
	}()

	fmt.Println("completeDemo函数开始")

	// 多个 defer
	defer fmt.Println("defer 1 (最后执行)")
	defer fmt.Println("defer 2 (其次执行)")

	fmt.Println("主函数开始")

	nestedPanic()

	fmt.Println("主函数结束")
}

// 嵌套函数中的panic
func nestedPanic() {
	defer fmt.Println("nestedPanic 的 defer")

	fmt.Println("nestedPanic 开始")
	panic("nestedPanic 中的错误")
	fmt.Println("nestedPanic 结束")
}

// 错误处理的实践
// 1. 使用errors.Is()判断错误
func checkErrorType() {
	err := fmt.Errorf("原始错误: %w", os.ErrNotExist)

	// 检查是否是某个特定错误
	if errors.Is(err, os.ErrNotExist) {
		fmt.Println("文件不存在")
	}

	// 检查是否是自定义错误类型
	if errors.Is(err, ErrInvalidAge) {
		fmt.Println("年龄错误")
	}
	fmt.Println()
}

// 2. 只对真正异常的情况使用panic
func validateConfig(config map[string]string) {
	if config == nil {
		// 配置错误，这是真正的异常，应该panic
		panic("config 不能为nil")
	}

	// 普通验证错误，返回error，检验是否存在
	if _, ok := config["port"]; !ok {  // ok 是 bool 类型
		fmt.Println("警告: port 未配置，使用默认值")
	}

	// 检验是否为空
	portVal, ok := config["port"]
	// 键不存在 或 键存在但值为空 → 触发警告
	if !ok || portVal == "" {
		fmt.Println("警告: port 不能为空或未配置")
	}
}

// 注意：os.ErrNotExist 需要导入 "os" 包
// 为了代码完整性，这里模拟一个
var os = struct {
	ErrNotExist error
}{
	ErrNotExist: errors.New("file does not exist"),
}
