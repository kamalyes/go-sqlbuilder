# EnhancedRepository

EnhancedRepository 继承 BaseRepository 的所有功能，并提供了额外的便利方法。

## 创建 EnhancedRepository

```go
enhanced := repository.NewEnhancedRepository[User](handler, logger, "users")
```

## 字段查询

### FindByField - 按单字段查询多条

```go
// 查找所有状态为 active 的用户
users, err := enhanced.FindByField(ctx, "status", "active")

// 查找所有年龄为 25 的用户
users, err := enhanced.FindByField(ctx, "age", 25)
```

### FindByFieldWithOrder - 带排序的字段查询

```go
// 按创建时间倒序
users, err := enhanced.FindByFieldWithOrder(ctx, "status", "active", "created_at", "DESC")

// 按名字升序
users, err := enhanced.FindByFieldWithOrder(ctx, "city", "Beijing", "name", "ASC")
```

### FindByFieldWithLimit - 带数量限制的字段查询

```go
// 查找前 10 个活跃用户
users, err := enhanced.FindByFieldWithLimit(ctx, "status", "active", 10)

// 查找前 5 个高级会员
users, err := enhanced.FindByFieldWithLimit(ctx, "vip_level", 3, 5)
```

### FindByFieldWithCursor - 游标分页查询

游标分页适合大数据量场景，比传统分页性能更好。

```go
// 第一页（20 条）
users, cursor, err := enhanced.FindByFieldWithCursor(
    ctx,
    "status",     // 字段名
    "active",     // 字段值
    "",           // 游标（首次为空）
    20,           // 每页数量
    "id",         // 排序字段
    "ASC",        // 排序方向
)

// 下一页
nextUsers, nextCursor, err := enhanced.FindByFieldWithCursor(
    ctx,
    "status",
    "active",
    cursor,       // 使用上次返回的游标
    20,
    "id",
    "ASC",
)

// 没有更多数据时，cursor 为空字符串
if nextCursor == "" {
    fmt.Println("没有更多数据了")
}
```

**游标分页示例：加载全部数据**

```go
var allUsers []*User
cursor := ""

for {
    users, newCursor, err := enhanced.FindByFieldWithCursor(
        ctx, "status", "active", cursor, 100, "id", "ASC",
    )
    if err != nil {
        return err
    }
    
    allUsers = append(allUsers, users...)
    
    if newCursor == "" {
        break  // 没有更多数据
    }
    cursor = newCursor
}

fmt.Printf("共加载 %d 条数据\n", len(allUsers))
```

## 字段操作

### IncrementField - 字段自增

原子性地增加字段值，避免并发问题。

```go
// 用户积分 +10
err := enhanced.IncrementField(ctx, 1, "points", 10)

// 文章浏览量 +1
err := enhanced.IncrementField(ctx, articleID, "view_count", 1)

// 库存数量 +50
err := enhanced.IncrementField(ctx, productID, "stock", 50)
```

### DecrementField - 字段自减

```go
// 用户积分 -5
err := enhanced.DecrementField(ctx, 1, "points", 5)

// 库存数量 -1
err := enhanced.DecrementField(ctx, productID, "stock", 1)

// 剩余次数 -1
err := enhanced.DecrementField(ctx, userID, "remaining_attempts", 1)
```

**注意**：自减可能导致负数，建议在数据库层面设置约束或在应用层检查。

```go
// 安全的库存扣减
if err := enhanced.DecrementField(ctx, productID, "stock", quantity); err != nil {
    return errors.New("库存扣减失败")
}

// 验证库存是否充足
product, _ := repo.Get(ctx, productID)
if product.Stock < 0 {
    // 回滚或报警
}
```

### ToggleField - 布尔字段切换

```go
// 切换用户激活状态
err := enhanced.ToggleField(ctx, userID, "is_active")

// 切换文章发布状态
err := enhanced.ToggleField(ctx, articleID, "is_published")

// 切换通知开关
err := enhanced.ToggleField(ctx, settingID, "email_notification")
```

## 实战示例

### 示例 1：文章系统

```go
type Article struct {
    ID          uint
    Title       string
    AuthorID    uint
    Status      string
    ViewCount   int
    LikeCount   int
    IsPublished bool
}

enhanced := repository.NewEnhancedRepository[Article](handler, logger, "articles")

// 查找作者的所有已发布文章
articles, err := enhanced.FindByField(ctx, "author_id", authorID)

// 获取热门文章（浏览量排序，前 10 篇）
hotArticles, err := enhanced.FindByFieldWithOrder(
    ctx, "is_published", true, "view_count", "DESC",
)

// 增加浏览量
err = enhanced.IncrementField(ctx, articleID, "view_count", 1)

// 点赞
err = enhanced.IncrementField(ctx, articleID, "like_count", 1)

// 发布文章
err = enhanced.ToggleField(ctx, articleID, "is_published")
```

### 示例 2：电商库存管理

```go
type Product struct {
    ID       uint
    Name     string
    Stock    int
    Category string
    Price    float64
}

enhanced := repository.NewEnhancedRepository[Product](handler, logger, "products")

// 查找某分类的商品
products, err := enhanced.FindByField(ctx, "category", "electronics")

// 加载所有商品（游标分页）
var allProducts []*Product
cursor := ""

for {
    batch, newCursor, err := enhanced.FindByFieldWithCursor(
        ctx, "category", "electronics", cursor, 100, "id", "ASC",
    )
    if err != nil {
        break
    }
    
    allProducts = append(allProducts, batch...)
    
    if newCursor == "" {
        break
    }
    cursor = newCursor
}

// 下单：扣减库存
err = enhanced.DecrementField(ctx, productID, "stock", quantity)

// 退货：增加库存
err = enhanced.IncrementField(ctx, productID, "stock", returnQuantity)
```

### 示例 3：用户积分系统

```go
type User struct {
    ID            uint
    Username      string
    Points        int
    Level         int
    IsVIP         bool
    DailyCheckIn  bool
}

enhanced := repository.NewEnhancedRepository[User](handler, logger, "users")

// 签到奖励
func DailyCheckIn(ctx context.Context, userID uint) error {
    // 增加积分
    if err := enhanced.IncrementField(ctx, userID, "points", 10); err != nil {
        return err
    }
    
    // 标记已签到
    user, err := repo.Get(ctx, userID)
    if err != nil {
        return err
    }
    
    user.DailyCheckIn = true
    _, err = repo.Update(ctx, user)
    return err
}

// 消费积分
func ConsumePoints(ctx context.Context, userID uint, points int) error {
    // 先检查积分是否足够
    user, err := repo.Get(ctx, userID)
    if err != nil {
        return err
    }
    
    if user.Points < points {
        return errors.New("积分不足")
    }
    
    // 扣减积分
    return enhanced.DecrementField(ctx, userID, "points", points)
}

// 升级 VIP
func UpgradeToVIP(ctx context.Context, userID uint) error {
    user, err := repo.Get(ctx, userID)
    if err != nil {
        return err
    }
    
    if user.IsVIP {
        return errors.New("已是 VIP")
    }
    
    return enhanced.ToggleField(ctx, userID, "is_vip")
}
```

### 示例 4：评论系统

```go
type Comment struct {
    ID        uint
    ArticleID uint
    UserID    uint
    Content   string
    LikeCount int
    IsHidden  bool
}

enhanced := repository.NewEnhancedRepository[Comment](handler, logger, "comments")

// 获取文章评论（分页）
comments, cursor, err := enhanced.FindByFieldWithCursor(
    ctx, "article_id", articleID, "", 20, "created_at", "DESC",
)

// 点赞评论
err = enhanced.IncrementField(ctx, commentID, "like_count", 1)

// 取消点赞
err = enhanced.DecrementField(ctx, commentID, "like_count", 1)

// 隐藏评论
err = enhanced.ToggleField(ctx, commentID, "is_hidden")
```

### 示例 5：任务系统

```go
type Task struct {
    ID              uint
    UserID          uint
    Status          string
    Priority        int
    RetryCount      int
    MaxRetries      int
    IsAutoRetry     bool
}

enhanced := repository.NewEnhancedRepository[Task](handler, logger, "tasks")

// 查找用户待处理任务
tasks, err := enhanced.FindByField(ctx, "user_id", userID)

// 任务重试次数 +1
func RetryTask(ctx context.Context, taskID uint) error {
    task, err := repo.Get(ctx, taskID)
    if err != nil {
        return err
    }
    
    if task.RetryCount >= task.MaxRetries {
        return errors.New("超过最大重试次数")
    }
    
    return enhanced.IncrementField(ctx, taskID, "retry_count", 1)
}

// 切换自动重试
err = enhanced.ToggleField(ctx, taskID, "is_auto_retry")
```

## 性能对比

### 游标分页 vs 传统分页

**传统分页**（大偏移量时性能差）：

```go
// 第 1000 页，每页 100 条
// OFFSET 99900，数据库需要扫描并跳过 99900 行
query := repository.NewQuery().Limit(100).Offset(99900)
```

**游标分页**（性能稳定）：

```go
// 无论在哪一页，都是基于索引的范围查询
users, cursor, err := enhanced.FindByFieldWithCursor(
    ctx, "status", "active", cursor, 100, "id", "ASC",
)
```

### 原子操作 vs 读后写

**读后写**（有并发问题）：

```go
// ❌ 不推荐
user, _ := repo.Get(ctx, userID)
user.Points += 10
repo.Update(ctx, user)
```

**原子操作**（线程安全）：

```go
// ✅ 推荐
enhanced.IncrementField(ctx, userID, "points", 10)
```

## 最佳实践

1. **使用游标分页处理大数据集**：避免大偏移量导致的性能问题
2. **使用原子操作避免并发冲突**：IncrementField、DecrementField 是线程安全的
3. **组合使用 BaseRepository 和 EnhancedRepository**：复杂查询用 BaseRepository，便利操作用 EnhancedRepository

## 📚 相关文档

- 📖 [CRUD 操作](./CRUD-OPERATIONS.md) - BaseRepository 基础 CRUD 功能
- 🚀 [便捷查询方法](./CONVENIENCE-METHODS.md) - 简化的查询方法
- 🔍 [高级查询](./ADVANCED-QUERIES.md) - 复杂查询构建
- 🎯 [FilterGroup](./FILTERGROUP.md) - 复杂条件组合
