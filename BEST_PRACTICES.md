# Go SQL Builder - 推荐与不推荐写法指南

## 📖 目录

- [常量定义](#常量定义)
- [CRUD操作](#crud操作)
- [错误处理](#错误处理)
- [并发安全](#并发安全)
- [性能优化](#性能优化)
- [架构设计](#架构设计)
- [测试最佳实践](#测试最佳实践)

---

## 🔧 常量定义

### ✅ 推荐写法

```go
// 定义常量集中管理
const (
    // 数据库操作类型
    OperationTypeCreate = "create"
    OperationTypeUpdate = "update" 
    OperationTypeDelete = "delete"
    OperationTypeUpsert = "upsert"
    
    // 钩子事件类型
    HookBeforeCreate     = "beforeCreate"
    HookAfterCreate      = "afterCreate"
    HookBeforeUpdate     = "beforeUpdate"
    HookAfterUpdate      = "afterUpdate"
    HookBeforeDelete     = "beforeDelete"
    HookAfterDelete      = "afterDelete"
    HookBeforeBatchUpsert = "beforeBatchUpsert"
    HookAfterBatchUpsert  = "afterBatchUpsert"
    
    // 审计字段名称
    AuditFieldCreatedAt = "created_at"
    AuditFieldUpdatedAt = "updated_at"
    AuditFieldDeletedAt = "deleted_at"
    AuditFieldVersion   = "version"
    
    // SQL操作符
    OperatorEqual        = "="
    OperatorNotEqual     = "!="
    OperatorGreater      = ">"
    OperatorGreaterEqual = ">="
    OperatorLess         = "<"
    OperatorLessEqual    = "<="
    OperatorLike         = "LIKE"
    OperatorIn           = "IN"
    OperatorNotIn        = "NOT IN"
    OperatorIsNull       = "IS NULL"
    OperatorIsNotNull    = "IS NOT NULL"
    OperatorBetween      = "BETWEEN"
    
    // 排序方向
    OrderDirectionAsc  = "ASC"
    OrderDirectionDesc = "DESC"
    
    // 默认配置
    DefaultBatchSize    = 1000
    DefaultTimeout      = 30 * time.Second
    DefaultPageSize     = 20
    DefaultMaxRetries   = 3
)

// 使用常量
func (eb *EnhancedBuilder) SmartCreate(ctx context.Context, data map[string]interface{}, options *CreateOptions) (*CreateResult, error) {
    // 执行前置钩子
    if err := eb.executeHooks(HookBeforeCreate, data); err != nil {
        return nil, err
    }
    
    // 添加审计字段
    eb.addAuditFields(data, OperationTypeCreate)
    
    // ...其他逻辑
}
```

### ❌ 不推荐写法

```go
// 硬编码字符串 - 容易出错，难以维护
func (eb *EnhancedBuilder) SmartCreate(ctx context.Context, data map[string]interface{}, options *CreateOptions) (*CreateResult, error) {
    // ❌ 硬编码的钩子名称
    if err := eb.executeHooks("beforeCreate", data); err != nil {
        return nil, err
    }
    
    // ❌ 硬编码的字段名
    data["created_at"] = time.Now()
    data["updated_at"] = time.Now()
    
    // ❌ 硬编码的数值
    query := builder.Limit(1000) // 魔法数字
    
    return result, nil
}

// ❌ 重复的字符串常量
func applyFilter(filter *Filter) {
    switch filter.Operator {
    case "=":     // 重复定义
    case "!=":    // 重复定义
    case ">":     // 重复定义
    case "LIKE":  // 重复定义
    // ...
    }
}
```

---

## 📝 CRUD操作

### ✅ 推荐写法 - 查询操作

```go
// 使用Builder模式，类型安全，支持链式调用
func GetActiveUsers(ctx context.Context, ageRange [2]int) ([]*User, error) {
    builder, err := sqlbuilder.NewEnhanced(db)
    if err != nil {
        return nil, err
    }
    
    options := &sqlbuilder.FindOptions{
        Filters: []*sqlbuilder.EnhancedFilter{
            {Field: "status", Operator: OperatorEqual, Value: "active"},
            {Field: "age", Operator: OperatorBetween, Value: ageRange},
            {Field: "email", Operator: OperatorIsNotNull, Value: nil},
        },
        Orders: []*sqlbuilder.OrderOption{
            {Field: AuditFieldCreatedAt, Direction: OrderDirectionDesc},
        },
        Limit:      DefaultPageSize,
        CountTotal: true,
    }
    
    result, err := builder.SmartFind(ctx, options)
    if err != nil {
        return nil, errors.NewErrorf(errors.ErrorCodeDBError, "查询活跃用户失败: %v", err)
    }
    
    return convertToUsers(result.Records), nil
}

// 使用Repository模式 - 更高层次抽象
func (ur *UserRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
    var user User
    err := ur.db.WithContext(ctx).
        Where("email = ?", email).
        Where("deleted_at IS NULL").
        First(&user).Error
        
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errors.NewError(errors.ErrorCodeNotFound, "用户不存在")
        }
        return nil, errors.NewErrorf(errors.ErrorCodeDBError, "查询用户失败: %v", err)
    }
    
    return &user, nil
}
```

### ❌ 不推荐写法 - 查询操作

```go
// ❌ 直接SQL拼接 - SQL注入风险
func GetUsersBad(name string) ([]*User, error) {
    query := fmt.Sprintf("SELECT * FROM users WHERE name = '%s'", name) // SQL注入风险!
    rows, err := db.Query(query)
    // ...
}

// ❌ 不处理错误
func GetUsersBad2() []*User {
    var users []*User
    db.Find(&users) // 忽略错误
    return users
}

// ❌ 没有分页，可能查询海量数据
func GetAllUsersBad() ([]*User, error) {
    var users []*User
    err := db.Find(&users).Error // 可能返回百万条记录
    return users, err
}
```

### ✅ 推荐写法 - 创建操作

```go
// 带验证和事务的创建操作
func CreateUserWithProfile(ctx context.Context, userData *User, profileData *Profile) (*User, error) {
    // 数据验证
    if err := validateUserData(userData); err != nil {
        return nil, errors.NewErrorf(errors.ErrorCodeInvalidInput, "用户数据验证失败: %v", err)
    }
    
    builder, err := sqlbuilder.New(db)
    if err != nil {
        return nil, err
    }
    
    var createdUser *User
    err = builder.Transaction(func(tx *sqlbuilder.Builder) error {
        // 创建用户
        userMap := map[string]interface{}{
            "name":                userData.Name,
            "email":               userData.Email,
            "password_hash":       hashPassword(userData.Password),
            AuditFieldCreatedAt:   time.Now(),
            AuditFieldUpdatedAt:   time.Now(),
            AuditFieldVersion:     1,
        }
        
        userID, err := tx.WithContext(ctx).
            Table("users").
            InsertGetID(userMap)
        if err != nil {
            return errors.NewErrorf(errors.ErrorCodeDBFailedInsert, "创建用户失败: %v", err)
        }
        
        // 创建用户档案
        profileMap := map[string]interface{}{
            "user_id":             userID,
            "bio":                 profileData.Bio,
            "avatar_url":          profileData.AvatarURL,
            AuditFieldCreatedAt:   time.Now(),
        }
        
        _, err = tx.WithContext(ctx).
            Table("user_profiles").
            InsertGetID(profileMap)
        if err != nil {
            return errors.NewErrorf(errors.ErrorCodeDBFailedInsert, "创建用户档案失败: %v", err)
        }
        
        userData.ID = userID
        createdUser = userData
        return nil
    })
    
    return createdUser, err
}

// 批量创建 - 高性能
func CreateUsersInBatch(ctx context.Context, users []*User) error {
    if len(users) == 0 {
        return nil
    }
    
    // 分批处理，避免单次操作数据量过大
    batchSize := DefaultBatchSize
    for i := 0; i < len(users); i += batchSize {
        end := i + batchSize
        if end > len(users) {
            end = len(users)
        }
        
        batch := users[i:end]
        data := make([]map[string]interface{}, len(batch))
        
        for j, user := range batch {
            if err := validateUserData(user); err != nil {
                return errors.NewErrorf(errors.ErrorCodeInvalidInput, "批次中用户数据验证失败[%d]: %v", i+j, err)
            }
            
            data[j] = map[string]interface{}{
                "name":               user.Name,
                "email":              user.Email,
                "password_hash":      hashPassword(user.Password),
                AuditFieldCreatedAt:  time.Now(),
                AuditFieldUpdatedAt:  time.Now(),
                AuditFieldVersion:    1,
            }
        }
        
        builder, err := sqlbuilder.New(db)
        if err != nil {
            return err
        }
        
        if err := builder.WithContext(ctx).Table("users").BatchInsert(data); err != nil {
            return errors.NewErrorf(errors.ErrorCodeDBFailedInsert, "批量创建用户失败: %v", err)
        }
    }
    
    return nil
}
```

### ❌ 不推荐写法 - 创建操作

```go
// ❌ 不验证输入数据
func CreateUserBad(user *User) error {
    data := map[string]interface{}{
        "name":  user.Name,  // 可能为空
        "email": user.Email, // 可能为空或格式错误
        "age":   user.Age,   // 可能为负数
    }
    _, err := db.Create(data).Error
    return err // 没有包装错误
}

// ❌ 在循环中逐条插入
func CreateUsersBad(users []*User) error {
    for _, user := range users {
        if err := CreateUserBad(user); err != nil {
            return err // 每次都是单独事务，性能差
        }
    }
    return nil
}

// ❌ 不处理关联数据
func CreateUserWithProfileBad(user *User, profile *Profile) error {
    // 分别创建，没有事务保护
    if err := db.Create(user).Error; err != nil {
        return err
    }
    
    profile.UserID = user.ID
    return db.Create(profile).Error // 如果失败，用户已经创建，数据不一致
}
```

### ✅ 推荐写法 - 更新操作

```go
// 带乐观锁的更新操作
func UpdateUserSafely(ctx context.Context, userID int64, updates map[string]interface{}, version int64) error {
    if len(updates) == 0 {
        return errors.NewError(errors.ErrorCodeInvalidInput, "没有需要更新的字段")
    }
    
    // 验证更新数据
    if err := validateUpdateData(updates); err != nil {
        return errors.NewErrorf(errors.ErrorCodeInvalidInput, "更新数据验证失败: %v", err)
    }
    
    builder, err := sqlbuilder.NewEnhanced(db)
    if err != nil {
        return err
    }
    
    builder.AddAuditFields(AuditFieldUpdatedAt)
    
    options := &sqlbuilder.UpdateOptions{
        Version: version, // 乐观锁
    }
    
    result, err := builder.SmartUpdate(ctx, userID, updates, options)
    if err != nil {
        return err
    }
    
    if result.RowsAffected == 0 {
        return errors.NewError(errors.ErrorCodeNotFound, "用户不存在或已被其他进程修改")
    }
    
    return nil
}

// 条件更新
func UpdateUsersByStatus(ctx context.Context, oldStatus, newStatus string) (int64, error) {
    updates := map[string]interface{}{
        "status":              newStatus,
        AuditFieldUpdatedAt:   time.Now(),
    }
    
    builder, err := sqlbuilder.New(db)
    if err != nil {
        return 0, err
    }
    
    result, err := builder.WithContext(ctx).
        Table("users").
        Where("status", OperatorEqual, oldStatus).
        WhereNull(AuditFieldDeletedAt).
        Update(updates)
    
    if err != nil {
        return 0, errors.NewErrorf(errors.ErrorCodeDBFailedUpdate, "批量更新用户状态失败: %v", err)
    }
    
    affected, _ := result.RowsAffected()
    return affected, nil
}
```

### ❌ 不推荐写法 - 更新操作

```go
// ❌ 直接更新整个结构体
func UpdateUserBad(user *User) error {
    return db.Save(user).Error // 可能会覆盖不应该改变的字段
}

// ❌ 没有WHERE条件的更新 - 非常危险!
func UpdateAllUsersBad(status string) error {
    updates := map[string]interface{}{"status": status}
    return db.Model(&User{}).Updates(updates).Error // 更新所有记录！
}

// ❌ 不检查影响行数
func UpdateUserBad2(userID int64, name string) error {
    result := db.Model(&User{}).Where("id = ?", userID).Update("name", name)
    return result.Error // 不知道是否真的更新了记录
}

// ❌ 没有乐观锁的并发更新
func UpdateUserConcurrentBad(userID int64, updates map[string]interface{}) error {
    // 多个进程同时更新可能导致数据覆盖
    return db.Model(&User{}).Where("id = ?", userID).Updates(updates).Error
}
```

### ✅ 推荐写法 - 删除操作

```go
// 软删除 - 生产环境推荐
func SoftDeleteUser(ctx context.Context, userID int64, operatorID int64) error {
    updates := map[string]interface{}{
        AuditFieldDeletedAt:  time.Now(),
        AuditFieldUpdatedAt:  time.Now(),
        "deleted_by":         operatorID, // 记录删除者
    }
    
    builder, err := sqlbuilder.New(db)
    if err != nil {
        return err
    }
    
    result, err := builder.WithContext(ctx).
        Table("users").
        Where("id", OperatorEqual, userID).
        WhereNull(AuditFieldDeletedAt). // 确保不是已删除的记录
        Update(updates)
        
    if err != nil {
        return errors.NewErrorf(errors.ErrorCodeDBFailedDelete, "软删除用户失败: %v", err)
    }
    
    affected, _ := result.RowsAffected()
    if affected == 0 {
        return errors.NewError(errors.ErrorCodeNotFound, "用户不存在或已被删除")
    }
    
    return nil
}

// 硬删除 - 需要特殊权限和审计
func HardDeleteUser(ctx context.Context, userID int64, operatorID int64) error {
    // 记录删除操作到审计日志
    auditLog := map[string]interface{}{
        "action":      "hard_delete_user",
        "target_id":   userID,
        "operator_id": operatorID,
        "timestamp":   time.Now(),
    }
    
    builder, err := sqlbuilder.New(db)
    if err != nil {
        return err
    }
    
    return builder.Transaction(func(tx *sqlbuilder.Builder) error {
        // 1. 记录审计日志
        _, err := tx.WithContext(ctx).Table("audit_logs").InsertGetID(auditLog)
        if err != nil {
            return errors.NewErrorf(errors.ErrorCodeDBError, "记录审计日志失败: %v", err)
        }
        
        // 2. 删除关联数据
        _, err = tx.Table("user_profiles").Where("user_id", OperatorEqual, userID).Delete().Exec()
        if err != nil {
            return errors.NewErrorf(errors.ErrorCodeDBFailedDelete, "删除用户档案失败: %v", err)
        }
        
        // 3. 删除主记录
        result, err := tx.Table("users").Where("id", OperatorEqual, userID).Delete().Exec()
        if err != nil {
            return errors.NewErrorf(errors.ErrorCodeDBFailedDelete, "删除用户失败: %v", err)
        }
        
        affected, _ := result.RowsAffected()
        if affected == 0 {
            return errors.NewError(errors.ErrorCodeNotFound, "用户不存在")
        }
        
        return nil
    })
}

// 批量软删除
func BatchSoftDeleteUsers(ctx context.Context, userIDs []int64, operatorID int64) error {
    if len(userIDs) == 0 {
        return nil
    }
    
    updates := map[string]interface{}{
        AuditFieldDeletedAt: time.Now(),
        AuditFieldUpdatedAt: time.Now(),
        "deleted_by":        operatorID,
    }
    
    builder, err := sqlbuilder.New(db)
    if err != nil {
        return err
    }
    
    // 转换为interface{}切片
    ids := make([]interface{}, len(userIDs))
    for i, id := range userIDs {
        ids[i] = id
    }
    
    result, err := builder.WithContext(ctx).
        Table("users").
        WhereIn("id", ids...).
        WhereNull(AuditFieldDeletedAt).
        Update(updates)
        
    if err != nil {
        return errors.NewErrorf(errors.ErrorCodeDBFailedDelete, "批量软删除用户失败: %v", err)
    }
    
    affected, _ := result.RowsAffected()
    if affected != int64(len(userIDs)) {
        return errors.NewErrorf(errors.ErrorCodePartialFailure, "期望删除%d个用户，实际删除%d个", len(userIDs), affected)
    }
    
    return nil
}
```

### ❌ 不推荐写法 - 删除操作

```go
// ❌ 直接硬删除 - 数据无法恢复
func DeleteUserBad(userID int64) error {
    return db.Delete(&User{}, userID).Error // 数据永久丢失
}

// ❌ 没有WHERE条件 - 非常危险!
func DeleteAllUsersBad() error {
    return db.Delete(&User{}).Error // 删除所有用户！
}

// ❌ 不在事务中处理关联删除
func DeleteUserWithDataBad(userID int64) error {
    // 分别删除，可能导致数据不一致
    db.Where("user_id = ?", userID).Delete(&UserProfile{})
    db.Where("user_id = ?", userID).Delete(&UserSetting{})
    return db.Delete(&User{}, userID).Error
}

// ❌ 没有记录删除操作
func DeleteUserNoAuditBad(userID int64) error {
    // 没有记录谁删除了什么，无法追溯
    return db.Delete(&User{}, userID).Error
}
```

---

## 🚨 错误处理

### ✅ 推荐写法

```go
// 使用自定义错误类型
type AppError struct {
    Code    ErrorCode `json:"code"`
    Message string    `json:"message"`
    Details string    `json:"details,omitempty"`
    Cause   error     `json:"-"`
}

const (
    // 错误代码常量
    ErrorCodeSuccess         ErrorCode = 0
    ErrorCodeNotFound        ErrorCode = 1001
    ErrorCodeAlreadyExist    ErrorCode = 1002
    ErrorCodeInvalidInput    ErrorCode = 1003
    ErrorCodeDBError         ErrorCode = 2001
    ErrorCodeDBFailedInsert  ErrorCode = 2002
    ErrorCodeDBFailedUpdate  ErrorCode = 2003
    ErrorCodeDBFailedDelete  ErrorCode = 2004
)

// 错误包装和传播
func (s *UserService) GetUser(ctx context.Context, userID int64) (*User, error) {
    user, err := s.repo.GetByID(ctx, userID)
    if err != nil {
        if IsErrorCode(err, ErrorCodeNotFound) {
            return nil, NewError(ErrorCodeNotFound, "用户不存在")
        }
        return nil, NewErrorf(ErrorCodeDBError, "获取用户失败: %v", err)
    }
    return user, nil
}

// 错误恢复机制
func (s *UserService) ProcessWithRetry(ctx context.Context, userID int64) error {
    for i := 0; i < DefaultMaxRetries; i++ {
        if err := s.processUser(ctx, userID); err != nil {
            if !isRetryableError(err) {
                return err // 不可重试错误
            }
            
            if i == DefaultMaxRetries-1 {
                return NewErrorf(ErrorCodeOperationFailed, 
                    "处理失败，已重试%d次: %v", DefaultMaxRetries, err)
            }
            
            // 指数退避
            time.Sleep(time.Duration(i+1) * time.Second)
            continue
        }
        return nil
    }
    return nil
}
```

### ❌ 不推荐写法

```go
// ❌ 忽略错误
func GetUserBad(userID int64) *User {
    user, _ := repo.GetByID(userID) // 忽略错误
    return user
}

// ❌ 不提供错误上下文
func ProcessUserBad(userID int64) error {
    err := someOperation(userID)
    if err != nil {
        return err // 没有添加上下文
    }
    return nil
}

// ❌ 使用panic处理业务错误
func GetUserPanicBad(userID int64) *User {
    user, err := repo.GetByID(userID)
    if err != nil {
        panic(err) // 不应该在业务逻辑中使用panic
    }
    return user
}
```

---

## 🔒 并发安全

### ✅ 推荐写法

```go
// 使用读写锁保护共享状态
type SafeCache struct {
    mu    sync.RWMutex
    cache map[string]*CacheItem
}

type CacheItem struct {
    Value     interface{}
    ExpiresAt time.Time
}

func (c *SafeCache) Get(key string) (interface{}, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    item, exists := c.cache[key]
    if !exists {
        return nil, false
    }
    
    if time.Now().After(item.ExpiresAt) {
        return nil, false // 已过期
    }
    
    return item.Value, true
}

func (c *SafeCache) Set(key string, value interface{}, ttl time.Duration) {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    c.cache[key] = &CacheItem{
        Value:     value,
        ExpiresAt: time.Now().Add(ttl),
    }
}

// 使用context控制超时和取消
func (s *UserService) ProcessWithTimeout(ctx context.Context, userID int64) error {
    ctx, cancel := context.WithTimeout(ctx, DefaultTimeout)
    defer cancel()
    
    done := make(chan error, 1)
    go func() {
        done <- s.doHeavyWork(userID)
    }()
    
    select {
    case err := <-done:
        return err
    case <-ctx.Done():
        return NewError(ErrorCodeTimeout, "处理超时")
    }
}

// 安全的goroutine使用
func (s *UserService) ProcessUsersConcurrently(userIDs []int64) error {
    const maxConcurrency = 10
    sem := make(chan struct{}, maxConcurrency)
    errCh := make(chan error, len(userIDs))
    
    var wg sync.WaitGroup
    
    for _, userID := range userIDs {
        wg.Add(1)
        go func(id int64) {
            defer wg.Done()
            
            sem <- struct{}{} // 获取信号量
            defer func() { <-sem }() // 释放信号量
            
            if err := s.processUser(context.Background(), id); err != nil {
                errCh <- NewErrorf(ErrorCodeOperationFailed, "处理用户%d失败: %v", id, err)
                return
            }
        }(userID)
    }
    
    wg.Wait()
    close(errCh)
    
    // 收集错误
    var errors []error
    for err := range errCh {
        errors = append(errors, err)
    }
    
    if len(errors) > 0 {
        return NewErrorf(ErrorCodePartialFailure, "处理完成，但有%d个错误", len(errors))
    }
    
    return nil
}
```

### ❌ 不推荐写法

```go
// ❌ 没有并发保护的共享状态
var globalCounter int // 竞态条件

func IncrementCounterBad() {
    globalCounter++ // 竞态条件
}

// ❌ 没有超时控制
func ProcessRequestBad() error {
    result := <-someChannel // 可能无限阻塞
    return processResult(result)
}

// ❌ 不安全的goroutine使用
func ProcessUsersBad(userIDs []int64) {
    for _, userID := range userIDs {
        go func(id int64) {
            processUser(id) // 没有错误处理
        }(userID)
    }
    // 没有等待goroutine完成
}
```

---

## 🚀 性能优化

### ✅ 推荐写法

```go
// 使用对象池减少GC压力
var (
    queryBuilderPool = sync.Pool{
        New: func() interface{} {
            return &strings.Builder{}
        },
    }
    
    userSlicePool = sync.Pool{
        New: func() interface{} {
            return make([]*User, 0, DefaultBatchSize)
        },
    }
)

func BuildComplexQuery(filters []Filter) string {
    builder := queryBuilderPool.Get().(*strings.Builder)
    defer func() {
        builder.Reset()
        queryBuilderPool.Put(builder)
    }()
    
    builder.WriteString("SELECT * FROM users WHERE 1=1")
    for _, filter := range filters {
        builder.WriteString(" AND ")
        builder.WriteString(filter.ToSQL())
    }
    
    return builder.String()
}

// 预分配切片容量
func ProcessUsers(users []*User) []*ProcessedUser {
    results := make([]*ProcessedUser, 0, len(users)) // 预分配容量
    
    for _, user := range users {
        if processed := processUser(user); processed != nil {
            results = append(results, processed)
        }
    }
    
    return results
}

// 批量操作减少数据库调用
func (r *UserRepository) BatchCreateOptimized(ctx context.Context, users []*User) error {
    if len(users) == 0 {
        return nil
    }
    
    // 分批处理
    for i := 0; i < len(users); i += DefaultBatchSize {
        end := i + DefaultBatchSize
        if end > len(users) {
            end = len(users)
        }
        
        batch := users[i:end]
        data := make([]map[string]interface{}, len(batch))
        
        for j, user := range batch {
            data[j] = map[string]interface{}{
                "name":                user.Name,
                "email":               user.Email,
                AuditFieldCreatedAt:   time.Now(),
                AuditFieldUpdatedAt:   time.Now(),
            }
        }
        
        if err := r.db.WithContext(ctx).CreateInBatches(data, DefaultBatchSize).Error; err != nil {
            return NewErrorf(ErrorCodeDBFailedInsert, "批量创建失败: %v", err)
        }
    }
    
    return nil
}

// 使用索引优化查询
func (r *UserRepository) FindActiveUsersOptimized(ctx context.Context, limit int) ([]*User, error) {
    users := userSlicePool.Get().([]*User)
    defer func() {
        users = users[:0] // 重置切片长度
        userSlicePool.Put(users)
    }()
    
    // 使用覆盖索引避免回表
    err := r.db.WithContext(ctx).
        Select("id, name, email, status"). // 只选择需要的字段
        Where("status = ?", "active").     // 使用索引字段
        Where("deleted_at IS NULL").       // 使用索引字段
        Order("created_at DESC").          // 使用索引排序
        Limit(limit).
        Find(&users).Error
    
    if err != nil {
        return nil, NewErrorf(ErrorCodeDBError, "查询活跃用户失败: %v", err)
    }
    
    // 复制结果以避免池对象被修改
    result := make([]*User, len(users))
    copy(result, users)
    
    return result, nil
}
```

### ❌ 不推荐写法

```go
// ❌ 频繁的内存分配
func ProcessDataBad(items []string) []string {
    var result []string // 没有预分配容量
    for _, item := range items {
        result = append(result, strings.ToUpper(item))
        // 每次append可能触发重新分配
    }
    return result
}

// ❌ N+1查询问题
func GetUsersWithProfilesBad(userIDs []int64) ([]*UserWithProfile, error) {
    var results []*UserWithProfile
    
    for _, userID := range userIDs {
        // 每个用户一次查询 - N+1问题
        user, _ := getUserByID(userID)
        profile, _ := getProfileByUserID(userID)
        
        results = append(results, &UserWithProfile{
            User:    user,
            Profile: profile,
        })
    }
    
    return results, nil
}

// ❌ 查询所有字段
func FindUsersBad() ([]*User, error) {
    var users []*User
    // 查询所有字段，包括大字段
    err := db.Find(&users).Error
    return users, err
}

// ❌ 不使用索引
func FindUsersByNameBad(name string) ([]*User, error) {
    var users []*User
    // LIKE查询不使用索引
    err := db.Where("UPPER(name) LIKE ?", "%"+strings.ToUpper(name)+"%").Find(&users).Error
    return users, err
}
```

---

## 🏗️ 架构设计

### ✅ 推荐写法

```go
// 使用接口隔离原则
type UserReader interface {
    GetByID(ctx context.Context, id int64) (*User, error)
    GetByEmail(ctx context.Context, email string) (*User, error)
    List(ctx context.Context, filters *UserFilters) ([]*User, error)
}

type UserWriter interface {
    Create(ctx context.Context, user *User) (*User, error)
    Update(ctx context.Context, id int64, updates map[string]interface{}) error
    Delete(ctx context.Context, id int64) error
}

type UserRepository interface {
    UserReader
    UserWriter
}

// 依赖注入
type UserService struct {
    repo   UserRepository
    logger Logger
    cache  Cache
    config *Config
}

func NewUserService(repo UserRepository, logger Logger, cache Cache, config *Config) *UserService {
    return &UserService{
        repo:   repo,
        logger: logger,
        cache:  cache,
        config: config,
    }
}

func (s *UserService) CreateUser(ctx context.Context, req *CreateUserRequest) (*User, error) {
    // 使用常量和配置
    if len(req.Name) < s.config.MinNameLength {
        return nil, NewError(ErrorCodeInvalidInput, "用户名长度不足")
    }
    
    // 记录日志
    s.logger.InfoContext(ctx, "创建用户", 
        "name", req.Name, 
        "email", req.Email)
    
    user := &User{
        Name:  req.Name,
        Email: req.Email,
        Status: UserStatusActive, // 使用常量
    }
    
    result, err := s.repo.Create(ctx, user)
    if err != nil {
        s.logger.ErrorContext(ctx, "创建用户失败", "error", err)
        return nil, err
    }
    
    // 清除相关缓存
    s.cache.Delete(fmt.Sprintf("user:email:%s", user.Email))
    
    return result, nil
}

// 使用策略模式
type ValidationStrategy interface {
    Validate(user *User) error
}

type EmailValidationStrategy struct{}

func (v *EmailValidationStrategy) Validate(user *User) error {
    if !isValidEmail(user.Email) {
        return NewError(ErrorCodeInvalidInput, "邮箱格式无效")
    }
    return nil
}

type UserValidator struct {
    strategies []ValidationStrategy
}

func (v *UserValidator) Validate(user *User) error {
    for _, strategy := range v.strategies {
        if err := strategy.Validate(user); err != nil {
            return err
        }
    }
    return nil
}
```

### ❌ 不推荐写法

```go
// ❌ 违反单一职责原则
type UserHandler struct {
    // 混合了HTTP、业务逻辑、数据访问
    db     *gorm.DB
    redis  *redis.Client
    logger *log.Logger
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    // 直接在HTTP处理器中写业务逻辑
    var req CreateUserRequest
    json.NewDecoder(r.Body).Decode(&req)
    
    // 硬编码的验证
    if len(req.Name) < 2 { // 魔法数字
        http.Error(w, "Name too short", 400)
        return
    }
    
    user := &User{Name: req.Name}
    h.db.Create(user) // 没有错误处理
    
    json.NewEncoder(w).Encode(user)
}

// ❌ 全局变量
var (
    DB    *gorm.DB    // 全局数据库连接
    Cache *redis.Client // 全局缓存
)

func CreateUser(user *User) error {
    return DB.Create(user).Error // 紧耦合
}

// ❌ 上帝对象
type UserService struct {
    // 处理所有用户相关的操作，职责过多
    db           *gorm.DB
    cache        *redis.Client
    emailService *EmailService
    smsService   *SMSService
    paymentService *PaymentService
    // ... 更多依赖
}

func (s *UserService) DoEverything() {
    // 一个方法做太多事情
    // 创建用户
    // 发送邮件
    // 发送短信
    // 处理支付
    // 记录日志
    // 更新缓存
    // ...
}
```

---

## 🧪 测试最佳实践

### ✅ 推荐写法

```go
// 表驱动测试
func TestUserValidation(t *testing.T) {
    tests := []struct {
        name    string
        user    *User
        wantErr bool
        errCode ErrorCode
    }{
        {
            name: "valid_user",
            user: &User{
                Name:  "张三",
                Email: "zhangsan@example.com",
                Age:   25,
            },
            wantErr: false,
        },
        {
            name: "invalid_email",
            user: &User{
                Name:  "张三",
                Email: "invalid-email",
                Age:   25,
            },
            wantErr: true,
            errCode: ErrorCodeInvalidInput,
        },
        {
            name: "empty_name",
            user: &User{
                Name:  "",
                Email: "zhangsan@example.com",
                Age:   25,
            },
            wantErr: true,
            errCode: ErrorCodeInvalidInput,
        },
        {
            name: "negative_age",
            user: &User{
                Name:  "张三",
                Email: "zhangsan@example.com",
                Age:   -1,
            },
            wantErr: true,
            errCode: ErrorCodeInvalidInput,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            validator := NewUserValidator()
            err := validator.Validate(tt.user)
            
            if tt.wantErr {
                require.Error(t, err)
                if tt.errCode != 0 {
                    var appErr *AppError
                    require.True(t, errors.As(err, &appErr))
                    assert.Equal(t, tt.errCode, appErr.Code)
                }
            } else {
                require.NoError(t, err)
            }
        })
    }
}

// 使用Mock进行单元测试
func TestUserService_CreateUser(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    mockRepo := NewMockUserRepository(ctrl)
    mockLogger := NewMockLogger(ctrl)
    mockCache := NewMockCache(ctrl)
    
    service := NewUserService(mockRepo, mockLogger, mockCache, &Config{
        MinNameLength: DefaultMinNameLength,
    })
    
    tests := []struct {
        name    string
        request *CreateUserRequest
        setup   func()
        wantErr bool
        errCode ErrorCode
    }{
        {
            name: "success",
            request: &CreateUserRequest{
                Name:  "张三",
                Email: "zhangsan@example.com",
            },
            setup: func() {
                expectedUser := &User{
                    Name:   "张三",
                    Email:  "zhangsan@example.com",
                    Status: UserStatusActive,
                }
                returnUser := &User{
                    ID:     1,
                    Name:   "张三",
                    Email:  "zhangsan@example.com",
                    Status: UserStatusActive,
                }
                
                mockRepo.EXPECT().
                    Create(gomock.Any(), expectedUser).
                    Return(returnUser, nil)
                    
                mockLogger.EXPECT().
                    InfoContext(gomock.Any(), "创建用户", gomock.Any())
                    
                mockCache.EXPECT().
                    Delete("user:email:zhangsan@example.com")
            },
            wantErr: false,
        },
        {
            name: "invalid_name_length",
            request: &CreateUserRequest{
                Name:  "a", // 太短
                Email: "test@example.com",
            },
            setup: func() {
                // 不期望调用repo和cache
            },
            wantErr: true,
            errCode: ErrorCodeInvalidInput,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            tt.setup()
            
            user, err := service.CreateUser(context.Background(), tt.request)
            
            if tt.wantErr {
                require.Error(t, err)
                assert.Nil(t, user)
                if tt.errCode != 0 {
                    var appErr *AppError
                    require.True(t, errors.As(err, &appErr))
                    assert.Equal(t, tt.errCode, appErr.Code)
                }
            } else {
                require.NoError(t, err)
                assert.NotNil(t, user)
            }
        })
    }
}

// 集成测试
func TestUserRepository_Integration(t *testing.T) {
    if testing.Short() {
        t.Skip("跳过集成测试")
    }
    
    db := setupTestDB(t)
    defer teardownTestDB(t, db)
    
    repo := NewUserRepository(db)
    ctx := context.Background()
    
    t.Run("create_and_get_user", func(t *testing.T) {
        user := &User{
            Name:   "测试用户",
            Email:  "test@example.com",
            Status: UserStatusActive,
        }
        
        // 创建用户
        created, err := repo.Create(ctx, user)
        require.NoError(t, err)
        assert.NotZero(t, created.ID)
        assert.Equal(t, user.Name, created.Name)
        assert.Equal(t, user.Email, created.Email)
        
        // 获取用户
        retrieved, err := repo.GetByID(ctx, created.ID)
        require.NoError(t, err)
        assert.Equal(t, created.ID, retrieved.ID)
        assert.Equal(t, created.Name, retrieved.Name)
        assert.Equal(t, created.Email, retrieved.Email)
    })
}
```

### ❌ 不推荐写法

```go
// ❌ 没有测试覆盖
// UserService 没有对应的测试用例

// ❌ 依赖外部资源的测试
func TestUserService_Bad(t *testing.T) {
    // 直接连接真实数据库
    db, _ := gorm.Open("mysql", "root:password@/test_db")
    service := NewUserService(db, nil, nil)
    
    // 依赖外部API
    user, err := service.CreateUserFromExternalAPI("https://api.example.com/user")
    assert.NoError(t, err)
}

// ❌ 测试中有副作用
func TestBadSideEffect(t *testing.T) {
    // 修改全局状态
    GlobalConfig.Environment = "test"
    
    // 创建文件但不清理
    file, _ := os.Create("test-file.txt")
    defer file.Close()
    
    // 测试逻辑...
    
    // 没有恢复全局状态
}

// ❌ 测试名称不清晰
func TestUser(t *testing.T) {
    // 测试什么？不清楚
}

func TestUserStuff(t *testing.T) {
    // 测试用户的什么功能？不明确
}

// ❌ 没有验证错误类型
func TestUserValidationBad(t *testing.T) {
    user := &User{Name: ""}
    err := ValidateUser(user)
    
    // 只检查是否有错误，不检查错误类型
    assert.Error(t, err)
}
```

---

## 📌 总结

### 🎯 核心原则

1. **使用常量定义**：避免魔法数字和硬编码字符串
2. **统一错误处理**：使用自定义错误类型，提供上下文信息
3. **接口隔离**：遵循SOLID原则，便于测试和扩展
4. **并发安全**：使用适当的同步机制保护共享状态
5. **性能优化**：预分配内存、使用对象池、批量操作
6. **全面测试**：单元测试、集成测试、表驱动测试

### 🚀 最佳实践清单

- [ ] 所有字符串和数值都使用常量定义
- [ ] 实现统一的错误处理机制
- [ ] 使用软删除保护数据
- [ ] 实现乐观锁防止并发冲突
- [ ] 添加数据验证和审计
- [ ] 使用事务处理复合操作
- [ ] 实现批量操作提高性能
- [ ] 添加上下文传递和超时控制
- [ ] 编写全面的单元测试
- [ ] 使用Mock隔离外部依赖

通过遵循这些推荐写法，您可以构建出高质量、可维护、高性能的Go应用程序。
