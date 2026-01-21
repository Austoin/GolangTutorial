package main

// 本文件演示 Go 语言的条件语句：if 和 switch

import (
	"fmt"
	"math"
)

func main() {
	// ========== if 语句 ==========

	// 1. 基础 if 语句
	age := 18
	if age >= 18 {
		fmt.Println("成年人") // 条件为 true 时执行
	}

	// 2. if-else 语句
	temperature := 25
	if temperature > 30 {
		fmt.Println("很热")
	} else {
		fmt.Println("不热") // 条件为 false 时执行
	}

	// 3. if-else if-else 多条件
	score := 85
	if score >= 90 {
		fmt.Println("优秀")
	} else if score >= 80 {
		fmt.Println("良好") // 这个分支会被执行
	} else if score >= 60 {
		fmt.Println("及格")
	} else {
		fmt.Println("不及格")
	}

	// 4. if 语句初始化（推荐写法）
	// 可以在 if 语句中声明变量，作用域仅限于 if-else 块
	if num := 10; num%2 == 0 {
		fmt.Printf("%d 是偶数\n", num)
	} else {
		fmt.Printf("%d 是奇数\n", num)
	}
	// 注意：num 在这里无法访问，作用域限制在 if-else 块内

	// 5. 多个条件（使用 && 和 ||）
	x := 5
	if x > 0 && x < 10 {
		fmt.Printf("%d 在 0 到 10 之间\n", x)
	}

	y := 15
	if y < 0 || y > 10 {
		fmt.Printf("%d 不在 0 到 10 之间\n", y)
	}

	// 6. 否定条件
	isActive := true
	if !isActive {
		fmt.Println("未激活")
	} else {
		fmt.Println("已激活") // 这个分支会被执行
	}

	// 7. 检查错误（Go 常见模式）
	// 通常函数返回 (result, error)，检查 error 是否为 nil
	result, err := sqrt(16)
	if err != nil {
		fmt.Printf("错误: %v\n", err)
	} else {
		fmt.Printf("平方根: %.2f\n", result)
	}

	// 8. 嵌套 if
	a, b, c := 3, 5, 7
	if a < b {
		if a < c {
			fmt.Printf("%d 是最小值\n", a)
		}
	}

	// ========== switch 语句 ==========

	// 1. 基础 switch
	day := 3
	switch day {
	case 1:
		fmt.Println("星期一")
	case 2:
		fmt.Println("星期二")
	case 3:
		fmt.Println("星期三") // 这个分支会被执行
	case 4:
		fmt.Println("星期四")
	default:
		fmt.Println("其他") // 没有匹配时执行
	}

	// 2. 多个 case 合并
	month := 12
	switch month {
	case 1, 2, 3:
		fmt.Println("第一季度")
	case 4, 5, 6:
		fmt.Println("第二季度")
	case 7, 8, 9:
		fmt.Println("第三季度")
	case 10, 11, 12:
		fmt.Println("第四季度") // 这个分支会被执行
	}

	// 3. switch 不带表达式（相当于多个 if-else）
	grade := 'B'
	switch {
	case grade == 'A':
		fmt.Println("优秀")
	case grade == 'B':
		fmt.Println("良好") // 这个分支会被执行
	case grade == 'C':
		fmt.Println("一般")
	default:
		fmt.Println("未知")
	}

	// 4. fallthrough 穿透
	// fallthrough 会强制执行下一个 case（无论条件是否满足）
	num := 2
	switch num {
	case 1:
		fmt.Println("case 1")
		fallthrough
	case 2:
		fmt.Println("case 2 (fallthrough)") // 这个会执行
		fallthrough
	case 3:
		fmt.Println("case 3 (fallthrough)") // 这个也会执行
	default:
		fmt.Println("default")
	}

	// 5. Type Switch（类型断言）
	var interfaceVar interface{} = "hello"
	switch v := interfaceVar.(type) {
	case string:
		fmt.Printf("字符串: %s\n", v)
	case int:
		fmt.Printf("整数: %d\n", v)
	case float64:
		fmt.Printf("浮点数: %.2f\n", v)
	default:
		fmt.Printf("未知类型: %T\n", v)
	}

	// 6. switch 初始化
	switch os := "linux"; os {
	case "windows":
		fmt.Println("Windows 系统")
	case "linux":
		fmt.Println("Linux 系统") // 这个分支会被执行
	case "macos":
		fmt.Println("macOS 系统")
	}

	// 7. 表达式作为 case
	value := 15
	switch {
	case value > 0 && value <= 10:
		fmt.Println("1-10")
	case value > 10 && value <= 20:
		fmt.Println("11-20") // 这个分支会被执行
	case value > 20:
		fmt.Println("20+")
	}

	// ========== if vs switch 选择 ==========
	// - if: 适用于复杂条件判断
	// - switch: 适用于单一变量的多个值判断

	// 推荐写法：使用 switch 替代多个 if-else
	status := 404
	switch status {
	case 200:
		fmt.Println("成功")
	case 400:
		fmt.Println("请求错误")
	case 404:
		fmt.Println("资源不存在") // 这个分支会被执行
	case 500:
		fmt.Println("服务器错误")
	default:
		fmt.Println("未知状态")
	}

	// ========== 比较运算符 ==========
	// == 相等, != 不相等, > 大于, < 小于, >= 大于等于, <= 小于等于

	a1, b1 := 5, 5
	fmt.Printf("\n比较运算:\n")
	fmt.Printf("  %d == %d: %v\n", a1, b1, a1 == b1)
	fmt.Printf("  %d != %d: %v\n", a1, b1, a1 != b1)
	fmt.Printf("  %d > %d: %v\n", a1, b1, a1 > b1)
	fmt.Printf("  %d < %d: %v\n", a1, b1, a1 < b1)
	fmt.Printf("  %d >= %d: %v\n", a1, b1, a1 >= b1)
	fmt.Printf("  %d <= %d: %v\n", a1, b1, a1 <= b1)
}

// sqrt 计算平方根，如果参数为负则返回错误
func sqrt(n float64) (float64, error) {
	if n < 0 {
		return 0, fmt.Errorf("负数没有实数平方根")
	}
	return math.Sqrt(n), nil
}

// ========== 总结 ==========
// 1. if 语句：if、if-else、if-else if-else
// 2. if 可以包含初始化语句
// 3. switch 语句：case 匹配、fallthrough 穿透
// 4. switch 可以不带表达式（多个 if-else）
// 5. Type Switch 用于判断接口值的实际类型
// 6. 推荐使用 switch 替代多个 if-else
