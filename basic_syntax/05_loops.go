package main

// 本文件演示 Go 语言的循环语句

import "fmt"

func main() {
	// ========== for 循环 ==========

	// 1. 基本 for 循环（类似 C 语言）
	// for 初始化; 条件; 增量 { }
	for i := 0; i < 5; i++ {
		fmt.Printf("i = %d\n", i)
	}

	// 2. 省略初始化（while 循环效果）
	// Go 没有 while 关键字，但可以通过省略初始化和增量实现
	counter := 0
	for counter < 3 {
		fmt.Printf("counter = %d\n", counter)
		counter++
	}

	// 3. 无限循环
	// 省略所有三个部分，创建无限循环
	// 通常配合 break 或 return 使用
	// for {
	//     fmt.Println("无限循环")
	//     break  // 必须有退出条件
	// }

	// 4. for range 遍历（最常用的遍历方式）
	// 用于遍历数组、切片、映射、字符串、通道等

	// 遍历切片
	nums := []int{1, 2, 3, 4, 5}
	for index, value := range nums {
		fmt.Printf("索引: %d, 值: %d\n", index, value)
	}

	// 省略索引（使用空白标识符 _）
	sum := 0
	for _, num := range nums {
		sum += num
	}
	fmt.Printf("切片元素之和: %d\n", sum)

	// 5. 遍历字符串
	// 字符串是字节数组，range 会按字节遍历
	str := "Go语言"
	for i, ch := range str {
		fmt.Printf("索引: %d, 字符: %c, 字节: %d\n", i, ch, str[i])
	}

	// 正确遍历 Unicode 字符
	for i, ch := range str {
		fmt.Printf("字符: %c, Unicode 码点: U+%04X\n", ch, ch)
	}

	// 6. 遍历映射（无序）
	ages := map[string]int{
		"张三": 25,
		"李四": 30,
		"王五": 28,
	}
	for name, age := range ages {
		fmt.Printf("%s 的年龄是 %d\n", name, age)
	}

	// 7. 遍历通道
	// ch := make(chan int, 3)
	// ch <- 1
	// ch <- 2
	// ch <- 3
	// close(ch)
	// for num := range ch {
	//     fmt.Println(num)
	// }

	// ========== break 和 continue ==========

	// break 退出循环
	fmt.Println("\nbreak 示例:")
	for i := 1; i <= 10; i++ {
		if i == 5 {
			break // 当 i 等于 5 时退出循环
		}
		fmt.Printf("%d ", i)
	}
	fmt.Println("\n循环结束")

	// continue 跳过本次循环，继续下一次
	fmt.Println("\ncontinue 示例:")
	for i := 1; i <= 5; i++ {
		if i == 3 {
			continue // 当 i 等于 3 时跳过本次循环
		}
		fmt.Printf("%d ", i)
	}
	fmt.Println("\n循环结束")

	// ========== 嵌套循环 ==========

	fmt.Println("\n嵌套循环示例:")
	for i := 1; i <= 3; i++ {
		for j := 1; j <= 3; j++ {
			fmt.Printf("(%d,%d) ", i, j)
		}
		fmt.Println()
	}

	// ========== 标签和 goto ==========

	// break + 标签（退出多层循环）
	fmt.Println("\nbreak + 标签示例:")
outer:
	for i := 1; i <= 3; i++ {
		for j := 1; j <= 3; j++ {
			if i == 2 && j == 2 {
				break outer // 直接退出外层循环
			}
			fmt.Printf("(%d,%d) ", i, j)
		}
		fmt.Println()
	}
	fmt.Println("已跳出嵌套循环")

	// continue + 标签
	fmt.Println("\ncontinue + 标签示例:")
outer2:
	for i := 1; i <= 3; i++ {
		for j := 1; j <= 3; j++ {
			if j == 2 {
				continue outer2 // 跳到外层循环的下一个迭代
			}
			fmt.Printf("(%d,%d) ", i, j)
		}
		fmt.Println()
	}

	// goto（不推荐使用，但需要了解）
	fmt.Println("\ngoto 示例:")
	goto label1
	fmt.Println("这行不会执行") // 跳过
label1:
	fmt.Println("跳转到这里")

	// ========== 循环控制技巧 ==========

	// 1. 反向遍历
	for i := len(nums) - 1; i >= 0; i-- {
		fmt.Printf("反向: %d ", nums[i])
	}
	fmt.Println()

	// 2. 带步长的遍历
	for i := 0; i < 10; i += 2 {
		fmt.Printf("步长2: %d ", i)
	}
	fmt.Println()

	// 3. 找到第一个满足条件的元素
	target := 3
	found := false
	for _, num := range nums {
		if num == target {
			found = true
			break
		}
	}
	fmt.Printf("找到 %d: %v\n", target, found)

	// 4. 统计满足条件的元素数量
	count := 0
	for _, num := range nums {
		if num > 2 {
			count++
		}
	}
	fmt.Printf("大于2的元素数量: %d\n", count)

	// 5. 构建新切片
	evenNums := []int{}
	for _, num := range nums {
		if num%2 == 0 {
			evenNums = append(evenNums, num)
		}
	}
	fmt.Printf("偶数切片: %v\n", evenNums)

	// ========== 常见错误 ==========

	// 1. 循环变量重用（Go 1.22+ 修复了这个问题）
	// 旧版本：for i := range nums { go func() { fmt.Println(i) }() }
	// 新版本：每个迭代都有独立的变量

	// 2. 忘记初始化
	// for i < 5 {  // i 未初始化，编译错误
	//     fmt.Println(i)
	// }

	// 3. 浮点数比较
	// for i := 0.1; i < 0.9; i += 0.1 {  // 浮点数精度问题
	//     fmt.Println(i)
	// }
}

// ========== 总结 ==========
// 1. for 循环是 Go 中唯一的循环关键字
// 2. for range 用于遍历各种数据结构
// 3. break 退出循环，continue 跳过本次迭代
// 4. break/continue 可配合标签使用
// 5. goto 可以跳转到标签处（谨慎使用）
// 6. 注意浮点数循环的精度问题
