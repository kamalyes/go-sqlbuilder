# Go SQLBuilder Enhanced - 通用数据库适配器

[![Go Version](https://img.shields.io/badge/Go-1.19+-00ADD8?style=for-the-badge&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-blue?style=for-the-badge)](LICENSE)
[![Build Status](https://img.shields.io/badge/Build-Passing-brightgreen?style=for-the-badge)](https://github.com)

一个功能强大的Go语言SQL查询构建器，**支持所有主流数据库框架和ORM**！通过统一的适配器层，无论您使用的是SQLX、GORM、XORM、Beego ORM、Ent还是原生database/sql，都能享受同样流畅的API体验。

## ✨ 核心特性

### 🔌 通用适配器架构
- ✅ **SQLX** - 高性能SQL扩展库
- ✅ **GORM** - 全功能ORM框架  
- 🚧 **XORM** - 简单强大的ORM
- 🚧 **Beego ORM** - 企业级ORM
- 🚧 **Ent** - Facebook的实体框架
- 🚧 **原生database/sql** - Go标准库
- 🚧 **PGX** - PostgreSQL专用驱动
- 🚧 **Bun** - 高性能PostgreSQL ORM
- 🚧 **Squirrel** - SQL查询构建器

### 🗄️ 多数据库支持  
- **MySQL** / **MariaDB**
- **PostgreSQL** 
- **SQLite**
- **Oracle** (计划中)
- **SQL Server** (计划中) 
- **ClickHouse** (计划中)
- **TiDB** (计划中)

### 🔗 统一API设计
- **无限链式调用** - 流畅的API设计
- **自动适配器检测** - 智能识别数据库类型
- **类型安全** - 完善的类型定义
- **上下文支持** - 完整的context.Context集成
- **事务支持** - 统一的事务管理接口

## 📦 安装

```bash
go get github.com/kamalyes/go-sqlbuilder
```

## 🚀 快速开始

### 自动检测适配器

```go
package main

import (
    "log"
    "github.com/jmoiron/sqlx"
    _ "github.com/go-sql-driver/mysql"
    sqlbuilder "github.com/kamalyes/go-sqlbuilder"
)

func main() {
    // 连接数据库 (任意框架)
    db, err := sqlx.Connect("mysql", "user:password@tcp(localhost:3306)/testdb")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // 自动检测并创建适配器
    builder, err := sqlbuilder.NewEnhancedSQLBuilder(db)
    if err != nil {
        log.Fatal(err)
    }
    defer builder.Close()

    // 统一的API - 无论底层使用什么框架
    users := builder.
        Table("users").
        Select("id", "name", "email").
        Where("status = ?", "active").
        Where("age > ?", 18).
        OrderBy("created_at", "DESC").
        Limit(10)

    // 检查适配器信息
    log.Printf("Using: %s (%s)", builder.GetAdapterName(), builder.GetAdapterType())
    log.Printf("Dialect: %s", builder.GetDialect())
}
```

### 指定特定适配器

```go
// 使用SQLX
builderSQLX, err := sqlbuilder.NewWithSQLX(sqlxDB)

// 使用GORM  
builderGORM, err := sqlbuilder.NewWithGORM(gormDB)

// 使用XORM
builderXORM, err := sqlbuilder.NewWithXORM(xormEngine)

// 使用Beego ORM
builderBeego, err := sqlbuilder.NewWithBeegoORM(beegoOrmer)

// 使用Ent
builderEnt, err := sqlbuilder.NewWithEnt(entClient)
```

## 🔧 核心概念

### 适配器模式

我们的架构基于适配器模式，为每个数据库框架提供统一的接口：

```go
// UniversalAdapter - 统一适配器接口
type UniversalAdapter interface {
    // 基础操作
    Query(ctx context.Context, query string, args ...interface{}) (ResultSet, error)
    Exec(ctx context.Context, query string, args ...interface{}) (ExecResult, error)
    
    // 批量操作
    BatchInsert(ctx context.Context, table string, data []map[string]interface{}) error
    BatchUpdate(ctx context.Context, table string, data []map[string]interface{}, whereColumns []string) error
    
    // 事务支持
    BeginTx(ctx context.Context, opts *TxOptions) (Transaction, error)
    
    // 功能检测
    SupportsORM() bool
    SupportsUpsert() bool
    SupportsBulkInsert() bool
}
```

### 自动适配器检测

系统能智能识别您传入的数据库实例类型：

```go
func AutoDetectAdapter(instance interface{}) (UniversalAdapter, error) {
    instanceType := reflect.TypeOf(instance)
    
    // 遍历所有注册的适配器工厂
    for _, factory := range registeredFactories {
        if factory.CanHandle(instanceType) {
            return factory.Create(instance)
        }
    }
    
    return nil, fmt.Errorf("unsupported database type: %T", instance)
}
```

## 💡 使用示例

### 基础查询

```go
// 复杂查询构建
query := builder.
    Table("users u").
    Select("u.id", "u.name", "p.title").
    LeftJoin("profiles p", "p.user_id = u.id").
    Where("u.status = ?", "active").
    Where("u.age BETWEEN ? AND ?", 18, 65).
    GroupBy("u.department").
    Having("COUNT(*) > ?", 5).
    OrderBy("u.created_at", "DESC").
    Limit(20).
    Offset(0)
```

### 批量操作

```go
// 批量插入
users := []map[string]interface{}{
    {"name": "Alice", "email": "alice@example.com", "age": 25},
    {"name": "Bob", "email": "bob@example.com", "age": 30},
}

err := builder.Table("users").BatchInsert(users)

// 批量更新
updates := []map[string]interface{}{
    {"id": 1, "status": "updated", "last_login": time.Now()},
    {"id": 2, "status": "updated", "last_login": time.Now()},
}

err = builder.Table("users").BatchUpdate(updates, []string{"id"})
```

### 事务处理

```go
// 开始事务
tx, err := builder.BeginTx(nil)
if err != nil {
    return err
}

// 在事务中执行操作
err = tx.Table("users").BatchInsert(newUsers)
if err != nil {
    tx.Rollback()
    return err
}

err = tx.Table("logs").BatchInsert(logEntries) 
if err != nil {
    tx.Rollback()
    return err
}

// 提交事务
return tx.Commit()
```

### 功能检测

```go
// 检查适配器支持的功能
if builder.SupportsFeature("orm") {
    ormInstance := builder.GetORMInstance()
    // 使用ORM特有功能
}

if builder.SupportsFeature("upsert") {
    // 使用Upsert操作
}

if builder.SupportsFeature("bulk_insert") {
    // 使用批量插入优化
}
```

## 🎯 框架对比

| 功能 | SQLX | GORM | XORM | Beego ORM | Ent |
|------|------|------|------|-----------|-----|
| ORM支持 | ❌ | ✅ | ✅ | ✅ | ✅ |
| 查询构建器 | ✅ | ✅ | ✅ | ✅ | ✅ |
| 事务支持 | ✅ | ✅ | ✅ | ✅ | ✅ |
| 连接池 | ✅ | ✅ | ✅ | ✅ | ✅ |
| 批量插入 | ✅ | ✅ | ✅ | ✅ | ✅ |
| Upsert | 🔶 | ✅ | ✅ | ✅ | ✅ |
| 代码生成 | ❌ | ❌ | ✅ | ❌ | ✅ |
| 性能 | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐⭐⭐ |

> 🔶 取决于底层数据库支持

## 🏗️ 架构设计

```
┌─────────────────────────────────────────────────────────┐
│                Enhanced SQLBuilder                       │
├─────────────────────────────────────────────────────────┤
│  统一API层 │ Table() │ Select() │ Where() │ Join() │... │
├─────────────────────────────────────────────────────────┤
│                  UniversalAdapter                       │
├─────────────────────────────────────────────────────────┤
│ SQLX    │ GORM     │ XORM    │ Beego   │ Ent    │ ...   │
│ Adapter │ Adapter  │ Adapter │ Adapter │ Adapter│       │
├─────────────────────────────────────────────────────────┤
│ mysql   │ postgres │ sqlite  │ oracle  │ mssql  │ ...   │
│ driver  │ driver   │ driver  │ driver  │ driver │       │
└─────────────────────────────────────────────────────────┘
```

## 🔮 路线图

### 已完成 ✅
- [x] 通用适配器架构设计
- [x] SQLX适配器实现
- [x] GORM适配器实现  
- [x] 自动适配器检测
- [x] 统一API设计
- [x] 事务支持
- [x] 批量操作支持

### 进行中 🚧
- [ ] XORM适配器实现
- [ ] Beego ORM适配器实现
- [ ] Ent适配器实现
- [ ] 原生database/sql适配器
- [ ] 数据库方言引擎

### 计划中 📋
- [ ] 高级查询功能（窗口函数、CTE、递归查询）
- [ ] 更多数据库支持（Oracle、SQL Server、ClickHouse）
- [ ] 性能优化和缓存机制
- [ ] 查询计划分析
- [ ] 可视化工具
- [ ] 插件系统

## 🤝 贡献指南

我们欢迎所有形式的贡献！

### 添加新的适配器

1. 在 `adapter/` 目录创建新的适配器文件
2. 实现 `UniversalAdapter` 接口
3. 创建对应的工厂类
4. 在注册中心注册新适配器
5. 添加测试用例

### 示例：添加新适配器

```go
// your_orm_adapter.go
type YourORMAdapter struct {
    db *yourorm.DB
}

func (a *YourORMAdapter) Query(ctx context.Context, query string, args ...interface{}) (ResultSet, error) {
    // 实现查询逻辑
}

// 实现其他接口方法...

// your_orm_factory.go  
type YourORMAdapterFactory struct{}

func (f *YourORMAdapterFactory) Create(instance interface{}) (UniversalAdapter, error) {
    if db, ok := instance.(*yourorm.DB); ok {
        return &YourORMAdapter{db: db}, nil
    }
    return nil, fmt.Errorf("unsupported instance type")
}
```

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 💬 社区

- 📧 Email: 501893067@qq.com
- 💬 Discussions: [GitHub Discussions](https://github.com/kamalyes/go-sqlbuilder/discussions)
- 🐛 Bug Reports: [GitHub Issues](https://github.com/kamalyes/go-sqlbuilder/issues)

---

**⭐ 如果这个项目对您有帮助，请给我们一个星标！**