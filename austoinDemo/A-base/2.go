package mybase

import "fmt"

func Lesson2() {

	// 方式1：var 声明
	var name string = "Austoin"
	var age int = 22

	// 方式2：简短声明 :=
	city := "重庆"
	height := 1.75

	// 常量声明
	const country = "中国"
	const pi = 3.14

	// 打印变量
	fmt.Println("姓名：", name)
	fmt.Println("年龄：", age)
	fmt.Println("城市：", city)
	fmt.Println("身高：", height)
	fmt.Println("国家：", country)
	fmt.Printf("圆周率：%.2f\n", pi)

	// 练习：修改值
	age = 23
	fmt.Println("明年年龄：", age)
}
