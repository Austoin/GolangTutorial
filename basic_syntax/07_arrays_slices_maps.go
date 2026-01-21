package main

// 本文件演示 Go 语言的数组、切片和映射

import "fmt"

func main() {
	// ========== 数组 ==========

	// 1. 数组声明和初始化
	// 数组长度固定，类型确定
	// var 数组名 [长度]类型

	// 方式1：声明并初始化
	var arr1 [5]int = [5]int{1, 2, 3, 4, 5}
	fmt.Printf("arr1: %v\n", arr1)

	// 方式2：简短声明
	arr2 := [3]string{"Go", "Python", "Java"}
	fmt.Printf("arr2: %v\n", arr2)

	// 方式3：部分初始化（其余为类型零值）
	arr3 := [5]int{1, 2} // [1, 2, 0, 0, 0]
	fmt.Printf("arr3: %v\n", arr3)

	// 方式4：指定位置初始化
	arr4 := [5]int{0: 10, 2: 30} // [10, 0, 30, 0, 0]
	fmt.Printf("arr4: %v\n", arr4)

	// 方式5：使用 ... 让编译器推断长度
	arr5 := [...]int{1, 2, 3, 4, 5} // 长度自动推断为 5
	fmt.Printf("arr5: %v (长度: %d)\n", arr5, len(arr5))

	// 2. 访问数组元素
	fmt.Printf("arr1[0] = %d\n", arr1[0])
	fmt.Printf("arr1[4] = %d\n", arr1[4])

	// 3. 修改数组元素
	arr1[0] = 100
	fmt.Printf("修改后 arr1: %v\n", arr1)

	// 4. 数组遍历
	fmt.Println("\n数组遍历:")
	for i := 0; i < len(arr1); i++ {
		fmt.Printf("arr1[%d] = %d\n", i, arr1[i])
	}

	// 使用 for range 遍历
	fmt.Println("\nfor range 遍历:")
	for index, value := range arr1 {
		fmt.Printf("索引: %d, 值: %d\n", index, value)
	}

	// 5. 多维数组
	var matrix [3][3]int = [3][3]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}
	fmt.Printf("matrix: %v\n", matrix)

	// 遍历多维数组
	fmt.Println("\n多维数组遍历:")
	for i := 0; i < len(matrix); i++ {
		for j := 0; j < len(matrix[i]); j++ {
			fmt.Printf("matrix[%d][%d] = %d ", i, j, matrix[i][j])
		}
		fmt.Println()
	}

	// 6. 数组是值类型
	// 数组赋值会复制整个数组
	original := [3]int{1, 2, 3}
	copy := original                       // 复制数组
	copy[0] = 100                          // 修改副本
	fmt.Printf("original: %v\n", original) // [1, 2, 3]
	fmt.Printf("copy: %v\n", copy)         // [100, 2, 3]

	// ========== 切片 ==========

	// 7. 切片声明
	// 切片是数组的视图，长度可变
	// var 切片名 []类型

	var slice1 []int // nil 切片，零值
	fmt.Printf("slice1: %v, 长度: %d, 容量: %d\n", slice1, len(slice1), cap(slice1))

	// 8. 使用 make 创建切片
	// make([]类型, 长度, 容量)
	slice2 := make([]int, 5)    // 长度5，容量5
	slice3 := make([]int, 3, 5) // 长度3，容量5
	fmt.Printf("slice2: %v, len=%d, cap=%d\n", slice2, len(slice2), cap(slice2))
	fmt.Printf("slice3: %v, len=%d, cap=%d\n", slice3, len(slice3), cap(slice3))

	// 9. 使用字面量创建切片
	slice4 := []int{1, 2, 3, 4, 5}
	fmt.Printf("slice4: %v, len=%d, cap=%d\n", slice4, len(slice4), cap(slice4))

	// 10. 切片是引用类型
	// 切片引用同一个底层数组
	s1 := []int{1, 2, 3, 4, 5}
	s2 := s1[1:3] // 切片 s2 指向 s1 的子数组
	fmt.Printf("s1: %v, s2: %v\n", s1, s2)

	s2[0] = 100 // 修改 s2 会影响 s1
	fmt.Printf("修改后 s1: %v, s2: %v\n", s1, s2)

	// 11. 切片操作
	nums := []int{1, 2, 3, 4, 5}

	// 切取子切片 [start:end]
	// 包含 start，不包含 end
	sub1 := nums[1:4] // [2, 3, 4]
	sub2 := nums[:3]  // [1, 2, 3] 从头开始
	sub3 := nums[2:]  // [3, 4, 5] 到末尾
	sub4 := nums[:]   // [1, 2, 3, 4, 5] 整个切片

	fmt.Printf("sub1: %v\n", sub1)
	fmt.Printf("sub2: %v\n", sub2)
	fmt.Printf("sub3: %v\n", sub3)
	fmt.Printf("sub4: %v\n", sub4)

	// 12. append 追加元素
	// 切片会自动扩容
	nums = append(nums, 6, 7, 8)
	fmt.Printf("append 后: %v\n", nums)

	// 使用 append 合并切片
	nums2 := []int{9, 10}
	nums = append(nums, nums2...)
	fmt.Printf("合并后: %v\n", nums)

	// 13. copy 复制切片
	src := []int{1, 2, 3}
	dst := make([]int, len(src))
	n := copy(dst, src) // 返回复制的元素数量
	fmt.Printf("copy: src=%v, dst=%v, n=%d\n", src, dst, n)

	// 14. 删除切片元素
	nums = append(nums[:2], nums[3:]...) // 删除索引2的元素
	fmt.Printf("删除后: %v\n", nums)

	// 15. 切片扩容
	// 当容量不足时，Go 会自动扩容（通常翻倍）
	fmt.Println("\n切片扩容示例:")
	s := make([]int, 0, 2) // 初始容量2
	fmt.Printf("初始: len=%d, cap=%d\n", len(s), cap(s))

	s = append(s, 1, 2)
	fmt.Printf("添加2个: len=%d, cap=%d\n", len(s), cap(s))

	s = append(s, 3)
	fmt.Printf("再添加1个: len=%d, cap=%d\n", len(s), cap(s))

	// ========== 映射 ==========

	// 16. 映射声明
	// 映射是键值对的无序集合
	// var 映射名 map[键类型]值类型

	var m1 map[string]int // nil 映射
	fmt.Printf("m1: %v\n", m1)

	// 17. 使用 make 创建映射
	m2 := make(map[string]int)
	fmt.Printf("m2: %v\n", m2)

	// 18. 使用字面量创建映射
	m3 := map[string]int{
		"张三": 25,
		"李四": 30,
		"王五": 28,
	}
	fmt.Printf("m3: %v\n", m3)

	// 19. 添加/修改元素
	m2["Go"] = 90
	m2["Python"] = 85
	m2["Java"] = 80
	fmt.Printf("m2: %v\n", m2)

	// 20. 获取元素
	// 映射返回两个值：值和是否存在
	score, exists := m2["Python"]
	if exists {
		fmt.Printf("Python 分数: %d\n", score)
	} else {
		fmt.Println("Python 不存在")
	}

	// 获取不存在的键返回零值
	fmt.Printf("C++ 分数: %d (零值)\n", m2["C++"])

	// 21. 删除元素
	delete(m2, "Java")
	fmt.Printf("删除后 m2: %v\n", m2)

	// 22. 遍历映射
	fmt.Println("\n映射遍历:")
	for key, value := range m3 {
		fmt.Printf("%s: %d岁\n", key, value)
	}

	// 只遍历键
	for key := range m3 {
		fmt.Printf("键: %s\n", key)
	}

	// 只遍历值
	for _, value := range m3 {
		fmt.Printf("值: %d\n", value)
	}

	// 23. 映射是引用类型
	originalMap := map[string]int{"A": 1}
	refMap := originalMap
	refMap["B"] = 2
	fmt.Printf("originalMap: %v\n", originalMap) // {A:1, B:2}
	fmt.Printf("refMap: %v\n", refMap)           // {A:1, B:2}

	// 24. 检查键是否存在（comma-ok 模式）
	if score, ok := m3["张三"]; ok {
		fmt.Printf("找到张三: %d岁\n", score)
	} else {
		fmt.Println("未找到张三")
	}

	// ========== 实际应用示例 ==========

	// 25. 统计字符出现频率
	text := "hello world"
	freq := make(map[rune]int)
	for _, ch := range text {
		freq[ch]++
	}
	fmt.Printf("\n字符频率: %v\n", freq)

	// 26. 过滤切片
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	evens := filter(numbers, func(n int) bool {
		return n%2 == 0
	})
	fmt.Printf("偶数: %v\n", evens)
}

// filter 过滤函数
func filter(numbers []int, f func(int) bool) []int {
	result := []int{}
	for _, n := range numbers {
		if f(n) {
			result = append(result, n)
		}
	}
	return result
}

// ========== 总结 ==========
// 1. 数组：固定长度，值类型
// 2. 切片：可变长度，引用类型，底层是数组
// 3. 映射：键值对，无序，引用类型
// 4. append 自动扩容切片
// 5. copy 复制切片元素
// 6. range 遍历切片和映射
// 7. 切片和映射都是引用类型
