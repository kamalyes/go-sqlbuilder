# Go-SQLBuilder: 项目概览与 V3 演进

**版本**: 2.5 (V2 已完成, V3 规划中)
**最后更新**: 2025-11-13

## 1. 项目概览

`go-sqlbuilder` 是一个生产级的 Go 语言 SQL 查询构建器，其核心设计目标是提供一个**好用、易用、现代、高性能且易于替换**的数据访问层。它通过统一的适配器层，支持如 SQLX、GORM 等多种主流数据库框架，并整合了强大的 Repository 模式，为开发者提供流畅、类型安全的数据库操作体验。

### 核心特性 (V2)

- **通用适配器**: 智能检测并适配 GORM、SQLX 等多种数据库实例。
- **链式查询构建**: 提供无限链式调用的 Fluent API，代码直观易读。
- **高级查询参数**: 超过 20 种便捷方法 (`AddEQ`, `AddLike`, `AddIn` 等)，简化复杂查询。
- **内置缓存层**: 支持 Redis 和内存缓存，提供自动缓存键生成、TTL 管理和性能统计。
- **统一错误处理**: 定义了 48 个标准错误码，简化错误处理逻辑。
- **Repository 模式**: 整合 `go-data-repository` 的能力，提供泛型 `Repository[T]` 接口。
- **生产就绪**: 拥有超过 50 个单元测试，覆盖所有核心功能，确保稳定性。

---

## 2. V2 架构回顾

V2 版本的架构实现了清晰的关注点分离，主要分为四层：

```bash
┌─────────────────────────────────────────────────────────────┐
│                     应用层 (Application Layer)               │
│         (Web 应用, API 服务, 命令行工具等)                   │
└────────────────────────┬────────────────────────────────────┘
                         │
┌────────────────────────┴────────────────────────────────────┐
│                  核心层 (Core Layer)                         │
│  - Builder: SQL 构建与执行                                  │
│  - CachedBuilder: 缓存装饰器                                │
└────────────────┬─────────────────────────┬──────────────────┘
                 │                         │
     ┌───────────┴──────────┐   ┌──────────┴──────────┐
     │                      │   │                     │
┌────┴────────────────┐ ┌──┴──────────────┐  ┌──────┴──────────┐
│ 适配层 (Adapter)    │ │  服务层 (Service) │  │  缓存层 (Cache)   │
│  - SQLx Adapter     │ │ - Query Params  │  │ - Store 接口    │
│  - GORM Adapter     │ │ - Error Handling│  │ - Redis/Mock    │
└─────────────────────┘ └─────────────────┘  └─────────────────┘
```

尽管 V2 架构功能完善，但在多个模块（`query/`, `builder.go`, `repository/`）中仍存在**操作符定义、过滤逻辑、排序和分页处理的重复**，这为 V3 的统一重构提供了动机。

---

## 3. V3 架构演进：设计与规划

V3 架构的核心目标是**消除重复、统一抽象**，引入更清晰的分层模型。

### V3 核心设计原则

1. **架构分层**: 引入独立的**执行层 (Executor)**、**编译层 (Compiler)** 和 **中间件层 (Middleware)**。
2. **常量统一**: 所有常量（操作符、错误码、SQL 关键字等）集中到 `constant/` 目录。
3. **日志先行**: 设计统一的 `logger/` 接口，默认集成 `go-logger`，支持任意实现。
4. **核心抽象**: 创建 `core/` 包，提供统一的 `FilterBuilder`, `OrderBuilder`, `PaginationBuilder` 和 `FilterApplier`。

### V3 架构图

```text
┌─────────────────────────────────────────┐
│   应用层 (Application Layer)            │
└────────────────────┬────────────────────┘
                     ↓
┌─────────────────────────────────────────┐
│   仓储层 (Repository Layer) ✅ 保留      │
└────────────────────┬────────────────────┘
                     ↓
┌─────────────────────────────────────────┐
│   查询构建层 (Query Builder Layer)      │
└────────────────────┬────────────────────┘
                     ↓
┌─────────────────────────────────────────┐
│   中间件层 (Middleware Layer) ⭐ 新增   │
└────────────────────┬────────────────────┘
                     ↓
┌─────────────────────────────────────────┐
│   执行引擎 (Execution Engine) ⭐ 新增  │
└────────────────────┬────────────────────┘
                     ↓
┌─────────────────────────────────────────┐
│   编译层 (Compiler Layer) ⭐ 新增      │
└────────────────────┬────────────────────┘
                     ↓
┌─────────────────────────────────────────┐
│   适配器层 (Adapter Layer) ✅ 保留      │
└────────────────────┬────────────────────┘
                     ↓
┌─────────────────────────────────────────┐
│   数据库驱动层 (Driver Layer)           │
└─────────────────────────────────────────┘
```

### V3 实施计划

#### Phase 1: 基础设施 (✅ 已完成)

- **`constant/`**: 创建 11 个文件，统一管理操作符、错误码、SQL 关键字、配置等所有常量。
- **`logger/`**: 创建日志接口、`go-logger` 适配器、`NoOp` 实现和日志工厂。
- **`core/` (基础)**: 定义 `Filter`, `OrderBy`, `Pagination` 等核心类型，并创建 `FilterBuilder`, `OrderBuilder`, `PaginationBuilder`。

#### Phase 2: GORM 集成 (✅ 已完成)

- **`core/filter_applier.go`**: 实现将 `core` 包中构建的查询条件（过滤、排序、分页）无缝应用到 GORM 查询。

#### Phase 3: 执行引擎 (Executor)

- 创建 `executor/` 包，负责 SQL 的最终执行、事务管理、钩子（Hooks）和缓存调用。

#### Phase 4: 中间件系统 (Middleware)

- 创建 `middleware/` 包，实现日志、性能指标、重试、超时等可插拔的中间件。

#### Phase 5: 编译层 (Compiler)

- 创建 `compiler/` 包，处理 SQL 方言（Dialect）差异，并为未来的查询优化做准备。

#### Phase 6: 模块集成与兼容性 (✅ 已完成)

- 通过 `builder_enhancer.go` 和 `compat.go` 建立 V2 到 V3 的桥梁，确保 100% 向后兼容。

---

## 4. 快速使用指南 (V2 API)

### 基础查询

```go
import "github.com/kamalyes/go-sqlbuilder/unified" // 推荐使用统一包

// 1. 初始化 (自动检测 GORM/SQLX)
builder, _ := unified.New(db)

// 2. 执行查询
var users []User
err := builder.Table("users").
    Select("id", "name", "email").
    Where("status", 1).
    OrderBy("created_at", "DESC").
    Limit(10).
    Find(&users)
```

### 高级查询 (便捷方法)

```go
import "github.com/kamalyes/go-sqlbuilder/query"

param := query.NewParam().
    AddEQ("status", 1).
    AddGT("age", 18).
    AddLike("name", "John").
    AddIn("category", 1, 2, 3).
    AddOrder("created_at", "DESC").
    SetPage(1, 20)

// 生成 WHERE 子句和参数
whereSQL, args := param.BuildWhereClause()
```

### 缓存查询

```go
import "github.com/kamalyes/go-sqlbuilder/cache"

// 1. 创建缓存存储
store := cache.NewRedisStore(redisClient, "myapp:")

// 2. 创建带缓存的 Builder
cachedBuilder, _ := unified.NewCachedBuilder(db, store, nil)

// 3. 执行缓存查询 (自动处理缓存读写)
var users []User
err := cachedBuilder.Table("users").
    Where("status", 1).
    WithTTL(10 * time.Minute). // 自定义本次查询的 TTL
    Find(&users)
```

### 事务处理

```go
err := builder.Transaction(func(tx *unified.Builder) error {
    // 在事务中执行多个操作
    if err := tx.Table("users").Set("balance", -100).Where("id", 1).Exec(); err != nil {
        return err // 返回 error 会自动回滚
    }
    if err := tx.Table("logs").Insert(logData).Exec(); err != nil {
        return err
    }
    return nil // 返回 nil 会自动提交
})
```

---

## 5. 现代化与质量保证

- **并发安全**: V3 设计中已规划对 Builder 和缓存模块添加锁保护 (`sync.RWMutex`, `sync.Map`)。
- **性能优化**: 推荐使用 `strings.Builder` 进行 SQL 拼接，并对切片进行预分配。
- **CI/CD**: 已规划完整的 GitHub Actions 工作流，包括自动化测试、基准测试和代码质量检查。
- **测试体系**: V2 已包含 50+ 单元测试。V3 规划了更完善的基准测试和集成测试（使用 Docker Compose）。

这个文档为您浓缩了 `go-sqlbuilder` 项目的精华，指明了其从一个功能完备的 V2 版本向一个架构更优、高度抽象的 V3 版本的演进路径。
