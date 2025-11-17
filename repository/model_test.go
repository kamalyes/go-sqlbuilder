/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-17 16:30:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-17 16:30:00
 * @FilePath: \go-sqlbuilder\repository\model_test.go
 * @Description: 模型测试文件
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package repository

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// User 测试用户模型
type User struct {
	BaseModel
	Name  string `json:"name" gorm:"type:varchar(100);not null;comment:用户名"`
	Email string `json:"email" gorm:"type:varchar(100);uniqueIndex;comment:邮箱"`
	Age   int    `json:"age" gorm:"comment:年龄"`
}

// Product 测试产品模型
type Product struct {
	LightModel
	Name  string  `json:"name" gorm:"type:varchar(200);not null;comment:产品名称"`
	Price float64 `json:"price" gorm:"type:decimal(10,2);comment:价格"`
}

// AuditUser 测试审计用户模型（包含审计字段）
type AuditUser struct {
	AuditModel
	Name  string `json:"name" gorm:"type:varchar(100);not null;comment:用户名"`
	Email string `json:"email" gorm:"type:varchar(100);uniqueIndex;comment:邮箱"`
}

// ModelTestSuite 模型测试套件
type ModelTestSuite struct {
	suite.Suite
	db *gorm.DB
}

// SetupSuite 测试套件初始化
func (s *ModelTestSuite) SetupSuite() {
	// 使用内存数据库
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(s.T(), err)
	s.db = db

	// 自动迁移
	err = db.AutoMigrate(&User{}, &Product{}, &AuditUser{})
	assert.NoError(s.T(), err)
}

// TearDownSuite 测试套件清理
func (s *ModelTestSuite) TearDownSuite() {
	sqlDB, _ := s.db.DB()
	if sqlDB != nil {
		sqlDB.Close()
	}
}

// TestBaseModel_SettersAndGetters 测试基础模型的设置和获取方法
func (s *ModelTestSuite) TestBaseModel_SettersAndGetters() {
	user := &User{
		Name:  "张三",
		Email: "zhangsan@example.com",
		Age:   25,
	}

	// 测试设置备注
	user.SetRemark("测试用户")
	assert.Equal(s.T(), "测试用户", user.Remark)

	// 测试启用/禁用
	user.Enable()
	assert.True(s.T(), user.IsEnabled())
	assert.Equal(s.T(), int8(1), user.Status)

	user.Disable()
	assert.False(s.T(), user.IsEnabled())
	assert.Equal(s.T(), int8(0), user.Status)

	// 测试时间设置
	now := time.Now()
	user.SetCreatedAt(now)
	user.SetUpdatedAt(now)
	assert.Equal(s.T(), now, user.CreatedAt)
	assert.Equal(s.T(), now, user.UpdatedAt)

	// 测试新记录判断
	assert.True(s.T(), user.IsNew())
	assert.Equal(s.T(), uint(0), user.GetID())
}

// TestBaseModel_Create 测试创建记录
func (s *ModelTestSuite) TestBaseModel_Create() {
	user := &User{
		Name:  "李四",
		Email: "lisi@example.com",
		Age:   30,
	}
	user.SetRemark("新用户")
	user.Enable()

	// 创建记录
	result := s.db.Create(user)
	assert.NoError(s.T(), result.Error)
	assert.NotZero(s.T(), user.ID)
	assert.False(s.T(), user.IsNew())
	assert.Equal(s.T(), 1, user.GetVersion())
	assert.True(s.T(), user.IsEnabled())
}

// TestBaseModel_Update 测试更新记录（版本号自增）
func (s *ModelTestSuite) TestBaseModel_Update() {
	// 创建用户
	user := &User{
		Name:  "王五",
		Email: "wangwu@example.com",
		Age:   28,
	}
	s.db.Create(user)

	// 记录初始版本号
	initialVersion := user.Version
	assert.Equal(s.T(), 1, initialVersion)

	// 更新用户
	user.Age = 29
	result := s.db.Save(user)
	assert.NoError(s.T(), result.Error)

	// 验证版本号自增
	assert.Equal(s.T(), initialVersion+1, user.Version)
	assert.Equal(s.T(), 2, user.GetVersion())
}

// TestBaseModel_SoftDelete 测试软删除
func (s *ModelTestSuite) TestBaseModel_SoftDelete() {
	// 创建用户
	user := &User{
		Name:  "赵六",
		Email: "zhaoliu@example.com",
		Age:   35,
	}
	s.db.Create(user)

	userID := user.ID
	assert.False(s.T(), user.IsDeleted())

	// 软删除
	result := s.db.Delete(user)
	assert.NoError(s.T(), result.Error)

	// 验证软删除
	var deletedUser User
	err := s.db.Unscoped().First(&deletedUser, userID).Error
	assert.NoError(s.T(), err)
	assert.True(s.T(), deletedUser.IsDeleted())
	assert.True(s.T(), deletedUser.DeletedAt.Valid)

	// 普通查询应该找不到
	err = s.db.First(&User{}, userID).Error
	assert.Error(s.T(), err)
	assert.Equal(s.T(), gorm.ErrRecordNotFound, err)
}

// TestBaseModel_Restore 测试恢复软删除
func (s *ModelTestSuite) TestBaseModel_Restore() {
	// 创建并删除用户
	user := &User{
		Name:  "孙七",
		Email: "sunqi@example.com",
		Age:   40,
	}
	s.db.Create(user)
	userID := user.ID
	s.db.Delete(user)

	// 验证已删除
	var deletedUser User
	s.db.Unscoped().First(&deletedUser, userID)
	assert.True(s.T(), deletedUser.IsDeleted())

	// 恢复记录
	result := s.db.Model(&User{}).Unscoped().Where("id = ?", userID).Update("deleted_at", nil)
	assert.NoError(s.T(), result.Error)

	// 验证已恢复
	var restoredUser User
	err := s.db.First(&restoredUser, userID).Error
	assert.NoError(s.T(), err)
	assert.False(s.T(), restoredUser.IsDeleted())
}

// TestBaseModel_StatusManagement 测试状态管理
func (s *ModelTestSuite) TestBaseModel_StatusManagement() {
	user := &User{
		Name:  "周八",
		Email: "zhouba@example.com",
		Age:   22,
	}

	// 默认启用
	s.db.Create(user)
	assert.True(s.T(), user.IsEnabled())

	// 禁用
	user.Disable()
	s.db.Save(user)

	var updatedUser User
	s.db.First(&updatedUser, user.ID)
	assert.False(s.T(), updatedUser.IsEnabled())
	assert.Equal(s.T(), int8(0), updatedUser.Status)

	// 重新启用
	updatedUser.Enable()
	s.db.Save(&updatedUser)

	var enabledUser User
	s.db.First(&enabledUser, user.ID)
	assert.True(s.T(), enabledUser.IsEnabled())
	assert.Equal(s.T(), int8(1), enabledUser.Status)
}

// TestLightModel 测试轻量级模型
func (s *ModelTestSuite) TestLightModel() {
	product := &Product{
		Name:  "笔记本电脑",
		Price: 5999.99,
	}

	// 测试新记录
	assert.True(s.T(), product.IsNew())

	// 创建记录
	result := s.db.Create(product)
	assert.NoError(s.T(), result.Error)
	assert.NotZero(s.T(), product.ID)
	assert.False(s.T(), product.IsNew())
	assert.Equal(s.T(), product.ID, product.GetID())
	assert.True(s.T(), product.IsEnabled())

	// LightModel没有软删除，是硬删除
	productID := product.ID
	s.db.Delete(product)
	var deletedProduct Product
	err := s.db.First(&deletedProduct, productID).Error
	assert.Error(s.T(), err) // 应该找不到，因为是硬删除
}

// TestBaseModel_BatchOperations 测试批量操作
func (s *ModelTestSuite) TestBaseModel_BatchOperations() {
	users := []*User{
		{Name: "用户1", Email: "user1@example.com", Age: 20},
		{Name: "用户2", Email: "user2@example.com", Age: 21},
		{Name: "用户3", Email: "user3@example.com", Age: 22},
	}

	// 设置统一备注
	for _, user := range users {
		user.SetRemark("批量创建用户")
		user.Enable()
	}

	// 批量创建
	result := s.db.Create(&users)
	assert.NoError(s.T(), result.Error)
	assert.Equal(s.T(), int64(3), result.RowsAffected)

	// 验证所有记录
	for _, user := range users {
		assert.NotZero(s.T(), user.ID)
		assert.True(s.T(), user.IsEnabled())
		assert.Equal(s.T(), "批量创建用户", user.Remark)
	}

	// 批量更新状态
	result = s.db.Model(&User{}).Where("age >= ?", 20).Update("status", 0)
	assert.NoError(s.T(), result.Error)

	// 验证更新
	var disabledUsers []User
	s.db.Where("age >= ?", 20).Find(&disabledUsers)
	for _, user := range disabledUsers {
		assert.False(s.T(), user.IsEnabled())
	}
}

// TestBaseModel_QueryWithFilters 测试带过滤条件的查询
func (s *ModelTestSuite) TestBaseModel_QueryWithFilters() {
	// 清理所有现有数据以避免测试间干扰
	s.db.Unscoped().Where("1=1").Delete(&User{})

	// 创建测试数据
	userA := &User{Name: "测试A", Email: "testa@example.com", Age: 25}
	s.db.Create(userA)
	// 确保启用（虽然默认就是启用）
	userA.Enable()
	s.db.Save(userA)

	userB := &User{Name: "测试B", Email: "testb@example.com", Age: 30}
	s.db.Create(userB)
	// 禁用用户B
	userB.Disable()
	s.db.Save(userB)

	userC := &User{Name: "测试C", Email: "testc@example.com", Age: 35}
	s.db.Create(userC)
	// 确保启用
	userC.Enable()
	s.db.Save(userC)

	// 查询启用的用户
	var enabledUsers []User
	s.db.Where("status = ?", 1).Find(&enabledUsers)
	assert.Len(s.T(), enabledUsers, 2)
	assert.Equal(s.T(), "测试A", enabledUsers[0].Name)
	assert.Equal(s.T(), "测试C", enabledUsers[1].Name)

	// 查询年龄大于25的用户
	var olderUsers []User
	s.db.Where("age > ?", 25).Find(&olderUsers)
	assert.Len(s.T(), olderUsers, 2)

	// 组合查询（启用且年龄>=30）
	var filteredUsers []User
	s.db.Where("status = ? AND age >= ?", 1, 30).Find(&filteredUsers)
	assert.Len(s.T(), filteredUsers, 1)
	assert.Equal(s.T(), "测试C", filteredUsers[0].Name)
}

// TestBaseModel_Concurrency 测试并发更新（乐观锁）
func (s *ModelTestSuite) TestBaseModel_Concurrency() {
	// 创建用户
	user := &User{
		Name:  "并发测试",
		Email: "concurrent@example.com",
		Age:   25,
	}
	s.db.Create(user)

	// 模拟两个并发更新
	var user1, user2 User
	s.db.First(&user1, user.ID)
	s.db.First(&user2, user.ID)

	assert.Equal(s.T(), user1.Version, user2.Version)
	initialVersion := user1.Version

	// 第一个更新成功
	user1.Age = 26
	result1 := s.db.Save(&user1)
	assert.NoError(s.T(), result1.Error)
	assert.Equal(s.T(), initialVersion+1, user1.Version)

	// 第二个更新也会成功（GORM默认行为），但版本号会继续增加
	user2.Age = 27
	result2 := s.db.Save(&user2)
	assert.NoError(s.T(), result2.Error)

	// 验证最终状态
	var finalUser User
	s.db.First(&finalUser, user.ID)
	assert.Equal(s.T(), 27, finalUser.Age) // 后面的更新覆盖
	// 版本号至少增加了一次
	assert.GreaterOrEqual(s.T(), finalUser.Version, initialVersion+1)
}

// TestModelInterfaces 测试模型接口实现
func (s *ModelTestSuite) TestModelInterfaces() {
	user := &User{
		Name:  "接口测试",
		Email: "interface@example.com",
		Age:   30,
	}

	// 测试 ModelInterface
	var m ModelInterface = user
	assert.True(s.T(), m.IsNew())

	// 测试 VersionedModel
	var vm VersionedModel = user
	assert.Equal(s.T(), 0, vm.GetVersion()) // 未保存前为0

	// 测试 StatusModel
	var sm StatusModel = user
	sm.Enable()
	assert.True(s.T(), sm.IsEnabled())
	sm.Disable()
	assert.False(s.T(), sm.IsEnabled())

	// 测试 RemarkableModel
	var rm RemarkableModel = user
	rm.SetRemark("接口测试备注")

	// 创建后测试
	s.db.Create(user)
	assert.False(s.T(), m.IsNew())
	assert.Equal(s.T(), 1, vm.GetVersion())

	// 测试 SoftDeletableModel
	var sdm SoftDeletableModel = user
	assert.False(s.T(), sdm.IsDeleted())
}

// TestUUIDModel 测试UUID模型
func (s *ModelTestSuite) TestUUIDModel() {
	type Document struct {
		UUIDModel
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	// 注意：实际使用时需要自己生成UUID
	doc := &Document{
		UUIDModel: UUIDModel{ID: "550e8400-e29b-41d4-a716-446655440000"},
		Title:     "测试文档",
		Content:   "这是一个测试文档",
	}

	assert.False(s.T(), doc.IsNew())
	assert.Equal(s.T(), "550e8400-e29b-41d4-a716-446655440000", doc.GetID())
}

// TestTimestampModel 测试时间戳模型
func (s *ModelTestSuite) TestTimestampModel() {
	type Log struct {
		TimestampModel
		Message string `json:"message"`
		Level   string `json:"level"`
	}

	log := &Log{
		Message: "测试日志",
		Level:   "INFO",
	}

	// 时间戳模型不应该有时间值（创建前）
	assert.True(s.T(), log.CreatedAt.IsZero())
	assert.True(s.T(), log.UpdatedAt.IsZero())
}

// TestAuditModel 测试审计模型
func (s *ModelTestSuite) TestAuditModel() {
	// 创建审计用户
	auditUser := &AuditUser{
		Name:  "审计用户",
		Email: "audit@example.com",
	}
	auditUser.SetCreatedBy(1001)
	auditUser.SetUpdatedBy(1001)
	auditUser.Enable()
	auditUser.SetRemark("需要审计的用户")

	// 创建记录
	result := s.db.Create(auditUser)
	assert.NoError(s.T(), result.Error)
	assert.NotZero(s.T(), auditUser.ID)
	assert.Equal(s.T(), uint(1001), auditUser.GetCreatedBy())
	assert.Equal(s.T(), uint(1001), auditUser.GetUpdatedBy())

	// 更新记录
	auditUser.Name = "更新后的审计用户"
	auditUser.SetUpdatedBy(1002)
	result = s.db.Save(auditUser)
	assert.NoError(s.T(), result.Error)
	assert.Equal(s.T(), uint(1001), auditUser.GetCreatedBy()) // 创建人不变
	assert.Equal(s.T(), uint(1002), auditUser.GetUpdatedBy()) // 更新人改变

	// 验证审计接口
	var am AuditableModel = auditUser
	assert.Equal(s.T(), uint(1001), am.GetCreatedBy())
	assert.Equal(s.T(), uint(1002), am.GetUpdatedBy())
}

// TestRepositoryIntegration 测试Repository集成
func (s *ModelTestSuite) TestRepositoryIntegration() {
	ctx := context.Background()

	// 创建BaseRepository
	dbHandler := &mockDBHandler{db: s.db}
	userRepo := NewBaseRepository[User](dbHandler, "users")

	// 创建用户
	user := &User{
		Name:  "Repository测试",
		Email: "repo@example.com",
		Age:   28,
	}
	user.Enable()
	user.SetRemark("通过Repository创建")

	createdUser, err := userRepo.Create(ctx, user)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), createdUser)
	assert.NotZero(s.T(), createdUser.ID)
	assert.True(s.T(), createdUser.IsEnabled())
	assert.Equal(s.T(), "通过Repository创建", createdUser.Remark)

	// 获取用户
	fetchedUser, err := userRepo.Get(ctx, createdUser.ID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), createdUser.ID, fetchedUser.ID)
	assert.Equal(s.T(), createdUser.Name, fetchedUser.Name)

	// 更新用户
	fetchedUser.Age = 29
	fetchedUser.Disable()
	_, err = userRepo.Update(ctx, fetchedUser)
	assert.NoError(s.T(), err)

	// 验证更新
	updatedUser, err := userRepo.Get(ctx, fetchedUser.ID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 29, updatedUser.Age)
	assert.False(s.T(), updatedUser.IsEnabled())
	assert.Greater(s.T(), updatedUser.Version, 1)
}

// mockDBHandler 模拟数据库处理器
type mockDBHandler struct {
	db *gorm.DB
}

func (m *mockDBHandler) DB() *gorm.DB {
	return m.db
}

func (m *mockDBHandler) Transaction(ctx context.Context, fn func(context.Context) error) error {
	return m.db.Transaction(func(tx *gorm.DB) error {
		return fn(ctx)
	})
}

// TestSuite 运行测试套件
func TestModelTestSuite(t *testing.T) {
	suite.Run(t, new(ModelTestSuite))
}
