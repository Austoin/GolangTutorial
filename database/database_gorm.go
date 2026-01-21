// database/database_gorm.go
// GORM 数据库操作 - 详细注释版

package main

import (
	"fmt"
	"log"
	"time"

	// 导入 GORM 和数据库驱动
	// GORM 是 Go 语言中最流行的 ORM 框架
	// 官方网站：https://gorm.io
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// ====== 数据模型定义 ======

// User 用户模型
// GORM 使用结构体标签来定义映射规则
type User struct {
	// gorm:"primaryKey" 标记为主键
	ID uint `gorm:"primaryKey"`

	// gorm:"size:50;uniqueIndex" 定义字段大小和唯一索引
	Username string `gorm:"size:50;uniqueIndex"`
	Email    string `gorm:"size:100;uniqueIndex"`

	// gorm:"-" 表示忽略此字段，不映射到数据库
	Password string `gorm:"-"` // 不存储密码明文

	// gorm:"column:password_hash" 自定义列名
	PasswordHash string `gorm:"column:password_hash"`

	// gorm:"autoCreateTime" 自动设置创建时间为当前时间
	CreatedAt time.Time `gorm:"autoCreateTime"`

	// gorm:"autoUpdateTime" 自动设置更新时间为当前时间
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

	// gorm:"-" 忽略此字段
	Age int `gorm:"-"` // 不存储年龄，只在内存中使用

	// 关联关系
	// gorm:"foreignKey:UserID" 定义外键
	Posts []Post `gorm:"foreignKey:UserID"` // 一对多关系
}

// Post 帖子模型
type Post struct {
	ID      uint   `gorm:"primaryKey"`
	Title   string `gorm:"size:255"`
	Content string `gorm:"type:text"`

	// 外键
	UserID uint `gorm:"index"` // 自动创建索引

	// 关联关系
	// gorm:"references:ID" 指定引用的列
	User User `gorm:"references:ID"` // 属于 User

	// 软删除
	DeletedAt gorm.DeletedAt `gorm:"index"` // 软删除支持
}

// Comment 评论模型
type Comment struct {
	ID      uint   `gorm:"primaryKey"`
	Content string `gorm:"type:text"`

	// 多对一关系
	UserID uint `gorm:"index"`
	User   User `gorm:"references:ID"`

	PostID uint `gorm:"index"`
	Post   Post `gorm:"references:ID"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
}

// ====== 数据库连接 ======

// Database GORM 数据库封装
type Database struct {
	db *gorm.DB // GORM DB 实例
}

// NewDatabase 创建数据库连接
func NewDatabase(dsn string) (*Database, error) {
	// 1. 配置 GORM
	// Config 结构体包含各种配置选项
	config := &gorm.Config{
		// NamingStrategy 命名策略
		// 用于自动生成表名、列名等
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "t_",  // 表前缀
			SingularTable: false, // 不使用单数表名
		},

		// Logger 日志配置
		// NewLogger 创建自定义日志
		// LogLevel 日志级别
		Logger: logger.Default.LogMode(logger.Info),

		// SkipDefaultTransaction 跳过默认事务
		// 对于大量写入操作，关闭事务可能提高性能
		SkipDefaultTransaction: false,

		// PrepareStmt 预编译 SQL
		// 在执行前预编译，提高性能
		PrepareStmt: true,
	}

	// 2. 打开数据库连接
	// gorm.Open 接受 Dialector 和 Config
	// Dialector 是数据库驱动的抽象
	db, err := gorm.Open(mysql.Open(dsn), config)
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %w", err)
	}

	// 3. 配置连接池（通过 *gorm.DB 访问底层 sql.DB）
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("获取底层连接失败: %w", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	log.Println("GORM 数据库连接成功")

	return &Database{db: db}, nil
}

// Close 关闭数据库连接
func (d *Database) Close() error {
	sqlDB, err := d.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// DB 获取数据库实例
func (d *Database) DB() *gorm.DB {
	return d.db
}

// ====== 自动迁移 ======

// AutoMigrate 自动迁移数据库表结构
// 根据模型定义创建或更新表
func (d *Database) AutoMigrate() error {
	// AutoMigrate 根据结构体变化自动更新表结构
	// 只会添加新列，不会删除或修改已有列
	log.Println("开始自动迁移...")

	err := d.db.AutoMigrate(
		&User{},    // 迁移 User 表
		&Post{},    // 迁移 Post 表
		&Comment{}, // 迁移 Comment 表
	)

	if err != nil {
		return fmt.Errorf("自动迁移失败: %w", err)
	}

	log.Println("自动迁移完成")
	return nil
}

// ====== 创建操作 ======

// CreateUser 创建用户
func (d *Database) CreateUser(user *User) error {
	// 1. 创建记录
	// Create 返回 *gorm.DB，可以链式调用
	result := d.db.Create(user)

	// 2. 检查错误
	if result.Error != nil {
		return fmt.Errorf("创建用户失败: %w", result.Error)
	}

	// 3. 获取插入的 ID
	// user.ID 会被自动填充
	log.Printf("创建用户成功，ID: %d", user.ID)

	return nil
}

// CreateUsers 批量创建用户
func (d *Database) CreateUsers(users []User) error {
	// 1. 批量创建
	// CreateInBatches 分批创建，避免内存溢出
	result := d.db.CreateInBatches(users, 100)

	if result.Error != nil {
		return fmt.Errorf("批量创建用户失败: %w", result.Error)
	}

	log.Printf("批量创建用户成功，数量: %d", result.RowsAffected)
	return nil
}

// CreateUserWithPosts 创建用户及帖子
func (d *Database) CreateUserWithPosts(user *User, posts []Post) error {
	// 1. 创建用户
	if err := d.db.Create(user).Error; err != nil {
		return err
	}

	// 2. 设置外键
	for i := range posts {
		posts[i].UserID = user.ID
	}

	// 3. 创建帖子
	if err := d.db.Create(&posts).Error; err != nil {
		return err
	}

	return nil
}

// ====== 查询操作 ======

// GetUserByID 根据 ID 查询用户
func (d *Database) GetUserByID(id uint) (*User, error) {
	var user User

	// 1. 查询单条记录
	// First 找到第一条匹配的记录
	// 如果找不到，返回 gorm.ErrRecordNotFound 错误
	result := d.db.First(&user, id)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // 没找到
		}
		return nil, result.Error
	}

	return &user, nil
}

// GetUserByUsername 根据用户名查询用户
func (d *Database) GetUserByUsername(username string) (*User, error) {
	var user User

	// 使用 Where 条件查询
	result := d.db.Where("username = ?", username).First(&user)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}

	return &user, nil
}

// GetAllUsers 查询所有用户
func (d *Database) GetAllUsers() ([]User, error) {
	var users []User

	// 查询所有记录
	result := d.db.Find(&users)

	if result.Error != nil {
		return nil, result.Error
	}

	return users, nil
}

// GetUsersByCondition 条件查询
func (d *Database) GetUsersByCondition(conditions map[string]interface{}) ([]User, error) {
	var users []User

	// 使用 Where 条件
	// 支持链式调用
	result := d.db.Where(conditions).Find(&users)

	if result.Error != nil {
		return nil, result.Error
	}

	return users, nil
}

// GetUsersByEmailPrefix 按邮箱前缀查询
func (d *Database) GetUsersByEmailPrefix(prefix string) ([]User, error) {
	var users []User

	// 使用 LIKE 进行模糊查询
	result := d.db.Where("email LIKE ?", prefix+"%").Find(&users)

	if result.Error != nil {
		return nil, result.Error
	}

	return users, nil
}

// GetUserWithPosts 获取用户及其帖子
func (d *Database) GetUserWithPosts(id uint) (*User, error) {
	var user User

	// Preload 预加载关联数据
	// 这样可以一次性获取用户及其所有帖子
	result := d.db.Preload("Posts").First(&user, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

// GetUserPosts 获取用户的帖子
func (d *Database) GetUserPosts(userID uint) ([]Post, error) {
	var posts []Post

	result := d.db.Where("user_id = ?", userID).Find(&posts)

	if result.Error != nil {
		return nil, result.Error
	}

	return posts, nil
}

// GetUserPostsWithComments 获取用户的帖子及其评论
func (d *Database) GetUserPostsWithComments(userID uint) ([]Post, error) {
	var posts []Post

	// 预加载多层关联
	result := d.db.Preload("Comments").Preload("Comments.User").
		Where("user_id = ?", userID).Find(&posts)

	if result.Error != nil {
		return nil, result.Error
	}

	return posts, nil
}

// CountUsers 统计用户数量
func (d *Database) CountUsers() (int64, error) {
	var count int64

	result := d.db.Model(&User{}).Count(&count)
	return count, result.Error
}

// ====== 更新操作 ======

// UpdateUser 更新用户
func (d *Database) UpdateUser(user *User) error {
	// 1. 保存更新
	// Save 会更新所有字段
	result := d.db.Save(user)

	if result.Error != nil {
		return fmt.Errorf("更新用户失败: %w", result.Error)
	}

	return nil
}

// UpdateUserField 更新用户单个字段
func (d *Database) UpdateUserField(id uint, field string, value interface{}) error {
	result := d.db.Model(&User{}).Where("id = ?", id).Update(field, value)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("用户不存在: %d", id)
	}

	return nil
}

// UpdateUserEmail 更新用户邮箱
func (d *Database) UpdateUserEmail(id uint, email string) error {
	result := d.db.Model(&User{}).Where("id = ?", id).Update("email", email)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("用户不存在: %d", id)
	}

	return nil
}

// UpdateUsersByCondition 批量更新
func (d *Database) UpdateUsersByCondition(condition map[string]interface{}, updates map[string]interface{}) error {
	result := d.db.Model(&User{}).Where(condition).Updates(updates)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

// ====== 删除操作 ======

// DeleteUser 删除用户（软删除）
func (d *Database) DeleteUser(id uint) error {
	// 如果模型包含 DeletedAt 字段，Delete 默认执行软删除
	result := d.db.Delete(&User{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("用户不存在: %d", id)
	}

	return nil
}

// DeleteUserPermanently 永久删除用户
func (d *Database) DeleteUserPermanently(id uint) error {
	// Unscoped 忽略软删除字段
	result := d.db.Unscoped().Delete(&User{}, id)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

// DeleteUsersByCondition 批量删除
func (d *Database) DeleteUsersByCondition(condition map[string]interface{}) error {
	result := d.db.Where(condition).Delete(&User{})

	if result.Error != nil {
		return result.Error
	}

	return nil
}

// ====== 原生 SQL ======

// QueryRaw 原生查询
func (d *Database) QueryRaw(query string, args ...interface{}) ([]map[string]interface{}, error) {
	// Row 返回一行数据
	// Rows 返回多行数据
	rows, err := d.db.Raw(query, args...).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// 获取列名
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	// 遍历行
	var results []map[string]interface{}
	for rows.Next() {
		// 创建切片来存储值
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))

		for i := range values {
			valuePtrs[i] = &values[i]
		}

		// 扫描行
		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		// 构建结果映射
		rowMap := make(map[string]interface{})
		for i, col := range columns {
			rowMap[col] = values[i]
		}

		results = append(results, rowMap)
	}

	return results, nil
}

// ExecRaw 原生执行
func (d *Database) ExecRaw(query string, args ...interface{}) error {
	return d.db.Exec(query, args...).Error
}

// ====== 事务 ======

// TransferMoney 转账示例
func (d *Database) TransferMoney(fromID, toID uint, amount float64) error {
	// 使用事务
	return d.db.Transaction(func(tx *gorm.DB) error {
		// 扣款
		result := tx.Model(&User{}).Where("id = ? AND balance >= ?", fromID, amount).
			Update("balance", gorm.Expr("balance - ?", amount))
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return fmt.Errorf("余额不足")
		}

		// 加款
		if err := tx.Model(&User{}).Where("id = ?", toID).
			Update("balance", gorm.Expr("balance + ?", amount)).Error; err != nil {
			return err
		}

		return nil
	})
}

// ====== 钩子函数 ======

/*
GORM 提供了丰富的钩子函数，可以在 CRUD 操作前后执行自定义逻辑：

创建前：
  - BeforeCreate(db *gorm.DB) error

创建后：
  - AfterCreate(tx *gorm.DB) error

更新前：
  - BeforeUpdate(db *gorm.DB) error

更新后：
  - AfterUpdate(tx *gorm.DB) error

删除前：
  - BeforeDelete(db *gorm.DB) error

删除后：
  - AfterDelete(tx *gorm.DB) error

示例：
func (u *User) BeforeCreate(tx *gorm.DB) error {
    // 加密密码
    u.PasswordHash = hashPassword(u.Password)
    return nil
}
*/

// ====== 主函数 ======

func main() {
	fmt.Println("=== GORM 数据库操作示例 ===")

	// 1. 连接数据库
	dsn := "root:password@tcp(localhost:3306)/testdb?charset=utf8mb4&parseTime=True"

	db, err := NewDatabase(dsn)
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}
	defer db.Close()

	// 2. 自动迁移
	if err := db.AutoMigrate(); err != nil {
		log.Fatalf("迁移失败: %v", err)
	}

	// 3. 创建测试数据
	users := []User{
		{Username: "alice", Email: "alice@example.com"},
		{Username: "bob", Email: "bob@example.com"},
		{Username: "charlie", Email: "charlie@example.com"},
	}

	if err := db.CreateUsers(users); err != nil {
		log.Printf("创建用户失败: %v", err)
	}

	// 4. 查询测试
	user, _ := db.GetUserByUsername("alice")
	if user != nil {
		fmt.Printf("查询到用户: %s (%s)\n", user.Username, user.Email)
	}

	// 5. 预加载测试
	userWithPosts, _ := db.GetUserWithPosts(user.ID)
	if userWithPosts != nil {
		fmt.Printf("用户 %s 有 %d 篇帖子\n", userWithPosts.Username, len(userWithPosts.Posts))
	}

	// 6. 更新测试
	if user != nil {
		user.Email = "alice.new@example.com"
		if err := db.UpdateUser(user); err != nil {
			log.Printf("更新失败: %v", err)
		}
	}

	// 7. 删除测试
	if err := db.DeleteUser(3); err != nil {
		log.Printf("删除失败: %v", err)
	}

	// 8. 统计
	count, _ := db.CountUsers()
	fmt.Printf("当前用户数量: %d\n", count)

	fmt.Println("GORM 操作示例完成")
}
