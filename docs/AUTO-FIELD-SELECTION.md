# 自动字段选择功能

## 概述

自动字段选择功能允许 Repository 根据 Model 的 struct tags 自动生成查询字段列表，避免使用 `SELECT *`，提升查询性能。

## 特性

- ✅ **自动字段提取**：基于 struct tags（gorm/json）自动识别字段
- ✅ **字段缓存**：提取的字段会被缓存，避免重复反射
- ✅ **智能模式**：支持手动 Select/Omit 覆盖自动字段
- ✅ **动态切换**：可随时启用/禁用自动字段模式
- ✅ **性能优化**：避免查询不需要的大字段

## 使用方法

### 1. 创建时启用自动字段

```go
import (
    "github.com/kamalyes/go-sqlbuilder/repository"
)

// 定义 Model
type User struct {
    ID        uint      `json:"id" gorm:"column:id"`
    Name      string    `json:"name" gorm:"column:name"`
    Email     string    `json:"email" gorm:"column:email"`
    Password  string    `json:"password" gorm:"column:password"`
    Content   string    `json:"content" gorm:"column:content;type:text"`
    CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
}

// 创建启用自动字段的 Repository
repo := repository.NewBaseRepository[User](
    dbHandler,
    logger,
    "users",
    repository.WithAutoFields[User](),  // 启用自动字段
)

// 查询会自动使用模型字段，而不是 SELECT *
users, err := repo.List(ctx, repository.NewQuery())
// 生成 SQL: SELECT id, name, email, password, content, created_at FROM users
```

### 2. 动态启用/禁用

```go
// 创建普通 Repository
repo := repository.NewBaseRepository[User](dbHandler, logger, "users")

// 动态启用自动字段
repo.EnableAutoFields()
users, err := repo.List(ctx, repository.NewQuery())
// 使用自动字段

// 动态禁用自动字段
repo.DisableAutoFields()
users, err = repo.List(ctx, repository.NewQuery())
// 使用 SELECT *

// 检查状态
if repo.IsAutoFieldsEnabled() {
    fmt.Println("自动字段模式已启用")
}
```

### 3. 与 Omit 结合使用

自动字段模式下，`Omit` 会智能地从模型字段中排除指定字段：

```go
repo := repository.NewBaseRepository[User](
    dbHandler, logger, "users",
    repository.WithAutoFields[User](),
)

// 排除敏感字段
query := repository.NewQuery().Omit("password", "content")
users, err := repo.List(ctx, query)
// 生成 SQL: SELECT id, name, email, created_at FROM users
```

### 4. 使用便捷方法

```go
// 排除敏感字段
query := repository.NewQuery().OmitSensitive()
// 自动排除: password, secret, token, api_key, access_token, refresh_token

// 排除大字段（优化性能）
query := repository.NewQuery().OmitLargeFields()
// 自动排除: content, description, detail, data, payload, body, remark

// 结合使用
users, err := repo.List(ctx, 
    repository.NewQuery().
        OmitLargeFields().
        AddFilter(repository.NewEqFilter("status", "active")),
)
```

### 5. 手动 Select 优先级更高

即使启用了自动字段，手动指定的 `Select` 仍然优先：

```go
repo := repository.NewBaseRepository[User](
    dbHandler, logger, "users",
    repository.WithAutoFields[User](),
)

// 手动指定字段（覆盖自动字段）
query := repository.NewQuery().Select("id", "name", "email")
users, err := repo.List(ctx, query)
// 生成 SQL: SELECT id, name, email FROM users
```

### 6. 获取模型字段

```go
// 获取缓存的模型字段
fields := repo.GetModelFields()
// 返回: []string{"id", "name", "email", "password", "content", "created_at"}

// 字段会在首次调用时提取并缓存
```

## 工作原理

### 字段提取规则

1. **优先级**：GORM column tag > JSON tag > 字段名蛇形转换
2. **跳过规则**：未导出字段、`gorm:"-"` 标记的字段
3. **缓存机制**：字段提取结果会被缓存，避免重复反射

### 示例

```go
type Product struct {
    ID          uint   `json:"id" gorm:"column:product_id"`           // → product_id
    ProductName string `json:"name" gorm:"column:product_name"`       // → product_name  
    Price       int    `json:"price"`                                 // → price
    Stock       int    `gorm:"column:stock_count"`                    // → stock_count
    Internal    string `json:"-"`                                     // 跳过
    notExported string `json:"hidden"`                                // 跳过（未导出）
}

// 提取结果：["product_id", "product_name", "price", "stock_count"]
```

## 性能对比

### 字段组装开销（纯逻辑性能）

以下是字段选择逻辑本身的性能开销（不包括数据库查询时间）：

```
BenchmarkFieldSelection_Overhead/GetStructFields-8                      1190356        1905 ns/op       400 B/op      13 allocs/op
BenchmarkFieldSelection_Overhead/GetStructFields_Cached-8            1000000000       0.34 ns/op         0 B/op       0 allocs/op
BenchmarkFieldSelection_Overhead/BuildSelectClause_NoOmit-8            2906066         835 ns/op       560 B/op      13 allocs/op
BenchmarkFieldSelection_Overhead/BuildSelectClause_WithOmit-8          3485412         693 ns/op       488 B/op      11 allocs/op
BenchmarkFieldSelection_Overhead/FilterFields-8                        6100388         416 ns/op       240 B/op       4 allocs/op
BenchmarkFieldSelection_Overhead/ApplyFieldSelection_Disabled-8     1000000000        2.32 ns/op         0 B/op       0 allocs/op
BenchmarkFieldSelection_Overhead/ApplyFieldSelection_AutoFields-8     52904686        48.3 ns/op        24 B/op       1 allocs/op
BenchmarkFieldSelection_Overhead/ApplyFieldSelection_AutoFields_WithOmit-8  4922030   482 ns/op       264 B/op       5 allocs/op
BenchmarkFieldSelection_Overhead/ApplyFieldSelection_ManualSelect-8   54158833        44.0 ns/op        24 B/op       1 allocs/op
BenchmarkFieldSelection_Overhead/Query_Construction_SelectAll-8    1000000000        0.27 ns/op         0 B/op       0 allocs/op
BenchmarkFieldSelection_Overhead/Query_Construction_WithSelect-8     60652315        36.7 ns/op        48 B/op       1 allocs/op
BenchmarkFieldSelection_Overhead/Query_Construction_WithOmit-8       82077939        30.3 ns/op        32 B/op       1 allocs/op
BenchmarkFieldSelection_Overhead/EnableDisable_AutoFields-8       1000000000         1.36 ns/op         0 B/op       0 allocs/op
```

**关键指标**：
- ✅ **字段缓存极快**：缓存命中仅 0.34 纳秒，零内存分配
- ✅ **自动字段开销很小**：仅 48 纳秒/次操作，24 字节内存
- ✅ **启用/禁用无开销**：1.36 纳秒，可动态切换
- ⚠️ **首次字段提取**：1905 纳秒（仅执行一次，后续使用缓存）

### 不使用自动字段（SELECT *）

```sql
-- 查询所有字段（包括大字段）
SELECT * FROM users WHERE status = 'active';
```

**问题**：
- 查询不需要的大字段（content, description 等）
- 浪费网络带宽
- 增加内存占用

### 使用自动字段 + Omit

```go
query := repository.NewQuery().
    OmitLargeFields().
    AddFilter(repository.NewEqFilter("status", "active"))
```

```sql
-- 只查询需要的字段
SELECT id, name, email, created_at, updated_at FROM users WHERE status = 'active';
```

**优势**：
- 减少数据传输量
- 降低内存占用
- 提升查询性能

### 性能结论

1. **字段组装开销微乎其微**：自动字段选择的逻辑开销仅 48 纳秒，远小于数据库查询时间（通常几毫秒到几十毫秒）
2. **缓存效果显著**：字段提取结果被缓存后，访问速度接近零开销（0.34 纳秒）
3. **实际收益明显**：通过避免查询不必要的大字段，可显著减少网络传输和内存占用
4. **建议**：专注于优化 SQL 查询本身（索引、条件筛选等），字段选择逻辑的开销可以忽略不计

## 实战示例

### 场景1：用户列表（排除密码）

```go
type User struct {
    ID        uint   `json:"id" gorm:"column:id"`
    Name      string `json:"name" gorm:"column:name"`
    Email     string `json:"email" gorm:"column:email"`
    Password  string `json:"password" gorm:"column:password"`
    CreatedAt time.Time `json:"created_at"`
}

repo := repository.NewBaseRepository[User](
    dbHandler, logger, "users",
    repository.WithAutoFields[User](),
)

// 查询用户列表，排除密码
query := repository.NewQuery().
    Omit("password").
    AddOrder("created_at", "DESC").
    Limit(20)

users, err := repo.List(ctx, query)
// SQL: SELECT id, name, email, created_at FROM users ORDER BY created_at DESC LIMIT 20
```

### 场景2：文章列表（排除内容）

```go
type Article struct {
    ID          uint   `json:"id"`
    Title       string `json:"title"`
    Summary     string `json:"summary"`
    Content     string `json:"content" gorm:"type:text"`  // 大字段
    Author      string `json:"author"`
    PublishedAt time.Time `json:"published_at"`
}

repo := repository.NewBaseRepository[Article](
    dbHandler, logger, "articles",
    repository.WithAutoFields[Article](),
)

// 列表页不需要完整内容
query := repository.NewQuery().
    Omit("content").  // 排除大字段
    AddFilter(repository.NewEqFilter("status", "published")).
    AddOrder("published_at", "DESC")

articles, err := repo.List(ctx, query)
// SQL: SELECT id, title, summary, author, published_at FROM articles 
//      WHERE status = 'published' ORDER BY published_at DESC
```

### 场景3：API响应优化

```go
type Product struct {
    ID          uint   `json:"id"`
    Name        string `json:"name"`
    Price       int    `json:"price"`
    Description string `json:"description" gorm:"type:text"`
    Images      string `json:"images" gorm:"type:json"`
    Stock       int    `json:"stock"`
}

repo := repository.NewBaseRepository[Product](
    dbHandler, logger, "products",
    repository.WithAutoFields[Product](),
)

// 列表API：只返回基本信息
listQuery := repository.NewQuery().
    OmitLargeFields().  // 排除 description, images
    AddFilter(repository.NewEqFilter("available", true))

products, _ := repo.List(ctx, listQuery)

// 详情API：返回完整信息
detailQuery := repository.NewQuery()  // 使用自动字段（包含所有字段）
product, _ := repo.Get(ctx, productID)
```

## 最佳实践

### 1. 默认启用自动字段

```go
// 推荐：为所有 Repository 启用自动字段
repo := repository.NewBaseRepository[User](
    dbHandler, logger, "users",
    repository.WithAutoFields[User](),
)
```

### 2. 列表查询排除大字段

```go
// 列表查询时排除不必要的大字段
listQuery := repository.NewQuery().
    OmitLargeFields().
    Limit(20)
```

### 3. 详情查询使用完整字段

```go
// 详情查询使用所有字段
detailQuery := repository.NewQuery()
user, err := repo.Get(ctx, userID)
```

### 4. API响应分层

```go
// 简单列表
type UserListItem struct {
    ID    uint   `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

// 详细信息
type UserDetail struct {
    ID        uint      `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    Bio       string    `json:"bio"`
    Settings  string    `json:"settings"`
    CreatedAt time.Time `json:"created_at"`
}

// 根据不同场景使用不同 Repository
listRepo := repository.NewBaseRepository[UserListItem](...)
detailRepo := repository.NewBaseRepository[UserDetail](...)
```

## 注意事项

1. **字段缓存**：字段在首次使用时提取并缓存，Model 结构变更后需重启应用
2. **优先级**：手动 `Select` > 自动字段 + `Omit` > 自动字段 > `SELECT *`
3. **性能**：
   - 字段组装开销：48 纳秒（可忽略）
   - 字段缓存访问：0.34 纳秒（接近零开销）
   - 首次字段提取：1905 纳秒（仅执行一次）
   - 建议：专注优化 SQL 查询本身，字段选择逻辑开销可忽略
4. **兼容性**：与现有代码完全兼容，可渐进式启用

## 总结

自动字段选择功能提供了一种优雅的方式来优化查询性能，避免 `SELECT *` 的弊端。性能测试表明，字段选择逻辑本身的开销微乎其微（48 纳秒），远小于数据库查询时间。通过合理使用 `WithAutoFields`、`Omit`、`OmitLargeFields` 等功能，可以显著提升应用性能和代码质量，而无需担心额外的性能损耗。
