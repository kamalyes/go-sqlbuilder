# 数据库迁移器 (Migrator)

数据库迁移器提供了一套完整的数据库迁移解决方案，支持自动迁移表结构、创建索引和添加表注释。

## 目录

- [快速开始](#快速开始)
- [核心概念](#核心概念)
- [配置说明](#配置说明)
- [API 参考](#api-参考)
- [使用示例](#使用示例)
- [最佳实践](#最佳实践)

## 快速开始

### 最简单的使用方式

```go
import (
    "github.com/kamalyes/go-sqlbuilder/db"
    "gorm.io/gorm"
)

// 快速迁移单个或多个模型
err := db.QuickMigrate(gormDB, &User{}, &Order{}, &Product{})
```

### 完整迁移（包含索引和注释）

```go
config := &db.MigratorConfig{
    Models: []interface{}{
        &User{},
        &Order{},
    },
    Indexes: []db.IndexDefinition{
        // 使用便捷函数，索引名自动生成
        db.NewUniqueIndex("users", "email"),            // => idx_users_email_unique
        db.NewIndex("orders", "user_id"),               // => idx_orders_user_id
        db.NewIndex("orders", "user_id", "status"),     // => idx_orders_user_id_status
    },
    Comments: []db.TableComment{
        {Table: "users", Comment: "用户信息表"},
        {Table: "orders", Comment: "订单记录表"},
    },
}

err := db.QuickAutoMigrate(gormDB, config)
```

## 核心概念

### 迁移流程

`Migrator` 执行迁移时按以下顺序进行：

1. **表结构迁移** - 使用 GORM 的 AutoMigrate 创建/更新表结构
2. **索引创建** - 创建自定义索引（支持普通索引和唯一索引）
3. **表注释添加** - 为表添加注释（支持 MySQL 和 PostgreSQL）

### 错误处理策略

- `SkipIndexOnError`: 索引创建失败时是否继续执行（默认 `true`）
- `SkipCommentOnError`: 注释添加失败时是否继续执行（默认 `true`）

这种设计允许迁移过程在遇到非致命错误时继续执行。

## 配置说明

### MigratorConfig

| 字段 | 类型 | 说明 | 默认值 |
|------|------|------|--------|
| `Models` | `[]interface{}` | 需要迁移的 GORM 模型列表 | `nil` |
| `Indexes` | `[]IndexDefinition` | 自定义索引定义列表 | `nil` |
| `Comments` | `[]TableComment` | 表注释定义列表 | `nil` |
| `Logger` | `logger.ILogger` | 日志记录器 | 自动创建 |
| `SkipIndexOnError` | `bool` | 索引失败时跳过 | `true` |
| `SkipCommentOnError` | `bool` | 注释失败时跳过 | `true` |

### IndexDefinition

| 字段 | 类型 | 说明 | 示例 |
|------|------|------|------|
| `Table` | `string` | 表名 | `"users"` |
| `Name` | `string` | 索引名称（**可选，为空时自动生成**） | `"idx_users_email"` |
| `Columns` | `string` | 列定义 | `"(email)"` 或 `"(col1, col2)"` |
| `Unique` | `bool` | 是否唯一索引 | `true` |

#### 索引名自动生成规则

当 `Name` 字段为空时，系统会根据以下规范自动生成索引名：

| 索引类型 | 命名规则 | 示例 |
|---------|---------|------|
| 单列普通索引 | `idx_{表名}_{列名}` | `idx_users_email` |
| 单列唯一索引 | `idx_{表名}_{列名}_unique` | `idx_users_email_unique` |
| 复合索引 | `idx_{表名}_{列1}_{列2}` | `idx_orders_user_id_status` |
| 复合唯一索引 | `idx_{表名}_{列1}_{列2}_unique` | `idx_orders_user_id_order_no_unique` |

### TableComment

| 字段 | 类型 | 说明 | 示例 |
|------|------|------|------|
| `Table` | `string` | 表名 | `"users"` |
| `Comment` | `string` | 注释内容 | `"用户信息表"` |

## API 参考

### 构造函数

#### NewMigrator

```go
func NewMigrator(db *gorm.DB, config *MigratorConfig) *Migrator
```

创建迁移器实例。如果 `config` 为 `nil`，将使用默认配置。

### 索引构建便捷函数 🔥

#### NewIndex

```go
func NewIndex(table string, columns ...string) IndexDefinition
```

创建普通索引定义，自动生成索引名。

```go
// 单列索引：idx_users_email
idx := db.NewIndex("users", "email")

// 复合索引：idx_orders_user_id_status
idx := db.NewIndex("orders", "user_id", "status")
```

#### NewUniqueIndex

```go
func NewUniqueIndex(table string, columns ...string) IndexDefinition
```

创建唯一索引定义，自动生成索引名（带 `_unique` 后缀）。

```go
// 唯一索引：idx_users_email_unique
idx := db.NewUniqueIndex("users", "email")
```

#### NewIndexWithName

```go
func NewIndexWithName(table, name, columns string, unique bool) IndexDefinition
```

创建带自定义名称的索引定义（需要完全控制索引名时使用）。

```go
idx := db.NewIndexWithName("users", "my_custom_idx", "(email, name)", false)
```

### 迁移方法

#### AutoMigrate

```go
func (m *Migrator) AutoMigrate() error
```

执行完整的自动迁移流程：表结构 → 索引 → 注释。

#### MigrateModels

```go
func (m *Migrator) MigrateModels() error
```

仅迁移模型表结构，不创建索引和注释。

#### CreateIndexes

```go
func (m *Migrator) CreateIndexes() error
```

创建配置中定义的所有索引。

#### AddComments

```go
func (m *Migrator) AddComments() error
```

添加配置中定义的所有表注释。

### 工具方法

#### HasTable

```go
func (m *Migrator) HasTable(table string) bool
```

检查单个表是否存在。

#### CheckTablesExist

```go
func (m *Migrator) CheckTablesExist(tables ...string) map[string]bool
```

批量检查多个表是否存在，返回表名到存在状态的映射。

#### DropTables

```go
func (m *Migrator) DropTables(tables ...string) error
```

删除指定的表（⚠️ 危险操作，请谨慎使用）。

### 便捷函数

#### QuickMigrate

```go
func QuickMigrate(db *gorm.DB, models ...interface{}) error
```

快速迁移模型，仅创建表结构，不处理索引和注释。

#### QuickAutoMigrate

```go
func QuickAutoMigrate(db *gorm.DB, config *MigratorConfig) error
```

使用配置执行完整迁移。

## 使用示例

### 示例 1: 基础模型迁移

```go
// 定义模型
type User struct {
    ID        uint      `gorm:"primaryKey"`
    Name      string    `gorm:"column:name;size:100"`
    Email     string    `gorm:"column:email;size:255;unique"`
    CreatedAt time.Time `gorm:"column:created_at"`
}

func (User) TableName() string {
    return "users"
}

// 执行迁移
err := db.QuickMigrate(gormDB, &User{})
if err != nil {
    log.Fatalf("迁移失败: %v", err)
}
```

### 示例 2: 带复合索引的迁移（推荐方式 🔥）

```go
config := &db.MigratorConfig{
    Models: []interface{}{&Order{}},
    Indexes: []db.IndexDefinition{
        // 使用便捷函数，自动生成规范的索引名
        db.NewIndex("orders", "status"),                    // => idx_orders_status
        db.NewIndex("orders", "user_id", "created_at"),     // => idx_orders_user_id_created_at
        db.NewUniqueIndex("orders", "order_no"),            // => idx_orders_order_no_unique
    },
}

migrator := db.NewMigrator(gormDB, config)
err := migrator.AutoMigrate()
```

### 示例 3: 传统方式（手动指定索引名）

```go
config := &db.MigratorConfig{
    Models: []interface{}{&Order{}},
    Indexes: []db.IndexDefinition{
        // 手动指定索引名
        {
            Table:   "orders",
            Name:    "idx_orders_status",
            Columns: "(status)",
        },
        // 复合索引
        {
            Table:   "orders",
            Name:    "idx_orders_user_created",
            Columns: "(user_id, created_at DESC)",
        },
        // 唯一索引
        {
            Table:   "orders",
            Name:    "idx_orders_order_no",
            Columns: "(order_no)",
            Unique:  true,
        },
    },
}

migrator := db.NewMigrator(gormDB, config)
err := migrator.AutoMigrate()
```

### 示例 3: 自定义日志记录器

```go
import "github.com/kamalyes/go-logger"

customLogger := logger.NewLogger(&logger.Config{
    Level: "debug",
})

config := &db.MigratorConfig{
    Models: []interface{}{&User{}},
    Logger: customLogger,
}

migrator := db.NewMigrator(gormDB, config)
err := migrator.AutoMigrate()
```

### 示例 4: 严格模式（失败即停止）

```go
config := &db.MigratorConfig{
    Models: []interface{}{&User{}, &Order{}},
    Indexes: []db.IndexDefinition{
        {Table: "users", Name: "idx_email", Columns: "(email)"},
    },
    // 设置为 false，任何错误都会停止迁移
    SkipIndexOnError:   false,
    SkipCommentOnError: false,
}

migrator := db.NewMigrator(gormDB, config)
if err := migrator.AutoMigrate(); err != nil {
    log.Fatalf("迁移失败: %v", err)
}
```

### 示例 5: 检查和删除表

```go
migrator := db.NewMigrator(gormDB, nil)

// 检查单个表
if migrator.HasTable("users") {
    fmt.Println("users 表存在")
}

// 批量检查
result := migrator.CheckTablesExist("users", "orders", "products")
for table, exists := range result {
    fmt.Printf("%s: %v\n", table, exists)
}

// 删除表（谨慎使用！）
err := migrator.DropTables("temp_table", "old_backup")
```

### 示例 6: 分步迁移

```go
config := &db.MigratorConfig{
    Models: []interface{}{&User{}, &Order{}},
    Indexes: []db.IndexDefinition{
        {Table: "users", Name: "idx_email", Columns: "(email)"},
    },
    Comments: []db.TableComment{
        {Table: "users", Comment: "用户表"},
    },
}

migrator := db.NewMigrator(gormDB, config)

// 分步执行，便于调试
fmt.Println("步骤 1: 迁移表结构")
if err := migrator.MigrateModels(); err != nil {
    log.Fatalf("表迁移失败: %v", err)
}

fmt.Println("步骤 2: 创建索引")
if err := migrator.CreateIndexes(); err != nil {
    log.Printf("索引创建有错误: %v", err)
}

fmt.Println("步骤 3: 添加注释")
if err := migrator.AddComments(); err != nil {
    log.Printf("注释添加有错误: %v", err)
}
```

## 最佳实践

### 1. 使用便捷函数创建索引（推荐）

使用 `NewIndex` 和 `NewUniqueIndex` 可以自动生成规范的索引名，避免手动命名错误：

```go
// ✅ 推荐：使用便捷函数
indexes := []db.IndexDefinition{
    db.NewIndex("users", "email"),                  // idx_users_email
    db.NewIndex("orders", "user_id", "status"),     // idx_orders_user_id_status
    db.NewUniqueIndex("orders", "order_no"),        // idx_orders_order_no_unique
}

// ❌ 不推荐：手动命名容易出错
indexes := []db.IndexDefinition{
    {Table: "users", Name: "idx_user_email", Columns: "(email)"},      // 表名用了单数
    {Table: "orders", Name: "index_orders_status", Columns: "(status)"}, // 前缀不统一
}
```

### 2. 索引命名规范（自动生成）

系统自动采用统一的索引命名规范：

| 类型 | 格式 | 示例 |
|------|------|------|
| 普通索引 | `idx_{表名}_{列名}` | `idx_users_email` |
| 唯一索引 | `idx_{表名}_{列名}_unique` | `idx_users_email_unique` |
| 复合索引 | `idx_{表名}_{列1}_{列2}` | `idx_orders_user_id_status` |

### 2. 生产环境建议

```go
// 生产环境推荐配置
config := &db.MigratorConfig{
    Models: models,
    Indexes: indexes,
    Comments: comments,
    // 生产环境建议开启跳过，避免因索引问题导致部署失败
    SkipIndexOnError:   true,
    SkipCommentOnError: true,
}
```

### 3. 迁移前检查

```go
migrator := db.NewMigrator(gormDB, config)

// 检查关键表是否已存在
criticalTables := []string{"users", "orders"}
result := migrator.CheckTablesExist(criticalTables...)

for _, table := range criticalTables {
    if result[table] {
        log.Printf("⚠️ 表 %s 已存在，将执行增量迁移", table)
    }
}

// 执行迁移
err := migrator.AutoMigrate()
```

### 4. 数据库兼容性

| 特性 | MySQL | PostgreSQL | SQLite | SQL Server |
|------|-------|------------|--------|------------|
| 表迁移 | ✅ | ✅ | ✅ | ✅ |
| 索引创建 | ✅ | ✅ | ✅ | ✅ |
| 表注释 | ✅ | ✅ | ❌ | ❌ |

> 💡 **提示**: 表注释在 SQLite 和 SQL Server 中会被静默跳过，不会报错。

## 📚 相关文档

- 📘 [快速开始](./QUICKSTART.md) - 5 分钟上手指南
- 📖 [仓储基础](./REPOSITORY-BASICS.md) - 学习基础 CRUD 操作
- 🏗️ [模型定义](./MODELS.md) - 定义数据模型
- 📒 [错误处理](./ERROR-HANDLING.md) - 错误管理和日志记录
