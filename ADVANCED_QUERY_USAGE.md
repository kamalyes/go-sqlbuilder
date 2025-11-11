# Advanced Query 高级查询 - 便捷 API 使用指南

## 概述

`AdvancedQueryParam` 提供了简化的便捷方法（Convenience Methods），无需传递操作符参数，直接通过方法名称表达查询意图。

## 基础过滤方法

### 等于 (EQ)

```go
// 方式1：使用便捷方法（推荐）
aq := NewAdvancedQueryParam().
    AddEQ("status", "active")

// 方式2：使用标准方法
aq := NewAdvancedQueryParam().
    AddFilter("status", OP_EQ, "active")
```

### 比较操作符

```go
aq := NewAdvancedQueryParam().
    AddGT("age", 18).              // 大于 >
    AddGTE("score", 80).           // 大于等于 >=
    AddLT("price", 1000).          // 小于 <
    AddLTE("quantity", 100).       // 小于等于 <=
    AddNEQ("status", "deleted")    // 不等于 !=
```

### 字符串匹配

```go
aq := NewAdvancedQueryParam().
    AddLike("name", "john").              // 全模糊 LIKE %john%
    AddStartsWith("username", "user").    // 前缀匹配 LIKE user%
    AddEndsWith("email", "@qq.com")       // 后缀匹配 LIKE %@qq.com
```

### IN 操作符

```go
aq := NewAdvancedQueryParam().
    AddIn("status", "active", "pending", "processing")

aq := NewAdvancedQueryParam().
    AddIn("user_id", 1, 2, 3, 4, 5)
```

## OR 条件方法

在前一个过滤条件后添加 OR 逻辑的过滤条件：

```go
aq := NewAdvancedQueryParam().
    AddEQ("status", "active").
    AddOrEQ("status", "pending").    // 将前一条件的逻辑改为 OR
    AddOrGT("priority", 5).
    AddOrLike("title", "urgent")
```

## 排序方法

```go
aq := NewAdvancedQueryParam().
    AddOrderAsc("created_at").       // 升序
    AddOrderDesc("score").           // 降序
    AddOrder("name", "ASC")          // 通用方法
```

## 分页方法

```go
aq := NewAdvancedQueryParam().
    SetPage(2, 20)  // 第2页，每页20条
```

## 时间范围

```go
startTime := time.Now().Add(-24 * time.Hour)
endTime := time.Now()

aq := NewAdvancedQueryParam().
    AddTimeRange("created_at", startTime, endTime)
```

## MySQL FIND_IN_SET

```go
aq := NewAdvancedQueryParam().
    AddFindInSet("tags", "golang", "database", "cache")
```

## 完整示例

```go
// 查询：状态为 active 或 pending 的用户，年龄大于18，名字包含"john"，按创建时间降序
aq := NewAdvancedQueryParam().
    AddEQ("status", "active").
    AddOrEQ("status", "pending").
    AddGT("age", 18).
    AddLike("name", "john").
    AddOrderDesc("created_at").
    SetPage(1, 20).
    SetSelectFields("id", "name", "email", "age")

// 构建 WHERE 子句
where, args := aq.BuildWhereClause()
// SQL: WHERE status = ? AND age > ? AND name LIKE ?
// args: ["active", 18, "%john%"]
```

## 高级特性

### 字段选择

```go
aq := NewAdvancedQueryParam().
    AddEQ("status", "active").
    SetSelectFields("id", "name", "email")  // 只查询这些字段
```

### 去重

```go
aq := NewAdvancedQueryParam().
    AddEQ("category", "book").
    SetDistinct(true)  // SELECT DISTINCT
```

### HAVING 子句

```go
aq := NewAdvancedQueryParam().
    AddEQ("status", "active").
    AddHaving("COUNT(*) > 5").
    AddHaving("SUM(amount) > 100")
```

## 方法速查表

| 方法 | 功能 | 示例 |
|------|------|------|
| `AddEQ(field, value)` | 等于 | `AddEQ("status", "active")` |
| `AddGT(field, value)` | 大于 | `AddGT("age", 18)` |
| `AddGTE(field, value)` | 大于等于 | `AddGTE("score", 80)` |
| `AddLT(field, value)` | 小于 | `AddLT("price", 1000)` |
| `AddLTE(field, value)` | 小于等于 | `AddLTE("quantity", 100)` |
| `AddNEQ(field, value)` | 不等于 | `AddNEQ("deleted_at", nil)` |
| `AddLike(field, value)` | 全模糊 | `AddLike("name", "john")` |
| `AddStartsWith(field, value)` | 前缀匹配 | `AddStartsWith("code", "ABC")` |
| `AddEndsWith(field, value)` | 后缀匹配 | `AddEndsWith("domain", "qq.com")` |
| `AddIn(field, values...)` | IN 列表 | `AddIn("status", "a", "b")` |
| `AddOrEQ(field, value)` | OR 等于 | `AddOrEQ("status", "pending")` |
| `AddOrGT(field, value)` | OR 大于 | `AddOrGT("price", 100)` |
| `AddOrLike(field, value)` | OR 模糊 | `AddOrLike("title", "sale")` |
| `AddOrIn(field, values...)` | OR IN | `AddOrIn("type", 1, 2, 3)` |
| `AddOrderAsc(field)` | 升序 | `AddOrderAsc("created_at")` |
| `AddOrderDesc(field)` | 降序 | `AddOrderDesc("score")` |
| `SetPage(page, size)` | 分页 | `SetPage(2, 20)` |
| `SetDistinct(bool)` | 去重 | `SetDistinct(true)` |
| `SetSelectFields(...)` | 字段选择 | `SetSelectFields("id", "name")` |
| `AddTimeRange(field, start, end)` | 时间范围 | `AddTimeRange("created_at", t1, t2)` |
| `AddFindInSet(field, values...)` | FIND_IN_SET | `AddFindInSet("tags", "go", "db")` |
| `AddHaving(clause)` | HAVING 条件 | `AddHaving("COUNT(*) > 5")` |

## 链式调用特性

所有便捷方法都支持链式调用，使代码更流畅：

```go
aq := NewAdvancedQueryParam().
    AddEQ("status", "active").
    AddGT("age", 18).
    AddLike("name", "test").
    AddIn("category", "A", "B").
    AddOrderDesc("created_at").
    SetPage(1, 20)
```

## 与缓存集成

```go
// 创建带缓存的查询构建器
cachedBuilder, err := NewCachedBuilder(dbInstance, cache, CacheConfig{
    Enabled:   true,
    TTL:       1 * time.Hour,
    KeyPrefix: "sqlbuilder:",
})

if err != nil {
    return err
}

// 使用 AdvancedQueryParam 构建查询
aq := NewAdvancedQueryParam().
    AddEQ("status", "active").
    AddGT("score", 80).
    SetPage(1, 20)

// 应用查询参数到构建器
// 然后执行缓存查询
result, err := cachedBuilder.GetCached(&users)
```

## 日志输出

项目使用 `go-logger` 进行日志记录。在构建复杂查询时，会自动记录调试信息：

```go
// 日志会显示查询构建过程
logger.SetGlobalLevel(logger.DEBUG)

aq := NewAdvancedQueryParam().
    AddEQ("status", "active").
    BuildWhereClause()  // 会输出调试日志
```

## 最佳实践

1. **链式调用**：使用链式调用使代码更易读
2. **便捷方法优先**：优先使用 `AddEQ` 等便捷方法而不是 `AddFilter`
3. **分页设置**：在最后设置分页选项
4. **字段选择**：明确指定需要的字段以提高性能
5. **时间范围**：对时间范围查询使用 `AddTimeRange`

## 常见问题

### Q: `AddEQ` 和 `AddEQFilter` 的区别？
A: 没有功能区别，`AddEQ` 是 `AddEQFilter` 的简洁别名，推荐使用 `AddEQ`。

### Q: 如何组合 AND 和 OR 条件？
A: 第一个 `Add*` 方法总是 AND，之后的条件可以用 `AddOr*` 方法变为 OR：
```go
aq.AddEQ("a", 1).        // AND a = 1
    AddOrEQ("b", 2).     // OR b = 2
    AddOrEQ("c", 3)      // OR c = 3
```

### Q: 时间范围查询如何使用？
A: 使用 `AddTimeRange` 方法，传入字段名和开始/结束时间：
```go
aq.AddTimeRange("created_at", startTime, endTime)
```

## 相关文档

- [Cache Redis 缓存集成](./cache.go)
- [Builder SQL 构建器](./builder.go)
- [Interfaces 接口定义](./interfaces.go)
