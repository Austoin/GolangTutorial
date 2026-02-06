package mybase 

import "fmt"

func sayHello() {
	fmt.Println("Hello from Lesson 6!")
}

// 有参数
func greet(name string) {
	fmt.Println("Hello " + name)
}

// 有返回值
// func add(a int, b int) int {
// 	return a + b
// }

// 多返回值
func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, fmt.Errorf("除数不能为零")
	}
	return a/b, nil // 返回商和 nil 错误
}

// 参数传递
func add(a, b int) int { return a + b }
// 变长参数
func sum(nums ...int) int {
	total := 0
	for _, num := range nums{
		total += num
	}
	return total
}

// 递归
func factorial(n int) int { //阶乘
	if n <= 1 {
		return 1
	}
	return n * factorial(n-1)
}

// 斐波那契
func fibonacci(n int) int {
	if n <= 1{
		return n
	}
	return fibonacci(n-1) + fibonacci(n-2)
}

func Lesson6() {
	sayHello()
	greet("Austoin")

	// sum := add(1, 2)
	// fmt.Printf("1 + 2 = %d\n", sum)

	quotient, err := divide(10, 2)
	if err != nil {
		fmt.Println("错误：", err)
	} else {
		fmt.Printf("10 / 2 = %.2f\n", quotient)
	}

	fmt.Printf("sum(1, 2, 3, 4) = %d\n", sum(1, 2, 3, 4))


}