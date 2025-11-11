# Architecture - Go-SQLBuilder 架构设计文档

## 系统架构图

```bash
┌─────────────────────────────────────────────────────────────┐
│                     应用层                                   │
│  Web应用、API服务、命令行工具等                             │
└────────────────────────┬────────────────────────────────────┘
                         │
┌────────────────────────┴────────────────────────────────────┐
│                  SQLBuilder 核心层                           │
│  ┌──────────────────────────────────────────────────────┐   │
│  │        Builder (SQL构建和执行)                        │   │
│  │  ├── Select/Insert/Update/Delete                    │   │
│  │  ├── Join/GroupBy/Having/OrderBy                    │   │
│  │  └── Transaction Support                             │   │
│  └──────────────────────┬───────────────────────────────┘   │
│                         │                                     │
│  ┌──────────────────────┴───────────────────────────────┐   │
│  │      CachedBuilder (性能优化层)                       │   │
│  │  ├── 自动缓存键生成 (MD5)                              │   │
│  │  ├── TTL自动失效                                      │   │
│  │  └── 缓存管理                                         │   │
│  └──────────────────────────────────────────────────────┘   │
└────────────────┬─────────────────────────┬──────────────────┘
                 │                         │
     ┌───────────┴──────────┐   ┌──────────┴──────────┐
     │                      │   │                     │
┌────┴────────────────┐ ┌──┴──────────────┐  ┌──────┴──────────┐
│   Framework Layer   │ │  Service Layer  │  │  Cache Layer    │
│  ┌──────────────┐   │ │ ┌─────────────┐ │  │ ┌─────────────┐ │
│  │ SQLx Adapter │   │ │ │Query Params │ │  │ │Store Inter  │ │
│  │ GORM Adapter │   │ │ │Error Handle │ │  │ │MockStore    │ │
│  │PostgreSQL    │   │ │ │Filter Build │ │  │ │RedisStore   │ │
│  │MySQL         │   │ │ │Pagination   │ │  │ │Manager      │ │
│  │SQLite        │   │ │ └─────────────┘ │  │ └─────────────┘ │
│  └──────────────┘   │ └─────────────────┘  └─────────────────┘
└─────────────────────┘
```

## 分层设计

### 1. 应用层（Application Layer）

- 业务逻辑实现
- API路由和控制器
- 数据验证和转换

### 2. 核心层（Core Layer）

**Builder** - SQL构建引擎

- 执行SQL生成和执行的主要职责
- 支持链式调用
- 管理SQL各个组件（SELECT、WHERE等）

**CachedBuilder** - 缓存优化层

- 包装Builder添加缓存能力
- 自动生成缓存键
- 管理TTL和失效

### 3. 服务层（Service Layer）

**query包** - 高级查询参数

- 参数构建和组合
- WHERE子句生成
- 20+便捷方法

**errors包** - 错误处理

- 统一错误定义（48种）
- 错误转换和格式化
- 错误码映射

### 4. 基础设施层（Infrastructure Layer）

**Framework Layer** - 框架适配

- SQLX适配器
- GORM适配器
- 多数据库支持（MySQL、PostgreSQL、SQLite）

**Cache Layer** - 缓存管理

- Store接口定义
- MockStore（测试）
- RedisStore（生产）
- CacheManager（统计管理）

## 核心设计模式

### 1. 适配器模式（Adapter Pattern）

```bash
┌─────────────┐     ┌─────────────┐
│   SQLX DB   │     │   GORM DB   │
└─────┬───────┘     └─────┬───────┘
      │                   │
      └─────────┬─────────┘
                │
         ┌──────┴──────┐
         │   Adapter   │
         │  Interface  │
         └──────┬──────┘
                │
            ┌───┴──────────┐
            │              │
      ┌─────▼──────┐  ┌───▼──────┐
      │SqlxAdapter │  │GormAdapter│
      └────────────┘  └───────────┘
```

**目的**: 统一多个ORM框架的接口，使Builder可以支持任何数据库框架

### 2. 构建器模式（Builder Pattern）

```go
// 链式调用的流畅设计
builder.Table("users").
    Select("id", "name").
    Where("age", ">", 18).
    OrderBy("name", "ASC").
    Limit(10).
    Find(&users)
```

**特点**:

- 每个方法返回自身（*Builder）
- 支持任意组合和顺序
- 可读性强，表达意图明确

### 3. 工厂模式（Factory Pattern）

```go
// 工厂函数集合
builder := sqlbuilder.New(db)
cachedBuilder, _ := sqlbuilder.NewCachedBuilder(db, store, config)
param := query.NewParam()
store := cache.NewMockStore()
store := cache.NewRedisStore(client, prefix)
err := errors.NewError(code, message)
```

**优势**:

- 统一的对象创建方式
- 参数验证
- 默认配置

### 4. 策略模式（Strategy Pattern）

```bash
┌──────────────┐
│ CachedBuilder│
└────────┬─────┘
         │ uses
    ┌────┴─────┐
    │  Store   │ (Strategy)
    └────┬─────┘
         │ can be
    ┌────┴─────────────┐
    │                  │
┌───┴────┐      ┌──────┴──┐
│MockStore│      │RedisStore│
└────────┘      └──────────┘
```

**应用**: 缓存实现可以切换，不影响上层代码

### 5. 模板方法模式（Template Method Pattern）

```go
// CachedBuilder.GetCached 的流程
1. 生成缓存键（MD5）
2. 尝试从缓存获取
3. 如果miss，执行SQL查询
4. 结果存入缓存
5. 返回结果
```

## 数据流转

### 查询流程

```bash
用户代码
   │
   ▼
┌──────────────┐
│Builder.Find()│
└──────┬───────┘
       │
       ▼
┌─────────────────────┐
│生成SQL:             │
│SELECT * FROM users..│
└──────┬──────────────┘
       │
       ▼
┌──────────────────────┐
│参数化（占位符？）     │
│SQL注入防护           │
└──────┬───────────────┘
       │
       ▼
┌──────────────────┐
│Adapter执行       │
│通过ORM框架       │
└──────┬───────────┘
       │
       ▼
┌──────────────┐
│返回结果集     │
└──────┬───────┘
       │
       ▼
┌──────────────┐
│映射到结构体   │
└──────────────┘
```

### 缓存查询流程

```bash
用户代码
   │
   ▼
┌──────────────────┐
│CachedBuilder.Get │
└──────┬───────────┘
       │
       ▼
┌──────────────────────┐
│生成缓存键 (MD5)      │
└──────┬───────────────┘
       │
       ▼
┌──────────────────────┐
│查询缓存Store         │
└──────┬───────────────┘
       │
   ┌───┴──────┐
   │           │
   ▼ Hit       ▼ Miss
┌────────┐  ┌──────────────┐
│返回值   │  │执行SQL查询   │
└────────┘  └──────┬───────┘
               │
               ▼
            ┌─────────────┐
            │存入缓存     │
            │(TTL失效)   │
            └──────┬──────┘
                   │
                   ▼
            ┌────────────┐
            │返回结果    │
            └────────────┘
```

## 错误处理流程

```bash
SQL执行 ─────┐
Adapter操作  ├──→ Error发生
超时/失败    ┘
             │
             ▼
        ┌──────────────┐
        │错误对象返回  │
        │(error type) │
        └──────┬───────┘
               │
    ┌──────────┴──────────┐
    │                     │
    ▼                     ▼
┌─────────────────┐  ┌──────────────┐
│标准error        │  │AppError      │
│（未分类）       │  │（分类编码）  │
└─────────────────┘  └─────┬────────┘
                           │
                    ┌──────┴──────┐
                    │             │
              ┌─────▼───┐   ┌────▼────┐
              │.Error() │   │.String()│
              │接口     │   │接口     │
              └─────────┘   └─────────┘
```

## 包依赖关系

```bash
sqlbuilder (root)
    │
    ├── adapter (通过interfaces)
    │   ├── sqlx 库
    │   └── gorm 库
    │
    ├── cache 包
    │   ├── interface (Store, RedisClientInterface)
    │   ├── config
    │   ├── manager
    │   ├── mock (for testing)
    │   └── redis
    │
    ├── query 包
    │   ├── operator (常量定义)
    │   ├── filter
    │   ├── pagination
    │   ├── option
    │   └── param
    │
    ├── errors 包
    │   ├── code (48种错误码)
    │   └── error (AppError实现)
    │
    └── advanced_query (向后兼容)
        └── (转发到query包)
```

**依赖特点**:

- ✅ 低耦合 - 包之间通过接口通信
- ✅ 高内聚 - 相关功能聚集在一起
- ✅ 易扩展 - 新增功能无需修改现有代码

## 并发安全性

### Builder

- ✅ 无全局状态，实例安全
- ✅ 推荐每个goroutine创建独立Builder实例
- ✅ 底层数据库连接池自动管理

### CachedBuilder

- ✅ Store接口实现线程安全
- ✅ RedisStore通过Redis原子操作保证
- ✅ MockStore使用内存安全的数据结构

### Cache Package

- ✅ RedisStore：通过redis库的并发安全
- ✅ MockStore：通过Go的map互斥锁保护
- ✅ Manager：原子操作记录统计

## 性能考量

### 缓存层优化

```bash
查询频率：100次/秒
缓存命中率：80%
└─ 数据库查询：20次/秒（减少80%）
└─ 缓存查询：80次/秒（<1ms延迟）
```

### 参数化查询优化

```bash
数据库预编译
└─ 首次解析SQL
└─ 后续直接执行
└─ 减少解析时间
```

### 连接池优化

```bash
连接复用
└─ SQLX/GORM内置连接池
└─ 减少连接开销
└─ 提高并发能力
```

## 扩展点

### 1. 新增数据库支持

```go
// 实现UniversalAdapterInterface
type CustomDBAdapter struct { ... }
func (a *CustomDBAdapter) Query() (sql.Rows, error) { ... }
// 与Builder兼容
```

### 2. 新增缓存实现

```go
// 实现cache.Store接口
type MyCustomStore struct { ... }
func (s *MyCustomStore) Get(ctx, key) (string, error) { ... }
// 可即插即用
```

### 3. 新增查询操作符

```go
// 在query/operator.go中添加
const OP_CUSTOM Operator = "CUSTOM_OP"
// 在query/param.go中添加便捷方法
func (p *Param) AddCustom(...) *Param { ... }
```

## 安全特性

### SQL注入防护

```bash
所有WHERE条件 ──→ 参数化查询 ──→ ?占位符 ──→ 参数绑定
直接用户输入 ──→ 防止注入 ──→ 数据库执行 ──→ 100%安全
```

### 错误信息安全

```bash
内部错误 ──→ ErrorCode ──→ 日志记录 ──→ 不暴露细节
用户端 ──→ 统一错误码 ──→ 安全响应
```
