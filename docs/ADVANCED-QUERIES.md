# 高级查询

本文档介绍如何使用 Query、Filter 和 FilterGroup 构建复杂查询。

> 💡 **提示**：如果您是初学者，建议先阅读 [CRUD 操作](./CRUD-OPERATIONS.md) 和 [便捷查询方法](./CONVENIENCE-METHODS.md) 了解基础用法。本文档专注于复杂查询场景。

## Query 对象

Query 是查询条件的容器，支持过滤、排序、分页等功能。

### 基础用法

```go
query := repository.NewQuery().
    AddFilter(repository.NewEqFilter("status", "active")).
    AddOrder("created_at", "DESC").
    Limit(10).
    Offset(20)
    
users, err := repo.List(ctx, query)
```

### 预加载关联

```go
query := repository.NewQuery().
    AddPreload("Profile").
    AddPreload("Orders.Items")  // 嵌套预加载
    
users, err := repo.List(ctx, query)
```

### 字段选择

```go
// 只查询指定字段
query := repository.NewQuery().
    Select([]string{"id", "name", "email"})
    
// 排除指定字段
query := repository.NewQuery().
    Omit([]string{"password", "secret_key"})
```

### 分组和聚合

```go
query := repository.NewQuery().
    GroupBy("status").
    Having("COUNT(*) > ?", 10)
```

## Filter 过滤器

Filter 提供了丰富的查询条件构建器。

### 等值查询

```go
// 等于
filter := repository.NewEqFilter("status", "active")

// 不等于
filter := repository.NewNeFilter("status", "deleted")
```

### 比较查询

```go
// 大于
filter := repository.NewGtFilter("age", 18)

// 大于等于
filter := repository.NewGteFilter("age", 18)

// 小于
filter := repository.NewLtFilter("price", 100)

// 小于等于
filter := repository.NewLteFilter("price", 100)
```

### 范围查询

```go
// IN 查询
filter := repository.NewInFilter("id", []interface{}{1, 2, 3, 4, 5})

// NOT IN
filter := repository.NewNotInFilter("status", []interface{}{"deleted", "banned"})

// BETWEEN
filter := repository.NewBetweenFilter("age", 18, 65)

// NOT BETWEEN
filter := repository.NewNotBetweenFilter("price", 100, 200)
```

### 模糊查询

```go
// LIKE '%keyword%'
filter := repository.NewLikeFilter("name", "张")

// 前缀匹配 'prefix%'
filter := repository.NewStartsWithFilter("code", "ORD")

// 后缀匹配 '%suffix'
filter := repository.NewEndsWithFilter("email", "@example.com")

// 包含匹配 '%keyword%'（与 Like 相同，但语义更清晰）
filter := repository.NewContainsFilter("description", "重要")
```

### NULL 查询

```go
// IS NULL
filter := repository.NewIsNullFilter("deleted_at")

// IS NOT NULL
filter := repository.NewIsNotNullFilter("email")
```

### 时间查询

```go
// 时间范围查询
startTime := time.Now().AddDate(0, -1, 0)
endTime := time.Now()
filter := repository.NewBetweenFilter("created_at", startTime, endTime)

// 日期查询（忽略时间部分）
filter := repository.NewDateFilter("created_at", time.Date(2023, 12, 25, 0, 0, 0, 0, time.UTC))
```

## 便捷查询构建方法 🔥

从 v1.2+ 版本开始，Query 提供了更多便捷的链式调用方法，让查询构建更加简洁和直观。

### 基础条件构建

```go
// 使用便捷方法构建查询
query := repository.NewQuery().
    AddEqual("status", 1).                    // 等于条件
    AddNotEqual("type", "deleted").           // 不等于条件
    AddLike("name", "test").                  // 包含匹配
    AddStartsWith("code", "ORD").             // 前缀匹配
    AddEndsWith("email", "@example.com").     // 后缀匹配
    AddIn("category_id", 1, 2, 3).           // IN 条件
    AddNotIn("tag", "spam", "adult").         // NOT IN 条件
    AddGreaterThan("age", 18).                // 大于
    AddGreaterEqual("score", 60).             // 大于等于
    AddLessThan("price", 1000).               // 小于
    AddLessEqual("quantity", 100).            // 小于等于
    AddBetween("age", 18, 65).                // 范围条件
    AddIsNull("deleted_at").                  // NULL 检查
    AddIsNotNull("email")                     // NOT NULL 检查

users, err := repo.List(ctx, query)
```

### 时间便捷方法

```go
// 时间相关的便捷方法
query := repository.NewQuery().
    AddTimeAfter("created_at", time.Now().AddDate(0, -1, 0)).  // 时间晚于
    AddTimeBefore("updated_at", time.Now()).                    // 时间早于
    AddTimeBetween("login_at", start, end).                     // 时间范围
    AddToday("register_date").                                  // 今天
    AddThisWeek("activity_date").                              // 本周
    AddThisMonth("payment_date").                              // 本月
    AddThisYear("join_date")                                   // 今年

users, err := repo.List(ctx, query)
```

### 排序和分页简化

```go
// 排序和分页的便捷方法
query := repository.NewQuery().
    AddEqual("status", 1).
    AddOrderAsc("name").          // 升序排序
    AddOrderDesc("created_at").   // 降序排序
    SetDistinct().               // 去重
    Page(1, 20).                // 分页
    Take(50).                   // 限制数量
    Skip(100)                   // 跳过数量

users, err := repo.List(ctx, query)
```

### 实际应用示例

#### 电商订单查询

```go
// 现在可以这样链式调用构建查询条件
query := repository.NewQuery().
    AddEqual("status", 1).
    AddLike("name", "test").
    AddTimeAfter("created_at", time.Now().AddDate(0, -1, 0)).
    AddIn("category_id", 1, 2, 3).
    AddOrderDesc("created_at").
    Take(10)

orders, err := repo.List(ctx, query)
```

#### 用户活动分析

```go
// 使用时间便捷方法
query := repository.NewQuery().
    AddEqual("status", 1).
    AddThisMonth("created_at").
    AddOrderDesc("id")

activeUsers, err := repo.List(ctx, query)
```

#### 复杂业务查询

```go
// 复杂条件组合
query := repository.NewQuery().
    AddEqual("status", 1).
    AddStartsWith("name", "user_").
    AddBetween("age", 18, 65).
    AddIsNotNull("email").
    Page(1, 20)

eligibleUsers, err := repo.List(ctx, query)
```

#### 报表统计查询

```go
// 财务报表查询示例
startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
endDate := time.Date(2025, 12, 31, 23, 59, 59, 0, time.UTC)

query := repository.NewQuery().
    AddEqual("tenant_id", 1001).                              // 租户过滤
    AddIn("status", "paid", "processing", "completed").       // 状态过滤
    AddGreaterEqual("amount", 100).                           // 最小金额
    AddLessEqual("amount", 50000).                            // 最大金额
    AddStartsWith("order_no", "ORD2025").                     // 订单号前缀
    AddIsNotNull("customer_email").                           // 必须有邮箱
    AddTimeBetween("created_at", startDate, endDate).         // 时间范围
    AddLike("shipping_address", "北京").                       // 地址包含
    AddNotIn("payment_method", "cash", "check").              // 排除支付方式
    AddOrderDesc("created_at").                               // 按创建时间降序
    AddOrderAsc("amount").                                    // 按金额升序
    Page(1, 20)                                               // 分页

reports, err := repo.List(ctx, query)
```

### 方法对比

| 传统方式 | 便捷方法 | 说明 |
|---------|----------|------|
| `AddFilter(NewEqFilter("status", 1))` | `AddEqual("status", 1)` | 更简洁直观 |
| `AddFilter(NewLikeFilter("name", "%test%"))` | `AddLike("name", "test")` | 自动添加通配符 |
| `AddOrder("created_at", "DESC")` | `AddOrderDesc("created_at")` | 语义更清晰 |
| `WithPaging(1, 20)` | `Page(1, 20)` | 更简短 |
| `Limit(10)` | `Take(10)` | 更直观的语义 |

### 链式调用的优势

1. **可读性强**: 代码读起来像自然语言
2. **类型安全**: 编译时检查，减少运行时错误
3. **智能处理**: 自动忽略空值和无效参数
4. **性能优化**: 减少对象创建和内存分配
5. **向后兼容**: 不影响现有的 Filter 方式

### 注意事项

- 时间相关方法会自动忽略零值时间
- 字符串匹配方法会自动忽略空字符串
- IN/NOT IN 方法会自动忽略空切片
- 所有方法都返回同一个 Query 对象，支持链式调用

**实战示例 - API 动态搜索：**

```go
type ProductSearchParams struct {
    Keyword    string   `json:"keyword"`
    Category   int64    `json:"category"`
    MinPrice   float64  `json:"min_price"`
    MaxPrice   float64  `json:"max_price"`
    Status     string   `json:"status"`
    Tags       []int64  `json:"tags"`
}

func (s *ProductService) Search(ctx context.Context, params ProductSearchParams) ([]*Product, error) {
    query := repository.NewQuery().
        AddLikeFilterIfNotEmpty("name", params.Keyword).
        AddEqFilterIfNotEmpty("category_id", params.Category).
        AddGteFilterIfNotEmpty("price", params.MinPrice).
        AddLteFilterIfNotEmpty("price", params.MaxPrice).
        AddEqFilterIfNotEmpty("status", params.Status).
        AddInFilterIfNotEmpty("tag_id", params.Tags).
        AddOrderDesc("created_at")
    
    return s.repo.List(ctx, query)
}
```

### AddSafeOrder - 安全排序（防 SQL 注入）

使用白名单机制防止恶意排序字段

```go
// 定义允许排序的字段白名单
allowedFields := []string{"id", "created_at", "price", "sales"}

query := repository.NewQuery().
    AddSafeOrder(userInputField, userInputDirection, allowedFields)

// 如果 userInputField 不在白名单中，会被忽略
```

**实战示例 - API 动态排序：**

```go
func (s *ProductService) ListProducts(ctx context.Context, sortField, sortDir string) ([]*Product, error) {
    // 白名单防止 SQL 注入
    allowedSorts := []string{"id", "name", "price", "created_at", "sales", "rating"}
    
    query := repository.NewQuery().
        AddEqFilter("status", "active").
        AddSafeOrder(sortField, sortDir, allowedSorts).
        Limit(20)
    
    return s.repo.List(ctx, query)
}
```

### 字段选择辅助方法

#### OmitSensitive - 排除敏感字段

```go
// 排除密码、令牌等敏感信息
query := repository.NewQuery().OmitSensitive()

// 实际调用: query.Omit([]string{"password", "token", "secret_key", "api_key"})
```

#### OmitLargeFields - 排除大字段

```go
// 排除大文本、二进制等大字段，提升查询性能
query := repository.NewQuery().OmitLargeFields()

// 实际调用: query.Omit([]string{"content", "description", "data", "blob"})
```

**实战示例 - 列表 API：**

```go
func (s *UserService) GetUserList(ctx context.Context) ([]*User, error) {
    query := repository.NewQuery().
        OmitSensitive().        // 排除密码等敏感字段
        OmitLargeFields().      // 排除简介等大字段
        AddEqFilter("status", "active").
        AddOrderDesc("created_at").
        Limit(50)
    
    return s.repo.List(ctx, query)
}
```

### 完整动态查询示例

```go
type UserSearchRequest struct {
    Keyword       string    `json:"keyword"`
    Status        string    `json:"status"`
    Role          string    `json:"role"`
    MinAge        int       `json:"min_age"`
    MaxAge        int       `json:"max_age"`
    RegisterStart time.Time `json:"register_start"`
    RegisterEnd   time.Time `json:"register_end"`
    Tags          []int64   `json:"tags"`
    ExcludeTags   []int64   `json:"exclude_tags"`
    SortField     string    `json:"sort_field"`
    SortDir       string    `json:"sort_dir"`
    Page          int       `json:"page"`
    PageSize      int       `json:"page_size"`
}

func (s *UserService) DynamicSearch(ctx context.Context, req UserSearchRequest) ([]*User, int64, error) {
    // 构建动态查询
    query := repository.NewQuery().
        AddLikeFilterIfNotEmpty("name", req.Keyword).
        AddEqFilterIfNotEmpty("status", req.Status).
        AddEqFilterIfNotEmpty("role", req.Role).
        AddGteFilterIfNotEmpty("age", req.MinAge).
        AddLteFilterIfNotEmpty("age", req.MaxAge).
        AddBetweenFilterIfNotEmpty("created_at", req.RegisterStart, req.RegisterEnd).
        AddInFilterIfNotEmpty("tag_id", req.Tags).
        AddNotInFilterIfNotEmpty("tag_id", req.ExcludeTags).
        OmitSensitive().  // 排除敏感字段
        OmitLargeFields() // 排除大字段
    
    // 安全排序
    allowedSorts := []string{"id", "name", "created_at", "age", "status"}
    query.AddSafeOrder(req.SortField, req.SortDir, allowedSorts)
    
    // 分页
    if req.PageSize > 0 {
        query.Limit(req.PageSize).Offset((req.Page - 1) * req.PageSize)
    }
    
    // 查询
    users, err := s.repo.List(ctx, query)
    if err != nil {
        return nil, 0, err
    }
    
    // 统计总数
    total, err := s.repo.Count(ctx, query)
    if err != nil {
        return nil, 0, err
    }
    
    return users, total, nil
}
```

## 其他高级过滤器

```go
// NOT LIKE
filter := repository.NewNotLikeFilter("email", "@temp.com")

// 自定义 LIKE 模式
filter := repository.NewLikeFilter("phone", "138%")  // 以 138 开头
filter := repository.NewLikeFilter("email", "%@gmail.com")  // Gmail 邮箱
```

### NULL 查询

```go
// IS NULL
filter := repository.NewIsNullFilter("deleted_at")

// IS NOT NULL
filter := repository.NewIsNotNullFilter("verified_at")
```

### 时间范围查询

```go
// 今天
filter := repository.NewTodayFilter("created_at")

// 本周
filter := repository.NewThisWeekFilter("created_at")

// 本月
filter := repository.NewThisMonthFilter("created_at")

// 本年
filter := repository.NewThisYearFilter("created_at")

// 昨天
filter := repository.NewYesterdayFilter("created_at")

// 上周
filter := repository.NewLastWeekFilter("created_at")

// 上月
filter := repository.NewLastMonthFilter("created_at")
```

### 自定义 SQL

```go
// 自定义条件
filter := repository.NewCustomFilter("YEAR(created_at) = ?", 2024)

// 复杂条件
filter := repository.NewCustomFilter(
    "price * quantity > ? AND discount < ?",
    1000, 0.5,
)
```

## FilterGroup - 复杂条件组合

FilterGroup 支持 AND/OR 逻辑组合和嵌套。

### 基础 AND 条件

```go
// (status = 'active' AND age > 18)
group := repository.NewFilterGroup(constants.AND_AND).
    AddFilter(repository.NewEqFilter("status", "active")).
    AddFilter(repository.NewGtFilter("age", 18))
    
query := repository.NewQuery().SetFilterGroup(group)
```

### 基础 OR 条件

```go
// (status = 'active' OR status = 'verified')
group := repository.NewFilterGroup(constants.AND_OR).
    AddFilter(repository.NewEqFilter("status", "active")).
    AddFilter(repository.NewEqFilter("status", "verified"))
```

### 嵌套条件

```go
// (status = 'active' AND (age > 18 OR vip_level > 0))
innerGroup := repository.NewFilterGroup(constants.AND_OR).
    AddFilter(repository.NewGtFilter("age", 18)).
    AddFilter(repository.NewGtFilter("vip_level", 0))
    
outerGroup := repository.NewFilterGroup(constants.AND_AND).
    AddFilter(repository.NewEqFilter("status", "active")).
    AddGroup(innerGroup)
    
query := repository.NewQuery().SetFilterGroup(outerGroup)
```

### 复杂示例：电商搜索

```go
// ((category = 'electronics' OR category = 'books') 
//  AND price BETWEEN 100 AND 1000 
//  AND (brand = 'Apple' OR rating >= 4.5))

categoryGroup := repository.NewFilterGroup(constants.AND_OR).
    AddFilter(repository.NewEqFilter("category", "electronics")).
    AddFilter(repository.NewEqFilter("category", "books"))
    
brandGroup := repository.NewFilterGroup(constants.AND_OR).
    AddFilter(repository.NewEqFilter("brand", "Apple")).
    AddFilter(repository.NewGteFilter("rating", 4.5))
    
mainGroup := repository.NewFilterGroup(constants.AND_AND).
    AddGroup(categoryGroup).
    AddFilter(repository.NewBetweenFilter("price", 100, 1000)).
    AddGroup(brandGroup)
    
query := repository.NewQuery().
    SetFilterGroup(mainGroup).
    AddOrder("price", "ASC")
    
products, err := repo.List(ctx, query)
```

### 复杂示例：用户筛选

```go
// ((status = 'active' AND verified = true)
//  OR (status = 'trial' AND created_at > 7 days ago))
//  AND (age > 18 OR is_student = true)

activeGroup := repository.NewFilterGroup(constants.AND_AND).
    AddFilter(repository.NewEqFilter("status", "active")).
    AddFilter(repository.NewEqFilter("verified", true))
    
trialGroup := repository.NewFilterGroup(constants.AND_AND).
    AddFilter(repository.NewEqFilter("status", "trial")).
    AddFilter(repository.NewGtFilter("created_at", time.Now().AddDate(0, 0, -7)))
    
statusGroup := repository.NewFilterGroup(constants.AND_OR).
    AddGroup(activeGroup).
    AddGroup(trialGroup)
    
ageGroup := repository.NewFilterGroup(constants.AND_OR).
    AddFilter(repository.NewGtFilter("age", 18)).
    AddFilter(repository.NewEqFilter("is_student", true))
    
finalGroup := repository.NewFilterGroup(constants.AND_AND).
    AddGroup(statusGroup).
    AddGroup(ageGroup)
    
query := repository.NewQuery().SetFilterGroup(finalGroup)
users, err := repo.List(ctx, query)
```

## 排序

### 单字段排序

```go
query := repository.NewQuery().AddOrder("created_at", "DESC")
```

### 多字段排序

```go
query := repository.NewQuery().
    AddOrder("status", "ASC").
    AddOrder("created_at", "DESC")
```

### 使用常量

```go
query := repository.NewQuery().
    AddOrder("priority", constants.Desc).
    AddOrder("created_at", constants.Asc)
```

### 自定义排序

```go
query := repository.NewQuery().
    AddOrder("FIELD(status, 'urgent', 'high', 'normal', 'low')", "")
```

## 分页

### 基础分页

```go
query := repository.NewQuery().
    Limit(20).
    Offset(40)  // 第 3 页，每页 20 条
```

### 使用 Paging 对象

```go
paging := &repository.Paging{
    Page:     2,
    PageSize: 20,
}

users, paging, err := repo.ListWithPagination(ctx, query, paging)

fmt.Printf("总记录: %d\n", paging.Total)
fmt.Printf("总页数: %d\n", paging.GetTotalPages())
fmt.Printf("是否有下一页: %v\n", paging.HasNext())
fmt.Printf("是否有上一页: %v\n", paging.HasPrev())
```

### 游标分页

```go
// 使用 EnhancedRepository
enhanced := repository.NewEnhancedRepository[User](handler, logger, "users")

// 第一页
users, cursor, err := enhanced.FindByFieldWithCursor(ctx, "status", "active", "", 20, "id", "ASC")

// 下一页
nextUsers, nextCursor, err := enhanced.FindByFieldWithCursor(ctx, "status", "active", cursor, 20, "id", "ASC")
```

## 性能优化

### 只查询需要的字段

```go
query := repository.NewQuery().
    Select([]string{"id", "name", "email"})  // 减少数据传输
```

### 使用索引字段排序

```go
// 在索引字段上排序
query := repository.NewQuery().AddOrder("created_at", "DESC")  // created_at 有索引
```

### 批量查询优化

```go
// 使用 IN 查询代替多次查询
ids := []interface{}{1, 2, 3, 4, 5}
filter := repository.NewInFilter("id", ids)
users, err := repo.List(ctx, repository.NewQuery().AddFilter(filter))
```

### 预加载优化

```go
// 一次性加载关联数据，避免 N+1 问题
query := repository.NewQuery().
    AddPreload("Profile").
    AddPreload("Orders")
```

## 查询状态检查

Query 提供了一组状态检查方法，用于判断是否已设置各类查询条件：

```go
query := repository.NewQuery().
    AddEqFilter("status", "active").
    AddOrderDesc("created_at").
    WithPaging(1, 20).
    Select("id", "name").
    Limit(50)

query.HasFilters()       // true - 是否有过滤条件（含 FilterGroup）
query.HasPagination()    // true - 是否设置了分页
query.HasOrders()        // true - 是否设置了排序
query.HasSelectFields()  // true - 是否指定了查询字段
query.HasOmitFields()    // false - 是否指定了排除字段
query.HasGroupBy()       // false - 是否设置了分组
query.HasHaving()        // false - 是否设置了 HAVING 条件
query.IsLimited()        // true - 是否设置了 Limit
query.IsOffset()         // false - 是否设置了 Offset
```

**实战示例 - 条件应用中间件：**

```go
func applyDefaults(query *repository.Query) *repository.Query {
    if !query.HasOrders() {
        query.AddOrderDesc("created_at")
    }
    if !query.HasPagination() && !query.IsLimited() {
        query.Limit(100) // 防止全表查询
    }
    return query
}
```

## 查询重置

当需要复用 Query 对象但修改部分条件时，可以使用重置方法：

```go
query := repository.NewQuery().
    AddEqFilter("status", "active").
    AddOrderDesc("created_at").
    WithPaging(1, 20)

// 重置过滤条件（清除 Filters 和 FilterGroup）
query.ResetFilters()

// 重置排序条件
query.ResetOrders()

// 重置分页设置
query.ResetPagination()
```

**实战示例 - 复用查询模板：**

```go
// 基础查询模板
base := repository.NewQuery().
    AddEqFilter("tenant_id", tenantID).
    OmitSensitive().
    AddOrderDesc("created_at")

// 第一次查询：活跃用户
activeUsers, _ := repo.List(ctx, base.Clone().AddEqFilter("status", "active"))

// 第二次查询：需要重新构建过滤条件
query := base.Clone().ResetFilters().
    AddEqFilter("tenant_id", tenantID).
    AddEqFilter("status", "deleted")
deletedUsers, _ := repo.List(ctx, query)
```

## 调试查询

```go
// 启用 SQL 日志（在开发环境）
db.Debug().Model(&User{})...
```

## 📚 相关文档

- 🎯 [FilterGroup 详细文档](./FILTERGROUP.md) - FilterGroup 完整用法
- 🚀 [EnhancedRepository](./ENHANCED-REPOSITORY.md) - 更多便利方法
- 📖 [Repository 基础](./REPOSITORY-BASICS.md) - 基础 CRUD 操作
