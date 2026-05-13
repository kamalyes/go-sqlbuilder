# 多租户作用域查询系统

## 📖 概述

Scope 模块提供了多租户场景下的数据访问控制能力，支持 OPS 域和租户域两种访问模式，以及全局、地区级、平台级、租户级四种作用域。

### 核心能力

- ✅ **OPS 域**：OPS 全局管理员、OPS 租户管理员
- ✅ **租户域**：租户全局用户、地区级用户、平台级用户
- ✅ **自动字段映射**：默认 `tenant_id`、`platform_id`、`region_code`
- ✅ **自定义配置**：通过 WithXxx 选项覆盖默认值
- ✅ **性能优化**：自动选择 EQ/IN 操作符
- ✅ **旧条件清理**：自动移除旧作用域字段避免重复

***

## 🚀 快速开始

### 1. 基本使用

```go
package main

import (
    "context"
    
    "github.com/kamalyes/go-sqlbuilder/repository"
    "github.com/kamalyes/go-sqlbuilder/scope"
)

func main() {
    // 创建作用域数据（使用默认字段映射）
    data := scope.NewScopeData()
    data.Domain = 1  // 租户域
    data.TenantID = "T001"
    data.ScopeEntries = []*scope.ScopeEntry{
        {ScopeType: 2, RegionCodes: []string{"MM", "TH"}},  // 地区级
    }
    
    // 应用作用域到查询
    query := repository.NewQuery()
    result := scope.ApplySQLScope(query, data)
    
    // 使用 result.FilterGroup 进行查询
    // ...
}
```

### 2. 与 Repository 集成

```go
// 创建 Repository
repo := repository.NewBaseRepository[User](handler, logger, "users")

// 构建作用域数据
data := scope.NewScopeData()
data.Domain = 1
data.TenantID = "T001"
data.ScopeEntries = []*scope.ScopeEntry{
    {ScopeType: 3, RegionPlatforms: []*scope.RegionPlatformEntry{
        {RegionCode: "MM", PlatformIds: []string{"P1", "P2"}},
    }},
}

// 应用作用域并查询
query := repository.NewQuery()
scopedQuery := scope.ApplySQLScope(query, data)
users, err := repo.List(ctx, scopedQuery)
```

***

## 🔗 与 go-scope-provider 集成

[go-scope-provider](https://github.com/kamalyes/go-scope-provider) 是独立的作用域提供者包，提供从 gRPC 上下文解析作用域数据的能力，不依赖任何业务 protobuf 定义。

### 1. 安装依赖

```bash
go get github.com/kamalyes/go-scope-provider
```

### 2. 从 Context 提取作用域

```go
package main

import (
    "context"
    
    "github.com/kamalyes/go-scope-provider/provider"
    "github.com/kamalyes/go-sqlbuilder/repository"
    "github.com/kamalyes/go-sqlbuilder/scope"
)

func QueryWithScope(ctx context.Context) {
    // 方式一：一步到位，从上下文解析并应用到查询
    query := repository.NewQuery()
    result := provider.ApplyScopeQuery(ctx, query)
    
    // 方式二：先解析再应用
    data := provider.ResolveScopeData(ctx)
    result = scope.ApplySQLScope(query, data)
    
    // ...
}
```

### 3. 注入 Payload 到上下文

```go
// 方式一：直接注入 Payload 对象
payload := &provider.Payload{
    Domain:   1,
    TenantID: "T001",
    RoleCode: "owner",
    ScopeBindings: []*provider.ScopeEntry{
        {ScopeType: 2, RegionCodes: []string{"MM", "TH"}},
    },
}
ctx := provider.WithPayload(context.Background(), payload)

// 方式二：通过上下文键注入基础信息
ctx := provider.WithScopeContext(context.Background(), 1, "T001", "owner")
```

### 4. gRPC 拦截器

在 gRPC 服务端注册拦截器，自动从 metadata 中的 `x-auth-payload` 解析载荷：

```go
import (
    "github.com/kamalyes/go-scope-provider/provider"
    "google.golang.org/grpc"
)

func main() {
    server := grpc.NewServer(
        grpc.UnaryInterceptor(provider.UnaryPayloadInterceptor()),
        grpc.StreamInterceptor(provider.StreamPayloadInterceptor()),
    )
    // ...
}
```

### 5. Payload JSON 格式

载荷以 base64 编码的 JSON 存放在 gRPC metadata 中，格式如下：

```json
{
    "domain": 1,
    "tenant_id": "T001",
    "user_id": "U001",
    "role_code": "owner",
    "scope_bindings": [
        {
            "scope_type": 2,
            "region_codes": ["MM", "TH"]
        },
        {
            "scope_type": 3,
            "region_platforms": [
                {
                    "region_code": "SG",
                    "platform_ids": ["P1", "P2"]
                }
            ]
        }
    ]
}
```

### 6. 自定义配置

```go
// 传递自定义 scope 选项
data := provider.ResolveScopeData(ctx,
    scope.WithTenantIDField("org_id"),
    scope.WithRegionCodeField("rc"),
)

// 或在 ApplyScopeQuery 中使用
result := provider.ApplyScopeQuery(ctx, query,
    scope.WithFieldMapping(scope.FieldMapping{
        TenantIDField:   "org_id",
        PlatformIDField: "pid",
        RegionCodeField: "rc",
    }),
)
```

***

## 📊 作用域层级

| 作用域类型 | ScopeType 默认值 | 适用域          | 访问范围      |
| ----- | ------------- | ------------ | --------- |
| 全局    | 1             | OPS / Tenant | 所有数据      |
| 地区级   | 2             | Tenant       | 指定地区所有数据  |
| 平台级   | 3             | Tenant       | 指定地区+平台数据 |
| 租户级   | 4             | OPS          | 指定租户数据    |

### 域类型

| 域     | Domain 默认值 | 说明          |
| ----- | ---------- | ----------- |
| 租户域   | 1          | 租户内部用户访问    |
| OPS 域 | 2          | OPS 平台管理员访问 |

***

## 🔧 配置选项

### 默认配置

```go
// DefaultConfig 返回的默认值
Config{
    Domain: DomainConfig{
        TenantDomainValue: 1,  // 租户域
        OpsDomainValue:    2,  // OPS 域
    },
    ScopeType: ScopeTypeConfig{
        GlobalValue:   1,  // 全局作用域
        RegionValue:   2,  // 地区级
        PlatformValue: 3,  // 平台级
        TenantValue:   4,  // 租户级
    },
    Mapping: FieldMapping{
        TenantIDField:   "tenant_id",    // 租户ID字段
        PlatformIDField: "platform_id",  // 平台ID字段
        RegionCodeField: "region_code",  // 地区编码字段
    },
}
```

### 自定义配置

```go
// 方式一：自定义字段名
data := scope.NewScopeData(
    scope.WithTenantIDField("org_id"),
    scope.WithRegionCodeField("rc"),
    scope.WithPlatformIDField("pid"),
)

// 方式二：自定义域值
data := scope.NewScopeData(
    scope.WithDomainConfig(10, 20),  // Tenant=10, OPS=20
)

// 方式三：自定义作用域类型值
data := scope.NewScopeData(
    scope.WithScopeTypeConfig(100, 200, 300, 400),  // Global=100, Region=200, Platform=300, Tenant=400
)

// 方式四：一次性设置所有字段映射
data := scope.NewScopeData(
    scope.WithFieldMapping(scope.FieldMapping{
        TenantIDField:   "org_id",
        PlatformIDField: "pid",
        RegionCodeField: "rc",
    }),
)
```

***

## 📝 使用示例

### OPS 全局管理员

```go
// 场景：OPS 全局管理员，可访问所有数据
data := scope.NewScopeData()
data.Domain = 2  // OPS 域
data.ScopeEntries = []*scope.ScopeEntry{
    {ScopeType: 1},  // 全局作用域
}

query := repository.NewQuery()
result := scope.ApplySQLScope(query, data)
// result.FilterGroup == nil（不加任何过滤）
```

### OPS 租户管理员

```go
// 场景：OPS 管理指定租户 T1、T2
data := scope.NewScopeData()
data.Domain = 2  // OPS 域
data.ScopeEntries = []*scope.ScopeEntry{
    {ScopeType: 4, TenantIds: []string{"T1", "T2"}},  // 租户级
}

query := repository.NewQuery()
result := scope.ApplySQLScope(query, data)
// 生成：AND { tenant_id IN ('T1', 'T2') }
```

### 租户全局用户

```go
// 场景：租户 Owner，可访问租户下所有数据
data := scope.NewScopeData()
data.Domain = 1  // 租户域
data.TenantID = "T001"
data.ScopeEntries = []*scope.ScopeEntry{
    {ScopeType: 1},  // 全局作用域
}

query := repository.NewQuery()
result := scope.ApplySQLScope(query, data)
// 生成：AND { tenant_id = 'T001' }
```

### 租户地区级用户

```go
// 场景：用户可访问 MM、TH 两个地区
data := scope.NewScopeData()
data.Domain = 1
data.TenantID = "T001"
data.ScopeEntries = []*scope.ScopeEntry{
    {ScopeType: 2, RegionCodes: []string{"MM", "TH"}},
}

query := repository.NewQuery()
result := scope.ApplySQLScope(query, data)
// 生成：AND { tenant_id = 'T001', OR { AND { region_code IN ('MM', 'TH') } } }
```

### 租户平台级用户

```go
// 场景：用户可访问 MM 地区的 P1、P2 平台，SG 地区的 P1 平台
data := scope.NewScopeData()
data.Domain = 1
data.TenantID = "T001"
data.ScopeEntries = []*scope.ScopeEntry{
    {ScopeType: 3, RegionPlatforms: []*scope.RegionPlatformEntry{
        {RegionCode: "MM", PlatformIds: []string{"P1", "P2"}},
        {RegionCode: "SG", PlatformIds: []string{"P1"}},
    }},
}

query := repository.NewQuery()
result := scope.ApplySQLScope(query, data)
// 生成：AND {
//   tenant_id = 'T001',
//   OR {
//     AND { region_code = 'MM', platform_id IN ('P1', 'P2') },
//     AND { region_code = 'SG', platform_id = 'P1' }
//   }
// }
```

### 混合作用域

```go
// 场景：用户同时拥有地区级和平台级权限
data := scope.NewScopeData()
data.Domain = 1
data.TenantID = "T001"
data.ScopeEntries = []*scope.ScopeEntry{
    {ScopeType: 2, RegionCodes: []string{"MM"}},  // MM 地区所有平台
    {ScopeType: 3, RegionPlatforms: []*scope.RegionPlatformEntry{
        {RegionCode: "SG", PlatformIds: []string{"P1"}},  // SG 地区 P1 平台
    }},
}

query := repository.NewQuery()
result := scope.ApplySQLScope(query, data)
// 两种条目合并到同一个 OR 组
```

***

## ⚡ 性能优化

### 自动选择 EQ/IN 操作符

```go
// 单值自动使用 EQ
data.ScopeEntries = []*scope.ScopeEntry{
    {ScopeType: 2, RegionCodes: []string{"MM"}},
}
// 生成：region_code = 'MM'

// 多值自动使用 IN
data.ScopeEntries = []*scope.ScopeEntry{
    {ScopeType: 2, RegionCodes: []string{"MM", "TH", "SG"}},
}
// 生成：region_code IN ('MM', 'TH', 'SG')
```

### 单地区优化路径

当平台级条目只有一个地区且同时配置了 `RegionCodeField` 和 `PlatformIDField` 时，自动走优化路径：

```go
// 单地区 + 双字段映射 → 合并到同一个 AND 组
data.ScopeEntries = []*scope.ScopeEntry{
    {ScopeType: 3, RegionPlatforms: []*scope.RegionPlatformEntry{
        {RegionCode: "MM", PlatformIds: []string{"P1", "P2"}},
    }},
}
// 生成：AND { region_code = 'MM', platform_id IN ('P1', 'P2') }
```

***

## 🔄 旧条件清理

当查询中已存在作用域字段过滤时，会自动移除旧条件：

```go
query := repository.NewQuery()
query.Filters = []*repository.Filter{
    {Field: "tenant_id", Operator: "EQ", Value: "old"},
    {Field: "region_code", Operator: "EQ", Value: "old"},
    {Field: "name", Operator: "EQ", Value: "keep"},  // 非作用域字段
}

result := scope.ApplySQLScope(query, data)
// result.Filters 仅保留 name 字段
// result.FilterGroup 包含新的作用域条件
```

***

## 📐 数据结构

### ScopeEntry

```go
type ScopeEntry struct {
    ScopeType       int32                   // 作用域类型
    RegionCodes     []string                // 地区编码列表（地区级）
    RegionPlatforms []*RegionPlatformEntry  // 地区-平台绑定（平台级）
    TenantIds       []string                // 租户ID列表（OPS 租户级）
}
```

### RegionPlatformEntry

```go
type RegionPlatformEntry struct {
    RegionCode  string    // 地区编码
    PlatformIds []string  // 平台ID列表（空=该地区所有平台）
}
```

### ScopeData

```go
type ScopeData struct {
    Domain       int32          // 当前域值
    TenantID     string         // 当前租户ID
    IsOwner      bool           // 是否为租户 Owner
    ScopeEntries []*ScopeEntry  // 作用域条目列表
    Config       Config         // 作用域配置
}
```

***

## 🧪 测试覆盖

所有测试文件均包含中文注释，说明场景和预期结果：

| 测试文件              | 覆盖内容             |
| ----------------- | ---------------- |
| `config_test.go`  | 配置、默认值、Option 函数 |
| `types_test.go`   | 数据类型、判断方法        |
| `adapter_test.go` | SQL 适配器、各种作用域场景  |

运行测试：

```bash
cd go-sqlbuilder/scope
go test -coverprofile cover.out
go tool cover -html=cover.out
```

***

## 📚 相关文档

- [SCOPE.md](../scope/SCOPE.md) - 作用域设计文档（含 Mermaid 流程图）
- [FILTER-GROUP.md](./FILTER-GROUP.md) - 条件组合高级用法
- [QUERY.md](./QUERY.md) - Query 对象使用

***

## ❓ 常见问题

### Q: 如何清空默认字段映射？

```go
data := scope.NewScopeData(
    scope.WithFieldMapping(scope.FieldMapping{}),  // 全部清空
)
```

### Q: 如何只清空某个字段？

```go
data := scope.NewScopeData(
    scope.WithRegionCodeField(""),  // 只清空地区字段
)
```

### Q: 平台级条目的 PlatformIds 为空是什么意思？

表示该地区所有平台，生成的条件仅包含 `region_code = 'X'`。

### Q: 如何判断当前用户是否有全局权限？

```go
if data.HasGlobalScope() {
    // 全局权限，不加过滤
}
```

### Q: 如何获取用户可访问的所有地区/平台？

```go
regionCodes := data.AllRegionCodes()   // 所有地区编码
platformIds := data.AllPlatformIds()   // 所有平台ID
```

