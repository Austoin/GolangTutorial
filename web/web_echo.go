// web/web_echo.go
// Echo Web 框架示例 - 详细注释版

package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

// ====== Echo 框架基础 ======
/*
Echo 是另一个流行的 Go Web 框架，以高性能和简洁的 API 著称。

主要特点：
1. 高性能 - 性能可与 Gin 媲美
2. 路由清晰 - API 设计直观
3. 中间件支持 - 强大的中间件机制
4. 自动 TLS - 轻松启用 HTTPS
5. JSON 支持 - 便捷的 JSON 处理

安装：
  go get -u github.com/labstack/echo/v4
*/

// ====== 数据模型 ======

// User 用户模型
type User struct {
	ID       uint   `json:"id" validate:"required"`
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Age      int    `json:"age" validate:"gte=0,lte=150"`
}

// Post 帖子模型
type Post struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title" validate:"required,min=1,max=200"`
	Content   string    `json:"content" validate:"required"`
	AuthorID  uint      `json:"author_id"`
	CreatedAt time.Time `json:"created_at"`
}

// ErrorResponse 错误响应
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// ====== 创建应用 ======

func createApp() *echo.Echo {
	// 1. 创建 Echo 实例
	// echo.New() 创建新的 Echo 实例
	e := echo.New()

	// 2. 配置实例
	// e.HideBanner = true // 隐藏启动横幅
	// e.HidePort = true   // 隐藏端口显示

	// 3. 添加全局中间件
	e.Use(LoggerMiddleware())
	e.Use(RecoveryMiddleware())

	// 4. 配置错误处理
	e.HTTPErrorHandler = customErrorHandler

	return e
}

// ====== 路由配置 ======

func setupRoutes(e *echo.Echo) {
	// 1. 健康检查
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// 2. 根路由
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Welcome to Echo Web Framework!")
	})

	// 3. API 路由组
	api := e.Group("/api/v1")

	// 用户路由
	api.POST("/users", createUserHandler)
	api.GET("/users", listUsersHandler)
	api.GET("/users/:id", getUserHandler)
	api.PUT("/users/:id", updateUserHandler)
	api.DELETE("/users/:id", deleteUserHandler)

	// 帖子路由
	api.POST("/posts", createPostHandler)
	api.GET("/posts", listPostsHandler)
	api.GET("/posts/:id", getPostHandler)

	// 4. V2 API 路由组
	apiV2 := e.Group("/api/v2")
	apiV2.GET("/users", listUsersV2Handler)
	apiV2.GET("/users/:id", getUserV2Handler)
}

// ====== 路由处理器 ======

// createUserHandler 创建用户
// POST /api/v1/users
func createUserHandler(c echo.Context) error {
	// 1. 绑定请求体到结构体
	// Bind 方法自动解析 JSON、Form 等格式
	var user User
	if err := c.Bind(&user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	// 2. 验证数据（使用自定义验证）
	if err := validateUser(&user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// 3. 处理业务逻辑
	user.ID = 1
	user.CreatedAt = time.Now()

	// 4. 返回响应
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "User created successfully",
		"user":    user,
	})
}

// listUsersHandler 获取用户列表
// GET /api/v1/users
func listUsersHandler(c echo.Context) error {
	// 1. 获取查询参数
	page := c.QueryParam("page")
	pageSize := c.QueryParam("page_size")

	// 2. 解析参数（带默认值）
	pageNum := 1
	if p, err := strconv.Atoi(page); err == nil && p > 0 {
		pageNum = p
	}

	limit := 10
	if ps, err := strconv.Atoi(pageSize); err == nil && ps > 0 && ps <= 100 {
		limit = ps
	}

	// 3. 返回模拟数据
	users := []User{
		{ID: 1, Username: "alice", Email: "alice@example.com", Age: 25},
		{ID: 2, Username: "bob", Email: "bob@example.com", Age: 30},
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data":      users,
		"page":      pageNum,
		"page_size": limit,
		"total":     2,
	})
}

// getUserHandler 获取单个用户
// GET /api/v1/users/:id
func getUserHandler(c echo.Context) error {
	// 1. 获取路径参数
	id := c.Param("id")

	// 2. 解析 ID
	userID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID")
	}

	// 3. 返回模拟数据
	if userID == 1 {
		return c.JSON(http.StatusOK, User{
			ID:       1,
			Username: "alice",
			Email:    "alice@example.com",
			Age:      25,
		})
	}

	// 4. 返回 404
	return echo.NewHTTPError(http.StatusNotFound, "User not found")
}

// updateUserHandler 更新用户
// PUT /api/v1/users/:id
func updateUserHandler(c echo.Context) error {
	id := c.Param("id")

	var user User
	if err := c.Bind(&user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "User updated successfully",
		"user":    user,
	})
}

// deleteUserHandler 删除用户
// DELETE /api/v1/users/:id
func deleteUserHandler(c echo.Context) error {
	id := c.Param("id")

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "User deleted successfully",
		"id":      id,
	})
}

// ====== 帖子路由处理器 ======

// createPostHandler 创建帖子
// POST /api/v1/posts
func createPostHandler(c echo.Context) error {
	var post Post
	if err := c.Bind(&post); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	post.ID = 1
	post.CreatedAt = time.Now()

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Post created successfully",
		"post":    post,
	})
}

// listPostsHandler 获取帖子列表
// GET /api/v1/posts
func listPostsHandler(c echo.Context) error {
	posts := []Post{
		{ID: 1, Title: "First Post", Content: "Hello World!", AuthorID: 1},
		{ID: 2, Title: "Second Post", Content: "Echo is great", AuthorID: 2},
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data":  posts,
		"total": 2,
	})
}

// getPostHandler 获取单个帖子
// GET /api/v1/posts/:id
func getPostHandler(c echo.Context) error {
	id := c.Param("id")

	return c.JSON(http.StatusOK, Post{
		ID:        1,
		Title:     "First Post",
		Content:   "Hello World!",
		AuthorID:  1,
		CreatedAt: time.Now(),
	})
}

// ====== V2 路由处理器 ======

// listUsersV2Handler 获取用户列表 V2
// GET /api/v2/users
func listUsersV2Handler(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"version": "v2",
		"users": []User{
			{ID: 1, Username: "alice", Email: "alice@example.com", Age: 25},
		},
	})
}

// getUserV2Handler 获取单个用户 V2
// GET /api/v2/users/:id
func getUserV2Handler(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"version": "v2",
		"user": User{
			ID:       1,
			Username: "alice",
			Email:    "alice@example.com",
			Age:      25,
		},
	})
}

// ====== 中间件 ======

// LoggerMiddleware 日志中间件
func LoggerMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// 请求开始时间
			start := time.Now()

			// 处理请求
			err := next(c)

			// 请求处理完成后
			duration := time.Since(start)

			// 打印日志
			println(c.Request().Method, c.Request().URL.Path,
				c.Response().Status, duration.String())

			return err
		}
	}
}

// RecoveryMiddleware 恢复中间件
func RecoveryMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			defer func() {
				if r := recover(); r != nil {
					err := echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
					c.Error(err)
				}
			}()
			return next(c)
		}
	}
}

// AuthMiddleware 认证中间件
func AuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := c.Request().Header.Get("Authorization")

			if token == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Authorization token required")
			}

			// 验证 token
			c.Set("user_id", 1)

			return next(c)
		}
	}
}

// CORSMiddleware 跨域中间件
func CORSMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// 设置响应头
			c.Response().Header().Set("Access-Control-Allow-Origin", "*")
			c.Response().Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
			c.Response().Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			// 处理 OPTIONS 预检请求
			if c.Request().Method == http.MethodOptions {
				return c.NoContent(http.StatusOK)
			}

			return next(c)
		}
	}
}

// RateLimitMiddleware 限流中间件
func RateLimitMiddleware() echo.MiddlewareFunc {
	// 实际实现可以使用 golang.org/x/time/rate
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// 简化的限流逻辑
			// 实际应该使用令牌桶或漏桶算法

			return next(c)
		}
	}
}

// ====== 自定义错误处理 ======

// customErrorHandler 自定义错误处理器
func customErrorHandler(err error, c echo.Context) {
	// 1. 检查是否是 HTTP 错误
	httpErr, ok := err.(*echo.HTTPError)
	if ok {
		c.JSON(httpErr.Code, ErrorResponse{
			Error:   httpErr.Message.(string),
			Message: "An error occurred",
		})
		return
	}

	// 2. 其他错误
	c.JSON(http.StatusInternalServerError, ErrorResponse{
		Error:   "Internal server error",
		Message: err.Error(),
	})
}

// ====== 数据验证 ======

// validateUser 验证用户数据
func validateUser(user *User) error {
	if user.Username == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Username is required")
	}
	if len(user.Username) < 3 {
		return echo.NewHTTPError(http.StatusBadRequest, "Username must be at least 3 characters")
	}
	if user.Email == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Email is required")
	}
	if user.Age < 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Age must be non-negative")
	}
	return nil
}

// ====== 静态文件服务 ======

func staticFileHandler(e *echo.Echo) {
	// 提供静态文件服务
	e.Static("/static", "./static")
	e.File("/favicon.ico", "./static/favicon.ico")
}

// ====== 文件上传 ======

func uploadHandler(e *echo.Echo) {
	e.POST("/upload", func(c echo.Context) error {
		// 1. 获取文件
		file, err := c.FormFile("file")
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		// 2. 打开文件
		src, err := file.Open()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		defer src.Close()

		// 3. 保存文件（示例）
		// dst, err := os.Create("./uploads/" + file.Filename)
		// if err != nil {
		//     return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		// }
		// defer dst.Close()
		// io.Copy(dst, src)

		return c.JSON(http.StatusOK, map[string]interface{}{
			"message":  "File uploaded successfully",
			"filename": file.Filename,
			"size":     file.Size,
		})
	})
}

// ====== 重定向 ======

func redirectHandler(e *echo.Echo) {
	e.GET("/old", func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, "/new")
	})

	e.GET("/new", func(c echo.Context) error {
		return c.String(http.StatusOK, "This is the new page!")
	})
}

// ====== 自定义 404 处理器 ======

func customNotFoundHandler(e *echo.Echo) {
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		if httpErr, ok := err.(*echo.HTTPError); ok && httpErr.Code == http.StatusNotFound {
			c.JSON(http.StatusNotFound, map[string]interface{}{
				"error":  "Page not found",
				"path":   c.Request().URL.Path,
				"method": c.Request().Method,
			})
			return
		}

		// 调用默认错误处理
		echo.DefaultHTTPErrorHandler(err, c)
	}
}

// ====== 主函数 ======

func main() {
	fmt.Println("=== Echo Web 框架示例 ===")

	// 1. 创建应用
	e := createApp()

	// 2. 配置路由
	setupRoutes(e)

	// 3. 配置静态文件
	staticFileHandler(e)

	// 4. 配置文件上传
	uploadHandler(e)

	// 5. 配置重定向
	redirectHandler(e)

	// 6. 配置自定义 404
	customNotFoundHandler(e)

	// 7. 添加中间件到特定路由
	e.GET("/protected", AuthMiddleware(), func(c echo.Context) error {
		userID := c.Get("user_id")
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "Protected content",
			"user_id": userID,
		})
	})

	// 8. 启动服务器
	// e.Start() 启动服务器
	// 使用 StartTLS 可以启用 TLS（HTTPS）
	e.Logger.Fatal(e.Start(":8080"))
}
