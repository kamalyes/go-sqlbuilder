# Query 便捷方法使用示例

本文档展示了 Query 新增便捷方法的实际使用场景和最佳实践。

## 🚀 快速上手

### 基础链式调用

```go
// 传统方式 vs 便捷方式
// ❌ 传统方式（较繁琐）
query := repository.NewQuery().
    AddFilter(repository.NewEqFilter("status", 1)).
    AddFilter(repository.NewLikeFilter("name", "%test%")).
    AddFilter(repository.NewGtFilter("created_at", time.Now().AddDate(0, -1, 0))).
    AddFilter(repository.NewInFilter("category_id", []interface{}{1, 2, 3})).
    AddOrder("created_at", "DESC").
    Limit(10)

// ✅ 便捷方式（推荐）
query := repository.NewQuery().
    AddEqual("status", 1).
    AddLike("name", "test").
    AddTimeAfter("created_at", time.Now().AddDate(0, -1, 0)).
    AddIn("category_id", 1, 2, 3).
    AddOrderDesc("created_at").
    Take(10)
```

## 📝 实际业务场景

### 用户管理系统

```go
// 查询活跃用户
func GetActiveUsers(repo *repository.BaseRepository[User], ctx context.Context) ([]*User, error) {
    return repo.List(ctx, repository.NewQuery().
        AddEqual("status", "active").
        AddIsNotNull("email").
        AddThisMonth("last_login_at").
        AddOrderDesc("last_login_at").
        Take(100))
}

// 搜索用户
func SearchUsers(repo *repository.BaseRepository[User], ctx context.Context, 
    keyword string, ageMin, ageMax int, statuses []string) ([]*User, error) {
    
    query := repository.NewQuery()
    
    // 条件式添加查询条件
    if keyword != "" {
        query.AddLike("name", keyword).
              AddLike("email", keyword)  // 多个 LIKE 条件会用 AND 连接
    }
    
    if ageMin > 0 {
        query.AddGreaterEqual("age", ageMin)
    }
    
    if ageMax > 0 {
        query.AddLessEqual("age", ageMax)
    }
    
    if len(statuses) > 0 {
        // 将 []string 转换为 []interface{}
        interfaces := make([]interface{}, len(statuses))
        for i, s := range statuses {
            interfaces[i] = s
        }
        query.AddIn("status", interfaces...)
    }
    
    return repo.List(ctx, query.AddOrderDesc("created_at"))
}
```

### 电商订单系统

```go
// 订单统计查询
func GetOrderStats(repo *repository.BaseRepository[Order], ctx context.Context) ([]*Order, error) {
    startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
    endDate := time.Date(2025, 12, 31, 23, 59, 59, 0, time.UTC)

    return repo.List(ctx, repository.NewQuery().
        AddEqual("tenant_id", 1001).                              // 租户过滤
        AddIn("status", "pending", "processing", "shipped").      // 状态过滤
        AddGreaterEqual("amount", 100).                           // 最小金额
        AddLessEqual("amount", 50000).                            // 最大金额
        AddStartsWith("order_no", "ORD2025").                     // 订单号前缀
        AddIsNotNull("customer_email").                           // 必须有邮箱
        AddTimeBetween("created_at", startDate, endTime).         // 时间范围
        AddLike("shipping_address", "北京").                       // 地址包含
        AddNotIn("payment_method", "cash_on_delivery", "check").  // 排除支付方式
        AddOrderDesc("created_at").                               // 按创建时间降序
        AddOrderAsc("amount"))                                    // 按金额升序
}

// 今日订单
func GetTodayOrders(repo *repository.BaseRepository[Order], ctx context.Context) ([]*Order, error) {
    return repo.List(ctx, repository.NewQuery().
        AddToday("created_at").
        AddNotEqual("status", "cancelled").
        AddOrderDesc("created_at"))
}

// 本月待处理订单
func GetPendingOrdersThisMonth(repo *repository.BaseRepository[Order], ctx context.Context) ([]*Order, error) {
    return repo.List(ctx, repository.NewQuery().
        AddThisMonth("created_at").
        AddIn("status", "pending", "processing").
        AddOrderAsc("created_at"))
}
```

### 内容管理系统

```go
// 文章搜索
func SearchArticles(repo *repository.BaseRepository[Article], ctx context.Context, 
    title, content string, categoryIds []int, published bool) ([]*Article, error) {
    
    query := repository.NewQuery()
    
    // 标题搜索
    if title != "" {
        query.AddLike("title", title)
    }
    
    // 内容搜索
    if content != "" {
        query.AddLike("content", content)
    }
    
    // 分类过滤
    if len(categoryIds) > 0 {
        interfaces := make([]interface{}, len(categoryIds))
        for i, id := range categoryIds {
            interfaces[i] = id
        }
        query.AddIn("category_id", interfaces...)
    }
    
    // 发布状态
    query.AddEqual("published", published)
    
    // 排序：置顶文章优先，然后按发布时间降序
    return repo.List(ctx, query.
        AddOrderDesc("is_top").
        AddOrderDesc("published_at"))
}

// 热门文章（本周）
func GetPopularArticlesThisWeek(repo *repository.BaseRepository[Article], ctx context.Context) ([]*Article, error) {
    return repo.List(ctx, repository.NewQuery().
        AddThisWeek("created_at").
        AddEqual("published", true).
        AddGreaterThan("view_count", 100).
        AddOrderDesc("view_count").
        Take(20))
}
```

### 日志分析系统

```go
// 错误日志查询
func GetErrorLogs(repo *repository.BaseRepository[Log], ctx context.Context, 
    level string, hours int) ([]*Log, error) {
    
    startTime := time.Now().Add(-time.Duration(hours) * time.Hour)
    
    return repo.List(ctx, repository.NewQuery().
        AddEqual("level", level).
        AddTimeAfter("created_at", startTime).
        AddIsNotNull("error_message").
        AddOrderDesc("created_at"))
}

// 今日系统监控
func GetTodaySystemLogs(repo *repository.BaseRepository[Log], ctx context.Context) ([]*Log, error) {
    return repo.List(ctx, repository.NewQuery().
        AddToday("created_at").
        AddIn("level", "ERROR", "WARN", "FATAL").
        AddOrderDesc("created_at"))
}
```

## 🕐 时间查询专题

### 时间便捷方法

```go
// 时间范围查询
query := repository.NewQuery().
    AddTimeAfter("created_at", time.Now().AddDate(0, 0, -7)).  // 7天前之后
    AddTimeBefore("updated_at", time.Now())                     // 现在之前

// 时间段查询
start := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
end := time.Date(2025, 3, 31, 23, 59, 59, 0, time.UTC)
query := repository.NewQuery().
    AddTimeBetween("created_at", start, end)  // Q1季度

// 相对时间查询
query := repository.NewQuery().
    AddToday("login_at")        // 今天登录的用户
    
query := repository.NewQuery().
    AddThisWeek("order_date")   // 本周的订单
    
query := repository.NewQuery().
    AddThisMonth("register_date") // 本月注册的用户
    
query := repository.NewQuery().
    AddThisYear("created_at")   // 今年创建的记录
```

### 报表查询示例

```go
// 月度销售报表
func GetMonthlySalesReport(repo *repository.BaseRepository[Order], ctx context.Context, 
    year int, month time.Month) ([]*Order, error) {
    
    start := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
    end := start.AddDate(0, 1, -1).Add(23*time.Hour + 59*time.Minute + 59*time.Second)
    
    return repo.List(ctx, repository.NewQuery().
        AddTimeBetween("created_at", start, end).
        AddIn("status", "completed", "shipped").
        AddGreaterThan("amount", 0).
        AddOrderDesc("amount"))
}

// 季度用户增长分析
func GetQuarterlyUserGrowth(repo *repository.BaseRepository[User], ctx context.Context, 
    year int, quarter int) ([]*User, error) {
    
    var start, end time.Time
    switch quarter {
    case 1:
        start = time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
        end = time.Date(year, 3, 31, 23, 59, 59, 0, time.UTC)
    case 2:
        start = time.Date(year, 4, 1, 0, 0, 0, 0, time.UTC)
        end = time.Date(year, 6, 30, 23, 59, 59, 0, time.UTC)
    // ... 其他季度
    }
    
    return repo.List(ctx, repository.NewQuery().
        AddTimeBetween("created_at", start, end).
        AddEqual("status", "active").
        AddOrderAsc("created_at"))
}
```

## 🔍 分页和排序

### 智能分页

```go
// 简单分页
func GetUsersList(repo *repository.BaseRepository[User], ctx context.Context, 
    page, pageSize int) ([]*User, *repository.Pagination, error) {
    
    query := repository.NewQuery().
        AddEqual("status", "active").
        AddOrderDesc("created_at").
        Page(page, pageSize)  // 自动处理边界值
    
    return repo.ListWithPagination(ctx, query, nil)
}

// 游标分页（性能更好）
func GetUsersWithCursor(repo *repository.BaseRepository[User], ctx context.Context, 
    lastId uint, limit int) ([]*User, error) {
    
    query := repository.NewQuery().
        AddEqual("status", "active")
    
    if lastId > 0 {
        query.AddGreaterThan("id", lastId)  // 游标分页
    }
    
    return repo.List(ctx, query.
        AddOrderAsc("id").
        Take(limit))
}
```

### 复杂排序

```go
// 多字段排序
func GetProductsByPopularity(repo *repository.BaseRepository[Product], ctx context.Context) ([]*Product, error) {
    return repo.List(ctx, repository.NewQuery().
        AddEqual("status", "active").
        AddOrderDesc("is_featured").     // 推荐商品优先
        AddOrderDesc("sales_count").     // 销量降序
        AddOrderAsc("price").            // 价格升序
        AddOrderDesc("created_at"))      // 最新优先
}
```

## 🎯 性能优化提示

### 避免 N+1 查询

```go
// ✅ 使用预加载
query := repository.NewQuery().
    AddEqual("status", "active").
    // 预加载关联数据
    AddPreload("Profile").
    AddPreload("Orders")

users, err := repo.ListWithPreloads(ctx, query, "Profile", "Orders")
```

### 只查询需要的字段

```go
// ✅ 只查询必要字段
query := repository.NewQuery().
    AddEqual("status", "active").
    Select([]string{"id", "name", "email"})  // 减少数据传输

users, err := repo.List(ctx, query)
```

### 使用索引友好的查询

```go
// ✅ 在索引字段上进行范围查询
query := repository.NewQuery().
    AddGreaterEqual("created_at", startTime).   // created_at 有索引
    AddLessEqual("created_at", endTime).
    AddOrderDesc("created_at")                  // 利用索引排序
```

## 🧪 测试示例

```go
func TestQueryBuilder(t *testing.T) {
    // 测试链式调用
    query := repository.NewQuery().
        AddEqual("status", 1).
        AddLike("name", "test").
        AddIn("category_id", 1, 2, 3).
        AddOrderDesc("created_at").
        Take(10)
    
    // 验证过滤条件
    assert.Equal(t, 3, len(query.Filters))
    assert.Equal(t, "status", query.Filters[0].Field)
    assert.Equal(t, constants.OP_EQ, query.Filters[0].Operator)
    
    // 验证排序
    assert.Equal(t, 1, len(query.Orders))
    assert.Equal(t, "created_at", query.Orders[0].Field)
    assert.Equal(t, "DESC", query.Orders[0].Direction)
    
    // 验证限制
    assert.NotNil(t, query.LimitValue)
    assert.Equal(t, 10, *query.LimitValue)
}
```

## 📚 更多资源

- 🔍 [高级查询文档](./ADVANCED-QUERIES.MD) - 完整的查询功能
- 🎯 [FilterGroup 使用](./FILTERGROUP.MD) - 复杂条件组合
- 🚀 [EnhancedRepository](./ENHANCED-REPOSITORY.MD) - 增强功能
- 📖 [API 参考](./API-REFERENCE.MD) - 完整的 API 文档