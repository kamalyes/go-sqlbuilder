<!--
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-23 22:53:25
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-23 22:55:12
 * @FilePath: \go-sqlbuilder\README.md
 * @Description: 
 * 
 * Copyright (c) 2025 by kamalyes, All Rights Reserved. 
-->
# go-sqlbuilder

[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/kamalyes/go-sqlbuilder)](https://github.com/kamalyes/go-sqlbuilder)
[![GoDoc](https://godoc.org/github.com/kamalyes/go-sqlbuilder?status.svg)](https://godoc.org/github.com/kamalyes/go-sqlbuilder)
[![License](https://img.shields.io/github/license/kamalyes/go-sqlbuilder)](https://github.com/kamalyes/go-sqlbuilder/blob/main/LICENSE)

一个功能丰富、高性能的 Go 语言 GORM 仓储层封装库，提供类型安全的 CRUD 操作、复杂查询构建和便利方法。

## ✨ 特性

- 🚀 **仓储模式**：泛型 BaseRepository 和 EnhancedRepository，类型安全的 CRUD 操作
- 🔍 **高级查询**：FilterGroup 支持复杂的 AND/OR 条件组合和无限嵌套
- 🎯 **类型安全**：完全的泛型支持，编译时类型检查
- 📊 **性能优化**：批量操作、游标分页、原子字段更新
- 🔐 **错误处理**：集成 go-toolbox/errorx 的结构化错误管理
- 📝 **审计追踪**：内置审计字段（created_by, updated_by）
- 🛠️ **便利方法**：常用操作的快捷方法

## 📦 安装

```bash
go get github.com/kamalyes/go-sqlbuilder
```

## 🚀 快速开始

```go
package main

import (
    "context"
    "log"
    
    sqlbuilder "github.com/kamalyes/go-sqlbuilder"
    "github.com/kamalyes/go-logger"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

type User struct {
    sqlbuilder.BaseModel
    Name   string `gorm:"type:varchar(100)"`
    Email  string `gorm:"type:varchar(100);uniqueIndex"`
    Age    int    `gorm:"type:int"`
    Status string `gorm:"type:varchar(20)"`
}

func main() {
    // 1. 连接数据库
    dsn := "user:password@tcp(127.0.0.1:3306)/testdb?charset=utf8mb4&parseTime=True"
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal(err)
    }
    
    // 2. 创建 Handler
    handler, err := sqlbuilder.NewGormHandler(db)
    if err != nil {
        log.Fatal(err)
    }
    
    // 3. 创建 Repository
    logger := logger.NewLogger(nil)
    repo := sqlbuilder.NewBaseRepository[User](handler, logger, "users")
    
    ctx := context.Background()
    
    // 4. 使用 Repository
    
    // 创建
    user := &User{Name: "张三", Email: "zhangsan@example.com", Age: 25}
    created, err := repo.Create(ctx, user)
    
    // 查询
    user, err = repo.Get(ctx, 1)
    
    // 更新
    user.Age = 26
    updated, err := repo.Update(ctx, user)
    
    // 删除
    err = repo.Delete(ctx, 1)
}
```

## 📚 核心功能

### BaseRepository - 完整 CRUD

提供所有基础数据库操作：

```go
repo := sqlbuilder.NewBaseRepository[User](handler, logger, "users")

// 创建操作
user, err := repo.Create(ctx, &User{Name: "Alice"})
err = repo.CreateBatch(ctx, users...)
created, isNew, err := repo.CreateIfNotExists(ctx, user, "email")

// 查询操作
user, err := repo.Get(ctx, 1)
users, err := repo.GetAll(ctx)
users, paging, err := repo.ListWithPagination(ctx, query, paging)

// 更新操作
user, err := repo.Update(ctx, user)
err = repo.UpdateFields(ctx, 1, map[string]interface{}{"age": 30})

// 删除操作
err = repo.Delete(ctx, 1)
err = repo.SoftDelete(ctx, 1, "deleted_at", time.Now())

// 统计操作
count, err := repo.Count(ctx)
exists, err := repo.Exists(ctx, filter)
```

### EnhancedRepository - 便利方法

扩展方法提供更便捷的操作：

```go
enhanced := sqlbuilder.NewEnhancedRepository[User](handler, logger, "users")

// 字段查询
users, err := enhanced.FindByField(ctx, "status", "active")

// 游标分页（大数据量性能更好）
users, cursor, err := enhanced.FindByFieldWithCursor(ctx, "status", "active", "", 20, "id", "ASC")

// 原子字段操作
err = enhanced.IncrementField(ctx, 1, "points", 10)      // 积分 +10
err = enhanced.DecrementField(ctx, 1, "stock", 5)        // 库存 -5
err = enhanced.ToggleField(ctx, 1, "is_active")          // 切换布尔值
```

### FilterGroup - 复杂查询

支持任意复杂的 AND/OR 条件组合：

```go
// 查询条件：(status = 'active' OR status = 'trial') AND age > 18
statusGroup := sqlbuilder.NewFilterGroup(sqlbuilder.OR).
    AddFilter(sqlbuilder.NewEqFilter("status", "active")).
    AddFilter(sqlbuilder.NewEqFilter("status", "trial"))

mainGroup := sqlbuilder.NewFilterGroup(sqlbuilder.AND).
    AddGroup(statusGroup).
    AddFilter(sqlbuilder.NewGtFilter("age", 18))

query := sqlbuilder.NewQuery().SetFilterGroup(mainGroup)
users, err := repo.List(ctx, query)
```

### 丰富的过滤器

```go
// 比较操作
NewEqFilter("status", "active")              // =
NewGtFilter("age", 18)                       // >
NewBetweenFilter("price", 100, 1000)         // BETWEEN

// 范围操作
NewInFilter("id", []interface{}{1, 2, 3})    // IN
NewNotInFilter("status", []interface{}{"deleted", "banned"})

// 模糊查询
NewLikeFilter("name", "张")                  // LIKE '%张%'

// NULL 检查
NewIsNullFilter("deleted_at")                // IS NULL
NewIsNotNullFilter("verified_at")            // IS NOT NULL

// 时间范围
NewTodayFilter("created_at")                 // 今天
NewThisWeekFilter("created_at")              // 本周
NewLastMonthFilter("created_at")             // 上月

// 自定义 SQL
NewCustomFilter("YEAR(created_at) = ?", 2024)
```

## 📖 完整文档

我们提供了详细的模块化文档：

- 📘 [快速开始](docs/QUICKSTART.MD) - 5 分钟上手指南
- 📗 [Repository 基础](docs/REPOSITORY-BASICS.MD) - 完整 CRUD 操作详解
- 📙 [高级查询](docs/ADVANCED-QUERIES.MD) - Query、Filter 和复杂查询
- 📕 [FilterGroup 完整指南](docs/FILTERGROUP.MD) - 复杂条件组合和嵌套
- 📓 [EnhancedRepository](docs/ENHANCED-REPOSITORY.MD) - 便利方法详解
- 📔 [模型定义](docs/MODELS.MD) - BaseModel、AuditModel 使用
- 📒 [错误处理](docs/ERROR-HANDLING.MD) - 错误管理和日志记录
- 🔄 [Context 使用指南](docs/CONTEXT-USAGE.MD) - 上下文、超时控制、日志追踪

## 🏗️ 内置模型

### BaseModel

包含基础字段的模型：

```go
type User struct {
    sqlbuilder.BaseModel  // ID, CreatedAt, UpdatedAt, DeletedAt
    Name  string
    Email string
}
```

### AuditModel

包含审计字段的模型：

```go
type Article struct {
    sqlbuilder.AuditModel  // BaseModel + CreatedBy, UpdatedBy
    Title   string
    Content string
}
```

## ⚙️ 配置选项

```go
repo := sqlbuilder.NewBaseRepository[User](
    handler,
    logger,
    "users",
    sqlbuilder.WithBatchSize[User](200),              // 批处理大小
    sqlbuilder.WithTimeout[User](60),                 // 超时时间
    sqlbuilder.WithDefaultPreloads[User]("Profile"),  // 默认预加载
    sqlbuilder.WithDefaultOrder[User]("id DESC"),     // 默认排序
)
```

## 🧪 测试

```bash
# 运行所有测试
go test ./... -v

# 测试覆盖率
go test ./... -cover

# 运行特定测试
go test -v -run TestBaseRepository
```

## 📦 依赖

- [GORM](https://gorm.io/) - 数据库 ORM
- [go-toolbox](https://github.com/kamalyes/go-toolbox) - 错误处理（errorx）和工具（mathx）
- [go-logger](https://github.com/kamalyes/go-logger) - 结构化日志


## 🤝 贡献

欢迎贡献！请随时提交 Pull Request。

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 开启 Pull Request

## 📄 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件

## 👨‍💻 作者

Kamal Yang ([@kamalyes](https://github.com/kamalyes))

## 🙏 致谢

- GORM 团队提供的优秀 ORM
- Go 社区的灵感和支持
