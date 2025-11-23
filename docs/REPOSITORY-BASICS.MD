# Repository 基础操作

BaseRepository 提供了完整的 CRUD 操作方法�?

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
// 默认批次大小�?100
err := repo.BulkCreate(ctx, users)

// 自定义批次大�?
err := repo.BulkCreate(ctx, users, 50)
```

### CreateIfNotExists - 不存在则创建

```go
user := &User{Email: "test@example.com", Name: "Test"}
created, isNew, err := repo.CreateIfNotExists(ctx, user, "email")
if isNew {
    log.Println("创建了新用户")
} else {
    log.Println("用户已存�?)
}
```

### CreateOrUpdate - 创建或更�?

```go
user := &User{Email: "test@example.com", Name: "Updated Name"}
result, created, err := repo.CreateOrUpdate(ctx, user, "email")
```

## 查询操作

### Get - �?ID 查询

```go
user, err := repo.Get(ctx, 1)
```

### GetWithPreloads - 带预加载的查�?

```go
user, err := repo.GetWithPreloads(ctx, 1, "Profile", "Orders")
```

### GetByFields - 按多个字段查�?

```go
user, err := repo.GetByFields(ctx, map[string]interface{}{
    "email":  "test@example.com",
    "status": "active",
})
```

### GetByFilter - 按单个过滤条件查�?

```go
filter := repository.NewEqFilter("email", "test@example.com")
user, err := repo.GetByFilter(ctx, filter)
```

### GetByFilters - 按多个过滤条件查�?

```go
filters := []*repository.Filter{
    repository.NewEqFilter("status", "active"),
    repository.NewGtFilter("age", 18),
}
user, err := repo.GetByFilters(ctx, filters...)
```

### List - 列表查询

```go
query := repository.NewQuery().
    AddFilter(repository.NewEqFilter("status", "active")).
    AddOrder("created_at", "DESC").
    Limit(10)
    
users, err := repo.List(ctx, query)
```

### ListWithPagination - 分页查询

```go
query := repository.NewQuery().
    AddFilter(repository.NewEqFilter("status", "active"))
    
paging := &repository.Paging{
    Page:     1,
    PageSize: 20,
}

users, paging, err := repo.ListWithPagination(ctx, query, paging)
log.Printf("总记录数: %d, 总页�? %d", paging.Total, paging.GetTotalPages())
```

### GetAll - 获取所有记�?

```go
users, err := repo.GetAll(ctx)
```

### First / Last - 获取第一�?最后一�?

```go
// 第一�?
firstUser, err := repo.First(ctx)

// 带条件的第一�?
activeFirst, err := repo.First(ctx, repository.NewEqFilter("status", "active"))

// 最后一�?
lastUser, err := repo.Last(ctx)
```

### FindOne - 查找单条（不存在返回 nil�?

```go
user, err := repo.FindOne(ctx, repository.NewEqFilter("email", "test@example.com"))
if user == nil {
    log.Println("未找到用�?)
}
```

## 更新操作

### Update - 更新整个记录

```go
user.Name = "新名�?
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
    "name":       "新名�?,
    "updated_at": time.Now(),
})
```

### UpdateFieldsByFilters - 按条件更新字�?

```go
filters := []*repository.Filter{
    repository.NewEqFilter("status", "inactive"),
}
fields := map[string]interface{}{
    "status": "active",
}
err := repo.UpdateFieldsByFilters(ctx, fields, filters...)
```

### UpdateByFilters - 按条件更�?

```go
updateData := &User{Status: "archived"}
filters := []*repository.Filter{
    repository.NewLtFilter("last_login", time.Now().AddDate(0, -6, 0)),
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

### DeleteByFilters - 按条件删�?

```go
filters := []*repository.Filter{
    repository.NewEqFilter("status", "deleted"),
    repository.NewLtFilter("created_at", time.Now().AddDate(-1, 0, 0)),
}
err := repo.DeleteByFilters(ctx, filters...)
```

### 软删�?

```go
// 使用 deleted_at 字段
err := repo.SoftDelete(ctx, 1, "deleted_at", time.Now())

// 批量软删�?
ids := []interface{}{1, 2, 3}
err := repo.SoftDeleteBatch(ctx, ids, "deleted_at", time.Now())

// 按条件软删除
filters := []*repository.Filter{
    repository.NewEqFilter("status", "inactive"),
}
err := repo.SoftDeleteByFilters(ctx, "deleted_at", time.Now(), filters...)
```

### 恢复软删�?

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

// 按条件计�?
count, err := repo.Count(ctx, repository.NewEqFilter("status", "active"))
```

### Exists - 检查是否存�?

```go
exists, err := repo.Exists(ctx, repository.NewEqFilter("email", "test@example.com"))
```

### CountByField - 按字段分组计�?

```go
countMap, err := repo.CountByField(ctx, "status")
// 结果: map[interface{}]int64{"active": 100, "inactive": 50}
```

### Pluck - 提取字段�?

```go
emails, err := repo.Pluck(ctx, "email", repository.NewEqFilter("status", "active"))
```

### Distinct - 获取去重�?

```go
statuses, err := repo.Distinct(ctx, "status")
```

## 事务操作

```go
err := repo.Transaction(ctx, func(tx repository.Transaction[User]) error {
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
repo := repository.NewBaseRepository[User](
    handler,
    logger,
    "users",
    repository.WithBatchSize[User](200),              // 批处理大�?
    repository.WithTimeout[User](60),                 // 查询超时(�?
    repository.WithReadOnly[User](),                  // 只读模式
    repository.WithDefaultPreloads[User]("Profile"),  // 默认预加�?
    repository.WithDefaultOrder[User]("id DESC"),     // 默认排序
)
```

## 下一�?

- 🔍 [高级查询](./ADVANCED-QUERIES.MD) - Query �?Filter 的高级用�?
- 🎯 [FilterGroup](./FILTERGROUP.MD) - 复杂 WHERE 条件构建
- 🚀 [EnhancedRepository](./ENHANCED-REPOSITORY.MD) - 便利方法详解
