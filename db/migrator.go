/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-30 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-30 09:01:19
 * @FilePath: \go-sqlbuilder\db\migrator.go
 * @Description: 通用数据库迁移器 - 支持自动迁移、索引创建、表注释
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package db

import (
	"fmt"
	"github.com/kamalyes/go-logger"
	"gorm.io/gorm"
	"regexp"
	"strings"
)

// IndexDefinition 索引定义
type IndexDefinition struct {
	Table   string   // 表名
	Name    string   // 索引名（可选，为空时自动生成）
	Columns string   // 列定义（如 "(col1, col2)" 或 "col1 DESC"）
	Unique  bool     // 是否唯一索引
	columns []string // 内部使用：解析后的列名列表
}

// GenerateIndexName 根据命名规范自动生成索引名
// 规范：idx_{表名}_{列名} 或 idx_{表名}_{列1}_{列2}_unique
func (idx *IndexDefinition) GenerateIndexName() string {
	if idx.Name != "" {
		return idx.Name
	}

	// 解析列名
	cols := idx.parseColumns()
	if len(cols) == 0 {
		return fmt.Sprintf("idx_%s_auto", idx.Table)
	}

	// 生成索引名
	name := fmt.Sprintf("idx_%s_%s", idx.Table, strings.Join(cols, "_"))

	// 唯一索引添加后缀
	if idx.Unique {
		name += "_unique"
	}

	return name
}

// parseColumns 解析列定义，提取纯列名
func (idx *IndexDefinition) parseColumns() []string {
	if len(idx.columns) > 0 {
		return idx.columns
	}

	// 移除括号
	cols := strings.Trim(idx.Columns, "()")

	// 按逗号分割
	parts := strings.Split(cols, ",")

	var result []string
	for _, part := range parts {
		// 清理空格
		part = strings.TrimSpace(part)

		// 移除排序关键字 (ASC, DESC)
		part = regexp.MustCompile(`(?i)\s+(ASC|DESC)$`).ReplaceAllString(part, "")

		// 移除其他修饰符
		part = strings.TrimSpace(part)

		if part != "" {
			result = append(result, part)
		}
	}

	idx.columns = result
	return result
}

// --- 索引构建器便捷方法 ---

// NewIndex 创建普通索引定义（自动生成索引名）
// 示例: NewIndex("users", "email") => idx_users_email
// 示例: NewIndex("users", "user_id", "status") => idx_users_user_id_status
func NewIndex(table string, columns ...string) IndexDefinition {
	return IndexDefinition{
		Table:   table,
		Columns: formatColumns(columns),
		Unique:  false,
	}
}

// NewUniqueIndex 创建唯一索引定义（自动生成索引名）
// 示例: NewUniqueIndex("users", "email") => idx_users_email_unique
func NewUniqueIndex(table string, columns ...string) IndexDefinition {
	return IndexDefinition{
		Table:   table,
		Columns: formatColumns(columns),
		Unique:  true,
	}
}

// NewIndexWithName 创建带自定义名称的索引定义
func NewIndexWithName(table, name, columns string, unique bool) IndexDefinition {
	return IndexDefinition{
		Table:   table,
		Name:    name,
		Columns: columns,
		Unique:  unique,
	}
}

// NewIndexDesc 创建带降序的索引定义
// 示例: NewIndexDesc("messages", "created_at") => INDEX idx_messages_created_at ON messages (created_at DESC)
func NewIndexDesc(table string, columns ...string) IndexDefinition {
	// 给每个列添加 DESC
	descCols := make([]string, len(columns))
	for i, col := range columns {
		descCols[i] = col + " DESC"
	}

	return IndexDefinition{
		Table:   table,
		Columns: formatColumns(columns), // 原始列名用于生成索引名
		Unique:  false,
		columns: columns, // 保存原始列名
	}
}

// formatColumns 格式化列名为 SQL 格式
func formatColumns(columns []string) string {
	if len(columns) == 0 {
		return ""
	}
	return "(" + strings.Join(columns, ", ") + ")"
}

// TableComment 表注释定义
type TableComment struct {
	Table   string // 表名
	Comment string // 注释内容
}

// MigratorConfig 迁移器配置
type MigratorConfig struct {
	// Models 需要迁移的模型列表
	Models []interface{}
	// Indexes 需要创建的索引列表
	Indexes []IndexDefinition
	// Comments 表注释列表
	Comments []TableComment
	// Logger 日志记录器（可选，默认使用 go-logger 创建新实例）
	Logger logger.ILogger
	// SkipIndexOnError 索引创建失败时是否跳过（默认 true）
	SkipIndexOnError bool
	// SkipCommentOnError 注释添加失败时是否跳过（默认 true）
	SkipCommentOnError bool
}

// Migrator 数据库迁移器
type Migrator struct {
	db     *gorm.DB
	config *MigratorConfig
	logger logger.ILogger
}

// NewMigrator 创建迁移器
func NewMigrator(db *gorm.DB, config *MigratorConfig) *Migrator {
	if config == nil {
		config = &MigratorConfig{
			SkipIndexOnError:   true,
			SkipCommentOnError: true,
		}
	}

	log := config.Logger
	if log == nil {
		log = logger.NewLogger(nil)
	}

	return &Migrator{
		db:     db,
		config: config,
		logger: log,
	}
}

// AutoMigrate 执行完整的自动迁移（表结构 + 索引 + 注释）
func (m *Migrator) AutoMigrate() error {
	m.logger.Info("🗄️ 开始数据库自动迁移...")

	// 1. 迁移表结构
	if err := m.MigrateModels(); err != nil {
		return err
	}

	// 2. 创建索引
	if err := m.CreateIndexes(); err != nil && !m.config.SkipIndexOnError {
		return err
	}

	// 3. 添加表注释
	if err := m.AddComments(); err != nil && !m.config.SkipCommentOnError {
		return err
	}

	m.logger.Info("✅ 数据库自动迁移完成")
	return nil
}

// MigrateModels 仅迁移模型表结构
func (m *Migrator) MigrateModels() error {
	if len(m.config.Models) == 0 {
		m.logger.Warn("⚠️ 没有需要迁移的模型")
		return nil
	}

	for _, model := range m.config.Models {
		if err := m.db.AutoMigrate(model); err != nil {
			return fmt.Errorf("迁移模型 %T 失败: %w", model, err)
		}
		m.logger.Info("✅ 模型 %T 迁移成功", model)
	}

	return nil
}

// CreateIndexes 创建所有定义的索引
func (m *Migrator) CreateIndexes() error {
	if len(m.config.Indexes) == 0 {
		return nil
	}

	m.logger.Info("📑 创建数据库索引...")

	var lastErr error
	for _, idx := range m.config.Indexes {
		indexName := idx.GenerateIndexName()
		if err := m.createIndex(idx); err != nil {
			m.logger.Warn("创建索引 %s 失败: %v", indexName, err)
			lastErr = err
			if !m.config.SkipIndexOnError {
				return err
			}
		} else {
			m.logger.Debug("✅ 索引 %s 创建成功", indexName)
		}
	}

	return lastErr
}

// createIndex 创建单个索引
func (m *Migrator) createIndex(idx IndexDefinition) error {
	// 自动生成索引名（如果未指定）
	indexName := idx.GenerateIndexName()

	// 先检查索引是否存在（兼容 MySQL/PostgreSQL/SQLite）
	if m.hasIndex(idx.Table, indexName) {
		m.logger.Debug("索引 %s 已存在，跳过创建", indexName)
		return nil
	}

	var sql string
	if idx.Unique {
		sql = fmt.Sprintf("CREATE UNIQUE INDEX %s ON %s %s", indexName, idx.Table, idx.Columns)
	} else {
		sql = fmt.Sprintf("CREATE INDEX %s ON %s %s", indexName, idx.Table, idx.Columns)
	}

	return m.db.Exec(sql).Error
}

// hasIndex 检查索引是否存在
func (m *Migrator) hasIndex(table, indexName string) bool {
	return m.db.Migrator().HasIndex(table, indexName)
}

// AddComments 添加所有表注释
func (m *Migrator) AddComments() error {
	if len(m.config.Comments) == 0 {
		return nil
	}

	m.logger.Info("💬 添加表注释...")

	dialector := m.db.Dialector.Name()
	var lastErr error

	for _, c := range m.config.Comments {
		if err := m.addComment(c, dialector); err != nil {
			m.logger.Debug("表 %s 注释添加失败: %v", c.Table, err)
			lastErr = err
			if !m.config.SkipCommentOnError {
				return err
			}
		} else {
			m.logger.Debug("✅ 表 %s 注释添加成功", c.Table)
		}
	}

	return lastErr
}

// addComment 添加单个表注释
func (m *Migrator) addComment(c TableComment, dialector string) error {
	var sql string
	switch dialector {
	case "mysql":
		sql = fmt.Sprintf("ALTER TABLE %s COMMENT = '%s'", c.Table, c.Comment)
	case "postgres":
		sql = fmt.Sprintf("COMMENT ON TABLE %s IS '%s'", c.Table, c.Comment)
	default:
		// 其他数据库不支持表注释
		return nil
	}

	return m.db.Exec(sql).Error
}

// DropTables 删除指定的表（危险操作）
func (m *Migrator) DropTables(tables ...string) error {
	m.logger.Warn("⚠️ 准备删除数据表: %v", tables)

	for _, table := range tables {
		if err := m.db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", table)).Error; err != nil {
			m.logger.Error("❌ 删除表 %s 失败: %v", table, err)
			return err
		}
		m.logger.Info("🗑️ 表 %s 已删除", table)
	}

	return nil
}

// DropTablesWithModels 根据模型删除表（危险操作）
// 自动从模型中获取表名，支持 TableName() 方法
func (m *Migrator) DropTablesWithModels(models ...interface{}) error {
	tables := m.parseModelTables(models...)
	return m.DropTables(tables...)
}

// CheckTablesExist 检查表是否存在
func (m *Migrator) CheckTablesExist(tables ...string) map[string]bool {
	result := make(map[string]bool)
	for _, table := range tables {
		result[table] = m.db.Migrator().HasTable(table)
	}
	return result
}

// CheckTablesExistWithModels 根据模型检查表是否存在
func (m *Migrator) CheckTablesExistWithModels(models ...interface{}) map[string]bool {
	tables := m.parseModelTables(models...)
	return m.CheckTablesExist(tables...)
}

// HasTable 检查单个表是否存在
func (m *Migrator) HasTable(table string) bool {
	return m.db.Migrator().HasTable(table)
}

// HasTableWithModel 根据模型检查表是否存在
func (m *Migrator) HasTableWithModel(model interface{}) bool {
	tables := m.parseModelTables(model)
	if len(tables) == 0 {
		return false
	}
	return m.HasTable(tables[0])
}

// parseModelTables 解析模型获取表名列表
func (m *Migrator) parseModelTables(models ...interface{}) []string {
	tables := make([]string, 0, len(models))
	for _, model := range models {
		stmt := &gorm.Statement{DB: m.db}
		if err := stmt.Parse(model); err != nil {
			m.logger.Warn("解析模型 %T 失败: %v", model, err)
			continue
		}
		tables = append(tables, stmt.Table)
	}
	return tables
}

// GetTableName 获取模型对应的表名
func (m *Migrator) GetTableName(model interface{}) string {
	tables := m.parseModelTables(model)
	if len(tables) == 0 {
		return ""
	}
	return tables[0]
}

// --- 便捷函数 ---

// QuickMigrate 快速迁移（仅模型，无索引和注释）
func QuickMigrate(db *gorm.DB, models ...interface{}) error {
	migrator := NewMigrator(db, &MigratorConfig{
		Models: models,
	})
	return migrator.MigrateModels()
}

// QuickAutoMigrate 快速完整迁移
func QuickAutoMigrate(db *gorm.DB, config *MigratorConfig) error {
	migrator := NewMigrator(db, config)
	return migrator.AutoMigrate()
}
