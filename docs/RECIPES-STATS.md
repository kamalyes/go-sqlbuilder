# 业务实战：统计报表

复制即用的统计代码。

## 基础统计

```go
// GetUserStats 用户统计
func (s *UserService) GetUserStats(ctx context.Context) (*UserStats, error) {
    total, _ := s.repo.Count(ctx, repository.NewQuery())
    active, _ := s.repo.Count(ctx, repository.NewQuery().AddEqual("status", "active"))
    today, _ := s.repo.Count(ctx, repository.NewQuery().AddToday("created_at"))
    
    return &UserStats{
        Total:    total,
        Active:   active,
        TodayNew: today,
    }, nil
}

type UserStats struct {
    Total    int64 `json:"total"`
    Active   int64 `json:"active"`
    TodayNew int64 `json:"today_new"`
}
```

## 分组统计

```go
// CountByStatus 按状态分组统计
func (s *UserService) CountByStatus(ctx context.Context) (map[string]int64, error) {
    type Result struct {
        Status string
        Count  int64
    }
    
    var results []Result
    err := s.repo.GetDB().Model(&User{}).
        Select("status, COUNT(*) as count").
        Group("status").
        Scan(&results).Error
    
    stats := make(map[string]int64)
    for _, r := range results {
        stats[r.Status] = r.Count
    }
    return stats, err
}
```

## 时间维度统计

```go
// GetDailyNewUsers 每日新增统计
func (s *UserService) GetDailyNewUsers(ctx context.Context, days int) ([]*DailyStats, error) {
    if days <= 0 {
        days = 30
    }
    
    type Result struct {
        Date  string
        Count int64
    }
    
    var results []Result
    err := s.repo.GetDB().Model(&User{}).
        Select("DATE(created_at) as date, COUNT(*) as count").
        Where("created_at >= ?", time.Now().AddDate(0, 0, -days)).
        Group("DATE(created_at)").
        Order("date ASC").
        Scan(&results).Error
    
    stats := make([]*DailyStats, len(results))
    for i, r := range results {
        stats[i] = &DailyStats{Date: r.Date, Count: r.Count}
    }
    return stats, err
}

type DailyStats struct {
    Date  string `json:"date"`
    Count int64  `json:"count"`
}
```

## 仪表盘数据（并发优化）

```go
// GetDashboardStats 仪表盘统计（并发查询）
func (s *UserService) GetDashboardStats(ctx context.Context) (*DashboardStats, error) {
    queries := repository.ConcurrentQuery{
        {Name: "total", Fn: func(ctx context.Context) (interface{}, error) {
            return s.repo.Count(ctx, repository.NewQuery())
        }},
        {Name: "today", Fn: func(ctx context.Context) (interface{}, error) {
            return s.repo.Count(ctx, repository.NewQuery().AddToday("created_at"))
        }},
        {Name: "week", Fn: func(ctx context.Context) (interface{}, error) {
            return s.repo.Count(ctx, repository.NewQuery().AddThisWeek("created_at"))
        }},
        {Name: "month", Fn: func(ctx context.Context) (interface{}, error) {
            return s.repo.Count(ctx, repository.NewQuery().AddThisMonth("created_at"))
        }},
    }
    
    results := queries.Execute(ctx)
    
    return &DashboardStats{
        Total:    results["total"].(int64),
        TodayNew: results["today"].(int64),
        WeekNew:  results["week"].(int64),
        MonthNew: results["month"].(int64),
    }, nil
}

type DashboardStats struct {
    Total    int64 `json:"total"`
    TodayNew int64 `json:"today_new"`
    WeekNew  int64 `json:"week_new"`
    MonthNew int64 `json:"month_new"`
}
```

## 条件聚合统计

```go
// GetOrderStats 订单统计（使用条件聚合）
func (s *OrderService) GetOrderStats(ctx context.Context, startDate, endDate time.Time) ([]*OrderStats, error) {
    type Result struct {
        Date        string
        TotalCount  int64
        TotalAmount float64
        PaidCount   int64
        PaidAmount  float64
    }
    
    var results []Result
    err := s.repo.GetDB().Raw(`
        SELECT 
            DATE(created_at) as date,
            COUNT(*) as total_count,
            SUM(total_amount) as total_amount,
            SUM(CASE WHEN status = 'paid' THEN 1 ELSE 0 END) as paid_count,
            SUM(CASE WHEN status = 'paid' THEN total_amount ELSE 0 END) as paid_amount
        FROM orders
        WHERE created_at BETWEEN ? AND ?
        GROUP BY DATE(created_at)
        ORDER BY date ASC
    `, startDate, endDate).Scan(&results).Error
    
    stats := make([]*OrderStats, len(results))
    for i, r := range results {
        stats[i] = &OrderStats{
            Date:         r.Date,
            TotalCount:   r.TotalCount,
            TotalAmount:  r.TotalAmount,
            PaidCount:    r.PaidCount,
            PaidAmount:   r.PaidAmount,
            UnpaidCount:  r.TotalCount - r.PaidCount,
            UnpaidAmount: r.TotalAmount - r.PaidAmount,
        }
    }
    return stats, err
}

type OrderStats struct {
    Date         string  `json:"date"`
    TotalCount   int64   `json:"total_count"`
    TotalAmount  float64 `json:"total_amount"`
    PaidCount    int64   `json:"paid_count"`
    PaidAmount   float64 `json:"paid_amount"`
    UnpaidCount  int64   `json:"unpaid_count"`
    UnpaidAmount float64 `json:"unpaid_amount"`
}
```
