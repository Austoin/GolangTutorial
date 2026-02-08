package mybase 

import (
	"fmt"
	"time"
)

func Lesson11() {
	packageDemo()
	importDemo()
	moduleDemo()
}

func packageDemo() {
	// 包内调用：同一包中的函数可以直接调用
	fmt.Println("调用同包的 sayHello 函数")
	sayHello()
	fmt.Println()
}

// sayHello 包内函数
// 首字母小写 = 私有，只能在本包内使用(在别的包声明过了，就注释了，不然报错)
// func sayHello() {
// 	fmt.Println("  Hello from same package!")
// }

// SayHi 公开函数
// 首字母大写 = 公开，可以被其他包调用
func SayHi() {
	fmt.Println("  Hi from public function!")
}

func importDemo() {
	// 标准库导入
	fmt.Println("当前时间:", time.Now().Format("2006-01-02 15:04:05"))

	// 格式化输出
	name := "Austoin"
	age := 22
	fmt.Printf("姓名：%s, 年龄：%d\n", name, age)

	// fmt 常用函数
	// Print 不换行
	fmt.Print("Print不换行 ")
	// Println 换行
	fmt.Println("Println换行")
	// Printf 格式化
	fmt.Printf("Printf格式化: %s\n", "hello")

	fmt.Println()
}

func moduleDemo() {
	// 模块的作用：
	// 1. 标识项目（类似Java的包名）
	// 2. 管理依赖（哪些外部库）
	// 3. 版本控制（依赖的版本）
	// replace 的作用：
	// 开发时，代码在本地，没有上传GitHub
	// 用 replace 告诉Go："这个模块虽然在远程路径里，但实际上在本地"
	fmt.Println()
}

// 牢记
// 1. 包名规范
// - 简洁：有意义但不冗长
// - 与目录名一致（惯例）
// - 避免使用下划线或驼峰（用小写字母）
// 2. 导入规范
// - 只导入使用的包
// - 使用标准导入分组（标准库、第三方、本地）
// import (
//     "fmt"                    // 标准库
//     "time"
//     "github.com/gin-gonic/gin"  // 第三方
//     abaselib "本地路径"      // 本地（带别名）
// )
// 3. 可见性规范
// - 暴露的API：首字母大写
// - 内部实现：首字母小写
// - 封装原则：只暴露必要的