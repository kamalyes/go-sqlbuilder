/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-10 01:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-10 23:34:00
 * @FilePath: \go-sqlbuilder\interfaces.go
 * @Description: Database interfaces for multiple ORM compatibility (sqlx, gorm, etc.)
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package sqlbuilder

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"time"

	"github.com/jmoiron/sqlx"
	"gorm.io/gorm"
)

// ==================== 核心数据库接口 ====================

// UniversalAdapterInterface 通用适配器接口 - 统一所有框架
type UniversalAdapterInterface interface {
	DatabaseInterface

	// 适配器信息
	GetAdapterType() string
	GetAdapterName() string
	GetDialect() string

	// 批量操作
	BatchInsert(ctx context.Context, table string, data []map[string]interface{}) error
	BatchUpdate(ctx context.Context, table string, data []map[string]interface{}, whereColumns []string) error

	// 功能检测
	SupportsORM() bool
	SupportsUpsert() bool
	SupportsBulkInsert() bool
	SupportsReturning() bool

	// 获取底层实例
	GetInstance() interface{}

	// 连接统计
	GetStats() ConnectionStats
}

// ConnectionStats 连接统计信息
type ConnectionStats struct {
	OpenConnections   int
	InUse             int
	Idle              int
	WaitCount         int64
	WaitDuration      time.Duration
	MaxIdleClosed     int64
	MaxLifetimeClosed int64
}

// DatabaseInterface 数据库核心接口
type DatabaseInterface interface {
	QueryerInterface
	ExecerInterface
	TransactionInterface
	PreparedInterface
	PingInterface
	CloseInterface
}

// QueryerInterface 查询接口
type QueryerInterface interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

// ExecerInterface 执行接口
type ExecerInterface interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

// TransactionInterface 事务接口
type TransactionInterface interface {
	ExecerInterface
	Begin() (TransactionInterface, error)
	BeginTx(ctx context.Context, opts *sql.TxOptions) (TransactionInterface, error)
	Commit() error
	Rollback() error
}

// PreparedInterface 预处理语句接口
type PreparedInterface interface {
	Prepare(query string) (StatementInterface, error)
	PrepareContext(ctx context.Context, query string) (StatementInterface, error)
}

// StatementInterface 语句接口
type StatementInterface interface {
	Exec(args ...interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, args ...interface{}) (sql.Result, error)
	Query(args ...interface{}) (*sql.Rows, error)
	QueryContext(ctx context.Context, args ...interface{}) (*sql.Rows, error)
	QueryRow(args ...interface{}) *sql.Row
	QueryRowContext(ctx context.Context, args ...interface{}) *sql.Row
	Close() error
}

// PingInterface Ping接口
type PingInterface interface {
	Ping() error
	PingContext(ctx context.Context) error
}

// CloseInterface 关闭接口
type CloseInterface interface {
	Close() error
}

// ==================== SQLX 特定接口 ====================

// SqlxInterface SQLX特定接口
type SqlxInterface interface {
	DatabaseInterface
	SqlxQueryInterface
	SqlxNamedInterface
	SqlxGetInterface
	SqlxSelectInterface
}

// SqlxQueryInterface SQLX查询扩展
type SqlxQueryInterface interface {
	Queryx(query string, args ...interface{}) (*sqlx.Rows, error)
	QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error)
	QueryRowx(query string, args ...interface{}) *sqlx.Row
	QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row
}

// SqlxNamedInterface SQLX命名参数接口
type SqlxNamedInterface interface {
	NamedExec(query string, arg interface{}) (sql.Result, error)
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
	NamedQuery(query string, arg interface{}) (*sqlx.Rows, error)
	NamedQueryContext(ctx context.Context, query string, arg interface{}) (*sqlx.Rows, error)
}

// SqlxGetInterface SQLX获取单个结果接口
type SqlxGetInterface interface {
	Get(dest interface{}, query string, args ...interface{}) error
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

// SqlxSelectInterface SQLX选择多个结果接口
type SqlxSelectInterface interface {
	Select(dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

// ==================== GORM 特定接口 ====================

// GormInterface GORM特定接口
type GormInterface interface {
	GormQueryInterface
	GormTransactionInterface
	GormMigrationInterface
	GormModelInterface
	GormAssociationInterface
}

// GormQueryInterface GORM查询接口
type GormQueryInterface interface {
	Model(value interface{}) GormInterface
	Table(name string, args ...interface{}) GormInterface
	Distinct(args ...interface{}) GormInterface
	Select(query interface{}, args ...interface{}) GormInterface
	Omit(columns ...string) GormInterface
	Where(query interface{}, args ...interface{}) GormInterface
	Not(query interface{}, args ...interface{}) GormInterface
	Or(query interface{}, args ...interface{}) GormInterface
	Joins(query string, args ...interface{}) GormInterface
	Group(name string) GormInterface
	Having(query interface{}, args ...interface{}) GormInterface
	Order(value interface{}) GormInterface
	Limit(limit int) GormInterface
	Offset(offset int) GormInterface
	Scopes(funcs ...func(GormInterface) GormInterface) GormInterface
}

// GormTransactionInterface GORM事务接口
type GormTransactionInterface interface {
	Transaction(fc func(tx GormInterface) error, opts ...*sql.TxOptions) error
	Begin(opts ...*sql.TxOptions) GormInterface
	Commit() GormInterface
	Rollback() GormInterface
	SavePoint(name string) GormInterface
	RollbackTo(name string) GormInterface
}

// GormMigrationInterface GORM迁移接口
type GormMigrationInterface interface {
	AutoMigrate(dst ...interface{}) error
	Migrator() gorm.Migrator
}

// GormModelInterface GORM模型操作接口
type GormModelInterface interface {
	Create(value interface{}) GormInterface
	CreateInBatches(value interface{}, batchSize int) GormInterface
	Save(value interface{}) GormInterface
	First(dest interface{}, conds ...interface{}) GormInterface
	Take(dest interface{}, conds ...interface{}) GormInterface
	Last(dest interface{}, conds ...interface{}) GormInterface
	Find(dest interface{}, conds ...interface{}) GormInterface
	FindInBatches(dest interface{}, batchSize int, fc func(tx GormInterface, batch int) error) error
	FirstOrInit(dest interface{}, conds ...interface{}) GormInterface
	FirstOrCreate(dest interface{}, conds ...interface{}) GormInterface
	Update(column string, value interface{}) GormInterface
	Updates(values interface{}) GormInterface
	UpdateColumn(column string, value interface{}) GormInterface
	UpdateColumns(values interface{}) GormInterface
	Delete(value interface{}, conds ...interface{}) GormInterface
	Count(count *int64) GormInterface
	Row() *sql.Row
	Rows() (*sql.Rows, error)
	Scan(dest interface{}) GormInterface
	Pluck(column string, dest interface{}) GormInterface
	ScanRows(rows *sql.Rows, dest interface{}) error
}

// GormAssociationInterface GORM关联接口
type GormAssociationInterface interface {
	Association(column string) *gorm.Association
}

// ==================== 通用查询构建器接口 ====================

// QueryBuilderInterface 查询构建器核心接口
type QueryBuilderInterface interface {
	// 基本查询方法
	Select(fields ...interface{}) QueryBuilderInterface
	Table(table interface{}) QueryBuilderInterface
	Where(args ...interface{}) QueryBuilderInterface
	OrWhere(args ...interface{}) QueryBuilderInterface
	WhereIn(field string, values []interface{}) QueryBuilderInterface
	WhereNotIn(field string, values []interface{}) QueryBuilderInterface
	WhereBetween(field string, start, end interface{}) QueryBuilderInterface
	WhereNull(field string) QueryBuilderInterface
	WhereNotNull(field string) QueryBuilderInterface
	WhereLike(field, pattern string) QueryBuilderInterface
	WhereExists(subQuery QueryBuilderInterface) QueryBuilderInterface

	// JOIN操作
	Join(table, condition string) QueryBuilderInterface
	LeftJoin(table, condition string) QueryBuilderInterface
	RightJoin(table, condition string) QueryBuilderInterface
	InnerJoin(table, condition string) QueryBuilderInterface

	// 分组和排序
	GroupBy(fields ...string) QueryBuilderInterface
	Having(args ...interface{}) QueryBuilderInterface
	OrderBy(field string, direction ...string) QueryBuilderInterface
	OrderByDesc(field string) QueryBuilderInterface

	// 分页和限制
	Limit(limit int64) QueryBuilderInterface
	Offset(offset int64) QueryBuilderInterface
	Page(page, pageSize int64) QueryBuilderInterface

	// 插入、更新、删除
	Insert(data interface{}) QueryBuilderInterface
	Update(data interface{}) QueryBuilderInterface
	Delete() QueryBuilderInterface

	// 执行方法
	ToSQL() (string, []interface{})
	Exec() (sql.Result, error)
	Get(dest interface{}) error
	Find(dest interface{}) error
	Count(column ...string) (int64, error)
	Exists() (bool, error)

	// 工具方法
	Clone() QueryBuilderInterface
	Debug(enable ...bool) QueryBuilderInterface
}

// ==================== 结果处理接口 ====================

// ResultInterface 结果处理接口
type ResultInterface interface {
	// 单条记录
	First(dest interface{}) error
	FirstOrFail(dest interface{}) error
	FirstOrDefault(dest interface{}, defaultValue interface{}) error

	// 多条记录
	Get(dest interface{}) error
	All(dest interface{}) error
	Chunk(size int, callback func(records interface{}) error) error

	// 映射结果
	ToMap() (map[string]interface{}, error)
	ToMaps() ([]map[string]interface{}, error)
	ToStruct(dest interface{}) error
	ToStructs(dest interface{}) error

	// 聚合函数
	Count(column ...string) (int64, error)
	Sum(column string) (interface{}, error)
	Avg(column string) (interface{}, error)
	Max(column string) (interface{}, error)
	Min(column string) (interface{}, error)

	// 存在性检查
	Exists() (bool, error)
	DoesntExist() (bool, error)

	// 单列值
	Pluck(column string, dest interface{}) error
	PluckMap(keyColumn, valueColumn string) (map[interface{}]interface{}, error)
	Value(column string) (interface{}, error)
}

// ==================== 高级查询接口 ====================

// AdvancedQueryInterface 高级查询接口
type AdvancedQueryInterface interface {
	// 子查询
	SubQuery(alias string) QueryBuilderInterface
	WhereSubQuery(field, operator string, subQuery QueryBuilderInterface) QueryBuilderInterface

	// 联合查询
	Union(query QueryBuilderInterface) QueryBuilderInterface
	UnionAll(query QueryBuilderInterface) QueryBuilderInterface

	// 条件构建
	When(condition bool, callback func(QueryBuilderInterface) QueryBuilderInterface) QueryBuilderInterface
	Unless(condition bool, callback func(QueryBuilderInterface) QueryBuilderInterface) QueryBuilderInterface

	// 作用域
	Scope(callback func(QueryBuilderInterface) QueryBuilderInterface) QueryBuilderInterface

	// 原始SQL
	Raw(sql string, bindings ...interface{}) QueryBuilderInterface
	SelectRaw(expression string, bindings ...interface{}) QueryBuilderInterface
	WhereRaw(sql string, bindings ...interface{}) QueryBuilderInterface
	OrWhereRaw(sql string, bindings ...interface{}) QueryBuilderInterface
	HavingRaw(sql string, bindings ...interface{}) QueryBuilderInterface
	OrderByRaw(sql string) QueryBuilderInterface

	// 高级WHERE
	WhereColumn(first, operator, second string) QueryBuilderInterface
	OrWhereColumn(first, operator, second string) QueryBuilderInterface
	WhereDate(column string, operator string, value interface{}) QueryBuilderInterface
	WhereTime(column string, operator string, value interface{}) QueryBuilderInterface
	WhereYear(column string, operator string, value interface{}) QueryBuilderInterface
	WhereMonth(column string, operator string, value interface{}) QueryBuilderInterface
	WhereDay(column string, operator string, value interface{}) QueryBuilderInterface

	// JSON查询（MySQL/PostgreSQL）
	WhereJson(column, path string, value interface{}) QueryBuilderInterface
	WhereJsonContains(column string, value interface{}) QueryBuilderInterface
	WhereJsonLength(column string, operator string, value int) QueryBuilderInterface
}

// ==================== 批量操作接口 ====================

// BatchOperationInterface 批量操作接口
type BatchOperationInterface interface {
	// 批量插入
	InsertBatch(data []map[string]interface{}) (sql.Result, error)
	InsertBatchSize(data []map[string]interface{}, batchSize int) error

	// 批量更新
	UpdateBatch(data []map[string]interface{}, primaryKey string) error
	UpdateBatchSize(data []map[string]interface{}, primaryKey string, batchSize int) error

	// 批量删除
	DeleteBatch(conditions []map[string]interface{}) error
	DeleteBatchSize(conditions []map[string]interface{}, batchSize int) error

	// Upsert操作
	Upsert(data interface{}, conflictFields []string, updateFields []string) error
	UpsertBatch(data []map[string]interface{}, conflictFields []string, updateFields []string) error

	// 批量查询
	FindBatch(queries []QueryBuilderInterface) ([]interface{}, error)
	FindBatchMaps(queries []QueryBuilderInterface) ([][]map[string]interface{}, error)
}

// ==================== 缓存接口 ====================

// CacheInterface 缓存接口
type CacheInterface interface {
	// 查询缓存
	Remember(key string, ttl int, callback func() (interface{}, error)) (interface{}, error)
	RememberForever(key string, callback func() (interface{}, error)) (interface{}, error)

	// 缓存管理
	Forget(key string) error
	Flush() error

	// 缓存配置
	CacheKey(key string) QueryBuilderInterface
	CacheTTL(ttl int) QueryBuilderInterface
	NoCache() QueryBuilderInterface
}

// ==================== 事件接口 ====================

// EventInterface 事件接口
type EventInterface interface {
	// 查询事件
	BeforeQuery(callback func(sql string, bindings []interface{})) QueryBuilderInterface
	AfterQuery(callback func(sql string, bindings []interface{}, result interface{})) QueryBuilderInterface

	// 执行事件
	BeforeExec(callback func(sql string, bindings []interface{})) QueryBuilderInterface
	AfterExec(callback func(sql string, bindings []interface{}, result sql.Result)) QueryBuilderInterface

	// 事务事件
	BeforeTransaction(callback func()) QueryBuilderInterface
	AfterTransaction(callback func(committed bool)) QueryBuilderInterface
}

// ==================== 模型接口 ====================

// ModelInterface 模型接口
type ModelInterface interface {
	// 模型绑定
	Model(model interface{}) QueryBuilderInterface
	With(relations ...string) QueryBuilderInterface
	Without(relations ...string) QueryBuilderInterface

	// 软删除
	WithTrashed() QueryBuilderInterface
	OnlyTrashed() QueryBuilderInterface
	WithoutTrashed() QueryBuilderInterface
	Restore() error
	ForceDelete() error

	// 时间戳
	Touch() error
	TouchQuietly() error

	// 属性访问
	GetAttribute(key string) interface{}
	SetAttribute(key string, value interface{}) ModelInterface
	GetAttributes() map[string]interface{}
	SetAttributes(attributes map[string]interface{}) ModelInterface

	// 模型状态
	IsDirty(attributes ...string) bool
	IsClean(attributes ...string) bool
	GetDirty() map[string]interface{}
	GetOriginal(key ...string) interface{}

	// 序列化
	ToJSON() ([]byte, error)
	ToMap() map[string]interface{}
	ToArray() []interface{}
}

// ==================== 关联接口 ====================

// RelationInterface 关联接口
type RelationInterface interface {
	// 一对一
	HasOne(related interface{}, foreignKey ...string) RelationInterface
	BelongsTo(related interface{}, foreignKey ...string) RelationInterface

	// 一对多
	HasMany(related interface{}, foreignKey ...string) RelationInterface

	// 多对多
	BelongsToMany(related interface{}, pivot ...string) RelationInterface

	// 多态关联
	MorphTo(name ...string) RelationInterface
	MorphOne(related interface{}, name string) RelationInterface
	MorphMany(related interface{}, name string) RelationInterface
	MorphToMany(related interface{}, name string) RelationInterface
	MorphByMany(related interface{}, name string) RelationInterface

	// 关联操作
	Associate(model interface{}) error
	Dissociate() error
	Attach(ids interface{}, attributes ...map[string]interface{}) error
	Detach(ids ...interface{}) error
	Sync(ids interface{}) error
	SyncWithoutDetaching(ids interface{}) error
}

// ==================== 验证接口 ====================

// ValidationInterface 验证接口
type ValidationInterface interface {
	// 验证规则
	Validate(rules map[string]string) error
	ValidateOrFail(rules map[string]string) error

	// 自定义验证
	ValidateWith(validator func(data interface{}) error) error

	// 验证消息
	SetValidationMessages(messages map[string]string) QueryBuilderInterface
	GetValidationErrors() []string

	// 字段验证
	Required(fields ...string) QueryBuilderInterface
	Unique(field string, ignoreId ...interface{}) QueryBuilderInterface
	Exists(field string, table string) QueryBuilderInterface
}

// ==================== 连接池接口 ====================

// ConnectionPoolInterface 连接池接口
type ConnectionPoolInterface interface {
	// 连接管理
	GetConnection() (DatabaseInterface, error)
	ReleaseConnection(conn DatabaseInterface) error

	// 连接池配置
	SetMaxOpenConns(max int) ConnectionPoolInterface
	SetMaxIdleConns(max int) ConnectionPoolInterface
	SetConnMaxLifetime(duration int) ConnectionPoolInterface
	SetConnMaxIdleTime(duration int) ConnectionPoolInterface

	// 连接池状态
	Stats() sql.DBStats

	// 健康检查
	HealthCheck() error

	// 关闭连接池
	Close() error
}

// ==================== 多数据库接口 ====================

// MultiDatabaseInterface 多数据库接口
type MultiDatabaseInterface interface {
	// 数据库切换
	Connection(name string) QueryBuilderInterface
	On(connection string) QueryBuilderInterface

	// 配置管理
	AddConnection(name string, config interface{}) error
	RemoveConnection(name string) error
	GetConnectionNames() []string

	// 默认连接
	SetDefaultConnection(name string) MultiDatabaseInterface
	GetDefaultConnection() string

	// 连接测试
	TestConnection(name string) error
	TestAllConnections() map[string]error
}

// ==================== 查询日志接口 ====================

// QueryLogInterface 查询日志接口
type QueryLogInterface interface {
	// 日志记录
	EnableQueryLog() QueryBuilderInterface
	DisableQueryLog() QueryBuilderInterface

	// 获取日志
	GetQueryLog() []QueryLog
	FlushQueryLog() QueryBuilderInterface

	// 日志配置
	LogQueries(enable bool) QueryBuilderInterface
	SetLogger(logger QueryLogger) QueryBuilderInterface
}

// QueryLog 查询日志结构
type QueryLog struct {
	SQL      string        `json:"sql"`
	Bindings []interface{} `json:"bindings"`
	Time     float64       `json:"time"`
	Context  string        `json:"context,omitempty"`
}

// QueryLogger 查询日志记录器接口
type QueryLogger interface {
	Log(log QueryLog)
	SetLevel(level string)
	GetLogs() []QueryLog
	Clear()
}

// ==================== 错误处理接口 ====================

// ErrorHandlerInterface 错误处理接口
type ErrorHandlerInterface interface {
	// 错误处理
	OnError(callback func(error) error) QueryBuilderInterface
	IgnoreErrors() QueryBuilderInterface
	FailOnError() QueryBuilderInterface

	// 错误重试
	Retry(maxAttempts int) QueryBuilderInterface
	RetryWhen(condition func(error) bool, maxAttempts int) QueryBuilderInterface

	// 错误转换
	TransformError(transformer func(error) error) QueryBuilderInterface

	// 错误记录
	LogErrors(enable bool) QueryBuilderInterface
}

// ==================== 性能监控接口 ====================

// PerformanceInterface 性能监控接口
type PerformanceInterface interface {
	// 查询分析
	Explain() ([]map[string]interface{}, error)
	Profile() (ProfileResult, error)

	// 性能指标
	GetMetrics() Metrics
	ResetMetrics() PerformanceInterface

	// 慢查询检测
	SlowQueryThreshold(threshold float64) PerformanceInterface
	OnSlowQuery(callback func(QueryLog)) PerformanceInterface
}

// ProfileResult 性能分析结果
type ProfileResult struct {
	QueryTime    float64                  `json:"query_time"`
	PrepareTime  float64                  `json:"prepare_time"`
	ExecuteTime  float64                  `json:"execute_time"`
	RowsAffected int64                    `json:"rows_affected"`
	RowsReturned int64                    `json:"rows_returned"`
	Explain      []map[string]interface{} `json:"explain"`
}

// Metrics 性能指标
type Metrics struct {
	TotalQueries     int64    `json:"total_queries"`
	TotalTime        float64  `json:"total_time"`
	AverageTime      float64  `json:"average_time"`
	SlowestQuery     QueryLog `json:"slowest_query"`
	FastestQuery     QueryLog `json:"fastest_query"`
	QueriesPerSecond float64  `json:"queries_per_second"`
	ErrorCount       int64    `json:"error_count"`
	CacheHitRate     float64  `json:"cache_hit_rate"`
}

// ==================== 扩展接口 ====================

// ExtensionInterface 扩展接口
type ExtensionInterface interface {
	// 扩展注册
	RegisterExtension(name string, extension Extension) error
	GetExtension(name string) (Extension, bool)
	RemoveExtension(name string) bool

	// 中间件
	UseMiddleware(middleware ...Middleware) QueryBuilderInterface

	// 宏定义
	Macro(name string, callback interface{}) error
	CallMacro(name string, args ...interface{}) (interface{}, error)
}

// Extension 扩展接口
type Extension interface {
	Name() string
	Install(builder QueryBuilderInterface) error
	Uninstall(builder QueryBuilderInterface) error
}

// Middleware 中间件接口
type Middleware interface {
	Handle(next func() (interface{}, error)) (interface{}, error)
}

// ==================== 驱动适配器接口 ====================

// DriverAdapterInterface 驱动适配器接口
type DriverAdapterInterface interface {
	// 驱动信息
	DriverName() string
	SupportsFeature(feature string) bool

	// SQL方言
	QuoteIdentifier(identifier string) string
	QuoteString(str string) string
	BuildLimit(offset, limit int64) string
	BuildUpsert(table string, data map[string]interface{}, conflictFields []string) (string, []interface{})

	// 数据类型转换
	ConvertValue(value interface{}) (driver.Value, error)
	ConvertScanValue(src interface{}) (interface{}, error)

	// 特殊操作
	LastInsertId(result sql.Result) (int64, error)
	RowsAffected(result sql.Result) (int64, error)
}
