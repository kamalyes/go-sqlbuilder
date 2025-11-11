# 完整使用指南 - Go-SQLBuilder v2.0

## 目录

1. [快速开始](#快速开始)
2. [基础查询](#基础查询)
3. [高级查询](#高级查询)
4. [缓存管理](#缓存管理)
5. [错误处理](#错误处理)
6. [最佳实践](#最佳实践)

---

## 快速开始

### 安装

```bash
go get github.com/kamalyes/go-sqlbuilder
```

### 基本使用

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
    db, err := sqlx.Connect("mysql", "user:pwd@tcp(localhost:3306)/testdb")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // 创建Builder
    builder := sqlbuilder.New(db)

    // 执行查询
    var users []User
    err = builder.Table("users").
        Select("id", "name", "email").
        Where("status", 1).
        Find(&users)
}
```

---

## 基础查询

### SELECT 查询

#### 1. 简单查询

```go
builder.Table("users").
    Select("id", "name", "email").
    Where("id", 1).
    First(&user)
```

#### 2. 多条件查询

```go
builder.Table("users").
    Select("*").
    Where("status", 1).
    Where("age", ">", 18).
    Where("city", "Beijing").
    Find(&users)
```

#### 3. IN 条件

```go
builder.Table("users").
    WhereIn("id", []interface{}{1, 2, 3, 4, 5}).
    Find(&users)
```

#### 4. BETWEEN 条件

```go
builder.Table("orders").
    WhereBetween("created_at", "2024-01-01", "2024-12-31").
    Find(&orders)
```

#### 5. IS NULL 条件

```go
builder.Table("users").
    WhereNull("deleted_at").
    Find(&activeUsers)
```

#### 6. LIKE 条件

```go
builder.Table("users").
    WhereLike("name", "%john%").
    Find(&users)
```

#### 7. 排序

```go
builder.Table("users").
    OrderBy("created_at", "DESC").
    OrderBy("name", "ASC").
    Find(&users)
```

#### 8. LIMIT 和 OFFSET

```go
builder.Table("users").
    Limit(10).
    Offset(20).
    Find(&users)
```

#### 9. 分页

```go
builder.Table("users").
    Paginate(2, 20).  // 第2页，每页20条
    Find(&users)
```

### INSERT 操作

#### 1. 单条插入

```go
data := map[string]interface{}{
    "name": "John",
    "email": "john@example.com",
    "age": 25,
}
result, err := builder.Table("users").
    Insert(data).
    Exec()
lastID := result.LastInsertId()
```

#### 2. 批量插入

```go
dataList := []map[string]interface{}{
    {"name": "John", "email": "john@example.com"},
    {"name": "Jane", "email": "jane@example.com"},
    {"name": "Bob", "email": "bob@example.com"},
}
for _, data := range dataList {
    builder.Table("users").Insert(data).Exec()
}
```

### UPDATE 操作

```go
data := map[string]interface{}{
    "name": "Updated Name",
    "updated_at": time.Now(),
}
result, err := builder.Table("users").
    Where("id", 1).
    Update(data).
    Exec()
```

### DELETE 操作

```go
result, err := builder.Table("users").
    Where("id", 1).
    Delete().
    Exec()
```

---

## 高级查询

### 使用 query.Param

```go
import "github.com/kamalyes/go-sqlbuilder/query"

// 创建高级查询参数
param := query.NewParam().
    AddEQ("status", 1).
    AddGT("age", 18).
    AddLike("name", "John").
    AddIn("category", 1, 2, 3).
    AddOrder("created_at", "DESC").
    SetPage(1, 20)

// 构建WHERE子句
whereSQL, args := param.BuildWhereClause()
// whereSQL: "WHERE status = ? AND age > ? AND name LIKE ? AND category IN (?, ?, ?)"
// args: [1, 18, "%John%", 1, 2, 3]
```

### 所有便捷方法

#### 比较操作

```go
param := query.NewParam().
    AddEQ("id", 1).           // 等于
    AddNEQ("status", 0).      // 不等于
    AddGT("price", 100).      // 大于
    AddGTE("score", 80).      // 大于等于
    AddLT("age", 60).         // 小于
    AddLTE("stock", 10)       // 小于等于
```

#### 字符串操作

```go
param := query.NewParam().
    AddLike("name", "john").              // 全模糊 %john%
    AddStartsWith("username", "admin").   // 前缀 admin%
    AddEndsWith("email", "@qq.com")       // 后缀 %@qq.com
```

#### 集合操作

```go
param := query.NewParam().
    AddIn("status", 1, 2, 3).
    AddFindInSet("tags", "hot", "new")    // MySQL FIND_IN_SET
```

#### 时间范围

```go
param := query.NewParam().
    AddTimeRange("created_at",
        time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
        time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC),
    )
```

#### 分页和排序

```go
param := query.NewParam().
    AddOrder("created_at", "DESC").
    AddOrder("id", "ASC").
    SetPage(1, 20).
    SetDistinct(true)
```

#### OR 条件

```go
param := query.NewParam().
    AddEQ("role", "admin").
    AddOrEQ("role", "moderator").
    AddOrGT("vip_level", 5)
```

---

## 缓存管理

### 基础缓存使用

```go
import (
    "github.com/kamalyes/go-sqlbuilder/cache"
)

// 创建缓存存储
store := cache.NewMockStore()  // 开发/测试
// 或
store := cache.NewRedisStore(redisClient, "myapp:")

// 创建带缓存的Builder
cachedBuilder, err := sqlbuilder.NewCachedBuilder(db, store, nil)
```

### 缓存查询

```go
// 带自动缓存的查询
result, err := cachedBuilder.GetCached(ctx, sql, args...)

// 获取第一行
row, err := cachedBuilder.FirstCached(ctx, sql, args...)

// 获取计数
count, err := cachedBuilder.CountCached(ctx, sql, args...)
```

### 自定义 TTL

```go
// 为此查询设置30分钟的缓存
cachedBuilder.WithTTL(30 * time.Minute)
result, err := cachedBuilder.GetCached(ctx, sql, args...)
```

### 缓存配置

```go
config := cache.NewConfig().
    SetEnabled(true).
    SetKeyPrefix("myapp:").
    SetDefaultTTL(1 * time.Hour).
    SetMaxSize(10000)
```

### 缓存管理器

```go
manager := cache.NewManager(store)

// 记录访问
manager.RecordHit()
manager.RecordMiss()

// 获取统计
stats := manager.GetStats()
fmt.Printf("命中率: %.2f%%\n", stats.HitRate * 100)

// 使缓存失效
manager.InvalidatePattern(ctx, "myapp:users:*")
```

---

## 错误处理

### 错误码使用

```go
import "github.com/kamalyes/go-sqlbuilder/errors"

// 创建错误
err := errors.NewError(errors.ErrCodePageNumberInvalid, "page must be > 0")

// 使用格式化
err := errors.NewErrorf(errors.ErrCodeInvalidTableName, "table %s not found", tableName)
```

### 错误检查

```go
if errors.IsErrorCode(err, errors.ErrCodePageNumberInvalid) {
    // 处理分页错误
    log.Printf("分页错误: %v", err)
}

// 提取错误码
code := errors.GetErrorCode(err)
message := errors.ErrorCodeString(code)
```

### 标准错误码

```go
// 构建器错误
errors.ErrCodeBuilderNotInitialized
errors.ErrCodeInvalidTableName
errors.ErrCodeInvalidFieldName

// 缓存错误
errors.ErrCodeCacheStoreNotFound
errors.ErrCodeCacheKeyNotFound
errors.ErrCodeCacheExpired

// 查询错误
errors.ErrCodePageNumberInvalid
errors.ErrCodePageSizeInvalid
errors.ErrCodeInvalidOperator

// Redis错误
errors.ErrCodeRedisConnFailed
errors.ErrCodeRedisOperationFailed
```

---

## 最佳实践

### 1. 使用工厂函数

badgo
// ❌ 不要直接创建结构体
builder := &sqlbuilder.Builder{...}

// ✅ 使用工厂函数
builder := sqlbuilder.New(db)
cachedBuilder, _ := sqlbuilder.NewCachedBuilder(db, store, config)

### 2. 总是关闭资源

go
db, _ := sqlx.Connect(...)
defer db.Close()

builder := sqlbuilder.New(db)
// 使用builder

### 3. 使用链式调用

go
// ✅ 推荐：链式调用，清晰易读
builder.Table("users").
    Select("id", "name").
    Where("status", 1).
    OrderBy("created_at", "DESC").
    Limit(10).
    Find(&users)

### 4. 利用高级查询参数

go
// ✅ 推荐：使用query.Param构建复杂查询
param := query.NewParam().
    AddGT("age", 18).
    AddLike("name", "John").
    AddIn("status", 1, 2).
    SetPage(1, 20)

whereSQL, args := param.BuildWhereClause()
// 在Builder中使用
builder.Table("users").WhereRaw(whereSQL, args...).Find(&users)

### 5. 缓存重要查询

go
// ✅ 为频繁查询的数据使用缓存
cachedBuilder, _ := sqlbuilder.NewCachedBuilder(db, redisStore, config)

// 用户经常查询自己的个人资料
userID := 123
sql := fmt.Sprintf("SELECT * FROM users WHERE id = %d", userID)
user, _ := cachedBuilder.FirstCached(ctx, sql)

### 6. 错误处理

go
// ✅ 正确的错误处理
if err != nil {
    if errors.IsErrorCode(err, errors.ErrCodeCacheStoreNotFound) {
        // 缓存不可用，降级使用数据库
        result, _ := builder.Query(sql).First()
    } else {
        // 其他错误
        log.Printf("错误: %v", err)
    }
}

### 7. 事务处理

go
// 开始事务
tx, err := builder.Begin()
if err != nil {
    return err
}
defer tx.Rollback()

// 执行多个操作
if err := tx.Table("orders").Insert(order).Exec(); err != nil {
    return err
}
if err := tx.Table("inventories").Update(data).Exec(); err != nil {
    return err
}

// 提交事务
if err := tx.Commit(); err != nil {
    return err
}

### 8. 性能优化

go
// ✅ 使用查询优化
builder.Table("users").
    Select("id", "name", "email").      // 只选择需要的字段
    Where("deleted_at", nil).            // 过滤已删除的记录
    Limit(100).                          // 添加LIMIT
    Find(&users)

// ✅ 使用缓存减少查询
cachedBuilder.GetCached(ctx, sql).WithTTL(5*time.Minute)

---

## 常见问题

### Q: 如何处理 NULL 值?

go
// 查询 NULL
builder.WhereNull("deleted_at").Find(&users)

// 查询非 NULL
builder.WhereNotNull("deleted_at").Find(&users)

### Q: 如何进行 JOIN 操作?

go
builder.Table("users").
    LeftJoin("orders", "users.id", "=", "orders.user_id").
    Select("users.id", "users.name", "orders.order_id").
    Find(&results)

### Q: 如何支持多个数据库?

go
db1 := sqlx.Connect("mysql", "dsn1")
db2 := sqlx.Connect("mysql", "dsn2")

builder1 := sqlbuilder.New(db1)
builder2 := sqlbuilder.New(db2)

### Q: 缓存如何失效?

go
// 手动失效
store.Delete(ctx, cacheKey)

// 模式失效
manager.InvalidatePattern(ctx, "myapp:users:*")

// TTL自动失效
cachedBuilder.WithTTL(1 * time.Hour)  // 1小时后自动过期

---

## 参考资源

- [项目分析文档](PROJECT_ANALYSIS.md)
- [架构设计文档](ARCHITECTURE.md)
- [高级查询使用指南](ADVANCED_QUERY_USAGE.md)
- [GitHub 仓库](https://github.com/kamalyes/go-sqlbuilder)
