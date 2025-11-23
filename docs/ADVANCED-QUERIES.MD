# 高级查询

本文档介绍如何使用 Query、Filter 和 FilterGroup 构建复杂查询。

## Query 对象

Query 是查询条件的容器，支持过滤、排序、分页等功能。

### 基础用法

```go
query := sqlbuilder.NewQuery().
    AddFilter(sqlbuilder.NewEqFilter("status", "active")).
    AddOrder("created_at", "DESC").
    Limit(10).
    Offset(20)
    
users, err := repo.List(ctx, query)
```

### 预加载关联

```go
query := sqlbuilder.NewQuery().
    AddPreload("Profile").
    AddPreload("Orders.Items")  // 嵌套预加载
    
users, err := repo.List(ctx, query)
```

### 字段选择

```go
// 只查询指定字段
query := sqlbuilder.NewQuery().
    Select([]string{"id", "name", "email"})
    
// 排除指定字段
query := sqlbuilder.NewQuery().
    Omit([]string{"password", "secret_key"})
```

### 分组和聚合

```go
query := sqlbuilder.NewQuery().
    GroupBy("status").
    Having("COUNT(*) > ?", 10)
```

## Filter 过滤器

Filter 提供了丰富的查询条件构建器。

### 等值查询

```go
// 等于
filter := sqlbuilder.NewEqFilter("status", "active")

// 不等于
filter := sqlbuilder.NewNeFilter("status", "deleted")
```

### 比较查询

```go
// 大于
filter := sqlbuilder.NewGtFilter("age", 18)

// 大于等于
filter := sqlbuilder.NewGteFilter("age", 18)

// 小于
filter := sqlbuilder.NewLtFilter("price", 100)

// 小于等于
filter := sqlbuilder.NewLteFilter("price", 100)
```

### 范围查询

```go
// IN 查询
filter := sqlbuilder.NewInFilter("id", []interface{}{1, 2, 3, 4, 5})

// NOT IN
filter := sqlbuilder.NewNotInFilter("status", []interface{}{"deleted", "banned"})

// BETWEEN
filter := sqlbuilder.NewBetweenFilter("age", 18, 65)

// NOT BETWEEN
filter := sqlbuilder.NewNotBetweenFilter("price", 100, 200)
```

### 模糊查询

```go
// LIKE '%keyword%'
filter := sqlbuilder.NewLikeFilter("name", "张")

// NOT LIKE
filter := sqlbuilder.NewNotLikeFilter("email", "@temp.com")

// 自定义 LIKE 模式
filter := sqlbuilder.NewLikeFilter("phone", "138%")  // 以 138 开头
filter := sqlbuilder.NewLikeFilter("email", "%@gmail.com")  // Gmail 邮箱
```

### NULL 查询

```go
// IS NULL
filter := sqlbuilder.NewIsNullFilter("deleted_at")

// IS NOT NULL
filter := sqlbuilder.NewIsNotNullFilter("verified_at")
```

### 时间范围查询

```go
// 今天
filter := sqlbuilder.NewTodayFilter("created_at")

// 本周
filter := sqlbuilder.NewThisWeekFilter("created_at")

// 本月
filter := sqlbuilder.NewThisMonthFilter("created_at")

// 本年
filter := sqlbuilder.NewThisYearFilter("created_at")

// 昨天
filter := sqlbuilder.NewYesterdayFilter("created_at")

// 上周
filter := sqlbuilder.NewLastWeekFilter("created_at")

// 上月
filter := sqlbuilder.NewLastMonthFilter("created_at")
```

### 自定义 SQL

```go
// 自定义条件
filter := sqlbuilder.NewCustomFilter("YEAR(created_at) = ?", 2024)

// 复杂条件
filter := sqlbuilder.NewCustomFilter(
    "price * quantity > ? AND discount < ?",
    1000, 0.5,
)
```

## FilterGroup - 复杂条件组合

FilterGroup 支持 AND/OR 逻辑组合和嵌套。

### 基础 AND 条件

```go
// (status = 'active' AND age > 18)
group := sqlbuilder.NewFilterGroup(sqlbuilder.AND).
    AddFilter(sqlbuilder.NewEqFilter("status", "active")).
    AddFilter(sqlbuilder.NewGtFilter("age", 18))
    
query := sqlbuilder.NewQuery().SetFilterGroup(group)
```

### 基础 OR 条件

```go
// (status = 'active' OR status = 'verified')
group := sqlbuilder.NewFilterGroup(sqlbuilder.OR).
    AddFilter(sqlbuilder.NewEqFilter("status", "active")).
    AddFilter(sqlbuilder.NewEqFilter("status", "verified"))
```

### 嵌套条件

```go
// (status = 'active' AND (age > 18 OR vip_level > 0))
innerGroup := sqlbuilder.NewFilterGroup(sqlbuilder.OR).
    AddFilter(sqlbuilder.NewGtFilter("age", 18)).
    AddFilter(sqlbuilder.NewGtFilter("vip_level", 0))
    
outerGroup := sqlbuilder.NewFilterGroup(sqlbuilder.AND).
    AddFilter(sqlbuilder.NewEqFilter("status", "active")).
    AddGroup(innerGroup)
    
query := sqlbuilder.NewQuery().SetFilterGroup(outerGroup)
```

### 复杂示例：电商搜索

```go
// ((category = 'electronics' OR category = 'books') 
//  AND price BETWEEN 100 AND 1000 
//  AND (brand = 'Apple' OR rating >= 4.5))

categoryGroup := sqlbuilder.NewFilterGroup(sqlbuilder.OR).
    AddFilter(sqlbuilder.NewEqFilter("category", "electronics")).
    AddFilter(sqlbuilder.NewEqFilter("category", "books"))
    
brandGroup := sqlbuilder.NewFilterGroup(sqlbuilder.OR).
    AddFilter(sqlbuilder.NewEqFilter("brand", "Apple")).
    AddFilter(sqlbuilder.NewGteFilter("rating", 4.5))
    
mainGroup := sqlbuilder.NewFilterGroup(sqlbuilder.AND).
    AddGroup(categoryGroup).
    AddFilter(sqlbuilder.NewBetweenFilter("price", 100, 1000)).
    AddGroup(brandGroup)
    
query := sqlbuilder.NewQuery().
    SetFilterGroup(mainGroup).
    AddOrder("price", "ASC")
    
products, err := repo.List(ctx, query)
```

### 复杂示例：用户筛选

```go
// ((status = 'active' AND verified = true)
//  OR (status = 'trial' AND created_at > 7 days ago))
//  AND (age > 18 OR is_student = true)

activeGroup := sqlbuilder.NewFilterGroup(sqlbuilder.AND).
    AddFilter(sqlbuilder.NewEqFilter("status", "active")).
    AddFilter(sqlbuilder.NewEqFilter("verified", true))
    
trialGroup := sqlbuilder.NewFilterGroup(sqlbuilder.AND).
    AddFilter(sqlbuilder.NewEqFilter("status", "trial")).
    AddFilter(sqlbuilder.NewGtFilter("created_at", time.Now().AddDate(0, 0, -7)))
    
statusGroup := sqlbuilder.NewFilterGroup(sqlbuilder.OR).
    AddGroup(activeGroup).
    AddGroup(trialGroup)
    
ageGroup := sqlbuilder.NewFilterGroup(sqlbuilder.OR).
    AddFilter(sqlbuilder.NewGtFilter("age", 18)).
    AddFilter(sqlbuilder.NewEqFilter("is_student", true))
    
finalGroup := sqlbuilder.NewFilterGroup(sqlbuilder.AND).
    AddGroup(statusGroup).
    AddGroup(ageGroup)
    
query := sqlbuilder.NewQuery().SetFilterGroup(finalGroup)
users, err := repo.List(ctx, query)
```

## 排序

### 单字段排序

```go
query := sqlbuilder.NewQuery().AddOrder("created_at", "DESC")
```

### 多字段排序

```go
query := sqlbuilder.NewQuery().
    AddOrder("status", "ASC").
    AddOrder("created_at", "DESC")
```

### 使用常量

```go
query := sqlbuilder.NewQuery().
    AddOrder("priority", sqlbuilder.Desc).
    AddOrder("created_at", sqlbuilder.Asc)
```

### 自定义排序

```go
query := sqlbuilder.NewQuery().
    AddOrder("FIELD(status, 'urgent', 'high', 'normal', 'low')", "")
```

## 分页

### 基础分页

```go
query := sqlbuilder.NewQuery().
    Limit(20).
    Offset(40)  // 第 3 页，每页 20 条
```

### 使用 Paging 对象

```go
paging := &sqlbuilder.Paging{
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
enhanced := sqlbuilder.NewEnhancedRepository[User](handler, logger, "users")

// 第一页
users, cursor, err := enhanced.FindByFieldWithCursor(ctx, "status", "active", "", 20, "id", "ASC")

// 下一页
nextUsers, nextCursor, err := enhanced.FindByFieldWithCursor(ctx, "status", "active", cursor, 20, "id", "ASC")
```

## 性能优化

### 只查询需要的字段

```go
query := sqlbuilder.NewQuery().
    Select([]string{"id", "name", "email"})  // 减少数据传输
```

### 使用索引字段排序

```go
// 在索引字段上排序
query := sqlbuilder.NewQuery().AddOrder("created_at", "DESC")  // created_at 有索引
```

### 批量查询优化

```go
// 使用 IN 查询代替多次查询
ids := []interface{}{1, 2, 3, 4, 5}
filter := sqlbuilder.NewInFilter("id", ids)
users, err := repo.List(ctx, sqlbuilder.NewQuery().AddFilter(filter))
```

### 预加载优化

```go
// 一次性加载关联数据，避免 N+1 问题
query := sqlbuilder.NewQuery().
    AddPreload("Profile").
    AddPreload("Orders")
```

## 调试查询

```go
// 启用 SQL 日志（在开发环境）
db.Debug().Model(&User{})...
```

## 下一步

- 🎯 [FilterGroup 详细文档](./FILTERGROUP.MD) - FilterGroup 完整用法
- 🚀 [EnhancedRepository](./ENHANCED-REPOSITORY.MD) - 更多便利方法
- 📖 [Repository 基础](./REPOSITORY-BASICS.MD) - 基础 CRUD 操作
