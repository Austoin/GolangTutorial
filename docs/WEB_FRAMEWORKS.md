# Go Web 框架详解

## 目录
- [1. Web 框架概述](#1-web-框架概述)
- [2. Gin 框架](#2-gin-框架)
- [3. Echo 框架](#3-echo-框架)
- [4. 框架对比](#4-框架对比)
- [5. 最佳实践](#5-最佳实践)

---

## 1. Web 框架概述

### 1.1 Go Web 框架生态

```
┌─────────────────────────────────────────────────────────────┐
│                      Go Web 框架                             │
├─────────────┬─────────────┬─────────────┬─────────────────┤
│   Gin       │   Echo      │   Fiber     │   Chi           │
│  (最流行)   │  (高性能)   │  (Express)  │  (轻量级)       │
├─────────────┼─────────────┼─────────────┼─────────────────┤
│  Martini    │  Buffalo    │  Go-Chi     │  Iris           │
│  (已停止)   │  (全栈)     │  (路由)     │  (自称最快)     │
└─────────────┴─────────────┴─────────────┴─────────────────┘
```

### 1.2 框架选择建议

| 框架 | 适用场景 | 特点 |
|------|---------|------|
| **Gin** | API 服务、高并发项目 | 路由速度快、中间件丰富、生态完善 |
| **Echo** | 高性能 API、Web 应用 | 极简设计、自动 TLS、内置验证 |
| **Fiber** | 迁移自 Node.js | Express 风格、内存占用低 |
| **Chi** | 微服务 | 轻量、兼容 net/http |
| **Buffalo** | 全栈 Web 应用 | 包含前端工作流 |

---

## 2. Gin 框架

### 2.1 安装和基础用法

```bash
# 安装
go get -u github.com/gin-gonic/gin
```

```go
package main

import (
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
)

// Person 结构体用于绑定 JSON
type Person struct {
    Name    string `json:"name" binding:"required"`
    Age     int    `json:"age" binding:"gte=0,lte=150"`
    Email   string `json:"email" binding:"email"`
}

func main() {
    // 创建 Gin 路由引擎
    engine := gin.Default()

    // 1. 基础路由
    engine.GET("/", func(c *gin.Context) {
        c.String(http.StatusOK, "Hello, Gin!")
    })

    // 2. 路由参数
    engine.GET("/user/:name", func(c *gin.Context) {
        name := c.Param("name")
        c.String(http.StatusOK, "Hello, %s!", name)
    })

    // 3. 查询参数
    engine.GET("/search", func(c *gin.Context) {
        name := c.Query("name")           // c.Request.URL.Query().Get("name")
        page := c.DefaultQuery("page", "1") // 带默认值
        c.String(http.StatusOK, "Name: %s, Page: %s", name, page)
    })

    // 4. POST JSON 参数
    engine.POST("/json", func(c *gin.Context) {
        var person Person
        
        // 绑定 JSON 到结构体
        if err := c.ShouldBindJSON(&person); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{
                "error": err.Error(),
            })
            return
        }

        c.JSON(http.StatusOK, gin.H{
            "message": "success",
            "data":    person,
        })
    })

    // 5. POST 表单参数
    engine.POST("/form", func(c *gin.Context) {
        username := c.PostForm("username")
        password := c.DefaultPostForm("password", "default")

        c.JSON(http.StatusOK, gin.H{
            "username": username,
            "password": password,
        })
    })

    // 6. 路径参数 + 查询参数
    engine.GET("/posts/:id/comments/:comment_id", func(c *gin.Context) {
        postID := c.Param("id")
        commentID := c.Param("comment_id")
        c.String(http.StatusOK, "Post: %s, Comment: %s", postID, commentID)
    })

    // 7. 上传文件
    engine.POST("/upload", func(c *gin.Context) {
        file, err := c.FormFile("file")
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        // 保存文件
        c.SaveUploadedFile(file, "./"+file.Filename)
        
        c.JSON(http.StatusOK, gin.H{
            "filename": file.Filename,
            "size":     file.Size,
        })
    })

    // 8. 多文件上传
    engine.POST("/upload/multiple", func(c *gin.Context) {
        form, _ := c.MultipartForm()
        files := form.File["files[]"]

        for _, file := range files {
            c.SaveUploadedFile(file, "./"+file.Filename)
        }

        c.JSON(http.StatusOK, gin.H{
            "count":   len(files),
            "message": "uploaded",
        })
    })

    // 9. 重定向
    engine.GET("/old", func(c *gin.Context) {
        c.Redirect(http.StatusMovedPermanently, "/new")
    })

    // 10. 异步请求
    engine.GET("/async", func(c *gin.Context) {
        // 复制上下文
        cCp := c.Copy()
        
        go func() {
            time.Sleep(2 * time.Second)
            cCp.String(http.StatusOK, "Async response")
        }()
        
        c.String(http.StatusOK, "Request accepted")
    })

    // 启动服务器
    engine.Run(":8080")
}
```

### 2.2 路由组

```go
package main

import (
    "github.com/gin-gonic/gin"
)

func main() {
    engine := gin.Default()

    // 创建路由组
    v1 := engine.Group("/api/v1")
    v2 := engine.Group("/api/v2")

    // v1 路由组
    v1.GET("/users", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "v1 users"})
    })
    v1.POST("/users", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "create v1 user"})
    })

    // v2 路由组（带中间件）
    v2.Use(middleware1())
    v2.GET("/users", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "v2 users"})
    })
    v2.POST("/users", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "create v2 user"})
    })

    // 嵌套路由组
    admin := v1.Group("/admin")
    admin.GET("/dashboard", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "admin dashboard"})
    })

    engine.Run(":8080")
}

func middleware1() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 处理中间件逻辑
        c.Next()
    }
}
```

### 2.3 中间件

```go
package main

import (
    "fmt"
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
)

func main() {
    engine := gin.New()

    // 1. 使用内置中间件
    engine.Use(gin.Recovery())      // 恢复 Panic
    engine.Use(gin.Logger())         // 日志
    engine.Use(CORSMiddleware())    // 自定义 CORS

    // 2. 自定义中间件：请求日志
    func() {
        engine.Use(func(c *gin.Context) {
            start := time.Now()
            path := c.Request.URL.Path

            // 处理请求
            c.Next()

            // 记录日志
            latency := time.Since(start)
            status := c.Writer.Status()
            fmt.Printf("%s %s %d %v\n", 
                c.Request.Method, path, status, latency)
        })
    }()

    // 3. 认证中间件
    engine.Use(func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                "error": "missing token",
            })
            return
        }
        // 验证 token...
        c.Next()
    })

    // 4. 限流中间件
    engine.Use(rateLimitMiddleware(100)) // 每秒 100 请求

    // 5. 业务中间件
    engine.GET("/protected", func(c *gin.Context) {
        userID := c.GetString("user_id")
        c.JSON(200, gin.H{"user_id": userID})
    })

    engine.Run(":8080")
}

// CORS 中间件
func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        c.Writer.Header().Set("Access-Control-Allow-Headers", 
            "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
        c.Writer.Header().Set("Access-Control-Allow-Methods", 
            "POST, OPTIONS, GET, PUT, DELETE")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }

        c.Next()
    }
}

// 简单限流中间件
func rateLimitMiddleware(requestsPerSecond int) gin.HandlerFunc {
    limiter := time.NewTicker(time.Second / time.Duration(requestsPerSecond))
    
    return func(c *gin.Context) {
        <-limiter.C
        c.Next()
    }
}
```

### 2.4 参数绑定和验证

```go
package main

import (
    "net/http"
    "strconv"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/go-playground/validator/v10"
)

// LoginRequest 登录请求
type LoginRequest struct {
    Username string `json:"username" binding:"required,min=3,max=20"`
    Password string `json:"password" binding:"required,min=6"`
}

// UserRequest 用户请求（带自定义验证）
type UserRequest struct {
    Name     string    `json:"name" binding:"required"`
    Email    string    `json:"email" binding:"required,email"`
    Age      int       `json:"age" binding:"gte=0,lte=150"`
    Birthday time.Time `json:"birthday" binding:"required"`
    Phone    string    `json:"phone" binding:"required,e164"` // E.164 格式
}

// RegisterRequest 注册请求
type RegisterRequest struct {
    Username        string `json:"username" binding:"required,min=3,max=20,alphanum"`
    Email           string `json:"email" binding:"required,email"`
    Password        string `json:"password" binding:"required,min=6"`
    ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
}

// 自定义验证器
var validate = validator.New()

func init() {
    // 注册自定义验证器
    validate.RegisterValidation("phone", validatePhone)
}

func validatePhone(fl validator.FieldLevel) bool {
    phone := fl.Field().String()
    return len(phone) == 11
}

func main() {
    engine := gin.Default()

    // ShouldBind 系列方法（推荐）
    engine.POST("/login", func(c *gin.Context) {
        var req LoginRequest
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{
                "error": err.Error(),
            })
            return
        }

        // 验证
        if err := validate.Struct(req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{
                "error": err.Error(),
            })
            return
        }

        c.JSON(http.StatusOK, gin.H{
            "message": "login success",
            "user":    req.Username,
        })
    })

    // Bind 系列方法（不推荐，会设置 400 状态码）
    engine.POST("/bind", func(c *gin.Context) {
        var req LoginRequest
        c.BindJSON(&req) // 如果验证失败会自动返回 400
        // ...
    })

    // 手动验证
    engine.POST("/register", func(c *gin.Context) {
        var req RegisterRequest
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{
                "error": err.Error(),
            })
            return
        }

        c.JSON(http.StatusOK, gin.H{
            "message": "register success",
        })
    })

    // 绑定 URI 参数
    engine.GET("/user/:id", func(c *gin.Context) {
        type URI struct {
            ID int `uri:"id" binding:"required,min=1"`
        }
        var uri URI
        if err := c.ShouldBindUri(&uri); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusOK, gin.H{"id": uri.ID})
    })

    // 绑定查询参数
    engine.GET("/search", func(c *gin.Context) {
        type Query struct {
            Page     int    `form:"page" binding:"min=1"`
            PageSize int    `form:"page_size" binding:"min=1,max=100"`
            Keyword  string `form:"keyword"`
        }
        var query Query
        if err := c.ShouldBindQuery(&query); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusOK, query)
    })

    // 绑定 Header
    engine.GET("/header", func(c *gin.Context) {
        type Header struct {
            ContentType string `header:"Content-Type"`
            Authorization string `header:"Authorization"`
        }
        var header Header
        c.ShouldBindHeader(&header)
        c.JSON(http.StatusOK, header)
    })

    engine.Run(":8080")
}
```

### 2.5 文件上传和处理

```go
package main

import (
    "fmt"
    "io"
    "net/http"
    "os"
    "path/filepath"

    "github.com/gin-gonic/gin"
)

func main() {
    engine := gin.Default()

    // 单文件上传
    engine.POST("/upload", func(c *gin.Context) {
        file, err := c.FormFile("file")
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        // 创建上传目录
        uploadDir := "./uploads"
        os.MkdirAll(uploadDir, 0755)

        // 生成文件名
        filename := fmt.Sprintf("%d_%s", 
            time.Now().UnixNano(), file.Filename)
        filepath := filepath.Join(uploadDir, filename)

        // 保存文件
        if err := c.SaveUploadedFile(file, filepath); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "error": err.Error(),
            })
            return
        }

        c.JSON(http.StatusOK, gin.H{
            "filename": filename,
            "size":     file.Size,
            "path":     filepath,
        })
    })

    // 多文件上传
    engine.POST("/upload/multiple", func(c *gin.Context) {
        form, err := c.MultipartForm()
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        files := form.File["files"]
        var uploaded []map[string]interface{}

        for _, file := range files {
            filename := fmt.Sprintf("%d_%s", 
                time.Now().UnixNano(), file.Filename)
            filepath := filepath.Join("./uploads", filename)
            
            c.SaveUploadedFile(file, filepath)

            uploaded = append(uploaded, map[string]interface{}{
                "filename": filename,
                "size":     file.Size,
            })
        }

        c.JSON(http.StatusOK, gin.H{
            "count":  len(files),
            "files":  uploaded,
        })
    })

    // 流式上传（适合大文件）
    engine.POST("/upload/stream", func(c *gin.Context) {
        file, err := c.FormFile("file")
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        src, err := file.Open()
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "error": err.Error(),
            })
            return
        }
        defer src.Close()

        // 创建目标文件
        dst, err := os.Create("./uploads/" + file.Filename)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "error": err.Error(),
            })
            return
        }
        defer dst.Close()

        // 复制文件
        if _, err := io.Copy(dst, src); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "error": err.Error(),
            })
            return
        }

        c.JSON(http.StatusOK, gin.H{"message": "uploaded"})
    })

    // 提供静态文件服务
    engine.Static("/static", "./public")
    engine.StaticFile("/favicon.ico", "./public/favicon.ico")

    engine.Run(":8080")
}
```

### 2.6 响应渲染

```go
package main

import (
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
)

func main() {
    engine := gin.Default()

    // 1. JSON 响应
    engine.GET("/json", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "message": "success",
            "data": map[string]interface{}{
                "name":  "张三",
                "age":   25,
                "items": []string{"a", "b", "c"},
            },
        })
    })

    // 2. XML 响应
    engine.GET("/xml", func(c *gin.Context) {
        c.XML(http.StatusOK, gin.H{
            "status": "ok",
            "data":   "some data",
        })
    })

    // 3. YAML 响应
    engine.GET("/yaml", func(c *gin.Context) {
        c.YAML(http.StatusOK, gin.H{
            "status": "ok",
            "data":   "some data",
        })
    })

    // 4. ProtoBuf 响应
    engine.GET("/protobuf", func(c *gin.Context) {
        // 需要定义 protobuf 结构
        c.ProtoBuf(http.StatusOK, nil)
    })

    // 5. HTML 渲染
    engine.LoadHTMLGlob("templates/*")
    engine.GET("/html", func(c *gin.Context) {
        c.HTML(http.StatusOK, "index.html", gin.H{
            "title": "Hello Gin",
            "data":  "World",
        })
    })

    // 6. 重定向
    engine.GET("/redirect", func(c *gin.Context) {
        c.Redirect(http.StatusMovedPermanently, "https://example.com")
    })

    // 7. 纯文本响应
    engine.GET("/text", func(c *gin.Context) {
        c.String(http.StatusOK, "Hello, %s!", "Gin")
    })

    // 8. 文件下载
    engine.GET("/download", func(c *gin.Context) {
        c.Header("Content-Disposition", 
            "attachment; filename=file.txt")
        c.File("./uploads/file.txt")
    })

    // 9. 流式响应
    engine.GET("/stream", func(c *gin.Context) {
        c.Stream(func(w io.Writer) bool {
            w.Write([]byte("chunk "))
            time.Sleep(time.Second)
            w.Write([]byte("data\n"))
            return false // 返回 false 结束流
        })
    })

    // 10. Data 响应（字节数组）
    engine.GET("/data", func(c *gin.Context) {
        c.Data(http.StatusOK, "application/octet-stream", 
            []byte("binary data"))
    })

    engine.Run(":8080")
}
```

### 2.7 Gin 最佳实践

```go
package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"

    "github.com/gin-gonic/gin"
    "gopkg.in/natefinch/lumberjack.v2"
)

// Config 应用配置
type Config struct {
    Port    string
    Mode    string // debug, release, test
    LogFile string
}

// Response 统一响应结构
type Response struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}

// HandlerFunc 统一处理器类型
type HandlerFunc func(*gin.Context) Response

func main() {
    config := Config{
        Port:    ":8080",
        Mode:    "release",
        LogFile: "./logs/app.log",
    }

    // 设置 Gin 模式
    gin.SetMode(config.Mode)

    // 配置日志
    gin.DefaultWriter = &lumberjack.Logger{
        Filename:   config.LogFile,
        MaxSize:    100, // MB
        MaxBackups: 7,
        MaxAge:     30,  // days
        Compress:   true,
    }

    engine := gin.New()
    engine.Use(gin.Recovery())
    engine.Use(requestLogger())

    // 健康检查
    engine.GET("/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, Response{
            Code:    0,
            Message: "ok",
        })
    })

    // API 路由
    v1 := engine.Group("/api/v1")
    v1.GET("/users", listUsers)
    v1.GET("/users/:id", getUser)
    v1.POST("/users", createUser)

    // 优雅关闭
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    go func() {
        <-quit
        log.Println("Shutting down server...")
    }()

    log.Printf("Server starting on %s", config.Port)
    engine.Run(config.Port)
}

func requestLogger() gin.HandlerFunc {
    return gin.LoggerWithConfig(gin.LoggerConfig{
        SkipPaths: []string{"/health"},
    })
}

// 统一错误处理
func handleError(c *gin.Context, err error) {
    log.Printf("Error: %v", err)
    c.JSON(http.StatusInternalServerError, Response{
        Code:    500,
        Message: err.Error(),
    })
}

// 示例处理器
func listUsers(c *gin.Context) Response {
    page := c.DefaultQuery("page", "1")
    pageSize := c.DefaultQuery("page_size", "10")
    return Response{
        Code:    0,
        Message: "success",
        Data: map[string]interface{}{
            "page":      page,
            "page_size": pageSize,
            "items":     []string{"user1", "user2"},
        },
    }
}

func getUser(c *gin.Context) Response {
    id := c.Param("id")
    return Response{
        Code:    0,
        Message: "success",
        Data: map[string]interface{}{
            "id":   id,
            "name": "张三",
        },
    }
}

func createUser(c *gin.Context) Response {
    var input map[string]interface{}
    c.ShouldBindJSON(&input)
    return Response{
        Code:    0,
        Message: "created",
        Data:    input,
    }
}

// 使用统一处理器
func useHandler(handler HandlerFunc) gin.HandlerFunc {
    return func(c *gin.Context) {
        resp := handler(c)
        c.JSON(http.StatusOK, resp)
    }
}
```

---

## 3. Echo 框架

### 3.1 安装和基础用法

```bash
# 安装
go get -u github.com/labstack/echo/v4
```

```go
package main

import (
    "net/http"
    "strconv"

    "github.com/labstack/echo/v4"
)

type User struct {
    ID       int    `json:"id" param:"id"`
    Name     string `json:"name" query:"name"`
    Username string `json:"username"`
    Password string `json:"-"`
}

func main() {
    // 创建 Echo 实例
    e := echo.New()
    e.HideBanner = true

    // 1. 基础路由
    e.GET("/", func(c echo.Context) error {
        return c.String(http.StatusOK, "Hello, Echo!")
    })

    // 2. 路径参数
    e.GET("/user/:id", func(c echo.Context) error {
        id := c.Param("id")
        return c.String(http.StatusOK, "User ID: "+id)
    })

    // 3. 查询参数
    e.GET("/search", func(c echo.Context) error {
        name := c.QueryParam("name")
        page := c.QueryParam("page")
        return c.String(http.StatusOK, 
            "Name: "+name+", Page: "+page)
    })

    // 4. POST 请求
    e.POST("/users", func(c echo.Context) error {
        var user User
        if err := c.Bind(&user); err != nil {
            return err
        }
        return c.JSON(http.StatusOK, user)
    })

    // 5. 路径参数 + 查询参数
    e.GET("/posts/:id/comments/:comment_id", func(c echo.Context) error {
        postID := c.Param("id")
        commentID := c.Param("comment_id")
        return c.String(http.StatusOK, 
            "Post: "+postID+", Comment: "+commentID)
    })

    // 6. 表单参数
    e.POST("/login", func(c echo.Context) error {
        username := c.FormValue("username")
        password := c.FormValue("password")
        return c.String(http.StatusOK, 
            "Username: "+username+", Password: "+password)
    })

    // 7. 文件上传
    e.POST("/upload", func(c echo.Context) error {
        file, err := c.FormFile("file")
        if err != nil {
            return err
        }
        return c.String(http.StatusOK, 
            "File: "+file.Filename+", Size: "+strconv.FormatInt(file.Size, 10))
    })

    // 8. JSON 响应
    e.GET("/json", func(c echo.Context) error {
        return c.JSON(http.StatusOK, map[string]interface{}{
            "message": "success",
            "data":    "some data",
        })
    })

    // 9. XML 响应
    e.GET("/xml", func(c echo.Context) error {
        return c.XML(http.StatusOK, map[string]string{
            "status": "ok",
        })
    })

    // 10. 重定向
    e.GET("/old", func(c echo.Context) error {
        return c.Redirect(http.StatusMovedPermanently, "/new")
    })

    // 11. HTML 响应
    e.GET("/html", func(c echo.Context) error {
        return c.HTML(http.StatusOK, "<h1>Hello, Echo!</h1>")
    })

    // 启动服务器
    e.Logger.Fatal(e.Start(":8080"))
}
```

### 3.2 中间件

```go
package main

import (
    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
)

func main() {
    e := echo.New()

    // 1. 内置中间件
    e.Use(middleware.Logger())           // 日志
    e.Use(middleware.Recover())          // 恢复 Panic
    e.Use(middleware.RequestID())        // 请求 ID
    e.Use(middleware.CORSWithConfig(
        middleware.CORSConfig{
            AllowOrigins: []string{"*"},
            AllowMethods: []string{echo.GET, echo.POST},
        },
    ))

    // 2. Gzip 压缩
    e.Use(middleware.GzipWithConfig(
        middleware.GzipConfig{
            Level: 5,
        },
    ))

    // 3. 速率限制
    e.Use(middleware.RateLimiterWithConfig(
        middleware.RateLimiterConfig{
            Skipper: middleware.DefaultSkipper,
            Store: middleware.NewRateLimiterMemoryStoreWithConfig(
                middleware.RateLimiterMemoryStoreConfig{
                    Rate:      10,             // 每秒 10 请求
                    Burst:     20,             // 突发 20
                    ExpiresIn: 3 * time.Minute,
                },
            ),
            ErrorFormatter: func(code int, message string, 
                info interface{}) string {
                return message
            },
            DenyHandler: func(context echo.Context, 
                request *http.Request, limit *middleware.RateLimitResult) error {
                return context.JSON(http.StatusTooManyRequests, 
                    map[string]string{"error": "rate limited"})
            },
        },
    ))

    // 4. 认证中间件
    e.Use(middleware.JWTWithConfig(
        middleware.JWTConfig{
            SigningKey: []byte("secret-key"),
            TokenLookup: "header:Authorization",
        },
    ))

    // 5. 自定义中间件
    e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            // 前置处理
            println("Before")
            
            err := next(c)
            
            // 后置处理
            println("After")
            
            return err
        }
    })

    // 6. 中间件组
    api := e.Group("/api")
    api.Use(middleware.Logger())
    api.GET("/users", func(c echo.Context) error {
        return c.String(http.StatusOK, "users")
    })

    e.Logger.Fatal(e.Start(":8080"))
}
```

### 3.3 参数验证

```go
package main

import (
    "net/http"
    "reflect"
    "time"

    "github.com/go-playground/validator/v10"
    "github.com/labstack/echo/v4"
)

// User 用户验证
type User struct {
    Name     string `validate:"required,min=2,max=50"`
    Email    string `validate:"required,email"`
    Age      int    `validate:"gte=0,lte=150"`
    Password string `validate:"required,min=6"`
}

func main() {
    e := echo.New()

    // 自定义验证器
    validate := validator.New()
    
    // 注册自定义验证器
    validate.RegisterValidation("phone", validatePhone)

    // 验证中间件
    e.POST("/users", func(c echo.Context) error {
        var user User
        
        // 绑定请求体
        if err := c.Bind(&user); err != nil {
            return c.JSON(http.StatusBadRequest, 
                map[string]string{"error": err.Error()})
        }

        // 验证
        if err := validate.Struct(user); err != nil {
            errors := err.(validator.ValidationErrors)
            return c.JSON(http.StatusBadRequest, 
                formatValidationErrors(errors))
        }

        return c.JSON(http.StatusOK, user)
    })

    e.Logger.Fatal(e.Start(":8080"))
}

func validatePhone(fl validator.FieldLevel) bool {
    phone := fl.Field().String()
    return len(phone) == 11
}

func formatValidationErrors(errors validator.ValidationErrors) map[string]string {
    result := make(map[string]string)
    for _, err := range errors {
        field := err.Field()
        tag := err.Tag()
        result[field] = formatErrorMessage(field, tag)
    }
    return result
}

func formatErrorMessage(field, tag string) string {
    messages := map[string]string{
        "required": "不能为空",
        "email":    "格式不正确",
        "min":      "太短",
        "max":      "太长",
        "gte":      "太小",
        "lte":      "太大",
    }
    if msg, ok := messages[tag]; ok {
        return field + " " + msg
    }
    return "格式不正确"
}
```

### 3.4 Echo 路由高级

```go
package main

import (
    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
)

func main() {
    e := echo.New()

    // 1. 路由前缀
    admin := e.Group("/admin")
    admin.GET("/dashboard", func(c echo.Context) error {
        return c.String(http.StatusOK, "Admin Dashboard")
    })

    // 2. 路由中间件
    e.GET("/protected", func(c echo.Context) error {
        return c.String(http.StatusOK, "Protected Resource")
    }, middleware.JWT([]byte("secret")))

    // 3. 路径参数验证
    e.GET("/user/:id", func(c echo.Context) error {
        id := c.Param("id")
        // 验证 ID 是否为数字
        return c.String(http.StatusOK, "User: "+id)
    })

    // 4. 正则路由
    e.GET("/users/:id", func(c echo.Context) error {
        id := c.Param("id")
        return c.String(http.StatusOK, "User: "+id)
    })

    // 5. 静态文件
    e.Static("/static", "public")

    // 6. 文件下载
    e.GET("/download", func(c echo.Context) error {
        return c.Attachment("files/demo.txt", "demo.txt")
    })

    // 7. 服务端发送事件 (SSE)
    e.GET("/sse", func(c echo.Context) error {
        c.Response().Header().Set("Content-Type", 
            "text/event-stream")
        c.Response().Header().Set("Cache-Control", "no-cache")
        c.Response().Header().Set("Connection", "keep-alive")

        for i := 0; i < 10; i++ {
            c.SSEvent("message", "data: "+strconv.Itoa(i))
            c.Response().Flush()
            time.Sleep(time.Second)
        }
        return nil
    })

    e.Logger.Fatal(e.Start(":8080"))
}
```

### 3.5 优雅关闭

```go
package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/labstack/echo/v4"
)

func main() {
    e := echo.New()

    e.GET("/", func(c echo.Context) error {
        return c.String(http.StatusOK, "Hello")
    })

    // 启动服务器
    go func() {
        if err := e.Start(":8080"); err != nil && 
            err != http.ErrServerClosed {
            log.Fatal("Server error:", err)
        }
    }()

    // 等待中断信号
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    log.Println("Shutting down server...")

    // 设置超时
    ctx, cancel := context.WithTimeout(context.Background(), 
        10*time.Second)
    defer cancel()

    // 关闭服务器
    if err := e.Shutdown(ctx); err != nil {
        log.Fatal("Server forced to shutdown:", err)
    }

    log.Println("Server exited")
}
```

---

## 4. 框架对比

### 4.1 性能对比

| 框架 | 请求/秒 | 延迟 (P99) | 内存占用 |
|------|---------|-----------|---------|
| Gin | ~50,000 | ~200μs | 低 |
| Echo | ~45,000 | ~220μs | 低 |
| Fiber | ~70,000 | ~150μs | 很低 |
| Chi | ~40,000 | ~250μs | 很低 |

### 4.2 功能对比

| 功能 | Gin | Echo | Fiber |
|------|-----|------|-------|
| 路由 | ✓ | ✓ | ✓ |
| 中间件 | ✓ | ✓ | ✓ |
| 参数绑定 | ✓ | ✓ | ✓ |
| 验证 | ✓ | ✓ | ✓ |
| HTML 渲染 | ✓ | ✓ | ✓ |
| 自动 TLS | ✗ | ✓ | ✗ |
| WebSocket | ✓ | ✓ | ✓ |
| SSE | ✓ | ✓ | ✓ |
| 速率限制 | ✓ | ✓ | ✓ |
| 压缩 | ✓ | ✓ | ✓ |

### 4.3 选择建议

**选择 Gin 如果：**
- 需要最好的生态支持
- 开发 API 服务
- 需要丰富的中间件
- 团队已有 Gin 使用经验

**选择 Echo 如果：**
- 需要极简 API
- 需要自动 HTTPS
- 需要内置验证
- 追求性能和简洁

**选择 Fiber 如果：**
- 从 Node.js 迁移
- 需要极低内存占用
- 需要 Express 风格的 API
- 追求最高性能

---

## 5. 最佳实践

### 5.1 项目结构

```
project/
├── cmd/
│   └── main.go
├── internal/
│   ├── config/
│   ├── handler/
│   ├── middleware/
│   ├── model/
│   └── service/
├── pkg/
│   ├── logger/
│   └── validator/
├── router/
│   └── router.go
├── config.yaml
├── go.mod
└── README.md
```

### 5.2 配置管理

```go
package config

import (
    "os"
    "github.com/spf13/viper"
)

type Config struct {
    Server struct {
        Host string `mapstructure:"host"`
        Port int    `mapstructure:"port"`
        Mode string `mapstructure:"mode"`
    } `mapstructure:"server"`
    Database struct {
        Host     string `mapstructure:"host"`
        Port     int    `mapstructure:"port"`
        Username string `mapstructure:"username"`
        Password string `mapstructure:"password"`
        Name     string `mapstructure:"name"`
    } `mapstructure:"database"`
}

func Load(path string) (*Config, error) {
    viper.SetConfigFile(path)
    viper.SetConfigType("yaml")

    viper.AutomaticEnv()

    if err := viper.ReadInConfig(); err != nil {
        return nil, err
    }

    var config Config
    if err := viper.Unmarshal(&config); err != nil {
        return nil, err
    }

    return &config, nil
}
```

### 5.3 日志管理

```go
package logger

import (
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

var log *zap.Logger

func Init(level string) error {
    config := zap.NewProductionConfig()
    if level == "debug" {
        config = zap.NewDevelopmentConfig()
    }
    
    config.EncoderConfig.TimeKey = "timestamp"
    config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

    var err error
    log, err = config.Build()
    return err
}

func Info(message string, fields ...zap.Field) {
    log.Info(message, fields...)
}

func Error(message string, fields ...zap.Field) {
    log.Error(message, fields...)
}

func Fatal(message string, fields ...zap.Field) {
    log.Fatal(message, fields...)
}
```

### 5.4 统一响应

```go
package response

import (
    "net/http"

    "github.com/labstack/echo/v4"
)

type Response struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
    Error   interface{} `json:"error,omitempty"`
}

func Success(c echo.Context, data interface{}) error {
    return c.JSON(http.StatusOK, Response{
        Code:    0,
        Message: "success",
        Data:    data,
    })
}

func Error(c echo.Context, code int, message string) error {
    return c.JSON(http.StatusOK, Response{
        Code:    code,
        Message: message,
    })
}

func ValidationError(c echo.Context, err error) error {
    return c.JSON(http.StatusBadRequest, Response{
        Code:    400,
        Message: "validation error",
        Error:   err,
    })
}
```
