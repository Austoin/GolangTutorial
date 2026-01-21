// database/database_sql.go
// Go 原生 SQL 数据库操作 - 详细注释版

package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	// 导入数据库驱动
	// Go 使用数据库驱动分离的架构
	// 常用的驱动：
	//   - github.com/go-sql-driver/mysql (MySQL)
	//   - github.com/lib/pq (PostgreSQL)
	//   - github.com/mattn/go-sqlite3 (SQLite)
	//   - github.com/denisenkom/go-mssqldb (SQL Server)
	_ "github.com/go-sql-driver/mysql"
)

// ====== 数据库连接基础 ======

// User 结构体表示用户表的数据模型
// 使用标签（tag）来映射数据库列名
type User struct {
	ID        int64     `json:"id"`         // 用户唯一标识
	Username  string    `json:"username"`   // 用户名
	Email     string    `json:"email"`      // 邮箱
	Password  string    `json:"-"`          // 密码不序列化
	CreatedAt time.Time `json:"created_at"` // 创建时间
	UpdatedAt time.Time `json:"updated_at"` // 更新时间
}

// UserModel 数据库操作封装
type UserModel struct {
	db *sql.DB // 数据库连接实例
}

// NewUserModel 创建用户模型
func NewUserModel(dsn string) (*UserModel, error) {
	// 1. 打开数据库连接
	// Open 返回一个 sql.DB 实例，它代表一个连接池
	// 参数：驱动名称、数据源名称（DSN）
	// DSN 格式因驱动而异
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("打开数据库失败: %w", err)
	}

	// 2. 配置连接池
	// SetMaxOpenConns 设置最大打开连接数
	// 0 表示无限制（受系统资源限制）
	db.SetMaxOpenConns(25)

	// SetMaxIdleConns 设置最大空闲连接数
	// 这些连接会被保留，避免频繁创建和销毁
	db.SetMaxIdleConns(5)

	// SetConnMaxLifetime 设置连接的最大存活时间
	// 超过这个时间的连接会被替换
	db.SetConnMaxLifetime(5 * time.Minute)

	// 3. 验证连接
	// Ping 检查数据库是否可达
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("数据库连接失败: %w", err)
	}

	log.Println("数据库连接成功")

	return &UserModel{db: db}, nil
}

// Close 关闭数据库连接
func (m *UserModel) Close() error {
	// Close 关闭数据库连接，释放所有资源
	// 已有的连接会等待完成，新请求会失败
	return m.db.Close()
}

// ====== 创建表 ======

// CreateTable 创建用户表
func (m *UserModel) CreateTable() error {
	// 1. 定义建表 SQL
	// 使用 ? 作为占位符，防止 SQL 注入
	// 注意：不同数据库的占位符可能不同
	//   MySQL: ?
	//   PostgreSQL: $1, $2...
	//   SQLite: ? 或 $name
	query := `
		CREATE TABLE IF NOT EXISTS users (
			id BIGINT AUTO_INCREMENT PRIMARY KEY,
			username VARCHAR(50) NOT NULL UNIQUE,
			email VARCHAR(100) NOT NULL UNIQUE,
			password VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			INDEX idx_username (username),
			INDEX idx_email (email)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
	`

	// 2. 执行 SQL
	// Exec 执行不返回行的 SQL，如 INSERT、UPDATE、DELETE、CREATE
	// 返回 Result 接口，包含 LastInsertId 和 RowsAffected
	result, err := m.db.Exec(query)
	if err != nil {
		return fmt.Errorf("创建表失败: %w", err)
	}

	// 3. 检查影响行数（可选）
	// RowsAffected 返回受影响的行数
	_ = result.RowsAffected()

	log.Println("用户表创建成功")
	return nil
}

// DropTable 删除用户表
func (m *UserModel) DropTable() error {
	query := "DROP TABLE IF EXISTS users"
	_, err := m.db.Exec(query)
	return err
}

// ====== 插入数据 ======

// InsertUser 插入单个用户
func (m *UserModel) InsertUser(user *User) (int64, error) {
	// 1. 准备 SQL 语句
	// Prepare 预编译 SQL，提高执行效率
	// 对于多次执行的语句，预编译可以提升性能
	stmt, err := m.db.Prepare(`
		INSERT INTO users (username, email, password)
		VALUES (?, ?, ?)
	`)
	if err != nil {
		return 0, fmt.Errorf("准备语句失败: %w", err)
	}
	defer stmt.Close() // 释放资源

	// 2. 执行插入
	// Exec 也可以直接执行，不需要预编译
	result, err := stmt.Exec(user.Username, user.Email, user.Password)
	if err != nil {
		return 0, fmt.Errorf("插入失败: %w", err)
	}

	// 3. 获取插入的 ID
	// LastInsertId 返回最后插入的自增 ID
	lastID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("获取插入ID失败: %w", err)
	}

	return lastID, nil
}

// InsertUsers 批量插入用户
func (m *UserModel) InsertUsers(users []User) error {
	// 使用事务进行批量操作
	// 事务可以确保操作的原子性
	tx, err := m.db.Begin()
	if err != nil {
		return fmt.Errorf("开始事务失败: %w", err)
	}
	defer tx.Rollback() // 如果失败，回滚事务

	// 准备语句
	stmt, err := tx.Prepare(`
		INSERT INTO users (username, email, password)
		VALUES (?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// 批量执行
	for _, user := range users {
		_, err := stmt.Exec(user.Username, user.Email, user.Password)
		if err != nil {
			return fmt.Errorf("插入用户 %s 失败: %w", user.Username, err)
		}
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("提交事务失败: %w", err)
	}

	return nil
}

// ====== 查询数据 ======

// GetUserByID 根据 ID 查询用户
func (m *UserModel) GetUserByID(id int64) (*User, error) {
	// 1. 查询单行数据
	// QueryRow 查询一行数据，返回 *sql.Row
	query := "SELECT id, username, email, password, created_at, updated_at FROM users WHERE id = ?"
	row := m.db.QueryRow(query, id)

	// 2. 扫描数据到结构体
	// Scan 自动将列值转换为目标类型
	// 注意：参数数量和类型必须匹配
	user := &User{}
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password,
		&user.CreatedAt, &user.UpdatedAt)

	// 3. 处理查询结果
	if err == sql.ErrNoRows {
		// 没有找到记录
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("扫描数据失败: %w", err)
	}

	return user, nil
}

// GetUserByUsername 根据用户名查询用户
func (m *UserModel) GetUserByUsername(username string) (*User, error) {
	query := "SELECT id, username, email, password, created_at, updated_at FROM users WHERE username = ?"
	row := m.db.QueryRow(query, username)

	user := &User{}
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password,
		&user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetAllUsers 查询所有用户
func (m *UserModel) GetAllUsers() ([]User, error) {
	// 1. 查询多行数据
	// Query 返回 *sql.Rows，包含所有匹配的行
	query := "SELECT id, username, email, password, created_at, updated_at FROM users ORDER BY id"
	rows, err := m.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("查询失败: %w", err)
	}
	defer rows.Close() // 关闭 Rows，释放资源

	// 2. 遍历结果集
	var users []User
	for rows.Next() {
		user := User{}
		// 3. 扫描每一行
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Password,
			&user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("扫描行失败: %w", err)
		}
		users = append(users, user)
	}

	// 4. 检查遍历错误
	// 即使成功遍历完所有行，也应该检查是否有错误
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("遍历结果失败: %w", err)
	}

	return users, nil
}

// GetUsersByEmailPrefix 按邮箱前缀查询用户
func (m *UserModel) GetUsersByEmailPrefix(prefix string) ([]User, error) {
	// 使用 LIKE 进行模糊查询
	// % 匹配任意字符序列
	query := "SELECT id, username, email, password, created_at, updated_at FROM users WHERE email LIKE ?"

	// 执行查询
	rows, err := m.db.Query(query, prefix+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		user := User{}
		if err := rows.Scan(&user.ID, &user.Username, &user.Email,
			&user.Password, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, rows.Err()
}

// CountUsers 统计用户数量
func (m *UserModel) CountUsers() (int64, error) {
	query := "SELECT COUNT(*) FROM users"
	var count int64
	err := m.db.QueryRow(query).Scan(&count)
	return count, err
}

// ====== 更新数据 ======

// UpdateUser 更新用户信息
func (m *UserModel) UpdateUser(user *User) error {
	query := `
		UPDATE users 
		SET username = ?, email = ?, password = ?, updated_at = NOW()
		WHERE id = ?
	`
	result, err := m.db.Exec(query, user.Username, user.Email, user.Password, user.ID)
	if err != nil {
		return fmt.Errorf("更新失败: %w", err)
	}

	// 检查是否有行被更新
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("用户不存在: ID=%d", user.ID)
	}

	return nil
}

// UpdatePassword 更新用户密码
func (m *UserModel) UpdatePassword(id int64, newPassword string) error {
	query := "UPDATE users SET password = ?, updated_at = NOW() WHERE id = ?"
	_, err := m.db.Exec(query, newPassword, id)
	return err
}

// ====== 删除数据 ======

// DeleteUserByID 根据 ID 删除用户
func (m *UserModel) DeleteUserByID(id int64) error {
	query := "DELETE FROM users WHERE id = ?"
	result, err := m.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("删除失败: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("用户不存在: ID=%d", id)
	}

	return nil
}

// DeleteUserByUsername 根据用户名删除用户
func (m *UserModel) DeleteUserByUsername(username string) error {
	query := "DELETE FROM users WHERE username = ?"
	result, err := m.db.Exec(query, username)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("用户不存在: username=%s", username)
	}

	return nil
}

// ====== 事务操作 ======

// TransferMoney 转账示例（使用事务）
func (m *UserModel) TransferMoney(fromID, toID int64, amount float64) error {
	// 1. 开始事务
	// Begin 返回 tx，用于执行事务内的操作
	tx, err := m.db.Begin()
	if err != nil {
		return fmt.Errorf("开始事务失败: %w", err)
	}

	// 确保回滚
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// 2. 扣除转出账户余额
	// FOR UPDATE 锁定行，防止并发修改
	query1 := "UPDATE users SET balance = balance - ? WHERE id = ? AND balance >= ?"
	result1, err := tx.Exec(query1, amount, fromID, amount)
	if err != nil {
		return fmt.Errorf("扣除余额失败: %w", err)
	}

	rows1, err := result1.RowsAffected()
	if err != nil {
		return err
	}

	if rows1 == 0 {
		return fmt.Errorf("余额不足或账户不存在: ID=%d", fromID)
	}

	// 3. 转入账户增加余额
	query2 := "UPDATE users SET balance = balance + ? WHERE id = ?"
	_, err = tx.Exec(query2, amount, toID)
	if err != nil {
		return fmt.Errorf("增加余额失败: %w", err)
	}

	// 4. 提交事务
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("提交事务失败: %w", err)
	}

	return nil
}

// ====== 错误处理 ======

// HandleSQLError 处理 SQL 错误
// 根据错误类型返回友好的错误信息
func HandleSQLError(err error) error {
	if err == nil {
		return nil
	}

	// 检查是否是数据库连接错误
	if err == sql.ErrConnDone {
		return fmt.Errorf("数据库连接已关闭，请重新连接")
	}

	// 检查是否是事务相关错误
	if err == sql.ErrTxDone {
		return fmt.Errorf("事务已完成或已回滚")
	}

	// 检查是否是 No Rows 错误
	if err == sql.ErrNoRows {
		return fmt.Errorf("未找到匹配的数据")
	}

	// 处理 MySQL 特有的错误
	// 这里可以添加更多的错误码判断
	return err
}

// ====== 主函数 ======

func main() {
	fmt.Println("=== Go 原生 SQL 数据库操作示例 ===")

	// 1. 创建数据库连接
	// DSN 格式：用户名:密码@协议(地址)/数据库名?参数
	dsn := "root:password@tcp(localhost:3306)/testdb?charset=utf8mb4&parseTime=True"

	model, err := NewUserModel(dsn)
	if err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}
	defer model.Close()

	// 2. 创建表
	if err := model.CreateTable(); err != nil {
		log.Fatalf("创建表失败: %v", err)
	}

	// 3. 插入测试数据
	users := []User{
		{Username: "alice", Email: "alice@example.com", Password: "pass123"},
		{Username: "bob", Email: "bob@example.com", Password: "pass456"},
		{Username: "charlie", Email: "charlie@example.com", Password: "pass789"},
	}

	if err := model.InsertUsers(users); err != nil {
		log.Printf("批量插入失败: %v", err)
	}

	// 4. 查询测试
	// 查询单个用户
	user, err := model.GetUserByID(1)
	if err != nil {
		log.Printf("查询用户失败: %v", err)
	} else if user != nil {
		fmt.Printf("查询到用户: %s (%s)\n", user.Username, user.Email)
	}

	// 查询所有用户
	allUsers, err := model.GetAllUsers()
	if err != nil {
		log.Printf("查询所有用户失败: %v", err)
	} else {
		fmt.Printf("共有 %d 个用户:\n", len(allUsers))
		for _, u := range allUsers {
			fmt.Printf("  - %s: %s\n", u.Username, u.Email)
		}
	}

	// 5. 更新测试
	user, _ = model.GetUserByUsername("alice")
	if user != nil {
		user.Email = "alice.new@example.com"
		if err := model.UpdateUser(user); err != nil {
			log.Printf("更新用户失败: %v", err)
		} else {
			fmt.Println("更新用户成功")
		}
	}

	// 6. 删除测试
	if err := model.DeleteUserByUsername("charlie"); err != nil {
		log.Printf("删除用户失败: %v", err)
	} else {
		fmt.Println("删除用户成功")
	}

	// 7. 统计
	count, _ := model.CountUsers()
	fmt.Printf("当前用户数量: %d\n", count)

	// 8. 清理（可选）
	// model.DropTable()

	fmt.Println("数据库操作示例完成")
}
