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
activeUsers, err := repo.List(ctx, 
    repository.NewQuery().AddEqual("status", "active"))

// 复合查询（便捷链式调用）
query := repository.NewQuery().
    AddEqual("status", "active").
    AddGreaterThan("age", 18).
    AddLike("name", "张").
    AddOrderDesc("created_at").
    Take(10)
    
users, err := repo.List(ctx, query)

// 分页查询
query := repository.NewQuery().
    AddEqual("status", "active").
    AddThisMonth("created_at").  // 本月注册的用户
    AddOrderDesc("created_at").
    Page(1, 20)  // 第1页，每页20条
    
users, pagination, err := repo.ListWithPagination(ctx, query, nil)

// 时间范围查询
lastWeek := time.Now().AddDate(0, 0, -7)
recentUsers, err := repo.List(ctx, 
    repository.NewQuery().
        AddTimeAfter("created_at", lastWeek).
        AddIn("status", "active", "pending"))
```

### 4. 便捷查询方法速览

```go
// 基础条件
query.AddEqual("field", value)           // field = value
query.AddNotEqual("field", value)        // field != value
query.AddLike("field", "keyword")        // field LIKE '%keyword%'
query.AddStartsWith("field", "prefix")   // field LIKE 'prefix%'
query.AddEndsWith("field", "suffix")     // field LIKE '%suffix'

// 范围条件
query.AddIn("field", 1, 2, 3)           // field IN (1,2,3)
query.AddNotIn("field", 1, 2)           // field NOT IN (1,2)
query.AddBetween("field", 1, 100)       // field BETWEEN 1 AND 100
query.AddGreaterThan("field", 10)       // field > 10
query.AddLessEqual("field", 50)         // field <= 50

// 时间条件
query.AddTimeAfter("date", time)        // date > time
query.AddTimeBefore("date", time)       // date < time
query.AddToday("date")                  // 今天
query.AddThisWeek("date")               // 本周
query.AddThisMonth("date")              // 本月

// 排序分页
query.AddOrderAsc("field")              // ORDER BY field ASC
query.AddOrderDesc("field")             // ORDER BY field DESC
query.Page(1, 20)                      // 分页
query.Take(10)                          // LIMIT 10
query.Skip(20)                          // OFFSET 20
```
```

## 下一步

- 📖 [基础操作指南](./REPOSITORY-BASICS.MD) - 学习所有CRUD 方法
- 🔍 [高级查询](./ADVANCED-QUERIES.MD) - 复杂查询和过滤
- 🎯 [FilterGroup 使用](./FILTERGROUP.MD) - 构建复杂 WHERE 条件
- 🏗️ [模型定义](./MODELS.MD) - 使用内置模型

