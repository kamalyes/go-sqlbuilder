# 便捷查询方法

BaseRepository 提供了一系列便捷方法，简化常见的查询、更新、删除操作，避免手动构建 Query 和 Filter。

## 📖 目录

- [查询方法](#查询方法)
- [分页方法](#分页方法)
- [统计方法](#统计方法)
- [更新方法](#更新方法)
- [删除方法](#删除方法)
- [保存方法](#保存方法)

---

## 查询方法

### FindWhere - 字段值对查询

简化的等值查询，支持多个字段。

```go
package main

import (
    "context"
    "github.com/kamalyes/go-sqlbuilder/repository"
)

// WHERE email = 'test@example.com' AND status = 'active'
func findActiveUserByEmail(repo repository.IBaseRepository[User, uint64]) ([]*User, error) {
    ctx := context.Background()
    
    users, err := repo.FindWhere(ctx, 
        "email", "test@example.com", 
        "status", "active")
    
    return users, err
}
```

**参数说明：**
- 参数为可变参数，按 `field, value, field, value...` 顺序传递
- 所有条件之间使用 AND 连接
- 所有条件都是等值查询（=）

### FindWhereOp - 字段操作符查询

支持自定义操作符的查询。

```go
import "github.com/kamalyes/go-sqlbuilder/constants"

// WHERE age > 18 AND status = 'active'
func findAdultActiveUsers(repo repository.IBaseRepository[User, uint64]) ([]*User, error) {
    ctx := context.Background()
    
    users, err := repo.FindWhereOp(ctx, 
        "age", constants.OP_GT, 18,
        "status", constants.OP_EQ, "active")
    
    return users, err
}
```

**参数说明：**
- 参数按 `field, operator, value, field, operator, value...` 顺序传递
- 支持所有操作符：`OP_EQ`, `OP_GT`, `OP_LT`, `OP_GTE`, `OP_LTE`, `OP_LIKE`, `OP_IN` 等

### FindOneWhere - 查询单条记录

返回符合条件的第一条记录。

```go
// WHERE email = 'test@example.com' (返回第一条)
func findUserByEmail(repo repository.IBaseRepository[User, uint64], email string) (*User, error) {
    ctx := context.Background()
    
    user, err := repo.FindOneWhere(ctx, "email", email)
    return user, err
}
```

### FindOneWhereOp - 操作符查询单条

支持操作符的单条记录查询。

```go
// WHERE age >= 18 (返回第一条)
func findOneAdultUser(repo repository.IBaseRepository[User, uint64]) (*User, error) {
    ctx := context.Background()
    
    user, err := repo.FindOneWhereOp(ctx, "age", constants.OP_GTE, 18)
    return user, err
}
```

---

## 分页方法

### Paginate - 分页查询

简化的分页查询，自动计算总数。

```go
// 第1页，每页20条，WHERE status = 'active'
func paginateActiveUsers(repo repository.IBaseRepository[User, uint64]) ([]*User, *repository.Pagination, error) {
    ctx := context.Background()
    
    users, pagination, err := repo.Paginate(ctx, 1, 20, "status", "active")
    if err != nil {
        return nil, nil, err
    }
    
    // 访问分页信息
    fmt.Printf("总记录数: %d\n", pagination.Total)
    fmt.Printf("总页数: %d\n", pagination.GetTotalPages())
    fmt.Printf("当前页: %d\n", pagination.Page)
    
    return users, pagination, nil
}
```

**返回值说明：**
- 第一个返回值：查询结果列表
- 第二个返回值：分页信息（包含 Total、Page、PageSize）
- 第三个返回值：错误信息

### PaginateOp - 分页 + 操作符查询

支持自定义操作符的分页查询。

```go
// 第2页，每页50条，WHERE age > 18
func paginateAdultUsers(repo repository.IBaseRepository[User, uint64]) ([]*User, *repository.Pagination, error) {
    ctx := context.Background()
    
    users, pagination, err := repo.PaginateOp(ctx, 2, 50, 
        "age", constants.OP_GT, 18)
    
    return users, pagination, err
}
```

---

## 统计方法

### CountWhere - 统计数量

快速统计符合条件的记录数。

```go
// COUNT(*) WHERE status = 'active'
func countActiveUsers(repo repository.IBaseRepository[User, uint64]) (int64, error) {
    ctx := context.Background()
    
    count, err := repo.CountWhere(ctx, "status", "active")
    return count, err
}
```

### CountWhereOp - 操作符统计

支持自定义操作符的统计。

```go
// COUNT(*) WHERE age >= 18
func countAdultUsers(repo repository.IBaseRepository[User, uint64]) (int64, error) {
    ctx := context.Background()
    
    count, err := repo.CountWhereOp(ctx, "age", constants.OP_GTE, 18)
    return count, err
}
```

### ExistsWhere - 检查记录是否存在

快速检查符合条件的记录是否存在。

```go
// 检查邮箱是否已存在
func isEmailExists(repo repository.IBaseRepository[User, uint64], email string) (bool, error) {
    ctx := context.Background()
    
    exists, err := repo.ExistsWhere(ctx, "email", email)
    return exists, err
}
```

---

## 更新方法

### UpdateWhere - 条件更新

根据条件批量更新字段。

```go
// UPDATE users SET status = 'inactive' WHERE last_login < '2023-01-01'
func deactivateInactiveUsers(repo repository.IBaseRepository[User, uint64], lastLoginThreshold time.Time) error {
    ctx := context.Background()
    
    updates := map[string]interface{}{
        "status": "inactive",
    }
    
    err := repo.UpdateWhere(ctx, updates, "last_login", lastLoginThreshold)
    return err
}
```

**参数说明：**
- 第一个参数：Context
- 第二个参数：更新字段 map
- 后续参数：`field, value` 对，用于构建 WHERE 条件

### UpdateWhereOp - 操作符更新

支持自定义操作符的条件更新。

```go
// UPDATE users SET vip_level = 0 WHERE total_spent < 1000
func downgradeVIPLevel(repo repository.IBaseRepository[User, uint64]) error {
    ctx := context.Background()
    
    updates := map[string]interface{}{
        "vip_level": 0,
    }
    
    err := repo.UpdateWhereOp(ctx, updates, 
        "total_spent", constants.OP_LT, 1000)
    
    return err
}
```

---

## 删除方法

### DeleteWhere - 条件删除

根据条件批量删除记录。

```go
// DELETE WHERE status = 'inactive' AND last_login < '2023-01-01'
func deleteOldInactiveUsers(repo repository.IBaseRepository[User, uint64], threshold time.Time) error {
    ctx := context.Background()
    
    err := repo.DeleteWhere(ctx, 
        "status", "inactive",
        "last_login", threshold)
    
    return err
}
```

### DeleteWhereOp - 操作符删除

支持自定义操作符的条件删除。

```go
// DELETE WHERE age < 18
func deleteUnderageUsers(repo repository.IBaseRepository[User, uint64]) error {
    ctx := context.Background()
    
    err := repo.DeleteWhereOp(ctx, "age", constants.OP_LT, 18)
    return err
}
```

---

## 保存方法

### Save - 保存（创建或更新）

根据 ID 判断是创建还是更新。

```go
func saveUser(repo repository.IBaseRepository[User, uint64], user *User) error {
    ctx := context.Background()
    
    // ID 为 0 则创建，否则更新
    savedUser, err := repo.Save(ctx, user)
    if err != nil {
        return err
    }
    
    fmt.Printf("保存成功: %+v\n", savedUser)
    return nil
}
```

**使用场景：**

```go
// 场景 1: 创建新用户
newUser := &User{
    Name:  "张三",
    Email: "test@example.com",
}
repo.Save(ctx, newUser) // 会执行 INSERT

// 场景 2: 更新现有用户
existingUser := &User{
    ID:    1,
    Name:  "张三（已更新）",
    Email: "updated@example.com",
}
repo.Save(ctx, existingUser) // 会执行 UPDATE
```

---

## 💡 使用技巧

### 1. 组合使用便捷方法

```go
func processUsers(repo repository.IBaseRepository[User, uint64]) error {
    ctx := context.Background()
    
    // 1. 检查是否存在
    exists, _ := repo.ExistsWhere(ctx, "email", "test@example.com")
    if !exists {
        return fmt.Errorf("用户不存在")
    }
    
    // 2. 查询用户
    user, _ := repo.FindOneWhere(ctx, "email", "test@example.com")
    
    // 3. 更新状态
    updates := map[string]interface{}{"status": "active"}
    repo.UpdateWhere(ctx, updates, "email", "test@example.com")
    
    return nil
}
```

### 2. 分页查询最佳实践

```go
func paginateWithSort(repo repository.IBaseRepository[User, uint64], page, pageSize int32) {
    ctx := context.Background()
    
    // 便捷方法只支持简单查询，复杂查询请使用 Query
    // 如果需要排序，建议使用 Query + ListWithPagination
    
    query := repository.NewQuery().
        AddEqFilter("status", "active").
        WithSort("created_at", constants.SORT_DESC)
    
    pagination := &repository.Pagination{Page: page, PageSize: pageSize}
    users, paginationResult, err := repo.ListWithPagination(ctx, query, pagination)
    
    // 处理结果...
}
```

### 3. 批量操作优化

```go
func bulkUpdateActiveUsers(repo repository.IBaseRepository[User, uint64]) error {
    ctx := context.Background()
    
    // 先统计数量
    count, _ := repo.CountWhere(ctx, "status", "active")
    fmt.Printf("将要更新 %d 条记录\n", count)
    
    // 执行批量更新
    updates := map[string]interface{}{
        "last_checked": time.Now(),
    }
    err := repo.UpdateWhere(ctx, updates, "status", "active")
    
    return err
}
```

---

## ⚠️ 注意事项

### 1. 参数顺序

便捷方法的参数必须按特定顺序传递：

- `FindWhere`: `field, value, field, value...`
- `FindWhereOp`: `field, operator, value, field, operator, value...`

错误示例：
```go
// ❌ 错误：参数顺序不对
repo.FindWhere(ctx, "active", "status")

// ✅ 正确
repo.FindWhere(ctx, "status", "active")
```

### 2. 分页返回值顺序

`Paginate` 方法返回值顺序为：`users, pagination, err`

```go
// ✅ 正确
users, pagination, err := repo.Paginate(ctx, 1, 20, "status", "active")

// ❌ 错误
pagination, users, err := repo.Paginate(ctx, 1, 20, "status", "active")
```

### 3. 复杂查询使用 Query

便捷方法适用于简单场景，复杂查询（如 OR 条件、子查询、JOIN）请使用 Query 对象：

```go
// 复杂查询：需要 OR 条件
query := repository.NewQuery().
    SetFilterGroup(
        repository.NewFilterGroup(constants.AND_OR).
            AddFilter(repository.NewEqFilter("status", "active")).
            AddFilter(repository.NewEqFilter("status", "pending")))

users, err := repo.List(ctx, query)
```

---

## 📚 相关文档

- [CRUD 操作](./CRUD-OPERATIONS.md) - 基础 CRUD 方法
- [工具方法](./UTILITY-METHODS.md) - Count、Exists、Pluck 等
- [过滤条件](./FILTERS.md) - Filter 详细说明
- [排序和分页](./SORTING-AND-PAGINATION.md) - 分页详解
