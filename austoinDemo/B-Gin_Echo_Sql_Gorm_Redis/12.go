package myafter

import (
	"fmt"
	"net/http"
	// "os/user"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// Gin Web 框架
func Lesson12() {
	// 设置 gin 模式
	gin.SetMode(gin.DebugMode)

	// 创建路由器
	router := gin.Default()

	// 基础路由
	basicDemo(router)

	// 参数获取
	paramDemo(router)

	// Json 绑定
	jsonDemo(router)

	// 路由分组
	groupDemo(router)

	// 中间件
	middlewareDemo(router)

	// 404处理（客户端）
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"error":  "页面未找到",
			"path":   c.Request.URL.Path,
			"method": c.Request.Method,
		})
	})

	fmt.Println("Gin 服务器启动")
	fmt.Println("访问 http://localhost:8080")
	fmt.Println("按 Ctrl+C 停止服务器")

	// 启动服务器
	router.Run(":8080")
}

// 基础路由
func basicDemo(router *gin.Engine) { // gin.Engine 路由引擎
	// GET - 查询数据
	router.GET("/", func(c *gin.Context) { // gin.Context 请求上下文
		c.String(http.StatusOK, "Welcome to Gin Web Framework!")
	})

	router.GET("/hello", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, Gin!")
	})

	// POST - 创建数据
	router.POST("/submit", func(c *gin.Context) {
		c.String(http.StatusOK, "数据提交成功")
	})

	// PUT - 更新数据
	router.PUT("/update", func(c *gin.Context) {
		c.String(http.StatusOK, "数据更新成功")
	})

	// DELETE - 删除数据
	router.DELETE("/delete", func(c *gin.Context) {
		c.String(http.StatusOK, "数据删除成功")
	})
}

// 参数获取
func paramDemo(router *gin.Engine) {
	// 路径参数: id
	router.GET("/users/:id", func(c *gin.Context) {
		id := c.Param("id") // 获取单个路径参数的方法
		c.String(http.StatusOK, "用户ID: %s", id)
	})

	// 多个路径参数
	router.GET("/posts/:postID/comments/:commentID", func(c *gin.Context) {
		postID := c.Param("postID")
		commentID := c.Param("commentID")
		c.String(http.StatusOK, "帖子%s的评论%s", postID, commentID)
	})

	// 通配符 *
	router.GET("/files/*filepath", func(c *gin.Context) {
		filepath := c.Param("filepath")
		c.String(http.StatusOK, "文件路径: %s", filepath)
	})

	// 查询参数 ?page=1&size=10
	router.GET("/search", func(c *gin.Context) {
		// DefaultQuery 获取，不存在返回默认值
		page := c.DefaultQuery("page", "1")
		size := c.DefaultQuery("size", "10")
		keyword := c.Query("keyword")
		// 示例：访问/seach?keyword=go，返回搜索: go, 页码: 1, 每页: 10

		c.String(http.StatusOK, "搜索: %s, 页码: %s, 每页: %s", keyword, page, size)
	})

	// 获取所有擦寻参数
	router.GET("/filter", func(c *gin.Context) {
		keys := c.QueryMap("filter")
		c.String(http.StatusOK, "筛选参数: %v", keys)
	})

	// Header参数
	router.GET("/header", func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		contentType := c.GetHeader("Content-Type")

		if auth == "" {
			auth = "无 token"
		}

		c.JSON(http.StatusOK, gin.H{
			"Authorization": auth,
			"Content-Type":  contentType,
		})
	})

	router.POST("/form", func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.DefaultPostForm("password", "123456")

		c.JSON(http.StatusOK, gin.H{
			"username": username,
			"password": password,
		})
	})
	// 示例：提交表单username=admin&password=888888;
	// 返回{"username":"admin","password":"888888"}；
	// 仅提交username=test，返回{"username":"test","password":"123456"}
}

// // Json 绑定实现部分

// 用户结构体
type User struct {
	ID       uint   `json:"id" binding:"required"`
	Username string `json:"username" binding:"required,min=3,max=20"`
	Email    string `json:"email" binding:"required,email"`
	Age      int    `json:"age" binding:"gte=0,lte=150"`
}

// 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

func jsonDemo(router *gin.Engine) {
	// 创建用户 - POST /users
	router.POST("/users", func(c *gin.Context) {
		var user User

		// ShouldBindJSON 解析JSON并验证
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   err.Error(),
				"message": "参数验证失败",
			})
			return
		}
		// 模拟保存到数据库
		user.ID = 1

		c.JSON(http.StatusCreated, gin.H{
			"message": "用户创建成功",
			"data":    user,
		})
	})

	// 登录 - POST /login
	router.POST("/login", func(c *gin.Context) {
		var req LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		// 验证用户名密码（模拟）
		if req.Username == "admin" && req.Password == "123456" {
			c.JSON(http.StatusOK, gin.H{
				"message": "登录成功",
				"token":   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "用户名或密码错误",
			})
		}
	})

	// 返回JSON - GET /json
	router.GET("/json", func(c *gin.Context) {
		user := User{
			ID:       1,
			Username: "austoin",
			Email:    "austoin@example.com",
			Age:      25,
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "success",
			"data":    user,
		})
	})

	// 获取用户列表 - GET /users
	router.GET("/users", func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

		users := []User{
			{ID: 1, Username: "alice", Email: "alice@example.com", Age: 25},
			{ID: 2, Username: "bob", Email: "bob@example.com", Age: 30},
		}

		c.JSON(http.StatusOK, gin.H{
			"data":  users,
			"page":  page,
			"size":  size,
			"total": 2,
		})
	})
}

// 路由分组
func groupDemo(router *gin.Engine) {
	// API v1 分组
	v1 := router.Group("/api/v1")
	{
		v1.GET("users", func(c *gin.Context) {
			c.String(http.StatusOK, "用户列表 - API v1")
		})

		v1.GET("/users/:id", func(c *gin.Context) {
			c.String(http.StatusOK, "用户详情 - API v1, ID: "+c.Param("id"))
		})

		v1.POST("/users", func(c *gin.Context) {
			c.String(http.StatusOK, "创建用户 - API v1")
		})

		v1.PUT("/users/:id", func(c *gin.Context) {
			c.String(http.StatusOK, "更新用户 - API v1, ID: "+c.Param("id"))
		})

		v1.DELETE("/users/:id", func(c *gin.Context) {
			c.String(http.StatusOK, "删除用户 - API v1, ID: "+c.Param("id"))
		})
	}

	// API v2 分组（新版本）
	v2 := router.Group("/api/v2")
	{
		v2.GET("/users", func(c *gin.Context) {
			c.String(http.StatusOK, "用户列表 - API v2 (新版)")
		})
		v2.GET("/users/:id", func(c *gin.Context) {
			c.String(http.StatusOK, "用户详情 - API v2, ID: "+c.Param("id"))
		})
	}

	// Admin 分组（带前缀）
	admin := router.Group("/admin")
	{
		admin.GET("/dashboard", func(c *gin.Context) {
			c.String(http.StatusOK, "管理后台首页")
		})
		admin.GET("/settings", func(c *gin.Context) {
			c.String(http.StatusOK, "系统设置")
		})
	}
}

// 中间件
// Logger 日志中间件：记录请求的时间、方法、路径、状态码、耗时
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		method := c.Request.Method
		path := c.Request.URL.Path
		// 放行请求，执行后续处理函数
		c.Next()
		// 请求结束后记录日志
		latency := time.Since(start)
		status := c.Writer.Status()
		fmt.Printf("[%s] %s %s | %d | %v\n",
			time.Now().Format("2006-01-02 15:04:05"),
			method, path, status, latency)
	}
}

// AuthMiddleware 认证中间件：校验Token，保护敏感路由
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "需要认证，请提供Authorization header",
			})
			return
		}
		// 简化版Token验证（实际项目需替换为JWT/数据库校验）
		if token != "Bearer mytoken123" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "无效的token",
			})
			return
		}
		// 验证通过，设置用户信息到上下文
		c.Set("user_id", 1)
		c.Set("username", "admin")
		c.Next()
	}
}

// CorsMiddleware 跨域中间件：解决前端跨域请求问题
func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 设置跨域响应头
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		// 处理OPTIONS预检请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

func middlewareDemo(router *gin.Engine) {
	// 全局中间件：所有路由都会执行
	router.Use(Logger())
	router.Use(CorsMiddleware())

	// 公开路由：无需认证
	router.GET("/public", func(c *gin.Context) {
		c.String(http.StatusOK, "这是公开页面，无需认证")
	})
	router.GET("/public/json", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "公开数据",
			"data":    []string{"item1", "item2", "item3"},
		})
	})

	// 受保护的路由分组：需要认证
	protected := router.Group("/protected")
	protected.Use(AuthMiddleware()) // 分组内路由都需经过认证中间件
	{
		protected.GET("/profile", func(c *gin.Context) {
			// 正确接收c.Get()的两个返回值（忽略exists）
			userId, _ := c.Get("user_id")
			username, _ := c.Get("username")

			c.JSON(http.StatusOK, gin.H{
				"message":  "受保护的-profile",
				"user_id":  userId,
				"username": username,
			})
		})

		protected.GET("/settings", func(c *gin.Context) {
			// 接收两个返回值，并判断值是否存在
			userId, exists := c.Get("user_id")
			if !exists {
				userId = "未知用户"
			}
			c.String(http.StatusOK, "受保护的-settings,用户ID: %v", userId)
		})

		protected.GET("/dashboard", func(c *gin.Context) {
			// 先获取值，再使用（修复编译错误）
			userId, _ := c.Get("user_id")
			c.JSON(http.StatusOK, gin.H{
				"message":    "管理后台首页",
				"user_id":    userId,
				"role":       "admin",
				"last_login": time.Now().Format("2006-01-02 15:04:05"),
			})
		})
	}
}

// gin.Engine是一个复杂的结构体（包含路由规则表、中间件列表、配置项等）；
// 传指针可以避免拷贝整个结构体（性能优化）
// 指针传递能保证在basicDemo中对router的修改（注册路由）；

//router.GET(path, handler)：路由注册语法
//第一个参数"/"：静态路由路径（根路径），表示 “访问域名 / IP + 端口的根地址（如http://localhost:8080/）”；
//第二个参数：匿名处理函数（func(c *gin.Context)），是请求的 “业务处理逻辑”，当匹配到GET /时执行。

//c.String(http.StatusOK, ...)：响应返回
//c是*gin.Context类型的参数（请求上下文），是处理请求 / 响应的 “唯一载体”；
//c.String()：Gin 提供的快捷方法，返回 “纯文本格式” 的响应；
//http.StatusOK：HTTP 标准状态码（值为 200），表示 “请求成功处理”，替代硬编码200更规范；
//第二个参数：要返回的字符串内容，客户端会收到这个文本。
