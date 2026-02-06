package mybase 

import (
	"fmt"
)

func Lesson7() {
	// 数组
	scores := [5]int{90, 85, 78, 92, 88}
	fmt.Println("数组：", scores, "长度：", len(scores))

	// for 循环遍历
	for i := 0; i < len(scores); i++ {
		fmt.Printf("第%d个成绩: %d\n", i+1, scores[i])
	}

	// for range 循环遍历
	for i, value := range scores {
		fmt.Printf("索引：%d, 值：%d\n", i, value)
	}

	// 切片
	// 方式1：直接声明
	nums := []int{1, 2, 3, 4, 5}
	fmt.Printf("切片：%v, 长度：%d, 容量：%d\n", nums, len(nums), cap(nums))

	// 方式2：使用 make 创建切片
	arr := make([]int, 3, 5) // 长度3，容量5
	fmt.Printf("make 切片:%v, len=%d, cap=%d\n", arr, len(arr), cap(arr))

	// 切取子切片
	sub := nums[1:4] // 包含索引1到3的元素
	fmt.Printf("子切片：%v\n", sub)
	fmt.Printf("%v\n", sub[:])

	nums = append(nums, 6, 7)  // 追加元素
	fmt.Printf("nums 追加元素后：%v\n", nums)

	// 合并两个切片
	extra := []int{8, 9, 10}
	nums = append(nums, extra...)
	fmt.Printf("合并切片后：%v\n", nums)

	// 映射（Map）
	age := map[string]int{
		"Alice": 30,
		"Bob":   25,
		"Charlie": 35,
	}
	fmt.Printf("映射：%v\n", age)

	// 添加新键值对
	age["David"] = 28
	fmt.Printf("添加新键值对后：%v\n", age)

	// 获取元素
	name := "Alice"
	if age, exists := age[name]; exists { // 判断键是否存在
		fmt.Printf("%s 的年龄是 %d\n", name, age)
	} else {
		fmt.Printf("%s 不存在\n", name)
	}
	// 获取不存在的键返回零值
	fmt.Printf("Eve 的年龄是 %d\n", age["Eve"]) // 返回0，因为Eve不存在

	// 删除键值对
	delete(age, "Bob")
	fmt.Printf("删除键值对后：%v\n", age)

	// 遍历映射
	for key, value := range age {
		fmt.Printf("age[%s] = %d\n", key, value)
	}

	// 综合应用 - 统计分数段
	allScores := []int{95, 82, 78, 88, 92, 65, 72, 85, 91, 68}
	groups := make(map[string]int)
	for _, score := range allScores {
		switch { // 等同 switch true
		case score >= 90:
			groups["优秀"]++
		case score >= 80:
			groups["良好"]++
		case score >= 70:
			groups["中等"]++
		default:
			groups["及格"]++
		}
	}
	fmt.Printf("分数段统计：%v\n", groups)

	// 过滤偶数
	allNums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	events := []int{}
	for _, value := range allNums {
		if value%2 == 0 {
			events = append(events, value)
		}
	}
	fmt.Printf("过滤偶数后的结果：%v\n", events)
}