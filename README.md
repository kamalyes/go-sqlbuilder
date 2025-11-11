# Go SQLBuilder - 高级SQL构建器 v2.0

[![Go Version](https://img.shields.io/badge/Go-1.19+-00ADD8?style=for-the-badge&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-blue?style=for-the-badge)](LICENSE)
[![Build Status](https://img.shields.io/badge/Build-Passing-brightgreen?style=for-the-badge)](https://github.com)
[![Tests](https://img.shields.io/badge/Tests-50%2B%20passing-brightgreen?style=for-the-badge)](https://github.com)

一个**生产级别**的SQL查询构建器，提供：

- 🔗 **无限链式调用** - 流畅的API设计
- � **模块化架构** - 独立的cache、query、errors包
- ⚡ **自动缓存** - Redis集成，自动TTL管理
- �️ **完整错误处理** - 48种标准错误码
- 📊 **全面测试** - 50+单元测试，100%通过率

## ✨ 核心特性

### Builder（SQL构建）

- 🔗 无限链式调用
- 📝 SELECT/INSERT/UPDATE/DELETE
- 🔀 JOIN、GROUP BY、HAVING、ORDER BY
- 🔄 事务支持
- 🛡️ 参数化查询（SQL注入防护）

### Cache（缓存管理）

- 💾 Redis集成
- ⏱️ 自动TTL失效
- 📈 命中率统计
- 🧪 完整的Mock实现

### Query（高级查询）

- � 20+便捷方法
- 🔍 灵活的过滤条件
- �📊 分页和排序
- 🎯 WHERE子句自动生成

### Errors（错误处理）

- 📋 48种标准错误码
- 📝 String()和Error()接口
- 🎯 错误分类管理

## � 文档速览

| 文档 | 说明 |
|------|------|
| [项目分析](PROJECT_ANALYSIS.md) | 完整的项目架构和功能分析 |
| [架构设计](ARCHITECTURE.md) | 设计模式、数据流、扩展点 |
| [使用指南](USAGE_GUIDE.md) | 详细的使用示例和最佳实践 |
| [高级查询](ADVANCED_QUERY_USAGE.md) | 20+便捷方法详解 |

## 📦 安装

```bash
go get github.com/kamalyes/go-sqlbuilder
```

## 🚀 快速开始

### 基础使用

```go
package main

import (
    "log"
    
    "github.com/jmoiron/sqlx"
    _ "github.com/go-sql-driver/mysql"
    sqlbuilder "github.com/kamalyes/go-sqlbuilder"
)

func main() {
    // 连接数据库
    db, err := sqlx.Connect("mysql", "user:password@tcp(localhost:3306)/testdb")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // 创建Builder实例
    builder := sqlbuilder.New(db)

    // 执行查询
    var users []User
    err = builder.Table("users").
        Select("id", "name", "email").
        Where("status", 1).
        Where("age", ">", 18).
        OrderBy("created_at", "DESC").
        Limit(10).
        Find(&users)
}
```

### 带缓存的查询

```go
import "github.com/kamalyes/go-sqlbuilder/cache"

// 创建缓存store
store := cache.NewMockStore()  // 或 cache.NewRedisStore(redisClient, "prefix:")

// 创建带缓存的Builder
cachedBuilder, _ := sqlbuilder.NewCachedBuilder(db, store, nil)

// 自动缓存查询结果
result, _ := cachedBuilder.GetCached(ctx, sql, args...)
```

### 高级查询参数

```go
import "github.com/kamalyes/go-sqlbuilder/query"

param := query.NewParam().
    AddEQ("status", 1).
    AddGT("age", 18).
    AddLike("name", "John").
    AddIn("category", 1, 2, 3).
    AddOrder("created_at", "DESC").
    SetPage(1, 20)

whereSQL, args := param.BuildWhereClause()
```

## 📖 详细使用说明

> 📚 **更多使用示例请查看** [使用指南](USAGE_GUIDE.md)
> 📊 **了解架构设计请查看** [架构设计](ARCHITECTURE.md)
> 🔍 **查看技术分析请查看** [项目分析](PROJECT_ANALYSIS.md)

### Builder实例

```go
// 连接数据库
db, err := sqlx.Connect("mysql", "user:password@tcp(host:3306)/dbname")
defer db.Close()

// 创建Builder
builder := sqlbuilder.New(db)

// 或使用GORM
import "gorm.io/gorm"
gormDB, err := gorm.Open(mysql.Open(dsn))
builder := sqlbuilder.New(gormDB)
```

### 查询 (SELECT)

```go
var users []User
builder.Table("users").
    Select("id", "name", "email").
    Where("status", 1).
    Where("age", ">", 18).
    OrderBy("created_at", "DESC").
    Limit(10).
    Find(&users)

// 单条记录
var user User
builder.Table("users").Where("id", 1).First(&user)

// 获取单个值
name, _ := builder.Table("users").Where("id", 1).Value("name")

// 计数
count, _ := builder.Table("users").Where("status", 1).Count()

// 存在检查
exists, _ := builder.Table("users").Where("id", 1).Exists()
```

### 插入 (INSERT)

```go
result, err := builder.Table("users").
    Insert(map[string]interface{}{
        "name":   "张三",
        "email":  "zhangsan@example.com",
        "status": 1,
    }).
    Exec()

id, _ := result.LastInsertId()

// 批量插入
data := []map[string]interface{}{
    {"name": "用户1", "email": "user1@test.com"},
    {"name": "用户2", "email": "user2@test.com"},
}
builder.Table("users").InsertBatch(data)
```

### 更新 (UPDATE)

```go
builder.Table("users").
    Where("id", 1).
    Update(map[string]interface{}{
        "name":       "新名字",
        "updated_at": time.Now(),
    }).
    Exec()

// 链式调用
builder.Table("users").
    Where("id", 1).
    Set("name", "新名字").
    Set("email", "new@example.com").
    Exec()

// 字段递增/递减
builder.Table("users").Where("id", 1).Increment("login_count", 1)
builder.Table("products").Where("id", 1).Decrement("stock", 5)
```

### 删除 (DELETE)

```go
builder.Table("users").
    Where("status", 0).
    Delete().
    Exec()

// 软删除
builder.Table("users").
    Where("id", 1).
    Set("deleted_at", time.Now()).
    Exec()
```

### 事务支持

```go
err := builder.Transaction(func(tx *sqlbuilder.SQLBuilder) error {
    // 创建用户
    result, _ := tx.Table("users").Insert(userData).Exec()
    
    userID, _ := result.LastInsertId()
    
    // 创建订单
    _, err := tx.Table("orders").Insert(orderData).Exec()
    
    return err
})
```

## 🔗 高级特性

### 缓存管理

```go
import "github.com/kamalyes/go-sqlbuilder/cache"

// Redis缓存
store := cache.NewRedisConfig("localhost:6379").
    WithPrefix("myapp:").
    Build()

cachedBuilder, _ := sqlbuilder.NewCachedBuilder(db, store, nil)

// 获取带缓存的结果
result, _ := cachedBuilder.GetCached(ctx, sql, args...)

// Mock测试用缓存
mockStore := cache.NewMockStore()
```

### 错误处理

```go
import "github.com/kamalyes/go-sqlbuilder/errors"

err := builder.Table("users").Where("id", 1).First(&user)

// 错误检查
if errors.IsErrorCode(err, errors.ErrorCodeKeyNotFound) {
    log.Println("缓存键未找到")
}

code := errors.GetErrorCode(err)
msg := errors.ErrorCodeString(code)
```

### 高级查询参数

```go
import "github.com/kamalyes/go-sqlbuilder/query"

// 构建复杂查询条件
param := query.NewParam().
    AddEQ("status", 1).              // status = 1
    AddGT("age", 18).                // age > 18
    AddLike("name", "John").         // name LIKE %John%
    AddIn("category", 1, 2, 3).      // category IN (1,2,3)
    AddBetween("price", 10, 100).    // price BETWEEN 10 AND 100
    AddOrder("created_at", "DESC").  // ORDER BY created_at DESC
    SetPage(1, 20)                   // LIMIT 20 OFFSET 0

whereSQL, args := param.BuildWhereClause()

// OR 条件
param2 := query.NewParam().
    AddEQ("role", "admin").
    AddOrEQ("permission_level", 10)
```

## 📊 模块架构

### 核心包结构

```
go-sqlbuilder/
├── builder.go              # SQL构建核心引擎 (670 lines)
├── builder_cached.go       # 缓存包装器 (173 lines)
├── adapters.go             # SQLX/GORM适配器 (1376 lines)
├── interfaces.go           # 接口定义
│
├── cache/                  # 缓存包
│   ├── interface.go        # Store接口
│   ├── config.go           # 配置管理
│   ├── redis.go            # Redis实现
│   ├── mock.go             # 测试Mock
│   └── manager.go          # 统计管理
│
├── query/                  # 查询参数包
│   ├── param.go            # 20+便捷方法
│   ├── filter.go           # 过滤条件
│   ├── operator.go         # 操作符定义
│   ├── pagination.go       # 分页支持
│   └── option.go           # 查询选项
│
└── errors/                 # 错误处理包
    ├── code.go             # 48种错误码
    └── error.go            # 错误结构体
```

## 📈 性能特性

- ⚡ **SQL缓存** - MD5 Cache Key自动生成，支持TTL失效
- 📊 **统计管理** - 缓存命中率、操作计数统计
- 🔄 **连接池** - 底层数据库驱动连接复用
- 🎯 **参数化查询** - 完全防止SQL注入
- 🧪 **完整测试** - 50+单元测试，100%通过率

## 🛠️ 支持的数据库

| 数据库 | 驱动 | 适配器 | 状态 |
|--------|------|--------|------|
| MySQL | github.com/go-sql-driver/mysql | SQLX | ✅ 生产就绪 |
| PostgreSQL | github.com/lib/pq | SQLX | ✅ 生产就绪 |
| SQLite | github.com/mattn/go-sqlite3 | SQLX | ✅ 生产就绪 |
| MySQL | GORM v1 | GORM | ✅ 支持 |
| PostgreSQL | GORM v2 | GORM v2 | ✅ 支持 |

## 🧪 测试

```bash
# 运行所有测试
go test ./... -v

# 运行特定包的测试
go test ./cache -v
go test ./query -v
go test ./errors -v

# 获取覆盖率报告
go test ./... -cover
```

## 🔐 安全特性

- 🛡️ **SQL注入防护** - 所有查询参数化
- 📝 **输入验证** - 严格的参数校验
- 🔒 **事务隔离** - 完善的事务管理
- 📊 **错误日志** - 详细的错误跟踪
- ✅ **类型安全** - 强类型检查

## 📚 文档导航

| 文档 | 描述 | 适合场景 |
|------|------|---------|
| [使用指南](USAGE_GUIDE.md) | 450+行，详细的使用示例 | 快速上手，常见用法 |
| [架构设计](ARCHITECTURE.md) | 350+行，设计模式和数据流 | 深度理解，二次开发 |
| [项目分析](PROJECT_ANALYSIS.md) | 350+行，完整的技术分析 | 全面掌握，参考手册 |
| [高级查询](ADVANCED_QUERY_USAGE.md) | 20+便捷方法详解 | 复杂条件查询 |

## 💡 最佳实践

1. **始终使用参数化查询** - 防止SQL注入
2. **利用事务处理** - 确保数据一致性
3. **合理使用缓存** - 提升查询性能  
4. **监控缓存统计** - 优化缓存策略
5. **错误处理** - 使用标准错误码
6. **批量操作** - 使用InsertBatch/UpdateBatch
7. **分页查询** - 避免一次加载大量数据
8. **建立适当索引** - 提升查询效率

## 🚀 快速链接

- 📖 [完整使用指南](USAGE_GUIDE.md) - 从这里开始
- 🏗️ [系统架构](ARCHITECTURE.md) - 了解项目设计  
- 📊 [项目分析](PROJECT_ANALYSIS.md) - 深入技术细节
- 🔍 [高级查询](ADVANCED_QUERY_USAGE.md) - 掌握便捷方法
- 📦 [Go Modules](go.mod) - 依赖管理
- ✅ [测试覆盖](comprehensive_test.go) - 质量保证

## 🤝 贡献

欢迎贡献代码！请遵循以下步骤：

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 创建 Pull Request

## � 报告问题

发现Bug或有功能建议？请提交Issue：

- 描述问题现象和复现步骤
- 提供最小化的代码示例
- 说明Go版本和数据库类型

## 🆘 支持

- 📧 Email: <501893067@qq.com>
- 🐛 Issues: [GitHub Issues](https://github.com/your-username/go-sqlbuilder/issues)
- 📖 文档: [Wiki](https://github.com/your-username/go-sqlbuilder/wiki)

## 🙏 致谢

感谢以下开源项目的启发：

- [jmoiron/sqlx](https://github.com/jmoiron/sqlx)
- [go-gorm/gorm](https://github.com/go-gorm/gorm)
- [Masterminds/squirrel](https://github.com/Masterminds/squirrel)

---

**最后更新:** 2024年
**版本:** v2.0 - 模块化架构
**测试状态:** 50+ 单元测试，100% 通过率
**生产就绪:** ✅ 完全可用于生产环境

⭐ 如果这个项目对你有帮助，请给我们一个星标！
