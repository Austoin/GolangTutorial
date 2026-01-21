package main

// 本文件演示 Go 语言中变量声明的各种方式

import "fmt"

// ========== 全局变量声明 ==========
// 全局变量（函数外部）只能使用 var 声明
// 不能使用 := 简短声明

// 方式1：声明单个变量
var globalVar1 int

// 方式2：声明并初始化
var globalVar2 string = "全局字符串"

// 方式3：类型推断（省略类型）
var globalVar3 = 3.14

// 方式4：批量声明
var (
	name      string  = "张三"
	age       int     = 25
	height    float64 = 175.5
	isStudent bool    = false
)

// main 函数
func main() {
	// ========== 局部变量声明 ==========

	// 方式1：var 声明单个变量
	// var 变量名 类型 = 初始值
	var num1 int = 10
	fmt.Printf("num1: %d, 类型: %T\n", num1, num1)

	// 方式2：var 声明不初始化（使用零值）
	// 每种类型都有零值：
	// int → 0, float → 0.0, string → "", bool → false
	var num2 int
	var str1 string
	var flag bool
	fmt.Printf("零值 - num2: %d, str1: '%s', flag: %v\n", num2, str1, flag)

	// 方式3：简短声明（推荐在函数内使用）
	// 编译器自动推断类型
	// 只能用于局部变量，不能用于全局变量
	name := "李四"
	age := 30
	height := 180.5
	isActive := true
	fmt.Printf("简短声明 - name: %s, age: %d, height: %.1f, isActive: %v\n",
		name, age, height, isActive)

	// 方式4：批量简短声明
	a, b, c := 1, 2.5, true
	fmt.Printf("批量声明 - a: %d, b: %.1f, c: %v\n", a, b, c)

	// 方式5：交换变量值
	x, y := 10, 20
	fmt.Printf("交换前 - x: %d, y: %d\n", x, y)
	x, y = y, x // 直接交换，无需临时变量
	fmt.Printf("交换后 - x: %d, y: %d\n", x, y)

	// ========== 变量作用域 ==========

	// 局部变量：声明在函数内部，作用域限于该函数
	functionScopedVar := "我是函数内变量"
	fmt.Println(functionScopedVar)

	// 块作用域：if、for 等代码块内
	if true {
		blockVar := "我是代码块内变量"
		fmt.Println(blockVar)
	}
	// blockVar 在这里不可访问

	// ========== 变量遮蔽 ==========
	// 内层声明的变量会遮蔽外层同名的变量
	shadowVar := "外层变量"
	fmt.Println("外层:", shadowVar)

	if true {
		shadowVar := "内层变量" // 这是一个新变量，遮蔽了外层的 shadowVar
		fmt.Println("内层:", shadowVar)
	}
	fmt.Println("外层:", shadowVar) // 仍然是外层的值

	// ========== 使用已声明的变量 ==========
	// Go 要求所有声明的变量都必须被使用
	// 否则会编译错误（防止无用的变量占用内存）
	usedVar := 100
	fmt.Println("使用的变量:", usedVar)

	// 如果确实需要声明但不立即使用，可以使用空白标识符 _
	_, unusedValue := 10, 20 // 10 被丢弃，只使用 20
	fmt.Println("只使用第二个值:", unusedValue)

	// ========== 常量与变量的区别 ==========

	// 常量：编译期确定，运行时不可修改
	// 使用 const 关键字声明
	const pi = 3.14159
	const maxUsers = 1000

	// 变量：运行时可以修改
	var counter int = 0
	counter = counter + 1 // 可以修改变量值
	fmt.Println("计数器:", counter)

	// 常量不能修改
	// pi = 3.14  // 这行会编译错误

	// ========== iota 常量生成器 ==========
	// iota 在常量声明块中从 0 开始递增
	const (
		monday    = iota // 0
		tuesday   = iota // 1
		wednesday = iota // 2
	)

	// 简写：如果类型相同可以省略
	const (
		red   = iota // 0
		green = iota // 1
		blue  = iota // 2
	)

	// 更简洁的写法
	const (
		cat  = iota // 0
		dog         // 1，继承上一个表达式的类型和值
		fish        // 2
	)

	// 结合位移操作
	const (
		bit0 = 1 << iota // 1 << 0 = 1
		bit1             // 1 << 1 = 2
		bit2             // 1 << 2 = 4
	)

	fmt.Printf("常量示例 - monday: %d, red: %d, cat: %d, bit2: %d\n",
		monday, red, cat, bit2)

	// ========== 全局变量使用 ==========
	fmt.Printf("全局变量 - name: %s, age: %d\n", name, age)

	// ========== 类型转换 ==========

	// Go 没有隐式类型转换，必须显式转换
	intToFloat := float64(num1)          // int → float64
	floatToInt := int(3.14)              // float64 → int（会截断小数）
	intToString := fmt.Sprintf("%d", 10) // int → string
	stringToInt := 42                    // string → int 需要 strconv 包

	fmt.Printf("类型转换 - intToFloat: %.2f, floatToInt: %d\n",
		intToFloat, floatToInt)
	fmt.Printf("intToString: %s\n", intToString)
}

// ========== 总结 ==========
// 1. 变量声明方式：var、:=、批量声明
// 2. 变量作用域：全局、局部、块作用域
// 3. 变量遮蔽：内层变量会遮蔽外层同名变量
// 4. 必须使用所有声明的变量
// 5. 常量使用 const，iota 生成常量序列
// 6. 必须显式进行类型转换
