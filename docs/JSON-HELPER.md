# JSON 辅助函数

## 概述
JSON 辅助函数用于数据库 JSON 字段与结构体之间的序列化和反序列化。

## SerializeJSON - 序列化

### 基础用法
```go
import "github.com/kamalyes/go-sqlbuilder/repository"

type UserProfile struct {
    Avatar   string `json:"avatar"`
    Nickname string `json:"nickname"`
    Bio      string `json:"bio"`
}

profile := &UserProfile{
    Avatar:   "https://example.com/avatar.jpg",
    Nickname: "张三",
    Bio:      "热爱编程",
}

// 序列化为 JSON 字符串
jsonStr, err := repository.SerializeJSON(profile)
// 结果: '{"avatar":"https://example.com/avatar.jpg","nickname":"张三","bio":"热爱编程"}'
```

### nil 值处理
```go
var profile *UserProfile = nil

jsonStr, err := repository.SerializeJSON(profile)
// 结果: "" (空字符串)
```

## DeserializeJSON - 反序列化

### 基础用法
```go
jsonStr := `{"avatar":"https://example.com/avatar.jpg","nickname":"张三","bio":"热爱编程"}`

profile, err := repository.DeserializeJSON[UserProfile](jsonStr)
if err != nil {
    return err
}

fmt.Printf("昵称: %s", profile.Nickname)
```

### 空字符串处理
```go
jsonStr := ""

profile, err := repository.DeserializeJSON[UserProfile](jsonStr)
// 结果: nil, nil
```

### 空对象处理
```go
jsonStr := "{}"

profile, err := repository.DeserializeJSON[UserProfile](jsonStr)
// 结果: nil, nil
```

## MustSerializeJSON - 忽略错误的序列化

```go
profile := &UserProfile{
    Avatar:   "https://example.com/avatar.jpg",
    Nickname: "张三",
}

// 不需要错误处理
jsonStr := repository.MustSerializeJSON(profile)
```

## MustDeserializeJSON - 忽略错误的反序列化

```go
jsonStr := `{"nickname":"张三"}`

// 不需要错误处理
profile := repository.MustDeserializeJSON[UserProfile](jsonStr)
if profile != nil {
    fmt.Printf("昵称: %s", profile.Nickname)
}
```

## 在模型中使用

```go
package main

import (
    "context"
    "github.com/kamalyes/go-sqlbuilder/repository"
)

type User struct {
    repository.BaseModel
    Name    string `json:"name"`
    Profile string `json:"profile" gorm:"type:json"` // JSON 字段
}

// 获取用户信息时反序列化 Profile
func getUserProfile(ctx context.Context, repo repository.IRepository[User], id uint) (*UserProfile, error) {
    user, err := repo.Get(ctx, id)
    if err != nil {
        return nil, err
    }
    
    return repository.DeserializeJSON[UserProfile](user.Profile)
}

// 保存用户时序列化 Profile
func saveUserProfile(ctx context.Context, repo repository.IRepository[User], userID uint, profile *UserProfile) error {
    jsonStr, err := repository.SerializeJSON(profile)
    if err != nil {
        return err
    }
    
    return repo.UpdateFields(ctx, userID, map[string]interface{}{
        "profile": jsonStr,
    })
}
```
