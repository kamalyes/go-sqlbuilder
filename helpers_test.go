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
