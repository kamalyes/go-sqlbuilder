/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-14 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-14 00:00:00
 * @FilePath: \go-sqlbuilder\examples\easycrud_examples.go
 * @Description: 简单CRUD使用示例 - 基于EnhancedBuilder的简化版本
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */

package examples

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/kamalyes/go-sqlbuilder"
)

// 简单CRUD使用示例 - 真正的简单易用！
func ExampleSimpleCRUD() {
	// 1. 连接数据库（这步都省不掉）
	db, err := sql.Open("mysql", "user:password@tcp(localhost:3306)/testdb")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 2. 创建简单CRUD操作器 - 基于EnhancedBuilder
	crud, err := sqlbuilder.NewSimple(db, "users")
	if err != nil {
		log.Fatal(err)
	}

	// ==================== 增删改查就这么简单！ ====================

	// ✅ 添加数据 - 1行代码搞定，自动处理时间戳
	user := map[string]interface{}{
		"name":  "张三",
		"email": "zhangsan@example.com",
		"age":   25,
	}

	err = crud.Add(user)
	if err != nil {
		fmt.Printf("添加失败: %v\n", err)
	} else {
		fmt.Println("✅ 添加成功！（自动添加了created_at和updated_at）")
	}

	// ✅ 查询数据 - 按ID查询
	userData, err := crud.Get(1)
	if err != nil {
		fmt.Printf("查询失败: %v\n", err)
	} else {
		fmt.Printf("✅ 查询成功: %+v\n", userData)
	}

	// ✅ 更新数据 - 按ID更新，自动处理updated_at
	updateData := map[string]interface{}{
		"age":  26,
		"name": "张三丰",
	}

	err = crud.Update(1, updateData)
	if err != nil {
		fmt.Printf("更新失败: %v\n", err)
	} else {
		fmt.Println("✅ 更新成功！（自动更新了updated_at）")
	}

	// ✅ 删除数据 - 软删除，自动设置deleted_at
	err = crud.Delete(1)
	if err != nil {
		fmt.Printf("删除失败: %v\n", err)
	} else {
		fmt.Println("✅ 软删除成功！（自动设置了deleted_at，数据仍在）")
	}

	// ✅ 获取列表 - 分页查询，自动过滤软删除数据
	userList, err := crud.List(1, 10) // 第1页，每页10条
	if err != nil {
		fmt.Printf("查询列表失败: %v\n", err)
	} else {
		fmt.Printf("✅ 查询到 %d 条数据（自动排除软删除数据）\n", len(userList))
	}

	// ✅ 统计数量 - 自动排除软删除数据
	count, err := crud.Count()
	if err != nil {
		fmt.Printf("统计失败: %v\n", err)
	} else {
		fmt.Printf("✅ 总共有 %d 条有效数据\n", count)
	}

	// ✅ 搜索功能
	searchResults, err := crud.Search("name", "张", 1, 10)
	if err != nil {
		fmt.Printf("搜索失败: %v\n", err)
	} else {
		fmt.Printf("✅ 搜索到 %d 条匹配数据\n", len(searchResults))
	}

	// ✅ 智能保存（有ID就更新，没ID就新增）
	newUser := map[string]interface{}{
		"name":  "李四",
		"email": "lisi@example.com",
		"age":   30,
	}
	err = crud.Save(newUser) // 没有ID，会执行新增
	if err != nil {
		fmt.Printf("智能保存失败: %v\n", err)
	} else {
		fmt.Println("✅ 智能保存成功（新增）！")
	}

	existingUser := map[string]interface{}{
		"id":   2,
		"name": "李四丰",
		"age":  31,
	}
	err = crud.Save(existingUser) // 有ID，会执行更新
	if err != nil {
		fmt.Printf("智能保存失败: %v\n", err)
	} else {
		fmt.Println("✅ 智能保存成功（更新）！")
	}
}

// 对比：复杂写法 vs 简单写法
func ComparisonExample() {
	db, _ := sql.Open("mysql", "user:password@tcp(localhost:3306)/testdb")
	defer db.Close()

	fmt.Println("=== ❌ 复杂的写法 ===")
	fmt.Println(`
	// 需要很多配置和复杂的操作
	builder, err := sqlbuilder.NewEnhanced(db)
	if err != nil {
		return err
	}
	
	builder.EnableSoftDelete(true).
		AddAuditFields("created_at", "updated_at", "deleted_at").
		AddHook(constant.HookEventBeforeCreate, func(ctx context.Context, data interface{}) error {
			// 手动处理时间戳...
			return nil
		})
	
	ctx := context.Background()
	options := &CreateOptions{SkipValidation: false, SkipHooks: false}
	result, err := builder.Table("users").SmartCreate(ctx, data, options)
	// ...还有复杂的错误处理
	`)

	fmt.Println("\n=== ✅ 现在的简单写法 ===")

	// 真正的简单用法
	crud, _ := sqlbuilder.NewSimple(db, "users")

	user := map[string]interface{}{
		"name":  "张三",
		"email": "zhangsan@example.com",
	}

	// 就一行代码！
	if err := crud.Add(user); err != nil {
		fmt.Printf("添加失败: %v\n", err)
	} else {
		fmt.Println("添加成功！")
	}

	fmt.Println("\n看到区别了吗？")
	fmt.Println("✅ 不需要复杂的配置")
	fmt.Println("✅ 不需要手动添加hooks")
	fmt.Println("✅ 不需要手动设置审计字段")
	fmt.Println("✅ 不需要关心context和options")
	fmt.Println("✅ 自动处理软删除")
	fmt.Println("✅ 自动处理时间戳")
	fmt.Println("✅ 错误信息是中文")
	fmt.Println("✅ 基于强大的EnhancedBuilder，高级功能随时可用")
}

// 真实项目使用示例
func RealWorldExample() {
	// 用户管理系统
	db, _ := sql.Open("mysql", "root:123456@tcp(localhost:3306)/myapp")
	defer db.Close()

	// 创建用户表操作器
	users, _ := sqlbuilder.NewSimple(db, "users")

	// 注册新用户 - 就这么简单！
	newUser := map[string]interface{}{
		"username": "john123",
		"email":    "john@email.com",
		"password": "hashedPassword",
		"status":   "active",
	}

	if err := users.Add(newUser); err != nil {
		fmt.Printf("用户注册失败: %v\n", err)
		return
	}
	fmt.Println("用户注册成功！")

	// 获取用户列表
	userList, err := users.List(1, 20)
	if err != nil {
		fmt.Printf("获取用户列表失败: %v\n", err)
	} else {
		fmt.Printf("用户列表：共 %d 个用户\n", len(userList))
		for _, user := range userList {
			fmt.Printf("- %s (%s)\n", user["username"], user["email"])
		}
	}

	// 搜索用户
	searchResults, err := users.Search("username", "john", 1, 10)
	if err != nil {
		fmt.Printf("搜索用户失败: %v\n", err)
	} else {
		fmt.Printf("搜索结果：找到 %d 个匹配用户\n", len(searchResults))
	}

	// 用户信息更新
	updateInfo := map[string]interface{}{
		"last_login":  "2025-11-14 10:30:00",
		"login_count": 1,
	}

	if err := users.Update(1, updateInfo); err != nil {
		fmt.Printf("更新用户信息失败: %v\n", err)
	} else {
		fmt.Println("用户信息更新成功！")
	}
}

// 使用指南
func UsageGuide() {
	fmt.Println("=== 🚀 简单CRUD使用指南 ===")
	fmt.Println()
	fmt.Println("1️⃣ 创建操作器:")
	fmt.Println("   crud, err := sqlbuilder.NewSimple(db, \"表名\")")
	fmt.Println()
	fmt.Println("2️⃣ 基础CRUD:")
	fmt.Println("   crud.Add(data)           // 添加（自动时间戳）")
	fmt.Println("   crud.Get(id)             // 按ID查询")
	fmt.Println("   crud.Update(id, data)    // 按ID更新（自动时间戳）")
	fmt.Println("   crud.Delete(id)          // 软删除（自动时间戳）")
	fmt.Println()
	fmt.Println("3️⃣ 列表和搜索:")
	fmt.Println("   crud.List(page, size)    // 分页列表（自动过滤软删除）")
	fmt.Println("   crud.Search(field, keyword, page, size)  // 搜索")
	fmt.Println("   crud.Count()             // 统计（排除软删除）")
	fmt.Println()
	fmt.Println("4️⃣ 智能操作:")
	fmt.Println("   crud.Save(data)          // 智能保存（有ID更新，无ID新增）")
	fmt.Println()
	fmt.Println("🎯 特色功能:")
	fmt.Println("   ✅ 基于强大的EnhancedBuilder")
	fmt.Println("   ✅ 自动时间戳（created_at, updated_at）")
	fmt.Println("   ✅ 自动软删除（deleted_at）")
	fmt.Println("   ✅ 中文友好错误信息")
	fmt.Println("   ✅ 零配置直接使用")
	fmt.Println("   ✅ 高级功能随时可用")
	fmt.Println()
	fmt.Println("💡 如需高级功能:")
	fmt.Println("   可以直接使用 EnhancedBuilder 的 SmartCreate、SmartUpdate 等方法")
	fmt.Println("   享受 hooks、validation、乐观锁等高级特性")
}
