# 🚀 简单CRUD - 基于EnhancedBuilder的最简单接口

> **您说得对！既然有了强大的EnhancedBuilder，为什么还要重复造轮子？**  
> 现在直接在EnhancedBuilder基础上提供简单易用的接口！

## 🎯 设计理念

**在强大的基础上做简化！** 基于EnhancedBuilder提供最简单的接口：

- ✅ **零配置使用** - 自动设置软删除、审计字段、时间戳hooks
- ✅ **自动时间戳** - 使用hook机制自动处理created_at、updated_at
- ✅ **中文错误信息** - 友好易懂的错误提示
- ✅ **保留高级功能** - 随时可以使用EnhancedBuilder的强大特性

## 🚀 5秒上手

```go
// 1. 创建简单操作器（基于EnhancedBuilder）
crud, err := sqlbuilder.NewSimple(db, "users")

// 2. 开始CRUD - 就这么简单！
err = crud.Add(map[string]interface{}{"name": "张三", "age": 25})
```

**就这样！你获得了简单接口 + 强大功能！**

## 📖 完整API

### 基础操作（99%的场景够用）

```go
// 添加数据（自动时间戳）
err := crud.Add(map[string]interface{}{
    "name": "张三",
    "email": "zhangsan@email.com",
    "age": 25,
})

// 按ID获取
user, err := crud.Get(1)

// 按ID更新（自动更新时间）
err := crud.Update(1, map[string]interface{}{
    "age": 26,
})

// 软删除（自动设置deleted_at）
err := crud.Delete(1)

// 分页列表（自动过滤软删除）
users, err := crud.List(1, 10) 

// 统计数量（排除软删除）
count, err := crud.Count()

// 搜索（模糊查询）
results, err := crud.Search("name", "张", 1, 10)

// 智能保存（有ID更新，无ID新增）
err := crud.Save(data)
```

### 高级功能随时可用

```go
// 创建简单操作器
crud, err := sqlbuilder.NewSimple(db, "users")

// 需要高级功能时，直接使用EnhancedBuilder方法
ctx := context.Background()
options := &sqlbuilder.CreateOptions{...}
result, err := crud.SmartCreate(ctx, data, options)

// 或者添加自定义validation、hooks等
crud.AddValidation("email", &sqlbuilder.EmailRule{})
crud.AddHook(constant.HookEventAfterCreate, myCustomHook)
```

## 🆚 对比说明

### ✅ 现在的优雅设计

```go
// 简单使用
crud, _ := sqlbuilder.NewSimple(db, "users")
err := crud.Add(data)  // 自动时间戳、软删除

// 高级使用
result, err := crud.SmartCreate(ctx, data, options)  // 完整功能
```

### ❌ 之前的重复造轮子

```go
// 简单功能：单独实现一套EasyCRUD
// 高级功能：再实现一套EnhancedBuilder
// 结果：代码重复，功能割裂
```

## 🛡️ 自动功能

你不需要关心这些，但它们都自动处理了：

- ✅ **自动时间戳**: Hook机制自动添加`created_at`和`updated_at`
- ✅ **软删除**: 自动设置`deleted_at`，查询时自动过滤
- ✅ **友好错误**: 错误信息是中文，直接易懂
- ✅ **参数检查**: 自动检查必要参数
- ✅ **分页保护**: 自动限制分页参数防止查询过多数据

## 🎮 真实使用场景

### 用户管理系统

```go
users, _ := sqlbuilder.NewSimple(db, "users")

// 用户注册
err := users.Add(map[string]interface{}{
    "username": "john123",
    "email": "john@email.com",
    "password": "hashedPassword",
})

// 用户登录验证（需要复杂查询时）
ctx := context.Background()
findOptions := &sqlbuilder.FindOptions{
    Filters: []*sqlbuilder.EnhancedFilter{
        {Field: "email", Operator: "=", Value: "john@email.com"},
        {Field: "status", Operator: "=", Value: "active"},
    },
}
result, err := users.SmartFind(ctx, findOptions)

// 更新登录时间
err := users.Update(userID, map[string]interface{}{
    "last_login": time.Now(),
})
```

## 🎯 架构优势

**基于EnhancedBuilder的分层设计：**

```
简单接口层    │  Add(), Get(), Update(), Delete()
            │  ↓ 调用
增强功能层    │  SmartCreate(), SmartUpdate(), SmartFind()
            │  ↓ 调用  
核心构建器    │  Table(), Where(), Insert(), Select()
```

**好处：**

- 🎯 **简单场景**：用简单接口，代码清爽
- 🚀 **复杂场景**：用增强接口，功能完整
- 🔧 **极端场景**：直接用核心构建器，完全控制

## 📦 安装使用

```bash
go get github.com/kamalyes/go-sqlbuilder
```

```go
import "github.com/kamalyes/go-sqlbuilder"

// 创建简单操作器
crud, err := sqlbuilder.NewSimple(db, "table_name")

// 开始使用！
err := crud.Add(data)
```

## 🎉 总结

这就是我们想要的数据库操作方式：

1. **简单时简单**: 一行代码搞定CRUD
2. **复杂时强大**: 完整的EnhancedBuilder功能随时可用
3. **无缝切换**: 同一个对象，简单和复杂方法并存
4. **避免重复**: 不再有两套代码做同样的事

**终于可以专注业务逻辑，简单时简单，复杂时强大！** 🎉
