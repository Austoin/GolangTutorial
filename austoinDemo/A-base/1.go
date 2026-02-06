package mybase

import "fmt"

func Lesson1() {
	// 使用 fmt.Println 打印
	fmt.Println("Hello, World!")

	// 使用 fmt.Printf 格式化输出
	name := "Go"
	version := 1.21
	fmt.Printf("这是一个 %s 程序，版本 %.1f\n", name, version)

	// 练习：打印你的信息
	fmt.Println("我的名字是：Austoin")
	fmt.Println("我正在学习 Go 语言")
}
