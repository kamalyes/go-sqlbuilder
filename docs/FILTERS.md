# 过滤条件（Filters）

Filter 是构建 WHERE 条件的核心，支持各种比较运算符和组合逻辑。

## 📖 目录

- [基本运算符](#基本运算符)
- [范围查询](#范围查询)
- [模糊查询](#模糊查询)
- [NULL 查询](#null-查询)
- [条件组合](#条件组合)
- [动态条件](#动态条件)

---

## 基本运算符

### 等值查询 (=)

```go
// WHERE status = 'active'
query := repository.NewQuery().
    AddEqFilter("status", "active")

users, err := repo.List(ctx, query)
```

### 不等值查询 (!=)

```go
// WHERE status != 'deleted'
query := repository.NewQuery().
    AddNeFilter("status", "deleted")

users, err := repo.List(ctx, query)
```

### 大于 (>)

```go
// WHERE age > 18
query := repository.NewQuery().
    AddGtFilter("age", 18)

users, err := repo.List(ctx, query)
```

### 大于等于 (>=)

```go
// WHERE age >= 18
query := repository.NewQuery().
    AddGteFilter("age", 18)

users, err := repo.List(ctx, query)
```

### 小于 (<)

```go
// WHERE age < 60
query := repository.NewQuery().
    AddLtFilter("age", 60)

users, err := repo.List(ctx, query)
```

### 小于等于 (<=)

```go
// WHERE age <= 60
query := repository.NewQuery().
    AddLteFilter("age", 60)

users, err := repo.List(ctx, query)
```

---

## 范围查询

### BETWEEN - 范围查询

```go
// WHERE age BETWEEN 18 AND 60
query := repository.NewQuery().
    AddBetweenFilter("age", 18, 60)

users, err := repo.List(ctx, query)
```

### IN - 包含查询

```go
// WHERE status IN ('active', 'verified', 'premium')
query := repository.NewQuery().
    AddInFilter("status", []interface{}{"active", "verified", "premium"})

users, err := repo.List(ctx, query)
```

### NOT IN - 不包含查询

```go
// WHERE status NOT IN ('deleted', 'banned')
query := repository.NewQuery().
    AddNotInFilter("status", []interface{}{"deleted", "banned"})

users, err := repo.List(ctx, query)
```

### NewInFilterSlice / NewNotInFilterSlice - 切片版本

当你已经有切片时，可以直接使用切片版本：

```go
// 使用切片参数
statusList := []interface{}{"active", "verified", "premium"}
query := repository.NewQuery().
    AddFilter(repository.NewInFilterSlice("status", statusList))

// NOT IN 切片版本
excludeList := []interface{}{"deleted", "banned", "suspended"}
query := repository.NewQuery().
    AddFilter(repository.NewNotInFilterSlice("status", excludeList))

users, err := repo.List(ctx, query)
```

---

## 模糊查询

### LIKE - 模糊匹配

```go
// WHERE name LIKE '%张三%'
query := repository.NewQuery().
    AddLikeFilter("name", "张三")

users, err := repo.List(ctx, query)
```

**注意：** `AddLikeFilter` 会自动在值两侧添加 `%`，无需手动添加。

### StartsWith - 前缀匹配

```go
// WHERE name LIKE '张%'
query := repository.NewQuery().
    AddFilter(repository.NewStartsWithFilter("name", "张"))

users, err := repo.List(ctx, query)
```

### EndsWith - 后缀匹配

```go
// WHERE email LIKE '%@gmail.com'
query := repository.NewQuery().
    AddFilter(repository.NewEndsWithFilter("email", "@gmail.com"))

users, err := repo.List(ctx, query)
```

### NOT LIKE - 不匹配

```go
// WHERE name NOT LIKE '%test%'
query := repository.NewQuery().
    AddFilter(repository.NewNotLikeFilter("name", "test"))

users, err := repo.List(ctx, query)
```

---

## NULL 查询

### IS NULL

```go
// WHERE deleted_at IS NULL
query := repository.NewQuery().
    AddIsNullFilter("deleted_at")

users, err := repo.List(ctx, query)
```

### IS NOT NULL

```go
// WHERE email IS NOT NULL
query := repository.NewQuery().
    AddIsNotNullFilter("email")

users, err := repo.List(ctx, query)
```

---

## 条件组合

### AND 条件（默认）

多个 Filter 默认使用 AND 连接。

```go
// WHERE status = 'active' AND age > 18 AND email IS NOT NULL
query := repository.NewQuery().
    AddEqFilter("status", "active").
    AddGtFilter("age", 18).
    AddIsNotNullFilter("email")

users, err := repo.List(ctx, query)
```

### OR 条件

使用 FilterGroup 实现 OR 逻辑。

```go
// WHERE status = 'active' OR status = 'verified'
filterGroup := repository.NewFilterGroup(constants.AND_OR).
    AddFilter(repository.NewEqFilter("status", "active")).
    AddFilter(repository.NewEqFilter("status", "verified"))

query := repository.NewQuery().SetFilterGroup(filterGroup)
users, err := repo.List(ctx, query)
```

详见 [FilterGroup 文档](./FILTERGROUP.md)。

---

## 动态条件

### 条件添加（非空值才添加）

```go
func searchUsers(repo repository.IBaseRepository[User, uint64], 
    name, email, status string, ageMin, ageMax int) ([]*User, error) {
    
    ctx := context.Background()
    query := repository.NewQuery()
    
    // 仅当值非空时添加过滤条件
    query.AddEqFilterIfNotEmpty("name", name).
          AddEqFilterIfNotEmpty("email", email).
          AddEqFilterIfNotEmpty("status", status)
    
    // 年龄范围
    if ageMin > 0 {
        query.AddGteFilter("age", ageMin)
    }
    if ageMax > 0 {
        query.AddLteFilter("age", ageMax)
    }
    
    return repo.List(ctx, query)
}
```

### 通用 Filter 构造器

使用 `NewFilter` 创建任意操作符的过滤器。

```go
import "github.com/kamalyes/go-sqlbuilder/constants"

query := repository.NewQuery().
    AddFilter(repository.NewFilter("age", constants.OP_GTE, 18)).
    AddFilter(repository.NewFilter("status", constants.OP_EQ, "active")).
    AddFilter(repository.NewFilter("city", constants.OP_IN, []interface{}{"北京", "上海"}))

users, err := repo.List(ctx, query)
```

**支持的操作符：**

- `OP_EQ` - 等于 (=)
- `OP_NEQ` - 不等于 (!=)
- `OP_GT` - 大于 (>)
- `OP_GTE` - 大于等于 (>=)
- `OP_LT` - 小于 (<)
- `OP_LTE` - 小于等于 (<=)
- `OP_LIKE` - 模糊匹配 (LIKE)
- `OP_NOT_LIKE` - 不匹配 (NOT LIKE)
- `OP_IN` - 包含 (IN)
- `OP_NOT_IN` - 不包含 (NOT IN)
- `OP_BETWEEN` - 范围 (BETWEEN)
- `OP_IS_NULL` - 为空 (IS NULL)
- `OP_IS_NOT_NULL` - 不为空 (IS NOT NULL)

---

## 💡 实战示例

### 示例 1: 用户搜索

```go
func advancedUserSearch(repo repository.IBaseRepository[User, uint64], 
    keyword string, 
    statuses []string, 
    ageRange [2]int, 
    includeDeleted bool) ([]*User, error) {
    
    ctx := context.Background()
    query := repository.NewQuery()
    
    // 关键词搜索（姓名或邮箱）
    if keyword != "" {
        keywordGroup := repository.NewFilterGroup(constants.AND_OR).
            AddFilter(repository.NewLikeFilter("name", keyword)).
            AddFilter(repository.NewLikeFilter("email", keyword))
        query.SetFilterGroup(keywordGroup)
    }
    
    // 状态筛选
    if len(statuses) > 0 {
        statusValues := make([]interface{}, len(statuses))
        for i, s := range statuses {
            statusValues[i] = s
        }
        query.AddInFilter("status", statusValues)
    }
    
    // 年龄范围
    if ageRange[0] > 0 {
        query.AddGteFilter("age", ageRange[0])
    }
    if ageRange[1] > 0 {
        query.AddLteFilter("age", ageRange[1])
    }
    
    // 是否包含已删除用户
    if !includeDeleted {
        query.AddIsNullFilter("deleted_at")
    }
    
    return repo.List(ctx, query)
}
```

### 示例 2: 商品筛选

```go
func filterProducts(repo repository.IBaseRepository[Product, uint64],
    categoryIDs []int,
    priceRange [2]float64,
    keyword string,
    inStock bool) ([]*Product, error) {
    
    ctx := context.Background()
    query := repository.NewQuery()
    
    // 分类筛选
    if len(categoryIDs) > 0 {
        categoryValues := make([]interface{}, len(categoryIDs))
        for i, id := range categoryIDs {
            categoryValues[i] = id
        }
        query.AddInFilter("category_id", categoryValues)
    }
    
    // 价格范围
    if priceRange[0] > 0 || priceRange[1] > 0 {
        if priceRange[0] > 0 && priceRange[1] > 0 {
            query.AddBetweenFilter("price", priceRange[0], priceRange[1])
        } else if priceRange[0] > 0 {
            query.AddGteFilter("price", priceRange[0])
        } else {
            query.AddLteFilter("price", priceRange[1])
        }
    }
    
    // 关键词搜索
    if keyword != "" {
        query.AddLikeFilter("name", keyword)
    }
    
    // 库存筛选
    if inStock {
        query.AddGtFilter("stock", 0)
    }
    
    return repo.List(ctx, query)
}
```

### 示例 3: 日期范围查询

```go
func getOrdersByDateRange(repo repository.IBaseRepository[Order, uint64],
    startDate, endDate time.Time,
    status string) ([]*Order, error) {
    
    ctx := context.Background()
    query := repository.NewQuery()
    
    // 日期范围
    if !startDate.IsZero() && !endDate.IsZero() {
        query.AddBetweenFilter("created_at", startDate, endDate)
    } else if !startDate.IsZero() {
        query.AddGteFilter("created_at", startDate)
    } else if !endDate.IsZero() {
        query.AddLteFilter("created_at", endDate)
    }
    
    // 状态筛选
    query.AddEqFilterIfNotEmpty("status", status)
    
    return repo.List(ctx, query)
}
```

---

## ⚠️ 注意事项

### 1. LIKE 查询自动添加通配符

```go
// ❌ 错误：手动添加 %
query.AddLikeFilter("name", "%张三%") // 会变成 LIKE '%%张三%%'

// ✅ 正确：自动添加
query.AddLikeFilter("name", "张三") // 会变成 LIKE '%张三%'
```

### 2. IN 查询参数类型

```go
// ✅ 正确：使用 []interface{}
query.AddInFilter("status", []interface{}{"active", "verified"})

// ❌ 错误：直接传递字符串切片可能导致类型问题
statuses := []string{"active", "verified"}
query.AddInFilter("status", statuses) // 可能出错
```

### 3. NULL 查询不需要值

```go
// ✅ 正确
query.AddIsNullFilter("deleted_at")

// ❌ 错误：不需要传值
query.AddEqFilter("deleted_at", nil) // 可能不生效
```

---

## 📚 相关文档

- [FilterGroup](./FILTERGROUP.md) - 复杂条件组合（AND/OR）
- [便捷查询方法](./CONVENIENCE-METHODS.md) - 简化的查询 API
- [时间查询](./TIME-QUERIES.md) - 时间相关的快捷查询
- [高级查询](./ADVANCED-QUERIES.md) - 子查询、JOIN 等
