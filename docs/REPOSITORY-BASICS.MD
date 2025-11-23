# Repository 基础操作

BaseRepository 提供了完整的 CRUD 操作方法。

## 创建操作

### Create - 创建单个记录

```go
user := &User{
    Name:  "张三",
    Email: "zhangsan@example.com",
    Age:   25,
}
created, err := repo.Create(ctx, user)
```

### CreateBatch - 批量创建

```go
users := []*User{
    {Name: "用户1", Email: "user1@test.com"},
    {Name: "用户2", Email: "user2@test.com"},
    {Name: "用户3", Email: "user3@test.com"},
}
err := repo.CreateBatch(ctx, users...)
```

### BulkCreate - 高性能批量创建

```go
// 默认批次大小为 100
err := repo.BulkCreate(ctx, users)

// 自定义批次大小
err := repo.BulkCreate(ctx, users, 50)
```

### CreateIfNotExists - 不存在则创建

```go
user := &User{Email: "test@example.com", Name: "Test"}
created, isNew, err := repo.CreateIfNotExists(ctx, user, "email")
if isNew {
    log.Println("创建了新用户")
} else {
    log.Println("用户已存在")
}
```

### CreateOrUpdate - 创建或更新

```go
user := &User{Email: "test@example.com", Name: "Updated Name"}
result, created, err := repo.CreateOrUpdate(ctx, user, "email")
```

## 查询操作

### Get - 按 ID 查询

```go
user, err := repo.Get(ctx, 1)
```

### GetWithPreloads - 带预加载的查询

```go
user, err := repo.GetWithPreloads(ctx, 1, "Profile", "Orders")
```

### GetByFields - 按多个字段查询

```go
user, err := repo.GetByFields(ctx, map[string]interface{}{
    "email":  "test@example.com",
    "status": "active",
})
```

### GetByFilter - 按单个过滤条件查询

```go
filter := sqlbuilder.NewEqFilter("email", "test@example.com")
user, err := repo.GetByFilter(ctx, filter)
```

### GetByFilters - 按多个过滤条件查询

```go
filters := []*sqlbuilder.Filter{
    sqlbuilder.NewEqFilter("status", "active"),
    sqlbuilder.NewGtFilter("age", 18),
}
user, err := repo.GetByFilters(ctx, filters...)
```

### List - 列表查询

```go
query := sqlbuilder.NewQuery().
    AddFilter(sqlbuilder.NewEqFilter("status", "active")).
    AddOrder("created_at", "DESC").
    Limit(10)
    
users, err := repo.List(ctx, query)
```

### ListWithPagination - 分页查询

```go
query := sqlbuilder.NewQuery().
    AddFilter(sqlbuilder.NewEqFilter("status", "active"))
    
paging := &sqlbuilder.Paging{
    Page:     1,
    PageSize: 20,
}

users, paging, err := repo.ListWithPagination(ctx, query, paging)
log.Printf("总记录数: %d, 总页数: %d", paging.Total, paging.GetTotalPages())
```

### GetAll - 获取所有记录

```go
users, err := repo.GetAll(ctx)
```

### First / Last - 获取第一条/最后一条

```go
// 第一条
firstUser, err := repo.First(ctx)

// 带条件的第一条
activeFirst, err := repo.First(ctx, sqlbuilder.NewEqFilter("status", "active"))

// 最后一条
lastUser, err := repo.Last(ctx)
```

### FindOne - 查找单条（不存在返回 nil）

```go
user, err := repo.FindOne(ctx, sqlbuilder.NewEqFilter("email", "test@example.com"))
if user == nil {
    log.Println("未找到用户")
}
```

## 更新操作

### Update - 更新整个记录

```go
user.Name = "新名字"
user.Age = 30
updated, err := repo.Update(ctx, user)
```

### UpdateBatch - 批量更新

```go
users[0].Name = "更新1"
users[1].Name = "更新2"
err := repo.UpdateBatch(ctx, users...)
```

### UpdateFields - 更新指定字段

```go
err := repo.UpdateFields(ctx, 1, map[string]interface{}{
    "name":       "新名字",
    "updated_at": time.Now(),
})
```

### UpdateFieldsByFilters - 按条件更新字段

```go
filters := []*sqlbuilder.Filter{
    sqlbuilder.NewEqFilter("status", "inactive"),
}
fields := map[string]interface{}{
    "status": "active",
}
err := repo.UpdateFieldsByFilters(ctx, fields, filters...)
```

### UpdateByFilters - 按条件更新

```go
updateData := &User{Status: "archived"}
filters := []*sqlbuilder.Filter{
    sqlbuilder.NewLtFilter("last_login", time.Now().AddDate(0, -6, 0)),
}
err := repo.UpdateByFilters(ctx, updateData, filters...)
```

## 删除操作

### Delete - 删除单个记录

```go
err := repo.Delete(ctx, 1)
```

### DeleteBatch - 批量删除

```go
ids := []interface{}{1, 2, 3, 4, 5}
err := repo.DeleteBatch(ctx, ids...)
```

### DeleteByFilters - 按条件删除

```go
filters := []*sqlbuilder.Filter{
    sqlbuilder.NewEqFilter("status", "deleted"),
    sqlbuilder.NewLtFilter("created_at", time.Now().AddDate(-1, 0, 0)),
}
err := repo.DeleteByFilters(ctx, filters...)
```

### 软删除

```go
// 使用 deleted_at 字段
err := repo.SoftDelete(ctx, 1, "deleted_at", time.Now())

// 批量软删除
ids := []interface{}{1, 2, 3}
err := repo.SoftDeleteBatch(ctx, ids, "deleted_at", time.Now())

// 按条件软删除
filters := []*sqlbuilder.Filter{
    sqlbuilder.NewEqFilter("status", "inactive"),
}
err := repo.SoftDeleteByFilters(ctx, "deleted_at", time.Now(), filters...)
```

### 恢复软删除

```go
// 恢复单个
err := repo.Restore(ctx, 1, "deleted_at", nil)

// 批量恢复
ids := []interface{}{1, 2, 3}
err := repo.RestoreBatch(ctx, ids, "deleted_at", nil)
```

## 统计操作

### Count - 计数

```go
// 总数
total, err := repo.Count(ctx)

// 按条件计数
count, err := repo.Count(ctx, sqlbuilder.NewEqFilter("status", "active"))
```

### Exists - 检查是否存在

```go
exists, err := repo.Exists(ctx, sqlbuilder.NewEqFilter("email", "test@example.com"))
```

### CountByField - 按字段分组计数

```go
countMap, err := repo.CountByField(ctx, "status")
// 结果: map[interface{}]int64{"active": 100, "inactive": 50}
```

### Pluck - 提取字段值

```go
emails, err := repo.Pluck(ctx, "email", sqlbuilder.NewEqFilter("status", "active"))
```

### Distinct - 获取去重值

```go
statuses, err := repo.Distinct(ctx, "status")
```

## 事务操作

```go
err := repo.Transaction(ctx, func(tx sqlbuilder.Transaction[User]) error {
    // 创建用户
    user := &User{Name: "Transaction User"}
    if err := tx.Create(ctx, user); err != nil {
        return err
    }
    
    // 批量创建
    users := []*User{
        {Name: "User1"},
        {Name: "User2"},
    }
    if err := tx.CreateBatch(ctx, users...); err != nil {
        return err
    }
    
    // 更新
    user.Status = "active"
    if err := tx.Update(ctx, user); err != nil {
        return err
    }
    
    return nil
})
```

## Repository 配置选项

```go
repo := sqlbuilder.NewBaseRepository[User](
    handler,
    logger,
    "users",
    sqlbuilder.WithBatchSize[User](200),              // 批处理大小
    sqlbuilder.WithTimeout[User](60),                 // 查询超时(秒)
    sqlbuilder.WithReadOnly[User](),                  // 只读模式
    sqlbuilder.WithDefaultPreloads[User]("Profile"),  // 默认预加载
    sqlbuilder.WithDefaultOrder[User]("id DESC"),     // 默认排序
)
```

## 下一步

- 🔍 [高级查询](./ADVANCED-QUERIES.MD) - Query 和 Filter 的高级用法
- 🎯 [FilterGroup](./FILTERGROUP.MD) - 复杂 WHERE 条件构建
- 🚀 [EnhancedRepository](./ENHANCED-REPOSITORY.MD) - 便利方法详解
