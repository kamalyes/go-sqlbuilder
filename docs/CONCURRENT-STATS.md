# 并发统计查询 (Concurrent Stats)

本模块提供了强大的并发查询和多表统计功能，支持高效的批量数据统计和聚合操作。

## 目录

- [并发查询执行器](#并发查询执行器)
- [多表统计构建器](#多表统计构建器)
- [时间分组统计](#时间分组统计)
- [条件聚合统计](#条件聚合统计)
- [最佳实践](#最佳实践)

## 并发查询执行器

`ConcurrentQueryExecutor` 提供了通用的并发查询执行框架，支持：

- 并发执行多个查询任务
- 自定义超时控制
- Worker池模式限制并发数
- 详细的执行结果和错误信息

### 基础用法

```go
import (
    "context"
    "github.com/kamalyes/go-sqlbuilder/repository"
)

// 创建执行器
executor := repository.NewConcurrentQueryExecutor(db).
    WithTimeout(30 * time.Second).
    WithWorkers(5) // 限制并发数为5

// 定义查询任务
tasks := []repository.ConcurrentQueryTask[int64]{
    {
        Name: "总用户数",
        Query: func(ctx context.Context) (int64, error) {
            var count int64
            err := db.WithContext(ctx).Model(&User{}).Count(&count).Error
            return count, err
        },
    },
    {
        Name: "活跃用户数",
        Query: func(ctx context.Context) (int64, error) {
            var count int64
            err := db.WithContext(ctx).Model(&User{}).
                Where("status = ?", "active").Count(&count).Error
            return count, err
        },
    },
}

// 执行查询
results, hasError := repository.ExecuteConcurrentQuery(executor, ctx, tasks)

// 处理结果
for _, result := range results {
    if result.Error != nil {
        log.Printf("查询 %s 失败: %v", result.Name, result.Error)
    } else {
        log.Printf("%s: %d", result.Name, result.Value)
    }
}
```

### 简化接口

对于返回相同类型的多个查询，可以使用简化接口：

```go
queries := map[string]func(ctx context.Context) (int64, error){
    "总数": func(ctx context.Context) (int64, error) {
        var count int64
        err := db.WithContext(ctx).Model(&User{}).Count(&count).Error
        return count, err
    },
    "求和": func(ctx context.Context) (int64, error) {
        var sum int64
        err := db.WithContext(ctx).Model(&User{}).
            Select("SUM(age)").Scan(&sum).Error
        return sum, err
    },
}

resultMap, hasError := repository.ConcurrentSimpleQuery(executor, ctx, queries)
// resultMap["总数"] = 1000
// resultMap["求和"] = 25000
```

## 多表统计构建器

`MultiTableStatsBuilder` 专门用于多表并发统计查询，提供了流式API。

### 基础统计

```go
builder := repository.NewMultiTableStatsBuilder(ctx, db, logger)
stats, err := builder.
    Count("users", "total_users").                    // COUNT(*)
    Count("posts", "total_posts").
    CountDistinct("posts", "user_id", "active_authors"). // COUNT(DISTINCT user_id)
    Sum("orders", "amount", "total_revenue").         // SUM(amount)
    Execute()

if err != nil {
    log.Fatal(err)
}

fmt.Printf("总用户数: %d\n", stats["total_users"])
fmt.Printf("总文章数: %d\n", stats["total_posts"])
fmt.Printf("活跃作者: %d\n", stats["active_authors"])
fmt.Printf("总收入: %d\n", stats["total_revenue"])
```

### 带条件的统计

```go
now := time.Now()
yesterday := now.Add(-24 * time.Hour)

stats, err := builder.
    WithTimeRange(yesterday, now).  // 时间范围过滤
    AddCondition("users", "status = ?", "active"). // 添加自定义条件
    Count("users", "active_users_today").
    Sum("orders", "amount", "today_revenue").
    Execute()
```

### 性能优化

```go
stats, err := builder.
    WithTimeout(60 * time.Second). // 设置查询超时
    WithWorkers(10).                // 限制并发数
    Count("users", "total_users").
    Count("posts", "total_posts").
    Count("comments", "total_comments").
    Execute()
```

### 获取详细结果

当需要处理部分查询失败的情况时：

```go
results, hasError, err := builder.
    Count("users", "total_users").
    Count("nonexistent_table", "will_fail").
    ExecuteWithDetails()

// err == nil，即使有查询失败
// hasError == true，表示有查询失败
// results 包含所有查询的详细结果

for _, result := range results {
    if result.Error != nil {
        log.Printf("查询 %s 失败: %v", result.Name, result.Error)
    } else {
        log.Printf("%s 成功: %d", result.Name, result.Value)
    }
}
```

## 时间分组统计

`TimeGroupStatsBuilder` 用于按时间分组的统计分析。

### 按天统计

```go
builder := repository.NewTimeGroupStatsBuilder(ctx, db, logger).
    SetModel(&User{}).
    SetTimeField("created_at").
    SetTimeRange(startDate, endDate).
    GroupByDay()

// 添加聚合字段
builder.Count("user_count").
    Sum("age", "total_age").
    Avg("age", "avg_age").
    Max("age", "max_age").
    Min("age", "min_age")

results, err := builder.Execute()

for _, row := range results {
    fmt.Printf("日期: %s, 用户数: %d, 平均年龄: %.2f\n",
        row["time_group"], row["user_count"], row["avg_age"])
}
```

### 按小时统计

```go
builder := repository.NewTimeGroupStatsBuilder(ctx, db, logger).
    SetModel(&Order{}).
    SetTimeField("created_at").
    GroupByHour().
    Sum("amount", "hourly_revenue").
    Count("order_count")

results, err := builder.Execute()
```

### 自定义时间格式

```go
builder.GroupByCustom("%Y-%m") // 按月统计
builder.GroupByCustom("%Y-W%W") // 按周统计
```

### 多维分组

```go
// 按时间和用户ID分组
builder := repository.NewTimeGroupStatsBuilder(ctx, db, logger).
    SetModel(&Post{}).
    GroupByDay().
    AddGroupByField("user_id"). // 额外的分组字段
    Count("post_count")

results, err := builder.Execute()
// 结果按 (日期, user_id) 分组
```

## 条件聚合统计

`ConditionalAggregateBuilder` 用于在同一查询中进行多个条件聚合。

### 计数聚合

```go
builder := repository.NewConditionalAggregateBuilder(ctx, db, logger).
    SetModel(&User{})

// 添加条件计数
builder.CountWhen("age >= 18", "adult_count").
    CountWhen("age < 18", "minor_count").
    CountWhen("status = 'active'", "active_count")

result, err := builder.Execute()

fmt.Printf("成年人: %d\n", result["adult_count"])
fmt.Printf("未成年: %d\n", result["minor_count"])
fmt.Printf("活跃用户: %d\n", result["active_count"])
```

### 值聚合

```go
builder := repository.NewConditionalAggregateBuilder(ctx, db, logger).
    SetModel(&Order{})

// 不同条件下的总金额
builder.SumWhen("status = 'paid'", "amount", "paid_amount").
    SumWhen("status = 'pending'", "amount", "pending_amount").
    AvgWhen("status = 'paid'", "amount", "avg_paid_amount")

result, err := builder.Execute()
```

### 复杂条件

```go
builder.CountWhen("age >= 18 AND status = 'active'", "active_adults").
    SumWhen("age >= 30 AND created_at > ?", "age", "senior_age_sum", yesterday).
    MaxWhen("is_premium = true", "score", "max_premium_score")
```

## 最佳实践

### 1. 设置合理的超时

```go
// 对于复杂查询，增加超时时间
builder := repository.NewMultiTableStatsBuilder(ctx, db, logger).
    WithTimeout(60 * time.Second)
```

### 2. 控制并发数

```go
// 避免过多并发导致数据库压力
builder.WithWorkers(10) // 限制为10个并发
```

### 3. 使用时间范围过滤

```go
// 减少查询数据量
builder.WithTimeRange(startTime, endTime)
```

### 4. 批量统计优化

```go
// 一次性执行多个统计，而不是多次查询
stats, err := builder.
    Count("users", "total_users").
    Count("posts", "total_posts").
    Count("comments", "total_comments").
    Sum("orders", "amount", "total_revenue").
    Execute()
```

### 5. 错误处理

```go
// 使用 ExecuteWithDetails 获取详细错误信息
results, hasError, err := builder.
    Count("table1", "count1").
    Count("table2", "count2").
    ExecuteWithDetails()

if err != nil {
    return err
}

if hasError {
    // 部分查询失败，记录日志
    for _, result := range results {
        if result.Error != nil {
            log.Printf("查询失败: %s - %v", result.Name, result.Error)
        }
    }
}
```

### 6. 索引优化

确保统计查询涉及的字段都有适当的索引：

```sql
-- 时间范围查询
CREATE INDEX idx_created_at ON users(created_at);

-- 条件过滤
CREATE INDEX idx_status ON users(status);

-- 时间+状态组合
CREATE INDEX idx_created_status ON users(created_at, status);
```

### 7. 大数据量处理

```go
// 对于大数据量，使用分页或限制
builder := repository.NewTimeGroupStatsBuilder(ctx, db, logger).
    SetModel(&User{}).
    GroupByDay().
    OrderBy("time_group", "DESC").
    Limit(30) // 只获取最近30天
```

## 性能建议

1. **避免全表扫描**：始终使用时间范围或条件过滤
2. **合理使用并发**：根据数据库连接池大小设置worker数
3. **索引优化**：为常用的过滤和分组字段建立索引
4. **结果缓存**：对于不常变化的统计数据，考虑缓存结果
5. **监控查询时间**：使用日志记录慢查询，及时优化

## 相关文档

- [并发查询详解](CONCURRENT-QUERIES.md)
- [查询示例](QUERY-EXAMPLES.md)
- [性能优化](PERFORMANCE.md)
