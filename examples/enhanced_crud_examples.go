/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-14 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-14 00:00:00
 * @FilePath: \go-sqlbuilder\examples\enhanced_crud_examples.go
 * @Description: 增强CRUD操作使用示例
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */

package examples

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/kamalyes/go-sqlbuilder"
)

// UserExample 用户示例
type UserExample struct {
	ID        int64      `json:"id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Age       int        `json:"age"`
	Status    string     `json:"status"`
	Version   int64      `json:"version"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// ==================== 推荐写法示例 ====================

// CreateUserExample 创建用户 - 推荐写法
func CreateUserExample(db interface{}) {
	fmt.Println("=== 创建用户 - 推荐写法 ===")

	// 创建增强构建器
	builder, err := sqlbuilder.NewEnhanced(db)
	if err != nil {
		log.Fatal(err)
	}

	// 配置增强功能
	builder.
		EnableSoftDelete(true).
		AddAuditFields("created_at", "updated_at", "deleted_at").
		AddValidation("name", sqlbuilder.RequiredRule{}).
		AddValidation("email", sqlbuilder.EmailRule{}).
		AddValidation("email", sqlbuilder.LengthRule{Min: 5, Max: 100}).
		AddHook("beforeCreate", func(ctx context.Context, data interface{}) error {
			fmt.Println("执行创建前钩子")
			return nil
		}).
		AddHook("afterCreate", func(ctx context.Context, data interface{}) error {
			fmt.Println("执行创建后钩子")
			return nil
		})

	// 创建用户数据
	userData := map[string]interface{}{
		"name":    "张三",
		"email":   "zhangsan@example.com",
		"age":     25,
		"status":  "active",
		"version": 1,
	}

	ctx := context.Background()
	result, err := builder.SmartCreate(ctx, userData, nil)
	if err != nil {
		log.Printf("创建用户失败: %v", err)
		return
	}

	fmt.Printf("用户创建成功，ID: %d\n", result.ID)
}

// UpdateUserExample 更新用户 - 推荐写法
func UpdateUserExample(db interface{}) {
	fmt.Println("=== 更新用户 - 推荐写法 ===")

	builder, err := sqlbuilder.NewEnhanced(db)
	if err != nil {
		log.Fatal(err)
	}

	// 配置增强功能
	builder.
		EnableSoftDelete(true).
		AddAuditFields("updated_at").
		AddValidation("email", sqlbuilder.EmailRule{})

	updateData := map[string]interface{}{
		"name":  "李四",
		"email": "lisi@example.com",
		"age":   30,
	}

	options := &sqlbuilder.UpdateOptions{
		Version: 1, // 乐观锁版本号
	}

	ctx := context.Background()
	result, err := builder.SmartUpdate(ctx, 1, updateData, options)
	if err != nil {
		log.Printf("更新用户失败: %v", err)
		return
	}

	fmt.Printf("用户更新成功，影响行数: %d\n", result.RowsAffected)
}

// QueryUserExample 查询用户 - 推荐写法
func QueryUserExample(db interface{}) {
	fmt.Println("=== 查询用户 - 推荐写法 ===")

	builder, err := sqlbuilder.NewEnhanced(db)
	if err != nil {
		log.Fatal(err)
	}

	builder.EnableSoftDelete(true)

	// 构建查询条件
	options := &sqlbuilder.FindOptions{
		Filters: []*sqlbuilder.EnhancedFilter{
			{Field: "status", Operator: "=", Value: "active"},
			{Field: "age", Operator: ">=", Value: 18},
			{Field: "email", Operator: "IS NOT NULL", Value: nil},
		},
		Orders: []*sqlbuilder.OrderOption{
			{Field: "created_at", Direction: "DESC"},
			{Field: "name", Direction: "ASC"},
		},
		Limit:      10,
		Offset:     0,
		CountTotal: true,
	}

	ctx := context.Background()
	result, err := builder.SmartFind(ctx, options)
	if err != nil {
		log.Printf("查询用户失败: %v", err)
		return
	}

	fmt.Printf("查询到 %d 条记录，总计 %d 条\n", len(result.Records), result.Total)
	for _, record := range result.Records {
		fmt.Printf("用户: %v\n", record)
	}
}

// DeleteUserExample 删除用户 - 推荐写法
func DeleteUserExample(db interface{}) {
	fmt.Println("=== 删除用户 - 推荐写法 ===")

	builder, err := sqlbuilder.NewEnhanced(db)
	if err != nil {
		log.Fatal(err)
	}

	builder.
		EnableSoftDelete(true).
		AddAuditFields("deleted_at")

	// 软删除（推荐）
	ctx := context.Background()
	result, err := builder.SmartDelete(ctx, 1, nil)
	if err != nil {
		log.Printf("软删除用户失败: %v", err)
		return
	}

	if result.SoftDelete {
		fmt.Printf("用户软删除成功，影响行数: %d\n", result.RowsAffected)
	}

	// 硬删除（谨慎使用）
	hardDeleteOptions := &sqlbuilder.DeleteOptions{Force: true}
	result2, err := builder.SmartDelete(ctx, 2, hardDeleteOptions)
	if err != nil {
		log.Printf("硬删除用户失败: %v", err)
		return
	}

	fmt.Printf("用户硬删除成功，影响行数: %d\n", result2.RowsAffected)
}

// BatchUpsertExample 批量Upsert - 推荐写法
func BatchUpsertExample(db interface{}) {
	fmt.Println("=== 批量Upsert - 推荐写法 ===")

	builder, err := sqlbuilder.NewEnhanced(db)
	if err != nil {
		log.Fatal(err)
	}

	builder.
		AddAuditFields("created_at", "updated_at")

	// 准备批量数据
	batchData := []map[string]interface{}{
		{
			"email":  "user1@example.com",
			"name":   "用户1",
			"age":    25,
			"status": "active",
		},
		{
			"email":  "user2@example.com",
			"name":   "用户2",
			"age":    30,
			"status": "active",
		},
		{
			"email":  "user3@example.com",
			"name":   "用户3",
			"age":    28,
			"status": "active",
		},
	}

	conflictFields := []string{"email"} // 以email作为冲突判断字段

	ctx := context.Background()
	err = builder.BatchUpsert(ctx, batchData, conflictFields)
	if err != nil {
		log.Printf("批量Upsert失败: %v", err)
		return
	}

	fmt.Printf("批量Upsert成功，处理了 %d 条记录\n", len(batchData))
}

// ==================== 不推荐写法示例 ====================

// CreateUserBadExample 创建用户 - 不推荐写法
func CreateUserBadExample(db interface{}) {
	fmt.Println("=== 创建用户 - 不推荐写法 ===")

	// ❌ 直接使用原生SQL，有SQL注入风险
	/*
		query := fmt.Sprintf("INSERT INTO users (name, email) VALUES ('%s', '%s')", name, email)
		db.Exec(query)
	*/

	// ❌ 不验证数据
	/*
		userData := map[string]interface{}{
			"name": "", // 空值
			"email": "invalid-email", // 无效邮箱
		}
		builder.Insert(userData).Exec()
	*/

	// ❌ 不处理错误
	/*
		builder.Insert(userData).Exec() // 忽略错误
	*/

	fmt.Println("这些都是不推荐的写法，请使用上面的推荐写法")
}

// ==================== 事务处理示例 ====================

// TransactionExample 事务处理 - 推荐写法
func TransactionExample(db interface{}) {
	fmt.Println("=== 事务处理 - 推荐写法 ===")

	builder, err := sqlbuilder.New(db)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	err = builder.Transaction(func(tx *sqlbuilder.Builder) error {
		// 在事务中执行多个操作

		// 1. 创建用户
		userData := map[string]interface{}{
			"name":       "事务用户",
			"email":      "transaction@example.com",
			"created_at": time.Now(),
		}

		userID, err := tx.WithContext(ctx).Table("users").InsertGetID(userData)
		if err != nil {
			return err // 自动回滚
		}

		// 2. 创建用户档案
		profileData := map[string]interface{}{
			"user_id":    userID,
			"bio":        "这是一个事务用户",
			"created_at": time.Now(),
		}

		_, err = tx.WithContext(ctx).Table("user_profiles").InsertGetID(profileData)
		if err != nil {
			return err // 自动回滚
		}

		// 如果所有操作都成功，事务会自动提交
		return nil
	})

	if err != nil {
		log.Printf("事务执行失败: %v", err)
		return
	}

	fmt.Println("事务执行成功")
}

// ==================== 高级查询示例 ====================

// AdvancedQueryExample 高级查询示例
func AdvancedQueryExample(db interface{}) {
	fmt.Println("=== 高级查询示例 ===")

	builder, err := sqlbuilder.New(db)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	// 复杂查询示例
	var users []UserExample
	err = builder.
		WithContext(ctx).
		Table("users").
		Select("users.*", "profiles.bio").
		LeftJoin("user_profiles as profiles", "users.id = profiles.user_id").
		Where("users.status", "=", "active").
		Where("users.age", ">=", 18).
		WhereIn("users.role", "admin", "user").
		WhereNotNull("users.email").
		WhereBetween("users.created_at",
			time.Now().AddDate(0, -1, 0),
			time.Now()).
		OrderByDesc("users.created_at").
		OrderBy("users.name").
		Limit(20).
		Get(&users)

	if err != nil {
		log.Printf("高级查询失败: %v", err)
		return
	}

	fmt.Printf("查询到 %d 个活跃用户\n", len(users))
}

// ==================== 性能优化示例 ====================

// PerformanceExample 性能优化示例
func PerformanceExample(db interface{}) {
	fmt.Println("=== 性能优化示例 ===")

	builder, err := sqlbuilder.New(db)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	// 1. 使用批量插入而不是循环单条插入
	batchData := make([]map[string]interface{}, 1000)
	for i := 0; i < 1000; i++ {
		batchData[i] = map[string]interface{}{
			"name":       fmt.Sprintf("用户%d", i),
			"email":      fmt.Sprintf("user%d@example.com", i),
			"age":        20 + (i % 50),
			"created_at": time.Now(),
		}
	}

	// 批量插入 - 高性能
	err = builder.WithContext(ctx).Table("users").BatchInsert(batchData)
	if err != nil {
		log.Printf("批量插入失败: %v", err)
		return
	}

	fmt.Println("批量插入1000条记录完成")

	// 2. 使用存在性检查而不是查询计数
	exists, err := builder.WithContext(ctx).
		Table("users").
		Where("email", "=", "user1@example.com").
		Exists()

	if err != nil {
		log.Printf("存在性检查失败: %v", err)
		return
	}

	if exists {
		fmt.Println("用户存在")
	}

	// 3. 使用分页而不是加载所有数据
	var users []UserExample
	err = builder.WithContext(ctx).
		Table("users").
		Where("status", "=", "active").
		OrderBy("id").
		Paginate(1, 20). // 第1页，每页20条
		Get(&users)

	if err != nil {
		log.Printf("分页查询失败: %v", err)
		return
	}

	fmt.Printf("分页查询到 %d 条记录\n", len(users))
}
