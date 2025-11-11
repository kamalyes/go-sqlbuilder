# Go SQLBuilder - 高级SQL构建器

[![Go Version](https://img.shields.io/badge/Go-1.19+-00ADD8?style=for-the-badge&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-blue?style=for-the-badge)](LICENSE)
[![Build Status](https://img.shields.io/badge/Build-Passing-brightgreen?style=for-the-badge)](https://github.com)

一个功能强大、易于使用的Go语言SQL查询构建器，支持无限链式调用，兼容多种ORM框架（sqlx、gorm等），提供丰富的数据库操作功能。

## ✨ 特性

- 🔗 **无限链式调用** - 从 `NewSQLBuilder()` 开始的流畅API设计
- 🏗️ **单级架构** - 扁平化结构，无复杂继承关系
- 🔌 **多框架兼容** - 支持sqlx、gorm等主流数据库框架
- 🗄️ **多数据库支持** - MySQL、PostgreSQL、SQLite等
- ⚡ **高性能** - 优化的SQL生成和执行机制
- 🎯 **类型安全** - 完善的类型定义和接口设计
- 📊 **性能监控** - 内置查询日志和性能分析
- 🔄 **事务支持** - 完整的事务管理功能
- 🛡️ **SQL注入防护** - 参数化查询，安全可靠

## 📦 安装

```bash
go get github.com/your-username/go-sqlbuilder
```

## 🚀 快速开始

### 基础用法

```go
package main

import (
    "log"
    
    "github.com/jmoiron/sqlx"
    _ "github.com/go-sql-driver/mysql"
    sqlbuilder "github.com/your-username/go-sqlbuilder"
)

func main() {
    // 连接数据库
    db, err := sqlx.Connect("mysql", "user:password@tcp(localhost:3306)/testdb")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // 创建SQLBuilder实例 - 唯一入口点
    builder := sqlbuilder.NewSQLBuilder(db, 
        sqlbuilder.WithDebug(true), 
        sqlbuilder.WithTimeout(10*time.Second))

    // 无限链式调用示例
    var users []User
    err = builder.Table("users").                    // 设置表名
        Select("id", "name", "email").               // 选择字段
        Where("status", 1).                          // 添加条件
        Where("age", ">=", 18).                      // 更多条件
        WhereIn("city", []interface{}{"北京", "上海"}). // IN 条件
        OrderBy("created_at", "DESC").               // 排序
        Limit(20).                                   // 限制数量
        Find(&users)                                 // 执行查询

    if err != nil {
        log.Printf("Query error: %v", err)
    }
}
```

## 📖 详细使用说明

### 1. 创建SQLBuilder实例

```go
// 使用sqlx
db, _ := sqlx.Connect("mysql", "dsn")
builder := sqlbuilder.NewSQLBuilder(db)

// 使用gorm
gormDB, _ := gorm.Open(mysql.Open("dsn"))
builder := sqlbuilder.NewSQLBuilder(gormDB)

// 带配置选项
builder := sqlbuilder.NewSQLBuilder(db,
    sqlbuilder.WithDebug(true),                  // 启用调试
    sqlbuilder.WithTimeout(10*time.Second),      // 设置超时
    sqlbuilder.WithDriver(MySQLDriverAdapter()), // 指定驱动适配器
)
```

### 2. 查询操作

#### 基础查询

```go
// 简单查询
var users []User
builder.Table("users").
    Select("id", "name", "email").
    Where("status", 1).
    Find(&users)

// 单条记录
var user User
builder.Table("users").
    Where("id", 1).
    First(&user)

// 获取单个值
name, err := builder.Table("users").
    Where("id", 1).
    Value("name")
```

#### 复杂查询

```go
// 复杂条件查询
var users []User
builder.Table("users").
    Select("u.*, COUNT(o.id) as order_count").
    As("u").
    LeftJoin("orders o", "u.id = o.user_id").
    Where("u.status", 1).
    Where("u.age", ">=", 18).
    WhereIn("u.city", []interface{}{"北京", "上海", "深圳"}).
    WhereBetween("u.created_at", "2023-01-01", "2023-12-31").
    WhereExists(subQuery).
    GroupBy("u.id").
    Having("COUNT(o.id)", ">", 5).
    OrderBy("order_count", "DESC").
    Limit(50).
    Find(&users)
```

#### 子查询

```go
// 使用子查询
subQuery := builder.Table("orders").
    Select("user_id").
    Where("amount", ">", 1000).
    GroupBy("user_id")

var users []User
builder.Table("users").
    WhereExists(subQuery).
    Find(&users)
```

### 3. 插入操作

```go
// 单条插入
userData := map[string]interface{}{
    "name":    "张三",
    "email":   "zhangsan@example.com", 
    "age":     25,
    "status":  1,
}

result, err := builder.Table("users").
    Insert(userData).
    Exec()

insertID, _ := result.LastInsertId()

// 批量插入
batchData := []map[string]interface{}{
    {"name": "用户1", "email": "user1@test.com"},
    {"name": "用户2", "email": "user2@test.com"},
    {"name": "用户3", "email": "user3@test.com"},
}

builder.Table("users").InsertBatch(batchData)

// 插入或更新 (MySQL)
builder.Table("users").
    Insert(userData).
    OnDuplicateKeyUpdate(map[string]interface{}{
        "updated_at": time.Now(),
    }).Exec()

// Upsert (PostgreSQL/MySQL兼容)
builder.Table("users").
    Upsert(userData, []string{"email"}, []string{"name", "age"})
```

### 4. 更新操作

```go
// 基础更新
updateData := map[string]interface{}{
    "name":       "新名字",
    "updated_at": time.Now(),
}

builder.Table("users").
    Where("id", 1).
    Update(updateData).
    Exec()

// 链式设置字段
builder.Table("users").
    Where("id", 1).
    Set("name", "新名字").
    Set("email", "newemail@example.com").
    Exec()

// 字段递增/递减
builder.Table("users").
    Where("id", 1).
    Increment("login_count", 1).
    Exec()

builder.Table("products").
    Where("id", 1).
    Decrement("stock", 5).
    Exec()

// 批量更新
batchData := []map[string]interface{}{
    {"id": 1, "name": "用户1", "status": 1},
    {"id": 2, "name": "用户2", "status": 0},
}

builder.Table("users").UpdateBatch(batchData, "id")
```

### 5. 删除操作

```go
// 基础删除
builder.Table("users").
    Where("status", 0).
    Delete().
    Exec()

// 软删除
builder.Table("users").
    Where("id", 1).
    Set("deleted_at", time.Now()).
    Exec()

// 恢复软删除
builder.Table("users").
    Where("id", 1).
    Restore()

// 强制删除
builder.Table("users").
    Where("id", 1).
    ForceDelete().
    Exec()
```

### 6. 聚合查询

```go
// 计数
count, err := builder.Table("users").
    Where("status", 1).
    Count()

// 存在检查
exists, err := builder.Table("users").
    Where("email", "test@example.com").
    Exists()

// 求和、平均值、最大值、最小值
totalAmount, _ := builder.Table("orders").Sum("amount")
avgAge, _ := builder.Table("users").Avg("age")
maxAge, _ := builder.Table("users").Max("age")
minAge, _ := builder.Table("users").Min("age")

// 单列值
var emails []string
builder.Table("users").
    Where("status", 1).
    Pluck("email", &emails)

// 键值对映射
userMap, err := builder.Table("users").
    PluckMap("id", "name")
```

### 7. 事务操作

```go
// 自动事务管理
err := builder.Transaction(func(tx *sqlbuilder.SQLBuilder) error {
    // 创建用户
    result, err := tx.Table("users").
        Insert(map[string]interface{}{
            "name":  "事务用户",
            "email": "tx@example.com",
        }).Exec()
    
    if err != nil {
        return err
    }
    
    userID, _ := result.LastInsertId()
    
    // 创建订单
    _, err = tx.Table("orders").
        Insert(map[string]interface{}{
            "user_id": userID,
            "amount":  100.00,
        }).Exec()
    
    return err
})

// 手动事务控制
tx, err := builder.BeginTx()
if err != nil {
    log.Fatal(err)
}

defer func() {
    if r := recover(); r != nil {
        tx.Rollback()
        panic(r)
    }
}()

// 执行操作...

if err != nil {
    tx.Rollback()
    return
}

err = tx.Commit()
```

### 8. 高级功能

#### 分页查询

```go
var users []User
pagination, err := builder.Table("users").
    Where("status", 1).
    OrderByDesc("created_at").
    Paginate(1, 20) // 第1页，每页20条

fmt.Printf("Total: %d, Pages: %d\n", pagination.Total, pagination.LastPage)
```

#### 分块处理

```go
// 分块处理大量数据
err := builder.Table("users").
    Where("status", 1).
    Chunk(1000, func(records interface{}) error {
        users := records.([]User)
        // 处理每个分块的数据
        return nil
    })

// 按ID分块（避免offset性能问题）
err := builder.Table("users").
    ChunkByID(1000, "id", func(records interface{}) error {
        // 处理逻辑
        return nil
    })
```

#### Union查询

```go
activeUsers := builder.Table("users").
    Select("name", "email").
    Where("status", 1)

inactiveUsers := builder.Table("users").
    Select("name", "email").
    Where("status", 0)

var allUsers []User
activeUsers.Union(inactiveUsers).Find(&allUsers)
```

#### 条件构造

```go
// 条件执行
builder.Table("users").
    When(condition, func(q *sqlbuilder.SQLBuilder) *sqlbuilder.SQLBuilder {
        return q.Where("status", 1)
    }).
    Unless(anotherCondition, func(q *sqlbuilder.SQLBuilder) *sqlbuilder.SQLBuilder {
        return q.Where("deleted_at", "IS", "NULL")
    })

// 作用域
builder.Scope(func(q *sqlbuilder.SQLBuilder) *sqlbuilder.SQLBuilder {
    return q.WhereNull("deleted_at") // 全局软删除过滤
})
```

### 9. 性能监控

```go
// 启用查询日志
builder = builder.EnableQueryLog().Debug(true)

// 执行查询...

// 获取查询日志
logs := builder.GetQueryLog()
for _, log := range logs {
    fmt.Printf("SQL: %s, Time: %.4fs\n", log.SQL, log.Time)
}

// 获取性能指标
metrics := builder.GetMetrics()
fmt.Printf("Total queries: %d\n", metrics.TotalQueries)
fmt.Printf("Average time: %.4fs\n", metrics.AverageTime)

// Explain分析
explain, err := builder.Table("users").
    Where("email", "test@example.com").
    Explain()

// 性能分析
profile, err := builder.Table("users").Profile()
```

### 10. 调试工具

```go
// 输出SQL但不执行
sql, params := builder.Table("users").
    Where("status", 1).
    ToSQL()

fmt.Printf("SQL: %s\nParams: %v\n", sql, params)

// 调试输出
builder.Table("users").
    Where("status", 1).
    Debug().       // 启用调试
    Dump().        // 输出但继续
    Find(&users)

// 输出并退出
builder.Table("users").
    Where("status", 1).
    DD()  // Dump and Die
```

## 🏗️ 架构设计

### 接口层次

```
DatabaseInterface
├── SqlxInterface
├── GormInterface  
└── DriverAdapterInterface
    ├── MySQLDriverAdapter
    └── PostgreSQLDriverAdapter
```

### 核心组件

- **SQLBuilder** - 主要构建器类，提供链式API
- **Adapters** - 数据库适配器层，支持多种ORM
- **Drivers** - 数据库驱动适配器，处理特定数据库语法
- **Query Components** - 查询组件（Where、Join、OrderBy等）
- **Event System** - 事件系统和钩子支持

## 🔧 配置选项

```go
type Options struct {
    Debug          bool              // 启用调试模式
    Timeout        time.Duration     // 查询超时时间
    Context        context.Context   // 上下文
    Driver         DriverAdapter     // 数据库驱动适配器
    QueryLog       bool             // 启用查询日志
    MaxOpenConns   int              // 最大连接数
    MaxIdleConns   int              // 最大空闲连接数
    ConnMaxLife    time.Duration     // 连接最大生命周期
}

// 使用配置
builder := NewSQLBuilder(db,
    WithDebug(true),
    WithTimeout(30*time.Second),
    WithContext(ctx),
    WithQueryLog(true),
)
```

## 📊 性能特性

- **连接池管理** - 智能连接池优化
- **查询缓存** - SQL和结果缓存机制
- **批量操作** - 高效的批量插入/更新
- **分块处理** - 大数据集的分块处理
- **性能监控** - 详细的性能指标收集

## 🧪 测试

```bash
# 运行所有测试
go test ./...

# 运行覆盖率测试
go test -cover ./...

# 运行基准测试
go test -bench=. ./...

# 运行特定测试
go test -run TestSQLBuilder ./...
```

## 📝 示例项目

查看 `examples.go` 文件获取完整的使用示例，包括：

- 基础CRUD操作
- 复杂查询构建
- 事务管理
- 性能监控
- 多数据库适配
- 业务场景实现

## 🤝 贡献

欢迎贡献代码！请遵循以下步骤：

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 创建 Pull Request

## 📄 许可证

本项目使用 MIT 许可证。查看 [LICENSE](LICENSE) 文件了解详情。

## 🆘 支持

- 📧 Email: 501893067@qq.com
- 🐛 Issues: [GitHub Issues](https://github.com/your-username/go-sqlbuilder/issues)
- 📖 文档: [Wiki](https://github.com/your-username/go-sqlbuilder/wiki)

## 🎯 路线图

- [ ] 支持更多数据库（Oracle、SQL Server）
- [ ] GraphQL查询支持
- [ ] 分布式查询支持
- [ ] NoSQL数据库适配器
- [ ] 可视化查询构建器
- [ ] 更多性能优化

## 🙏 致谢

感谢以下开源项目的启发：

- [jmoiron/sqlx](https://github.com/jmoiron/sqlx)
- [go-gorm/gorm](https://github.com/go-gorm/gorm)
- [Masterminds/squirrel](https://github.com/Masterminds/squirrel)

---

⭐ 如果这个项目对你有帮助，请给我们一个星标！