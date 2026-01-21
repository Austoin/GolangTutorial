package main

// 本文件演示 Go 语言的函数

import "fmt"

func main() {
	// ========== 函数调用 ==========

	// 1. 基本函数调用
	result := add(10, 20)
	fmt.Printf("add(10, 20) = %d\n", result)

	// 2. 多返回值函数调用
	sum, avg := calc(10, 20, 30)
	fmt.Printf("和: %d, 平均值: %.2f\n", sum, avg)

	// 3. 命名返回值函数调用
	name := getFullName("张", "三")
	fmt.Println("全名:", name)

	// 4. 变参函数调用
	nums := []int{1, 2, 3, 4, 5}
	sum2 := sumAll(nums...)
	fmt.Printf("sumAll(1,2,3,4,5) = %d\n", sum2)

	// ========== 函数作为参数 ==========

	// 5. 函数作为参数
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	evenSum := filterAndSum(numbers, func(n int) bool {
		return n%2 == 0 // 只累加偶数
	})
	fmt.Printf("偶数之和: %d\n", evenSum)

	// 6. 使用命名函数作为参数
	oddSum := filterAndSum(numbers, isOdd)
	fmt.Printf("奇数之和: %d\n", oddSum)

	// ========== 闭包 ==========

	// 7. 闭包示例
	counter := newCounter()
	fmt.Printf("计数器: %d\n", counter())
	fmt.Printf("计数器: %d\n", counter())
	fmt.Printf("计数器: %d\n", counter())

	// ========== defer ==========

	// 8. defer 延迟执行
	fmt.Println("\ndefer 示例:")
	deferExample()
	fmt.Println("deferExample 返回后")

	// 9. 多个 defer（后进先出）
	fmt.Println("\n多个 defer:")
	defer fmt.Println("最后执行") // 最后入栈，最先执行
	defer fmt.Println("中间执行") // 中间入栈，中间执行
	defer fmt.Println("最先执行") // 最先入栈，最后执行
	fmt.Println("先执行")

	// ========== 错误处理 ==========

	// 10. 返回错误
	value, err := safeDiv(10, 0)
	if err != nil {
		fmt.Printf("错误: %v\n", err)
	} else {
		fmt.Printf("结果: %d\n", value)
	}

	// ========== 递归 ==========

	// 11. 递归函数
	factorialResult := factorial(5)
	fmt.Printf("5! = %d\n", factorialResult)

	fibResult := fib(10)
	fmt.Printf("fib(10) = %d\n", fibResult)
}

// ========== 函数定义 ==========

// 1. 基本函数
// func 函数名(参数列表) 返回类型 { }
// 参数格式：参数名 类型，多个参数用逗号分隔
func add(a int, b int) int {
	return a + b // 返回结果
}

// 2. 多返回值函数
// Go 函数可以返回多个值
func calc(a, b, c int) (int, float64) {
	sum := a + b + c
	avg := float64(sum) / 3.0
	return sum, avg
}

// 3. 命名返回值
// 可以给返回值命名，在函数体中可以直接使用
func getFullName(firstName, lastName string) (fullName string) {
	// fullName 已经声明，可以直接使用
	fullName = firstName + lastName
	return // 可以省略返回值（裸 return）
}

// 4. 可变参数函数
// 最后一个参数可以使用 ...类型 表示接受任意数量的参数
func sumAll(nums ...int) int {
	total := 0
	// nums 是一个切片，可以像切片一样使用
	for _, num := range nums {
		total += num
	}
	return total
}

// 5. 函数作为类型
// 定义函数类型
type filterFunc func(n int) bool

// 6. 高阶函数（接受函数作为参数）
func filterAndSum(numbers []int, filter filterFunc) int {
	sum := 0
	for _, num := range numbers {
		if filter(num) {
			sum += num
		}
	}
	return sum
}

// 7. 命名函数作为参数
func isOdd(n int) bool {
	return n%2 != 0
}

// 8. 闭包函数
// 闭包是一个函数及其引用环境的组合
func newCounter() func() int {
	// count 是闭包的环境变量
	count := 0
	// 返回的函数可以访问和修改 count
	return func() int {
		count++
		return count
	}
}

// 9. defer 函数
// defer 延迟执行，在函数返回前执行
func deferExample() {
	fmt.Println("函数开始")

	// defer 会在函数返回前执行
	defer fmt.Println("资源清理") // 后进先出

	// defer 可以修改命名返回值
	result := 0
	defer func() {
		result = 100 // 修改命名返回值
		fmt.Println("defer 中修改 result")
	}()

	fmt.Println("函数结束")
	return result
}

// 10. 错误处理函数
func safeDiv(a, b int) (int, error) {
	if b == 0 {
		// 返回错误值
		return 0, fmt.Errorf("除数不能为零")
	}
	return a / b, nil // 返回nil表示没有错误
}

// 11. 递归函数
// 函数调用自身
func factorial(n int) int {
	// 基准条件（递归终止条件）
	if n <= 1 {
		return 1
	}
	// 递归调用
	return n * factorial(n-1)
}

// 12. 斐波那契数列（递归）
func fib(n int) int {
	if n <= 1 {
		return n
	}
	return fib(n-1) + fib(n-2)
}

// 13. 方法（函数与类型的关联）
// 在 Go 中，方法是一种特殊的函数
type Rectangle struct {
	width, height float64
}

// 定义 Rectangle 的方法
func (r Rectangle) Area() float64 {
	return r.width * r.height
}

// 指针接收者方法（可以修改原始值）
func (r *Rectangle) Scale(factor float64) {
	r.width *= factor
	r.height *= factor
}

// ========== 总结 ==========
// 1. 函数定义：func 函数名(参数) 返回类型 { }
// 2. 多返回值：Go 特有功能
// 3. 命名返回值：使代码更清晰
// 4. 可变参数：使用 ...类型
// 5. 闭包：函数 + 引用环境
// 6. defer：延迟执行，后进先出
// 7. 错误处理：返回 error 类型
// 8. 递归：函数调用自身
// 9. 方法：函数与类型的关联
