# 统计方法

## 概述

BaseRepository 提供了丰富的统计方法，用于数据分析和聚合操作。

## 计数方法

### Count - 统计记录数

```go
// 统计所有记录
count, err := repo.Count(ctx)

// 带条件统计
count, err := repo.Count(ctx, repository.NewEqFilter("status", "active"))

// 多条件统计
count, err := repo.Count(ctx,
    repository.NewEqFilter("status", "active"),
    repository.NewGteFilter("age", 18),
)
```

### CountByField - 按字段分组统计

```go
// 按状态统计用户数
stats, err := repo.CountByField(ctx, "status")
// 返回: [{Value: "active", Count: 100}, {Value: "inactive", Count: 50}]

// 带条件过滤的分组统计
query := repository.NewQuery().AddFilter(repository.NewGteFilter("created_at", "2024-01-01"))
stats, err := repo.CountByFieldWithQuery(ctx, "status", query)
```

## 聚合方法

### Sum - 求和

```go
// 计算所有订单的总金额
total, err := repo.Sum(ctx, "amount")

// 带条件求和
query := repository.NewQuery().AddFilter(repository.NewEqFilter("status", "completed"))
total, err := repo.SumWithQuery(ctx, "amount", query)
```

### Avg - 平均值

```go
// 计算平均年龄
avg, err := repo.Avg(ctx, "age")

// 带条件平均值
query := repository.NewQuery().AddFilter(repository.NewEqFilter("status", "active"))
avg, err := repo.AvgWithQuery(ctx, "age", query)
```

### Max - 最大值

```go
// 获取最大年龄
maxAge, err := repo.Max(ctx, "age")

// 带条件最大值
query := repository.NewQuery().AddFilter(repository.NewEqFilter("status", "active"))
maxAge, err := repo.MaxWithQuery(ctx, "age", query)
```

### Min - 最小值

```go
// 获取最小年龄
minAge, err := repo.Min(ctx, "age")

// 带条件最小值
query := repository.NewQuery().AddFilter(repository.NewEqFilter("status", "active"))
minAge, err := repo.MinWithQuery(ctx, "age", query)
```

## 数据提取

### Pluck - 提取字段值

```go
// 提取所有用户的邮箱
emails, err := repo.Pluck(ctx, "email")

// 带条件提取
query := repository.NewQuery().AddFilter(repository.NewEqFilter("status", "active"))
emails, err := repo.PluckWithQuery(ctx, "email", query)

// 去重提取
emails, err := repo.PluckUnique(ctx, "email")

// 带条件的去重提取
emails, err := repo.PluckUniqueWithQuery(ctx, "email", query)
```

## 完整示例

```go
package main

import (
    "context"
    "fmt"
    "github.com/kamalyes/go-sqlbuilder/repository"
)

// 获取用户统计
func getUserStats(ctx context.Context, repo repository.IRepository[User]) (*UserStats, error) {
    // 总用户数
    total, err := repo.Count(ctx)
    if err != nil {
        return nil, err
    }
    
    // 活跃用户数
    activeCount, err := repo.Count(ctx, repository.NewEqFilter("status", "active"))
    if err != nil {
        return nil, err
    }
    
    // 按状态分组统计
    statusStats, err := repo.CountByField(ctx, "status")
    if err != nil {
        return nil, err
    }
    
    // 平均年龄
    avgAge, err := repo.Avg(ctx, "age")
    if err != nil {
        return nil, err
    }
    
    return &UserStats{
        Total:         total,
        ActiveCount:   activeCount,
        StatusStats:   statusStats,
        AverageAge:    avgAge,
    }, nil
}

// 获取订单统计
func getOrderStats(ctx context.Context, repo repository.IRepository[Order]) (*OrderStats, error) {
    // 已完成订单总额
    completedQuery := repository.NewQuery().
        AddFilter(repository.NewEqFilter("status", "completed"))
    completedTotal, err := repo.SumWithQuery(ctx, "amount", completedQuery)
    if err != nil {
        return nil, err
    }
    
    // 平均订单金额
    avgAmount, err := repo.Avg(ctx, "amount")
    if err != nil {
        return nil, err
    }
    
    // 最大订单金额
    maxAmount, err := repo.Max(ctx, "amount")
    if err != nil {
        return nil, err
    }
    
    return &OrderStats{
        CompletedTotal: completedTotal,
        AverageAmount:  avgAmount,
        MaxAmount:      maxAmount,
    }, nil
}

// 提取所有邮箱
func getAllEmails(ctx context.Context, repo repository.IRepository[User]) ([]interface{}, error) {
    return repo.PluckUnique(ctx, "email")
}
```
