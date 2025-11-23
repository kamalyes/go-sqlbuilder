# 高级查询

本文档介绍如何使用 Query、Filter 和 FilterGroup 构建复杂查询。

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

## 调试查询

```go
// 启用 SQL 日志（在开发环境）
db.Debug().Model(&User{})...
```

## 下一步

- 🎯 [FilterGroup 详细文档](./FILTERGROUP.MD) - FilterGroup 完整用法
- 🚀 [EnhancedRepository](./ENHANCED-REPOSITORY.MD) - 更多便利方法
- 📖 [Repository 基础](./REPOSITORY-BASICS.MD) - 基础 CRUD 操作
