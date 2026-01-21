// basic_syntax/11_packages_modules.go
// Go 包和模块详解 - 详细注释版

package main

/*
Go 包和模块是组织代码的基本方式。

主要概念：
1. 包（Package）- 代码组织单元
2. 模块（Module）- 依赖管理单元
3. 工作区（Workspace）- 开发环境

包命名规则：
- 包名应简短、有意义
- 使用小写字母，不使用下划线
- 与目录名一致

模块命令：
  go mod init <module_path>     # 初始化模块
  go mod tidy                   # 整理依赖
  go mod download               # 下载依赖
  go mod verify                 # 验证依赖
  go mod edit                   # 编辑模块
*/

import (
	"fmt"
	"math"
	"math/cmplx"
	"math/rand"
	"time"
)

// ====== 包的组织 ======

/*
项目结构示例：
  myproject/
  ├── go.mod              # 模块定义
  ├── main.go             # 主包，主入口
  ├── api/                # API 包
  │   └── handler.go
  ├── models/             # 数据模型包
  │   └── user.go
  ├── utils/              # 工具包
  │   ├── helper.go
  │   └── validator.go
  └── services/           # 服务包
      └── user_service.go
*/

// ====== 导入包 ======

func importExamples() {
	// 1. 单个导入
	// import "fmt"

	// 2. 批量导入
	/*
		import (
			"fmt"
			"math"
			"time"
		)
	*/

	// 3. 点导入（不推荐）
	// import . "fmt"  // 直接使用 Println 而不是 fmt.Println

	// 4. 别名导入
	// import f "fmt"
	// f.Println("Hello")

	// 5. 下划线导入（仅执行 init 函数）
	// import _ "database/sql"  // 注册数据库驱动

	// 6. 函数内导入（局部导入）
	// func localImport() {
	//     import "log"
	//     log.Println("local")
	// }

	// 使用导入的包
	_ = rand.Int()     // math/rand - 随机数
	_ = time.Now()     // time - 时间处理
	_ = math.Sqrt(4)   // math - 数学函数
	_ = cmplx.Sqrt(-1) // math/cmplx - 复数
}

// ====== 导出和访问 ======

/*
在 Go 中，以大写字母开头的标识符是导出的（exported），
可以被其他包访问。

规则：
- 大写字母开头：public（导出）
- 小写字母开头：private（不导出）

示例：
  package mypackage

  var PublicVar = 10     // 导出变量
  var privateVar = 20    // 不导出变量

  func PublicFunc() {}    // 导出函数
  func privateFunc() {}   // 不导出函数
*/

// 导出的变量和函数
var Pi = 3.14159

// Add 是导出的函数
func Add(a, b int) int {
	return a + b
}

// max 是未导出的函数（仅在本包内使用）
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// ====== init 函数 ======

/*
init 函数在包被导入时自动执行：
- 每个包可以有多个 init 函数
- init 函数没有参数和返回值
- 执行顺序：按文件名字母顺序，同一文件按出现顺序

使用场景：
- 初始化包级变量
- 注册驱动程序
- 加载配置文件
*/

// 声明时初始化
var config string = loadConfig()

// init 函数示例 1
func init() {
	fmt.Println("init 函数 1 执行")
}

// init 函数示例 2
func init() {
	fmt.Println("init 函数 2 执行")
}

// loadConfig 模拟加载配置
func loadConfig() string {
	fmt.Println("加载配置...")
	return "default_config"
}

// ====== 模块依赖管理 ======

/*
go.mod 文件示例：
  module myproject

  go 1.21

  require (
      github.com/gin-gonic/gin v1.9.1
      gorm.io/gorm v1.25.5
  )

  replace (
      github.com/example/old => ./local/path
  )

  exclude (
      github.com/bad/dependency v1.0.0
  )
*/

// ====== 常用标准库包 ======

func standardLibraryExamples() {
	// fmt - 格式化 I/O
	fmt.Printf("格式化输出: %.2f\n", 3.14159)
	fmt.Println("Hello, World!")

	// strconv - 字符串转换
	// strconv.Atoi("123")  // 字符串转整数
	// strconv.Itoa(123)    // 整数转字符串

	// strings - 字符串操作
	// strings.Contains("hello", "ell")
	// strings.ToLower("HELLO")
	// strings.Split("a,b,c", ",")

	// os - 操作系统功能
	// os.Getwd()  // 获取当前目录
	// os.Exit(0)  // 退出程序

	// io - I/O 接口
	// io.Reader, io.Writer

	// sync - 同步原语
	// sync.Mutex, sync.WaitGroup, sync.Once

	// context - 上下文管理
	// context.Background(), context.WithCancel()

	// reflect - 反射
	// reflect.TypeOf(), reflect.ValueOf()

	// encoding/json - JSON 处理
	// json.Marshal(), json.Unmarshal()
}

// ====== 自定义包示例 ======

/*
假设我们有一个 utils 包：
  package utils

  // utils/helper.go

  // IsEmailValid 验证邮箱格式
  func IsEmailValid(email string) bool {
      // 简化的验证
      return len(email) > 3 && contains(email, '@')
  }

  // contains 检查字符串是否包含子串
  func contains(s, substr string) bool {
      return len(s) >= len(substr) &&
             (s == substr ||
              containsPrefix(s, substr) ||
              containsSuffix(s, substr))
  }

  // containsPrefix 检查前缀
  func containsPrefix(s, prefix string) bool {
      for i := 0; i <= len(s)-len(prefix); i++ {
          if s[i:i+len(prefix)] == prefix {
              return true
          }
      }
      return false
  }

  // containsSuffix 检查后缀
  func containsSuffix(s, suffix string) bool {
      if len(s) < len(suffix) {
          return false
      }
      return s[len(s)-len(suffix):] == suffix
  }
*/

// ====== go mod 命令 ======

/*
常用命令：
  go mod init <module>        # 初始化模块
  go mod tidy                 # 整理依赖
  go mod download             # 下载依赖到缓存
  go mod verify               # 验证依赖
  go mod graph                # 显示依赖图
  go mod why <package>        # 解释为什么需要某个包
  go mod edit -go=1.21        # 修改 go 版本
  go list -m all              # 列出所有模块
  go get <package>@latest     # 获取最新版本
  go get <package>@v1.2.3     # 获取指定版本
  go get <package>@upgrade    # 升级到最新版本
  go get <package>@downgrade  # 降级到最新版本
*/

// ====== Go 工作区（Workspace）=====

/*
Go 1.18+ 支持工作区，可以同时开发多个模块。

go.work 文件示例：
  go 1.21

  use (
      ./myproject
      ./mylib
  )

工作区命令：
  go work init
  go work use ./module1 ./module2
  go work edit -go=1.21
*/

// ====== 包的可见性规则 ======

/*
可见性规则：
1. 同一包内所有文件可以互相访问
2. 导出的标识符以大写字母开头
3. 未导出的标识符只能在本包内使用
4. 结构体字段同理（大写导出，小写不导出）

示例：
  package foo

  type MyStruct struct {
      PublicField  int    // 导出字段
      privateField string // 私有字段
  }

  func (m *MyStruct) PublicMethod() {}    // 导出方法
  func (m *MyStruct) privateMethod() {}   // 私有方法
*/

// MyStruct 演示可见性
type MyStruct struct {
	PublicField  int    // 导出字段
	privateField string // 私有字段
}

// NewMyStruct 构造函数（导出）
func NewMyStruct(public int, private string) *MyStruct {
	return &MyStruct{
		PublicField:  public,
		privateField: private,
	}
}

// PublicMethod 导出方法
func (m *MyStruct) PublicMethod() {
	fmt.Println("PublicMethod called")
}

// privateMethod 私有方法
func (m *MyStruct) privateMethod() {
	fmt.Println("privateMethod called")
}

// ====== Go 泛型（1.18+）=====

/*
泛型允许编写适用于多种类型的代码。

语法：
  func [T Type] FunctionName(param T) T { ... }

  type [T Type] StructType struct {
      Field T
  }
*/

// GenericFunction 泛型函数
func GenericFunction[T any](value T) T {
	return value
}

// Addable 接口约束
type Addable interface {
	int | int32 | int64 | float64 | string
}

// Add 泛型加法函数
func Add[T Addable](a, b T) T {
	return a + b
}

// MapFunction 泛型映射函数
func MapFunction[T any, U any](slice []T, f func(T) U) []U {
	result := make([]U, len(slice))
	for i, v := range slice {
		result[i] = f(v)
	}
	return result
}

// ====== 依赖管理最佳实践 ======

/*
最佳实践：
1. 使用语义版本控制（semver）
2. 定期更新依赖
3. 锁定依赖版本（go.sum）
4. 使用 replace 进行本地开发
5. 避免使用太宽松的版本范围
6. 定期运行 go mod tidy

版本号规则：
  vMAJOR.MINOR.PATCH
  - MAJOR: 不兼容的 API 变更
  - MINOR: 新功能（向后兼容）
  - PATCH: 修复 bug（向后兼容）
*/

// ====== 主函数 ======

func main() {
	fmt.Println("=== Go 包和模块详解 ===")

	// 1. 包导入示例
	fmt.Println("\n--- 包导入 ---")
	importExamples()

	// 2. 导出和访问示例
	fmt.Println("\n--- 导出和访问 ---")
	fmt.Println("Pi =", Pi)
	fmt.Println("Add(1, 2) =", Add(1, 2))
	fmt.Println("max(3, 5) =", max(3, 5))

	// 3. init 函数示例
	fmt.Println("\n--- init 函数 ---")
	fmt.Println("main 函数开始")

	// 4. 自定义包示例
	fmt.Println("\n--- 自定义包示例 ---")
	_ = NewMyStruct(42, "private")

	// 5. 泛型示例
	fmt.Println("\n--- 泛型示例 ---")
	fmt.Println("GenericFunction(42) =", GenericFunction(42))
	fmt.Println("GenericFunction(\"hello\") =", GenericFunction("hello"))
	fmt.Println("Add(1, 2) =", Add(1, 2))
	fmt.Println("Add(1.5, 2.5) =", Add(1.5, 2.5))
	fmt.Println("Add(\"Hello, \", \"World!\") =", Add("Hello, ", "World!"))

	// 6. 映射示例
	fmt.Println("\n--- 映射示例 ---")
	nums := []int{1, 2, 3, 4, 5}
	doubled := MapFunction(nums, func(n int) int { return n * 2 })
	fmt.Println("MapFunction:", doubled)

	// 7. 标准库示例
	fmt.Println("\n--- 标准库示例 ---")
	standardLibraryExamples()

	fmt.Println("\n包和模块示例完成")
}
