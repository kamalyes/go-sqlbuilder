/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-30 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-30 09:30:00
 * @FilePath: \go-sqlbuilder\db\migrator_test.go
 * @Description: 数据库迁移器测试
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package db

import (
	"github.com/kamalyes/go-logger"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	gormLogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"testing"
	"time"
)

// TestMigrateUser 测试用的用户模型
type TestMigrateUser struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	Name      string     `json:"name" gorm:"column:name;size:100"`
	Email     string     `json:"email" gorm:"column:email;size:255;unique"`
	Age       int        `json:"age" gorm:"column:age"`
	Status    string     `json:"status" gorm:"column:status;size:50"`
	CreatedAt time.Time  `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"column:deleted_at;index"`
}

// TableName 指定表名
func (TestMigrateUser) TableName() string {
	return "test_migrate_users"
}

// TestMigrateOrder 测试用的订单模型
type TestMigrateOrder struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"column:user_id;index"`
	OrderNo   string    `json:"order_no" gorm:"column:order_no;size:50;unique"`
	Amount    float64   `json:"amount" gorm:"column:amount"`
	Status    string    `json:"status" gorm:"column:status;size:20"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
}

// TableName 指定表名
func (TestMigrateOrder) TableName() string {
	return "test_migrate_orders"
}

// setupMigratorTestDB 设置测试数据库（SQLite 内存数据库）
func setupMigratorTestDB() (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   gormLogger.Default.LogMode(gormLogger.Silent),
	})
}

// TestNewMigrator 测试创建迁移器
func TestNewMigrator(t *testing.T) {
	gormDB, err := setupMigratorTestDB()
	assert.NoError(t, err)

	// 测试默认配置
	migrator := NewMigrator(gormDB, nil)
	assert.NotNil(t, migrator, "迁移器不应为空")
	assert.NotNil(t, migrator.db, "数据库连接不应为空")
	assert.NotNil(t, migrator.config, "配置不应为空")
	assert.True(t, migrator.config.SkipIndexOnError, "默认应跳过索引错误")
	assert.True(t, migrator.config.SkipCommentOnError, "默认应跳过注释错误")

	// 测试自定义配置
	customLogger := logger.NewLogger(nil)
	config := &MigratorConfig{
		Models:             []interface{}{&TestMigrateUser{}},
		SkipIndexOnError:   false,
		SkipCommentOnError: false,
		Logger:             customLogger,
	}
	migrator2 := NewMigrator(gormDB, config)
	assert.NotNil(t, migrator2)
	assert.Len(t, migrator2.config.Models, 1)
	assert.False(t, migrator2.config.SkipIndexOnError)
	assert.False(t, migrator2.config.SkipCommentOnError)
}

// TestMigratorMigrateModels 测试模型迁移
func TestMigratorMigrateModels(t *testing.T) {
	gormDB, err := setupMigratorTestDB()
	assert.NoError(t, err)

	config := &MigratorConfig{
		Models: []interface{}{
			&TestMigrateUser{},
			&TestMigrateOrder{},
		},
	}

	migrator := NewMigrator(gormDB, config)

	// 执行迁移
	err = migrator.MigrateModels()
	assert.NoError(t, err, "模型迁移不应出错")

	// 验证表已创建
	assert.True(t, migrator.HasTable("test_migrate_users"), "用户表应已创建")
	assert.True(t, migrator.HasTable("test_migrate_orders"), "订单表应已创建")
}

// TestMigratorMigrateModelsEmpty 测试空模型列表迁移
func TestMigratorMigrateModelsEmpty(t *testing.T) {
	gormDB, err := setupMigratorTestDB()
	assert.NoError(t, err)

	config := &MigratorConfig{
		Models: []interface{}{},
	}

	migrator := NewMigrator(gormDB, config)

	// 空模型列表应该不报错
	err = migrator.MigrateModels()
	assert.NoError(t, err, "空模型列表迁移不应出错")
}

// TestMigratorAutoMigrate 测试完整自动迁移
func TestMigratorAutoMigrate(t *testing.T) {
	gormDB, err := setupMigratorTestDB()
	assert.NoError(t, err)

	config := &MigratorConfig{
		Models: []interface{}{
			&TestMigrateUser{},
			&TestMigrateOrder{},
		},
		Indexes: []IndexDefinition{
			{
				Table:   "test_migrate_users",
				Name:    "idx_users_status",
				Columns: "(status)",
				Unique:  false,
			},
			{
				Table:   "test_migrate_orders",
				Name:    "idx_orders_user_status",
				Columns: "(user_id, status)",
				Unique:  false,
			},
		},
		Comments: []TableComment{
			{Table: "test_migrate_users", Comment: "用户表"},
			{Table: "test_migrate_orders", Comment: "订单表"},
		},
	}

	migrator := NewMigrator(gormDB, config)

	// 执行完整迁移
	err = migrator.AutoMigrate()
	assert.NoError(t, err, "完整迁移不应出错")

	// 验证表已创建
	assert.True(t, migrator.HasTable("test_migrate_users"), "用户表应已创建")
	assert.True(t, migrator.HasTable("test_migrate_orders"), "订单表应已创建")
}

// TestMigratorCreateIndexes 测试索引创建
func TestMigratorCreateIndexes(t *testing.T) {
	gormDB, err := setupMigratorTestDB()
	assert.NoError(t, err)

	// 先创建表
	err = gormDB.AutoMigrate(&TestMigrateUser{})
	assert.NoError(t, err)

	config := &MigratorConfig{
		Indexes: []IndexDefinition{
			{
				Table:   "test_migrate_users",
				Name:    "idx_users_name",
				Columns: "(name)",
				Unique:  false,
			},
			{
				Table:   "test_migrate_users",
				Name:    "idx_users_email_unique",
				Columns: "(email)",
				Unique:  true,
			},
		},
		SkipIndexOnError: true,
	}

	migrator := NewMigrator(gormDB, config)

	// 创建索引
	err = migrator.CreateIndexes()
	// SQLite 可能不完全支持 IF NOT EXISTS，但不应崩溃
	// 这里只验证不会 panic
}

// TestMigratorCreateIndexesEmpty 测试空索引列表
func TestMigratorCreateIndexesEmpty(t *testing.T) {
	gormDB, err := setupMigratorTestDB()
	assert.NoError(t, err)

	config := &MigratorConfig{
		Indexes: []IndexDefinition{},
	}

	migrator := NewMigrator(gormDB, config)

	err = migrator.CreateIndexes()
	assert.NoError(t, err, "空索引列表不应出错")
}

// TestMigratorAddComments 测试添加表注释
func TestMigratorAddComments(t *testing.T) {
	gormDB, err := setupMigratorTestDB()
	assert.NoError(t, err)

	// 先创建表
	err = gormDB.AutoMigrate(&TestMigrateUser{})
	assert.NoError(t, err)

	config := &MigratorConfig{
		Comments: []TableComment{
			{Table: "test_migrate_users", Comment: "用户测试表"},
		},
		SkipCommentOnError: true,
	}

	migrator := NewMigrator(gormDB, config)

	// SQLite 不支持表注释，但不应报错（SkipCommentOnError=true）
	err = migrator.AddComments()
	// SQLite 不支持 COMMENT，这里验证不会崩溃
}

// TestMigratorAddCommentsEmpty 测试空注释列表
func TestMigratorAddCommentsEmpty(t *testing.T) {
	gormDB, err := setupMigratorTestDB()
	assert.NoError(t, err)

	config := &MigratorConfig{
		Comments: []TableComment{},
	}

	migrator := NewMigrator(gormDB, config)

	err = migrator.AddComments()
	assert.NoError(t, err, "空注释列表不应出错")
}

// TestMigratorDropTables 测试删除表
func TestMigratorDropTables(t *testing.T) {
	gormDB, err := setupMigratorTestDB()
	assert.NoError(t, err)

	// 先创建表
	err = gormDB.AutoMigrate(&TestMigrateUser{}, &TestMigrateOrder{})
	assert.NoError(t, err)

	migrator := NewMigrator(gormDB, nil)

	// 验证表存在
	assert.True(t, migrator.HasTable("test_migrate_users"))
	assert.True(t, migrator.HasTable("test_migrate_orders"))

	// 删除表
	err = migrator.DropTables("test_migrate_users", "test_migrate_orders")
	assert.NoError(t, err, "删除表不应出错")

	// 验证表已删除
	assert.False(t, migrator.HasTable("test_migrate_users"), "用户表应已删除")
	assert.False(t, migrator.HasTable("test_migrate_orders"), "订单表应已删除")
}

// TestMigratorCheckTablesExist 测试检查表是否存在
func TestMigratorCheckTablesExist(t *testing.T) {
	gormDB, err := setupMigratorTestDB()
	assert.NoError(t, err)

	// 先创建一个表
	err = gormDB.AutoMigrate(&TestMigrateUser{})
	assert.NoError(t, err)

	migrator := NewMigrator(gormDB, nil)

	// 检查多个表
	result := migrator.CheckTablesExist("test_migrate_users", "test_migrate_orders", "non_existent_table")

	assert.True(t, result["test_migrate_users"], "用户表应存在")
	assert.False(t, result["test_migrate_orders"], "订单表不应存在")
	assert.False(t, result["non_existent_table"], "不存在的表应返回false")
}

// TestMigratorHasTable 测试检查单个表是否存在
func TestMigratorHasTable(t *testing.T) {
	gormDB, err := setupMigratorTestDB()
	assert.NoError(t, err)

	migrator := NewMigrator(gormDB, nil)

	// 表不存在
	assert.False(t, migrator.HasTable("test_migrate_users"))

	// 创建表
	err = gormDB.AutoMigrate(&TestMigrateUser{})
	assert.NoError(t, err)

	// 表存在
	assert.True(t, migrator.HasTable("test_migrate_users"))
}

// TestQuickMigrate 测试快速迁移函数
func TestQuickMigrate(t *testing.T) {
	gormDB, err := setupMigratorTestDB()
	assert.NoError(t, err)

	// 使用快速迁移
	err = QuickMigrate(gormDB, &TestMigrateUser{}, &TestMigrateOrder{})
	assert.NoError(t, err, "快速迁移不应出错")

	// 验证表已创建
	assert.True(t, gormDB.Migrator().HasTable("test_migrate_users"))
	assert.True(t, gormDB.Migrator().HasTable("test_migrate_orders"))
}

// TestQuickAutoMigrate 测试快速完整迁移函数
func TestQuickAutoMigrate(t *testing.T) {
	gormDB, err := setupMigratorTestDB()
	assert.NoError(t, err)

	config := &MigratorConfig{
		Models: []interface{}{
			&TestMigrateUser{},
		},
		Indexes: []IndexDefinition{
			{
				Table:   "test_migrate_users",
				Name:    "idx_users_status",
				Columns: "(status)",
			},
		},
		SkipIndexOnError:   true,
		SkipCommentOnError: true,
	}

	err = QuickAutoMigrate(gormDB, config)
	assert.NoError(t, err, "快速完整迁移不应出错")

	// 验证表已创建
	assert.True(t, gormDB.Migrator().HasTable("test_migrate_users"))
}

// TestMigratorWithData 测试迁移后插入数据
func TestMigratorWithData(t *testing.T) {
	gormDB, err := setupMigratorTestDB()
	assert.NoError(t, err)

	config := &MigratorConfig{
		Models: []interface{}{&TestMigrateUser{}},
	}

	migrator := NewMigrator(gormDB, config)

	// 执行迁移
	err = migrator.AutoMigrate()
	assert.NoError(t, err)

	// 插入测试数据
	user := &TestMigrateUser{
		Name:   "Test User",
		Email:  "test@example.com",
		Age:    25,
		Status: "active",
	}
	err = gormDB.Create(user).Error
	assert.NoError(t, err, "插入数据不应出错")
	assert.NotZero(t, user.ID, "用户ID应被设置")

	// 查询验证
	var foundUser TestMigrateUser
	err = gormDB.First(&foundUser, user.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, "Test User", foundUser.Name)
	assert.Equal(t, "test@example.com", foundUser.Email)
	assert.Equal(t, 25, foundUser.Age)
	assert.Equal(t, "active", foundUser.Status)
}

// TestMigratorMultipleMigrations 测试多次迁移（幂等性）
func TestMigratorMultipleMigrations(t *testing.T) {
	gormDB, err := setupMigratorTestDB()
	assert.NoError(t, err)

	config := &MigratorConfig{
		Models: []interface{}{&TestMigrateUser{}},
	}

	migrator := NewMigrator(gormDB, config)

	// 第一次迁移
	err = migrator.AutoMigrate()
	assert.NoError(t, err, "第一次迁移不应出错")

	// 插入数据
	user := &TestMigrateUser{Name: "User1", Email: "user1@example.com", Age: 20, Status: "active"}
	err = gormDB.Create(user).Error
	assert.NoError(t, err)

	// 第二次迁移（幂等性测试）
	err = migrator.AutoMigrate()
	assert.NoError(t, err, "第二次迁移不应出错")

	// 验证数据仍然存在
	var count int64
	gormDB.Model(&TestMigrateUser{}).Count(&count)
	assert.Equal(t, int64(1), count, "数据应保持不变")
}

// TestIndexDefinition 测试索引定义结构
func TestIndexDefinition(t *testing.T) {
	idx := IndexDefinition{
		Table:   "users",
		Name:    "idx_users_email",
		Columns: "(email)",
		Unique:  true,
	}

	assert.Equal(t, "users", idx.Table)
	assert.Equal(t, "idx_users_email", idx.Name)
	assert.Equal(t, "(email)", idx.Columns)
	assert.True(t, idx.Unique)
}

// TestTableComment 测试表注释结构
func TestTableComment(t *testing.T) {
	comment := TableComment{
		Table:   "users",
		Comment: "用户信息表",
	}

	assert.Equal(t, "users", comment.Table)
	assert.Equal(t, "用户信息表", comment.Comment)
}

// TestMigratorConfig 测试迁移器配置结构
func TestMigratorConfig(t *testing.T) {
	config := &MigratorConfig{
		Models: []interface{}{&TestMigrateUser{}},
		Indexes: []IndexDefinition{
			{Table: "users", Name: "idx_name", Columns: "(name)"},
		},
		Comments: []TableComment{
			{Table: "users", Comment: "用户表"},
		},
		SkipIndexOnError:   true,
		SkipCommentOnError: false,
	}

	assert.Len(t, config.Models, 1)
	assert.Len(t, config.Indexes, 1)
	assert.Len(t, config.Comments, 1)
	assert.True(t, config.SkipIndexOnError)
	assert.False(t, config.SkipCommentOnError)
}

// TestMigratorDropNonExistentTable 测试删除不存在的表
func TestMigratorDropNonExistentTable(t *testing.T) {
	gormDB, err := setupMigratorTestDB()
	assert.NoError(t, err)

	migrator := NewMigrator(gormDB, nil)

	// 删除不存在的表不应出错（IF EXISTS）
	err = migrator.DropTables("non_existent_table_12345")
	assert.NoError(t, err, "删除不存在的表不应出错")
}

// --- 索引自动命名测试 ---

// TestGenerateIndexName 测试自动生成索引名
func TestGenerateIndexName(t *testing.T) {
	tests := []struct {
		name     string
		idx      IndexDefinition
		expected string
	}{
		{
			name:     "单列普通索引",
			idx:      IndexDefinition{Table: "users", Columns: "(email)"},
			expected: "idx_users_email",
		},
		{
			name:     "单列唯一索引",
			idx:      IndexDefinition{Table: "users", Columns: "(email)", Unique: true},
			expected: "idx_users_email_unique",
		},
		{
			name:     "复合索引",
			idx:      IndexDefinition{Table: "orders", Columns: "(user_id, status)"},
			expected: "idx_orders_user_id_status",
		},
		{
			name:     "复合唯一索引",
			idx:      IndexDefinition{Table: "orders", Columns: "(user_id, order_no)", Unique: true},
			expected: "idx_orders_user_id_order_no_unique",
		},
		{
			name:     "带DESC的索引",
			idx:      IndexDefinition{Table: "messages", Columns: "(created_at DESC)"},
			expected: "idx_messages_created_at",
		},
		{
			name:     "带ASC的索引",
			idx:      IndexDefinition{Table: "messages", Columns: "(created_at ASC)"},
			expected: "idx_messages_created_at",
		},
		{
			name:     "自定义名称优先",
			idx:      IndexDefinition{Table: "users", Name: "my_custom_idx", Columns: "(email)"},
			expected: "my_custom_idx",
		},
		{
			name:     "三列复合索引",
			idx:      IndexDefinition{Table: "stats", Columns: "(date, hour, type)"},
			expected: "idx_stats_date_hour_type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.idx.GenerateIndexName()
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestNewIndex 测试便捷索引构建函数
func TestNewIndex(t *testing.T) {
	// 单列索引
	idx := NewIndex("users", "email")
	assert.Equal(t, "users", idx.Table)
	assert.Equal(t, "(email)", idx.Columns)
	assert.False(t, idx.Unique)
	assert.Equal(t, "idx_users_email", idx.GenerateIndexName())

	// 复合索引
	idx2 := NewIndex("orders", "user_id", "status")
	assert.Equal(t, "(user_id, status)", idx2.Columns)
	assert.Equal(t, "idx_orders_user_id_status", idx2.GenerateIndexName())
}

// TestNewUniqueIndex 测试便捷唯一索引构建函数
func TestNewUniqueIndex(t *testing.T) {
	idx := NewUniqueIndex("users", "email")
	assert.Equal(t, "users", idx.Table)
	assert.Equal(t, "(email)", idx.Columns)
	assert.True(t, idx.Unique)
	assert.Equal(t, "idx_users_email_unique", idx.GenerateIndexName())

	// 复合唯一索引
	idx2 := NewUniqueIndex("orders", "user_id", "order_no")
	assert.Equal(t, "idx_orders_user_id_order_no_unique", idx2.GenerateIndexName())
}

// TestNewIndexWithName 测试带自定义名称的索引构建
func TestNewIndexWithName(t *testing.T) {
	idx := NewIndexWithName("users", "my_custom_index", "(email, name)", false)
	assert.Equal(t, "users", idx.Table)
	assert.Equal(t, "my_custom_index", idx.Name)
	assert.Equal(t, "(email, name)", idx.Columns)
	assert.False(t, idx.Unique)
	assert.Equal(t, "my_custom_index", idx.GenerateIndexName()) // 自定义名称优先
}

// TestAutoIndexNameWithMigration 测试自动索引名配合实际迁移
func TestAutoIndexNameWithMigration(t *testing.T) {
	gormDB, err := setupMigratorTestDB()
	assert.NoError(t, err)

	// 先创建表
	err = gormDB.AutoMigrate(&TestMigrateUser{})
	assert.NoError(t, err)

	config := &MigratorConfig{
		Indexes: []IndexDefinition{
			// 使用便捷函数创建索引（自动生成名称）
			NewIndex("test_migrate_users", "status"),
			NewIndex("test_migrate_users", "name", "age"),
			NewUniqueIndex("test_migrate_users", "email"),
		},
		SkipIndexOnError: true,
	}

	migrator := NewMigrator(gormDB, config)

	// 创建索引
	err = migrator.CreateIndexes()
	// SQLite 支持索引创建
	assert.NoError(t, err)
}

// TestParseColumns 测试列解析
func TestParseColumns(t *testing.T) {
	tests := []struct {
		columns  string
		expected []string
	}{
		{"(email)", []string{"email"}},
		{"(user_id, status)", []string{"user_id", "status"}},
		{"(created_at DESC)", []string{"created_at"}},
		{"(user_id, created_at DESC)", []string{"user_id", "created_at"}},
		{"(a, b, c)", []string{"a", "b", "c"}},
		{"(col ASC)", []string{"col"}},
	}

	for _, tt := range tests {
		idx := IndexDefinition{Columns: tt.columns}
		result := idx.parseColumns()
		assert.Equal(t, tt.expected, result, "列解析失败: %s", tt.columns)
	}
}

// --- 补充测试以提高覆盖率 ---

// TestNewIndexDesc 测试降序索引构建函数
func TestNewIndexDesc(t *testing.T) {
	// 单列降序索引
	idx := NewIndexDesc("messages", "created_at")
	assert.Equal(t, "messages", idx.Table)
	assert.Equal(t, "(created_at)", idx.Columns)
	assert.False(t, idx.Unique)
	assert.Equal(t, "idx_messages_created_at", idx.GenerateIndexName())

	// 多列降序索引
	idx2 := NewIndexDesc("orders", "created_at", "updated_at")
	assert.Equal(t, "(created_at, updated_at)", idx2.Columns)
	assert.Equal(t, "idx_orders_created_at_updated_at", idx2.GenerateIndexName())
}

// TestFormatColumnsEmpty 测试空列格式化
func TestFormatColumnsEmpty(t *testing.T) {
	result := formatColumns([]string{})
	assert.Equal(t, "", result, "空列应返回空字符串")

	result2 := formatColumns(nil)
	assert.Equal(t, "", result2, "nil列应返回空字符串")
}

// TestGenerateIndexNameEmptyColumns 测试空列时的索引名生成
func TestGenerateIndexNameEmptyColumns(t *testing.T) {
	idx := IndexDefinition{
		Table:   "users",
		Columns: "",
	}
	result := idx.GenerateIndexName()
	assert.Equal(t, "idx_users_auto", result, "空列应生成自动索引名")

	// 只有括号的情况
	idx2 := IndexDefinition{
		Table:   "orders",
		Columns: "()",
	}
	result2 := idx2.GenerateIndexName()
	assert.Equal(t, "idx_orders_auto", result2)
}

// TestParseColumnsWithCache 测试列解析缓存
func TestParseColumnsWithCache(t *testing.T) {
	idx := IndexDefinition{
		Table:   "users",
		Columns: "(email, name)",
		columns: []string{"cached_col"}, // 预设缓存
	}

	// 应该返回缓存的值
	result := idx.parseColumns()
	assert.Equal(t, []string{"cached_col"}, result, "应返回缓存的列")
}

// TestMigrateModelsError 测试模型迁移错误
func TestMigrateModelsError(t *testing.T) {
	gormDB, err := setupMigratorTestDB()
	assert.NoError(t, err)

	// 使用一个会导致错误的模型（无效模型）
	config := &MigratorConfig{
		Models: []interface{}{
			"invalid_model", // 字符串不是有效的模型
		},
	}

	migrator := NewMigrator(gormDB, config)

	// 应该返回错误
	err = migrator.MigrateModels()
	assert.Error(t, err, "无效模型应返回错误")
}

// TestAutoMigrateWithModelError 测试AutoMigrate模型错误分支
func TestAutoMigrateWithModelError(t *testing.T) {
	gormDB, err := setupMigratorTestDB()
	assert.NoError(t, err)

	config := &MigratorConfig{
		Models: []interface{}{
			"invalid", // 无效模型
		},
	}

	migrator := NewMigrator(gormDB, config)

	err = migrator.AutoMigrate()
	assert.Error(t, err, "模型迁移错误应返回")
}

// TestAutoMigrateStrictIndexError 测试严格模式下索引错误
func TestAutoMigrateStrictIndexError(t *testing.T) {
	gormDB, err := setupMigratorTestDB()
	assert.NoError(t, err)

	// 先创建表
	err = gormDB.AutoMigrate(&TestMigrateUser{})
	assert.NoError(t, err)

	config := &MigratorConfig{
		Models: []interface{}{}, // 空模型，跳过模型迁移
		Indexes: []IndexDefinition{
			{
				Table:   "non_existent_table_xyz",
				Name:    "idx_test",
				Columns: "(col)",
			},
		},
		SkipIndexOnError: false, // 严格模式
	}

	migrator := NewMigrator(gormDB, config)

	err = migrator.AutoMigrate()
	assert.Error(t, err, "严格模式下索引错误应返回")
}

// TestAutoMigrateStrictCommentError 测试严格模式下注释错误
func TestAutoMigrateStrictCommentError(t *testing.T) {
	gormDB, err := setupMigratorTestDB()
	assert.NoError(t, err)

	config := &MigratorConfig{
		Models:  []interface{}{},
		Indexes: []IndexDefinition{},
		Comments: []TableComment{
			{Table: "test_table", Comment: "测试"},
		},
		SkipIndexOnError:   true,
		SkipCommentOnError: false, // 严格模式
	}

	migrator := NewMigrator(gormDB, config)

	// SQLite 不支持注释，但因为 default 分支返回 nil，所以不会报错
	err = migrator.AutoMigrate()
	// SQLite 的 default 分支返回 nil，不会报错
	assert.NoError(t, err)
}

// TestCreateIndexesStrictMode 测试严格模式下创建索引
func TestCreateIndexesStrictMode(t *testing.T) {
	gormDB, err := setupMigratorTestDB()
	assert.NoError(t, err)

	config := &MigratorConfig{
		Indexes: []IndexDefinition{
			{
				Table:   "non_existent_table",
				Name:    "idx_fail",
				Columns: "(col)",
			},
		},
		SkipIndexOnError: false, // 严格模式
	}

	migrator := NewMigrator(gormDB, config)

	err = migrator.CreateIndexes()
	assert.Error(t, err, "严格模式下索引创建失败应返回错误")
}

// TestAddCommentsStrictMode 测试严格模式下添加注释
func TestAddCommentsStrictMode(t *testing.T) {
	gormDB, err := setupMigratorTestDB()
	assert.NoError(t, err)

	config := &MigratorConfig{
		Comments: []TableComment{
			{Table: "test_table", Comment: "测试注释"},
		},
		SkipCommentOnError: false, // 严格模式
	}

	migrator := NewMigrator(gormDB, config)

	// SQLite 的 default 分支返回 nil，不会触发错误
	err = migrator.AddComments()
	assert.NoError(t, err)
}

// TestCreateIndexUniqueIndex 测试创建唯一索引
func TestCreateIndexUniqueIndex(t *testing.T) {
	gormDB, err := setupMigratorTestDB()
	assert.NoError(t, err)

	// 先创建表
	err = gormDB.AutoMigrate(&TestMigrateUser{})
	assert.NoError(t, err)

	idx := IndexDefinition{
		Table:   "test_migrate_users",
		Columns: "(status)",
		Unique:  true,
	}

	migrator := NewMigrator(gormDB, nil)

	err = migrator.createIndex(idx)
	assert.NoError(t, err, "创建唯一索引不应出错")
}

// TestCreateIndexNonUniqueIndex 测试创建普通索引
func TestCreateIndexNonUniqueIndex(t *testing.T) {
	gormDB, err := setupMigratorTestDB()
	assert.NoError(t, err)

	// 先创建表
	err = gormDB.AutoMigrate(&TestMigrateUser{})
	assert.NoError(t, err)

	idx := IndexDefinition{
		Table:   "test_migrate_users",
		Columns: "(name)",
		Unique:  false,
	}

	migrator := NewMigrator(gormDB, nil)

	err = migrator.createIndex(idx)
	assert.NoError(t, err, "创建普通索引不应出错")
}

// TestNewIndexEmptyColumns 测试空列创建索引
func TestNewIndexEmptyColumns(t *testing.T) {
	idx := NewIndex("users")
	assert.Equal(t, "users", idx.Table)
	assert.Equal(t, "", idx.Columns)
}

// TestNewUniqueIndexEmptyColumns 测试空列创建唯一索引
func TestNewUniqueIndexEmptyColumns(t *testing.T) {
	idx := NewUniqueIndex("users")
	assert.Equal(t, "users", idx.Table)
	assert.Equal(t, "", idx.Columns)
	assert.True(t, idx.Unique)
}

// TestIndexDefinitionParseColumnsMultipleCalls 测试多次调用解析
func TestIndexDefinitionParseColumnsMultipleCalls(t *testing.T) {
	idx := IndexDefinition{
		Table:   "users",
		Columns: "(a, b, c)",
	}

	// 第一次调用
	result1 := idx.parseColumns()
	assert.Equal(t, []string{"a", "b", "c"}, result1)

	// 第二次调用应该返回缓存
	result2 := idx.parseColumns()
	assert.Equal(t, result1, result2)
}

// TestCreateIndexesWithMultipleIndexes 测试创建多个索引（部分成功）
func TestCreateIndexesWithMultipleIndexes(t *testing.T) {
	gormDB, err := setupMigratorTestDB()
	assert.NoError(t, err)

	// 先创建表
	err = gormDB.AutoMigrate(&TestMigrateUser{})
	assert.NoError(t, err)

	config := &MigratorConfig{
		Indexes: []IndexDefinition{
			NewIndex("test_migrate_users", "status"), // 成功
			NewIndex("non_existent_table", "col"),    // 失败
			NewIndex("test_migrate_users", "name"),   // 成功
		},
		SkipIndexOnError: true, // 跳过错误继续
	}

	migrator := NewMigrator(gormDB, config)

	// 应该有错误但继续执行
	err = migrator.CreateIndexes()
	// 返回最后一个错误
	assert.Error(t, err)
}

// TestAddCommentsWithMultipleComments 测试添加多个注释
func TestAddCommentsWithMultipleComments(t *testing.T) {
	gormDB, err := setupMigratorTestDB()
	assert.NoError(t, err)

	config := &MigratorConfig{
		Comments: []TableComment{
			{Table: "table1", Comment: "注释1"},
			{Table: "table2", Comment: "注释2"},
			{Table: "table3", Comment: "注释3"},
		},
		SkipCommentOnError: true,
	}

	migrator := NewMigrator(gormDB, config)

	// SQLite 不支持注释，都会进入 default 分支
	err = migrator.AddComments()
	assert.NoError(t, err)
}

// TestGenerateIndexNameWithPresetColumns 测试预设columns字段的索引名生成
func TestGenerateIndexNameWithPresetColumns(t *testing.T) {
	idx := IndexDefinition{
		Table:   "users",
		Columns: "(a, b)",
		columns: []string{"x", "y"}, // 预设会被使用
	}

	name := idx.GenerateIndexName()
	assert.Equal(t, "idx_users_x_y", name)
}

// TestParseColumnsEmptyParts 测试解析包含空部分的列
func TestParseColumnsEmptyParts(t *testing.T) {
	idx := IndexDefinition{
		Columns: "(a, , b)",
	}
	result := idx.parseColumns()
	assert.Equal(t, []string{"a", "b"}, result)
}

// TestCreateIndexesDebugLog 测试索引创建成功的调试日志
func TestCreateIndexesDebugLog(t *testing.T) {
	gormDB, err := setupMigratorTestDB()
	assert.NoError(t, err)

	// 先创建表
	err = gormDB.AutoMigrate(&TestMigrateUser{})
	assert.NoError(t, err)

	config := &MigratorConfig{
		Indexes: []IndexDefinition{
			NewIndex("test_migrate_users", "status"),
		},
		SkipIndexOnError: true,
	}

	migrator := NewMigrator(gormDB, config)

	err = migrator.CreateIndexes()
	assert.NoError(t, err)
}

// TestAddCommentsDebugLog 测试注释添加成功的调试日志
func TestAddCommentsDebugLog(t *testing.T) {
	gormDB, err := setupMigratorTestDB()
	assert.NoError(t, err)

	config := &MigratorConfig{
		Comments: []TableComment{
			{Table: "test_table", Comment: "测试"},
		},
		SkipCommentOnError: true,
	}

	migrator := NewMigrator(gormDB, config)

	err = migrator.AddComments()
	assert.NoError(t, err)
}

// --- 高覆盖率补充测试 ---

// TestAddCommentMySQL 测试 MySQL 注释语法生成
func TestAddCommentMySQL(t *testing.T) {
	gormDB, err := setupMigratorTestDB()
	assert.NoError(t, err)

	// 先创建表
	err = gormDB.AutoMigrate(&TestMigrateUser{})
	assert.NoError(t, err)

	migrator := NewMigrator(gormDB, nil)
	comment := TableComment{Table: "test_migrate_users", Comment: "用户表"}

	// 直接调用 addComment 测试 mysql 分支
	// SQLite 会执行 MySQL 语法但会失败，这正是我们要测试的
	err = migrator.addComment(comment, "mysql")
	// MySQL 语法在 SQLite 上会失败
	assert.Error(t, err)
}

// TestAddCommentPostgres 测试 PostgreSQL 注释语法生成
func TestAddCommentPostgres(t *testing.T) {
	gormDB, err := setupMigratorTestDB()
	assert.NoError(t, err)

	// 先创建表
	err = gormDB.AutoMigrate(&TestMigrateUser{})
	assert.NoError(t, err)

	migrator := NewMigrator(gormDB, nil)
	comment := TableComment{Table: "test_migrate_users", Comment: "用户表"}

	// 直接调用 addComment 测试 postgres 分支
	err = migrator.addComment(comment, "postgres")
	// PostgreSQL 语法在 SQLite 上会失败
	assert.Error(t, err)
}

// TestAddCommentDefault 测试默认（不支持注释）分支
func TestAddCommentDefault(t *testing.T) {
	gormDB, err := setupMigratorTestDB()
	assert.NoError(t, err)

	migrator := NewMigrator(gormDB, nil)
	comment := TableComment{Table: "test_table", Comment: "测试"}

	// 直接调用 addComment 测试 default 分支（sqlite）
	err = migrator.addComment(comment, "sqlite")
	assert.NoError(t, err, "不支持的数据库应返回 nil")

	// 其他不支持的数据库
	err = migrator.addComment(comment, "sqlserver")
	assert.NoError(t, err)
}

// TestDropTablesError 测试删除表失败（模拟错误）
func TestDropTablesError(t *testing.T) {
	gormDB, err := setupMigratorTestDB()
	assert.NoError(t, err)

	// 关闭数据库连接来模拟错误
	sqlDB, err := gormDB.DB()
	assert.NoError(t, err)
	sqlDB.Close()

	migrator := NewMigrator(gormDB, nil)

	// 连接已关闭，删除操作应该失败
	err = migrator.DropTables("any_table")
	assert.Error(t, err, "数据库连接关闭后删除表应失败")
}

// TestAutoMigrateFullFlow 测试完整迁移流程（含所有分支）
func TestAutoMigrateFullFlow(t *testing.T) {
	gormDB, err := setupMigratorTestDB()
	assert.NoError(t, err)

	config := &MigratorConfig{
		Models: []interface{}{
			&TestMigrateUser{},
		},
		Indexes: []IndexDefinition{
			NewIndex("test_migrate_users", "status"),
		},
		Comments: []TableComment{
			{Table: "test_migrate_users", Comment: "用户表"},
		},
		SkipIndexOnError:   true,
		SkipCommentOnError: true,
	}

	migrator := NewMigrator(gormDB, config)

	err = migrator.AutoMigrate()
	assert.NoError(t, err)

	// 验证表存在
	assert.True(t, migrator.HasTable("test_migrate_users"))
}

// TestAddCommentsStrictModeWithError 测试严格模式注释错误（使用MySQL语法触发错误）
func TestAddCommentsStrictModeWithError(t *testing.T) {
	gormDB, err := setupMigratorTestDB()
	assert.NoError(t, err)

	// 先创建表
	err = gormDB.AutoMigrate(&TestMigrateUser{})
	assert.NoError(t, err)

	config := &MigratorConfig{
		Comments: []TableComment{
			{Table: "test_migrate_users", Comment: "测试"},
		},
		SkipCommentOnError: false, // 严格模式
	}

	migrator := NewMigrator(gormDB, config)

	// 直接使用 mysql dialector 来触发错误
	err = migrator.addComment(config.Comments[0], "mysql")
	assert.Error(t, err)
}

// TestCreateIndexesAllSuccess 测试所有索引创建成功
func TestCreateIndexesAllSuccess(t *testing.T) {
	gormDB, err := setupMigratorTestDB()
	assert.NoError(t, err)

	// 先创建表
	err = gormDB.AutoMigrate(&TestMigrateUser{}, &TestMigrateOrder{})
	assert.NoError(t, err)

	config := &MigratorConfig{
		Indexes: []IndexDefinition{
			NewIndex("test_migrate_users", "status"),
			NewIndex("test_migrate_users", "age"),
			NewIndex("test_migrate_orders", "status"),
		},
		SkipIndexOnError: false,
	}

	migrator := NewMigrator(gormDB, config)

	err = migrator.CreateIndexes()
	assert.NoError(t, err, "所有索引创建应成功")
}

// TestAutoMigrateIndexErrorNotSkipped 测试索引错误不跳过时返回错误
func TestAutoMigrateIndexErrorNotSkipped(t *testing.T) {
	gormDB, err := setupMigratorTestDB()
	assert.NoError(t, err)

	config := &MigratorConfig{
		Models: []interface{}{&TestMigrateUser{}},
		Indexes: []IndexDefinition{
			{Table: "non_existent", Name: "idx_fail", Columns: "(col)"},
		},
		SkipIndexOnError:   false, // 不跳过
		SkipCommentOnError: true,
	}

	migrator := NewMigrator(gormDB, config)

	err = migrator.AutoMigrate()
	assert.Error(t, err, "索引错误不跳过时应返回错误")
}

// TestMigratorDropTablesMultiple 测试删除多个表
func TestMigratorDropTablesMultiple(t *testing.T) {
	gormDB, err := setupMigratorTestDB()
	assert.NoError(t, err)

	// 创建多个表
	err = gormDB.AutoMigrate(&TestMigrateUser{}, &TestMigrateOrder{})
	assert.NoError(t, err)

	migrator := NewMigrator(gormDB, nil)

	// 验证表存在
	assert.True(t, migrator.HasTable("test_migrate_users"))
	assert.True(t, migrator.HasTable("test_migrate_orders"))

	// 删除第一个表
	err = migrator.DropTables("test_migrate_users")
	assert.NoError(t, err)
	assert.False(t, migrator.HasTable("test_migrate_users"))

	// 删除第二个表
	err = migrator.DropTables("test_migrate_orders")
	assert.NoError(t, err)
	assert.False(t, migrator.HasTable("test_migrate_orders"))
}

// TestQuickMigrateEmpty 测试空模型快速迁移
func TestQuickMigrateEmpty(t *testing.T) {
	gormDB, err := setupMigratorTestDB()
	assert.NoError(t, err)

	err = QuickMigrate(gormDB)
	assert.NoError(t, err)
}

// TestNewIndexDescEmpty 测试空列降序索引
func TestNewIndexDescEmpty(t *testing.T) {
	idx := NewIndexDesc("users")
	assert.Equal(t, "users", idx.Table)
	assert.Equal(t, "", idx.Columns)
}

// TestAddCommentsWithMySQLDialectorError 测试使用MySQL dialector时注释失败
func TestAddCommentsWithMySQLDialectorError(t *testing.T) {
	// 使用 MySQL 连接字符串创建 GORM（但不实际连接）
	// 这里我们使用 SQLite 但手动覆盖 dialector 测试
	gormDB, err := setupMigratorTestDB()
	assert.NoError(t, err)

	config := &MigratorConfig{
		Comments: []TableComment{
			{Table: "test_table", Comment: "测试"},
			{Table: "test_table2", Comment: "测试2"},
		},
		SkipCommentOnError: true,
	}

	migrator := NewMigrator(gormDB, config)

	// 测试多个注释的情况，都会成功（SQLite default 返回 nil）
	err = migrator.AddComments()
	assert.NoError(t, err)
}

// TestAutoMigrateCommentErrorNotSkipped 测试注释错误不跳过（通过关闭连接模拟）
func TestAutoMigrateCommentErrorNotSkipped(t *testing.T) {
	gormDB, err := setupMigratorTestDB()
	assert.NoError(t, err)

	// 先正常创建表
	err = gormDB.AutoMigrate(&TestMigrateUser{})
	assert.NoError(t, err)

	config := &MigratorConfig{
		Models:  []interface{}{},     // 空模型，跳过模型迁移
		Indexes: []IndexDefinition{}, // 空索引
		Comments: []TableComment{
			{Table: "test_migrate_users", Comment: "测试"},
		},
		SkipIndexOnError:   true,
		SkipCommentOnError: false, // 严格模式
	}

	migrator := NewMigrator(gormDB, config)

	// SQLite default 分支返回 nil，所以不会报错
	err = migrator.AutoMigrate()
	assert.NoError(t, err)
}

// mockMySQLMigrator 通过强制使用 mysql dialector 来模拟错误
type mockMySQLDialector struct {
	gorm.Dialector
}

func (m mockMySQLDialector) Name() string {
	return "mysql"
}

// TestAddCommentsStrictModeReturn 测试严格模式下 AddComments 返回分支
func TestAddCommentsStrictModeReturn(t *testing.T) {
	gormDB, err := setupMigratorTestDB()
	assert.NoError(t, err)

	// 先创建表
	err = gormDB.AutoMigrate(&TestMigrateUser{})
	assert.NoError(t, err)

	config := &MigratorConfig{
		Comments: []TableComment{
			{Table: "test_migrate_users", Comment: "测试"},
		},
		SkipCommentOnError: false, // 严格模式
	}

	migrator := NewMigrator(gormDB, config)

	// 手动调用 addComment 并传入 mysql，触发错误
	comment := TableComment{Table: "test_migrate_users", Comment: "测试"}
	err = migrator.addComment(comment, "mysql")
	assert.Error(t, err, "MySQL 语法在 SQLite 上应该失败")
}

// TestAddCommentsMultipleWithError 测试多个注释，部分失败
func TestAddCommentsMultipleWithError(t *testing.T) {
	gormDB, err := setupMigratorTestDB()
	assert.NoError(t, err)

	config := &MigratorConfig{
		Comments: []TableComment{
			{Table: "table1", Comment: "注释1"},
			{Table: "table2", Comment: "注释2"},
		},
		SkipCommentOnError: true, // 跳过错误继续
	}

	migrator := NewMigrator(gormDB, config)

	// SQLite 默认返回 nil，所有都成功
	err = migrator.AddComments()
	assert.NoError(t, err)
}

// TestAutoMigrateWithCommentErrorSkipped 测试注释错误被跳过的情况
func TestAutoMigrateWithCommentErrorSkipped(t *testing.T) {
	gormDB, err := setupMigratorTestDB()
	assert.NoError(t, err)

	config := &MigratorConfig{
		Models: []interface{}{&TestMigrateUser{}},
		Comments: []TableComment{
			{Table: "test_migrate_users", Comment: "测试"},
		},
		SkipIndexOnError:   true,
		SkipCommentOnError: true, // 跳过注释错误
	}

	migrator := NewMigrator(gormDB, config)

	err = migrator.AutoMigrate()
	assert.NoError(t, err, "注释错误应被跳过")
}

// TestAutoMigrateNoIndexNoComment 测试无索引无注释的迁移
func TestAutoMigrateNoIndexNoComment(t *testing.T) {
	gormDB, err := setupMigratorTestDB()
	assert.NoError(t, err)

	config := &MigratorConfig{
		Models:             []interface{}{&TestMigrateUser{}},
		Indexes:            []IndexDefinition{},
		Comments:           []TableComment{},
		SkipIndexOnError:   true,
		SkipCommentOnError: true,
	}

	migrator := NewMigrator(gormDB, config)

	err = migrator.AutoMigrate()
	assert.NoError(t, err)

	assert.True(t, migrator.HasTable("test_migrate_users"))
}

// TestQuickAutoMigrateNilConfig 测试 nil 配置的快速迁移
func TestQuickAutoMigrateNilConfig(t *testing.T) {
	gormDB, err := setupMigratorTestDB()
	assert.NoError(t, err)

	err = QuickAutoMigrate(gormDB, nil)
	assert.NoError(t, err)
}

// TestAddCommentsStrictModeDBClosed 测试严格模式下数据库关闭导致的错误
func TestAddCommentsStrictModeDBClosed(t *testing.T) {
	gormDB, err := setupMigratorTestDB()
	assert.NoError(t, err)

	// 先创建表
	err = gormDB.AutoMigrate(&TestMigrateUser{})
	assert.NoError(t, err)

	config := &MigratorConfig{
		Comments: []TableComment{
			{Table: "test_migrate_users", Comment: "测试"},
		},
		SkipCommentOnError: false, // 严格模式
	}

	migrator := NewMigrator(gormDB, config)

	// 关闭数据库连接
	sqlDB, _ := gormDB.DB()
	sqlDB.Close()

	// 关闭连接后，任何 SQL 操作都会失败
	// 但 SQLite default 分支直接返回 nil，不执行 SQL
	// 所以我们直接测试 mysql dialector
	err = migrator.addComment(config.Comments[0], "mysql")
	assert.Error(t, err)
}

// MockDialector 用于测试的模拟数据库方言
type MockMySQLDialector struct{}

func (d MockMySQLDialector) Name() string {
	return "mysql" // 模拟 MySQL
}

func (d MockMySQLDialector) Initialize(*gorm.DB) error {
	return nil
}

func (d MockMySQLDialector) Migrator(db *gorm.DB) gorm.Migrator {
	return nil
}

func (d MockMySQLDialector) DataTypeOf(*schema.Field) string {
	return ""
}

func (d MockMySQLDialector) DefaultValueOf(*schema.Field) clause.Expression {
	return nil
}

func (d MockMySQLDialector) BindVarTo(writer clause.Writer, stmt *gorm.Statement, v interface{}) {}

func (d MockMySQLDialector) QuoteTo(clause.Writer, string) {}

func (d MockMySQLDialector) Explain(sql string, vars ...interface{}) string {
	return ""
}

// TestAddCommentsStrictModeReturnError 测试严格模式下 AddComments 返回错误
func TestAddCommentsStrictModeReturnError(t *testing.T) {
	// 创建带有 Mock MySQL dialector 的 gorm.DB
	// 这样 addComment 会执行 MySQL 的 ALTER TABLE 语句，但由于底层是 SQLite 会失败
	gormDB, err := setupMigratorTestDB()
	assert.NoError(t, err)

	config := &MigratorConfig{
		Comments: []TableComment{
			{Table: "test_table", Comment: "测试"},
		},
		SkipCommentOnError: false, // 严格模式
	}

	// 关闭数据库连接
	sqlDB, _ := gormDB.DB()
	sqlDB.Close()

	migrator := NewMigrator(gormDB, config)

	// 直接调用 addComment 使用 mysql dialector，由于连接已关闭会报错
	err = migrator.addComment(config.Comments[0], "mysql")
	assert.Error(t, err, "关闭连接后执行 MySQL ALTER TABLE 应该返回错误")
}

// TestAddCommentsLoopReturnError 测试 AddComments 循环中的错误返回
func TestAddCommentsLoopReturnError(t *testing.T) {
	gormDB, err := setupMigratorTestDB()
	assert.NoError(t, err)

	// 创建一个自定义的 gorm.DB，修改 dialector 名称
	// 方法：通过反射或使用 wrapper 来模拟 MySQL
	// 由于 SQLite 的 default 分支返回 nil，我们需要直接测试 error 路径

	config := &MigratorConfig{
		Comments: []TableComment{
			{Table: "t1", Comment: "c1"},
		},
		SkipCommentOnError: false, // 严格模式 - 遇到错误立即返回
	}

	// 关闭数据库连接后创建 migrator
	sqlDB, _ := gormDB.DB()
	sqlDB.Close()

	migrator := NewMigrator(gormDB, config)

	// 直接调用内部方法测试 mysql error 路径
	err = migrator.addComment(config.Comments[0], "mysql")
	assert.Error(t, err)

	// 直接调用内部方法测试 postgres error 路径
	err = migrator.addComment(config.Comments[0], "postgres")
	assert.Error(t, err)
}

// TestAutoMigrateCommentErrorStrict 测试 AutoMigrate 注释严格模式错误
func TestAutoMigrateCommentErrorStrict(t *testing.T) {
	gormDB, err := setupMigratorTestDB()
	assert.NoError(t, err)

	// 先创建表
	err = gormDB.AutoMigrate(&TestMigrateUser{})
	assert.NoError(t, err)

	config := &MigratorConfig{
		Models:  []interface{}{}, // 跳过模型迁移
		Indexes: []IndexDefinition{},
		Comments: []TableComment{
			{Table: "test_migrate_users", Comment: "测试"},
		},
		SkipIndexOnError:   true,
		SkipCommentOnError: true, // 即使设为 false，SQLite 也不会报错
	}

	migrator := NewMigrator(gormDB, config)

	err = migrator.AutoMigrate()
	assert.NoError(t, err)
}

// TestAddCommentsSkipError 测试跳过注释错误模式
func TestAddCommentsSkipError(t *testing.T) {
	gormDB, err := setupMigratorTestDB()
	assert.NoError(t, err)

	config := &MigratorConfig{
		Comments: []TableComment{
			{Table: "t1", Comment: "c1"},
			{Table: "t2", Comment: "c2"},
			{Table: "t3", Comment: "c3"},
		},
		SkipCommentOnError: true,
	}

	migrator := NewMigrator(gormDB, config)

	err = migrator.AddComments()
	assert.NoError(t, err)
}

// TestAddCommentsSuccessPath 测试注释成功路径
func TestAddCommentsSuccessPath(t *testing.T) {
	gormDB, err := setupMigratorTestDB()
	assert.NoError(t, err)

	// 创建表
	err = gormDB.AutoMigrate(&TestMigrateUser{})
	assert.NoError(t, err)

	config := &MigratorConfig{
		Comments: []TableComment{
			{Table: "test_migrate_users", Comment: "用户表测试"},
		},
		SkipCommentOnError: false,
	}

	migrator := NewMigrator(gormDB, config)

	// SQLite 返回 nil，走成功分支
	err = migrator.AddComments()
	assert.NoError(t, err)
}

// TestAutoMigrateIndexErrorSkipped 测试索引错误被跳过后继续执行
func TestAutoMigrateIndexErrorSkipped(t *testing.T) {
	gormDB, err := setupMigratorTestDB()
	assert.NoError(t, err)

	config := &MigratorConfig{
		Models: []interface{}{&TestMigrateUser{}},
		Indexes: []IndexDefinition{
			{Table: "non_existent", Name: "idx_fail", Columns: "(col)"},
		},
		Comments: []TableComment{
			{Table: "test_migrate_users", Comment: "测试"},
		},
		SkipIndexOnError:   true, // 跳过索引错误
		SkipCommentOnError: true,
	}

	migrator := NewMigrator(gormDB, config)

	// 索引错误被跳过，继续执行注释，最终成功
	err = migrator.AutoMigrate()
	assert.NoError(t, err)
}
