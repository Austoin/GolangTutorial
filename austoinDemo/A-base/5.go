package mybase

import (
	"fmt"
)

func Lesson5() {
	for i := 1; i < 5; i++ {
		fmt.Printf("%d ", i)
		if i == 4 {
			fmt.Println()
		}
	}

	// while 循环
	cnt := 0
	for cnt < 3 {
		fmt.Printf("计数器：%d\n", cnt)
		cnt++
	}

	// for range 遍历切片
	nums := []int{10, 20, 30, 40}
	for index, value := range nums {
		fmt.Printf("索引：%d, 值：%d\n", index, value)
	}

	// 计算切片和
	sum := 0
	for _, value := range nums {
		sum += value
	}
	fmt.Printf("切片和：%d\n", sum)

	// 遍历字符串
	str := "Hello, GO!"
	for i, ch := range str {
		fmt.Printf("索引：%d, 字符：%c, Unicode: %U+%04X\n", i, ch, ch)
	}

	// 遍历映射
	scores := map[string]int{
		"Alice":   85,
		"Bob":     90,
		"Charlie": 78,
	}
	for name, score := range scores {
		fmt.Printf("%s 的成绩是 %d\n", name, score)
	}

	// break 和 continue
	for i := 1; i <= 10; i++ {
		if i == 3 {
			continue // 跳过当前循环
		}
		if i == 9 {
			break // 退出循环
		}
	}

	// 嵌套循环
	for i := 1; i <= 3; i++ {
		for j := 1; j <= 3; j++ {
			fmt.Printf("(%d, %d) ", i, j)
		}
	}
	fmt.Printf("\n------\n")

	// 九九乘法表
	for i := 1; i <= 9; i++ {
		for j := 1; j <= i; j++ {
			fmt.Printf("%d*%d=%d ", j, i, i*j)
			if j == i {
				fmt.Println()
			}
		}
	}

	// 构建新切片
	eventNums := []int{}
	for _, value := range nums {
		if value > 15 {
			eventNums = append(eventNums, value)
		}
	}
	fmt.Printf("\n大于15的数:%v\n\n", eventNums)

	// 反向遍历
	for i := len(nums) - 1; i >= 0; i-- {
		fmt.Printf("%d ", nums[i])
	}
	fmt.Println()

	// 带步长
	for i := 0; i <= 10; i += 2 {
		fmt.Printf("%d ", i)
	}
	fmt.Println()
}
