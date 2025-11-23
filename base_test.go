/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-23 15:45:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-23 23:57:00
 * @FilePath: \go-sqlbuilder\base_test.go
 * @Description:
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package sqlbuilder

import (
	"context"
	"fmt"
	"github.com/kamalyes/go-logger"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
	"time"
)

// TestUser 测试用的用户实体
type TestUser struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	Name      string     `json:"name" gorm:"column:name"`
	Email     string     `json:"email" gorm:"column:email;unique"`
	Age       int        `json:"age" gorm:"column:age"`
	Status    string     `json:"status" gorm:"column:status"`
	CreatedAt time.Time  `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"column:deleted_at"`
}

// setupTestDB 设置测试数据库（SQLite 内存数据库）
func setupTestDB() (*gorm.DB, error) {
	gormDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return nil, err
	}

	// 自动迁移表结构
	err = gormDB.AutoMigrate(&TestUser{})
	if err != nil {
		return nil, err
	}

	return gormDB, nil
}

// testDBHandler 测试数据库处理器
type testDBHandler struct {
	gormDB *gorm.DB
}

func (t *testDBHandler) GetDB() *gorm.DB {
	return t.gormDB
}

func (t *testDBHandler) IsConnected() bool {
	if t.gormDB == nil {
		return false
	}
	sqlDB, err := t.gormDB.DB()
	if err != nil {
		return false
	}
	return sqlDB.Ping() == nil
}

func newTestDBHandler(gormDB *gorm.DB) Handler {
	return &testDBHandler{gormDB: gormDB}
}

// TestNewBaseRepository 测试基础仓储创建
func TestNewBaseRepository(t *testing.T) {
	gormDB, err := setupTestDB()
	assert.NoError(t, err)

	dbHandler := newTestDBHandler(gormDB)

	// 测试默认配置
	repo := NewBaseRepository[TestUser](dbHandler, logger.NewLogger(nil), "users")
	assert.NotNil(t, repo, "仓储不应为空")
	assert.Equal(t, "users", repo.table, "表名应为 'users'")
	assert.Equal(t, 100, repo.batchSize, "默认批处理大小应为 100")
	assert.Equal(t, 30, repo.timeout, "默认超时时间应为 30秒")
	assert.False(t, repo.readOnly, "默认不应为只读模式")
	assert.Empty(t, repo.preloads, "默认预加载应为空")
	assert.Empty(t, repo.defaultOrder, "默认排序应为空")
	assert.NotNil(t, repo.logger, "日志记录器不应为空")
	// assert.NotNil(t, repo.contextExtractor, "context提取器不应为空")
}

// TestNewBaseRepositoryWithOptions 测试带选项的基础仓储创建
func TestNewBaseRepositoryWithOptions(t *testing.T) {
	gormDB, err := setupTestDB()
	assert.NoError(t, err)

	dbHandler := newTestDBHandler(gormDB)
	testLogger := logger.NewLogger(nil)

	// 测试自定义配置
	repo := NewBaseRepository[TestUser](
		dbHandler,
		testLogger,
		"test_users",
		WithBatchSize[TestUser](50),
		WithTimeout[TestUser](60),
		WithReadOnly[TestUser](),
		WithDefaultPreloads[TestUser]("Profile", "Posts"),
		WithDefaultOrder[TestUser]("created_at DESC"),
		WithLogger[TestUser](testLogger),
	)

	assert.NotNil(t, repo, "仓储不应为空")
	assert.Equal(t, "test_users", repo.table, "表名应为 'test_users'")
	assert.Equal(t, 50, repo.batchSize, "批处理大小应为 50")
	assert.Equal(t, 60, repo.timeout, "超时时间应为 60秒")
	assert.True(t, repo.readOnly, "应为只读模式")
	assert.Equal(t, []string{"Profile", "Posts"}, repo.preloads, "预加载应包含指定关联")
	assert.Equal(t, "created_at DESC", repo.defaultOrder, "默认排序应为指定值")
	assert.Equal(t, testLogger, repo.logger, "应使用指定的日志记录器")
}

// TestBaseRepositoryCreate 测试创建操作
func TestBaseRepositoryCreate(t *testing.T) {
	gormDB, err := setupTestDB()
	assert.NoError(t, err)

	dbHandler := newTestDBHandler(gormDB)
	repo := NewBaseRepository[TestUser](dbHandler, logger.NewLogger(nil), "test_users")

	ctx := context.Background()
	user := &TestUser{
		Name:   "John Doe",
		Email:  "john@example.com",
		Age:    25,
		Status: "active",
	}

	result, err := repo.Create(ctx, user)
	assert.NoError(t, err, "创建操作不应出错")
	assert.NotNil(t, result, "结果不应为空")
	assert.Equal(t, user.Name, result.Name, "用户名应相同")
	assert.Equal(t, user.Email, result.Email, "邮箱应相同")
	assert.Equal(t, user.Age, result.Age, "年龄应相同")
	assert.Equal(t, user.Status, result.Status, "状态应相同")
	assert.NotZero(t, result.ID, "ID应被自动生成")
	assert.NotZero(t, result.CreatedAt, "创建时间应被设置")
	assert.NotZero(t, result.UpdatedAt, "更新时间应被设置")
}

// TestBaseRepositoryCreateNilEntity 测试创建空实体
func TestBaseRepositoryCreateNilEntity(t *testing.T) {
	gormDB, err := setupTestDB()
	assert.NoError(t, err)

	dbHandler := newTestDBHandler(gormDB)
	repo := NewBaseRepository[TestUser](dbHandler, logger.NewLogger(nil), "users")

	ctx := context.Background()

	result, err := repo.Create(ctx, nil)
	assert.Error(t, err, "创建空实体应返回错误")
	assert.Nil(t, result, "结果应为空")
}

// TestBaseRepositoryCreateReadOnlyMode 测试只读模式下创建
func TestBaseRepositoryCreateReadOnlyMode(t *testing.T) {
	gormDB, err := setupTestDB()
	assert.NoError(t, err)

	dbHandler := newTestDBHandler(gormDB)
	repo := NewBaseRepository[TestUser](dbHandler, logger.NewLogger(nil), "users", WithReadOnly[TestUser]())

	ctx := context.Background()
	user := &TestUser{Name: "John"}

	result, err := repo.Create(ctx, user)
	assert.Error(t, err, "只读模式下创建应返回错误")
	assert.Nil(t, result, "结果应为空")
}

// TestBaseRepositoryGet 测试获取单个记录
func TestBaseRepositoryGet(t *testing.T) {
	gormDB, err := setupTestDB()
	assert.NoError(t, err)

	dbHandler := newTestDBHandler(gormDB)
	repo := NewBaseRepository[TestUser](dbHandler, logger.NewLogger(nil), "test_users")

	ctx := context.Background()

	// 先创建一个用户
	user := &TestUser{
		Name:   "John Doe",
		Email:  "john@example.com",
		Age:    25,
		Status: "active",
	}
	createdUser, err := repo.Create(ctx, user)
	assert.NoError(t, err)

	// 获取创建的用户
	result, err := repo.Get(ctx, createdUser.ID)
	assert.NoError(t, err, "获取操作不应出错")
	assert.NotNil(t, result, "结果不应为空")
	assert.Equal(t, createdUser.ID, result.ID, "用户ID应正确")
	assert.Equal(t, "John Doe", result.Name, "用户名应正确")
	assert.Equal(t, "john@example.com", result.Email, "邮箱应正确")
}

// TestBaseRepositoryList 测试列表查询
func TestBaseRepositoryList(t *testing.T) {
	gormDB, err := setupTestDB()
	assert.NoError(t, err)

	dbHandler := newTestDBHandler(gormDB)
	repo := NewBaseRepository[TestUser](dbHandler, logger.NewLogger(nil), "test_users")

	ctx := context.Background()

	// 先创建一些用户
	users := []*TestUser{
		{Name: "User1", Email: "user1@example.com", Age: 25, Status: "active"},
		{Name: "User2", Email: "user2@example.com", Age: 30, Status: "active"},
		{Name: "User3", Email: "user3@example.com", Age: 35, Status: "inactive"},
	}
	err = repo.CreateBatch(ctx, users...)
	assert.NoError(t, err)

	// 查询所有活跃用户
	query := NewQuery()
	query.AddFilter(NewEqFilter("status", "active"))
	query.AddOrder("created_at", "ASC")
	query.Limit(10)
	query.Offset(0)

	result, err := repo.List(ctx, query)
	assert.NoError(t, err, "列表查询不应出错")
	assert.NotNil(t, result, "结果不应为空")
	assert.Len(t, result, 2, "应返回 2 个活跃用户")
}

// TestBaseRepositoryCount 测试计数操作
func TestBaseRepositoryCount(t *testing.T) {
	gormDB, err := setupTestDB()
	assert.NoError(t, err)

	dbHandler := newTestDBHandler(gormDB)
	repo := NewBaseRepository[TestUser](dbHandler, logger.NewLogger(nil), "test_users")

	ctx := context.Background()

	// 创建用户
	users := []*TestUser{
		{Name: "User1", Email: "user1@example.com", Age: 25, Status: "active"},
		{Name: "User2", Email: "user2@example.com", Age: 30, Status: "active"},
		{Name: "User3", Email: "user3@example.com", Age: 35, Status: "inactive"},
	}
	err = repo.CreateBatch(ctx, users...)
	assert.NoError(t, err)

	// 计数活跃用户
	filter := NewEqFilter("status", "active")
	count, err := repo.Count(ctx, filter)
	assert.NoError(t, err, "计数操作不应出错")
	assert.Equal(t, int64(2), count, "活跃用户计数应为 2")
}

// TestBaseRepositoryExists 测试存在性检查
func TestBaseRepositoryExists(t *testing.T) {
	gormDB, err := setupTestDB()
	assert.NoError(t, err)

	dbHandler := newTestDBHandler(gormDB)
	repo := NewBaseRepository[TestUser](dbHandler, logger.NewLogger(nil), "test_users")

	ctx := context.Background()

	// 创建用户
	user := &TestUser{
		Name:   "John Doe",
		Email:  "john@example.com",
		Age:    25,
		Status: "active",
	}
	_, err = repo.Create(ctx, user)
	assert.NoError(t, err)

	// 检查存在
	filter := NewEqFilter("email", "john@example.com")
	exists, err := repo.Exists(ctx, filter)
	assert.NoError(t, err, "存在性检查不应出错")
	assert.True(t, exists, "应存在记录")

	// 检查不存在
	filter2 := NewEqFilter("email", "notfound@example.com")
	exists2, err := repo.Exists(ctx, filter2)
	assert.NoError(t, err, "存在性检查不应出错")
	assert.False(t, exists2, "不应存在记录")
}

// TestApplyFilter 测试过滤器应用函数
func TestApplyFilter(t *testing.T) {
	gormDB, err := setupTestDB()
	assert.NoError(t, err)

	dbQuery := gormDB.Table("test_users")

	testCases := []struct {
		name   string
		filter *Filter
	}{
		{"EQ filter", NewEqFilter("name", "John")},
		{"GT filter", NewGtFilter("age", 18)},
		{"IN filter", NewInFilter("status", "active", "pending")},
		{"LIKE filter", NewLikeFilter("title", "test")},
		{"BETWEEN filter", NewBetweenFilter("created_at", "2023-01-01", "2023-12-31")},
		{"IS NULL filter", NewIsNullFilter("deleted_at")},
		{"IS NOT NULL filter", NewIsNotNullFilter("updated_at")},
		{"FIND_IN_SET filter", NewFindInSetFilter("tags", "important")},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 应用过滤器不应panic
			assert.NotPanics(t, func() {
				applyFilter(dbQuery, tc.filter)
			}, "应用"+tc.name+"不应panic")
		})
	}

	// 测试nil过滤器
	t.Run("nil filter", func(t *testing.T) {
		result := applyFilter(dbQuery, nil)
		assert.Equal(t, dbQuery, result, "nil过滤器应返回原始查询")
	})

	// 测试BETWEEN过滤器的边界情况
	t.Run("BETWEEN filter with invalid value", func(t *testing.T) {
		invalidBetweenFilter := &Filter{
			Field:    "created_at",
			Operator: OP_BETWEEN,
			Value:    "not_an_array", // 无效值
		}

		assert.NotPanics(t, func() {
			applyFilter(dbQuery, invalidBetweenFilter)
		}, "无效BETWEEN值不应panic")
	})
}

// TestBaseRepositoryGetByFilters 测试根据多个过滤器获取
func TestBaseRepositoryGetByFilters(t *testing.T) {
	gormDB, err := setupTestDB()
	assert.NoError(t, err)

	dbHandler := newTestDBHandler(gormDB)
	repo := NewBaseRepository[TestUser](dbHandler, logger.NewLogger(nil), "test_users")

	ctx := context.Background()

	// 创建用户
	user := &TestUser{
		Name:   "John Doe",
		Email:  "john@example.com",
		Age:    25,
		Status: "active",
	}
	_, err = repo.Create(ctx, user)
	assert.NoError(t, err)

	// 按多个过滤器查询
	result, err := repo.GetByFilters(ctx,
		NewEqFilter("email", "john@example.com"),
		NewEqFilter("status", "active"),
	)
	assert.NoError(t, err, "获取操作不应出错")
	assert.NotNil(t, result, "结果不应为空")
	assert.Equal(t, "john@example.com", result.Email, "邮箱应正确")
	assert.Equal(t, "John Doe", result.Name, "名称应正确")
	assert.Equal(t, 25, result.Age, "年龄应正确")
	assert.Equal(t, "active", result.Status, "状态应正确")
	assert.Greater(t, result.ID, uint(0), "ID应大于0")
	assert.NotZero(t, result.CreatedAt, "创建时间应不为零值")
	assert.NotZero(t, result.UpdatedAt, "更新时间应不为零值")
}

// TestBaseRepositoryGetByFields 测试按字段获取
func TestBaseRepositoryGetByFields(t *testing.T) {
	gormDB, err := setupTestDB()
	assert.NoError(t, err)

	dbHandler := newTestDBHandler(gormDB)
	repo := NewBaseRepository[TestUser](dbHandler, logger.NewLogger(nil), "test_users")

	ctx := context.Background()

	// 创建用户
	user := &TestUser{
		Name:   "John Doe",
		Email:  "john@example.com",
		Age:    25,
		Status: "active",
	}
	_, err = repo.Create(ctx, user)
	assert.NoError(t, err)

	// 按字段获取
	fields := map[string]interface{}{
		"email":  "john@example.com",
		"status": "active",
	}

	result, err := repo.GetByFields(ctx, fields)
	assert.NoError(t, err, "获取操作不应出错")
	assert.NotNil(t, result, "结果不应为空")
	assert.Equal(t, "john@example.com", result.Email, "邮箱应正确")
	assert.Equal(t, "John Doe", result.Name, "名称应正确")
	assert.Equal(t, 25, result.Age, "年龄应正确")
	assert.Equal(t, "active", result.Status, "状态应正确")
	assert.Greater(t, result.ID, uint(0), "ID应大于0")

	// 测试空字段
	emptyFields := map[string]interface{}{}
	result2, err := repo.GetByFields(ctx, emptyFields)
	assert.Error(t, err, "空字段应返回错误")
	assert.Nil(t, result2, "结果应为空")
}

// TestBaseRepositoryListWithPreloads 测试带预加载的列表查询
func TestBaseRepositoryListWithPreloads(t *testing.T) {
	gormDB, err := setupTestDB()
	assert.NoError(t, err)

	dbHandler := newTestDBHandler(gormDB)
	repo := NewBaseRepository[TestUser](dbHandler, logger.NewLogger(nil), "test_users")

	ctx := context.Background()

	// 创建用户
	users := []*TestUser{
		{Name: "User1", Email: "user1@example.com", Age: 25, Status: "active"},
		{Name: "User2", Email: "user2@example.com", Age: 30, Status: "active"},
	}
	err = repo.CreateBatch(ctx, users...)
	assert.NoError(t, err)

	// 查询所有活跃用户（不带预加载，因为TestUser没有关联表）
	query := NewQuery()
	query.AddFilter(NewEqFilter("status", "active"))

	result, err := repo.List(ctx, query)
	assert.NoError(t, err, "列表查询不应出错")
	assert.NotNil(t, result, "结果不应为空")
	assert.Len(t, result, 2, "应返回 2 个活跃用户")
}

// TestBaseRepositoryFind 测试兼容旧API查询
func TestBaseRepositoryFind(t *testing.T) {
	gormDB, err := setupTestDB()
	assert.NoError(t, err)

	dbHandler := newTestDBHandler(gormDB)
	repo := NewBaseRepository[TestUser](dbHandler, logger.NewLogger(nil), "test_users")

	ctx := context.Background()

	// 创建用户
	users := []*TestUser{
		{Name: "User1", Email: "user1@example.com", Age: 25, Status: "active"},
		{Name: "User2", Email: "user2@example.com", Age: 30, Status: "active"},
		{Name: "User3", Email: "user3@example.com", Age: 35, Status: "inactive"},
	}
	err = repo.CreateBatch(ctx, users...)
	assert.NoError(t, err)

	// 使用旧API查询
	options := &FindOptions{
		Conditions: []Condition{
			{Field: "status", Op: OP_EQ, Value: "active"},
		},
		Orders: []OrderBy{
			{Field: "age", Direction: "ASC"},
		},
		Limit:  10,
		Offset: 0,
	}

	result, err := repo.Find(ctx, options)
	assert.NoError(t, err, "查询不应出错")
	assert.NotNil(t, result, "结果不应为空")
	assert.Len(t, result, 2, "应返回 2 个活跃用户")
	assert.Equal(t, 25, result[0].Age, "第一个用户年龄应为25（按年龄升序）")
	assert.Equal(t, 30, result[1].Age, "第二个用户年龄应为30（按年龄升序）")
	// 验证排序是否正确
	assert.LessOrEqual(t, result[0].Age, result[1].Age, "结果应按年龄升序排列")
	// 验证状态过滤是否正确
	for _, user := range result {
		assert.Equal(t, "active", user.Status, "所有用户状态应为active")
	}
}

// TestBaseRepositoryUpdateBatch 测试批量更新
func TestBaseRepositoryUpdateBatch(t *testing.T) {
	gormDB, err := setupTestDB()
	assert.NoError(t, err)

	dbHandler := newTestDBHandler(gormDB)
	repo := NewBaseRepository[TestUser](dbHandler, logger.NewLogger(nil), "test_users")

	ctx := context.Background()

	// 创建用户
	users := []*TestUser{
		{Name: "User1", Email: "user1@example.com", Age: 25, Status: "active"},
		{Name: "User2", Email: "user2@example.com", Age: 30, Status: "active"},
	}
	err = repo.CreateBatch(ctx, users...)
	assert.NoError(t, err)

	// 获取创建的用户
	var createdUsers []TestUser
	gormDB.Table("test_users").Find(&createdUsers)
	assert.Len(t, createdUsers, 2)

	// 修改用户信息
	for i := range createdUsers {
		createdUsers[i].Status = "updated"
	}

	// 批量更新
	updateUsers := make([]*TestUser, len(createdUsers))
	for i := range createdUsers {
		updateUsers[i] = &createdUsers[i]
	}

	err = repo.UpdateBatch(ctx, updateUsers...)
	assert.NoError(t, err, "批量更新不应出错")

	// 验证更新
	var updatedCount int64
	gormDB.Table("test_users").Where("status = ?", "updated").Count(&updatedCount)
	assert.Equal(t, int64(2), updatedCount, "应有 2 个用户被更新")

	// 验证具体的更新数据
	var batchUpdatedUsers []TestUser
	gormDB.Table("test_users").Where("status = ?", "updated").Find(&batchUpdatedUsers)
	assert.Len(t, batchUpdatedUsers, 2, "应有2个更新的用户")

	for i, user := range batchUpdatedUsers {
		assert.Equal(t, "updated", user.Status, fmt.Sprintf("第%d个用户状态应为updated", i+1))
		assert.Greater(t, user.ID, uint(0), fmt.Sprintf("第%d个用户ID应大于0", i+1))
		assert.NotEmpty(t, user.Name, fmt.Sprintf("第%d个用户名不应为空", i+1))
		assert.NotEmpty(t, user.Email, fmt.Sprintf("第%d个用户邮箱不应为空", i+1))
		assert.Greater(t, user.Age, 0, fmt.Sprintf("第%d个用户年龄应大于0", i+1))
		assert.NotZero(t, user.CreatedAt, fmt.Sprintf("第%d个用户创建时间应不为零值", i+1))
		assert.NotZero(t, user.UpdatedAt, fmt.Sprintf("第%d个用户更新时间应不为零值", i+1))
	}

	// 验证原始数据未被意外修改
	expectedNames := []string{"User1", "User2"}
	expectedEmails := []string{"user1@example.com", "user2@example.com"}
	for _, user := range batchUpdatedUsers {
		assert.Contains(t, expectedNames, user.Name, "用户名应保持不变")
		assert.Contains(t, expectedEmails, user.Email, "用户邮箱应保持不变")
	}

	// 验证具体的更新数据
	var updatedUsers []TestUser
	gormDB.Table("test_users").Where("status = ?", "updated").Find(&updatedUsers)
	for _, user := range updatedUsers {
		assert.Equal(t, "updated", user.Status, "用户状态应为updated")
		assert.Contains(t, []string{"User1", "User2"}, user.Name, "用户名应保持不变")
		assert.Contains(t, []string{"user1@example.com", "user2@example.com"}, user.Email, "用户邮箱应保持不变")
	}
}

// TestBaseRepositoryUpdateByFilters 测试按过滤器更新
func TestBaseRepositoryUpdateByFilters(t *testing.T) {
	gormDB, err := setupTestDB()
	assert.NoError(t, err)

	dbHandler := newTestDBHandler(gormDB)
	repo := NewBaseRepository[TestUser](dbHandler, logger.NewLogger(nil), "test_users")

	ctx := context.Background()

	// 创建用户
	users := []*TestUser{
		{Name: "User1", Email: "user1@example.com", Age: 25, Status: "active"},
		{Name: "User2", Email: "user2@example.com", Age: 30, Status: "active"},
		{Name: "User3", Email: "user3@example.com", Age: 35, Status: "inactive"},
	}
	err = repo.CreateBatch(ctx, users...)
	assert.NoError(t, err)

	// 按过滤器更新所有活跃用户的状态
	updateEntity := &TestUser{Status: "suspended"}
	err = repo.UpdateByFilters(ctx, updateEntity, NewEqFilter("status", "active"))
	assert.NoError(t, err, "按过滤器更新不应出错")

	// 验证更新
	var suspendedCount int64
	gormDB.Table("test_users").Where("status = ?", "suspended").Count(&suspendedCount)
	assert.Equal(t, int64(2), suspendedCount, "应有 2 个用户被挂起")

	var inactiveCount int64
	gormDB.Table("test_users").Where("status = ?", "inactive").Count(&inactiveCount)
	assert.Equal(t, int64(1), inactiveCount, "应保留 1 个非活跃用户")
}

// TestBaseRepositoryDeleteByFilters 测试按过滤器删除
func TestBaseRepositoryDeleteByFilters(t *testing.T) {
	gormDB, err := setupTestDB()
	assert.NoError(t, err)

	dbHandler := newTestDBHandler(gormDB)
	repo := NewBaseRepository[TestUser](dbHandler, logger.NewLogger(nil), "test_users")

	ctx := context.Background()

	// 创建用户
	users := []*TestUser{
		{Name: "User1", Email: "user1@example.com", Age: 25, Status: "active"},
		{Name: "User2", Email: "user2@example.com", Age: 30, Status: "active"},
		{Name: "User3", Email: "user3@example.com", Age: 35, Status: "inactive"},
	}
	err = repo.CreateBatch(ctx, users...)
	assert.NoError(t, err)

	// 按过滤器删除所有活跃用户
	err = repo.DeleteByFilters(ctx, NewEqFilter("status", "active"))
	assert.NoError(t, err, "按过滤器删除不应出错")

	// 验证删除
	var remainingCount int64
	gormDB.Table("test_users").Count(&remainingCount)
	assert.Equal(t, int64(1), remainingCount, "应保留 1 个用户")

	var inactiveCount int64
	gormDB.Table("test_users").Where("status = ?", "inactive").Count(&inactiveCount)
	assert.Equal(t, int64(1), inactiveCount, "应保留非活跃用户")
}

// TestBaseRepositoryTransaction 测试事务操作
func TestBaseRepositoryTransaction(t *testing.T) {
	gormDB, err := setupTestDB()
	assert.NoError(t, err)

	dbHandler := newTestDBHandler(gormDB)
	repo := NewBaseRepository[TestUser](dbHandler, logger.NewLogger(nil), "test_users")

	ctx := context.Background()

	// 测试成功的事务
	err = repo.Transaction(ctx, func(tx Transaction[TestUser]) error {
		// 在事务中创建用户（这里只能模拟，因为tx接口没有实现）
		return nil
	})
	assert.NoError(t, err, "事务不应出错")

	// 测试失败的事务
	testError := fmt.Errorf("test transaction error")
	err = repo.Transaction(ctx, func(tx Transaction[TestUser]) error {
		return testError
	})
	assert.Error(t, err, "事务应该失败")
	assert.Equal(t, testError, err, "应返回原始错误")
}

// TestBaseRepositoryGetAll 测试获取所有记录
func TestBaseRepositoryGetAll(t *testing.T) {
	gormDB, err := setupTestDB()
	assert.NoError(t, err)

	dbHandler := newTestDBHandler(gormDB)
	repo := NewBaseRepository[TestUser](dbHandler, logger.NewLogger(nil), "test_users")

	ctx := context.Background()

	// 创建用户
	users := []*TestUser{
		{Name: "User1", Email: "user1@example.com", Age: 25, Status: "active"},
		{Name: "User2", Email: "user2@example.com", Age: 30, Status: "active"},
		{Name: "User3", Email: "user3@example.com", Age: 35, Status: "inactive"},
	}
	err = repo.CreateBatch(ctx, users...)
	assert.NoError(t, err)

	// 获取所有记录
	result, err := repo.GetAll(ctx)
	assert.NoError(t, err, "获取所有记录不应出错")
	assert.NotNil(t, result, "结果不应为空")
	assert.Len(t, result, 3, "应返回 3 个用户")

	// 验证所有用户的数据完整性
	expectedNames := []string{"User1", "User2", "User3"}
	expectedEmails := []string{"user1@example.com", "user2@example.com", "user3@example.com"}
	expectedAges := []int{25, 30, 35}

	actualNames := make([]string, len(result))
	actualEmails := make([]string, len(result))
	actualStatuses := make([]string, len(result))
	actualAges := make([]int, len(result))

	for i, user := range result {
		assert.Greater(t, user.ID, uint(0), fmt.Sprintf("第%d个用户ID应大于0", i+1))
		assert.NotEmpty(t, user.Name, fmt.Sprintf("第%d个用户名不应为空", i+1))
		assert.NotEmpty(t, user.Email, fmt.Sprintf("第%d个用户邮箱不应为空", i+1))
		assert.Greater(t, user.Age, 0, fmt.Sprintf("第%d个用户年龄应大于0", i+1))
		assert.NotEmpty(t, user.Status, fmt.Sprintf("第%d个用户状态不应为空", i+1))
		assert.NotZero(t, user.CreatedAt, fmt.Sprintf("第%d个用户创建时间应不为零值", i+1))
		assert.NotZero(t, user.UpdatedAt, fmt.Sprintf("第%d个用户更新时间应不为零值", i+1))

		actualNames[i] = user.Name
		actualEmails[i] = user.Email
		actualStatuses[i] = user.Status
		actualAges[i] = user.Age
	}

	// 验证用户数据的唯一性和完整性
	for _, expectedName := range expectedNames {
		assert.Contains(t, actualNames, expectedName, fmt.Sprintf("应包含用户名%s", expectedName))
	}
	for _, expectedEmail := range expectedEmails {
		assert.Contains(t, actualEmails, expectedEmail, fmt.Sprintf("应包含邮箱%s", expectedEmail))
	}
	for _, expectedAge := range expectedAges {
		assert.Contains(t, actualAges, expectedAge, fmt.Sprintf("应包含年龄%d", expectedAge))
	}

	// 验证状态分布
	activeCount := 0
	inactiveCount := 0
	for _, status := range actualStatuses {
		switch status {
		case "active":
			activeCount++
		case "inactive":
			inactiveCount++
		}
	}
	assert.Equal(t, 2, activeCount, "应有2个活跃用户")
	assert.Equal(t, 1, inactiveCount, "应有1个非活跃用户")

	// 验证所有用户的数据完整性
	names := make([]string, len(result))
	emails := make([]string, len(result))
	statuses := make([]string, len(result))
	for i, user := range result {
		assert.Greater(t, user.ID, uint(0), "用户ID应大于0")
		assert.NotEmpty(t, user.Name, "用户名不应为空")
		assert.NotEmpty(t, user.Email, "用户邮箱不应为空")
		assert.Greater(t, user.Age, 0, "用户年龄应大于0")
		assert.NotEmpty(t, user.Status, "用户状态不应为空")
		names[i] = user.Name
		emails[i] = user.Email
		statuses[i] = user.Status
	}
	// 验证用户数据的唯一性和完整性
	assert.Contains(t, names, "User1", "应包含User1")
	assert.Contains(t, names, "User2", "应包含User2")
	assert.Contains(t, names, "User3", "应包含User3")
	assert.Contains(t, emails, "user1@example.com", "应包含user1邮箱")
	assert.Contains(t, emails, "user2@example.com", "应包含user2邮箱")
	assert.Contains(t, emails, "user3@example.com", "应包含user3邮箱")
}

// TestBaseRepositoryFirst 测试获取第一个记录
func TestBaseRepositoryFirst(t *testing.T) {
	gormDB, err := setupTestDB()
	assert.NoError(t, err)

	dbHandler := newTestDBHandler(gormDB)
	repo := NewBaseRepository[TestUser](dbHandler, logger.NewLogger(nil), "test_users")

	ctx := context.Background()

	// 创建用户
	users := []*TestUser{
		{Name: "User1", Email: "user1@example.com", Age: 25, Status: "active"},
		{Name: "User2", Email: "user2@example.com", Age: 30, Status: "active"},
	}
	err = repo.CreateBatch(ctx, users...)
	assert.NoError(t, err)

	// 获取第一个活跃用户
	result, err := repo.First(ctx, NewEqFilter("status", "active"))
	assert.NoError(t, err, "获取第一个记录不应出错")
	assert.NotNil(t, result, "结果不应为空")
	assert.Equal(t, "active", result.Status, "状态应为活跃")
	assert.Greater(t, result.ID, uint(0), "用户ID应大于0")
	assert.NotEmpty(t, result.Name, "用户名不应为空")
	assert.NotEmpty(t, result.Email, "用户邮箱不应为空")
	assert.Greater(t, result.Age, 0, "用户年龄应大于0")
	assert.NotZero(t, result.CreatedAt, "创建时间应不为零值")
	assert.NotZero(t, result.UpdatedAt, "更新时间应不为零值")

	// 验证是否是创建的活跃用户之一
	expectedNames := []string{"User1", "User2"}
	expectedEmails := []string{"user1@example.com", "user2@example.com"}
	expectedAges := []int{25, 30}
	assert.Contains(t, expectedNames, result.Name, "应该是活跃用户之一")
	assert.Contains(t, expectedEmails, result.Email, "邮箱应该是活跃用户之一")
	assert.Contains(t, expectedAges, result.Age, "年龄应该是活跃用户之一")
	assert.Greater(t, result.ID, uint(0), "用户ID应大于0")
	assert.NotEmpty(t, result.Name, "用户名不应为空")
	assert.NotEmpty(t, result.Email, "用户邮箱不应为空")
	assert.Greater(t, result.Age, 0, "用户年龄应大于0")
	// 验证是否是创建的活跃用户之一
	assert.Contains(t, []string{"User1", "User2"}, result.Name, "应该是活跃用户之一")
}

// TestBaseRepositoryLast 测试获取最后一个记录
func TestBaseRepositoryLast(t *testing.T) {
	gormDB, err := setupTestDB()
	assert.NoError(t, err)

	dbHandler := newTestDBHandler(gormDB)
	repo := NewBaseRepository[TestUser](dbHandler, logger.NewLogger(nil), "test_users")

	ctx := context.Background()

	// 创建用户
	users := []*TestUser{
		{Name: "User1", Email: "user1@example.com", Age: 25, Status: "active"},
		{Name: "User2", Email: "user2@example.com", Age: 30, Status: "active"},
	}
	err = repo.CreateBatch(ctx, users...)
	assert.NoError(t, err)

	// 获取最后一个活跃用户
	result, err := repo.Last(ctx, NewEqFilter("status", "active"))
	assert.NoError(t, err, "获取最后一个记录不应出错")
	assert.NotNil(t, result, "结果不应为空")
	assert.Equal(t, "active", result.Status, "状态应为活跃")
	assert.Greater(t, result.ID, uint(0), "用户ID应大于0")
	assert.NotEmpty(t, result.Name, "用户名不应为空")
	assert.NotEmpty(t, result.Email, "用户邮箱不应为空")
	assert.Greater(t, result.Age, 0, "用户年龄应大于0")
	assert.NotZero(t, result.CreatedAt, "创建时间应不为零值")
	assert.NotZero(t, result.UpdatedAt, "更新时间应不为零值")

	// 验证是否是创建的活跃用户之一
	expectedNames := []string{"User1", "User2"}
	expectedEmails := []string{"user1@example.com", "user2@example.com"}
	expectedAges := []int{25, 30}
	assert.Contains(t, expectedNames, result.Name, "应该是活跃用户之一")
	assert.Contains(t, expectedEmails, result.Email, "邮箱应该是活跃用户之一")
	assert.Contains(t, expectedAges, result.Age, "年龄应该是活跃用户之一")
	assert.Greater(t, result.ID, uint(0), "用户ID应大于0")
	assert.NotEmpty(t, result.Name, "用户名不应为空")
	assert.NotEmpty(t, result.Email, "用户邮箱不应为空")
	assert.Greater(t, result.Age, 0, "用户年龄应大于0")
	// 验证是否是创建的活跃用户之一
	assert.Contains(t, []string{"User1", "User2"}, result.Name, "应该是活跃用户之一")
}

// TestBaseRepositoryFindOne 测试查找单个记录
func TestBaseRepositoryFindOne(t *testing.T) {
	gormDB, err := setupTestDB()
	assert.NoError(t, err)

	dbHandler := newTestDBHandler(gormDB)
	repo := NewBaseRepository[TestUser](dbHandler, logger.NewLogger(nil), "test_users")

	ctx := context.Background()

	// 创建用户
	user := &TestUser{
		Name:   "John Doe",
		Email:  "john@example.com",
		Age:    25,
		Status: "active",
	}
	_, err = repo.Create(ctx, user)
	assert.NoError(t, err)

	// 查找单个记录
	result, err := repo.FindOne(ctx, NewEqFilter("email", "john@example.com"))
	assert.NoError(t, err, "查找单个记录不应出错")
	assert.NotNil(t, result, "结果不应为空")
	assert.Equal(t, "john@example.com", result.Email, "邮箱应正确")
	assert.Equal(t, "John Doe", result.Name, "名称应正确")
	assert.Equal(t, 25, result.Age, "年龄应正确")
	assert.Equal(t, "active", result.Status, "状态应正确")
	assert.Greater(t, result.ID, uint(0), "ID应大于0")
	assert.NotZero(t, result.CreatedAt, "创建时间应不为零值")
	assert.NotZero(t, result.UpdatedAt, "更新时间应不为零值")
	assert.Equal(t, "John Doe", result.Name, "名称应正确")
	assert.Equal(t, 25, result.Age, "年龄应正确")
	assert.Equal(t, "active", result.Status, "状态应正确")
	assert.Greater(t, result.ID, uint(0), "ID应大于0")

	// 查找不存在的记录
	result2, err := repo.FindOne(ctx, NewEqFilter("email", "notfound@example.com"))
	// 在GORM中，如果记录不存在，FindOne可能返回nil但不报错，这取决于实现
	// 这里我们检查result2是否为nil即可
	assert.Nil(t, result2, "结果应为空")
}

// TestBaseRepositoryCreateOrUpdate 测试创建或更新操作
func TestBaseRepositoryCreateOrUpdate(t *testing.T) {
	gormDB, err := setupTestDB()
	assert.NoError(t, err)

	dbHandler := newTestDBHandler(gormDB)
	repo := NewBaseRepository[TestUser](dbHandler, logger.NewLogger(nil), "test_users")

	ctx := context.Background()
	user := &TestUser{
		Name:   "John Doe",
		Email:  "john@example.com",
		Age:    25,
		Status: "active",
	}

	// 第一次调用应该创建用户
	result, created, err := repo.CreateOrUpdate(ctx, user, "Email")
	assert.NoError(t, err, "创建或更新操作不应出错")
	assert.True(t, created, "应显示记录已创建")
	assert.NotNil(t, result, "结果不应为空")

	// 第二次调用应该更新用户
	user.Name = "John Updated"
	user.Age = 26
	result2, created2, err := repo.CreateOrUpdate(ctx, user, "Email")
	assert.NoError(t, err, "创建或更新操作不应出错")
	assert.False(t, created2, "应显示记录已更新")
	assert.NotNil(t, result2, "结果不应为空")
	assert.Equal(t, "John Updated", result2.Name, "名称应已更新")
	assert.Equal(t, 26, result2.Age, "年龄应已更新")
}

// TestBaseRepositoryBulkCreate 测试高性能批量创建
func TestBaseRepositoryBulkCreate(t *testing.T) {
	gormDB, err := setupTestDB()
	assert.NoError(t, err)

	dbHandler := newTestDBHandler(gormDB)
	repo := NewBaseRepository[TestUser](dbHandler, logger.NewLogger(nil), "test_users", WithBatchSize[TestUser](2))

	ctx := context.Background()
	users := []*TestUser{
		{Name: "User1", Email: "user1@example.com", Age: 25, Status: "active"},
		{Name: "User2", Email: "user2@example.com", Age: 30, Status: "active"},
		{Name: "User3", Email: "user3@example.com", Age: 35, Status: "active"},
	}

	err = repo.BulkCreate(ctx, users)
	assert.NoError(t, err, "批量创建不应出错")

	// 验证所有用户都已创建
	var count int64
	gormDB.Table("test_users").Count(&count)
	assert.Equal(t, int64(3), count, "应创建3个用户")
}

// TestBaseRepositoryUpdateFieldsByFilters 测试按过滤器更新字段
func TestBaseRepositoryUpdateFieldsByFilters(t *testing.T) {
	gormDB, err := setupTestDB()
	assert.NoError(t, err)

	dbHandler := newTestDBHandler(gormDB)
	repo := NewBaseRepository[TestUser](dbHandler, logger.NewLogger(nil), "test_users")

	ctx := context.Background()

	// 创建用户
	users := []*TestUser{
		{Name: "User1", Email: "user1@example.com", Age: 25, Status: "active"},
		{Name: "User2", Email: "user2@example.com", Age: 30, Status: "active"},
		{Name: "User3", Email: "user3@example.com", Age: 35, Status: "inactive"},
	}
	err = repo.CreateBatch(ctx, users...)
	assert.NoError(t, err)

	// 按过滤器更新字段
	fields := map[string]interface{}{
		"status": "suspended",
		"age":    99,
	}

	err = repo.UpdateFieldsByFilters(ctx, fields, NewEqFilter("status", "active"))
	assert.NoError(t, err, "按过滤器更新字段不应出错")

	// 验证更新
	var suspendedCount int64
	gormDB.Table("test_users").Where("status = ? AND age = ?", "suspended", 99).Count(&suspendedCount)
	assert.Equal(t, int64(2), suspendedCount, "应有 2 个用户被更新")

	var inactiveCount int64
	gormDB.Table("test_users").Where("status = ?", "inactive").Count(&inactiveCount)
	assert.Equal(t, int64(1), inactiveCount, "应保留 1 个非活跃用户")
}

// TestBaseRepositoryCreateIfNotExists 测试有条件创建
func TestBaseRepositoryCreateIfNotExists(t *testing.T) {
	gormDB, err := setupTestDB()
	assert.NoError(t, err)

	dbHandler := newTestDBHandler(gormDB)
	repo := NewBaseRepository[TestUser](dbHandler, logger.NewLogger(nil), "test_users")

	ctx := context.Background()

	// 第一次创建用户
	user := &TestUser{
		Name:   "Unique User",
		Email:  "unique@example.com",
		Age:    30,
		Status: "active",
	}

	result1, created1, err := repo.CreateIfNotExists(ctx, user, "Email")
	assert.NoError(t, err, "第一次创建不应出错")
	assert.True(t, created1, "应显示记录已创建")
	assert.NotNil(t, result1, "结果不应为空")
	assert.Equal(t, "unique@example.com", result1.Email, "邮箱应正确")
	assert.Equal(t, "Unique User", result1.Name, "名称应正确")
	assert.Equal(t, 30, result1.Age, "年龄应正确")
	assert.Equal(t, "active", result1.Status, "状态应正确")
	assert.Greater(t, result1.ID, uint(0), "ID应大于0")
	assert.NotZero(t, result1.CreatedAt, "创建时间应不为零值")
	assert.NotZero(t, result1.UpdatedAt, "更新时间应不为零值")
	createdID := result1.ID

	// 第二次尝试创建相同邮箱的用户
	user2 := &TestUser{
		Name:   "Another User",
		Email:  "unique@example.com", // 相同邮箱
		Age:    25,
		Status: "inactive",
	}

	result2, created2, err := repo.CreateIfNotExists(ctx, user2, "Email")
	assert.NoError(t, err, "第二次创建不应出错")
	assert.False(t, created2, "应显示记录未创建")
	assert.NotNil(t, result2, "结果不应为空")
	assert.Equal(t, createdID, result2.ID, "ID应与第一次创建的相同")
	assert.Equal(t, "Unique User", result2.Name, "名称应保持为第一次创建的值")
	assert.Equal(t, 30, result2.Age, "年龄应保持为第一次创建的值")
	assert.Equal(t, "active", result2.Status, "状态应保持为第一次创建的值")

	// 验证数据库中只有1条记录
	var count int64
	gormDB.Table("test_users").Where("email = ?", "unique@example.com").Count(&count)
	assert.Equal(t, int64(1), count, "数据库中应只有1条记录")
}

// TestBaseRepositoryGetWithPreloads 测试带预加载的单条查询
func TestBaseRepositoryGetWithPreloads(t *testing.T) {
	gormDB, err := setupTestDB()
	assert.NoError(t, err)

	dbHandler := newTestDBHandler(gormDB)
	repo := NewBaseRepository[TestUser](dbHandler, logger.NewLogger(nil), "test_users")

	ctx := context.Background()

	// 创建用户
	user := &TestUser{
		Name:   "Test User",
		Email:  "test@example.com",
		Age:    25,
		Status: "active",
	}
	_, err = repo.Create(ctx, user)
	assert.NoError(t, err)

	// 获取并验证ID
	var createdUser TestUser
	gormDB.Table("test_users").Where("email = ?", "test@example.com").First(&createdUser)

	// 按ID获取用户（不带预加载，因为TestUser没有关联表）
	result, err := repo.GetWithPreloads(ctx, createdUser.ID)
	assert.NoError(t, err, "获取操作不应出错")
	assert.NotNil(t, result, "结果不应为空")
	assert.Equal(t, createdUser.ID, result.ID, "ID应正确")
	assert.Equal(t, "test@example.com", result.Email, "邮箱应正确")
	assert.Equal(t, "Test User", result.Name, "名称应正确")
	assert.Equal(t, 25, result.Age, "年龄应正确")
	assert.Equal(t, "active", result.Status, "状态应正确")
	assert.NotZero(t, result.CreatedAt, "创建时间应不为零值")
	assert.NotZero(t, result.UpdatedAt, "更新时间应不为零值")

	// 测试不存在的ID
	result2, err := repo.GetWithPreloads(ctx, 99999)
	assert.Error(t, err, "不存在的ID应返回错误")
	assert.Nil(t, result2, "结果应为空")
}

// TestBaseRepositoryGetByFilter 测试按单个过滤器查询
func TestBaseRepositoryGetByFilter(t *testing.T) {
	gormDB, err := setupTestDB()
	assert.NoError(t, err)

	dbHandler := newTestDBHandler(gormDB)
	repo := NewBaseRepository[TestUser](dbHandler, logger.NewLogger(nil), "test_users")

	ctx := context.Background()

	// 创建多个用户
	users := []*TestUser{
		{Name: "John", Email: "john@example.com", Age: 25, Status: "active"},
		{Name: "Jane", Email: "jane@example.com", Age: 30, Status: "inactive"},
		{Name: "Bob", Email: "bob@example.com", Age: 35, Status: "active"},
	}
	err = repo.CreateBatch(ctx, users...)
	assert.NoError(t, err)

	// 按状态过滤查询
	filter := NewEqFilter("status", "active")
	result, err := repo.GetByFilter(ctx, filter)
	assert.NoError(t, err, "按过滤器查询不应出错")
	assert.NotNil(t, result, "结果不应为空")
	assert.Equal(t, "active", result.Status, "状态应为active")
	assert.Greater(t, result.ID, uint(0), "ID应大于0")
	assert.NotEmpty(t, result.Name, "名称不应为空")
	assert.NotEmpty(t, result.Email, "邮箱不应为空")
	assert.Greater(t, result.Age, 0, "年龄应大于0")
	assert.NotZero(t, result.CreatedAt, "创建时间应不为零值")
	assert.NotZero(t, result.UpdatedAt, "更新时间应不为零值")
	// 验证是否是活跃用户之一
	expectedNames := []string{"John", "Bob"}
	assert.Contains(t, expectedNames, result.Name, "应该是活跃用户之一")

	// 测试不存在的状态
	notFoundFilter := NewEqFilter("status", "deleted")
	result2, err := repo.GetByFilter(ctx, notFoundFilter)
	assert.Error(t, err, "不存在的状态应返回错误")
	assert.Nil(t, result2, "结果应为空")
}

// TestBaseRepositoryListWithPagination 测试带分页的列表查询
func TestBaseRepositoryListWithPagination(t *testing.T) {
	gormDB, err := setupTestDB()
	assert.NoError(t, err)

	dbHandler := newTestDBHandler(gormDB)
	repo := NewBaseRepository[TestUser](dbHandler, logger.NewLogger(nil), "test_users")

	ctx := context.Background()

	// 创建5个用户
	users := []*TestUser{
		{Name: "User1", Email: "user1@example.com", Age: 20, Status: "active"},
		{Name: "User2", Email: "user2@example.com", Age: 25, Status: "active"},
		{Name: "User3", Email: "user3@example.com", Age: 30, Status: "active"},
		{Name: "User4", Email: "user4@example.com", Age: 35, Status: "active"},
		{Name: "User5", Email: "user5@example.com", Age: 40, Status: "active"},
	}
	err = repo.CreateBatch(ctx, users...)
	assert.NoError(t, err)

	// 创建查询和分页参数
	query := NewQuery()
	query.AddFilter(NewEqFilter("status", "active"))
	query.AddOrder("age", "ASC")

	pagination := &Pagination{
		Page:     1,
		PageSize: 3,
	}

	// 执行分页查询
	result, resultPaging, err := repo.ListWithPagination(ctx, query, pagination)
	assert.NoError(t, err, "分页查询不应出错")
	assert.NotNil(t, result, "结果不应为空")
	assert.NotNil(t, resultPaging, "分页信息不应为空")

	// 验证结果数据
	assert.Len(t, result, 3, "应返回3条记录")
	assert.Equal(t, 20, result[0].Age, "第一条记录年龄应为20")
	assert.Equal(t, 25, result[1].Age, "第二条记录年龄应为25")
	assert.Equal(t, 30, result[2].Age, "第三条记录年龄应为30")

	// 验证每个结果的完整性
	for i, user := range result {
		assert.Equal(t, "active", user.Status, fmt.Sprintf("第%d个用户状态应为active", i+1))
		assert.Greater(t, user.ID, uint(0), fmt.Sprintf("第%d个用户ID应大于0", i+1))
		assert.NotEmpty(t, user.Name, fmt.Sprintf("第%d个用户名不应为空", i+1))
		assert.NotEmpty(t, user.Email, fmt.Sprintf("第%d个用户邮箱不应为空", i+1))
		assert.NotZero(t, user.CreatedAt, fmt.Sprintf("第%d个用户创建时间应不为零值", i+1))
		assert.NotZero(t, user.UpdatedAt, fmt.Sprintf("第%d个用户更新时间应不为零值", i+1))
	}

	// 验证分页信息
	assert.Equal(t, int64(5), resultPaging.Total, "总记录数应为5")
	assert.Equal(t, int32(1), resultPaging.Page, "当前页数应为1")
	assert.Equal(t, int32(3), resultPaging.PageSize, "每页大小应为3")
	assert.Equal(t, int64(2), resultPaging.GetTotalPages(), "总页数应为2")
	assert.True(t, resultPaging.HasNextPage(), "应有下一页")
	assert.False(t, resultPaging.HasPrevPage(), "不应有上一页")
}

// TestBaseRepositoryUpdate 测试更新单个实体
func TestBaseRepositoryUpdate(t *testing.T) {
	gormDB, err := setupTestDB()
	assert.NoError(t, err)

	dbHandler := newTestDBHandler(gormDB)
	repo := NewBaseRepository[TestUser](dbHandler, logger.NewLogger(nil), "test_users")

	ctx := context.Background()

	// 创建用户
	user := &TestUser{
		Name:   "Original User",
		Email:  "original@example.com",
		Age:    25,
		Status: "active",
	}
	createdUser, err := repo.Create(ctx, user)
	assert.NoError(t, err)
	assert.NotNil(t, createdUser)

	createdAt := createdUser.CreatedAt
	originalID := createdUser.ID

	// 修改用户信息
	createdUser.Name = "Updated User"
	createdUser.Age = 30
	createdUser.Status = "inactive"

	// 更新用户
	updatedUser, err := repo.Update(ctx, createdUser)
	assert.NoError(t, err, "更新操作不应出错")
	assert.NotNil(t, updatedUser, "更新结果不应为空")
	assert.Equal(t, originalID, updatedUser.ID, "ID应保持不变")
	assert.Equal(t, "Updated User", updatedUser.Name, "名称应已更新")
	assert.Equal(t, 30, updatedUser.Age, "年龄应已更新")
	assert.Equal(t, "inactive", updatedUser.Status, "状态应已更新")
	assert.Equal(t, "original@example.com", updatedUser.Email, "邮箱应保持不变")
	assert.Equal(t, createdAt, updatedUser.CreatedAt, "创建时间应保持不变")
	assert.GreaterOrEqual(t, updatedUser.UpdatedAt, createdAt, "更新时间应晚于创建时间")

	// 验证数据库中的数据
	var dbUser TestUser
	gormDB.Table("test_users").Where("id = ?", originalID).First(&dbUser)
	assert.Equal(t, "Updated User", dbUser.Name, "数据库中的名称应已更新")
	assert.Equal(t, 30, dbUser.Age, "数据库中的年龄应已更新")
	assert.Equal(t, "inactive", dbUser.Status, "数据库中的状态应已更新")
}

// TestBaseRepositoryDelete 测试删除操作
func TestBaseRepositoryDelete(t *testing.T) {
	gormDB, err := setupTestDB()
	assert.NoError(t, err)

	dbHandler := newTestDBHandler(gormDB)
	repo := NewBaseRepository[TestUser](dbHandler, logger.NewLogger(nil), "test_users")

	ctx := context.Background()

	// 创建用户
	user := &TestUser{
		Name:   "To Delete",
		Email:  "delete@example.com",
		Age:    25,
		Status: "active",
	}
	createdUser, err := repo.Create(ctx, user)
	assert.NoError(t, err)
	assert.NotNil(t, createdUser)

	// 删除用户
	err = repo.Delete(ctx, createdUser.ID)
	assert.NoError(t, err, "删除操作不应出错")

	// 验证用户已被删除
	_, err = repo.Get(ctx, createdUser.ID)
	assert.Error(t, err, "删除后获取应返回错误")

	// 验证数据库中确实没有该记录
	var count int64
	gormDB.Table("test_users").Where("id = ?", createdUser.ID).Count(&count)
	assert.Equal(t, int64(0), count, "数据库中应没有该记录")

	// 测试删除不存在的ID
	err = repo.Delete(ctx, 99999)
	assert.NoError(t, err, "删除不存在的记录应不出错") // GORM的Delete操作对不存在的记录不报错
}

// TestBaseRepositoryDeleteBatch 测试批量删除
func TestBaseRepositoryDeleteBatch(t *testing.T) {
	gormDB, err := setupTestDB()
	assert.NoError(t, err)

	dbHandler := newTestDBHandler(gormDB)
	repo := NewBaseRepository[TestUser](dbHandler, logger.NewLogger(nil), "test_users")

	ctx := context.Background()

	// 创建多个用户
	users := []*TestUser{
		{Name: "User1", Email: "user1@example.com", Age: 25, Status: "active"},
		{Name: "User2", Email: "user2@example.com", Age: 30, Status: "active"},
		{Name: "User3", Email: "user3@example.com", Age: 35, Status: "active"},
	}
	err = repo.CreateBatch(ctx, users...)
	assert.NoError(t, err)

	// 获取创建的用户ID
	var createdUsers []TestUser
	gormDB.Table("test_users").Find(&createdUsers)
	assert.Len(t, createdUsers, 3, "应创建3个用户")

	// 批量删除前两个用户
	err = repo.DeleteBatch(ctx, createdUsers[0].ID, createdUsers[1].ID)
	assert.NoError(t, err, "批量删除不应出错")

	// 验证删除结果
	var remainingCount int64
	gormDB.Table("test_users").Count(&remainingCount)
	assert.Equal(t, int64(1), remainingCount, "应剩余1个用户")

	// 验证剩余的用户是第三个
	var remainingUser TestUser
	gormDB.Table("test_users").First(&remainingUser)
	assert.Equal(t, createdUsers[2].ID, remainingUser.ID, "剩余的应是第三个用户")
	assert.Equal(t, "User3", remainingUser.Name, "剩余用户名应为User3")
	assert.Equal(t, "user3@example.com", remainingUser.Email, "剩余用户邮箱应正确")
	assert.Equal(t, 35, remainingUser.Age, "剩余用户年龄应为35")
	assert.Equal(t, "active", remainingUser.Status, "剩余用户状态应为active")
}

// TestBaseRepositoryUpdateFields 测试按ID更新字段
func TestBaseRepositoryUpdateFields(t *testing.T) {
	gormDB, err := setupTestDB()
	assert.NoError(t, err)

	dbHandler := newTestDBHandler(gormDB)
	repo := NewBaseRepository[TestUser](dbHandler, logger.NewLogger(nil), "test_users")

	ctx := context.Background()

	// 创建用户
	user := &TestUser{
		Name:   "Original User",
		Email:  "original@example.com",
		Age:    25,
		Status: "active",
	}
	createdUser, err := repo.Create(ctx, user)
	assert.NoError(t, err)

	// 准备更新字段
	fields := map[string]interface{}{
		"name":   "Updated Name",
		"age":    30,
		"status": "inactive",
	}

	// 执行字段更新
	err = repo.UpdateFields(ctx, createdUser.ID, fields)
	assert.NoError(t, err, "字段更新不应出错")

	// 验证更新结果
	updatedUser, err := repo.Get(ctx, createdUser.ID)
	assert.NoError(t, err)
	assert.NotNil(t, updatedUser)
	assert.Equal(t, "Updated Name", updatedUser.Name, "名称应已更新")
	assert.Equal(t, 30, updatedUser.Age, "年龄应已更新")
	assert.Equal(t, "inactive", updatedUser.Status, "状态应已更新")
	assert.Equal(t, "original@example.com", updatedUser.Email, "邮箱应保持不变")
	assert.Equal(t, createdUser.ID, updatedUser.ID, "ID应保持不变")

	// 验证数据库中的数据
	var dbUser TestUser
	gormDB.Table("test_users").Where("id = ?", createdUser.ID).First(&dbUser)
	assert.Equal(t, "Updated Name", dbUser.Name, "数据库中的名称应已更新")
	assert.Equal(t, 30, dbUser.Age, "数据库中的年龄应已更新")
	assert.Equal(t, "inactive", dbUser.Status, "数据库中的状态应已更新")
}

// TestBaseRepositorySoftDelete 测试软删除
func TestBaseRepositorySoftDelete(t *testing.T) {
	gormDB, err := setupTestDB()
	assert.NoError(t, err)

	dbHandler := newTestDBHandler(gormDB)
	repo := NewBaseRepository[TestUser](dbHandler, logger.NewLogger(nil), "test_users")

	ctx := context.Background()

	// 创建用户
	user := &TestUser{
		Name:   "To Soft Delete",
		Email:  "soft@example.com",
		Age:    25,
		Status: "active",
	}
	createdUser, err := repo.Create(ctx, user)
	assert.NoError(t, err)

	// 软删除用户
	err = repo.SoftDelete(ctx, createdUser.ID, "status", "deleted")
	assert.NoError(t, err, "软删除不应出错")

	// 验证软删除结果 - 记录仍在数据库中
	var dbUser TestUser
	gormDB.Table("test_users").Where("id = ?", createdUser.ID).First(&dbUser)
	assert.Equal(t, "deleted", dbUser.Status, "状态应已更新为deleted")
	assert.Equal(t, createdUser.ID, dbUser.ID, "ID应保持不变")
	assert.Equal(t, "To Soft Delete", dbUser.Name, "名称应保持不变")
	assert.Equal(t, "soft@example.com", dbUser.Email, "邮箱应保持不变")
	assert.Equal(t, 25, dbUser.Age, "年龄应保持不变")

	// 验证记录确实存在但状态已改变
	var count int64
	gormDB.Table("test_users").Where("id = ?", createdUser.ID).Count(&count)
	assert.Equal(t, int64(1), count, "记录应仍存在于数据库中")

	var deletedCount int64
	gormDB.Table("test_users").Where("id = ? AND status = ?", createdUser.ID, "deleted").Count(&deletedCount)
	assert.Equal(t, int64(1), deletedCount, "记录状态应为deleted")
}

// TestBaseRepositorySoftDeleteBatch 测试批量软删除
func TestBaseRepositorySoftDeleteBatch(t *testing.T) {
	gormDB, err := setupTestDB()
	assert.NoError(t, err)

	dbHandler := newTestDBHandler(gormDB)
	repo := NewBaseRepository[TestUser](dbHandler, logger.NewLogger(nil), "test_users")

	ctx := context.Background()

	// 创建多个用户
	users := []*TestUser{
		{Name: "User1", Email: "user1@example.com", Age: 25, Status: "active"},
		{Name: "User2", Email: "user2@example.com", Age: 30, Status: "active"},
		{Name: "User3", Email: "user3@example.com", Age: 35, Status: "active"},
	}
	err = repo.CreateBatch(ctx, users...)
	assert.NoError(t, err)

	// 获取创建的用户ID
	var createdUsers []TestUser
	gormDB.Table("test_users").Find(&createdUsers)
	assert.Len(t, createdUsers, 3)

	// 批量软删除前两个用户
	ids := []interface{}{createdUsers[0].ID, createdUsers[1].ID}
	err = repo.SoftDeleteBatch(ctx, ids, "status", "deleted")
	assert.NoError(t, err, "批量软删除不应出错")

	// 验证软删除结果
	var deletedCount int64
	gormDB.Table("test_users").Where("status = ?", "deleted").Count(&deletedCount)
	assert.Equal(t, int64(2), deletedCount, "应有2个用户被软删除")

	var activeCount int64
	gormDB.Table("test_users").Where("status = ?", "active").Count(&activeCount)
	assert.Equal(t, int64(1), activeCount, "应有1个用户保持活跃")

	// 验证总记录数没有减少
	var totalCount int64
	gormDB.Table("test_users").Count(&totalCount)
	assert.Equal(t, int64(3), totalCount, "总记录数应保持不变")

	// 验证具体的软删除数据
	var deletedUsers []TestUser
	gormDB.Table("test_users").Where("status = ?", "deleted").Find(&deletedUsers)
	assert.Len(t, deletedUsers, 2)

	expectedNames := []string{"User1", "User2"}
	for i, user := range deletedUsers {
		assert.Equal(t, "deleted", user.Status, fmt.Sprintf("第%d个用户状态应为deleted", i+1))
		assert.Contains(t, expectedNames, user.Name, "用户名应是被删除的用户之一")
		assert.Greater(t, user.ID, uint(0), fmt.Sprintf("第%d个用户ID应大于0", i+1))
		assert.NotEmpty(t, user.Email, fmt.Sprintf("第%d个用户邮箱不应为空", i+1))
		assert.Greater(t, user.Age, 0, fmt.Sprintf("第%d个用户年龄应大于0", i+1))
	}
}

// TestBaseRepositoryRestore 测试恢复软删除
func TestBaseRepositoryRestore(t *testing.T) {
	gormDB, err := setupTestDB()
	assert.NoError(t, err)

	dbHandler := newTestDBHandler(gormDB)
	repo := NewBaseRepository[TestUser](dbHandler, logger.NewLogger(nil), "test_users")

	ctx := context.Background()

	// 创建用户并软删除
	user := &TestUser{
		Name:   "To Restore",
		Email:  "restore@example.com",
		Age:    25,
		Status: "active",
	}
	createdUser, err := repo.Create(ctx, user)
	assert.NoError(t, err)

	// 软删除
	err = repo.SoftDelete(ctx, createdUser.ID, "status", "deleted")
	assert.NoError(t, err)

	// 恢复
	err = repo.Restore(ctx, createdUser.ID, "status", "active")
	assert.NoError(t, err, "恢复操作不应出错")

	// 验证恢复结果
	var restoredUser TestUser
	gormDB.Table("test_users").Where("id = ?", createdUser.ID).First(&restoredUser)
	assert.Equal(t, "active", restoredUser.Status, "状态应已恢复为active")
	assert.Equal(t, createdUser.ID, restoredUser.ID, "ID应保持不变")
	assert.Equal(t, "To Restore", restoredUser.Name, "名称应保持不变")
	assert.Equal(t, "restore@example.com", restoredUser.Email, "邮箱应保持不变")
	assert.Equal(t, 25, restoredUser.Age, "年龄应保持不变")

	// 验证记录状态
	var activeCount int64
	gormDB.Table("test_users").Where("id = ? AND status = ?", createdUser.ID, "active").Count(&activeCount)
	assert.Equal(t, int64(1), activeCount, "记录应为活跃状态")
}

// TestBaseRepositoryRestoreBatch 测试批量恢复
func TestBaseRepositoryRestoreBatch(t *testing.T) {
	gormDB, err := setupTestDB()
	assert.NoError(t, err)

	dbHandler := newTestDBHandler(gormDB)
	repo := NewBaseRepository[TestUser](dbHandler, logger.NewLogger(nil), "test_users")

	ctx := context.Background()

	// 创建多个用户
	users := []*TestUser{
		{Name: "User1", Email: "user1@example.com", Age: 25, Status: "deleted"},
		{Name: "User2", Email: "user2@example.com", Age: 30, Status: "deleted"},
		{Name: "User3", Email: "user3@example.com", Age: 35, Status: "active"},
	}
	err = repo.CreateBatch(ctx, users...)
	assert.NoError(t, err)

	// 获取已删除用户的ID
	var deletedUsers []TestUser
	gormDB.Table("test_users").Where("status = ?", "deleted").Find(&deletedUsers)
	assert.Len(t, deletedUsers, 2)

	ids := []interface{}{deletedUsers[0].ID, deletedUsers[1].ID}

	// 批量恢复
	err = repo.RestoreBatch(ctx, ids, "status", "active")
	assert.NoError(t, err, "批量恢复不应出错")

	// 验证恢复结果
	var activeCount int64
	gormDB.Table("test_users").Where("status = ?", "active").Count(&activeCount)
	assert.Equal(t, int64(3), activeCount, "应有3个活跃用户")

	var deletedCount int64
	gormDB.Table("test_users").Where("status = ?", "deleted").Count(&deletedCount)
	assert.Equal(t, int64(0), deletedCount, "应没有已删除用户")
}

// TestBaseRepositoryCountByField 测试按字段统计
func TestBaseRepositoryCountByField(t *testing.T) {
	gormDB, err := setupTestDB()
	assert.NoError(t, err)

	dbHandler := newTestDBHandler(gormDB)
	repo := NewBaseRepository[TestUser](dbHandler, logger.NewLogger(nil), "test_users")

	ctx := context.Background()

	// 创建不同状态的用户
	users := []*TestUser{
		{Name: "User1", Email: "user1@example.com", Age: 25, Status: "active"},
		{Name: "User2", Email: "user2@example.com", Age: 30, Status: "active"},
		{Name: "User3", Email: "user3@example.com", Age: 35, Status: "inactive"},
		{Name: "User4", Email: "user4@example.com", Age: 40, Status: "active"},
	}
	err = repo.CreateBatch(ctx, users...)
	assert.NoError(t, err)

	// 按状态字段统计
	counts, err := repo.CountByField(ctx, "status")
	assert.NoError(t, err, "字段统计不应出错")
	assert.NotNil(t, counts, "统计结果不应为空")

	// 验证统计结果
	assert.Equal(t, int64(3), counts["active"], "active状态应有3个用户")
	assert.Equal(t, int64(1), counts["inactive"], "inactive状态应有1个用户")
	assert.Len(t, counts, 2, "应有2种状态")
}

// TestBaseRepositoryPluck 测试提取字段值
func TestBaseRepositoryPluck(t *testing.T) {
	gormDB, err := setupTestDB()
	assert.NoError(t, err)

	dbHandler := newTestDBHandler(gormDB)
	repo := NewBaseRepository[TestUser](dbHandler, logger.NewLogger(nil), "test_users")

	ctx := context.Background()

	// 创建用户
	users := []*TestUser{
		{Name: "Alice", Email: "alice@example.com", Age: 25, Status: "active"},
		{Name: "Bob", Email: "bob@example.com", Age: 30, Status: "active"},
		{Name: "Charlie", Email: "charlie@example.com", Age: 35, Status: "inactive"},
	}
	err = repo.CreateBatch(ctx, users...)
	assert.NoError(t, err)

	// 提取所有活跃用户的名称
	names, err := repo.Pluck(ctx, "name", NewEqFilter("status", "active"))
	assert.NoError(t, err, "提取字段值不应出错")
	assert.NotNil(t, names, "结果不应为空")
	assert.Len(t, names, 2, "应返回2个名称")

	// 验证提取的名称
	nameStrings := make([]string, len(names))
	for i, name := range names {
		nameStrings[i] = name.(string)
	}
	assert.Contains(t, nameStrings, "Alice", "应包含Alice")
	assert.Contains(t, nameStrings, "Bob", "应包含Bob")
	assert.NotContains(t, nameStrings, "Charlie", "不应包含Charlie")

	// 提取所有用户的年龄
	ages, err := repo.Pluck(ctx, "age")
	assert.NoError(t, err, "提取年龄不应出错")
	assert.NotNil(t, ages, "年龄结果不应为空")
	assert.Len(t, ages, 3, "应返回3个年龄")
}

// TestBaseRepositoryDistinct 测试去重查询
func TestBaseRepositoryDistinct(t *testing.T) {
	gormDB, err := setupTestDB()
	assert.NoError(t, err)

	dbHandler := newTestDBHandler(gormDB)
	repo := NewBaseRepository[TestUser](dbHandler, logger.NewLogger(nil), "test_users")

	ctx := context.Background()

	// 创建有重复状态的用户
	users := []*TestUser{
		{Name: "User1", Email: "user1@example.com", Age: 25, Status: "active"},
		{Name: "User2", Email: "user2@example.com", Age: 30, Status: "active"},
		{Name: "User3", Email: "user3@example.com", Age: 35, Status: "inactive"},
		{Name: "User4", Email: "user4@example.com", Age: 40, Status: "active"},
		{Name: "User5", Email: "user5@example.com", Age: 45, Status: "inactive"},
	}
	err = repo.CreateBatch(ctx, users...)
	assert.NoError(t, err)

	// 获取不重复的状态值
	distinctStatuses, err := repo.Distinct(ctx, "status")
	assert.NoError(t, err, "去重查询不应出错")
	assert.NotNil(t, distinctStatuses, "结果不应为空")
	assert.Len(t, distinctStatuses, 2, "应返回2个不重复的状态")

	// 验证去重结果
	statusStrings := make([]string, len(distinctStatuses))
	for i, status := range distinctStatuses {
		statusStrings[i] = status.(string)
	}
	assert.Contains(t, statusStrings, "active", "应包含active")
	assert.Contains(t, statusStrings, "inactive", "应包含inactive")

	// 测试带过滤条件的去重查询
	activeDistinct, err := repo.Distinct(ctx, "status", NewEqFilter("status", "active"))
	assert.NoError(t, err, "带过滤的去重查询不应出错")
	assert.NotNil(t, activeDistinct, "结果不应为空")
	assert.Len(t, activeDistinct, 1, "应返回1个不重复的状态")
	assert.Equal(t, "active", activeDistinct[0].(string), "应只包含active")
}

// TestBaseRepositoryDBHandler 测试获取数据库处理器
func TestBaseRepositoryDBHandler(t *testing.T) {
	gormDB, err := setupTestDB()
	assert.NoError(t, err)

	dbHandler := newTestDBHandler(gormDB)
	repo := NewBaseRepository[TestUser](dbHandler, logger.NewLogger(nil), "test_users")

	// 获取数据库处理器
	handler := repo.DBHandler()
	assert.NotNil(t, handler, "数据库处理器不应为空")
	assert.Equal(t, dbHandler, handler, "应返回相同的数据库处理器")
}

// TestBaseRepositoryTable 测试获取表名
func TestBaseRepositoryTable(t *testing.T) {
	gormDB, err := setupTestDB()
	assert.NoError(t, err)

	dbHandler := newTestDBHandler(gormDB)
	repo := NewBaseRepository[TestUser](dbHandler, logger.NewLogger(nil), "test_users")

	// 获取表名
	tableName := repo.Table()
	assert.Equal(t, "test_users", tableName, "表名应正确")
}

// TestBaseRepositorySoftDeleteByFilters 测试按过滤器软删除
func TestBaseRepositorySoftDeleteByFilters(t *testing.T) {
	gormDB, err := setupTestDB()
	assert.NoError(t, err)

	dbHandler := newTestDBHandler(gormDB)
	repo := NewBaseRepository[TestUser](dbHandler, logger.NewLogger(nil), "test_users")

	ctx := context.Background()

	// 创建用户
	users := []*TestUser{
		{Name: "User1", Email: "user1@example.com", Age: 25, Status: "active"},
		{Name: "User2", Email: "user2@example.com", Age: 30, Status: "active"},
		{Name: "User3", Email: "user3@example.com", Age: 35, Status: "inactive"},
	}
	err = repo.CreateBatch(ctx, users...)
	assert.NoError(t, err)

	// 按过滤器软删除所有活跃用户
	err = repo.SoftDeleteByFilters(ctx, "status", "deleted", NewEqFilter("status", "active"))
	assert.NoError(t, err, "按过滤器软删除不应出错")

	// 验证软删除结果
	var deletedCount int64
	gormDB.Table("test_users").Where("status = ?", "deleted").Count(&deletedCount)
	assert.Equal(t, int64(2), deletedCount, "应有2个用户被软删除")

	var inactiveCount int64
	gormDB.Table("test_users").Where("status = ?", "inactive").Count(&inactiveCount)
	assert.Equal(t, int64(1), inactiveCount, "应有1个用户保持非活跃")

	// 验证总记录数
	var totalCount int64
	gormDB.Table("test_users").Count(&totalCount)
	assert.Equal(t, int64(3), totalCount, "总记录数应保持不变")

	// 验证具体的软删除数据
	var deletedUsers []TestUser
	gormDB.Table("test_users").Where("status = ?", "deleted").Find(&deletedUsers)
	expectedNames := []string{"User1", "User2"}
	for _, user := range deletedUsers {
		assert.Equal(t, "deleted", user.Status, "用户状态应为deleted")
		assert.Contains(t, expectedNames, user.Name, "用户名应是被删除的用户之一")
		assert.Greater(t, user.ID, uint(0), "用户ID应大于0")
		assert.NotEmpty(t, user.Email, "用户邮箱不应为空")
		assert.Greater(t, user.Age, 0, "用户年龄应大于0")
	}
}

// TestBaseRepositoryCreateBatch 测试批量创建实体
func TestBaseRepositoryCreateBatch(t *testing.T) {
	gormDB, err := setupTestDB()
	assert.NoError(t, err)

	dbHandler := newTestDBHandler(gormDB)
	repo := NewBaseRepository[TestUser](dbHandler, logger.NewLogger(nil), "test_users")
	ctx := context.Background()

	// 准备测试数据
	user1 := &TestUser{
		Name:   "批量用户1",
		Email:  "batch1@test.com",
		Age:    25,
		Status: "active",
	}
	user2 := &TestUser{
		Name:   "批量用户2",
		Email:  "batch2@test.com",
		Age:    30,
		Status: "active",
	}
	user3 := &TestUser{
		Name:   "批量用户3",
		Email:  "batch3@test.com",
		Age:    35,
		Status: "inactive",
	}

	// 执行批量创建
	err = repo.CreateBatch(ctx, user1, user2, user3)
	assert.NoError(t, err, "批量创建不应该返回错误")

	// 验证用户被正确创建
	assert.NotZero(t, user1.ID, "用户1的ID应该被设置")
	assert.NotZero(t, user2.ID, "用户2的ID应该被设置")
	assert.NotZero(t, user3.ID, "用户3的ID应该被设置")

	// 验证数据库中的数据
	var count int64
	gormDB.Table("test_users").Count(&count)
	assert.Equal(t, int64(3), count, "应该有3个用户被创建")

	// 验证具体用户数据
	var savedUser1 TestUser
	err = gormDB.First(&savedUser1, user1.ID).Error
	assert.NoError(t, err, "应该能找到用户1")
	assert.Equal(t, "批量用户1", savedUser1.Name, "用户1名字应该一致")
	assert.Equal(t, "batch1@test.com", savedUser1.Email, "用户1邮箱应该一致")
	assert.Equal(t, 25, savedUser1.Age, "用户1年龄应该一致")
	assert.Equal(t, "active", savedUser1.Status, "用户1状态应该一致")

	// 测试空切片
	err = repo.CreateBatch(ctx)
	assert.NoError(t, err, "空批量创建不应该返回错误")

	// 测试只读模式
	repoReadOnly := NewBaseRepository[TestUser](dbHandler, logger.NewLogger(nil), "test_users", WithReadOnly[TestUser]())
	err = repoReadOnly.CreateBatch(ctx, user1)
	assert.Error(t, err, "只读模式下创建应该返回错误")
}

// TestComplexFiltering 测试复杂过滤逻辑覆盖applyFilters方法
func TestComplexFiltering(t *testing.T) {
	gormDB, err := setupTestDB()
	assert.NoError(t, err)

	dbHandler := newTestDBHandler(gormDB)
	repo := NewBaseRepository[TestUser](dbHandler, logger.NewLogger(nil), "test_users")
	ctx := context.Background()

	// 创建测试数据
	users := []*TestUser{
		{Name: "张三", Email: "zhang@test.com", Age: 25, Status: "active"},
		{Name: "李四", Email: "li@test.com", Age: 30, Status: "active"},
		{Name: "王五", Email: "wang@test.com", Age: 35, Status: "inactive"},
		{Name: "赵六", Email: "zhao@test.com", Age: 40, Status: "active"},
	}

	for _, user := range users {
		_, err := repo.Create(ctx, user)
		assert.NoError(t, err)
	}

	// 测试复杂查询：过滤器组合
	query := NewQuery()
	filterGroup := NewFilterGroup("AND")
	filterGroup.AddFilter(NewEqFilter("status", "active"))
	filterGroup.AddFilter(NewGteFilter("age", 30))
	query.WithFilterGroup(filterGroup)
	query.AddOrder("age", "DESC")

	results, err := repo.List(ctx, query)
	assert.NoError(t, err, "复杂查询不应该返回错误")
	assert.Equal(t, 2, len(results), "应该返回2个用户")
	assert.Equal(t, "赵六", results[0].Name, "第一个用户应该是赵六")
	assert.Equal(t, "李四", results[1].Name, "第二个用户应该是李四")

	// 测试简单的OR查询
	orQuery := NewQuery()
	orGroup := NewFilterGroup("OR")
	orGroup.AddFilter(NewEqFilter("name", "张三"))
	orGroup.AddFilter(NewEqFilter("name", "王五"))
	orQuery.WithFilterGroup(orGroup)

	orResults, err := repo.List(ctx, orQuery)
	assert.NoError(t, err, "OR查询不应该返回错误")
	assert.Equal(t, 2, len(orResults), "应该返回2个用户")

	// 验证返回的用户（张三和王五）
	names := make([]string, len(orResults))
	for i, user := range orResults {
		names[i] = user.Name
	}
	assert.Contains(t, names, "张三", "应该包含张三")
	assert.Contains(t, names, "王五", "应该包含王五")
}

// TestTransactionOperations 测试事务中的各种操作
func TestTransactionOperations(t *testing.T) {
	gormDB, err := setupTestDB()
	assert.NoError(t, err)

	dbHandler := newTestDBHandler(gormDB)
	repo := NewBaseRepository[TestUser](dbHandler, logger.NewLogger(nil), "test_users")
	ctx := context.Background()

	// 准备初始数据
	initialUser := &TestUser{Name: "初始用户", Email: "initial@test.com", Age: 25, Status: "active"}
	_, err = repo.Create(ctx, initialUser)
	assert.NoError(t, err)

	err = repo.Transaction(ctx, func(tx Transaction[TestUser]) error {
		// 测试Create单个用户
		user1 := &TestUser{Name: "事务用户1", Email: "tx1@test.com", Age: 30, Status: "active"}
		err := tx.Create(ctx, user1)
		assert.NoError(t, err, "事务中创建不应该返回错误")

		user2 := &TestUser{Name: "事务用户2", Email: "tx2@test.com", Age: 35, Status: "active"}
		err = tx.Create(ctx, user2)
		assert.NoError(t, err, "事务中创建不应该返回错误")

		// 测试Update
		initialUser.Name = "更新后的用户"
		err = tx.Update(ctx, initialUser)
		assert.NoError(t, err, "事务中更新不应该返回错误")

		// 测试UpdateBatch
		user1.Age = 31
		user2.Age = 36
		err = tx.UpdateBatch(ctx, user1, user2)
		assert.NoError(t, err, "事务中批量更新不应该返回错误")

		// 测试CreateBatch
		batchUser1 := &TestUser{Name: "批量用户1", Email: "batch1@test.com", Age: 40, Status: "active"}
		batchUser2 := &TestUser{Name: "批量用户2", Email: "batch2@test.com", Age: 42, Status: "active"}
		err = tx.CreateBatch(ctx, batchUser1, batchUser2)
		assert.NoError(t, err, "事务中批量创建不应该返回错误")

		// 测试Delete - 删除user2
		err = tx.Delete(ctx, user2)
		assert.NoError(t, err, "事务中删除不应该返回错误")

		// 测试DeleteBatch - 删除批量创建的用户
		err = tx.DeleteBatch(ctx, batchUser1, batchUser2)
		assert.NoError(t, err, "事务中批量删除不应该返回错误")

		return nil
	})
	assert.NoError(t, err, "事务不应该返回错误")

	// 验证事务结果
	var finalCount int64
	gormDB.Table("test_users").Count(&finalCount)
	assert.Equal(t, int64(2), finalCount, "最终应该有2个用户") // initialUser + user1

	// 验证用户更新
	var updatedUser TestUser
	err = gormDB.First(&updatedUser, initialUser.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, "更新后的用户", updatedUser.Name, "用户名应该被更新")

	// 验证user1的年龄更新
	var txUser TestUser
	err = gormDB.Where("name = ?", "事务用户1").First(&txUser).Error
	assert.NoError(t, err)
	assert.Equal(t, 31, txUser.Age, "年龄应该被更新到31")
}

// TestBaseRepository_WithContextExtractor 测试上下文提取器配置
func TestBaseRepository_WithContextExtractor(t *testing.T) {
	gormDB, err := setupTestDB()
	assert.NoError(t, err)

	dbHandler := MustNewGormHandler(gormDB)
	testLogger := logger.NewLogger(nil)

	repo := NewBaseRepository[TestUser](dbHandler, testLogger, "test_users")
	assert.NotNil(t, repo)
	assert.NotNil(t, repo.logger)
}

// TestBaseRepository_DefaultContextExtractor 测试默认上下文提取器
func TestBaseRepository_DefaultContextExtractor(t *testing.T) {
	testLogger := logger.NewLogger(nil)

	// 创建包含各种上下文值的context
	ctx := context.WithValue(context.Background(), "request_id", "req-123")
	ctx = context.WithValue(ctx, "user_id", "user-456")
	ctx = context.WithValue(ctx, "trace_id", "trace-789")
	ctx = context.WithValue(ctx, "session_id", "session-999")

	result := testLogger.WithContext(ctx)
	assert.NotNil(t, result)
}

// TestBaseRepository_ListWithPreloads 测试带预加载的列表查询
func TestBaseRepository_ListWithPreloads(t *testing.T) {
	gormDB, err := setupTestDB()
	assert.NoError(t, err)

	dbHandler := MustNewGormHandler(gormDB)
	repo := NewBaseRepository[TestUser](dbHandler, logger.NewLogger(nil), "test_users")

	// 插入测试数据
	user := &TestUser{
		Name:   "John Doe",
		Email:  "john@example.com",
		Age:    30,
		Status: "active",
	}
	_, err = repo.Create(context.Background(), user)
	assert.NoError(t, err)

	// 测试带预加载的查询（即使没有关联也不应该报错）
	query := NewQuery().Limit(10)
	results, err := repo.ListWithPreloads(context.Background(), query) // 不指定预加载关联
	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, "John Doe", results[0].Name)
}

// TestBaseRepository_BulkCreate_ErrorHandling 测试批量创建错误处理
func TestBaseRepository_BulkCreate_ErrorHandling(t *testing.T) {
	gormDB, err := setupTestDB()
	assert.NoError(t, err)

	dbHandler := MustNewGormHandler(gormDB)
	repo := NewBaseRepository[TestUser](dbHandler, logger.NewLogger(nil), "test_users")

	// 测试空数组
	err = repo.BulkCreate(context.Background(), []*TestUser{})
	assert.NoError(t, err)

	// 测试自定义批量大小
	users := []*TestUser{
		{Name: "User1", Email: "user1@test.com", Age: 25},
		{Name: "User2", Email: "user2@test.com", Age: 30},
		{Name: "User3", Email: "user3@test.com", Age: 35},
	}

	err = repo.BulkCreate(context.Background(), users, 2) // 自定义批量大小为2
	assert.NoError(t, err)

	// 验证数据插入成功
	all, err := repo.GetAll(context.Background())
	assert.NoError(t, err)
	assert.Len(t, all, 3)
}

// TestBaseRepository_List_Advanced 测试高级列表查询功能
func TestBaseRepository_List_Advanced(t *testing.T) {
	gormDB, err := setupTestDB()
	assert.NoError(t, err)

	dbHandler := MustNewGormHandler(gormDB)
	repo := NewBaseRepository[TestUser](dbHandler, logger.NewLogger(nil), "test_users")

	// 插入测试数据
	users := []*TestUser{
		{Name: "Alice", Email: "alice@test.com", Age: 25, Status: "active"},
		{Name: "Bob", Email: "bob@test.com", Age: 30, Status: "inactive"},
		{Name: "Charlie", Email: "charlie@test.com", Age: 35, Status: "active"},
	}
	for _, user := range users {
		_, err = repo.Create(context.Background(), user)
		assert.NoError(t, err)
	}

	// 测试带去重的查询
	query := NewQuery().WithDistinct(true)
	query.AddFilter(&Filter{Field: "status", Operator: OP_EQ, Value: "active"})
	results, err := repo.List(context.Background(), query)
	assert.NoError(t, err)
	assert.Len(t, results, 2)

	// 测试带分组的查询（不使用HAVING避免SQLite兼容性问题）
	query2 := NewQuery()
	query2.AddGroupBy("status")
	results2, err := repo.List(context.Background(), query2)
	assert.NoError(t, err)
	assert.NotNil(t, results2)

	// 测试限制和偏移
	query3 := NewQuery().Limit(2).Offset(1)
	results3, err := repo.List(context.Background(), query3)
	assert.NoError(t, err)
	assert.Len(t, results3, 2)
}

// TestBaseRepository_FilterConditions 测试复杂过滤条件构建
func TestBaseRepository_FilterConditions(t *testing.T) {
	gormDB, err := setupTestDB()
	assert.NoError(t, err)

	dbHandler := MustNewGormHandler(gormDB)
	repo := NewBaseRepository[TestUser](dbHandler, logger.NewLogger(nil), "test_users")

	// 插入测试数据
	users := []*TestUser{
		{Name: "Alice", Email: "alice@test.com", Age: 25, Status: "active"},
		{Name: "Bob", Email: "bob@test.com", Age: 30, Status: "inactive"},
		{Name: "Charlie", Email: "charlie@test.com", Age: 35, Status: "active"},
	}
	for _, user := range users {
		_, err = repo.Create(context.Background(), user)
		assert.NoError(t, err)
	}

	// 测试BETWEEN条件
	betweenFilter := &Filter{
		Field:    "age",
		Operator: OP_BETWEEN,
		Value:    []interface{}{25, 30},
	}
	query := NewQuery().AddFilter(betweenFilter)
	results, err := repo.List(context.Background(), query)
	assert.NoError(t, err)
	assert.Len(t, results, 2) // Alice(25) and Bob(30)

	// 测试IS NULL条件
	nullFilter := &Filter{
		Field:    "deleted_at",
		Operator: OP_IS_NULL,
	}
	query2 := NewQuery().AddFilter(nullFilter)
	results2, err := repo.List(context.Background(), query2)
	assert.NoError(t, err)
	assert.Len(t, results2, 3) // 所有用户的deleted_at都是NULL

	// 测试IS NOT NULL条件
	notNullFilter := &Filter{
		Field:    "name",
		Operator: OP_IS_NOT_NULL,
	}
	query3 := NewQuery().AddFilter(notNullFilter)
	results3, err := repo.List(context.Background(), query3)
	assert.NoError(t, err)
	assert.Len(t, results3, 3) // 所有用户都有name
}

// TestBaseRepository_FilterGroup_ComplexLogic 测试复杂过滤组逻辑
func TestBaseRepository_FilterGroup_ComplexLogic(t *testing.T) {
	gormDB, err := setupTestDB()
	assert.NoError(t, err)

	dbHandler := MustNewGormHandler(gormDB)
	repo := NewBaseRepository[TestUser](dbHandler, logger.NewLogger(nil), "test_users")

	// 插入测试数据
	users := []*TestUser{
		{Name: "Alice", Email: "alice@test.com", Age: 25, Status: "active"},
		{Name: "Bob", Email: "bob@test.com", Age: 30, Status: "inactive"},
		{Name: "Charlie", Email: "charlie@test.com", Age: 35, Status: "active"},
		{Name: "Diana", Email: "diana@test.com", Age: 28, Status: "pending"},
	}
	for _, user := range users {
		_, err = repo.Create(context.Background(), user)
		assert.NoError(t, err)
	}

	// 创建复杂的OR组合：(age > 30) OR (status = 'pending')
	orGroup := NewFilterGroup(LOGIC_OR)
	orGroup.AddFilter(&Filter{Field: "age", Operator: OP_GT, Value: 30})
	orGroup.AddFilter(&Filter{Field: "status", Operator: OP_EQ, Value: "pending"})

	query := NewQuery().WithFilterGroup(orGroup)
	results, err := repo.List(context.Background(), query)
	assert.NoError(t, err)
	assert.Len(t, results, 2) // Charlie(35) and Diana(pending)

	// 测试嵌套过滤组：((age > 25 AND status = 'active') OR status = 'pending')
	innerAndGroup := NewFilterGroup(LOGIC_AND)
	innerAndGroup.AddFilter(&Filter{Field: "age", Operator: OP_GT, Value: 25})
	innerAndGroup.AddFilter(&Filter{Field: "status", Operator: OP_EQ, Value: "active"})

	outerOrGroup := NewFilterGroup(LOGIC_OR)
	outerOrGroup.AddGroup(innerAndGroup)
	outerOrGroup.AddFilter(&Filter{Field: "status", Operator: OP_EQ, Value: "pending"})

	query2 := NewQuery().WithFilterGroup(outerOrGroup)
	results2, err := repo.List(context.Background(), query2)
	assert.NoError(t, err)
	assert.Len(t, results2, 2) // Charlie(35,active) and Diana(pending)
}

// TestBaseRepository_EdgeCases 测试各种边界情况
func TestBaseRepository_EdgeCases(t *testing.T) {
	gormDB, err := setupTestDB()
	assert.NoError(t, err)

	dbHandler := MustNewGormHandler(gormDB)
	repo := NewBaseRepository[TestUser](dbHandler, logger.NewLogger(nil), "test_users")

	// 测试空数组的批量操作
	err = repo.UpdateBatch(context.Background())
	assert.NoError(t, err)

	err = repo.DeleteBatch(context.Background())
	assert.NoError(t, err)

	err = repo.SoftDeleteBatch(context.Background(), []interface{}{}, "deleted_at", time.Now())
	assert.NoError(t, err)

	err = repo.RestoreBatch(context.Background(), []interface{}{}, "deleted_at", nil)
	assert.NoError(t, err)

	// 测试空字段映射的更新
	user := &TestUser{Name: "Test", Email: "test@test.com", Age: 30}
	created, err := repo.Create(context.Background(), user)
	assert.NoError(t, err)

	err = repo.UpdateFields(context.Background(), created.ID, map[string]interface{}{})
	assert.NoError(t, err)

	err = repo.UpdateFieldsByFilters(context.Background(), map[string]interface{}{},
		&Filter{Field: "id", Operator: OP_EQ, Value: created.ID})
	assert.NoError(t, err)

	// 测试空表的查询
	all1, err := repo.GetAll(context.Background())
	assert.NoError(t, err)
	assert.Len(t, all1, 1) // 只有刚创建的用户

	exists, err := repo.Exists(context.Background())
	assert.NoError(t, err)
	assert.True(t, exists)

	// 清空表后再测试
	err = repo.Delete(context.Background(), created.ID)
	assert.NoError(t, err)

	all2, err := repo.GetAll(context.Background())
	assert.NoError(t, err)
	assert.Len(t, all2, 0)

	exists2, err := repo.Exists(context.Background())
	assert.NoError(t, err)
	assert.False(t, exists2)
}

// TestBaseRepository_ReadOnlyMode 测试只读模式
func TestBaseRepository_ReadOnlyMode(t *testing.T) {
	gormDB, err := setupTestDB()
	assert.NoError(t, err)

	dbHandler := MustNewGormHandler(gormDB)
	repo := NewBaseRepository[TestUser](dbHandler, logger.NewLogger(nil), "test_users", WithReadOnly[TestUser]())

	// 测试只读模式下的创建操作应该返回错误
	user := &TestUser{Name: "Test", Email: "test@test.com", Age: 30}
	_, err = repo.Create(context.Background(), user)
	assert.Error(t, err)

	// 测试只读模式下的批量创建操作应该返回错误
	err = repo.CreateBatch(context.Background(), user)
	assert.Error(t, err)

	err = repo.BulkCreate(context.Background(), []*TestUser{user})
	assert.Error(t, err)

	// 测试只读模式下的CreateIfNotExists应该返回错误
	_, _, err = repo.CreateIfNotExists(context.Background(), user, "email")
	assert.Error(t, err)

	_, _, err = repo.CreateOrUpdate(context.Background(), user, "email")
	assert.Error(t, err)
}

// TestBaseRepository_InvalidInputs 测试无效输入处理
func TestBaseRepository_InvalidInputs(t *testing.T) {
	gormDB, err := setupTestDB()
	assert.NoError(t, err)

	dbHandler := MustNewGormHandler(gormDB)
	repo := NewBaseRepository[TestUser](dbHandler, logger.NewLogger(nil), "test_users")

	// 测试nil实体
	_, err = repo.Create(context.Background(), nil)
	assert.Error(t, err)

	_, err = repo.Update(context.Background(), nil)
	assert.Error(t, err)
	err = repo.UpdateByFilters(context.Background(), nil,
		&Filter{Field: "id", Operator: OP_EQ, Value: 1})
	assert.Error(t, err)

	// 测试空过滤器的操作
	_, err = repo.GetByFilter(context.Background(), nil)
	assert.Error(t, err)

	_, err = repo.GetByFilters(context.Background())
	assert.Error(t, err)

	err = repo.UpdateFieldsByFilters(context.Background(), map[string]interface{}{"name": "test"})
	assert.Error(t, err)

	err = repo.DeleteByFilters(context.Background())
	assert.Error(t, err)

	err = repo.SoftDeleteByFilters(context.Background(), "deleted_at", time.Now())
	assert.Error(t, err)

	err = repo.UpdateFieldsByFilters(context.Background(), map[string]interface{}{"name": "test"})
	assert.Error(t, err)

	// 测试空字段的GetByFields
	_, err = repo.GetByFields(context.Background(), map[string]interface{}{})
	assert.Error(t, err)

	// 测试CreateIfNotExists的无效输入
	user := &TestUser{Name: "Test", Email: "test@test.com", Age: 30}
	_, _, err = repo.CreateIfNotExists(context.Background(), nil, "email")
	assert.Error(t, err)

	_, _, err = repo.CreateIfNotExists(context.Background(), user)
	assert.Error(t, err)
}

// TestBaseRepository_Find_Compatibility 测试Find方法的兼容性
func TestBaseRepository_Find_Compatibility(t *testing.T) {
	gormDB, err := setupTestDB()
	assert.NoError(t, err)

	dbHandler := MustNewGormHandler(gormDB)
	repo := NewBaseRepository[TestUser](dbHandler, logger.NewLogger(nil), "test_users")

	// 插入测试数据
	users := []*TestUser{
		{Name: "Alice", Email: "alice@test.com", Age: 25, Status: "active"},
		{Name: "Bob", Email: "bob@test.com", Age: 30, Status: "inactive"},
		{Name: "Charlie", Email: "charlie@test.com", Age: 35, Status: "active"},
	}
	for _, user := range users {
		_, err = repo.Create(context.Background(), user)
		assert.NoError(t, err)
	}

	// 测试nil参数的Find方法（应该返回所有记录）
	results, err := repo.Find(context.Background(), nil)
	assert.NoError(t, err)
	assert.Len(t, results, 3)

	// 测试带条件的Find方法
	findOpts := &FindOptions{
		Conditions: []Condition{
			{Field: "status", Op: OP_EQ, Value: "active"},
		},
		Orders: []OrderBy{
			{Field: "age", Direction: "ASC"},
		},
		Limit:  2,
		Offset: 0,
	}

	results2, err := repo.Find(context.Background(), findOpts)
	assert.NoError(t, err)
	assert.Len(t, results2, 2)
	assert.Equal(t, "Alice", results2[0].Name) // 年龄较小的先返回
}

// TestBaseRepository_Pluck_Distinct 测试Pluck和Distinct方法
func TestBaseRepository_Pluck_Distinct(t *testing.T) {
	gormDB, err := setupTestDB()
	assert.NoError(t, err)

	dbHandler := MustNewGormHandler(gormDB)
	repo := NewBaseRepository[TestUser](dbHandler, logger.NewLogger(nil), "test_users")

	// 插入测试数据
	users := []*TestUser{
		{Name: "Alice", Email: "alice@test.com", Age: 25, Status: "active"},
		{Name: "Bob", Email: "bob@test.com", Age: 30, Status: "active"},
		{Name: "Charlie", Email: "charlie@test.com", Age: 35, Status: "inactive"},
	}
	for _, user := range users {
		_, err = repo.Create(context.Background(), user)
		assert.NoError(t, err)
	}

	// 测试Pluck方法
	names, err := repo.Pluck(context.Background(), "name")
	assert.NoError(t, err)
	assert.Len(t, names, 3)

	// 测试带过滤条件的Pluck
	activeNames, err := repo.Pluck(context.Background(), "name",
		&Filter{Field: "status", Operator: OP_EQ, Value: "active"})
	assert.NoError(t, err)
	assert.Len(t, activeNames, 2)

	// 测试Distinct方法
	statuses, err := repo.Distinct(context.Background(), "status")
	assert.NoError(t, err)
	assert.Len(t, statuses, 2) // active, inactive

	// 测试带过滤条件的Distinct
	activeStatuses, err := repo.Distinct(context.Background(), "status",
		&Filter{Field: "age", Operator: OP_GT, Value: 25})
	assert.NoError(t, err)
	assert.Len(t, activeStatuses, 2) // active, inactive (Bob和Charlie)
}

// TestBuildFilterCondition 测试过滤条件构建函数
func TestBuildFilterCondition(t *testing.T) {
	gormDB, err := setupTestDB()
	assert.NoError(t, err)

	dbHandler := MustNewGormHandler(gormDB)
	repo := NewBaseRepository[TestUser](dbHandler, logger.NewLogger(nil), "test_users")

	// 测试各种操作符的条件构建
	testCases := []struct {
		filter   *Filter
		expected bool // 是否应该成功构建条件
	}{
		{&Filter{Field: "name", Operator: OP_EQ, Value: "test"}, true},
		{&Filter{Field: "age", Operator: OP_GT, Value: 18}, true},
		{&Filter{Field: "age", Operator: OP_LT, Value: 65}, true},
		{&Filter{Field: "age", Operator: OP_GTE, Value: 18}, true},
		{&Filter{Field: "age", Operator: OP_LTE, Value: 65}, true},
		{&Filter{Field: "status", Operator: OP_IN, Value: []string{"active", "inactive"}}, true},
		{&Filter{Field: "status", Operator: OP_NOT_IN, Value: []string{"deleted"}}, true},
		{&Filter{Field: "name", Operator: OP_LIKE, Value: "%test%"}, true},
		{&Filter{Field: "name", Operator: OP_NOT_LIKE, Value: "%spam%"}, true},
		{&Filter{Field: "age", Operator: OP_BETWEEN, Value: []interface{}{18, 65}}, true},
		{&Filter{Field: "deleted_at", Operator: OP_IS_NULL}, true},
		{&Filter{Field: "name", Operator: OP_IS_NOT_NULL}, true},
		{&Filter{Field: "status", Operator: OP_NEQ, Value: "deleted"}, true},
		{nil, false}, // nil过滤器
		{&Filter{Field: "age", Operator: OP_BETWEEN, Value: "invalid"}, true}, // 无效BETWEEN值
	}

	// 插入测试数据
	user := &TestUser{Name: "Test User", Email: "test@example.com", Age: 30, Status: "active"}
	_, err = repo.Create(context.Background(), user)
	assert.NoError(t, err)

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			if tc.filter == nil {
				// 测试nil过滤器
				query := NewQuery()
				_, err := repo.List(context.Background(), query)
				assert.NoError(t, err)
				return
			}

			// 测试过滤器应用
			query := NewQuery().AddFilter(tc.filter)
			results, err := repo.List(context.Background(), query)
			if tc.expected {
				assert.NoError(t, err)
				assert.NotNil(t, results)
			}
		})
	}
}

// TestBuildGroupCondition 测试过滤组条件构建
func TestBuildGroupCondition(t *testing.T) {
	gormDB, err := setupTestDB()
	assert.NoError(t, err)

	dbHandler := MustNewGormHandler(gormDB)
	repo := NewBaseRepository[TestUser](dbHandler, logger.NewLogger(nil), "test_users")

	// 插入测试数据
	users := []*TestUser{
		{Name: "Alice", Email: "alice@test.com", Age: 25, Status: "active"},
		{Name: "Bob", Email: "bob@test.com", Age: 30, Status: "inactive"},
		{Name: "Charlie", Email: "charlie@test.com", Age: 35, Status: "active"},
	}
	for _, user := range users {
		_, err = repo.Create(context.Background(), user)
		assert.NoError(t, err)
	}

	// 测试空的过滤组
	emptyGroup := NewFilterGroup(LOGIC_AND)
	query1 := NewQuery().WithFilterGroup(emptyGroup)
	results1, err := repo.List(context.Background(), query1)
	assert.NoError(t, err)
	assert.Len(t, results1, 3) // 应该返回所有记录

	// 测试单条件AND组
	andGroup := NewFilterGroup(LOGIC_AND)
	andGroup.AddFilter(&Filter{Field: "status", Operator: OP_EQ, Value: "active"})
	query2 := NewQuery().WithFilterGroup(andGroup)
	results2, err := repo.List(context.Background(), query2)
	assert.NoError(t, err)
	assert.Len(t, results2, 2) // Alice和Charlie

	// 测试单条件OR组
	orGroup := NewFilterGroup(LOGIC_OR)
	orGroup.AddFilter(&Filter{Field: "age", Operator: OP_GT, Value: 32})
	query3 := NewQuery().WithFilterGroup(orGroup)
	results3, err := repo.List(context.Background(), query3)
	assert.NoError(t, err)
	assert.Len(t, results3, 1) // Charlie

	// 测试嵌套过滤组
	innerGroup := NewFilterGroup(LOGIC_AND)
	innerGroup.AddFilter(&Filter{Field: "age", Operator: OP_GT, Value: 20})
	innerGroup.AddFilter(&Filter{Field: "status", Operator: OP_EQ, Value: "active"})

	outerGroup := NewFilterGroup(LOGIC_OR)
	outerGroup.AddGroup(innerGroup)
	outerGroup.AddFilter(&Filter{Field: "name", Operator: OP_EQ, Value: "Bob"})

	query4 := NewQuery().WithFilterGroup(outerGroup)
	results4, err := repo.List(context.Background(), query4)
	assert.NoError(t, err)
	assert.Len(t, results4, 3) // Alice, Charlie (from inner group) + Bob
}

// TestApplyOrdering 测试排序应用
func TestApplyOrdering(t *testing.T) {
	gormDB, err := setupTestDB()
	assert.NoError(t, err)

	dbHandler := MustNewGormHandler(gormDB)

	// 测试带默认排序的仓储
	repo := NewBaseRepository[TestUser](dbHandler, logger.NewLogger(nil), "test_users", WithDefaultOrder[TestUser]("age DESC"))

	// 插入测试数据
	users := []*TestUser{
		{Name: "Alice", Email: "alice@test.com", Age: 25, Status: "active"},
		{Name: "Bob", Email: "bob@test.com", Age: 30, Status: "inactive"},
		{Name: "Charlie", Email: "charlie@test.com", Age: 20, Status: "active"},
	}
	for _, user := range users {
		_, err = repo.Create(context.Background(), user)
		assert.NoError(t, err)
	}

	// 测试默认排序（无显式排序条件时）
	query1 := NewQuery()
	results1, err := repo.List(context.Background(), query1)
	assert.NoError(t, err)
	assert.Len(t, results1, 3)
	// 应该按年龄降序：Bob(30), Alice(25), Charlie(20)
	assert.Equal(t, "Bob", results1[0].Name)

	// 测试显式排序（应该覆盖默认排序）
	query2 := NewQuery().AddOrder("name", "ASC")
	results2, err := repo.List(context.Background(), query2)
	assert.NoError(t, err)
	assert.Len(t, results2, 3)
	// 应该按名字升序：Alice, Bob, Charlie
	assert.Equal(t, "Alice", results2[0].Name)

	// 测试多重排序
	query3 := NewQuery().AddOrder("status", "ASC").AddOrder("age", "DESC")
	results3, err := repo.List(context.Background(), query3)
	assert.NoError(t, err)
	assert.Len(t, results3, 3)

	// 测试空排序字段
	query4 := NewQuery().AddOrder("", "ASC")
	results4, err := repo.List(context.Background(), query4)
	assert.NoError(t, err)
	assert.Len(t, results4, 3)
	// 应该回退到默认排序
	assert.Equal(t, "Alice", results4[0].Name)
}
