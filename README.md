# Go 基础语法详解

本文件夹包含 Go 语言的基础语法学习材料，从入门到进阶。**每个文件都包含详细的中文注释，每句代码都有解释**，非常适合 Go 语言初学者系统学习。

---

## 📁 目录结构

| 文件 | 主题 | 说明 |
|------|------|------|
| [01_hello_world.go](./01_hello_world.go) | Hello World | 第一个 Go 程序 |
| [02_variables.go](./02_variables.go) | 变量声明 | 变量声明、初始化、类型推断 |
| [03_basic_types.go](./03_basic_types.go) | 基本数据类型 | 整型、浮点型、字符串、布尔型 |
| [04_conditions.go](./04_conditions.go) | 条件语句 | if、switch、select 语句 |
| [05_loops.go](./05_loops.go) | 循环语句 | for 循环、range 遍历 |
| [06_functions.go](./06_functions.go) | 函数 | 函数定义、参数、返回值、递归 |
| [07_arrays_slices_maps.go](./07_arrays_slices_maps.go) | 数组、切片、映射 | 复合数据类型详解 |
| [08_structs_methods.go](./08_structs_methods.go) | 结构体和方法 | 结构体定义、方法接收者 |
| [09_concurrency.go](./09_concurrency.go) | 并发编程 | Goroutine、Channel、WaitGroup |
| [10_error_handling.go](./10_error_handling.go) | 错误处理 | 错误定义、捕获、传递 |
| [11_packages_modules.go](./11_packages_modules.go) | 包和模块 | 包组织、模块管理、导入 |

---

## 📦 模块与依赖管理

### 什么是模块（Module）？

**模块 = 项目 = 一个 go.mod 文件**

模块是 Go 语言**依赖管理的基本单位**，相当于其他语言的：
- Java 的 Maven/Gradle 项目
- Node.js 的 npm 包
- Python 的 pip 包

### 模块的作用

| 功能 | 说明 | 示例 |
|------|------|------|
| **标识项目** | 通过模块名标识你的项目 | `module GolangTutorial` |
| **管理依赖** | 自动下载和管理第三方库 | `require github.com/gin-gonic/gin v1.9.1` |
| **版本锁定** | 确保依赖版本一致 | `go.sum` 记录依赖哈希 |
| **可复现** | 任何人克隆代码都能得到相同环境 | `go mod download` |

### 目录结构（只需要一个模块）

```
GolangTutorial/          ← 项目根目录
├── go.mod               ← 只需要这里有 go.mod
├── go.sum               ← 依赖版本校验文件
├── basic_syntax/        ← 不需要 go.mod
├── database/            ← 不需要 go.mod
├── networking/          ← 不需要 go.mod
└── ...
```

---

## 🔧 go mod 命令详解

### 1. go mod init - 初始化模块

**作用：** 只创建 `go.mod` 文件，**不下载任何依赖**

```bash
go mod init GolangTutorial
```

**输出：**
```
go: creating new go.mod: module GolangTutorial
```

**生成的文件：**
```go
// go.mod
module GolangTutorial

go 1.21
```

**什么时候使用？**
- 首次创建项目时
- **只需要执行一次**

### 2. go mod tidy - 整理依赖

**作用：** 扫描代码中的 `import`，下载实际用到的依赖

```bash
go mod tidy
```

**输出：**
```
go: finding modules for packages...
go: downloading packages...
```

**更新后的文件：**
```go
// go.mod
module GolangTutorial

go 1.21

require (
    github.com/gin-gonic/gin v1.9.1
    gorm.io/gorm v1.25.5
)
```

**什么时候使用？**
- 添加了新的第三方库时
- 首次运行需要依赖的代码时
- 需要更新依赖版本时

### 3. go mod download - 下载依赖

**作用：** 下载所有依赖到本地缓存（不会修改 go.mod）

```bash
go mod download
```

### 4. go get - 下载指定包

**作用：** 下载指定版本的包并更新 go.mod

```bash
go get github.com/gin-gonic/gin@v1.9.1
go get github.com/redis/go-redis/v9@latest
```

---

## ❓ 常见问题

### Q1：需要每个目录都创建模块吗？

**答：不需要！** 只需要在**根目录（GolangTutorial/）**创建一个模块。

### Q2：go mod init 和 go mod tidy 的区别？

| 命令 | 作用 | 是否下载依赖 | 何时使用 |
|------|------|-------------|----------|
| `go mod init` | 只创建 go.mod | ❌ 不下载 | 项目首次创建 |
| `go mod tidy` | 扫描并下载依赖 | ✅ 下载 | 添加新依赖时 |

### Q3：为什么 basic_syntax 不需要下载依赖？

因为 **basic_syntax 只使用 Go 标准库**，不需要第三方库：

```go
import (
    "fmt"      // 格式化 I/O（标准库）
    "time"     // 时间处理（标准库）
    "math"     // 数学函数（标准库）
    "strings"  // 字符串操作（标准库）
)
```

### Q4：什么时候需要 go mod tidy？

当代码使用了**第三方库**时才需要：

```go
// 需要下载依赖
import (
    "github.com/gin-gonic/gin"   // Web 框架
    "gorm.io/gorm"               // 数据库 ORM
    "github.com/redis/go-redis/v9" // Redis 客户端
)

# 此时需要执行：
go mod tidy
```

### Q5：国内网络下载依赖失败？

设置国内代理：

```bash
go env -w GOPROXY=https://goproxy.cn,direct
go mod tidy
```

---

## 🚀 完整运行步骤

### 步骤 1：打开终端

**Windows:**
- 按 `Win + R`，输入 `cmd`，回车

**Mac:**
- 按 `Cmd + Space`，搜索 "Terminal"，回车

**Linux:**
- 按 `Ctrl + Alt + T`

### 步骤 2：进入项目目录

```bash
cd Desktop/GolangTutorial/basic_syntax
```

**验证：**
```bash
ls  # 应该看到 01_hello_world.go 等文件
```

### 步骤 3：初始化模块（首次运行需要）

```bash
# 回到根目录
cd ..

# 只在根目录初始化一次！
go mod init GolangTutorial
```

**说明：**
- `go mod init` 只创建 `go.mod`，**不下载依赖**
- basic_syntax 不需要第三方库，所以**不需要** `go mod tidy`

### 步骤 4：运行程序

```bash
cd basic_syntax
go run 01_hello_world.go
```

**输出：**
```
Hello, World!
Go 语言版本: go1.21.0
```

### 步骤 5：运行其他文件

```bash
go run 02_variables.go
go run 03_basic_types.go
```

---

## 📋 常用命令速查表

| 命令 | 说明 | 是否下载依赖 |
|------|------|-------------|
| `go mod init <模块名>` | 初始化模块 | ❌ |
| `go mod tidy` | 下载并整理依赖 | ✅ |
| `go mod download` | 下载所有依赖 | ✅ |
| `go get <包>@版本` | 下载指定包 | ✅ |
| `go run <文件>` | 编译并运行 | ❌ |
| `go build <文件>` | 编译为可执行文件 | ❌ |
| `go env -w GOPROXY=...` | 设置代理 | ❌ |

---

## 📚 学习路径建议

### 第一阶段：基础概念（⭐）
```
01 → 05
```
- Hello World、变量、数据类型、条件、循环

### 第二阶段：复合类型（⭐⭐）
```
06 → 08
```
- 函数、数组/切片/映射、结构体

### 第三阶段：高级特性（⭐⭐⭐）
```
09 → 11
```
- 并发、错误处理、包和模块

---

## 🛠️ 常见问题解决

### 问题 1：'go' 不是内部或外部命令

**解决方法：** 安装 Go 并重启终端
```bash
go version  # 验证安装
```

### 问题 2：go.mod file not found

**解决方法：**
```bash
cd ../
go mod init GolangTutorial
```

### 问题 3：依赖下载失败

**解决方法：**
```bash
go env -w GOPROXY=https://goproxy.cn,direct
go mod tidy
```

---

## 📖 配套资源

| 目录 | 内容 |
|------|------|
| [docs/](../docs/) | 详细文档和进阶教程 |
| [networking/](../networking/) | HTTP/TCP/UDP 网络编程示例 |
| [database/](../database/) | SQL/GORM/Redis 数据库操作 |
| [web/](../web/) | Gin/Echo Web 框架示例 |
| [microservices/](../microservices/) | gRPC 微服务示例 |
| [testing_example.go](../testing_example.go) | 单元测试示例 |

---

## 💡 学习技巧

1. **循序渐进**：按照文件编号顺序学习
2. **动手实践**：每个例子都运行一遍
3. **阅读注释**：每句代码都有中文注释
4. **修改测试**：尝试修改代码观察结果
5. **结合文档**：查看 docs/ 目录下的详细文档

---

## 📝 示例代码结构

```go
// 01_hello_world.go
// 本文件演示 Go 程序的基本结构

package main

import (
    "fmt"  // 标准库，不需要下载依赖
)

// 主函数是程序入口
func main() {
    fmt.Println("Hello, World!")  // 打印到控制台
}
```

---

## 🔄 工作流程图

```
开始
  │
  ▼
进入项目根目录 GolangTutorial/
  │
  ▼
检查是否有 go.mod？
  │
  ├── 否 ──→ go mod init GolangTutorial  （只创建，不下载）
  │
  └── 是
       │
       ▼
进入 basic_syntax 目录
       │
       ▼
运行程序
  │
  ├── go run 01_hello_world.go  （推荐，无需依赖）
  │
  └── go build -o hello ./basic_syntax/01_hello_world.go  （编译）
       │
       ▼
查看输出结果
```
