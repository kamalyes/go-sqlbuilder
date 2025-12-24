# CRUD 操作指南

本文档详细介绍 BaseRepository 的 CRUD（创建、读取、更新、删除）操作。

## 📖 目录

- [创建操作](#创建操作)
- [读取操作](#读取操作)
- [更新操作](#更新操作)
- [删除操作](#删除操作)

---

## 创建操作

### 1. Create - 创建记录

创建单条记录，返回包含自动填充字段（ID、CreatedAt 等）的实体。

```go
package main

import (
    "context"
    "github.com/kamalyes/go-sqlbuilder/repository"
)

func createUser(repo repository.IBaseRepository[User, uint64]) (*User, error) {
    ctx := context.Background()
    
    user := &User{
        Name:  "张三",
        Email: "test@example.com",
        Age:   25,
    }
    
    createdUser, err := repo.Create(ctx, user)
    if err != nil {
        return nil, err
    }
    
    // createdUser 包含自动填充的 ID、CreatedAt 等字段
    return createdUser, nil
}
```

### 2. CreateIfNotExists - 不存在则创建

根据唯一字段检查是否存在，不存在则创建。

```go
func createUserIfNotExists(repo repository.IBaseRepository[User, uint64]) (*User, bool, error) {
    ctx := context.Background()
    
    user := &User{
        Email: "test@example.com",
        Name:  "张三",
    }
    
    // 根据 email 字段检查是否存在
    createdUser, created, err := repo.CreateIfNotExists(ctx, user, "email")
    if err != nil {
        return nil, false, err
    }
    
    if created {
        fmt.Println("创建了新用户")
    } else {
        fmt.Println("用户已存在")
    }
    
    return createdUser, created, nil
}
```

### 3. CreateOrUpdate - 创建或更新

根据唯一字段查找，存在则更新，不存在则创建（Upsert）。

```go
func createOrUpdateUser(repo repository.IBaseRepository[User, uint64]) (*User, bool, error) {
    ctx := context.Background()
    
    user := &User{
        Email: "test@example.com",
        Name:  "张三",
        Age:   25,
    }
    
    // 根据 email 查找，存在则更新，不存在则创建
    savedUser, created, err := repo.CreateOrUpdate(ctx, user, "email")
    if err != nil {
        return nil, false, err
    }
    
    return savedUser, created, nil
}
```

### 4. CreateBatch - 批量创建

批量创建多条记录。

```go
func createMultipleUsers(repo repository.IBaseRepository[User, uint64]) error {
    ctx := context.Background()
    
    users := []*User{
        {Name: "张三", Email: "user1@example.com"},
        {Name: "李四", Email: "user2@example.com"},
        {Name: "王五", Email: "user3@example.com"},
    }
    
    err := repo.CreateBatch(ctx, users...)
    if err != nil {
        return err
    }
    
    return nil
}
```

### 5. BulkCreate - 大批量创建（分批）

适用于大批量数据导入，自动分批插入以避免单次插入数据过多。

```go
func bulkCreateUsers(repo repository.IBaseRepository[User, uint64]) error {
    ctx := context.Background()
    
    // 创建 10000 条记录
    users := make([]*User, 10000)
    for i := range users {
        users[i] = &User{
            Name:  fmt.Sprintf("User%d", i),
            Email: fmt.Sprintf("user%d@example.com", i),
        }
    }
    
    // 使用默认批次大小（1000 条/批）
    err := repo.BulkCreate(ctx, users)
    if err != nil {
        return err
    }
    
    return nil
}

func bulkCreateUsersCustomBatch(repo repository.IBaseRepository[User, uint64]) error {
    ctx := context.Background()
    
    users := make([]*User, 5000)
    for i := range users {
        users[i] = &User{Name: fmt.Sprintf("User%d", i)}
    }
    
    // 指定批次大小（每批 500 条）
    err := repo.BulkCreate(ctx, users, 500)
    return err
}
```

---

## 读取操作

### 1. Get - 根据 ID 查询

根据主键 ID 查询单条记录。

```go
func getUserByID(repo repository.IBaseRepository[User, uint64], id uint64) (*User, error) {
    ctx := context.Background()
    
    user, err := repo.Get(ctx, id)
    if err != nil {
        return nil, err
    }
    
    return user, nil
}
```

### 2. GetAll - 获取所有记录

查询表中所有记录（谨慎使用，数据量大时建议分页）。

```go
func getAllUsers(repo repository.IBaseRepository[User, uint64]) ([]*User, error) {
    ctx := context.Background()
    
    users, err := repo.GetAll(ctx)
    if err != nil {
        return nil, err
    }
    
    return users, nil
}
```

### 3. First - 获取第一条记录

获取符合条件的第一条记录（按主键正序）。

```go
func getFirstActiveUser(repo repository.IBaseRepository[User, uint64]) (*User, error) {
    ctx := context.Background()
    
    user, err := repo.First(ctx, repository.NewEqFilter("status", "active"))
    if err != nil {
        return nil, err
    }
    
    return user, nil
}
```

### 4. Last - 获取最后一条记录

获取符合条件的最后一条记录（按主键倒序）。

```go
func getLastActiveUser(repo repository.IBaseRepository[User, uint64]) (*User, error) {
    ctx := context.Background()
    
    user, err := repo.Last(ctx, repository.NewEqFilter("status", "active"))
    if err != nil {
        return nil, err
    }
    
    return user, nil
}
```

### 5. FindOne - 查找单条记录

查找符合条件的第一条记录，不存在返回 nil（不报错）。

```go
func findUserByEmail(repo repository.IBaseRepository[User, uint64], email string) (*User, error) {
    ctx := context.Background()
    
    user, err := repo.FindOne(ctx, repository.NewEqFilter("email", email))
    if err != nil {
        return nil, err
    }
    
    if user == nil {
        fmt.Println("用户不存在")
        return nil, nil
    }
    
    return user, nil
}
```

### 6. List - 列表查询

使用 Query 对象进行复杂查询。

```go
func listActiveUsers(repo repository.IBaseRepository[User, uint64]) ([]*User, error) {
    ctx := context.Background()
    
    query := repository.NewQuery().
        AddEqFilter("status", "active").
        WithSort("created_at", constants.SORT_DESC).
        Limit(100)
    
    users, err := repo.List(ctx, query)
    if err != nil {
        return nil, err
    }
    
    return users, nil
}
```

---

## 更新操作

### 1. Update - 更新记录

更新整个实体（需要主键 ID）。

```go
func updateUser(repo repository.IBaseRepository[User, uint64]) (*User, error) {
    ctx := context.Background()
    
    user := &User{
        ID:    1,
        Name:  "张三（已更新）",
        Email: "updated@example.com",
        Age:   26,
    }
    
    updatedUser, err := repo.Update(ctx, user)
    if err != nil {
        return nil, err
    }
    
    return updatedUser, nil
}
```

### 2. UpdateBatch - 批量更新

批量更新多条记录。

```go
func updateMultipleUsers(repo repository.IBaseRepository[User, uint64]) error {
    ctx := context.Background()
    
    users := []*User{
        {ID: 1, Name: "User1-Updated", Email: "user1@updated.com"},
        {ID: 2, Name: "User2-Updated", Email: "user2@updated.com"},
    }
    
    err := repo.UpdateBatch(ctx, users...)
    if err != nil {
        return err
    }
    
    return nil
}
```

### 3. UpdateFields - 更新指定字段

更新指定 ID 记录的特定字段。

```go
func updateUserFields(repo repository.IBaseRepository[User, uint64], userID uint64) error {
    ctx := context.Background()
    
    // 只更新 status 和 last_login 字段
    updates := map[string]interface{}{
        "status":     "active",
        "last_login": time.Now(),
    }
    
    err := repo.UpdateFields(ctx, userID, updates)
    if err != nil {
        return err
    }
    
    return nil
}
```

### 4. UpdateFieldsByFilters - 按条件更新字段

根据过滤条件批量更新字段。

```go
func updateInactiveUsers(repo repository.IBaseRepository[User, uint64]) error {
    ctx := context.Background()
    
    // UPDATE users SET status = 'inactive' WHERE age < 18
    updates := map[string]interface{}{
        "status": "inactive",
    }
    
    filters := []*repository.Filter{
        repository.NewLtFilter("age", 18),
    }
    
    err := repo.UpdateFieldsByFilters(ctx, updates, filters...)
    if err != nil {
        return err
    }
    
    return nil
}
```

---

## 删除操作

### 1. Delete - 删除记录

根据主键 ID 删除单条记录。

```go
func deleteUser(repo repository.IBaseRepository[User, uint64], userID uint64) error {
    ctx := context.Background()
    
    err := repo.Delete(ctx, userID)
    if err != nil {
        return err
    }
    
    return nil
}
```

### 2. DeleteBatch - 批量删除

根据多个主键 ID 批量删除记录。

```go
func deleteMultipleUsers(repo repository.IBaseRepository[User, uint64]) error {
    ctx := context.Background()
    
    // 删除 ID 为 1, 2, 3 的用户
    err := repo.DeleteBatch(ctx, 1, 2, 3)
    if err != nil {
        return err
    }
    
    return nil
}
```

### 3. DeleteByFilters - 按条件删除

根据过滤条件删除符合条件的记录。

```go
func deleteInactiveYoungUsers(repo repository.IBaseRepository[User, uint64]) error {
    ctx := context.Background()
    
    // DELETE WHERE status = 'inactive' AND age < 18
    filters := []*repository.Filter{
        repository.NewEqFilter("status", "inactive"),
        repository.NewLtFilter("age", 18),
    }
    
    err := repo.DeleteByFilters(ctx, filters...)
    if err != nil {
        return err
    }
    
    return nil
}
```

---

## 💡 最佳实践

### 1. 使用事务处理批量操作

```go
func createUsersWithTransaction(repo repository.IBaseRepository[User, uint64]) error {
    ctx := context.Background()
    
    return repo.Transaction(ctx, func(tx repository.Transaction[User]) error {
        user1 := &User{Name: "User1", Email: "user1@example.com"}
        if err := tx.Create(ctx, user1); err != nil {
            return err
        }
        
        user2 := &User{Name: "User2", Email: "user2@example.com"}
        if err := tx.Create(ctx, user2); err != nil {
            return err
        }
        
        return nil
    })
}
```

### 2. 检查记录是否存在

```go
func createUserSafely(repo repository.IBaseRepository[User, uint64], email string) error {
    ctx := context.Background()
    
    // 先检查是否存在
    exists, err := repo.Exists(ctx, repository.NewEqFilter("email", email))
    if err != nil {
        return err
    }
    
    if exists {
        return fmt.Errorf("邮箱 %s 已被使用", email)
    }
    
    // 创建用户
    user := &User{Email: email, Name: "新用户"}
    _, err = repo.Create(ctx, user)
    return err
}
```

### 3. 软删除

```go
func softDeleteUser(repo repository.IBaseRepository[User, uint64], userID uint64) error {
    ctx := context.Background()
    
    // 使用 deleted_at 字段软删除
    err := repo.SoftDelete(ctx, userID, "deleted_at", time.Now())
    return err
}
```

---

## 📚 相关文档

- [便捷查询方法](./CONVENIENCE-METHODS.md) - 简化的 CRUD 快捷方法
- [过滤条件](./FILTERS.md) - Filter 详细说明
- [Repository 基础](./REPOSITORY-BASICS.md) - Repository 基础概念
- [错误处理](./ERROR-HANDLING.md) - 错误处理最佳实践
