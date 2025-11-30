# 模型定义

go-sqlbuilder 提供了几个内置的基础模型，简化常见的数据表结构定义。

## BaseModel

最基础的模型，包含主键、时间戳、软删除、状态管理和版本控制。

### 定义

```go
type BaseModel struct {
    ID        uint           `json:"id" gorm:"primaryKey;autoIncrement;comment:自增主键"`
    Version   int            `json:"version" gorm:"default:1;comment:版本号"`
    CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime;comment:创建时间"`
    UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime;comment:更新时间"`
    DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index;comment:删除时间"`
    Remark    string         `json:"remark,omitempty" gorm:"type:varchar(500);comment:备注"`
    Status    int8           `json:"status" gorm:"default:1;index;comment:状态(1:启用 0:禁用)"`
}
```

### 使用

```go
type User struct {
    repository.BaseModel
    Username string `gorm:"type:varchar(50);uniqueIndex" json:"username"`
    Email    string `gorm:"type:varchar(100);uniqueIndex" json:"email"`
    Age      int    `gorm:"type:int" json:"age"`
}

// 自动迁移（必须）
db.AutoMigrate(&User{})

// 创建用户时，时间字段会自动填充
user := &User{
    Username: "zhangsan",
    Email:    "zhangsan@example.com",
    Age:      25,
}
repo.Create(ctx, user)
// user.CreatedAt 和 user.UpdatedAt 会自动设置
// user.Status 默认为 1（启用）
// user.Version 默认为 1

// 更新时，UpdatedAt 会自动更新
user.Age = 26
repo.Update(ctx, user)
// user.UpdatedAt 会自动更新为当前时间
// user.Version 会自动 +1（通过 BeforeUpdate 钩子）
```

### 字段说明

| 字段 | 类型 | 说明 |
|------|------|------|
| ID | uint | 主键，自增 |
| Version | int | 版本号，默认 1，更新时自动 +1 |
| CreatedAt | time.Time | 创建时间，自动填充 |
| UpdatedAt | time.Time | 更新时间，自动更新 |
| DeletedAt | gorm.DeletedAt | 软删除时间，默认 NULL |
| Remark | string | 备注，最大 500 字符 |
| Status | int8 | 状态，默认 1（启用），0（禁用）|

### 便利方法

```go
// 状态管理
user.Enable()              // 设置 Status = 1
user.Disable()             // 设置 Status = 0
isEnabled := user.IsEnabled()  // 检查是否启用

// 版本控制
version := user.GetVersion()   // 获取版本号

// 软删除检查
isDeleted := user.IsDeleted()  // 检查是否已删除

// 新记录检查
isNew := user.IsNew()          // ID == 0 时返回 true

// 备注
user.SetRemark("这是一个测试用户")
```

## AuditModel

审计模型，在 BaseModel 基础上增加创建者和更新者字段。

### 定义

```go
type AuditModel struct {
    BaseModel
    CreatedBy uint `gorm:"type:bigint;comment:创建人ID" json:"created_by"`
    UpdatedBy uint `gorm:"type:bigint;comment:更新人ID" json:"updated_by"`
}
```

### 使用

```go
type Article struct {
    repository.AuditModel
    Title   string `gorm:"type:varchar(200)" json:"title"`
    Content string `gorm:"type:text" json:"content"`
    Status  string `gorm:"type:varchar(20)" json:"status"`
}

// 创建文章时设置创建者
article := &Article{
    Title:   "我的文章",
    Content: "文章内容",
    Status:  "draft",
}
article.CreatedBy = currentUserID
article.UpdatedBy = currentUserID

repo.Create(ctx, article)

// 更新文章时设置更新者
article.Title = "更新后的标题"
article.UpdatedBy = currentUserID
repo.Update(ctx, article)
```

### 完整示例：权限管理

```go
type Permission struct {
    repository.AuditModel
    Name        string `gorm:"type:varchar(100);uniqueIndex" json:"name"`
    Description string `gorm:"type:varchar(255)" json:"description"`
    Module      string `gorm:"type:varchar(50)" json:"module"`
}

// 初始化时必须迁移
func init() {
    db.AutoMigrate(&Permission{})
}

func CreatePermission(ctx context.Context, name, desc, module string, userID uint) error {
    perm := &Permission{
        Name:        name,
        Description: desc,
        Module:      module,
    }
    perm.SetCreatedBy(userID)
    perm.SetUpdatedBy(userID)
    
    _, err := repo.Create(ctx, perm)
    // perm.CreatedAt 和 perm.UpdatedAt 已自动设置
    return err
}

func UpdatePermission(ctx context.Context, id uint, desc string, userID uint) error {
    perm, err := repo.Get(ctx, id)
    if err != nil {
        return err
    }
    
    perm.Description = desc
    perm.SetUpdatedBy(userID)
    
    _, err = repo.Update(ctx, perm)
    // perm.UpdatedAt 会自动更新
    // perm.Version 会自动 +1
    return err
}
```

### 完整示例：用户管理（展示所有特性）

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/kamalyes/go-logger"
    "github.com/kamalyes/go-sqlbuilder/db"
    "github.com/kamalyes/go-sqlbuilder/repository"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

type User struct {
    repository.BaseModel
    Username string `gorm:"type:varchar(50);uniqueIndex" json:"username"`
    Email    string `gorm:"type:varchar(100);uniqueIndex" json:"email"`
    Age      int    `gorm:"type:int" json:"age"`
}

func main() {
    // 1. 连接数据库
    dsn := "user:password@tcp(127.0.0.1:3306)/testdb?charset=utf8mb4&parseTime=True"
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal(err)
    }
    
    // 2. 自动迁移（必须！）
    if err := db.AutoMigrate(&User{}); err != nil {
        log.Fatal(err)
    }
    
    // 3. 创建 Repository
    handler, _ := db.NewGormHandler(gormDB)
    logger := logger.NewLogger(nil)
    repo := repository.NewBaseRepository[User](handler, logger, "users")
    
    ctx := context.Background()
    
    // 4. 创建用户
    user := &User{
        Username: "zhangsan",
        Email:    "zhangsan@example.com",
        Age:      25,
    }
    
    created, err := repo.Create(ctx, user)
    if err != nil {
        log.Fatal(err)
    }
    
    // ✅ 时间字段已自动填充
    fmt.Printf("创建成功:\n")
    fmt.Printf("  ID: %d\n", created.ID)
    fmt.Printf("  Username: %s\n", created.Username)
    fmt.Printf("  CreatedAt: %s\n", created.CreatedAt)  // 自动设置
    fmt.Printf("  UpdatedAt: %s\n", created.UpdatedAt)  // 自动设置
    fmt.Printf("  Status: %d\n", created.Status)        // 默认 1
    fmt.Printf("  Version: %d\n", created.Version)      // 默认 1
    
    // 5. 更新用户
    created.Age = 26
    created.SetRemark("更新了年龄")
    
    updated, err := repo.Update(ctx, created)
    if err != nil {
        log.Fatal(err)
    }
    
    // ✅ UpdatedAt 自动更新，Version 自动 +1
    fmt.Printf("\n更新成功:\n")
    fmt.Printf("  Age: %d\n", updated.Age)
    fmt.Printf("  UpdatedAt: %s\n", updated.UpdatedAt)  // 自动更新
    fmt.Printf("  Version: %d\n", updated.Version)      // 自动 +1 = 2
    fmt.Printf("  Remark: %s\n", updated.Remark)
    
    // 6. 状态管理
    updated.Disable()
    repo.Update(ctx, updated)
    fmt.Printf("\n用户已禁用: Status = %d\n", updated.Status)
    
    // 7. 软删除
    err = repo.SoftDelete(ctx, updated.ID, "deleted_at", time.Now())
    if err != nil {
        log.Fatal(err)
    }
    
    // 8. 检查是否已删除
    deletedUser, _ := db.Unscoped().First(&User{}, updated.ID)
    fmt.Printf("\n软删除检查:\n")
    fmt.Printf("  IsDeleted: %v\n", deletedUser.IsDeleted())
    fmt.Printf("  DeletedAt: %v\n", deletedUser.DeletedAt)
}
```

## 自定义模型

### 基础自定义

```go
type Product struct {
    ID          uint      `gorm:"primaryKey"`
    Name        string    `gorm:"type:varchar(100);not null"`
    Price       float64   `gorm:"type:decimal(10,2);not null"`
    Stock       int       `gorm:"type:int;default:0"`
    Category    string    `gorm:"type:varchar(50);index"`
    IsAvailable bool      `gorm:"type:tinyint(1);default:1"`
    CreatedAt   time.Time `gorm:"autoCreateTime"`
    UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}
```

### 组合使用

```go
type Order struct {
    repository.BaseModel
    UserID      uint           `gorm:"type:bigint;not null;index" json:"user_id"`
    OrderNo     string         `gorm:"type:varchar(50);uniqueIndex" json:"order_no"`
    TotalAmount float64        `gorm:"type:decimal(10,2)" json:"total_amount"`
    Status      string         `gorm:"type:varchar(20);index" json:"status"`
    PaidAt      *time.Time     `gorm:"type:datetime" json:"paid_at"`
    ShippedAt   *time.Time     `gorm:"type:datetime" json:"shipped_at"`
    CompletedAt *time.Time     `gorm:"type:datetime" json:"completed_at"`
}
```

### 嵌套关联

```go
type User struct {
    repository.BaseModel
    Username string    `gorm:"type:varchar(50);uniqueIndex" json:"username"`
    Email    string    `gorm:"type:varchar(100);uniqueIndex" json:"email"`
    Profile  Profile   `gorm:"foreignKey:UserID" json:"profile"`
    Orders   []Order   `gorm:"foreignKey:UserID" json:"orders"`
}

type Profile struct {
    ID        uint   `gorm:"primaryKey"`
    UserID    uint   `gorm:"type:bigint;uniqueIndex" json:"user_id"`
    Avatar    string `gorm:"type:varchar(255)" json:"avatar"`
    Bio       string `gorm:"type:text" json:"bio"`
    CreatedAt time.Time
    UpdatedAt time.Time
}

// 预加载关联数据
query := repository.NewQuery().AddPreload("Profile").AddPreload("Orders")
users, err := repo.List(ctx, query)
```

## 字段标签

### GORM 标签

```go
type Example struct {
    // 主键
    ID uint `gorm:"primaryKey;autoIncrement"`
    
    // 类型和长度
    Name string `gorm:"type:varchar(100)"`
    
    // 非空
    Email string `gorm:"not null"`
    
    // 唯一索引
    Username string `gorm:"uniqueIndex"`
    
    // 普通索引
    Status string `gorm:"index"`
    
    // 默认值
    Role string `gorm:"default:user"`
    
    // 注释
    Age int `gorm:"comment:用户年龄"`
    
    // 外键
    UserID uint `gorm:"foreignKey:UserID"`
    
    // 忽略字段（不存入数据库）
    TempData string `gorm:"-"`
}
```

### JSON 标签

```go
type User struct {
    ID       uint   `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`
    Password string `json:"-"`                    // 不序列化到 JSON
    Age      int    `json:"age,omitempty"`        // 零值时省略
    Profile  string `json:"profile,omitempty"`
}
```

## 其他内置模型

### SimpleModel - 简化模型

不包含软删除、版本控制和状态管理，仅保留基本字段。

```go
type SimpleModel struct {
    ID        uint      `json:"id" gorm:"primaryKey;autoIncrement;comment:自增主键"`
    CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime;comment:创建时间"`
    UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime;comment:更新时间"`
}

// 使用
type Log struct {
    repository.SimpleModel
    Message string `gorm:"type:text"`
    Level   string `gorm:"type:varchar(20)"`
}
```

### LightModel - 轻量级模型

包含状态管理，但不包含软删除和版本控制。

```go
type LightModel struct {
    ID        uint      `json:"id" gorm:"primaryKey;autoIncrement;comment:自增主键"`
    CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime;comment:创建时间"`
    UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime;comment:更新时间"`
    Status    int8      `json:"status" gorm:"default:1;index;comment:状态(1:启用 0:禁用)"`
}

// 使用
type Config struct {
    repository.LightModel
    Key   string `gorm:"type:varchar(100);uniqueIndex"`
    Value string `gorm:"type:text"`
}
```

### UUIDModel - UUID 主键模型

使用 UUID 作为主键，适合分布式系统。

```go
type UUIDModel struct {
    ID        string         `json:"id" gorm:"primaryKey;type:char(36);comment:UUID主键"`
    Version   int            `json:"version" gorm:"default:1;comment:版本号"`
    CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime;comment:创建时间"`
    UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime;comment:更新时间"`
    DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index;comment:删除时间"`
}

// 使用
import "github.com/google/uuid"

type Order struct {
    repository.UUIDModel
    OrderNo string  `gorm:"type:varchar(50);uniqueIndex"`
    Amount  float64 `gorm:"type:decimal(10,2)"`
}

// 创建时手动设置 UUID
order := &Order{
    OrderNo: "ORD123456",
    Amount:  99.99,
}
order.ID = uuid.New().String()
repo.Create(ctx, order)
```

## 最佳实践

### 1. 使用 BaseModel 作为基础（需要迁移）

```go
// ✅ 推荐：使用 BaseModel
type User struct {
    repository.BaseModel
    Username string `gorm:"type:varchar(50);uniqueIndex"`
    Email    string `gorm:"type:varchar(100);uniqueIndex"`
}

// ⚠️ 重要：必须进行数据库迁移
db.AutoMigrate(&User{})

// 这样 CreatedAt、UpdatedAt 才会自动填充
user := &User{Username: "test", Email: "test@test.com"}
repo.Create(ctx, user)
// ✅ user.CreatedAt 和 user.UpdatedAt 已自动设置

// ❌ 不推荐：手动定义所有字段（重复且容易出错）
type User struct {
    ID        uint
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt
    Username  string
    Email     string
}
```

### 1.1 时间自动填充的前提条件

GORM 的 `autoCreateTime` 和 `autoUpdateTime` 只在以下情况生效：

1. **已执行数据库迁移**：`db.AutoMigrate(&YourModel{})`
2. **使用 GORM 创建/更新**：通过 Repository 或 GORM 的 Create/Update 方法
3. **字段为零值**：CreatedAt/UpdatedAt 为零值时才自动填充

```go
// ✅ 正确：时间会自动填充
user := &User{Username: "test"}
repo.Create(ctx, user)

// ❌ 错误：手动设置了时间，不会自动填充
user := &User{
    Username:  "test",
    CreatedAt: time.Now(),  // 手动设置，GORM 不会覆盖
}
repo.Create(ctx, user)

// ✅ 如果需要手动设置时间
user := &User{Username: "test"}
user.SetCreatedAt(customTime)
user.SetUpdatedAt(customTime)
repo.Create(ctx, user)
```

### 2. 需要审计时使用 AuditModel

```go
// 重要数据使用 AuditModel
type Transaction struct {
    repository.AuditModel
    Amount float64
    Type   string
}

// 普通数据使用 BaseModel
type Log struct {
    repository.BaseModel
    Message string
    Level   string
}
```

### 3. 合理使用索引

```go
type User struct {
    repository.BaseModel
    Email    string `gorm:"uniqueIndex"`              // 唯一索引
    Username string `gorm:"uniqueIndex"`              // 唯一索引
    Status   string `gorm:"index"`                    // 普通索引
    Age      int    `gorm:"index:idx_age_status"`     // 组合索引
    City     string `gorm:"index:idx_age_status"`     // 组合索引
}
```

### 4. 使用枚举常量

```go
const (
    UserStatusActive   = "active"
    UserStatusInactive = "inactive"
    UserStatusBanned   = "banned"
)

type User struct {
    repository.BaseModel
    Username string `gorm:"type:varchar(50)"`
    Status   string `gorm:"type:varchar(20);default:active"`
}

func (u *User) IsActive() bool {
    return u.Status == UserStatusActive
}
```

### 5. 添加验证方法

```go
type User struct {
    repository.BaseModel
    Email    string `gorm:"type:varchar(100)"`
    Age      int    `gorm:"type:int"`
    Password string `gorm:"type:varchar(255)"`
}

func (u *User) Validate() error {
    if u.Email == "" {
        return errors.New("邮箱不能为空")
    }
    
    if u.Age < 18 {
        return errors.New("年龄必须大于 18")
    }
    
    if len(u.Password) < 6 {
        return errors.New("密码长度至少 6 位")
    }
    
    return nil
}

// 使用
user := &User{Email: "test@test.com", Age: 20, Password: "123456"}
if err := user.Validate(); err != nil {
    return err
}
repo.Create(ctx, user)
```

## 迁移

### 自动迁移

```go
db.AutoMigrate(&User{}, &Article{}, &Order{})
```

### 手动迁移

```go
// 创建表
db.Migrator().CreateTable(&User{})

// 删除表
db.Migrator().DropTable(&User{})

// 添加列
db.Migrator().AddColumn(&User{}, "nickname")

// 删除列
db.Migrator().DropColumn(&User{}, "nickname")

// 添加索引
db.Migrator().CreateIndex(&User{}, "Email")
```

## 📚 相关文档

- 📖 [Repository 基础](./REPOSITORY-BASICS.md) - 学习 CRUD 操作
- 🔍 [高级查询](./ADVANCED-QUERIES.md) - 复杂查询构建
- 🎯 [FilterGroup](./FILTERGROUP.md) - 复杂条件组合
