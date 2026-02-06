package mybase 

import "fmt"

type Person struct {
	Name string
	Age int		
}

type Student struct {
	Person  // 嵌套结构体
	StudentID string
	Score float64
}

// 值接收者：方法内部操作的是Person结构体的副本，不会修改原结构体
func (p Person) GrowUp1() {
	// 示例逻辑：年龄+1
	p.Age++
	fmt.Printf("方法内的副本年龄：%d\n", p.Age)
}

// 2. 补全指针接收者的GrowUp方法（核心）
// 指针接收者：(p *Person) 表示方法接收的是Person结构体的内存地址，操作的是原值
func (p *Person) GrowUp2() {
	p.Age++
	fmt.Printf("方法内的年龄（指针接收者）：%d\n", p.Age)
}

func Lesson8(){
	// 结构体和方法
	// 方式1: 声明并初始化
	p1 := Person {
		Name: "asutoin",
		Age: 22,
	}
	fmt.Printf("p1: %+v\n", p1)
	fmt.Printf("p1.Name: %s\n", p1.Name)
	fmt.Printf("p1.Age: %d\n", p1.Age)

	// 方式2: 直接赋值
	p2 := Person{"asutoin", 22}
	fmt.Printf("\np2: %+v\n", p2)

	// 修改字段
	p1.Age = 26
	fmt.Printf("修改后: %+v\n", p1)

	//方式3: 使用 new 创建（返回指针）
	p3 := new(Person)
	p3.Age = 30
	p3.Name = "Kail"
	fmt.Printf("用 new 创建的 p3: %+v\n", p3)

	// 练习 - 嵌套结构体
	s1 := Student {
		Person: Person{
			Name:"Tom",
			Age: 20,
		},
		StudentID: "S002",
		Score: 95.5,
	}
	fmt.Printf("s1: %+v\n", s1)
	// 嵌套了可以直接跳到 s1.Name 用
	fmt.Printf("Name: %s, StudengID: %s, Score: %.1f\n", s1.Name, s1.StudentID, s1.Score)

	// 值接收者：方法内部操作的是副本
	person1 := Person{Name: "Alice", Age: 30}
	fmt.Printf("调用前: %+v\n", person1)
	person1.GrowUp1() //调用值接收者方法
	fmt.Printf("调用后: %+v\n\n", person1)

	// 指针接收者：方法内部操作的是原值
	person2 := &Person{Name: "Bob", Age: 25}
	fmt.Printf("调用前: %+v\n", person2)
	person2.GrowUp2()
	fmt.Printf("调用后: %+v\n\n", person2)

	
	// 接口: 定义方法签名，具体实现由结构体提供
	people := []Person {
		{Name: "Charlie", Age: 35},
		{Name: "Diana", Age: 28},
	}

	for _, p := range people {
		fmt.Printf("%s 说: %s\n", p.Name, p.Introduce())
	}
	fmt.Println()

	Lesson8Interface()
}

// 定义接口
type Introducer interface {
	Introduce() string
}
// 实现接口方法
func (p Person) Introduce() string {
	return fmt.Sprintf("我叫%s, 今年%d岁", p.Name, p.Age)
}

// 接口的实际应用
func Lesson8Interface() {
	// 接口可以存储任何实现了该接口的类型
	var greeter Introducer

	p1 := Person{Name: "Eva", Age: 32}
	greeter = p1
	fmt.Printf("greeter.Introduce(): %s\n", greeter.Introduce())

	// 使用空接口可以存储任意类型
	var anything interface{}
	anything = 42
	fmt.Printf("anything: %v\n", anything)
	anything = "hello"
	fmt.Printf("anything: %v\n", anything)
	anything = []int{1, 2, 3}
	fmt.Printf("anything: %v\n", anything)
}