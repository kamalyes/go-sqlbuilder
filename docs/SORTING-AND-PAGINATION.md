# 排序和分页

本文档介绍如何使用排序和分页功能。

## 📖 目录

- [排序](#排序)
- [分页](#分页)
- [限制和偏移](#限制和偏移)

---

## 排序

### 单字段排序

```go
// ORDER BY created_at DESC
query := repository.NewQuery().
    WithSort("created_at", constants.SORT_DESC)

users, err := repo.List(ctx, query)
```

### 多字段排序

```go
// ORDER BY status ASC, created_at DESC
query := repository.NewQuery().
    WithSort("status", constants.SORT_ASC).
    WithSort("created_at", constants.SORT_DESC)

users, err := repo.List(ctx, query)
```

---

## 分页

### 使用 WithPaging

```go
// 第1页，每页20条
query := repository.NewQuery().
    WithPaging(1, 20)

users, err := repo.List(ctx, query)
```

### 使用 ListWithPagination

```go
query := repository.NewQuery().
    AddEqFilter("status", "active")

pagination := &repository.Pagination{Page: 1, PageSize: 20}
users, paginationResult, err := repo.ListWithPagination(ctx, query, pagination)

fmt.Printf("总页数: %d\n", paginationResult.GetTotalPages())
```

### 使用便捷方法 Paginate

```go
users, pagination, err := repo.Paginate(ctx, 1, 20, "status", "active")
```

---

## 限制和偏移

### Limit - 限制结果数量

```go
// LIMIT 10
query := repository.NewQuery().
    Limit(10)

users, err := repo.List(ctx, query)
```

### Offset - 跳过指定数量

```go
// OFFSET 20
query := repository.NewQuery().
    Offset(20)

users, err := repo.List(ctx, query)
```

---

## 📚 相关文档

- [便捷查询方法](./CONVENIENCE-METHODS.md) - Paginate 方法
- [CRUD 操作](./CRUD-OPERATIONS.md) - List 方法
