package mybase

import "fmt"

func Lesson3() {

	// 整型
	var age int = 22
	var savings uint = 100000
	fmt.Printf("年龄：%d\n", age)
	fmt.Printf("存款：%d\n", savings)

	// 浮点型
	var height float64 = 1.75
	var weight float64 = 65.5
	bmi := weight / (height * height)
	fmt.Printf("BMI：%.2f\n", bmi)

	// 布尔型
	var isStudent bool = true
	var hasJob bool = false
	fmt.Printf("是学生：%v\n", isStudent)
	fmt.Printf("有工作：%v\n", hasJob)
	fmt.Printf("是学生且是男性：%v\n", isStudent && true)
	fmt.Printf("有工作或是学生：%v\n", hasJob || isStudent)

	// 字符串
	var myName string = "Austoin"
	var myCity string = "重庆"
	message := myName + " 来自 " + myCity
	fmt.Println(message)

	// 字符串切片
	fmt.Printf("名字的第一个字符：%c\n", myName[0])

	// 零值验证
	var zeroInt int
	var zeroFloat float64
	var zeroBool bool
	var zeroString string
	fmt.Printf("int零值：%d\n", zeroInt)
	fmt.Printf("float零值：%.2f\n", zeroFloat)
	fmt.Printf("bool零值：%v\n", zeroBool)
	fmt.Printf("string零值：'%s'\n", zeroString)
}
