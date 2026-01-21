// web/web_gin.go
// Gin Web 框架示例 - 详细注释版

package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// ====== Gin 框架基础 ======
/*
Gin 是 Go 语言中最流行的 Web 框架之一。

主要特点：
1. 高性能 - 基于 httprouter，速度极快
2. 中间件支持 - 强大的中间件机制
3. 路由分组 - 方便的路由组织
4. JSON 解析 - 自动绑定和验证
5. 错误管理 - 优雅的错误处理

安装：
  go get -u github.com/gin-gonic/gin
*/

// ====== 数据模型 ======

// User 用户模型
type User struct {
	ID       uint   `json:"id" binding:"required"`
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Age      int    `json:"age" binding:"gte=0,lte=150"`
}

// Post 帖子模型
type Post struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title" binding:"required,min=1,max=200"`
	Content   string    `json:"content" binding:"required"`
	AuthorID  uint      `json:"author_id"`
	CreatedAt time.Time `json:"created_at"`
}

// ====== 创建路由 ======

func setupRouter() *gin.Engine {
	// 1. 创建 Gin 路由器
	// gin.Default() 创建带有默认中间件的路由器
	// gin.New() 创建不带中间件的路由器
	router := gin.Default()

	// 2. 配置全局中间件
	// Logger 中间件：记录请求日志
	// Recovery 中间件：从 panic 中恢复
	router.Use(gin.Recovery())

	// 3. 健康检查路由
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// 4. 根路由
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome to Gin Web Framework!")
	})

	// 5. 路由分组 - API v1
	v1 := router.Group("/api/v1")
	{
		// 用户相关路由
		v1.POST("/users", createUser)
		v1.GET("/users", listUsers)
		v1.GET("/users/:id", getUser)
		v1.PUT("/users/:id", updateUser)
		v1.DELETE("/users/:id", deleteUser)

		// 帖子相关路由
		v1.POST("/posts", createPost)
		v1.GET("/posts", listPosts)
		v1.GET("/posts/:id", getPost)
	}

	// 6. 路由分组 - API v2
	v2 := router.Group("/api/v2")
	{
		v2.GET("/users", listUsersV2)
		v2.GET("/users/:id", getUserV2)
	}

	return router
}

// ====== 路由处理器 ======

// createUser 创建用户
// POST /api/v1/users
func createUser(c *gin.Context) {
	// 1. 绑定 JSON 数据到结构体
	// ShouldBind 自动验证 binding 标签
	// 如果验证失败，返回 400 错误
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// 2. 处理业务逻辑
	user.ID = 1 // 模拟数据库生成 ID
	user.CreatedAt = time.Now()

	// 3. 返回响应
	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user":    user,
	})
}

// listUsers 获取用户列表
// GET /api/v1/users
func listUsers(c *gin.Context) {
	// 1. 获取查询参数
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "10")

	// 2. 解析参数
	pageNum, _ := strconv.Atoi(page)
	limit, _ := strconv.Atoi(pageSize)

	// 3. 验证参数
	if pageNum < 1 {
		pageNum = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// 4. 返回模拟数据
	users := []User{
		{ID: 1, Username: "alice", Email: "alice@example.com", Age: 25},
		{ID: 2, Username: "bob", Email: "bob@example.com", Age: 30},
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      users,
		"page":      pageNum,
		"page_size": limit,
		"total":     2,
	})
}

// getUser 获取单个用户
// GET /api/v1/users/:id
func getUser(c *gin.Context) {
	// 1. 获取路径参数
	id := c.Param("id")

	// 2. 解析 ID
	userID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	// 3. 返回模拟数据
	if userID == 1 {
		c.JSON(http.StatusOK, User{
			ID:       1,
			Username: "alice",
			Email:    "alice@example.com",
			Age:      25,
		})
		return
	}

	// 4. 返回 404
	c.JSON(http.StatusNotFound, gin.H{
		"error": "User not found",
	})
}

// updateUser 更新用户
// PUT /api/v1/users/:id
func updateUser(c *gin.Context) {
	id := c.Param("id")

	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User updated successfully",
		"user":    user,
	})
}

// deleteUser 删除用户
// DELETE /api/v1/users/:id
func deleteUser(c *gin.Context) {
	id := c.Param("id")

	c.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully",
		"id":      id,
	})
}

// ====== 帖子相关路由 ======

// createPost 创建帖子
// POST /api/v1/posts
func createPost(c *gin.Context) {
	var post Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	post.ID = 1
	post.CreatedAt = time.Now()

	c.JSON(http.StatusCreated, gin.H{
		"message": "Post created successfully",
		"post":    post,
	})
}

// listPosts 获取帖子列表
// GET /api/v1/posts
func listPosts(c *gin.Context) {
	posts := []Post{
		{ID: 1, Title: "First Post", Content: "Hello World!", AuthorID: 1},
		{ID: 2, Title: "Second Post", Content: "Gin is great", AuthorID: 2},
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  posts,
		"total": 2,
	})
}

// getPost 获取单个帖子
// GET /api/v1/posts/:id
func getPost(c *gin.Context) {
	id := c.Param("id")

	c.JSON(http.StatusOK, Post{
		ID:        1,
		Title:     "First Post",
		Content:   "Hello World!",
		AuthorID:  1,
		CreatedAt: time.Now(),
	})
}

// ====== V2 版本路由 ======

// listUsersV2 获取用户列表 V2
// GET /api/v2/users
func listUsersV2(c *gin.Context) {
	// 使用新的响应格式
	c.JSON(http.StatusOK, gin.H{
		"version": "v2",
		"users": []User{
			{ID: 1, Username: "alice", Email: "alice@example.com", Age: 25},
		},
	})
}

// getUserV2 获取单个用户 V2
// GET /api/v2/users/:id
func getUserV2(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"version": "v2",
		"user": User{
			ID:       1,
			Username: "alice",
			Email:    "alice@example.com",
			Age:      25,
		},
	})
}

// ====== 中间件示例 ======

// LoggerMiddleware 日志中间件
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 请求开始时间
		start := time.Now()

		// 处理请求
		c.Next()

		// 请求处理完成后
		duration := time.Since(start)

		// 记录日志
		gin.DefaultWriter.Write([]byte(
			c.Request.Method + " " +
				c.Request.URL.Path + " " +
				c.Writer.Header().Get("Content-Type") + " " +
				strconv.Itoa(c.Writer.Status()) + " " +
				duration.String() + "\n",
		))
	}
}

// AuthMiddleware 认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 Header 获取 token
		token := c.GetHeader("Authorization")

		// 验证 token
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization token required",
			})
			return
		}

		// 验证通过，设置用户信息到上下文
		c.Set("user_id", 1)
		c.Next()
	}
}

// RateLimitMiddleware 限流中间件
func RateLimitMiddleware() gin.HandlerFunc {
	// 使用令牌桶算法实现限流
	// 这里简化为固定计数
	return func(c *gin.Context) {
		// 检查请求频率
		// 实际实现可以使用 golang.org/x/time/rate

		c.Next()
	}
}

// ====== 静态文件服务 ======

func staticFileHandler(router *gin.Engine) {
	// 提供静态文件服务
	// 第一个参数是 URL 路径前缀
	// 第二个参数是文件目录
	router.Static("/static", "./static")
	router.StaticFile("/favicon.ico", "./static/favicon.ico")
}

// HTML 渲染
func htmlHandler(router *gin.Engine) {
	// 加载 HTML 模板
	router.LoadHTMLGlob("templates/*")

	// 渲染 HTML
	router.GET("/page", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Gin Web Framework",
		})
	})
}

// ====== 文件上传 ======

func uploadHandler(router *gin.Engine) {
	router.POST("/upload", func(c *gin.Context) {
		// 1. 获取文件
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		// 2. 保存文件
		filename := "./uploads/" + file.Filename
		if err := c.SaveUploadedFile(file, filename); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		// 3. 返回响应
		c.JSON(http.StatusOK, gin.H{
			"message":  "File uploaded successfully",
			"filename": file.Filename,
			"size":     file.Size,
		})
	})
}

// ====== 重定向 ======

func redirectHandler(router *gin.Engine) {
	// 301 重定向
	router.GET("/old", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/new")
	})

	// 302 重定向
	router.GET("/new", func(c *gin.Context) {
		c.String(http.StatusOK, "This is the new page!")
	})
}

// ====== 自定义 404 处理器 ======

func customNotFoundHandler(router *gin.Engine) {
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"error":  "Page not found",
			"path":   c.Request.URL.Path,
			"method": c.Request.Method,
		})
	})
}

// ====== 主函数 ======

func main() {
	fmt.Println("=== Gin Web 框架示例 ===")

	// 1. 创建路由
	router := setupRouter()

	// 2. 配置静态文件
	staticFileHandler(router)

	// 3. 配置 HTML 渲染
	// htmlHandler(router)

	// 4. 配置文件上传
	uploadHandler(router)

	// 5. 配置重定向
	redirectHandler(router)

	// 6. 配置自定义 404
	customNotFoundHandler(router)

	// 7. 添加中间件到特定路由
	router.GET("/protected", AuthMiddleware(), func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		c.JSON(http.StatusOK, gin.H{
			"message": "Protected content",
			"user_id": userID,
		})
	})

	// 8. 启动服务器
	// gin.Run() 等同于 http.ListenAndServe(":8080", router)
	router.Run(":8080")
	// 或指定地址：router.Run(":3000")
}
