# FilterGroup 完整指南

FilterGroup 是构建复WHERE 条件的强大工具，支持 AND/OR 逻辑组合和无限层级嵌套

## 基础概念

### AND 逻辑

所有条件必须同时满足（交集）

```go
import (
    "github.com/kamalyes/go-sqlbuilder/constants"
    "github.com/kamalyes/go-sqlbuilder/repository"
)

// WHERE status = 'active' AND age > 18
group := repository.NewFilterGroup(constants.AND_AND).
    AddFilter(repository.NewEqFilter("status", "active")).
    AddFilter(repository.NewGtFilter("age", 18))
```

### OR 逻辑

任一条件满足即可（并集）

```go
// WHERE status = 'active' OR status = 'verified'
group := repository.NewFilterGroup(constants.AND_OR).
    AddFilter(repository.NewEqFilter("status", "active")).
    AddFilter(repository.NewEqFilter("status", "verified"))
```

## 完整 API

### NewFilterGroup

创建过滤器组

```go
// AND 
andGroup := repository.NewFilterGroup(constants.AND_AND)

// OR 
orGroup := repository.NewFilterGroup(constants.AND_OR)
```

### AddFilter

添加单个过滤条件

```go
group.AddFilter(repository.NewEqFilter("status", "active"))
```

### AddFilters

批量添加多个过滤条件

```go
filters := []*repository.Filter{
    repository.NewEqFilter("status", "active"),
    repository.NewGtFilter("age", 18),
    repository.NewLikeFilter("name", "),
}
group.AddFilters(filters...)
```

### AddGroup

添加嵌套的子组

```go
subGroup := repository.NewFilterGroup(constants.AND_OR).
    AddFilter(repository.NewEqFilter("role", "admin")).
    AddFilter(repository.NewEqFilter("role", "moderator"))
    
mainGroup := repository.NewFilterGroup(constants.AND_AND).
    AddGroup(subGroup).
    AddFilter(repository.NewEqFilter("status", "active"))
```

### Apply

FilterGroup 应用GORM 查询

```go
db := handler.GetDB()
query := group.Apply(db)
```

## 实战示例

### 示例 1：电商商品筛

**需*：查找价格在 100-1000 之间的电子产品或图书，且评分大于 4.0 或销量大1000

```go
// (category IN ('electronics', 'books') AND price BETWEEN 100 AND 1000)
// AND (rating > 4.0 OR sales > 1000)

categoryFilter := repository.NewInFilter("category", []interface{}{"electronics", "books"})
priceFilter := repository.NewBetweenFilter("price", 100, 1000)

qualityGroup := repository.NewFilterGroup(constants.AND_OR).
    AddFilter(repository.NewGtFilter("rating", 4.0)).
    AddFilter(repository.NewGtFilter("sales", 1000))

mainGroup := repository.NewFilterGroup(constants.AND_AND).
    AddFilter(categoryFilter).
    AddFilter(priceFilter).
    AddGroup(qualityGroup)

query := repository.NewQuery().SetFilterGroup(mainGroup)
products, err := repo.List(ctx, query)
```

### 示例 2：用户权限筛

**需*：查找活跃的管理员或版主，且年龄大于 18 岁或已验证邮箱

```go
// ((role = 'admin' OR role = 'moderator') AND status = 'active')
// AND (age > 18 OR email_verified = true)

roleGroup := repository.NewFilterGroup(constants.AND_OR).
    AddFilter(repository.NewEqFilter("role", "admin")).
    AddFilter(repository.NewEqFilter("role", "moderator"))

userGroup := repository.NewFilterGroup(constants.AND_AND).
    AddGroup(roleGroup).
    AddFilter(repository.NewEqFilter("status", "active"))

verificationGroup := repository.NewFilterGroup(constants.AND_OR).
    AddFilter(repository.NewGtFilter("age", 18)).
    AddFilter(repository.NewEqFilter("email_verified", true))

finalGroup := repository.NewFilterGroup(constants.AND_AND).
    AddGroup(userGroup).
    AddGroup(verificationGroup)

query := repository.NewQuery().SetFilterGroup(finalGroup)
users, err := repo.List(ctx, query)
```

### 示例 3：订单状态查

**需*：查找待支付超过 30 分钟的订单，或已支付但未发货超过 24 小时的订单

```go
// (status = 'pending' AND created_at < NOW() - 30分钟)
// OR (status = 'paid' AND paid_at < NOW() - 24小时 AND shipped_at IS NULL)

thirtyMinAgo := time.Now().Add(-30 * time.Minute)
oneDayAgo := time.Now().Add(-24 * time.Hour)

pendingGroup := repository.NewFilterGroup(constants.AND_AND).
    AddFilter(repository.NewEqFilter("status", "pending")).
    AddFilter(repository.NewLtFilter("created_at", thirtyMinAgo))

paidGroup := repository.NewFilterGroup(constants.AND_AND).
    AddFilter(repository.NewEqFilter("status", "paid")).
    AddFilter(repository.NewLtFilter("paid_at", oneDayAgo)).
    AddFilter(repository.NewIsNullFilter("shipped_at"))

mainGroup := repository.NewFilterGroup(constants.AND_OR).
    AddGroup(pendingGroup).
    AddGroup(paidGroup)

query := repository.NewQuery().SetFilterGroup(mainGroup)
orders, err := repo.List(ctx, query)
```

### 示例 4：内容审核查

**需*：查找需要人工审核的内容（包含敏感词或被举报次数大于 3，且未被审核）

```go
// (has_sensitive_words = true OR report_count > 3)
// AND status = 'pending_review'
// AND reviewed_at IS NULL

flagGroup := repository.NewFilterGroup(constants.AND_OR).
    AddFilter(repository.NewEqFilter("has_sensitive_words", true)).
    AddFilter(repository.NewGtFilter("report_count", 3))

mainGroup := repository.NewFilterGroup(constants.AND_AND).
    AddGroup(flagGroup).
    AddFilter(repository.NewEqFilter("status", "pending_review")).
    AddFilter(repository.NewIsNullFilter("reviewed_at"))

query := repository.NewQuery().SetFilterGroup(mainGroup)
contents, err := repo.List(ctx, query)
```

### 示例 5：多条件搜索

**需*：搜索名字或邮箱包含关键词，且状态为活跃或试用中的用户

```go
// (name LIKE '%keyword%' OR email LIKE '%keyword%')
// AND (status = 'active' OR status = 'trial')

keyword := "张三"

searchGroup := repository.NewFilterGroup(constants.AND_OR).
    AddFilter(repository.NewLikeFilter("name", keyword)).
    AddFilter(repository.NewLikeFilter("email", keyword))

statusGroup := repository.NewFilterGroup(constants.AND_OR).
    AddFilter(repository.NewEqFilter("status", "active")).
    AddFilter(repository.NewEqFilter("status", "trial"))

mainGroup := repository.NewFilterGroup(constants.AND_AND).
    AddGroup(searchGroup).
    AddGroup(statusGroup)

query := repository.NewQuery().SetFilterGroup(mainGroup)
users, err := repo.List(ctx, query)
```

### 示例 6：时间范围组

**需*：查找本周创建的订单，或上月更新过的高优先级订单

```go
// (created_at >= THIS_WEEK_START AND created_at < THIS_WEEK_END)
// OR (updated_at >= LAST_MONTH_START AND updated_at < LAST_MONTH_END AND priority = 'high')

thisWeekGroup := repository.NewFilterGroup(constants.AND_AND).
    AddFilter(repository.NewThisWeekFilter("created_at"))

lastMonthGroup := repository.NewFilterGroup(constants.AND_AND).
    AddFilter(repository.NewLastMonthFilter("updated_at")).
    AddFilter(repository.NewEqFilter("priority", "high"))

mainGroup := repository.NewFilterGroup(constants.AND_OR).
    AddGroup(thisWeekGroup).
    AddGroup(lastMonthGroup)

query := repository.NewQuery().SetFilterGroup(mainGroup)
orders, err := repo.List(ctx, query)
```

### 示例 7：多级嵌

**需*：复杂的会员筛选（VIP 或高消费用户，且活跃或最近登录）

```go
// ((vip_level > 0 OR total_spent > 10000) AND status = 'active')
// AND (last_login > 7天前 OR login_count_this_month > 10)

memberGroup := repository.NewFilterGroup(constants.AND_OR).
    AddFilter(repository.NewGtFilter("vip_level", 0)).
    AddFilter(repository.NewGtFilter("total_spent", 10000))

memberStatusGroup := repository.NewFilterGroup(constants.AND_AND).
    AddGroup(memberGroup).
    AddFilter(repository.NewEqFilter("status", "active"))

activityGroup := repository.NewFilterGroup(constants.AND_OR).
    AddFilter(repository.NewGtFilter("last_login", time.Now().AddDate(0, 0, -7))).
    AddFilter(repository.NewGtFilter("login_count_this_month", 10))

finalGroup := repository.NewFilterGroup(constants.AND_AND).
    AddGroup(memberStatusGroup).
    AddGroup(activityGroup)

query := repository.NewQuery().SetFilterGroup(finalGroup)
members, err := repo.List(ctx, query)
```

## 最佳实

### 1. 命名清晰的变

```go
// 不好
g1 := repository.NewFilterGroup(constants.AND_OR)
g2 := repository.NewFilterGroup(constants.AND_AND)

// 
roleGroup := repository.NewFilterGroup(constants.AND_OR)
statusGroup := repository.NewFilterGroup(constants.AND_AND)
```

### 2. 逐步构建复杂条件

```go
// 先构建子
categoryGroup := repository.NewFilterGroup(constants.AND_OR).
    AddFilter(repository.NewEqFilter("category", "A")).
    AddFilter(repository.NewEqFilter("category", "B"))

priceGroup := repository.NewFilterGroup(constants.AND_AND).
    AddFilter(repository.NewGteFilter("price", 100)).
    AddFilter(repository.NewLteFilter("price", 1000))

// 再组
finalGroup := repository.NewFilterGroup(constants.AND_AND).
    AddGroup(categoryGroup).
    AddGroup(priceGroup)
```

### 3. 添加注释说明逻辑

```go
// 查找需要提醒的用户
// - 试用期即将到期（7天内
// - 或者付费会员即将到期（30天内
trialExpiringGroup := repository.NewFilterGroup(constants.AND_AND).
    AddFilter(repository.NewEqFilter("user_type", "trial")).
    AddFilter(repository.NewBetweenFilter("expires_at", 
        time.Now(), 
        time.Now().AddDate(0, 0, 7)))

paidExpiringGroup := repository.NewFilterGroup(constants.AND_AND).
    AddFilter(repository.NewEqFilter("user_type", "paid")).
    AddFilter(repository.NewBetweenFilter("expires_at", 
        time.Now(), 
        time.Now().AddDate(0, 0, 30)))

reminderGroup := repository.NewFilterGroup(constants.AND_OR).
    AddGroup(trialExpiringGroup).
    AddGroup(paidExpiringGroup)
```

### 4. 使用辅助函数封装常用条件

```go
// 封装常用的筛选逻辑
func ActiveUsersGroup() *repository.FilterGroup {
    return repository.NewFilterGroup(constants.AND_AND).
        AddFilter(repository.NewEqFilter("status", "active")).
        AddFilter(repository.NewIsNullFilter("deleted_at"))
}

func VIPUsersGroup() *repository.FilterGroup {
    return repository.NewFilterGroup(constants.AND_OR).
        AddFilter(repository.NewGtFilter("vip_level", 0)).
        AddFilter(repository.NewGtFilter("total_spent", 10000))
}

// 使用
mainGroup := repository.NewFilterGroup(constants.AND_AND).
    AddGroup(ActiveUsersGroup()).
    AddGroup(VIPUsersGroup())
```

## 性能建议

1. **索引优化**：确保过滤字段有索引
2. **避免过深嵌套**：通常 3 层以内最
3. **优先过滤数据量大的条*：将选择性高的条件放在前
4. **合理使用 OR**：OR 可能导致索引失效，考虑使用 UNION

## 调试技

### 打印生成SQL

```go
// 开发环境启Debug
db.Debug().Model(&User{}).Where(...)
```

### 查看 FilterGroup 结构

```go
// 打印 FilterGroup
fmt.Printf("%+v\n", mainGroup)
```

## 下一

- 📖 [Repository 基础](./REPOSITORY-BASICS.MD) - 基础 CRUD 操作
- 🔍 [高级查询](./ADVANCED-QUERIES.MD) - Query Filter 详解
- 🚀 [EnhancedRepository](./ENHANCED-REPOSITORY.MD) - 便利方法详解

