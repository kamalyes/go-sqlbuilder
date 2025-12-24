# 字段选择

本文档介绍如何使用 SELECT 和 DISTINCT 功能。

## 选择指定字段

```go
// SELECT id, name, email FROM users
query := repository.NewQuery().
    Select("id", "name", "email")

users, err := repo.List(ctx, query)
```

## 去重查询

```go
// SELECT DISTINCT status FROM users
query := repository.NewQuery().
    WithDistinct().
    Select("status")

users, err := repo.List(ctx, query)
```

## 组合使用

```go
// SELECT DISTINCT department FROM users WHERE status = 'active'
query := repository.NewQuery().
    WithDistinct().
    Select("department").
    AddEqFilter("status", "active")

users, err := repo.List(ctx, query)
```

---

## 📚 相关文档

- [CRUD 操作](./CRUD-OPERATIONS.md) - List 方法
- [工具方法](./UTILITY-METHODS.md) - Pluck, Distinct 方法
