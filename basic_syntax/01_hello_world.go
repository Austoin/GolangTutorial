// package main - 定义包名
// Go 程序必须有一个 main 包，这是程序的入口包
// 所有可执行的 Go 程序都必须包含一个 main 包
package main

// import "fmt" - 导入格式化包
// fmt 包提供了格式化 I/O 的功能
// 类似于 C 语言中的 printf，Python 中的 print
import (
	"fmt"
	"os"
)

// import "os" - 导入操作系统包
// 用于访问命令行参数和环境变量

// main 函数 - 程序的入口点
// 每个可执行的 Go 程序必须有一个 main 函数
// 程序执行从 main 函数开始
func main() {
	// ========== 基础输出 ==========

	// fmt.Println - 打印内容并自动换行
	// 自动在末尾添加换行符 \n
	// 参数可以是任意类型，函数会自动调用其 String() 方法或使用默认格式
	fmt.Println("Hello, World!") // 输出: Hello, World!
	fmt.Println("欢迎学习 Go 语言！")   // 输出中文支持
	fmt.Println(123)             // 输出数字
	fmt.Println(3.14)            // 输出浮点数
	fmt.Println(true)            // 输出布尔值

	// ========== 格式化输出 ==========

	// fmt.Printf - 格式化输出，不自动换行
	// %s - 字符串占位符
	// %d - 整数占位符
	// %f - 浮点数占位符
	// %v - 通用占位符（自动选择合适格式）
	name := "Go"    // := 声明并初始化变量
	version := 1.21 // 自动推导类型
	fmt.Printf("这是一个 %s 程序，版本 %.1f\n", name, version)

	// ========== 变量声明与输出 ==========

	// 方式1：var 声明
	var num1 int = 10
	// 方式2：简短声明（只能在函数内使用）
	num2 := 20
	// 方式3：多变量声明
	a, b := 5, 10

	fmt.Printf("num1 = %d, num2 = %d\n", num1, num2)
	fmt.Printf("a = %d, b = %d\n", a, b)

	// ========== 多行输出 ==========

	// 多次调用 Println 会自动换行
	fmt.Println("第一行")
	fmt.Println("第二行")
	fmt.Println("第三行")

	// ========== Print 和 Println 的区别 ==========

	// Print - 不自动添加换行符
	fmt.Print("这是 Print 输出，") // 不会换行
	fmt.Print("不会自动换行\n")     // 需要手动添加 \n

	// Println - 自动添加换行符，参数之间自动添加空格
	fmt.Println("这是 Println 输出，", "会自动添加空格和换行")

	// ========== Sprint 格式化 ==========

	// Sprint 系列函数返回字符串，不直接打印
	// Sprintf - 格式化后返回字符串
	message := fmt.Sprintf("Sprint 格式化：%s %d", "数字", 42)
	fmt.Println(message)

	// Sprintln - 自动添加空格和换行
	message2 := fmt.Sprintln("Sprintln", "会自动", "添加空格")
	fmt.Print(message2)

	// ========== 获取程序参数 ==========

	// os.Args 是一个字符串切片
	// os.Args[0] 是程序本身的路径
	// os.Args[1:] 是命令行参数
	fmt.Println("程序参数:", os.Args)

	// 遍历参数
	for i, arg := range os.Args {
		if i == 0 {
			continue // 跳过程序名
		}
		fmt.Printf("参数 %d: %s\n", i, arg)
	}

	// ========== 退出程序 ==========

	// os.Exit(0) - 立即退出程序
	// 参数 0 表示正常退出，非 0 表示异常退出
	// os.Exit(0)
}

// ========== 总结 ==========
// 1. package main - 程序入口包
// 2. func main() - 程序入口函数
// 3. fmt.Println() - 打印并换行
// 4. fmt.Printf() - 格式化打印
// 5. 变量声明：var 或 :=
// 6. os.Args - 获取命令行参数
