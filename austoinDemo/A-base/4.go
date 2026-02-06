package mybase

import "fmt"

func Lesson4() {

	// 练习1：if-else 成绩评级
	score := 85
	if score >= 90 {
		fmt.Println("优秀")
	} else if score >= 80 {
		fmt.Println("良好")
	} else if score >= 70 {
		fmt.Println("中等")
	} else if score >= 60 {
		fmt.Println("及格")
	} else {
		fmt.Println("不及格")
	}

	// 练习2：switch 星期判断
	day := 3
	switch day {
	case 1:
		fmt.Println("星期一")
	case 2:
		fmt.Println("星期二")
	case 3:
		fmt.Println("星期三")
	case 4:
		fmt.Println("星期四")
	case 5:
		fmt.Println("星期五")
	case 6:
		fmt.Println("星期六")
	case 7:
		fmt.Println("星期日")
	default:
		fmt.Println("无效日期")
	}

	// 练习3：switch 季节判断（多个case合并）
	month := 6
	switch month {
	case 3, 4, 5:
		fmt.Println("春季")
	case 6, 7, 8:
		fmt.Println("夏季")
	case 9, 10, 11:
		fmt.Println("秋季")
	case 12, 1, 2:
		fmt.Println("冬季")
	default:
		fmt.Println("无效月份")
	}

	// 练习4：多个条件判断
	age := 25
	isWeekend := false
	if age >= 18 && age <= 60 && !isWeekend {
		fmt.Println("工作日，成年人")
	} else {
		fmt.Println("非工作日或非成年人")
	}

	// 练习5：switch 实现计算器
	a, b := 10, 5
	operator := "+"
	switch operator {
	case "+":
		fmt.Printf("%d + %d = %d\n", a, b, a+b)
	case "-":
		fmt.Printf("%d - %d = %d\n", a, b, a-b)
	case "*":
		fmt.Printf("%d * %d = %d\n", a, b, a*b)
	case "/":
		fmt.Printf("%d / %d = %.2f\n", a, b, float64(a)/float64(b))
	default:
		fmt.Println("未知运算符")
	}

	// 练习6：BMI 分类
	h := 1.75
	w := 70.0
	calcBMI := w / (h * h)
	if calcBMI < 18.5 {
		fmt.Println("偏瘦")
	} else if calcBMI < 24 {
		fmt.Println("正常")
	} else if calcBMI < 28 {
		fmt.Println("偏胖")
	} else {
		fmt.Println("肥胖")
	}
}
