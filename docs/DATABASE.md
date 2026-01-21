# Go 数据库操作详解

## 目录
- [1. 数据库操作概述](#1-数据库操作概述)
- [2. 原生 SQL 操作](#2-原生-sql-操作)
- [3. GORM 基础](#3-gorm-基础)
- [4. GORM 高级特性](#4-gorm-高级特性)
- [5. 连接池与事务](#5-连接池与事务)
- [6. 最佳实践](#6-最佳实践)

---

## 1. 数据库操作概述

### 1.1 Go 数据库接口

Go 标准库定义了 `database/sql` 接口，提供了统一的数据库访问方式。

```go
// 核心接口定义
type DB struct {
    // 内部实现
}

type Rows struct {
    // 游标结果集
}

type Row struct {
    // 单行结果
}

type Stmt struct {
    // 预编译语句
}

type Tx struct {
    // 事务
}

// 常用方法
func (db *DB) Open(driverName, dataSourceName string) (*DB, error)
func (db *DB) Close() error
func (db *DB) Ping() error
func (db *DB) Exec(query string, args ...interface{}) (Result, error)
func (db *DB) Query(query string, args ...interface{}) (*Rows, error)
func (db *DB) QueryRow(query string, args ...interface{}) *Row
func (db *DB) Prepare(query string) (*Stmt, error)
func (db *DB) Begin() (*Tx, error)
```

### 1.2 常用数据库驱动

```go
// MySQL
import _ "github.com/go-sql-driver/mysql"

// PostgreSQL
import _ "github.com/lib/pq"

// SQLite
import _ "github.com/mattn/go-sqlite3"

// SQL Server
import _ "github.com/denisenkom/go-mssqldb"

// MongoDB (NoSQL)
import "go.mongodb.org/mongo-driver/mongo"
```

---

## 2. 原生 SQL 操作

### 2.1 数据库连接

```go
package main

import (
    "database/sql"
    "fmt"
    "log"
    "time"

    _ "github.com/go-sql-driver/mysql"
)

func main() {
    // MySQL 连接
    dsn := "user:password@tcp(localhost:3306)/database?charset=utf8mb4&parseTime=True&loc=Local"
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        log.Fatal("连接数据库失败:", err)
    }
    defer db.Close()

    // 配置连接池
    db.SetMaxOpenConns(25)           // 最大打开连接数
    db.SetMaxIdleConns(25)           // 最大空闲连接数
    db.SetConnMaxLifetime(5 * time.Minute) // 连接最大生命周期

    // 测试连接
    if err := db.Ping(); err != nil {
        log.Fatal("Ping 失败:", err)
    }
    fmt.Println("成功连接到 MySQL 数据库")
}
```

### 2.2 插入操作

```go
package main

import (
    "database/sql"
    "fmt"
    "log"

    _ "github.com/go-sql-driver/mysql"
)

type User struct {
    ID        int64
    Username  string
    Email     string
    Age       int
    CreatedAt time.Time
}

func main() {
    dsn := "user:password@tcp(localhost:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local"
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // 方法1：Exec 直接插入
    result, err := db.Exec(`
        INSERT INTO users (username, email, age, created_at)
        VALUES (?, ?, ?, ?)`, "张三", "zhangsan@example.com", 25, time.Now())
    if err != nil {
        log.Fatal("插入失败:", err)
    }

    // 获取插入的 ID
    id, err := result.LastInsertId()
    if err != nil {
        log.Fatal("获取 ID 失败:", err)
    }
    fmt.Printf("插入成功，ID: %d\n", id)

    // 方法2：使用预处理语句
    stmt, err := db.Prepare(`
        INSERT INTO users (username, email, age, created_at)
        VALUES (?, ?, ?, ?)`)
    if err != nil {
        log.Fatal("预处理失败:", err)
    }
    defer stmt.Close()

    result, err = stmt.Exec("李四", "lisi@example.com", 30, time.Now())
    id, err = result.LastInsertId()
    fmt.Printf("预处理插入成功，ID: %d\n", id)

    // 方法3：批量插入
    users := []User{
        {Username: "王五", Email: "wangwu@example.com", Age: 28},
        {Username: "赵六", Email: "zhaoliu@example.com", Age: 32},
    }

    tx, err := db.Begin()
    if err != nil {
        log.Fatal("开始事务失败:", err)
    }

    stmt, _ = tx.Prepare(`
        INSERT INTO users (username, email, age, created_at)
        VALUES (?, ?, ?, ?)`)
    defer stmt.Close()

    for _, user := range users {
        _, err = stmt.Exec(user.Username, user.Email, user.Age, time.Now())
        if err != nil {
            tx.Rollback()
            log.Fatal("批量插入失败:", err)
        }
    }

    tx.Commit()
    fmt.Println("批量插入成功")
}
```

### 2.3 查询操作

```go
package main

import (
    "database/sql"
    "fmt"
    "log"

    _ "github.com/go-sql-driver/mysql"
)

func main() {
    dsn := "user:password@tcp(localhost:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local"
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // 方法1：QueryRow 查询单行
    var user User
    err = db.QueryRow(`
        SELECT id, username, email, age, created_at
        FROM users WHERE id = ?`, 1).Scan(
        &user.ID, &user.Username, &user.Email, 
        &user.Age, &user.CreatedAt)
    if err != nil {
        if err == sql.ErrNoRows {
            fmt.Println("未找到用户")
        } else {
            log.Fatal("查询失败:", err)
        }
    } else {
        fmt.Printf("查询结果: %+v\n", user)
    }

    // 方法2：Query 查询多行
    rows, err := db.Query(`
        SELECT id, username, email, age
        FROM users WHERE age > ? ORDER BY age DESC`, 20)
    if err != nil {
        log.Fatal("查询失败:", err)
    }
    defer rows.Close()

    var users []User
    for rows.Next() {
        var user User
        err := rows.Scan(&user.ID, &user.Username, 
            &user.Email, &user.Age)
        if err != nil {
            log.Fatal("扫描行失败:", err)
        }
        users = append(users, user)
    }

    if err := rows.Err(); err != nil {
        log.Fatal("遍历结果集失败:", err)
    }

    fmt.Printf("查询到 %d 个用户:\n", len(users))
    for _, u := range users {
        fmt.Printf("  ID: %d, 用户名: %s, 邮箱: %s, 年龄: %d\n",
            u.ID, u.Username, u.Email, u.Age)
    }

    // 方法3：使用预处理语句查询
    stmt, err := db.Prepare(`
        SELECT id, username, email FROM users WHERE username LIKE ?`)
    if err != nil {
        log.Fatal("预处理失败:", err)
    }
    defer stmt.Close()

    rows, err = stmt.Query("%张%")
    if err != nil {
        log.Fatal("查询失败:", err)
    }
    defer rows.Close()

    for rows.Next() {
        var id int
        var username, email string
        rows.Scan(&id, &username, &email)
        fmt.Printf("模糊查询: %s (%s)\n", username, email)
    }
}
```

### 2.4 更新和删除操作

```go
package main

import (
    "database/sql"
    "fmt"
    "log"

    _ "github.com/go-sql-driver/mysql"
)

func main() {
    dsn := "user:password@tcp(localhost:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local"
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // 更新操作
    result, err := db.Exec(`
        UPDATE users SET age = age + 1 WHERE id = ?`, 1)
    if err != nil {
        log.Fatal("更新失败:", err)
    }

    affected, err := result.RowsAffected()
    if err != nil {
        log.Fatal("获取影响行数失败:", err)
    }
    fmt.Printf("更新了 %d 行\n", affected)

    // 批量更新（事务）
    tx, err := db.Begin()
    if err != nil {
        log.Fatal("开始事务失败:", err)
    }

    updates := []struct {
        id   int
        age  int
    }{
        {1, 26},
        {2, 31},
    }

    stmt, _ := tx.Prepare(`UPDATE users SET age = ? WHERE id = ?`)
    defer stmt.Close()

    for _, u := range updates {
        _, err = stmt.Exec(u.age, u.id)
        if err != nil {
            tx.Rollback()
            log.Fatal("批量更新失败:", err)
        }
    }

    tx.Commit()
    fmt.Println("批量更新成功")

    // 删除操作
    result, err = db.Exec(`DELETE FROM users WHERE id = ?`, 5)
    if err != nil {
        log.Fatal("删除失败:", err)
    }

    affected, err = result.RowsAffected()
    fmt.Printf("删除了 %d 行\n", affected)

    // 软删除（使用状态字段）
    _, err = db.Exec(`UPDATE users SET status = 'deleted' WHERE id = ?`, 3)
    if err != nil {
        log.Fatal("软删除失败:", err)
    }
}
```

### 2.5 事务处理

```go
package main

import (
    "database/sql"
    "fmt"
    "log"

    _ "github.com/go-sql-driver/mysql"
)

func transferMoney(db *sql.DB, fromID, toID int, amount float64) error {
    // 开始事务
    tx, err := db.Begin()
    if err != nil {
        return fmt.Errorf("开始事务失败: %v", err)
    }
    defer func() {
        if p := recover(); p != nil {
            tx.Rollback()
            panic(p)
        }
    }()

    // 查询转出人余额
    var fromBalance float64
    err = tx.QueryRow(`SELECT balance FROM accounts WHERE id = ? FOR UPDATE`,
        fromID).Scan(&fromBalance)
    if err != nil {
        tx.Rollback()
        return fmt.Errorf("查询转出人余额失败: %v", err)
    }

    if fromBalance < amount {
        tx.Rollback()
        return fmt.Errorf("余额不足")
    }

    // 扣减转出人余额
    _, err = tx.Exec(`UPDATE accounts SET balance = balance - ? WHERE id = ?`,
        amount, fromID)
    if err != nil {
        tx.Rollback()
        return fmt.Errorf("扣减余额失败: %v", err)
    }

    // 增加转入人余额
    _, err = tx.Exec(`UPDATE accounts SET balance = balance + ? WHERE id = ?`,
        amount, toID)
    if err != nil {
        tx.Rollback()
        return fmt.Errorf("增加余额失败: %v", err)
    }

    // 记录交易
    _, err = tx.Exec(`INSERT INTO transactions (from_id, to_id, amount) 
        VALUES (?, ?, ?)`, fromID, toID, amount)
    if err != nil {
        tx.Rollback()
        return fmt.Errorf("记录交易失败: %v", err)
    }

    // 提交事务
    if err = tx.Commit(); err != nil {
        return fmt.Errorf("提交事务失败: %v", err)
    }

    return nil
}

func main() {
    dsn := "user:password@tcp(localhost:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local"
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    err = transferMoney(db, 1, 2, 100.0)
    if err != nil {
        log.Printf("转账失败: %v", err)
    } else {
        fmt.Println("转账成功")
    }
}
```

---

## 3. GORM 基础

### 3.1 GORM 简介

GORM 是一个功能强大的 Go 语言 ORM 库，支持关联、事务、钩子等功能。

```go
// 安装
go get -u gorm.io/gorm
go get -u gorm.io/driver/mysql
```

### 3.2 模型定义

```go
package main

import (
    "fmt"
    "time"

    "gorm.io/gorm"
)

// BaseModel 包含通用字段
type BaseModel struct {
    ID        uint      `gorm:"primaryKey"`
    CreatedAt time.Time `gorm:"autoCreateTime"`
    UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// User 用户模型
type User struct {
    BaseModel
    Username  string    `gorm:"size:50;uniqueIndex;not null"`
    Email     string    `gorm:"size:100;uniqueIndex"`
    Password  string    `gorm:"size:255;not null"`
    Age       int       `gorm:"default:0"`
    Birthday  *time.Time
    Status    int8      `gorm:"default:1"` // 1: 正常, 0: 禁用
    Orders    []Order                         // 一对多关联
    CreditCard CreditCard `gorm:"foreignKey:UserID"` // 一对一关联
}

// Order 订单模型
type Order struct {
    gorm.Model
    UserID      uint    `gorm:"index"`
    OrderNo     string  `gorm:"size:50;uniqueIndex"`
    Amount      float64 `gorm:"type:decimal(10,2)"`
    Status      string  `gorm:"size:20;default:'pending'"`
    User        User    `gorm:"foreignKey:UserID"` // 关联
}

// CreditCard 信用卡模型
type CreditCard struct {
    BaseModel
    UserID     uint   `gorm:"uniqueIndex"`
    CardNumber string `gorm:"size:20"`
    ExpiryDate string `gorm:"size:10"`
    CVV        string `gorm:"size:4"`
}

// Product 产品模型（多对多关联）
type Product struct {
    gorm.Model
    Name        string    `gorm:"size:100"`
    Price       float64   `gorm:"type:decimal(10,2)"`
    CategoryID  uint
    Category    Category  `gorm:"foreignKey:CategoryID"`
    Tags        []Tag     `gorm:"many2many:product_tags;"`
}

// Category 分类模型
type Category struct {
    gorm.Model
    Name  string    `gorm:"size:50"`
    Products []Product
}

// Tag 标签模型
type Tag struct {
    gorm.Model
    Name    string    `gorm:"size:30;uniqueIndex"`
    Products []Product `gorm:"many2many:product_tags;"`
}

// 使用标签自定义字段
type Person struct {
    gorm.Model
    Name string `gorm:"<-:create"` // 允许创建和读取，但不允许更新
    Age  int    `gorm:"->:false"`  // 只允许写入
}
```

### 3.3 数据库连接

```go
package main

import (
    "fmt"
    "log"
    "time"

    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
    "gorm.io/gorm/schema"
)

func main() {
    dsn := "user:password@tcp(localhost:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local"

    // 配置选项
    config := &gorm.Config{
        // 命名策略
        NamingStrategy: schema.NamingStrategy{
            TablePrefix:   "t_",       // 表前缀
            SingularTable: true,       // 使用单数表名
        },
        // 日志配置
        Logger: logger.Default.LogMode(logger.Info),
        // 跳过默认事务
        SkipDefaultTransaction: false,
        // 批量插入大小
        CreateBatchSize: 1000,
    }

    db, err := gorm.Open(mysql.Open(dsn), config)
    if err != nil {
        log.Fatal("连接数据库失败:", err)
    }

    // 获取底层 *sql.DB
    sqlDB, err := db.DB()
    if err != nil {
        log.Fatal("获取 DB 失败:", err)
    }

    // 配置连接池
    sqlDB.SetMaxOpenConns(25)
    sqlDB.SetMaxIdleConns(25)
    sqlDB.SetConnMaxLifetime(5 * time.Minute)

    fmt.Println("成功连接到数据库")
}
```

### 3.4 CRUD 操作

```go
package main

import (
    "fmt"
    "log"

    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

type User struct {
    ID       uint   `gorm:"primaryKey"`
    Username string `gorm:"size:50"`
    Email    string `gorm:"size:100;uniqueIndex"`
    Age      int
}

func main() {
    dsn := "user:password@tcp(localhost:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local"
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal(err)
    }

    // ========== 创建 ==========
    
    // 方法1：创建单条记录
    user := User{Username: "张三", Email: "zhangsan@example.com", Age: 25}
    result := db.Create(&user)
    if result.Error != nil {
        log.Fatal("创建失败:", result.Error)
    }
    fmt.Printf("创建成功，ID: %d\n", user.ID)

    // 方法2：批量创建
    users := []User{
        {Username: "李四", Email: "lisi@example.com", Age: 30},
        {Username: "王五", Email: "wangwu@example.com", Age: 28},
    }
    result = db.CreateInBatches(users, 100)
    fmt.Printf("批量创建成功: %d 行\n", result.RowsAffected)

    // 方法3：使用 CreateWithAssociations 创建关联
    // user := User{Username: "赵六", Email: "zhaoliu@example.com", Age: 35}
    // db.Create(&user)

    // ========== 读取 ==========

    // 方法1：获取第一条记录
    var firstUser User
    db.First(&firstUser)
    fmt.Printf("第一条: %+v\n", firstUser)

    // 方法2：获取最后一条记录
    var lastUser User
    db.Last(&lastUser)
    fmt.Printf("最后一条: %+v\n", lastUser)

    // 方法3：根据主键获取
    var userByID User
    db.First(&userByID, 5)
    fmt.Printf("ID=5: %+v\n", userByID)

    // 方法4：条件查询
    var users []User
    db.Where("age > ? AND status = ?", 20, 1).Find(&users)
    fmt.Printf("条件查询: %d 用户\n", len(users))

    // 方法5：IN 查询
    db.Where("username IN ?", []string{"张三", "李四"}).Find(&users)

    // 方法6：LIKE 查询
    db.Where("username LIKE ?", "%张%").Find(&users)

    // 方法7：链式查询
    users = []User{}
    db.Where("age >= ?", 20).Where("status = ?", 1).Order("age DESC").Find(&users)

    // ========== 更新 ==========

    // 方法1：更新单字段
    db.Model(&User{}).Where("id = ?", 1).Update("age", 26)

    // 方法2：更新多字段
    db.Model(&User{ID: 1}).Updates(User{Username: "张三三", Age: 26})

    // 方法3：使用 Struct 更新（只会更新非零值字段）
    db.Model(&User{ID: 1}).Updates(User{Age: 27})

    // 方法4：使用 Map 更新
    db.Model(&User{}).Where("id = ?", 1).Updates(map[string]interface{}{
        "age":      28,
        "username": "张三四",
    })

    // ========== 删除 ==========

    // 方法1：删除记录
    db.Delete(&User{ID: 1})

    // 方法2：批量删除
    db.Where("age < ?", 18).Delete(&User{})

    // 方法3：物理删除（需要使用 Unscoped）
    // db.Unscoped().Delete(&User{ID: 2})
}
```

---

## 4. GORM 高级特性

### 4.1 关联关系

```go
package main

import (
    "fmt"
    "log"

    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

type User struct {
    ID        uint   `gorm:"primaryKey"`
    Username  string
    Email     string `gorm:"uniqueIndex"`
    Orders    []Order
}

type Order struct {
    ID        uint    `gorm:"primaryKey"`
    OrderNo   string  `gorm:"uniqueIndex"`
    Amount    float64
    UserID    uint
    User      User    `gorm:"foreignKey:UserID"`
}

func main() {
    dsn := "user:password@tcp(localhost:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local"
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal(err)
    }

    // ========== 一对多 (Has Many) ==========

    // 创建用户时同时创建订单
    user := User{
        Username: "张三",
        Email:    "zhangsan@example.com",
        Orders: []Order{
            {OrderNo: "ORD001", Amount: 100.0},
            {OrderNo: "ORD002", Amount: 200.0},
        },
    }
    db.Create(&user)

    // 查询用户及其订单
    var userWithOrders User
    db.Preload("Orders").First(&userWithOrders, user.ID)
    fmt.Printf("用户: %s, 订单数: %d\n", 
        userWithOrders.Username, len(userWithOrders.Orders))

    // ========== 一对一 (Belongs To / Has One) ==========

    // Belongs To: Order 属于 User
    var orders []Order
    db.Preload("User").Find(&orders)
    for _, order := range orders {
        fmt.Printf("订单 %s 属于用户 %s\n", order.OrderNo, order.User.Username)
    }

    // ========== 多对多 (Many To Many) ==========

    type Product struct {
        ID    uint   `gorm:"primaryKey"`
        Name  string
        Tags  []Tag  `gorm:"many2many:product_tags;"`
    }

    type Tag struct {
        ID     uint   `gorm:"primaryKey"`
        Name   string
        Products []Product `gorm:"many2many:product_tags;"`
    }

    // 创建产品并关联标签
    product := Product{
        Name: "iPhone",
        Tags: []Tag{
            {Name: "手机"},
            {Name: "苹果"},
            {Name: "高端"},
        },
    }
    db.Create(&product)

    // 查询产品及其标签
    var productWithTags Product
    db.Preload("Tags").First(&productWithTags, product.ID)
    fmt.Printf("产品: %s\n标签: ", productWithTags.Name)
    for _, tag := range productWithTags.Tags {
        fmt.Printf("%s ", tag.Name)
    }
    fmt.Println()

    // ========== 关联模式 ==========

    var user User
    db.First(&user, 1)

    // 替换关联
    user.Orders = []Order{
        {OrderNo: "ORD003", Amount: 300.0},
    }
    db.Save(&user)

    // 添加关联
    order := Order{OrderNo: "ORD004", Amount: 400.0}
    db.Model(&user).Append(&order)

    // 删除关联
    db.Model(&user).Association("Orders").Delete(order)

    // 清空关联
    // db.Model(&user).Association("Orders").Clear()

    // 统计关联数量
    var orderCount int64
    db.Model(&user).Association("Orders").Count(&orderCount)
    fmt.Printf("用户订单数: %d\n", orderCount)
}
```

### 4.2 钩子函数

```go
package main

import (
    "fmt"
    "log"

    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

type User struct {
    ID        uint   `gorm:"primaryKey"`
    Username  string
    Email     string
    Age       int
    CreatedAt gorm.DeletedAt // 软删除
}

// BeforeCreate 创建前钩子
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
    fmt.Println("BeforeCreate: 准备创建用户")
    
    // 设置默认值
    if u.Age == 0 {
        u.Age = 18
    }
    
    // 验证
    if u.Username == "" {
        err = fmt.Errorf("用户名不能为空")
    }
    return
}

// AfterCreate 创建后钩子
func (u *User) AfterCreate(tx *gorm.DB) (err error) {
    fmt.Printf("AfterCreate: 用户 %s 创建成功\n", u.Username)
    return nil
}

// BeforeUpdate 更新前钩子
func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
    fmt.Printf("BeforeUpdate: 准备更新用户 %s\n", u.Username)
    return
}

// AfterUpdate 更新后钩子
func (u *User) AfterUpdate(tx *gorm.DB) (err error) {
    fmt.Printf("AfterUpdate: 用户 %s 更新成功\n", u.Username)
    return
}

// BeforeDelete 删除前钩子
func (u *User) BeforeDelete(tx *gorm.DB) (err error) {
    fmt.Printf("BeforeDelete: 准备删除用户 %s\n", u.Username)
    return
}

// AfterFind 查询后钩子
func (u *User) AfterFind(tx *gorm.DB) (err error) {
    fmt.Printf("AfterFind: 查到用户 %s\n", u.Username)
    return
}

func main() {
    dsn := "user:password@tcp(localhost:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local"
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal(err)
    }

    // 创建用户（会触发 BeforeCreate 和 AfterCreate）
    user := User{Username: "测试用户", Email: "test@example.com"}
    db.Create(&user)

    // 查询用户（会触发 AfterFind）
    var foundUser User
    db.First(&foundUser, user.ID)

    // 更新用户（会触发 BeforeUpdate 和 AfterUpdate）
    db.Model(&foundUser).Update("age", 25)

    // 删除用户（会触发 BeforeDelete）
    db.Delete(&foundUser)
}
```

### 4.3 事务

```go
package main

import (
    "fmt"
    "log"

    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

func main() {
    dsn := "user:password@tcp(localhost:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local"
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal(err)
    }

    // ========== 方法1：使用 Transaction ==========
    
    err = db.Transaction(func(tx *gorm.DB) error {
        // 在事务中进行的操作
        var user User
        tx.First(&user, 1)
        
        // 更新
        if err := tx.Model(&user).Update("age", 30).Error; err != nil {
            return err
        }
        
        // 创建记录
        order := Order{OrderNo: "TX001", Amount: 100.0, UserID: user.ID}
        if err := tx.Create(&order).Error; err != nil {
            return err
        }
        
        return nil
    })

    if err != nil {
        log.Printf("事务失败: %v", err)
    } else {
        fmt.Println("事务成功")
    }

    // ========== 方法2：手动控制事务 ==========
    
    tx := db.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()

    if err := tx.Error; err != nil {
        log.Fatal("开始事务失败:", err)
    }

    // 操作1
    if err := tx.Exec("UPDATE users SET balance = balance - 100 WHERE id = 1").Error; err != nil {
        tx.Rollback()
        log.Fatal("操作1失败:", err)
    }

    // 操作2
    if err := tx.Exec("UPDATE users SET balance = balance + 100 WHERE id = 2").Error; err != nil {
        tx.Rollback()
        log.Fatal("操作2失败:", err)
    }

    if err := tx.Commit().Error; err != nil {
        log.Fatal("提交事务失败:", err)
    }

    fmt.Println("手动事务成功")
}
```

### 4.4 原生 SQL

```go
package main

import (
    "fmt"
    "log"

    "gorm.io/datormysql"
    "gorm.io/gorm"
)

func main() {
    dsn := "user:password@tcp(localhost:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local"
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal(err)
    }

    // ========== 原生 SQL 查询 ==========

    // 查询到 map
    var results []map[string]interface{}
    db.Raw("SELECT id, username, age FROM users WHERE age > ?", 20).Scan(&results)
    for _, row := range results {
        fmt.Printf("ID: %v, Username: %v, Age: %v\n", 
            row["id"], row["username"], row["age"])
    }

    // 查询到结构体
    var users []User
    db.Raw("SELECT * FROM users WHERE age > ?", 20).Scan(&users)

    // Raw 执行插入
    db.Exec("INSERT INTO users (username, email, age) VALUES (?, ?, ?)", 
        "新用户", "new@example.com", 25)

    // ========== 命名参数 ==========

    var user User
    db.Raw("SELECT * FROM users WHERE username = @username AND age > @age",
        map[string]interface{}{"username": "张三", "age": 20}).Scan(&user)

    // ========== 子查询 ==========

    var orders []Order
    subQuery := db.Model(&User{}).Select("id").Where("age > ?", 20)
    db.Where("user_id IN (?)", subQuery).Find(&orders)

    // ========== 连接查询 ==========

    type Result struct {
        Username string
        OrderNo  string
        Amount   float64
    }

    var results2 []Result
    db.Raw(`
        SELECT u.username, o.order_no, o.amount
        FROM users u
        LEFT JOIN orders o ON u.id = o.user_id
        WHERE u.age > ?
    `, 20).Scan(&results2)

    for _, r := range results2 {
        fmt.Printf("%s 的订单: %s, 金额: %.2f\n", 
            r.Username, r.OrderNo, r.Amount)
    }
}
```

---

## 5. 连接池与事务

### 5.1 连接池配置

```go
package main

import (
    "database/sql"
    "fmt"
    "log"
    "time"

    _ "github.com/go-sql-driver/mysql"
)

func main() {
    dsn := "user:password@tcp(localhost:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local"
    
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // 设置最大打开连接数
    db.SetMaxOpenConns(100)

    // 设置最大空闲连接数
    db.SetMaxIdleConns(10)

    // 设置连接最大生命周期
    db.SetConnMaxLifetime(1 * time.Hour)

    // 获取连接池状态
    stats := db.Stats()
    fmt.Printf("打开连接数: %d\n", stats.OpenConnections)
    fmt.Printf("空闲连接数: %d\n", stats.Idle)
    fmt.Printf("正在使用连接数: %d\n", stats.InUse)
    fmt.Printf("等待获取连接的 goroutine 数: %d\n", stats.WaitCount)
    fmt.Printf("连接被获取前等待时间超过 1 秒的次数: %d\n", stats.WaitDuration)
}
```

### 5.2 使用连接池的最佳实践

```go
package main

import (
    "database/sql"
    "fmt"
    "log"
    "sync"
    "time"

    _ "github.com/go-sql-driver/mysql"
)

// DBPool 数据库连接池包装
type DBPool struct {
    db *sql.DB
    mu sync.RWMutex
}

func NewDBPool(dsn string, maxOpen, maxIdle int, maxLifetime time.Duration) (*DBPool, error) {
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        return nil, err
    }

    db.SetMaxOpenConns(maxOpen)
    db.SetMaxIdleConns(maxIdle)
    db.SetConnMaxLifetime(maxLifetime)

    // 验证连接
    if err := db.Ping(); err != nil {
        return nil, err
    }

    return &DBPool{db: db}, nil
}

// GetDB 获取底层 DB
func (p *DBPool) GetDB() *sql.DB {
    p.mu.RLock()
    defer p.mu.RUnlock()
    return p.db
}

// Close 关闭连接池
func (p *DBPool) Close() error {
    p.mu.Lock()
    defer p.mu.Unlock()
    return p.db.Close()
}

// HealthCheck 健康检查
func (p *DBPool) HealthCheck() error {
    p.mu.RLock()
    defer p.mu.RUnlock()
    return p.db.Ping()
}

func main() {
    pool, err := NewDBPool(
        "user:password@tcp(localhost:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local",
        100,  // maxOpen
        20,   // maxIdle
        1*time.Hour, // maxLifetime
    )
    if err != nil {
        log.Fatal(err)
    }
    defer pool.Close()

    // 使用连接池
    db := pool.GetDB()
    var count int
    db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
    fmt.Printf("用户总数: %d\n", count)
}
```

---

## 6. 最佳实践

### 6.1 错误处理

```go
package main

import (
    "errors"
    "log"

    "gorm.io/gorm"
)

// IsRecordNotFoundError 检查是否是记录不存在错误
func IsRecordNotFoundError(err error) bool {
    return errors.Is(err, gorm.ErrRecordNotFound)
}

// HandleError 处理数据库错误
func HandleError(err error, operation string) {
    if err != nil {
        if IsRecordNotFoundError(err) {
            log.Printf("%s: 记录不存在", operation)
        } else {
            log.Printf("%s 失败: %v", operation, err)
        }
    }
}
```

### 6.2 SQL 注入防护

```go
package main

import (
    "database/sql"
    "fmt"
    "log"

    _ "github.com/go-sql-driver/mysql"
)

// 危险！不要这样写
func BadQuery(db *sql.DB, username string) {
    // SQL 注入漏洞
    query := fmt.Sprintf("SELECT * FROM users WHERE username = '%s'", username)
    db.QueryRow(query)
}

// 正确做法：使用参数化查询
func GoodQuery(db *sql.DB, username string) {
    // 使用占位符
    db.QueryRow("SELECT * FROM users WHERE username = ?", username)
}

// 使用预处理语句
func PreparedQuery(db *sql.DB, username string) {
    stmt, err := db.Prepare("SELECT * FROM users WHERE username = ?")
    if err != nil {
        log.Fatal(err)
    }
    defer stmt.Close()
    stmt.QueryRow(username)
}
```

### 6.3 数据库迁移

```go
package main

import (
    "log"

    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
    // 自动迁移（根据结构体创建/更新表）
    return db.AutoMigrate(
        &User{},
        &Order{},
        &Product{},
        &Category{},
        &Tag{},
    )
}

func ManualMigrate(db *gorm.DB) error {
    // 手动迁移
    // 创建表
    // db.Exec("CREATE TABLE IF NOT EXISTS users (...)")

    // 添加索引
    // db.Exec("CREATE INDEX idx_users_email ON users(email)")

    // 添加字段
    // db.Exec("ALTER TABLE users ADD COLUMN phone VARCHAR(20)")

    return nil
}

func main() {
    dsn := "user:password@tcp(localhost:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local"
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal(err)
    }

    if err := AutoMigrate(db); err != nil {
        log.Fatal("迁移失败:", err)
    }
    log.Println("迁移成功")
}
```

### 6.4 连接配置建议

```go
// MySQL 连接配置建议
const (
    // 连接池配置
    MaxOpenConns    = 100  // 根据服务器配置调整
    MaxIdleConns    = 10   // 通常设为 CPU 核心数的 2-4 倍
    ConnMaxLifetime = 1 * time.Hour
    
    // 超时配置
    ConnectTimeout = 10 * time.Second
    ReadTimeout    = 30 * time.Second
    WriteTimeout   = 30 * time.Second
)

// 连接字符串
dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%v&loc=Local&timeout=%s",
    config.User, config.Password, config.Host, config.Port, 
    config.Database, config.Charset, true, ConnectTimeout)
```

---

## 常见问题

1. **连接池耗尽**：增加 `MaxOpenConns` 或减少单个请求的连接占用时间
2. **死锁**：避免在事务中执行长时间操作
3. **SQL 注入**：始终使用参数化查询
4. **性能问题**：使用索引、优化查询、必要时使用缓存
5. **连接超时**：检查网络和数据库服务器配置
