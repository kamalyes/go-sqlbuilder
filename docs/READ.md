# 查询操作 (Read)

## 概述
查询操作用于从数据库检索记录，支持主键查询、条件查询、列表查询和分页查询

## 主键查询

### 根据 ID 查询
```go
// 查询 ID 为 1 的用户
user, err := repo.Get(ctx, 1)
if err != nil {
    return err
}
// 如果记录不存在，返回 gorm.ErrRecordNotFound
```

### 查询所有记录
```go
// 查询所有用户
users, err := repo.GetAll(ctx)
if err != nil {
    return err
}
```

## 条件查询

### 根据 Filter 查询单条
```go
// 根据邮箱查询单个用户
filter := repository.NewEqFilter("email", "zhangsan@example.com")
user, err := repo.GetByFilter(ctx, filter)
```

### 根据多个 Filters 查询单条
```go
// 根据多个条件查询
user, err := repo.GetByFilters(ctx,
    repository.NewEqFilter("status", "active"),
    repository.NewEqFilter("email", "test@example.com"),
)
```

## 列表查询

### 基础列表查询
```go
// 查询所有记录
users, err := repo.List(ctx, repository.NewQuery())

// 带条件的列表查询
query := repository.NewQuery().
    AddFilter(repository.NewEqFilter("status", "active")).
    AddOrder("created_at", "DESC")
users, err := repo.List(ctx, query)
```

### 带分页的列表查询
`page` 参数为**可选参数**，支持以下调用方式：

```go
// 方式1：不传 page，使用 query 中的分页设置（推荐）
query := repository.NewQuery().
    WithPaging(1, 20).
    AddFilter(repository.NewEqFilter("status", "active"))
users, pageInfo, err := repo.ListWithPagination(ctx, query)

// 方式2：不传 page 且 query 未设置分页，自动使用默认值（Page=1, PageSize=20）
query := repository.NewQuery().
    AddFilter(repository.NewEqFilter("status", "active"))
users, pageInfo, err := repo.ListWithPagination(ctx, query)

// 方式3：显式传入 page，覆盖 query 中的分页设置
query := repository.NewQuery().
    AddFilter(repository.NewEqFilter("status", "active"))
pagination := &repository.Pagination{Page: 1, PageSize: 20}
users, pageInfo, err := repo.ListWithPagination(ctx, query, pagination)

// pageInfo 包含：
// - Total: 总记录数
// - PageCount: 总页数
// - HasNext: 是否有下一页
// - HasPrev: 是否有上一页
```

**优先级规则**：显式传入的 `page` > `query.Pagination`（通过 `WithPaging` 等设置）> 默认值

## 存在性检查

### Exists
```go
// 检查是否存在满足条件的记录
exists, err := repo.Exists(ctx, repository.NewEqFilter("email", "test@example.com"))
if err != nil {
    return err
}
if exists {
    fmt.Println("邮箱已存在")
}
```

### ExistsByFilter
```go
// 使用 ExistsByFilter 方法
exists, err := repo.ExistsByFilter(ctx, repository.NewEqFilter("status", "active"))
```

## 获取首尾记录

### First
```go
// 获取第一条记录（默认按 ID 升序）
first, err := repo.First(ctx)

// 指定排序字段
first, err := repo.First(ctx, "created_at")

// 使用 Query 查询第一条记录（过滤、排序、FilterGroup 等条件都来自 Query）
query := repository.NewQuery().
    AddFilter(repository.NewEqFilter("status", "active")).
    AddOrder("created_at", "DESC")
first, err := repo.FirstWithQuery(ctx, query)
```

`FirstWithQuery` 适合在读取单条记录时复用列表查询里的 `Query` 条件它会应用 Query 中的过滤条件、条件组、排序等设置，并自动限制为 1 条记录

### Last
```go
// 获取最后一条记录
last, err := repo.Last(ctx)
last, err := repo.Last(ctx, "created_at")

// 带条件的 Last
query := repository.NewQuery().AddFilter(repository.NewEqFilter("status", "active"))
last, err := repo.LastWithQuery(ctx, query, "created_at")
```

## 字段提取

### Pluck
```go
// 提取所有用户的邮箱
emails, err := repo.Pluck(ctx, "email")

// 带条件提取
query := repository.NewQuery().AddFilter(repository.NewEqFilter("status", "active"))
emails, err := repo.PluckWithQuery(ctx, "email", query)

// 去重提取
emails, err := repo.PluckUnique(ctx, "email")
```

## 完整示例

```go
package main

import (
    "context"
    "github.com/kamalyes/go-sqlbuilder/repository"
)

// 根据 ID 获取用户
func getUserByID(ctx context.Context, repo repository.IRepository[User], id uint) (*User, error) {
    return repo.Get(ctx, id)
}

// 获取活跃用户列表
func getActiveUsers(ctx context.Context, repo repository.IRepository[User]) ([]*User, error) {
    query := repository.NewQuery().
        AddFilter(repository.NewEqFilter("status", "active")).
        AddOrder("created_at", "DESC")
    return repo.List(ctx, query)
}

// 分页查询用户
func getUsersWithPagination(ctx context.Context, repo repository.IRepository[User], page, pageSize int) ([]*User, *repository.Pagination, error) {
    query := repository.NewQuery().AddOrder("id", "DESC")
    pagination := &repository.Pagination{
        Page:     page,
        PageSize: pageSize,
    }
    return repo.ListWithPagination(ctx, query, pagination)
}

// 检查邮箱是否已存在
func isEmailExists(ctx context.Context, repo repository.IRepository[User], email string) (bool, error) {
    return repo.Exists(ctx, repository.NewEqFilter("email", email))
}
```
