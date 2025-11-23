/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-23 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-23 23:55:48
 * @FilePath: \go-sqlbuilder\helpers_test.go
 * @Description: 仓储辅助工具 - 软删除、查询辅助等功能
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package sqlbuilder

import (
	"context"
	"fmt"
	"github.com/kamalyes/go-logger"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// TestSoftDeleteHelpers 测试软删除帮助函数
func TestSoftDeleteHelpers(t *testing.T) {
	gormDB, err := setupTestDB()
	assert.NoError(t, err)

	dbHandler := MustNewGormHandler(gormDB)
	repo := NewBaseRepository[TestUser](dbHandler, logger.NewLogger(nil), "test_users")

	ctx := context.Background()

	// 插入测试数据
	users := []*TestUser{
		{Name: "Alice", Email: "alice@test.com", Age: 25, Status: "active"},
		{Name: "Bob", Email: "bob@test.com", Age: 30, Status: "active"},
		{Name: "Charlie", Email: "charlie@test.com", Age: 35, Status: "active"},
	}
	var userIDs []uint
	for _, user := range users {
		createdUser, err := repo.Create(ctx, user)
		assert.NoError(t, err)
		userIDs = append(userIDs, createdUser.ID)
	}

	// 软删除一些用户（直接操作数据库模拟）
	now := time.Now()
	err = gormDB.Model(&TestUser{}).Where("id IN ?", userIDs[:2]).Update("deleted_at", now).Error
	assert.NoError(t, err)

	// 测试 GetDeleted - 获取已删除记录
	deletedUsers, err := GetDeleted[TestUser](ctx, gormDB, &Query{})
	assert.NoError(t, err)
	assert.Len(t, deletedUsers, 2) // Alice 和 Bob

	// 测试 GetNonDeleted - 获取未删除记录
	nonDeletedUsers, err := GetNonDeleted[TestUser](ctx, gormDB, &Query{})
	assert.NoError(t, err)
	assert.Len(t, nonDeletedUsers, 1) // Charlie

	// 测试带过滤条件的已删除记录
	queryWithFilter := &Query{}
	queryWithFilter.AddFilter(&Filter{Field: "name", Operator: "=", Value: "Alice"})
	filteredDeleted, err := GetDeleted[TestUser](ctx, gormDB, queryWithFilter)
	assert.NoError(t, err)
	assert.Len(t, filteredDeleted, 1)
	assert.Equal(t, "Alice", filteredDeleted[0].Name)

	// 测试带过滤条件的未删除记录
	queryWithFilter2 := &Query{}
	queryWithFilter2.AddFilter(&Filter{Field: "age", Operator: ">", Value: 30})
	filteredNonDeleted, err := GetNonDeleted[TestUser](ctx, gormDB, queryWithFilter2)
	assert.NoError(t, err)
	assert.Len(t, filteredNonDeleted, 1)
	assert.Equal(t, "Charlie", filteredNonDeleted[0].Name)
}

// TestRestoreDeleted 测试恢复已删除记录
func TestRestoreDeleted(t *testing.T) {
	gormDB, err := setupTestDB()
	assert.NoError(t, err)

	dbHandler := MustNewGormHandler(gormDB)
	repo := NewBaseRepository[TestUser](dbHandler, logger.NewLogger(nil), "test_users")

	ctx := context.Background()

	// 创建用户
	user := &TestUser{Name: "David", Email: "david@test.com", Age: 28, Status: "active"}
	createdUser, err := repo.Create(ctx, user)
	assert.NoError(t, err)
	userID := createdUser.ID

	// 软删除用户
	now := time.Now()
	err = gormDB.Model(&TestUser{}).Where("id = ?", userID).Update("deleted_at", now).Error
	assert.NoError(t, err)

	// 验证用户已被软删除
	deletedUsers, err := GetDeleted[TestUser](ctx, gormDB, &Query{})
	assert.NoError(t, err)
	assert.Len(t, deletedUsers, 1)

	// 恢复用户
	err = RestoreDeleted[TestUser](ctx, gormDB, userID)
	assert.NoError(t, err)

	// 验证用户已被恢复
	restoredUser, err := repo.Get(ctx, userID)
	assert.NoError(t, err)
	assert.Equal(t, "David", restoredUser.Name)
	assert.Nil(t, restoredUser.DeletedAt) // 删除时间应该被清空

	// 验证已删除列表为空
	deletedUsersAfter, err := GetDeleted[TestUser](ctx, gormDB, &Query{})
	assert.NoError(t, err)
	assert.Len(t, deletedUsersAfter, 0)
}

// TestRestoreDeletedBatch 测试批量恢复已删除记录
func TestRestoreDeletedBatch(t *testing.T) {
	gormDB, err := setupTestDB()
	assert.NoError(t, err)

	dbHandler := MustNewGormHandler(gormDB)

	repo := NewBaseRepository[TestUser](dbHandler, logger.NewLogger(nil), "test_users")

	ctx := context.Background()

	// 创建多个用户
	users := []*TestUser{
		{Name: "Eva", Email: "eva@test.com", Age: 26, Status: "active"},
		{Name: "Frank", Email: "frank@test.com", Age: 32, Status: "active"},
		{Name: "Grace", Email: "grace@test.com", Age: 29, Status: "active"},
	}
	var userIDs []uint
	for _, user := range users {
		createdUser, err := repo.Create(ctx, user)
		assert.NoError(t, err)
		userIDs = append(userIDs, createdUser.ID)
	}

	// 软删除前两个用户
	now := time.Now()
	err = gormDB.Model(&TestUser{}).Where("id IN ?", userIDs[:2]).Update("deleted_at", now).Error
	assert.NoError(t, err)

	// 验证有2个被删除的用户
	deletedUsers, err := GetDeleted[TestUser](ctx, gormDB, &Query{})
	assert.NoError(t, err)
	assert.Len(t, deletedUsers, 2)

	// 批量恢复
	var idsToRestore []interface{}
	for _, id := range userIDs[:2] {
		idsToRestore = append(idsToRestore, id)
	}
	err = RestoreDeletedBatch[TestUser](ctx, gormDB, idsToRestore)
	assert.NoError(t, err)

	// 验证所有用户都已恢复
	for _, userID := range userIDs[:2] {
		restoredUser, err := repo.Get(ctx, userID)
		assert.NoError(t, err)
		assert.Nil(t, restoredUser.DeletedAt)
	}

	// 验证已删除列表为空
	deletedUsersAfter, err := GetDeleted[TestUser](ctx, gormDB, &Query{})
	assert.NoError(t, err)
	assert.Len(t, deletedUsersAfter, 0)
}

// TestPermanentlyDelete 测试永久删除
func TestPermanentlyDelete(t *testing.T) {
	gormDB, err := setupTestDB()
	assert.NoError(t, err)

	dbHandler := MustNewGormHandler(gormDB)

	repo := NewBaseRepository[TestUser](dbHandler, logger.NewLogger(nil), "test_users")

	ctx := context.Background()

	// 创建用户
	user := &TestUser{Name: "Henry", Email: "henry@test.com", Age: 33, Status: "active"}
	createdUser, err := repo.Create(ctx, user)
	assert.NoError(t, err)
	userID := createdUser.ID

	// 先软删除用户
	now := time.Now()
	err = gormDB.Model(&TestUser{}).Where("id = ?", userID).Update("deleted_at", now).Error
	assert.NoError(t, err)

	// 验证用户在已删除列表中
	deletedUsers, err := GetDeleted[TestUser](ctx, gormDB, &Query{})
	assert.NoError(t, err)
	assert.Len(t, deletedUsers, 1)

	// 永久删除
	err = PermanentlyDelete[TestUser](ctx, gormDB, userID)
	assert.NoError(t, err)

	// 验证用户不再存在于已删除列表中
	deletedUsersAfter, err := GetDeleted[TestUser](ctx, gormDB, &Query{})
	assert.NoError(t, err)
	assert.Len(t, deletedUsersAfter, 0)

	// 验证用户也不在正常列表中
	_, err = repo.Get(ctx, userID)
	assert.Error(t, err)
}

// TestPermanentlyDeleteBatch 测试批量永久删除
func TestPermanentlyDeleteBatch(t *testing.T) {
	gormDB, err := setupTestDB()
	assert.NoError(t, err)

	dbHandler := MustNewGormHandler(gormDB)

	repo := NewBaseRepository[TestUser](dbHandler, logger.NewLogger(nil), "test_users")

	ctx := context.Background()

	// 创建多个用户
	users := []*TestUser{
		{Name: "Ivy", Email: "ivy@test.com", Age: 24, Status: "active"},
		{Name: "Jack", Email: "jack@test.com", Age: 31, Status: "active"},
	}
	var userIDs []uint
	for _, user := range users {
		createdUser, err := repo.Create(ctx, user)
		assert.NoError(t, err)
		userIDs = append(userIDs, createdUser.ID)
	}

	// 软删除所有用户
	now := time.Now()
	err = gormDB.Model(&TestUser{}).Where("id IN ?", userIDs).Update("deleted_at", now).Error
	assert.NoError(t, err)

	// 验证有2个被删除的用户
	deletedUsers, err := GetDeleted[TestUser](ctx, gormDB, &Query{})
	assert.NoError(t, err)
	assert.Len(t, deletedUsers, 2)

	// 批量永久删除
	var idsToDelete []interface{}
	for _, id := range userIDs {
		idsToDelete = append(idsToDelete, id)
	}
	err = PermanentlyDeleteBatch[TestUser](ctx, gormDB, idsToDelete)
	assert.NoError(t, err)

	// 验证已删除列表为空
	deletedUsersAfter, err := GetDeleted[TestUser](ctx, gormDB, &Query{})
	assert.NoError(t, err)
	assert.Len(t, deletedUsersAfter, 0)

	// 验证用户完全不存在
	for _, userID := range userIDs {
		_, err = repo.Get(ctx, userID)
		assert.Error(t, err)
	}
}

// TestSoftDeleteHelpersWithPagination 测试带分页的软删除帮助函数
func TestSoftDeleteHelpersWithPagination(t *testing.T) {
	gormDB, err := setupTestDB()
	assert.NoError(t, err)

	dbHandler := MustNewGormHandler(gormDB)

	repo := NewBaseRepository[TestUser](dbHandler, logger.NewLogger(nil), "test_users")

	ctx := context.Background()

	// 插入大量测试数据
	users := []*TestUser{}
	for i := 1; i <= 15; i++ {
		user := &TestUser{
			Name:   fmt.Sprintf("User%d", i),
			Email:  fmt.Sprintf("user%d@test.com", i),
			Age:    20 + i,
			Status: "active",
		}
		users = append(users, user)
	}

	var allUserIDs []uint
	for _, user := range users {
		createdUser, err := repo.Create(ctx, user)
		assert.NoError(t, err)
		allUserIDs = append(allUserIDs, createdUser.ID)
	}

	// 软删除前10个用户
	now := time.Now()
	err = gormDB.Model(&TestUser{}).Where("id IN ?", allUserIDs[:10]).Update("deleted_at", now).Error
	assert.NoError(t, err)

	// 测试带分页的已删除记录查询
	queryWithPagination := &Query{}
	queryWithPagination.WithPaging(1, 5) // 第1页，每页5条

	deletedPage1, err := GetDeleted[TestUser](ctx, gormDB, queryWithPagination)
	assert.NoError(t, err)
	assert.Len(t, deletedPage1, 5)

	// 测试第2页
	queryWithPagination.WithPaging(2, 5) // 第2页，每页5条
	deletedPage2, err := GetDeleted[TestUser](ctx, gormDB, queryWithPagination)
	assert.NoError(t, err)
	assert.Len(t, deletedPage2, 5)

	// 测试带分页的未删除记录查询
	queryWithPagination2 := &Query{}
	queryWithPagination2.WithPaging(1, 3) // 第1页，每页3条

	nonDeletedPage1, err := GetNonDeleted[TestUser](ctx, gormDB, queryWithPagination2)
	assert.NoError(t, err)
	assert.Len(t, nonDeletedPage1, 3)

	// 第2页应该有2条记录（总共5条未删除，第1页3条，第2页2条）
	queryWithPagination2.WithPaging(2, 3)
	nonDeletedPage2, err := GetNonDeleted[TestUser](ctx, gormDB, queryWithPagination2)
	assert.NoError(t, err)
	assert.Len(t, nonDeletedPage2, 2)
}

// TestSoftDeleteHelpersWithOrdering 测试带排序的软删除帮助函数
func TestSoftDeleteHelpersWithOrdering(t *testing.T) {
	gormDB, err := setupTestDB()
	assert.NoError(t, err)

	dbHandler := MustNewGormHandler(gormDB)

	repo := NewBaseRepository[TestUser](dbHandler, logger.NewLogger(nil), "test_users")

	ctx := context.Background()

	// 插入测试数据
	users := []*TestUser{
		{Name: "Zoe", Email: "zoe@test.com", Age: 30, Status: "active"},
		{Name: "Alice", Email: "alice@test.com", Age: 25, Status: "active"},
		{Name: "Bob", Email: "bob@test.com", Age: 35, Status: "active"},
	}

	var userIDs []uint
	for _, user := range users {
		createdUser, err := repo.Create(ctx, user)
		assert.NoError(t, err)
		userIDs = append(userIDs, createdUser.ID)
	}

	// 软删除所有用户
	now := time.Now()
	err = gormDB.Model(&TestUser{}).Where("id IN ?", userIDs).Update("deleted_at", now).Error
	assert.NoError(t, err)

	// 测试按名字升序排序的已删除记录
	queryWithOrder := &Query{}
	queryWithOrder.AddOrder("name", "ASC")

	deletedOrderedByName, err := GetDeleted[TestUser](ctx, gormDB, queryWithOrder)
	assert.NoError(t, err)
	assert.Len(t, deletedOrderedByName, 3)
	assert.Equal(t, "Alice", deletedOrderedByName[0].Name)
	assert.Equal(t, "Bob", deletedOrderedByName[1].Name)
	assert.Equal(t, "Zoe", deletedOrderedByName[2].Name)

	// 测试按年龄降序排序的已删除记录
	queryWithOrder2 := &Query{}
	queryWithOrder2.AddOrder("age", "DESC")

	deletedOrderedByAge, err := GetDeleted[TestUser](ctx, gormDB, queryWithOrder2)
	assert.NoError(t, err)
	assert.Len(t, deletedOrderedByAge, 3)
	assert.Equal(t, "Bob", deletedOrderedByAge[0].Name)   // 35岁
	assert.Equal(t, "Zoe", deletedOrderedByAge[1].Name)   // 30岁
	assert.Equal(t, "Alice", deletedOrderedByAge[2].Name) // 25岁
}

// TestRepositoryWithSoftDelete_DeletedAt 测试使用 deleted_at 的软删除仓储
func TestRepositoryWithSoftDelete_DeletedAt(t *testing.T) {
	gormDB, err := setupTestDB()
	assert.NoError(t, err)

	dbHandler := MustNewGormHandler(gormDB)
	baseRepo := NewBaseRepository[TestUser](dbHandler, logger.NewLogger(nil), "test_users")

	// 创建软删除仓储
	repo := NewRepositoryWithSoftDelete(baseRepo)
	ctx := context.Background()

	// 创建测试数据
	users := []*TestUser{
		{Name: "User1", Email: "user1@soft.com", Age: 25, Status: "active"},
		{Name: "User2", Email: "user2@soft.com", Age: 30, Status: "active"},
		{Name: "User3", Email: "user3@soft.com", Age: 35, Status: "active"},
	}

	var userIDs []uint
	for _, user := range users {
		createdUser, err := baseRepo.Create(ctx, user)
		assert.NoError(t, err)
		userIDs = append(userIDs, createdUser.ID)
	}

	// 测试单个软删除
	err = repo.SoftDeleteWithDeletedAt(ctx, userIDs[0])
	assert.NoError(t, err)

	// 验证软删除成功
	deleted, err := GetDeleted[TestUser](ctx, gormDB, nil)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(deleted))
	assert.Equal(t, "User1", deleted[0].Name)

	// 测试恢复
	err = repo.RestoreWithDeletedAt(ctx, userIDs[0])
	assert.NoError(t, err)

	// 验证恢复成功
	deleted, err = GetDeleted[TestUser](ctx, gormDB, nil)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(deleted))

	// 测试批量软删除
	ids := []interface{}{userIDs[0], userIDs[1]}
	err = repo.SoftDeleteBatchWithDeletedAt(ctx, ids)
	assert.NoError(t, err)

	// 验证批量软删除成功
	deleted, err = GetDeleted[TestUser](ctx, gormDB, nil)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(deleted))

	// 测试批量恢复
	err = repo.RestoreBatchWithDeletedAt(ctx, ids)
	assert.NoError(t, err)

	// 验证批量恢复成功
	deleted, err = GetDeleted[TestUser](ctx, gormDB, nil)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(deleted))

	// 测试按条件软删除
	filters := []*Filter{
		{Field: "name", Operator: OP_EQ, Value: "User1"},
	}
	err = repo.SoftDeleteByFiltersWithDeletedAt(ctx, filters...)
	assert.NoError(t, err)

	// 验证按条件软删除成功
	deleted, err = GetDeleted[TestUser](ctx, gormDB, nil)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(deleted))
	assert.Equal(t, "User1", deleted[0].Name)
}

// TestRepositoryWithSoftDelete_ListMethods 测试列表查询方法
func TestRepositoryWithSoftDelete_ListMethods(t *testing.T) {
	gormDB, err := setupTestDB()
	assert.NoError(t, err)

	dbHandler := MustNewGormHandler(gormDB)
	baseRepo := NewBaseRepository[TestUser](dbHandler, logger.NewLogger(nil), "test_users")

	// 创建软删除仓储
	repo := NewRepositoryWithSoftDelete(baseRepo)
	ctx := context.Background()

	// 创建测试数据
	users := []*TestUser{
		{Name: "ActiveUser1", Email: "active1@list.com", Age: 25, Status: "active"},
		{Name: "ActiveUser2", Email: "active2@list.com", Age: 30, Status: "active"},
		{Name: "DeletedUser", Email: "deleted@list.com", Age: 35, Status: "inactive"},
	}

	var userIDs []uint
	for _, user := range users {
		createdUser, err := baseRepo.Create(ctx, user)
		assert.NoError(t, err)
		userIDs = append(userIDs, createdUser.ID)
	}

	// 软删除第三个用户
	err = repo.SoftDeleteWithDeletedAt(ctx, userIDs[2])
	assert.NoError(t, err)

	// 测试 ListNotDeleted
	query := NewQuery().AddOrder("id", "ASC")
	notDeleted, err := repo.ListNotDeleted(ctx, query)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(notDeleted), "应该有2个未删除的用户")
	if len(notDeleted) >= 2 {
		assert.Equal(t, "ActiveUser1", notDeleted[0].Name)
		assert.Equal(t, "ActiveUser2", notDeleted[1].Name)
	}

	// 测试 ListDeleted
	deletedQuery := NewQuery().AddOrder("id", "ASC")
	deleted, err := repo.ListDeleted(ctx, deletedQuery)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(deleted), 1, "应该至少有1个已删除的用户")
	if len(deleted) > 0 {
		assert.Equal(t, "DeletedUser", deleted[0].Name)
	}

	// 测试带过滤条件的 ListNotDeleted
	queryWithFilter := NewQuery().
		AddFilter(&Filter{
			Field:    "age",
			Operator: OP_GT,
			Value:    25,
		}).
		AddOrder("id", "ASC")
	notDeletedFiltered, err := repo.ListNotDeleted(ctx, queryWithFilter)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(notDeletedFiltered), 1, "应该至少有1个符合条件的用户")
	if len(notDeletedFiltered) > 0 {
		assert.Equal(t, "ActiveUser2", notDeletedFiltered[0].Name)
	}
}

// TestNewRepositoryWithSoftDelete 测试创建软删除仓储
func TestNewRepositoryWithSoftDelete(t *testing.T) {
	gormDB, err := setupTestDB()
	assert.NoError(t, err)

	dbHandler := MustNewGormHandler(gormDB)
	baseRepo := NewBaseRepository[TestUser](dbHandler, logger.NewLogger(nil), "test_users")

	// 测试创建软删除仓储
	repo := NewRepositoryWithSoftDelete(baseRepo)
	assert.NotNil(t, repo)
	assert.NotNil(t, repo.BaseRepository)
	assert.Equal(t, baseRepo, repo.BaseRepository)
}

// TestGormHandler_IsConnected 测试数据库连接检测
func TestGormHandler_IsConnected(t *testing.T) {
	gormDB, err := setupTestDB()
	assert.NoError(t, err)

	dbHandler := MustNewGormHandler(gormDB)

	// 测试正常连接
	assert.True(t, dbHandler.IsConnected())

	// 测试 nil DB
	nilHandler := &GormHandler{db: nil}
	assert.False(t, nilHandler.IsConnected())
}

// TestRepositoryWithSoftDelete_IsDeletedMethods 测试 is_deleted 字段的软删除方法
func TestRepositoryWithSoftDelete_IsDeletedMethods(t *testing.T) {
	gormDB, err := setupTestDB()
	assert.NoError(t, err)

	// 自动迁移创建表
	err = gormDB.AutoMigrate(&TestUserWithIsDeleted{})
	assert.NoError(t, err)

	dbHandler := MustNewGormHandler(gormDB)
	baseRepo := NewBaseRepository[TestUserWithIsDeleted](dbHandler, logger.NewLogger(nil), "test_users_is_deleted")
	repo := NewRepositoryWithSoftDelete(baseRepo)

	ctx := context.Background()

	// 创建测试数据
	user1 := &TestUserWithIsDeleted{Name: "User1", Email: "user1@test.com"}
	user2 := &TestUserWithIsDeleted{Name: "User2", Email: "user2@test.com"}
	user3 := &TestUserWithIsDeleted{Name: "User3", Email: "user3@test.com"}

	_, err = repo.Create(ctx, user1)
	assert.NoError(t, err)
	_, err = repo.Create(ctx, user2)
	assert.NoError(t, err)
	_, err = repo.Create(ctx, user3)
	assert.NoError(t, err)

	// 测试 SoftDeleteWithIsDeleted
	err = repo.SoftDeleteWithIsDeleted(ctx, user1.ID)
	assert.NoError(t, err)

	found, err := repo.BaseRepository.FindOne(ctx, NewEqFilter("id", user1.ID))
	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, 1, found.IsDeleted)

	// 测试 SoftDeleteBatchWithIsDeleted
	err = repo.SoftDeleteBatchWithIsDeleted(ctx, []interface{}{user2.ID})
	assert.NoError(t, err)

	found, err = repo.BaseRepository.FindOne(ctx, NewEqFilter("id", user2.ID))
	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, 1, found.IsDeleted)

	// 测试 ListDeletedByIsDeleted
	query := NewQuery()
	deleted, err := repo.ListDeletedByIsDeleted(ctx, query)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(deleted), 2)

	// 测试 ListNotDeletedByIsDeleted (新的 Query 实例)
	query2 := NewQuery()
	notDeleted, err := repo.ListNotDeletedByIsDeleted(ctx, query2)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(notDeleted)) // user3 未删除

	// 测试 RestoreWithIsDeleted
	err = repo.RestoreWithIsDeleted(ctx, user1.ID)
	assert.NoError(t, err)

	found, err = repo.BaseRepository.FindOne(ctx, NewEqFilter("id", user1.ID))
	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, 0, found.IsDeleted)

	// 测试 RestoreBatchWithIsDeleted
	err = repo.RestoreBatchWithIsDeleted(ctx, []interface{}{user2.ID})
	assert.NoError(t, err)

	found, err = repo.BaseRepository.FindOne(ctx, NewEqFilter("id", user2.ID))
	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, 0, found.IsDeleted)

	// 测试 SoftDeleteByFiltersWithIsDeleted
	filters := []*Filter{NewEqFilter("name", "User3")}
	err = repo.SoftDeleteByFiltersWithIsDeleted(ctx, filters...)
	assert.NoError(t, err)

	found, err = repo.BaseRepository.FindOne(ctx, NewEqFilter("id", user3.ID))
	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, 1, found.IsDeleted)
}

// TestUserWithIsDeleted is_deleted 字段的测试模型
type TestUserWithIsDeleted struct {
	ID        uint   `gorm:"primarykey"`
	Name      string `gorm:"size:100;not null"`
	Email     string `gorm:"size:100"`
	IsDeleted int    `gorm:"default:0;index"` // 0=未删除, 1=已删除
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (TestUserWithIsDeleted) TableName() string {
	return "test_users_is_deleted"
}

// TestRepositoryWithSoftDelete_ListDeletedAndNotDeleted 测试 ListDeleted 和 ListNotDeleted 方法
func TestRepositoryWithSoftDelete_ListDeletedAndNotDeleted(t *testing.T) {
	gormDB, err := setupTestDB()
	assert.NoError(t, err)

	dbHandler := MustNewGormHandler(gormDB)
	baseRepo := NewBaseRepository[TestUser](dbHandler, logger.NewLogger(nil), "test_users")
	repo := NewRepositoryWithSoftDelete(baseRepo)

	ctx := context.Background()

	// 创建测试数据
	user1 := &TestUser{Name: "User1", Email: "user1@test.com", Age: 25}
	user2 := &TestUser{Name: "User2", Email: "user2@test.com", Age: 30}

	_, err = repo.Create(ctx, user1)
	assert.NoError(t, err)
	_, err = repo.Create(ctx, user2)
	assert.NoError(t, err)

	// 软删除一个用户
	err = repo.SoftDeleteWithDeletedAt(ctx, user1.ID)
	assert.NoError(t, err)

	// 测试 ListDeleted (nil query)
	deleted, err := repo.ListDeleted(ctx, nil)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(deleted), 1)

	// 测试 ListNotDeleted (nil query)
	notDeleted, err := repo.ListNotDeleted(ctx, nil)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(notDeleted), 1)
}
