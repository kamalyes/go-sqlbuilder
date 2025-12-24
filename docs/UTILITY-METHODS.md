# 工具方法

BaseRepository 提供了一系列工具方法，用于统计、检查、提取数据等操作。

## 📖 目录

- [统计方法](#统计方法)
- [检查方法](#检查方法)
- [数据提取方法](#数据提取方法)

---

## 统计方法

### Count - 统计记录数

统计符合条件的记录总数。

```go
package main

import (
    "context"
    "github.com/kamalyes/go-sqlbuilder/repository"
)

// COUNT(*) WHERE status = 'active'
func countActiveUsers(repo repository.IBaseRepository[User, uint64]) (int64, error) {
    ctx := context.Background()
    
    count, err := repo.Count(ctx, repository.NewEqFilter("status", "active"))
    return count, err
}

// 多条件统计
func countAdultActiveUsers(repo repository.IBaseRepository[User, uint64]) (int64, error) {
    ctx := context.Background()
    
    filters := []*repository.Filter{
        repository.NewEqFilter("status", "active"),
        repository.NewGteFilter("age", 18),
    }
    
    count, err := repo.Count(ctx, filters...)
    return count, err
}
```

### CountByField - 按字段分组统计

按指定字段分组统计，返回每个分组的记录数。

```go
// GROUP BY status, COUNT(*)
func countUsersByStatus(repo repository.IBaseRepository[User, uint64]) (map[interface{}]int64, error) {
    ctx := context.Background()
    
    countMap, err := repo.CountByField(ctx, "status")
    if err != nil {
        return nil, err
    }
    
    // 输出结果
    for status, count := range countMap {
        fmt.Printf("状态 %v: %d 人\n", status, count)
    }
    
    // 示例输出：
    // 状态 active: 120 人
    // 状态 inactive: 30 人
    // 状态 deleted: 5 人
    
    return countMap, nil
}
```

**使用场景：**

- 统计各状态的用户数量
- 统计各部门的员工数量
- 统计各分类的商品数量

---

## 检查方法

### Exists - 检查记录是否存在

检查是否存在符合条件的记录。

```go
// 检查是否存在 status='active' 的记录
func hasActiveUsers(repo repository.IBaseRepository[User, uint64]) (bool, error) {
    ctx := context.Background()
    
    exists, err := repo.Exists(ctx, repository.NewEqFilter("status", "active"))
    return exists, err
}

// 检查邮箱是否已被使用
func isEmailTaken(repo repository.IBaseRepository[User, uint64], email string) (bool, error) {
    ctx := context.Background()
    
    exists, err := repo.Exists(ctx, repository.NewEqFilter("email", email))
    return exists, err
}
```

**性能优势：**

- `Exists` 比查询后判断更高效（使用 `SELECT EXISTS(...)` 或 `SELECT 1 LIMIT 1`）
- 不会加载实际数据，只返回布尔值

```go
// ❌ 低效：查询后判断
user, err := repo.FindOne(ctx, filter)
exists := user != nil

// ✅ 高效：使用 Exists
exists, err := repo.Exists(ctx, filter)
```

---

## 数据提取方法

### GetAll - 获取所有记录

获取表中所有记录。

```go
// SELECT * FROM users
func getAllUsers(repo repository.IBaseRepository[User, uint64]) ([]*User, error) {
    ctx := context.Background()
    
    users, err := repo.GetAll(ctx)
    return users, err
}
```

**⚠️ 注意：** 数据量大时请使用分页查询，避免一次性加载过多数据。

### Pluck - 提取单个字段的值列表

提取指定字段的所有值，返回切片。

```go
// 获取所有活跃用户的邮箱列表
func getActiveUserEmails(repo repository.IBaseRepository[User, uint64]) ([]interface{}, error) {
    ctx := context.Background()
    
    emails, err := repo.Pluck(ctx, "email", repository.NewEqFilter("status", "active"))
    if err != nil {
        return nil, err
    }
    
    // 输出结果
    for _, email := range emails {
        fmt.Println(email.(string))
    }
    
    // 示例输出：
    // user1@example.com
    // user2@example.com
    // user3@example.com
    
    return emails, nil
}

// 类型转换辅助函数
func pluckEmails(repo repository.IBaseRepository[User, uint64]) ([]string, error) {
    ctx := context.Background()
    
    values, err := repo.Pluck(ctx, "email")
    if err != nil {
        return nil, err
    }
    
    emails := make([]string, len(values))
    for i, v := range values {
        emails[i] = v.(string)
    }
    
    return emails, nil
}
```

**使用场景：**

- 获取 ID 列表用于批量操作
- 获取邮箱列表用于发送通知
- 获取名称列表用于下拉菜单

### Distinct - 获取去重字段值列表

获取指定字段的所有不重复值。

```go
// 获取所有不重复的部门名称
func getAllDepartments(repo repository.IBaseRepository[User, uint64]) ([]interface{}, error) {
    ctx := context.Background()
    
    departments, err := repo.Distinct(ctx, "department")
    if err != nil {
        return nil, err
    }
    
    // 输出结果
    for _, dept := range departments {
        fmt.Println(dept.(string))
    }
    
    // 示例输出（去重）：
    // IT
    // HR
    // Finance
    // Sales
    
    return departments, nil
}

// 带条件的去重查询
func getActiveUserDepartments(repo repository.IBaseRepository[User, uint64]) ([]interface{}, error) {
    ctx := context.Background()
    
    departments, err := repo.Distinct(ctx, "department", 
        repository.NewEqFilter("status", "active"))
    
    return departments, nil
}
```

**使用场景：**

- 获取所有标签列表
- 获取所有分类列表
- 获取所有状态值用于筛选器

---

## 💡 实战示例

### 示例 1: 仪表板统计

```go
func getDashboardStats(repo repository.IBaseRepository[User, uint64]) (map[string]interface{}, error) {
    ctx := context.Background()
    stats := make(map[string]interface{})
    
    // 总用户数
    totalUsers, err := repo.Count(ctx)
    if err != nil {
        return nil, err
    }
    stats["total_users"] = totalUsers
    
    // 活跃用户数
    activeUsers, err := repo.Count(ctx, repository.NewEqFilter("status", "active"))
    if err != nil {
        return nil, err
    }
    stats["active_users"] = activeUsers
    
    // 按状态分组统计
    statusCounts, err := repo.CountByField(ctx, "status")
    if err != nil {
        return nil, err
    }
    stats["status_breakdown"] = statusCounts
    
    // VIP 用户数
    vipUsers, err := repo.Count(ctx, repository.NewGtFilter("vip_level", 0))
    if err != nil {
        return nil, err
    }
    stats["vip_users"] = vipUsers
    
    return stats, nil
}
```

### 示例 2: 数据验证

```go
func validateUserData(repo repository.IBaseRepository[User, uint64], email, phone string) error {
    ctx := context.Background()
    
    // 检查邮箱是否已存在
    emailExists, err := repo.Exists(ctx, repository.NewEqFilter("email", email))
    if err != nil {
        return err
    }
    if emailExists {
        return fmt.Errorf("邮箱 %s 已被使用", email)
    }
    
    // 检查手机号是否已存在
    phoneExists, err := repo.Exists(ctx, repository.NewEqFilter("phone", phone))
    if err != nil {
        return err
    }
    if phoneExists {
        return fmt.Errorf("手机号 %s 已被使用", phone)
    }
    
    return nil
}
```

### 示例 3: 批量通知

```go
func sendNotificationToActiveUsers(repo repository.IBaseRepository[User, uint64]) error {
    ctx := context.Background()
    
    // 获取所有活跃用户的邮箱
    values, err := repo.Pluck(ctx, "email", repository.NewEqFilter("status", "active"))
    if err != nil {
        return err
    }
    
    // 转换为字符串切片
    emails := make([]string, len(values))
    for i, v := range values {
        emails[i] = v.(string)
    }
    
    // 发送通知
    for _, email := range emails {
        // 发送邮件逻辑
        fmt.Printf("发送通知到: %s\n", email)
    }
    
    return nil
}
```

### 示例 4: 数据分析

```go
func analyzeUserDistribution(repo repository.IBaseRepository[User, uint64]) error {
    ctx := context.Background()
    
    // 按部门统计用户数
    deptCounts, err := repo.CountByField(ctx, "department")
    if err != nil {
        return err
    }
    
    fmt.Println("部门分布:")
    for dept, count := range deptCounts {
        fmt.Printf("  %v: %d 人\n", dept, count)
    }
    
    // 按年龄段统计（需要自定义查询）
    ageGroups := map[string][2]int{
        "18-25": {18, 25},
        "26-35": {26, 35},
        "36-45": {36, 45},
        "46+":   {46, 999},
    }
    
    fmt.Println("\n年龄分布:")
    for group, ages := range ageGroups {
        filters := []*repository.Filter{
            repository.NewGteFilter("age", ages[0]),
            repository.NewLteFilter("age", ages[1]),
        }
        count, _ := repo.Count(ctx, filters...)
        fmt.Printf("  %s: %d 人\n", group, count)
    }
    
    return nil
}
```

---

## ⚙️ 性能优化建议

### 1. 使用 Exists 而不是查询

```go
// ❌ 低效
user, err := repo.FindOne(ctx, filter)
exists := user != nil

// ✅ 高效
exists, err := repo.Exists(ctx, filter)
```

### 2. 使用 Count 而不是查询全部

```go
// ❌ 低效：查询所有数据再统计
users, _ := repo.List(ctx, query)
count := len(users)

// ✅ 高效：直接使用 Count
count, err := repo.Count(ctx, filters...)
```

### 3. 使用 Pluck 提取需要的字段

```go
// ❌ 低效：查询所有字段
users, _ := repo.List(ctx, query)
for _, user := range users {
    emails = append(emails, user.Email)
}

// ✅ 高效：只提取需要的字段
emails, err := repo.Pluck(ctx, "email", filters...)
```

---

## 📚 相关文档

- [CRUD 操作](./CRUD-OPERATIONS.md) - 基础 CRUD 方法
- [便捷查询方法](./CONVENIENCE-METHODS.md) - 简化的查询 API
- [过滤条件](./FILTERS.md) - Filter 详细说明
- [性能优化建议](./PERFORMANCE-TIPS.md) - 性能优化技巧
