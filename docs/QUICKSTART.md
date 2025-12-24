# 快速开始

## 安装

```bash
go get github.com/kamalyes/go-sqlbuilder
```

## 基础示例

### 1. 创建 Repository

```go
package main

import (
    "context"
    "log"
    "time"
    
    "github.com/kamalyes/go-logger"
    "github.com/kamalyes/go-sqlbuilder/db"
    "github.com/kamalyes/go-sqlbuilder/repository"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

type User struct {
    ID        uint      `gorm:"primaryKey"`
    Name      string    `gorm:"type:varchar(100)"`
    Email     string    `gorm:"type:varchar(100);uniqueIndex"`
    Age       int       `gorm:"type:int"`
    Status    string    `gorm:"type:varchar(20)"`
    CreatedAt time.Time
}

func main() {
    // 1. 连接数据库
    dsn := "user:password@tcp(127.0.0.1:3306)/testdb?charset=utf8mb4&parseTime=True"
    gormDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal(err)
    }
    
    // 2. 创建 Handler
    handler, err := db.NewGormHandler(gormDB)
    if err != nil {
        log.Fatal(err)
    }
    
    // 3. 创建 Repository
    logger := logger.NewLogger(nil)
    repo := repository.NewBaseRepository[User](handler, logger, "users")
    
    // 4. 使用 Repository
    ctx := context.Background()
    
    // 创建用户
    user := &User{
        Name:   "张三",
        Email:  "zhangsan@example.com",
        Age:    25,
        Status: "active",
    }
    createdUser, err := repo.Create(ctx, user)
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("Created user: %+v", createdUser)
}
```

### 2. CRUD 操作

```go
// 创建
user := &User{Name: "李四", Email: "lisi@example.com", Age: 30}
created, err := repo.Create(ctx, user)

// 读取
user, err := repo.Get(ctx, 1) // 按ID查询

// 更新
user.Age = 31
updated, err := repo.Update(ctx, user)

// 删除
err = repo.Delete(ctx, 1)
```

### 3. 查询操作

```go
// 获取所有用户
users, err := repo.GetAll(ctx)

// 按条件查询（传统方式）
filter := repository.NewEqFilter("status", "active")
activeUsers, err := repo.List(ctx, repository.NewQuery().AddFilter(filter))

// 按条件查询（便捷方式） 🔥 推荐
activeUsers, err := repo.FindWhere(ctx, "status", "active")
```

## 下一步

- 📖 [CRUD 操作](./CRUD-OPERATIONS.md) - 完整的创建、读取、更新、删除方法
- 🚀 [便捷查询方法](./CONVENIENCE-METHODS.md) - 简化的查询 API
- 🔍 [过滤条件](./FILTERS.md) - 构建复杂查询条件
- 📊 [排序和分页](./SORTING-AND-PAGINATION.md) - 数据排序和分页