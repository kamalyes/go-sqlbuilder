# FilterGroup 完整指南

FilterGroup 是构建复杂 WHERE 条件的强大工具，支持 AND/OR 逻辑组合和无限层级嵌套。

## 基础概念

### AND 逻辑

所有条件必须同时满足（交集）。

```go
// WHERE status = 'active' AND age > 18
group := sqlbuilder.NewFilterGroup(sqlbuilder.AND).
    AddFilter(sqlbuilder.NewEqFilter("status", "active")).
    AddFilter(sqlbuilder.NewGtFilter("age", 18))
```

### OR 逻辑

任一条件满足即可（并集）。

```go
// WHERE status = 'active' OR status = 'verified'
group := sqlbuilder.NewFilterGroup(sqlbuilder.OR).
    AddFilter(sqlbuilder.NewEqFilter("status", "active")).
    AddFilter(sqlbuilder.NewEqFilter("status", "verified"))
```

## 完整 API

### NewFilterGroup

创建过滤器组。

```go
// AND 组
andGroup := sqlbuilder.NewFilterGroup(sqlbuilder.AND)

// OR 组
orGroup := sqlbuilder.NewFilterGroup(sqlbuilder.OR)
```

### AddFilter

添加单个过滤条件。

```go
group.AddFilter(sqlbuilder.NewEqFilter("status", "active"))
```

### AddFilters

批量添加多个过滤条件。

```go
filters := []*sqlbuilder.Filter{
    sqlbuilder.NewEqFilter("status", "active"),
    sqlbuilder.NewGtFilter("age", 18),
    sqlbuilder.NewLikeFilter("name", "张"),
}
group.AddFilters(filters...)
```

### AddGroup

添加嵌套的子组。

```go
subGroup := sqlbuilder.NewFilterGroup(sqlbuilder.OR).
    AddFilter(sqlbuilder.NewEqFilter("role", "admin")).
    AddFilter(sqlbuilder.NewEqFilter("role", "moderator"))
    
mainGroup := sqlbuilder.NewFilterGroup(sqlbuilder.AND).
    AddGroup(subGroup).
    AddFilter(sqlbuilder.NewEqFilter("status", "active"))
```

### Apply

将 FilterGroup 应用到 GORM 查询。

```go
db := handler.GetDB()
query := group.Apply(db)
```

## 实战示例

### 示例 1：电商商品筛选

**需求**：查找价格在 100-1000 之间的电子产品或图书，且评分大于 4.0 或销量大于 1000。

```go
// (category IN ('electronics', 'books') AND price BETWEEN 100 AND 1000)
// AND (rating > 4.0 OR sales > 1000)

categoryFilter := sqlbuilder.NewInFilter("category", []interface{}{"electronics", "books"})
priceFilter := sqlbuilder.NewBetweenFilter("price", 100, 1000)

qualityGroup := sqlbuilder.NewFilterGroup(sqlbuilder.OR).
    AddFilter(sqlbuilder.NewGtFilter("rating", 4.0)).
    AddFilter(sqlbuilder.NewGtFilter("sales", 1000))

mainGroup := sqlbuilder.NewFilterGroup(sqlbuilder.AND).
    AddFilter(categoryFilter).
    AddFilter(priceFilter).
    AddGroup(qualityGroup)

query := sqlbuilder.NewQuery().SetFilterGroup(mainGroup)
products, err := repo.List(ctx, query)
```

### 示例 2：用户权限筛选

**需求**：查找活跃的管理员或版主，且年龄大于 18 岁或已验证邮箱。

```go
// ((role = 'admin' OR role = 'moderator') AND status = 'active')
// AND (age > 18 OR email_verified = true)

roleGroup := sqlbuilder.NewFilterGroup(sqlbuilder.OR).
    AddFilter(sqlbuilder.NewEqFilter("role", "admin")).
    AddFilter(sqlbuilder.NewEqFilter("role", "moderator"))

userGroup := sqlbuilder.NewFilterGroup(sqlbuilder.AND).
    AddGroup(roleGroup).
    AddFilter(sqlbuilder.NewEqFilter("status", "active"))

verificationGroup := sqlbuilder.NewFilterGroup(sqlbuilder.OR).
    AddFilter(sqlbuilder.NewGtFilter("age", 18)).
    AddFilter(sqlbuilder.NewEqFilter("email_verified", true))

finalGroup := sqlbuilder.NewFilterGroup(sqlbuilder.AND).
    AddGroup(userGroup).
    AddGroup(verificationGroup)

query := sqlbuilder.NewQuery().SetFilterGroup(finalGroup)
users, err := repo.List(ctx, query)
```

### 示例 3：订单状态查询

**需求**：查找待支付超过 30 分钟的订单，或已支付但未发货超过 24 小时的订单。

```go
// (status = 'pending' AND created_at < NOW() - 30分钟)
// OR (status = 'paid' AND paid_at < NOW() - 24小时 AND shipped_at IS NULL)

thirtyMinAgo := time.Now().Add(-30 * time.Minute)
oneDayAgo := time.Now().Add(-24 * time.Hour)

pendingGroup := sqlbuilder.NewFilterGroup(sqlbuilder.AND).
    AddFilter(sqlbuilder.NewEqFilter("status", "pending")).
    AddFilter(sqlbuilder.NewLtFilter("created_at", thirtyMinAgo))

paidGroup := sqlbuilder.NewFilterGroup(sqlbuilder.AND).
    AddFilter(sqlbuilder.NewEqFilter("status", "paid")).
    AddFilter(sqlbuilder.NewLtFilter("paid_at", oneDayAgo)).
    AddFilter(sqlbuilder.NewIsNullFilter("shipped_at"))

mainGroup := sqlbuilder.NewFilterGroup(sqlbuilder.OR).
    AddGroup(pendingGroup).
    AddGroup(paidGroup)

query := sqlbuilder.NewQuery().SetFilterGroup(mainGroup)
orders, err := repo.List(ctx, query)
```

### 示例 4：内容审核查询

**需求**：查找需要人工审核的内容（包含敏感词或被举报次数大于 3，且未被审核）。

```go
// (has_sensitive_words = true OR report_count > 3)
// AND status = 'pending_review'
// AND reviewed_at IS NULL

flagGroup := sqlbuilder.NewFilterGroup(sqlbuilder.OR).
    AddFilter(sqlbuilder.NewEqFilter("has_sensitive_words", true)).
    AddFilter(sqlbuilder.NewGtFilter("report_count", 3))

mainGroup := sqlbuilder.NewFilterGroup(sqlbuilder.AND).
    AddGroup(flagGroup).
    AddFilter(sqlbuilder.NewEqFilter("status", "pending_review")).
    AddFilter(sqlbuilder.NewIsNullFilter("reviewed_at"))

query := sqlbuilder.NewQuery().SetFilterGroup(mainGroup)
contents, err := repo.List(ctx, query)
```

### 示例 5：多条件搜索

**需求**：搜索名字或邮箱包含关键词，且状态为活跃或试用中的用户。

```go
// (name LIKE '%keyword%' OR email LIKE '%keyword%')
// AND (status = 'active' OR status = 'trial')

keyword := "张三"

searchGroup := sqlbuilder.NewFilterGroup(sqlbuilder.OR).
    AddFilter(sqlbuilder.NewLikeFilter("name", keyword)).
    AddFilter(sqlbuilder.NewLikeFilter("email", keyword))

statusGroup := sqlbuilder.NewFilterGroup(sqlbuilder.OR).
    AddFilter(sqlbuilder.NewEqFilter("status", "active")).
    AddFilter(sqlbuilder.NewEqFilter("status", "trial"))

mainGroup := sqlbuilder.NewFilterGroup(sqlbuilder.AND).
    AddGroup(searchGroup).
    AddGroup(statusGroup)

query := sqlbuilder.NewQuery().SetFilterGroup(mainGroup)
users, err := repo.List(ctx, query)
```

### 示例 6：时间范围组合

**需求**：查找本周创建的订单，或上月更新过的高优先级订单。

```go
// (created_at >= THIS_WEEK_START AND created_at < THIS_WEEK_END)
// OR (updated_at >= LAST_MONTH_START AND updated_at < LAST_MONTH_END AND priority = 'high')

thisWeekGroup := sqlbuilder.NewFilterGroup(sqlbuilder.AND).
    AddFilter(sqlbuilder.NewThisWeekFilter("created_at"))

lastMonthGroup := sqlbuilder.NewFilterGroup(sqlbuilder.AND).
    AddFilter(sqlbuilder.NewLastMonthFilter("updated_at")).
    AddFilter(sqlbuilder.NewEqFilter("priority", "high"))

mainGroup := sqlbuilder.NewFilterGroup(sqlbuilder.OR).
    AddGroup(thisWeekGroup).
    AddGroup(lastMonthGroup)

query := sqlbuilder.NewQuery().SetFilterGroup(mainGroup)
orders, err := repo.List(ctx, query)
```

### 示例 7：多级嵌套

**需求**：复杂的会员筛选（VIP 或高消费用户，且活跃或最近登录）。

```go
// ((vip_level > 0 OR total_spent > 10000) AND status = 'active')
// AND (last_login > 7天前 OR login_count_this_month > 10)

memberGroup := sqlbuilder.NewFilterGroup(sqlbuilder.OR).
    AddFilter(sqlbuilder.NewGtFilter("vip_level", 0)).
    AddFilter(sqlbuilder.NewGtFilter("total_spent", 10000))

memberStatusGroup := sqlbuilder.NewFilterGroup(sqlbuilder.AND).
    AddGroup(memberGroup).
    AddFilter(sqlbuilder.NewEqFilter("status", "active"))

activityGroup := sqlbuilder.NewFilterGroup(sqlbuilder.OR).
    AddFilter(sqlbuilder.NewGtFilter("last_login", time.Now().AddDate(0, 0, -7))).
    AddFilter(sqlbuilder.NewGtFilter("login_count_this_month", 10))

finalGroup := sqlbuilder.NewFilterGroup(sqlbuilder.AND).
    AddGroup(memberStatusGroup).
    AddGroup(activityGroup)

query := sqlbuilder.NewQuery().SetFilterGroup(finalGroup)
members, err := repo.List(ctx, query)
```

## 最佳实践

### 1. 命名清晰的变量

```go
// ❌ 不好
g1 := sqlbuilder.NewFilterGroup(sqlbuilder.OR)
g2 := sqlbuilder.NewFilterGroup(sqlbuilder.AND)

// ✅ 好
roleGroup := sqlbuilder.NewFilterGroup(sqlbuilder.OR)
statusGroup := sqlbuilder.NewFilterGroup(sqlbuilder.AND)
```

### 2. 逐步构建复杂条件

```go
// 先构建子组
categoryGroup := sqlbuilder.NewFilterGroup(sqlbuilder.OR).
    AddFilter(sqlbuilder.NewEqFilter("category", "A")).
    AddFilter(sqlbuilder.NewEqFilter("category", "B"))

priceGroup := sqlbuilder.NewFilterGroup(sqlbuilder.AND).
    AddFilter(sqlbuilder.NewGteFilter("price", 100)).
    AddFilter(sqlbuilder.NewLteFilter("price", 1000))

// 再组合
finalGroup := sqlbuilder.NewFilterGroup(sqlbuilder.AND).
    AddGroup(categoryGroup).
    AddGroup(priceGroup)
```

### 3. 添加注释说明逻辑

```go
// 查找需要提醒的用户：
// - 试用期即将到期（7天内）
// - 或者付费会员即将到期（30天内）
trialExpiringGroup := sqlbuilder.NewFilterGroup(sqlbuilder.AND).
    AddFilter(sqlbuilder.NewEqFilter("user_type", "trial")).
    AddFilter(sqlbuilder.NewBetweenFilter("expires_at", 
        time.Now(), 
        time.Now().AddDate(0, 0, 7)))

paidExpiringGroup := sqlbuilder.NewFilterGroup(sqlbuilder.AND).
    AddFilter(sqlbuilder.NewEqFilter("user_type", "paid")).
    AddFilter(sqlbuilder.NewBetweenFilter("expires_at", 
        time.Now(), 
        time.Now().AddDate(0, 0, 30)))

reminderGroup := sqlbuilder.NewFilterGroup(sqlbuilder.OR).
    AddGroup(trialExpiringGroup).
    AddGroup(paidExpiringGroup)
```

### 4. 使用辅助函数封装常用条件

```go
// 封装常用的筛选逻辑
func ActiveUsersGroup() *sqlbuilder.FilterGroup {
    return sqlbuilder.NewFilterGroup(sqlbuilder.AND).
        AddFilter(sqlbuilder.NewEqFilter("status", "active")).
        AddFilter(sqlbuilder.NewIsNullFilter("deleted_at"))
}

func VIPUsersGroup() *sqlbuilder.FilterGroup {
    return sqlbuilder.NewFilterGroup(sqlbuilder.OR).
        AddFilter(sqlbuilder.NewGtFilter("vip_level", 0)).
        AddFilter(sqlbuilder.NewGtFilter("total_spent", 10000))
}

// 使用
mainGroup := sqlbuilder.NewFilterGroup(sqlbuilder.AND).
    AddGroup(ActiveUsersGroup()).
    AddGroup(VIPUsersGroup())
```

## 性能建议

1. **索引优化**：确保过滤字段有索引
2. **避免过深嵌套**：通常 3 层以内最佳
3. **优先过滤数据量大的条件**：将选择性高的条件放在前面
4. **合理使用 OR**：OR 可能导致索引失效，考虑使用 UNION

## 调试技巧

### 打印生成的 SQL

```go
// 开发环境启用 Debug
db.Debug().Model(&User{}).Where(...)
```

### 查看 FilterGroup 结构

```go
// 打印 FilterGroup
fmt.Printf("%+v\n", mainGroup)
```

## 下一步

- 📖 [Repository 基础](./REPOSITORY-BASICS.MD) - 基础 CRUD 操作
- 🔍 [高级查询](./ADVANCED-QUERIES.MD) - Query 和 Filter 详解
- 🚀 [EnhancedRepository](./ENHANCED-REPOSITORY.MD) - 便利方法详解
