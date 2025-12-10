# go-sqlbuilder

[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/kamalyes/go-sqlbuilder)](https://github.com/kamalyes/go-sqlbuilder)
[![GoDoc](https://godoc.org/github.com/kamalyes/go-sqlbuilder?status.svg)](https://godoc.org/github.com/kamalyes/go-sqlbuilder)
[![License](https://img.shields.io/github/license/kamalyes/go-sqlbuilder)](https://github.com/kamalyes/go-sqlbuilder/blob/main/LICENSE)

一个功能丰富、高性能的 Go 语言 GORM 仓储层封装库，提供类型安全的 CRUD 操作、复杂查询构建和便利方法。

## 📖 特性 & 文档导航

| 特性 | 说明 | 文档 |
|:-----|:-----|:----:|
| 🚀 **仓储模式** | 泛型 BaseRepository 和 EnhancedRepository，类型安全的 CRUD | [📘 快速开始](docs/QUICKSTART.md) |
| 🔍 **高级查询** | FilterGroup 支持复杂 AND/OR 条件组合和无限嵌套 | [📙 高级查询](docs/ADVANCED-QUERIES.md) |
| 🎯 **类型安全** | 完全的泛型支持，编译时类型检查 | [📗 Repository 基础](docs/REPOSITORY-BASICS.md) |
| ⚡ **自动字段选择** | 基于 struct tags 自动生成查询字段，避免 SELECT * | [⚡ 自动字段选择](docs/AUTO-FIELD-SELECTION.md) |
| 🔗 **链式查询** | 便捷方法支持链式调用构建查询条件 | [🔥 便捷查询示例](docs/QUERY-EXAMPLES.md) |
| 📊 **性能优化** | 批量操作、游标分页、原子字段更新、字段缓存 | [📓 EnhancedRepository](docs/ENHANCED-REPOSITORY.md) |
| 🔐 **错误处理** | 集成 go-toolbox/errorx 的结构化错误管理 | [📒 错误处理](docs/ERROR-HANDLING.md) |
| 📝 **审计追踪** | 内置审计字段（created_by, updated_by） | [📔 模型定义](docs/MODELS.md) |
| 🗄️ **数据库迁移** | 自动迁移、索引创建、表注释 | [🗄️ 数据库迁移器](docs/MIGRATOR.md) |
| 🔄 **上下文支持** | 超时控制、日志追踪、请求隔离 | [🔄 Context 使用指南](docs/CONTEXT-USAGE.md) |
| 🎛️ **复杂条件** | FilterGroup 支持无限嵌套的条件组合 | [📕 FilterGroup 指南](docs/FILTERGROUP.md) |
| 🚄 **并发统计** | 多表并发查询、时间分组统计、条件聚合 | [🚄 并发统计查询](docs/CONCURRENT-STATS.md) |

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
    
    "github.com/kamalyes/go-sqlbuilder/db"
    "github.com/kamalyes/go-sqlbuilder/repository"
    "github.com/kamalyes/go-logger"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

type User struct {
    repository.BaseModel
    Name   string ``gorm:"type:varchar(100)"``
    Email  string ``gorm:"type:varchar(100);uniqueIndex"``
    Age    int    ``gorm:"type:int"``
    Status string ``gorm:"type:varchar(20)"``
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
    
    ctx := context.Background()
    
    // 4. 使用 Repository
    user := &User{Name: "张三", Email: "zhangsan@example.com", Age: 25}
    created, err := repo.Create(ctx, user)
    user, err = repo.Get(ctx, 1)
    user.Age = 26
    updated, err := repo.Update(ctx, user)
    err = repo.Delete(ctx, 1)
}
```

## 📚 核心功能速览

<details>
<summary><b>🔧 BaseRepository - 完整 CRUD</b></summary>

```go
repo := repository.NewBaseRepository[User](handler, logger, "users",
    repository.WithAutoFields[User](),  // 🔥 自动字段选择
)

// 创建
user, err := repo.Create(ctx, &User{Name: "Alice"})
err = repo.CreateBatch(ctx, users...)

// 查询
user, err := repo.Get(ctx, 1)
users, paging, err := repo.ListWithPagination(ctx, query, paging)

// 更新
user, err := repo.Update(ctx, user)
err = repo.UpdateFields(ctx, 1, map[string]interface{}{"age": 30})

// 删除
err = repo.Delete(ctx, 1)
err = repo.SoftDelete(ctx, 1, "deleted_at", time.Now())
```

</details>

<details>
<summary><b>🔗 链式查询构建</b></summary>

```go
query := repository.NewQuery().
    AddEqual("status", 1).
    AddLike("name", "test").
    AddBetween("age", 18, 65).
    AddOrderDesc("created_at").
    Page(1, 20)

users, err := repo.List(ctx, query)
```

| 传统方式 | 便捷方法 |
|---------|----------|
| ``AddFilter(NewEqFilter(...))`` | ``AddEqual("field", value)`` |
| ``AddFilter(NewLikeFilter(...))`` | ``AddLike("field", "test")`` |
| ``AddOrder("field", "DESC")`` | ``AddOrderDesc("field")`` |

</details>

<details>
<summary><b>⚡ EnhancedRepository - 便利方法</b></summary>

```go
enhanced := repository.NewEnhancedRepository[User](handler, logger, "users")

users, err := enhanced.FindByField(ctx, "status", "active")
err = enhanced.IncrementField(ctx, 1, "points", 10)  // 原子 +10
err = enhanced.ToggleField(ctx, 1, "is_active")      // 切换布尔
```

</details>

<details>
<summary><b>🎛️ FilterGroup - 复杂查询</b></summary>

```go
// (status = 'active' OR status = 'trial') AND age > 18
statusGroup := repository.NewFilterGroup(constants.AND_OR).
    AddFilter(repository.NewEqFilter("status", "active")).
    AddFilter(repository.NewEqFilter("status", "trial"))

mainGroup := repository.NewFilterGroup(constants.AND_AND).
    AddGroup(statusGroup).
    AddFilter(repository.NewGtFilter("age", 18))

query := repository.NewQuery().SetFilterGroup(mainGroup)
```

</details>

## 🏗️ 内置模型

```go
// BaseModel - ID, CreatedAt, UpdatedAt, DeletedAt
type User struct {
    repository.BaseModel
    Name string
}

// AuditModel - BaseModel + CreatedBy, UpdatedBy
type Article struct {
    repository.AuditModel
    Title string
}
```

## ⚙️ 配置选项

```go
repo := repository.NewBaseRepository[User](handler, logger, "users",
    repository.WithAutoFields[User](),                // 自动字段选择
    repository.WithBatchSize[User](200),              // 批处理大小
    repository.WithTimeout[User](60),                 // 超时时间
    repository.WithDefaultPreloads[User]("Profile"),  // 默认预加载
    repository.WithDefaultOrder[User]("id DESC"),     // 默认排序
)
```

## 🧪 测试

```bash
go test ./... -v                    # 运行所有测试
go test ./... -cover                # 测试覆盖率
go test -v -run TestBaseRepository  # 运行特定测试
```

## 📦 依赖

- [GORM](https://gorm.io/) - 数据库 ORM
- [go-toolbox](https://github.com/kamalyes/go-toolbox) - 错误处理和工具
- [go-logger](https://github.com/kamalyes/go-logger) - 结构化日志

## 🤝 贡献

1. Fork 本仓库
2. 创建特性分支 (``git checkout -b feature/amazing-feature``)
3. 提交更改 (``git commit -m '✨ feat: Add amazing feature'``)
4. 推送到分支 (``git push origin feature/amazing-feature``)
5. 开启 Pull Request

## 📋 Git Commit Emoji 规范

<details>
<summary>点击展开 Emoji 规范表</summary>

| Emoji | 类型 | 说明 |
|:-----:|------|------|
| ✨ | feat | 新功能 |
| 🐛 | fix | 修复 bug |
| 📝 | docs | 文档更新 |
| ♻️ | refactor | 代码重构 |
| ⚡ | perf | 性能优化 |
| ✅ | test | 测试相关 |
| 🔧 | chore | 配置/构建 |
| 🚀 | deploy | 部署发布 |
| 🔒 | security | 安全修复 |
| 🔥 | remove | 删除代码 |
| 🗃️ | db | 数据库相关 |

**示例：** ``git commit -m "✨ feat(db): 新增 Migrator 迁移器"``

</details>

## 📄 许可证

MIT License - 详见 [LICENSE](LICENSE)

## 👨‍💻 作者

Kamal Yang ([@kamalyes](https://github.com/kamalyes))
