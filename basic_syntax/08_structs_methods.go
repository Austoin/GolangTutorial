package main

// 本文件演示 Go 语言的结构体和方法

import "fmt"

func main() {
	// ========== 结构体 ==========

	// 1. 结构体声明
	// 使用 type 关键字定义新的类型
	// 结构体是值类型，包含多个不同类型的字段

	// 定义 Person 结构体
	type Person struct {
		Name    string  // 姓名，字符串类型
		Age     int     // 年龄，整数类型
		Height  float64 // 身高，浮点类型
		Weight  float64 // 体重，浮点类型
		Email   string  // 邮箱，字符串类型
	}

	// 2. 创建结构体实例

	// 方式1：声明变量并初始化所有字段
	var p1 Person
	p1.Name = "张三"
	p1.Age = 25
	p1.Height = 175.5
	p1.Weight = 70.0
	p1.Email = "zhangsan@example.com"
	fmt.Printf("p1: %+v\n", p1)  // %+v 显示字段名和值

	// 方式2：简短声明，初始化所有字段
	p2 := Person{
		Name:   "李四",
		Age:    30,
		Height: 180.0,
		Weight: 75.0,
		Email:  "lisi@example.com",
	}
	fmt.Printf("p2: %+v\n", p2)

	// 方式3：只初始化部分字段（其余为类型零值）
	p3 := Person{Name: "王五"}
	fmt.Printf("p3: %+v\n", p3)  // Age, Height, Weight, Email 为零值

	// 方式4：按顺序初始化（不推荐，字段顺序变化会出错）
	p4 := Person{"赵六", 28, 165.5, 60.0, "zhaoliu@example.com"}
	fmt.Printf("p4: %+v\n", p4)

	// 3. 访问结构体字段
	fmt.Printf("p1.Name: %s\n", p1.Name)
	fmt.Printf("p1.Age: %d\n", p1.Age)

	// 4. 修改结构体字段
	p1.Age = 26
	fmt.Printf("修改后 p1: %+v\n", p1)

	// 5. 结构体是值类型
	// 赋值会复制整个结构体
	p5 := p2  // 复制 p2 到 p5
	p5.Age = 40
	fmt.Printf("p2: %+v\n", p2)  // p2 不变
	fmt.Printf("p5: %+v\n", p5)  // p5 的 Age 被修改

	// 6. 结构体指针
	// 使用 & 获取结构体地址
	ptr := &p1
	fmt.Printf("p1 的地址: %p\n", ptr)
	fmt.Printf("通过指针访问 Name: %s\n", ptr.Name)  // 自动解引用

	// 创建指针类型结构体
	p6 := &Person{Name: "孙七"}
	fmt.Printf("p6: %+v\n", p6)

	// 7. 使用 new 创建结构体指针
	// new(Type) 返回指向类型零值的指针
	p7 := new(Person)
	p7.Name = "周八"
	p7.Age = 35
	fmt.Printf("p7: %+v\n", p7)

	// 8. 匿名结构体
	// 临时使用，不需要先定义类型
	anon := struct {
		X int
		Y string
	}{
		X: 100,
		Y: "匿名",
	}
	fmt.Printf("匿名结构体: %+v\n", anon)

	// ========== 方法 ==========

	// 9. 方法定义
	// 方法是与特定类型关联的函数
	// 格式：func (接收者) 方法名(参数) 返回类型 { }

	// 定义 Rectangle 结构体
	type Rectangle struct {
		Width  float64
		Height float64
	}

	// 值接收者方法
	// 接收者是结构体的副本
	// 不能修改原始结构体的值
	// (r Rectangle) 表示 r 是 Rectangle 类型的接收者
	// 这种方法不会修改原始的 Rectangle
	// 因为传递的是值的副本
	// 适用于不需要修改原始数据的情况
	// 代码示例：
	// rect := Rectangle{10, 20}
	// area := rect.Area()  // 返回 200.0
	// 但不会改变 rect 的原始值
	// 因为 rect.Area() 接收的是 rect 的副本
	// 不会修改原始的 rect 结构体
	// 这就是值接收者方法的特点
	// 传递的是值的拷贝，不会影响原始数据
	// 适用于只读操作或不需要修改原始结构体的场景

	rect := Rectangle{Width: 10, Height: 20}
	area := rect.Area()
	fmt.Printf("矩形面积: %.2f\n", area)

	// 指针接收者方法
	// 接收者是指针，指向原始结构体
	// 可以修改原始结构体的值
	// 传递指针避免了复制大结构体的开销
	rect.Scale(2.0)
	fmt.Printf("缩放后 rect: %.2f x %.2f\n", rect.Width, rect.Height)

	// 10. 指针接收者 vs 值接收者
	// 指针接收者：
	//   - 可以修改原始值
	//   - 避免大结构体复制的开销
	//   - 需要确保指针不为 nil
	//
	// 值接收者：
	//   - 操作副本，不影响原始值
	//   - 简单，不需要检查 nil
	//   - 适合只读操作

	// 11. 方法与函数的区别
	// 函数：独立的代码块
	// 方法：与特定类型关联

	// 函数调用
	area2 := calculateArea(rect.Width, rect.Height)
	fmt.Printf("函数计算面积: %.2f\n", area2)

	// 方法调用
	area3 := rect.Area()
	fmt.Printf("方法计算面积: %.2f\n", area3)

	// ========== 结构体嵌入（模拟继承） ==========

	// 12. 结构体嵌入（匿名嵌入）
	// 类似其他语言的"继承"
	// Go 使用组合而非继承

	type Address struct {
		City    string
		State   string
		Country string
	}

	type Employee struct {
		Name    string
		Salary  float64
		Address // 匿名嵌入 Address
		// Address 中的字段 City, State, Country
		// 直接成为 Employee 的字段
	}

	emp := Employee{
		Name:   "张三",
		Salary: 50000,
		Address: Address{
			City:    "北京",
			State:   "北京",
			Country: "中国",
		},
	}

	fmt.Printf("员工: %s, 城市: %s\n", emp.Name, emp.City)  // 直接访问嵌入字段
	fmt.Printf("员工: %s, 国家: %s\n", emp.Name, emp.Address.Country)

	// 13. 方法继承
	// 嵌入结构体的方法也会被"继承"

	// 为 Address 添加方法
	(a Address) FullAddress() string {
		return fmt.Sprintf("%s, %s, %s", a.City, a.State, a.Country)
	}

	fmt.Printf("完整地址: %s\n", emp.FullAddress())

	// 14. 方法重写
	// 如果Employee有同名方法，会覆盖Address的方法

	type Manager struct {
		Employee
		Department string
	}

	// Employee 也有一个 PrintInfo 方法
	emp.PrintInfo()
	mgr := Manager{
		Employee: Employee{
			Name:   "李经理",
			Salary: 80000,
			Address: Address{
				City:    "上海",
				State:   "上海",
				Country: "中国",
			},
		},
		Department: "技术部",
	}

	mgr.PrintInfo()  // 调用 Manager 的 PrintInfo
	// emp.PrintInfo()  // 调用 Employee 的 PrintInfo

	// 15. 结构体标签（Tag）
	// 用于反射、JSON 序列化等

	type User struct {
		ID       int    `json:"id"`       // JSON 序列化为 "id"
		Username string `json:"username"` // JSON 序列化为 "username"
		Password string `json:"-"`        // 忽略此字段
		Email    string `json:"email,omitempty"`  // 空值时忽略
	}

	user := User{
		ID:       1,
		Username: "admin",
		Password: "secret",
		Email:    "",
	}

	fmt.Printf("User 结构体: %+v\n", user)
}

// ========== Rectangle 的方法 ==========

// Area 计算矩形面积
// 值接收者方法，不修改原始结构体
// 当调用 rect.Area() 时，r 是 rect 的副本
func (r Rectangle) Area() float64 {
	// r 是 Rectangle 的副本
	// 不能修改原始的 rect
	return r.Width * r.Height
}

// Scale 缩放矩形
// 指针接收者方法，修改原始结构体
// 当调用 rect.Scale(2.0) 时，r 指向原始 rect
// 可以修改 rect 的 Width 和 Height
func (r *Rectangle) Scale(factor float64) {
	// 通过指针修改原始值
	r.Width *= factor
	r.Height *= factor
}

// calculateArea 独立函数
// 计算矩形面积，不关联任何类型
func calculateArea(width, height float64) float64 {
	return width * height
}

// ========== Employee 的方法 ==========

// PrintInfo 打印员工信息
func (e Employee) PrintInfo() {
	fmt.Printf("员工: %s, 薪资: %.2f, 地址: %s\n", 
		e.Name, e.Salary, e.FullAddress())
}

// ========== Manager 的方法 ==========

// PrintInfo 重写 Employee 的 PrintInfo
func (m Manager) PrintInfo() {
	fmt.Printf("经理: %s, 部门: %s, 薪资: %.2f, 地址: %s\n",
		m.Name, m.Department, m.Salary, m.FullAddress())
}

// ========== 总结 ==========
// 1. 结构体是复合类型，包含多个字段
// 2. 结构体是值类型，赋值会复制
// 3. 方法是与特定类型关联的函数
// 4. 值接收者 vs 指针接收者
// 5. 结构体嵌入实现代码复用
// 6. 结构体标签用于元信息
