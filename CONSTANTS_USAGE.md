# Go-SQLBuilder 常量使用规范

## 项目概述

本项目已按功能模块将常量组织到 `constant/` 包中，避免硬编码字符串，提高代码可维护性和类型安全性。

## 常量文件组织结构

```bash
constant/
├── adapter.go      # 数据库适配器相关常量
├── config.go       # 配置相关常量  
├── error.go        # 错误码和错误信息常量
├── field.go        # 数据库字段名常量 (新增)
├── hook.go         # 钩子事件常量 (新增)
├── logger.go       # 日志相关常量
├── message.go      # 消息相关常量
├── middleware.go   # 中间件相关常量
├── operation.go    # 数据库操作类型常量 (新增)
├── operator.go     # SQL操作符常量
├── query.go        # 查询相关常量
├── sort.go         # 排序相关常量
├── sql.go          # SQL语句相关常量
└── validation.go   # 验证相关常量
```

## 推荐使用方式

### 1. 数据库操作类型 (operation.go)

```go
// ✅ 推荐写法
eb.addAuditFields(data, constant.OperationTypeCreate)
eb.addAuditFields(data, constant.OperationTypeUpdate)
eb.addAuditFields(data, constant.OperationTypeDelete)

// ❌ 不推荐写法
eb.addAuditFields(data, "create")
eb.addAuditFields(data, "update")
```

### 2. 钩子事件 (hook.go)

```go
// ✅ 推荐写法
builder.AddHook(constant.HookEventBeforeCreate, hookFunc)
builder.AddHook(constant.HookEventAfterUpdate, hookFunc)

// ❌ 不推荐写法
builder.AddHook("beforeCreate", hookFunc)
builder.AddHook("afterUpdate", hookFunc)
```

### 3. 数据库字段名 (field.go)

```go
// ✅ 推荐写法
builder.AddAuditFields(
    constant.FieldCreatedAt,
    constant.FieldUpdatedAt,
    constant.FieldDeletedAt,
)

// 查询中使用字段常量
query.WhereNull(constant.FieldDeletedAt)

// ❌ 不推荐写法
builder.AddAuditFields("created_at", "updated_at", "deleted_at")
query.WhereNull("deleted_at")
```

### 4. SQL操作符 (operator.go)

```go
// ✅ 推荐写法
filter := &EnhancedFilter{
    Field:    "age",
    Operator: string(constant.OP_GT),
    Value:    18,
}

filter2 := &EnhancedFilter{
    Field:    "email",
    Operator: string(constant.OP_LIKE),
    Value:    "%@example.com",
}

// ❌ 不推荐写法
filter := &EnhancedFilter{
    Field:    "age", 
    Operator: ">",
    Value:    18,
}
```

### 5. 排序方向 (sort.go)

```go
// ✅ 推荐写法
order := &OrderOption{
    Field:     constant.FieldCreatedAt,
    Direction: constant.OrderDESC,
}

// ❌ 不推荐写法
order := &OrderOption{
    Field:     "created_at",
    Direction: "DESC",
}
```

### 6. 配置默认值 (config.go)

```go
// ✅ 推荐写法
findOptions := &FindOptions{
    Limit: constant.DefaultPageSize,
}

// ❌ 不推荐写法
findOptions := &FindOptions{
    Limit: 10,
}
```

## 常量使用的好处

### 1. 类型安全

```go
// 编译时检查，避免运行时字符串错误
eb.addAuditFields(data, constant.OperationTypeCreate) // ✅ 
eb.addAuditFields(data, "creat")                      // ❌ 拼写错误
```

### 2. IDE智能提示

```go
// IDE可以自动补全 constant. 后的所有选项
constant.OP_         // 显示所有操作符常量
constant.Field       // 显示所有字段名常量
constant.HookEvent   // 显示所有钩子事件常量
```

### 3. 重构安全性

```go
// 修改常量定义后，所有使用位置自动更新
// 如果有遗漏，编译器会报错
```

### 4. 代码可读性

```go
// 常量名称自文档化，清晰表达意图
if operation == constant.OperationTypeCreate {  // ✅ 清晰
if operation == "create" {                      // ❌ 魔法字符串
```

## 导入和使用

```go
import (
    "github.com/kamalyes/go-sqlbuilder/constant"
)

// 在代码中使用常量
func example() {
    // 操作类型
    operationType := constant.OperationTypeCreate
    
    // 字段名
    fieldName := constant.FieldCreatedAt
    
    // 操作符  
    operator := constant.OP_EQ
    
    // 排序方向
    sortDir := constant.OrderDESC
    
    // 钩子事件
    hookEvent := constant.HookEventBeforeCreate
}
```

## 最佳实践总结

### ✅ 推荐做法

1. **始终使用常量**：优先使用 `constant` 包中的常量
2. **分类组织**：将相关常量放在同一文件中
3. **命名规范**：使用清晰、一致的命名模式
4. **文档说明**：为常量添加适当的注释
5. **类型定义**：为操作符等使用自定义类型增强类型安全

### ❌ 避免做法

1. **硬编码字符串**：避免在代码中直接使用字符串字面量
2. **魔法数字**：避免使用未命名的数字常量
3. **重复定义**：避免在多个地方定义相同的常量
4. **不当分类**：避免将无关常量混在一起
5. **缺少文档**：避免创建无注释的常量

## 代码示例

完整的使用示例请参考：

- `examples/enhanced_crud_with_constants.go` - 增强CRUD操作示例
- `enhanced_crud.go` - 主要实现文件

这些示例展示了如何在实际项目中正确使用常量，提高代码质量和维护性。
