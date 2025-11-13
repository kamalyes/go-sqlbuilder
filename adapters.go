/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-10 01:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-11 09:20:54
 * @FilePath: \go-sqlbuilder\adapters.go
 * @Description: Database adapters for sqlx, gorm and other ORM frameworks
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package sqlbuilder

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"gorm.io/gorm"

	"github.com/kamalyes/go-sqlbuilder/errors"
)

// ==================== SQLX 适配器 ====================

// SqlxAdapter SQLX数据库适配器 - 实现通用适配器接口
type SqlxAdapter struct {
	db   *sqlx.DB
	tx   *sqlx.Tx // 事务状态
	name string
}

// NewSqlxAdapter 创建SQLX适配器
func NewSqlxAdapter(db *sqlx.DB) *SqlxAdapter {
	return &SqlxAdapter{
		db:   db,
		name: "SQLX-Adapter",
	}
}

// NewSqlxTxAdapter 创建SQLX事务适配器
func NewSqlxTxAdapter(tx *sqlx.Tx) *SqlxAdapter {
	return &SqlxAdapter{
		tx:   tx,
		name: "SQLX-Transaction",
	}
}

// ==================== 通用适配器接口实现 ====================

// GetAdapterType 获取适配器类型
func (a *SqlxAdapter) GetAdapterType() string {
	return "SQLX"
}

// GetAdapterName 获取适配器名称
func (a *SqlxAdapter) GetAdapterName() string {
	return a.name
}

// GetDialect 获取数据库方言
func (a *SqlxAdapter) GetDialect() string {
	if a.db != nil {
		return a.db.DriverName()
	}
	return "unknown"
}

// SupportsORM ORM支持检测
func (a *SqlxAdapter) SupportsORM() bool {
	return false // SQLX是查询构建器，不是ORM
}

// SupportsUpsert Upsert支持检测
func (a *SqlxAdapter) SupportsUpsert() bool {
	dialect := a.GetDialect()
	// 根据数据库类型判断是否支持upsert
	switch dialect {
	case "mysql", "postgres", "sqlite3":
		return true
	default:
		return false
	}
}

// SupportsBulkInsert 批量插入支持检测
func (a *SqlxAdapter) SupportsBulkInsert() bool {
	return true
}

// SupportsReturning RETURNING语句支持检测
func (a *SqlxAdapter) SupportsReturning() bool {
	dialect := a.GetDialect()
	return dialect == "postgres" || dialect == "sqlite3"
}

// GetInstance 获取底层SQLX实例
func (a *SqlxAdapter) GetInstance() interface{} {
	if a.tx != nil {
		return a.tx
	}
	return a.db
}

// GetStats 获取连接统计
func (a *SqlxAdapter) GetStats() ConnectionStats {
	if a.db != nil {
		stats := a.db.Stats()
		return ConnectionStats{
			OpenConnections:   stats.OpenConnections,
			InUse:             stats.InUse,
			Idle:              stats.Idle,
			WaitCount:         stats.WaitCount,
			WaitDuration:      stats.WaitDuration,
			MaxIdleClosed:     stats.MaxIdleClosed,
			MaxLifetimeClosed: stats.MaxLifetimeClosed,
		}
	}
	return ConnectionStats{}
}

// ==================== 批量操作实现 ====================

// BatchInsert 批量插入
func (a *SqlxAdapter) BatchInsert(ctx context.Context, table string, data []map[string]interface{}) error {
	if len(data) == 0 {
		return nil
	}

	// 获取列名
	var columns []string
	for col := range data[0] {
		columns = append(columns, col)
	}

	// 构建批量插入SQL (优化: 使用 strings.Builder)
	var queryBuf strings.Builder
	queryBuf.WriteString("INSERT INTO ")
	queryBuf.WriteString(table)
	queryBuf.WriteString(" (")
	for i, col := range columns {
		if i > 0 {
			queryBuf.WriteString(", ")
		}
		queryBuf.WriteString(col)
	}
	queryBuf.WriteString(") VALUES ")

	var values []interface{}
	var placeholders []string

	for _, row := range data {
		var rowPlaceholders []string
		for _, col := range columns {
			values = append(values, row[col])
			rowPlaceholders = append(rowPlaceholders, "?")
		}
		placeholders = append(placeholders, "("+strings.Join(rowPlaceholders, ", ")+")")
	}

	for i, p := range placeholders {
		if i > 0 {
			queryBuf.WriteString(", ")
		}
		queryBuf.WriteString(p)
	}
	
	query := queryBuf.String()

	// 执行批量插入
	if a.tx != nil {
		_, err := a.tx.ExecContext(ctx, query, values...)
		return err
	} else if a.db != nil {
		_, err := a.db.ExecContext(ctx, query, values...)
		return err
	}

	return errors.NewError(errors.ErrorCodeNoDatabaseConn, "no database connection available")
}

// BatchUpdate 批量更新
func (a *SqlxAdapter) BatchUpdate(ctx context.Context, table string, data []map[string]interface{}, whereColumns []string) error {
	if len(data) == 0 {
		return nil
	}

	// 使用事务进行批量更新
	var tx *sqlx.Tx
	var err error
	var shouldCommit bool

	if a.tx != nil {
		tx = a.tx
	} else if a.db != nil {
		tx, err = a.db.BeginTxx(ctx, nil)
		if err != nil {
			return errors.NewErrorf(errors.ErrorCodeDBError, errors.MsgDatabaseOperationFailed+": %v", err)
		}
		shouldCommit = true
		defer func() {
			if err != nil && shouldCommit {
				tx.Rollback()
			}
		}()
	} else {
		return errors.NewError(errors.ErrorCodeNoDatabaseConn, "no database connection available")
	}

	for _, row := range data {
		var setClauses []string
		var setValues []interface{}
		var whereClauses []string
		var whereValues []interface{}

		for col, val := range row {
			isWhereColumn := false
			for _, whereCol := range whereColumns {
				if col == whereCol {
					isWhereColumn = true
					break
				}
			}

			if isWhereColumn {
				whereClauses = append(whereClauses, col+" = ?")
				whereValues = append(whereValues, val)
			} else {
				setClauses = append(setClauses, col+" = ?")
				setValues = append(setValues, val)
			}
		}

		if len(setClauses) == 0 || len(whereClauses) == 0 {
			continue
		}

		query := fmt.Sprintf("UPDATE %s SET %s WHERE %s",
			table,
			strings.Join(setClauses, ", "),
			strings.Join(whereClauses, " AND "))

		args := append(setValues, whereValues...)
		_, err = tx.ExecContext(ctx, query, args...)
		if err != nil {
			return errors.Wrap(err, errors.ErrorCodeDBFailedUpdate)
		}
	}

	if shouldCommit {
		return tx.Commit()
	}

	return nil
}

// 实现DatabaseInterface接口
func (a *SqlxAdapter) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return a.db.Query(query, args...)
}

func (a *SqlxAdapter) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return a.db.QueryContext(ctx, query, args...)
}

func (a *SqlxAdapter) QueryRow(query string, args ...interface{}) *sql.Row {
	return a.db.QueryRow(query, args...)
}

func (a *SqlxAdapter) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return a.db.QueryRowContext(ctx, query, args...)
}

func (a *SqlxAdapter) Exec(query string, args ...interface{}) (sql.Result, error) {
	return a.db.Exec(query, args...)
}

func (a *SqlxAdapter) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return a.db.ExecContext(ctx, query, args...)
}

func (a *SqlxAdapter) Begin() (TransactionInterface, error) {
	if a.db == nil {
		return nil, errors.NewError(errors.ErrorCodeCacheStoreNotFound, errors.MsgNoDatabaseConnection)
	}
	tx, err := a.db.Beginx()
	if err != nil {
		return nil, err
	}
	return NewSqlxTxAdapter(tx), nil
}

func (a *SqlxAdapter) BeginTx(ctx context.Context, opts *sql.TxOptions) (TransactionInterface, error) {
	if a.db == nil {
		return nil, errors.NewError(errors.ErrorCodeCacheStoreNotFound, errors.MsgNoDatabaseConnection)
	}
	tx, err := a.db.BeginTxx(ctx, opts)
	if err != nil {
		return nil, err
	}
	return NewSqlxTxAdapter(tx), nil
}

func (a *SqlxAdapter) Prepare(query string) (StatementInterface, error) {
	stmt, err := a.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	return NewSqlxStmtAdapter(stmt), nil
}

func (a *SqlxAdapter) PrepareContext(ctx context.Context, query string) (StatementInterface, error) {
	stmt, err := a.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	return NewSqlxStmtAdapter(stmt), nil
}

func (a *SqlxAdapter) Ping() error {
	return a.db.Ping()
}

func (a *SqlxAdapter) PingContext(ctx context.Context) error {
	return a.db.PingContext(ctx)
}

func (a *SqlxAdapter) Commit() error {
	return errors.NewError(errors.ErrorCodeBuilderNotInitialized, "not in a transaction")
}

func (a *SqlxAdapter) Rollback() error {
	return errors.NewError(errors.ErrorCodeBuilderNotInitialized, "not in a transaction")
}

func (a *SqlxAdapter) Close() error {
	return a.db.Close()
}

// 实现SqlxInterface特有方法
func (a *SqlxAdapter) Queryx(query string, args ...interface{}) (*sqlx.Rows, error) {
	return a.db.Queryx(query, args...)
}

func (a *SqlxAdapter) QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error) {
	return a.db.QueryxContext(ctx, query, args...)
}

func (a *SqlxAdapter) QueryRowx(query string, args ...interface{}) *sqlx.Row {
	return a.db.QueryRowx(query, args...)
}

func (a *SqlxAdapter) QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row {
	return a.db.QueryRowxContext(ctx, query, args...)
}

func (a *SqlxAdapter) NamedExec(query string, arg interface{}) (sql.Result, error) {
	return a.db.NamedExec(query, arg)
}

func (a *SqlxAdapter) NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error) {
	return a.db.NamedExecContext(ctx, query, arg)
}

func (a *SqlxAdapter) NamedQuery(query string, arg interface{}) (*sqlx.Rows, error) {
	return a.db.NamedQuery(query, arg)
}

func (a *SqlxAdapter) NamedQueryContext(ctx context.Context, query string, arg interface{}) (*sqlx.Rows, error) {
	return a.db.NamedQueryContext(ctx, query, arg)
}

func (a *SqlxAdapter) Get(dest interface{}, query string, args ...interface{}) error {
	return a.db.Get(dest, query, args...)
}

func (a *SqlxAdapter) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return a.db.GetContext(ctx, dest, query, args...)
}

func (a *SqlxAdapter) Select(dest interface{}, query string, args ...interface{}) error {
	return a.db.Select(dest, query, args...)
}

func (a *SqlxAdapter) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return a.db.SelectContext(ctx, dest, query, args...)
}

func (a *SqlxAdapter) GetDB() *sqlx.DB {
	return a.db
}

// ==================== SQLX 事务适配器 ====================

// SqlxTxAdapter SQLX事务适配器 (保留原有设计兼容性)
type SqlxTxAdapter struct {
	tx *sqlx.Tx
}

// NewSqlxTxAdapterLegacy 创建SQLX事务适配器 (兼容性)
func NewSqlxTxAdapterLegacy(tx *sqlx.Tx) *SqlxTxAdapter {
	return &SqlxTxAdapter{tx: tx}
}

func (a *SqlxTxAdapter) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return a.tx.Query(query, args...)
}

func (a *SqlxTxAdapter) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return a.tx.QueryContext(ctx, query, args...)
}

func (a *SqlxTxAdapter) QueryRow(query string, args ...interface{}) *sql.Row {
	return a.tx.QueryRow(query, args...)
}

func (a *SqlxTxAdapter) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return a.tx.QueryRowContext(ctx, query, args...)
}

func (a *SqlxTxAdapter) Exec(query string, args ...interface{}) (sql.Result, error) {
	return a.tx.Exec(query, args...)
}

func (a *SqlxTxAdapter) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return a.tx.ExecContext(ctx, query, args...)
}

func (a *SqlxTxAdapter) Commit() error {
	return a.tx.Commit()
}

func (a *SqlxTxAdapter) Rollback() error {
	return a.tx.Rollback()
}

func (a *SqlxTxAdapter) Begin() (TransactionInterface, error) {
	return nil, errors.NewError(errors.ErrorCodeNestedTransaction, "cannot begin transaction within transaction")
}

func (a *SqlxTxAdapter) BeginTx(ctx context.Context, opts *sql.TxOptions) (TransactionInterface, error) {
	return nil, errors.NewError(errors.ErrorCodeNestedTransaction, "cannot begin transaction within transaction")
}

func (a *SqlxTxAdapter) Get(dest interface{}, query string, args ...interface{}) error {
	return a.tx.Get(dest, query, args...)
}

func (a *SqlxTxAdapter) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return a.tx.GetContext(ctx, dest, query, args...)
}

func (a *SqlxTxAdapter) Select(dest interface{}, query string, args ...interface{}) error {
	return a.tx.Select(dest, query, args...)
}

func (a *SqlxTxAdapter) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return a.tx.SelectContext(ctx, dest, query, args...)
}

func (a *SqlxTxAdapter) GetTx() *sqlx.Tx {
	return a.tx
}

// ==================== SQLX 语句适配器 ====================

// SqlxStmtAdapter SQLX语句适配器
type SqlxStmtAdapter struct {
	stmt *sql.Stmt
}

// NewSqlxStmtAdapter 创建SQLX语句适配器
func NewSqlxStmtAdapter(stmt *sql.Stmt) *SqlxStmtAdapter {
	return &SqlxStmtAdapter{stmt: stmt}
}

func (a *SqlxStmtAdapter) Exec(args ...interface{}) (sql.Result, error) {
	return a.stmt.Exec(args...)
}

func (a *SqlxStmtAdapter) ExecContext(ctx context.Context, args ...interface{}) (sql.Result, error) {
	return a.stmt.ExecContext(ctx, args...)
}

func (a *SqlxStmtAdapter) Query(args ...interface{}) (*sql.Rows, error) {
	return a.stmt.Query(args...)
}

func (a *SqlxStmtAdapter) QueryContext(ctx context.Context, args ...interface{}) (*sql.Rows, error) {
	return a.stmt.QueryContext(ctx, args...)
}

func (a *SqlxStmtAdapter) QueryRow(args ...interface{}) *sql.Row {
	return a.stmt.QueryRow(args...)
}

func (a *SqlxStmtAdapter) QueryRowContext(ctx context.Context, args ...interface{}) *sql.Row {
	return a.stmt.QueryRowContext(ctx, args...)
}

func (a *SqlxStmtAdapter) Close() error {
	return a.stmt.Close()
}

// ==================== GORM 适配器 ====================

// GormAdapter GORM数据库适配器 - 实现通用适配器接口
type GormAdapter struct {
	db   *gorm.DB
	tx   *gorm.DB // GORM事务状态
	name string
}

// NewGormAdapter 创建GORM适配器
func NewGormAdapter(db *gorm.DB) *GormAdapter {
	return &GormAdapter{
		db:   db,
		name: "GORM-Adapter",
	}
}

// NewGormTxAdapter 创建GORM事务适配器
func NewGormTxAdapter(tx *gorm.DB) *GormAdapter {
	return &GormAdapter{
		tx:   tx,
		name: "GORM-Transaction",
	}
}

// ==================== 通用适配器接口实现 ====================

// GetAdapterType 获取适配器类型
func (a *GormAdapter) GetAdapterType() string {
	return "GORM"
}

// GetAdapterName 获取适配器名称
func (a *GormAdapter) GetAdapterName() string {
	return a.name
}

// GetDialect 获取数据库方言
func (a *GormAdapter) GetDialect() string {
	db := a.getDB()
	if db != nil {
		return db.Dialector.Name()
	}
	return "unknown"
}

// SupportsORM ORM支持检测
func (a *GormAdapter) SupportsORM() bool {
	return true // GORM是全功能ORM
}

// SupportsUpsert Upsert支持检测
func (a *GormAdapter) SupportsUpsert() bool {
	return true // GORM支持Upsert操作
}

// SupportsBulkInsert 批量插入支持检测
func (a *GormAdapter) SupportsBulkInsert() bool {
	return true
}

// SupportsReturning RETURNING语句支持检测
func (a *GormAdapter) SupportsReturning() bool {
	dialect := a.GetDialect()
	return dialect == "postgres" || dialect == "sqlite"
}

// GetInstance 获取底层GORM实例
func (a *GormAdapter) GetInstance() interface{} {
	return a.getDB()
}

// GetStats 获取连接统计
func (a *GormAdapter) GetStats() ConnectionStats {
	db := a.getDB()
	if db != nil {
		sqlDB, err := db.DB()
		if err == nil {
			stats := sqlDB.Stats()
			return ConnectionStats{
				OpenConnections:   stats.OpenConnections,
				InUse:             stats.InUse,
				Idle:              stats.Idle,
				WaitCount:         stats.WaitCount,
				WaitDuration:      stats.WaitDuration,
				MaxIdleClosed:     stats.MaxIdleClosed,
				MaxLifetimeClosed: stats.MaxLifetimeClosed,
			}
		}
	}
	return ConnectionStats{}
}

// getDB 获取当前活跃的数据库连接
func (a *GormAdapter) getDB() *gorm.DB {
	if a.tx != nil {
		return a.tx
	}
	return a.db
}

// ==================== 批量操作实现 ====================

// BatchInsert 批量插入
func (a *GormAdapter) BatchInsert(ctx context.Context, table string, data []map[string]interface{}) error {
	if len(data) == 0 {
		return nil
	}

	db := a.getDB()
	if db == nil {
		return errors.NewError(errors.ErrorCodeDBError, "no database connection available")
	}

	// 使用GORM的CreateInBatches功能
	return db.WithContext(ctx).Table(table).CreateInBatches(data, len(data)).Error
}

// BatchUpdate 批量更新
func (a *GormAdapter) BatchUpdate(ctx context.Context, table string, data []map[string]interface{}, whereColumns []string) error {
	if len(data) == 0 {
		return nil
	}

	db := a.getDB()
	if db == nil {
		return errors.NewError(errors.ErrorCodeDBError, "no database connection available")
	}

	// 使用事务进行批量更新
	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, row := range data {
			var setClauses []string
			var setValues []interface{}
			var whereClauses []string
			var whereValues []interface{}

			for col, val := range row {
				isWhereColumn := false
				for _, whereCol := range whereColumns {
					if col == whereCol {
						isWhereColumn = true
						break
					}
				}

				if isWhereColumn {
					whereClauses = append(whereClauses, fmt.Sprintf("%s = ?", col))
					whereValues = append(whereValues, val)
				} else {
					setClauses = append(setClauses, fmt.Sprintf("%s = ?", col))
					setValues = append(setValues, val)
				}
			}

			if len(setClauses) == 0 || len(whereClauses) == 0 {
				continue
			}

			query := fmt.Sprintf("UPDATE %s SET %s WHERE %s",
				table,
				strings.Join(setClauses, ", "),
				strings.Join(whereClauses, " AND "))

			args := append(setValues, whereValues...)
			if err := tx.Exec(query, args...).Error; err != nil {
				return errors.NewErrorf(errors.ErrorCodeDBFailedUpdate, errors.MsgFailedToExecuteUpdate+": %v", err)
			}
		}
		return nil
	})
}

// 基本数据库操作
func (a *GormAdapter) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return a.db.Raw(query, args...).Rows()
}

func (a *GormAdapter) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return a.db.WithContext(ctx).Raw(query, args...).Rows()
}

func (a *GormAdapter) QueryRow(query string, args ...interface{}) *sql.Row {
	return a.db.Raw(query, args...).Row()
}

func (a *GormAdapter) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return a.db.WithContext(ctx).Raw(query, args...).Row()
}

func (a *GormAdapter) Exec(query string, args ...interface{}) (sql.Result, error) {
	result := a.db.Exec(query, args...)
	return &GormResult{result: result}, result.Error
}

func (a *GormAdapter) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	result := a.db.WithContext(ctx).Exec(query, args...)
	return &GormResult{result: result}, result.Error
}

func (a *GormAdapter) Begin() (TransactionInterface, error) {
	tx := a.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	return NewGormTxAdapter(tx), nil
}

func (a *GormAdapter) BeginTx(ctx context.Context, opts *sql.TxOptions) (TransactionInterface, error) {
	tx := a.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	return NewGormTxAdapter(tx), nil
}

func (a *GormAdapter) Prepare(query string) (StatementInterface, error) {
	return nil, errors.NewError(errors.ErrorCodeAdapterNotSupported, "gorm does not support prepared statements directly")
}

func (a *GormAdapter) PrepareContext(ctx context.Context, query string) (StatementInterface, error) {
	return nil, errors.NewError(errors.ErrorCodeAdapterNotSupported, "gorm does not support prepared statements directly")
}

func (a *GormAdapter) Ping() error {
	sqlDB, err := a.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

func (a *GormAdapter) PingContext(ctx context.Context) error {
	sqlDB, err := a.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.PingContext(ctx)
}

func (a *GormAdapter) Close() error {
	sqlDB, err := a.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (a *GormAdapter) Commit() error {
	return errors.NewError(errors.ErrorCodeBuilderNotInitialized, "cannot commit on non-transaction adapter")
}

func (a *GormAdapter) Rollback() error {
	return errors.NewError(errors.ErrorCodeBuilderNotInitialized, "cannot rollback on non-transaction adapter")
}

func (a *GormAdapter) GetDB() *gorm.DB {
	return a.db
}

// GORM特有方法
func (a *GormAdapter) Model(value interface{}) *GormAdapter {
	return &GormAdapter{db: a.db.Model(value)}
}

func (a *GormAdapter) Table(name string, args ...interface{}) *GormAdapter {
	return &GormAdapter{db: a.db.Table(name, args...)}
}

func (a *GormAdapter) Select(query interface{}, args ...interface{}) *GormAdapter {
	return &GormAdapter{db: a.db.Select(query, args...)}
}

func (a *GormAdapter) Where(query interface{}, args ...interface{}) *GormAdapter {
	return &GormAdapter{db: a.db.Where(query, args...)}
}

func (a *GormAdapter) Or(query interface{}, args ...interface{}) *GormAdapter {
	return &GormAdapter{db: a.db.Or(query, args...)}
}

func (a *GormAdapter) Not(query interface{}, args ...interface{}) *GormAdapter {
	return &GormAdapter{db: a.db.Not(query, args...)}
}

func (a *GormAdapter) Joins(query string, args ...interface{}) *GormAdapter {
	return &GormAdapter{db: a.db.Joins(query, args...)}
}

func (a *GormAdapter) Group(name string) *GormAdapter {
	return &GormAdapter{db: a.db.Group(name)}
}

func (a *GormAdapter) Having(query interface{}, args ...interface{}) *GormAdapter {
	return &GormAdapter{db: a.db.Having(query, args...)}
}

func (a *GormAdapter) Order(value interface{}) *GormAdapter {
	return &GormAdapter{db: a.db.Order(value)}
}

func (a *GormAdapter) Limit(limit int) *GormAdapter {
	return &GormAdapter{db: a.db.Limit(limit)}
}

func (a *GormAdapter) Offset(offset int) *GormAdapter {
	return &GormAdapter{db: a.db.Offset(offset)}
}

func (a *GormAdapter) Create(value interface{}) error {
	return a.db.Create(value).Error
}

func (a *GormAdapter) Save(value interface{}) error {
	return a.db.Save(value).Error
}

func (a *GormAdapter) First(dest interface{}, conds ...interface{}) error {
	return a.db.First(dest, conds...).Error
}

func (a *GormAdapter) Find(dest interface{}, conds ...interface{}) error {
	return a.db.Find(dest, conds...).Error
}

func (a *GormAdapter) Update(column string, value interface{}) error {
	return a.db.Update(column, value).Error
}

func (a *GormAdapter) Updates(values interface{}) error {
	return a.db.Updates(values).Error
}

func (a *GormAdapter) Delete(value interface{}, conds ...interface{}) error {
	return a.db.Delete(value, conds...).Error
}

func (a *GormAdapter) Count(count *int64) error {
	return a.db.Count(count).Error
}

func (a *GormAdapter) Scan(dest interface{}) error {
	return a.db.Scan(dest).Error
}

func (a *GormAdapter) Pluck(column string, dest interface{}) error {
	return a.db.Pluck(column, dest).Error
}

// ==================== GORM 事务适配器 (兼容性) ====================

// GormTxAdapterLegacy GORM事务适配器 (保留兼容性)
type GormTxAdapterLegacy struct {
	tx *gorm.DB
}

// NewGormTxAdapterLegacy 创建GORM事务适配器 (兼容性)
func NewGormTxAdapterLegacy(tx *gorm.DB) *GormTxAdapterLegacy {
	return &GormTxAdapterLegacy{tx: tx}
}

func (a *GormTxAdapterLegacy) Commit() error {
	return a.tx.Commit().Error
}

func (a *GormTxAdapterLegacy) Rollback() error {
	return a.tx.Rollback().Error
}

func (a *GormTxAdapterLegacy) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return a.tx.Raw(query, args...).Rows()
}

func (a *GormTxAdapterLegacy) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return a.tx.WithContext(ctx).Raw(query, args...).Rows()
}

func (a *GormTxAdapterLegacy) QueryRow(query string, args ...interface{}) *sql.Row {
	return a.tx.Raw(query, args...).Row()
}

func (a *GormTxAdapterLegacy) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return a.tx.WithContext(ctx).Raw(query, args...).Row()
}

func (a *GormTxAdapterLegacy) Exec(query string, args ...interface{}) (sql.Result, error) {
	result := a.tx.Exec(query, args...)
	return &GormResult{result: result}, result.Error
}

func (a *GormTxAdapterLegacy) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	result := a.tx.WithContext(ctx).Exec(query, args...)
	return &GormResult{result: result}, result.Error
}

func (a *GormTxAdapterLegacy) Begin() (TransactionInterface, error) {
	return nil, errors.NewError(errors.ErrorCodeNestedTransaction, "cannot begin transaction within transaction")
}

func (a *GormTxAdapterLegacy) BeginTx(ctx context.Context, opts *sql.TxOptions) (TransactionInterface, error) {
	return nil, errors.NewError(errors.ErrorCodeNestedTransaction, "cannot begin transaction within transaction")
}

func (a *GormTxAdapterLegacy) GetTx() *gorm.DB {
	return a.tx
}

// ==================== GORM Result 适配器 ====================

// GormResult GORM结果适配器
type GormResult struct {
	result *gorm.DB
}

func (r *GormResult) LastInsertId() (int64, error) {
	return 0, errors.NewError(errors.ErrorCodeUnsupported, "gorm does not support LastInsertId, use returning clause")
}

func (r *GormResult) RowsAffected() (int64, error) {
	return r.result.RowsAffected, r.result.Error
}

// ==================== 驱动适配器实现 ====================

// MySQLDriverAdapter MySQL驱动适配器
type MySQLDriverAdapter struct{}

func NewMySQLDriverAdapter() *MySQLDriverAdapter {
	return &MySQLDriverAdapter{}
}

func (a *MySQLDriverAdapter) DriverName() string {
	return "mysql"
}

func (a *MySQLDriverAdapter) SupportsFeature(feature string) bool {
	supportedFeatures := map[string]bool{
		"upsert":           true,
		"returning":        false,
		"json":             true,
		"cte":              true,
		"window_functions": true,
		"full_text":        true,
	}
	return supportedFeatures[feature]
}

func (a *MySQLDriverAdapter) QuoteIdentifier(identifier string) string {
	return fmt.Sprintf("`%s`", strings.ReplaceAll(identifier, "`", "``"))
}

func (a *MySQLDriverAdapter) QuoteString(str string) string {
	return fmt.Sprintf("'%s'", strings.ReplaceAll(str, "'", "''"))
}

func (a *MySQLDriverAdapter) BuildLimit(offset, limit int64) string {
	if offset > 0 {
		return fmt.Sprintf(" LIMIT %d, %d", offset, limit)
	}
	return fmt.Sprintf(" LIMIT %d", limit)
}

func (a *MySQLDriverAdapter) BuildUpsert(table string, data map[string]interface{}, conflictFields []string) (string, []interface{}) {
	fields := make([]string, 0, len(data))
	placeholders := make([]string, 0, len(data))
	values := make([]interface{}, 0, len(data))
	updates := make([]string, 0, len(data))

	for field, value := range data {
		fields = append(fields, a.QuoteIdentifier(field))
		placeholders = append(placeholders, "?")
		values = append(values, value)
		updates = append(updates, fmt.Sprintf("%s = VALUES(%s)", a.QuoteIdentifier(field), a.QuoteIdentifier(field)))
	}

	sql := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) ON DUPLICATE KEY UPDATE %s",
		table,
		strings.Join(fields, ", "),
		strings.Join(placeholders, ", "),
		strings.Join(updates, ", "))

	return sql, values
}

func (a *MySQLDriverAdapter) ConvertValue(value interface{}) (driver.Value, error) {
	return driver.DefaultParameterConverter.ConvertValue(value)
}

func (a *MySQLDriverAdapter) ConvertScanValue(src interface{}) (interface{}, error) {
	return src, nil
}

func (a *MySQLDriverAdapter) LastInsertId(result sql.Result) (int64, error) {
	return result.LastInsertId()
}

func (a *MySQLDriverAdapter) RowsAffected(result sql.Result) (int64, error) {
	return result.RowsAffected()
}

// PostgreSQLDriverAdapter PostgreSQL驱动适配器
type PostgreSQLDriverAdapter struct{}

func NewPostgreSQLDriverAdapter() *PostgreSQLDriverAdapter {
	return &PostgreSQLDriverAdapter{}
}

func (a *PostgreSQLDriverAdapter) DriverName() string {
	return "postgres"
}

func (a *PostgreSQLDriverAdapter) SupportsFeature(feature string) bool {
	supportedFeatures := map[string]bool{
		"upsert":           true,
		"returning":        true,
		"json":             true,
		"cte":              true,
		"window_functions": true,
		"full_text":        true,
		"arrays":           true,
	}
	return supportedFeatures[feature]
}

func (a *PostgreSQLDriverAdapter) QuoteIdentifier(identifier string) string {
	return fmt.Sprintf(`"%s"`, strings.ReplaceAll(identifier, `"`, `""`))
}

func (a *PostgreSQLDriverAdapter) QuoteString(str string) string {
	return fmt.Sprintf("'%s'", strings.ReplaceAll(str, "'", "''"))
}

func (a *PostgreSQLDriverAdapter) BuildLimit(offset, limit int64) string {
	if offset > 0 {
		return fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)
	}
	return fmt.Sprintf(" LIMIT %d", limit)
}

func (a *PostgreSQLDriverAdapter) BuildUpsert(table string, data map[string]interface{}, conflictFields []string) (string, []interface{}) {
	fields := make([]string, 0, len(data))
	placeholders := make([]string, 0, len(data))
	values := make([]interface{}, 0, len(data))
	updates := make([]string, 0, len(data))

	i := 1
	for field, value := range data {
		fields = append(fields, a.QuoteIdentifier(field))
		placeholders = append(placeholders, fmt.Sprintf("$%d", i))
		values = append(values, value)
		updates = append(updates, fmt.Sprintf("%s = EXCLUDED.%s", a.QuoteIdentifier(field), a.QuoteIdentifier(field)))
		i++
	}

	conflictClause := ""
	if len(conflictFields) > 0 {
		quotedFields := make([]string, len(conflictFields))
		for i, field := range conflictFields {
			quotedFields[i] = a.QuoteIdentifier(field)
		}
		conflictClause = fmt.Sprintf("(%s)", strings.Join(quotedFields, ", "))
	}

	sql := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) ON CONFLICT %s DO UPDATE SET %s",
		table,
		strings.Join(fields, ", "),
		strings.Join(placeholders, ", "),
		conflictClause,
		strings.Join(updates, ", "))

	return sql, values
}

func (a *PostgreSQLDriverAdapter) ConvertValue(value interface{}) (driver.Value, error) {
	return driver.DefaultParameterConverter.ConvertValue(value)
}

func (a *PostgreSQLDriverAdapter) ConvertScanValue(src interface{}) (interface{}, error) {
	return src, nil
}

func (a *PostgreSQLDriverAdapter) LastInsertId(result sql.Result) (int64, error) {
	return 0, errors.NewError(errors.ErrorCodeDBError, errors.MsgPostgresNotSupportLastInsertId)
}

func (a *PostgreSQLDriverAdapter) RowsAffected(result sql.Result) (int64, error) {
	return result.RowsAffected()
}

// ==================== 适配器工厂 ====================

// AdapterFactory 适配器工厂
type AdapterFactory struct {
	adapters map[string]func() DriverAdapterInterface
}

// NewAdapterFactory 创建适配器工厂
func NewAdapterFactory() *AdapterFactory {
	factory := &AdapterFactory{
		adapters: make(map[string]func() DriverAdapterInterface),
	}

	// 注册默认适配器
	factory.Register("mysql", func() DriverAdapterInterface {
		return NewMySQLDriverAdapter()
	})

	factory.Register("postgres", func() DriverAdapterInterface {
		return NewPostgreSQLDriverAdapter()
	})

	return factory
}

// Register 注册适配器
func (f *AdapterFactory) Register(name string, creator func() DriverAdapterInterface) {
	f.adapters[name] = creator
}

// Create 创建适配器
func (f *AdapterFactory) Create(name string) (DriverAdapterInterface, error) {
	creator, exists := f.adapters[name]
	if !exists {
		return nil, errors.NewErrorf(errors.ErrorCodeAdapterNotSupported, "unknown adapter: %s", name)
	}
	return creator(), nil
}

// GetSupportedAdapters 获取支持的适配器列表
func (f *AdapterFactory) GetSupportedAdapters() []string {
	adapters := make([]string, 0, len(f.adapters))
	for name := range f.adapters {
		adapters = append(adapters, name)
	}
	return adapters
}

// ==================== 全局适配器实例 ====================

var (
	// 全局适配器工厂
	globalAdapterFactory = NewAdapterFactory()
)

// RegisterAdapter 注册全局适配器
func RegisterAdapter(name string, creator func() DriverAdapterInterface) {
	globalAdapterFactory.Register(name, creator)
}

// CreateAdapter 创建适配器实例
func CreateAdapter(name string) (DriverAdapterInterface, error) {
	return globalAdapterFactory.Create(name)
}

// GetSupportedAdapters 获取支持的适配器
func GetSupportedAdapters() []string {
	return globalAdapterFactory.GetSupportedAdapters()
}

// ==================== 适配器包装器 ====================

// DatabaseAdapterWrapper 包装旧的DatabaseInterface为新的UniversalAdapterInterface
type DatabaseAdapterWrapper struct {
	db   DatabaseInterface
	name string
}

// NewDatabaseAdapterWrapper 创建适配器包装器
func NewDatabaseAdapterWrapper(db DatabaseInterface, name string) *DatabaseAdapterWrapper {
	return &DatabaseAdapterWrapper{
		db:   db,
		name: name,
	}
}

// ==================== 通用适配器接口实现 ====================

func (w *DatabaseAdapterWrapper) GetAdapterType() string {
	return "Legacy"
}

func (w *DatabaseAdapterWrapper) GetAdapterName() string {
	if w.name != "" {
		return w.name
	}
	return "Database-Adapter-Wrapper"
}

func (w *DatabaseAdapterWrapper) GetDialect() string {
	return "unknown"
}

func (w *DatabaseAdapterWrapper) SupportsORM() bool {
	return false
}

func (w *DatabaseAdapterWrapper) SupportsUpsert() bool {
	return false
}

func (w *DatabaseAdapterWrapper) SupportsBulkInsert() bool {
	return true
}

func (w *DatabaseAdapterWrapper) SupportsReturning() bool {
	return false
}

func (w *DatabaseAdapterWrapper) GetInstance() interface{} {
	return w.db
}

func (w *DatabaseAdapterWrapper) GetStats() ConnectionStats {
	return ConnectionStats{}
}

// 委托给底层DatabaseInterface
func (w *DatabaseAdapterWrapper) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return w.db.Query(query, args...)
}

func (w *DatabaseAdapterWrapper) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return w.db.QueryContext(ctx, query, args...)
}

func (w *DatabaseAdapterWrapper) QueryRow(query string, args ...interface{}) *sql.Row {
	return w.db.QueryRow(query, args...)
}

func (w *DatabaseAdapterWrapper) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return w.db.QueryRowContext(ctx, query, args...)
}

func (w *DatabaseAdapterWrapper) Exec(query string, args ...interface{}) (sql.Result, error) {
	return w.db.Exec(query, args...)
}

func (w *DatabaseAdapterWrapper) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return w.db.ExecContext(ctx, query, args...)
}

func (w *DatabaseAdapterWrapper) Begin() (TransactionInterface, error) {
	return w.db.Begin()
}

func (w *DatabaseAdapterWrapper) BeginTx(ctx context.Context, opts *sql.TxOptions) (TransactionInterface, error) {
	return w.db.BeginTx(ctx, opts)
}

func (w *DatabaseAdapterWrapper) Commit() error {
	return w.db.Commit()
}

func (w *DatabaseAdapterWrapper) Rollback() error {
	return w.db.Rollback()
}

func (w *DatabaseAdapterWrapper) Prepare(query string) (StatementInterface, error) {
	return w.db.Prepare(query)
}

func (w *DatabaseAdapterWrapper) PrepareContext(ctx context.Context, query string) (StatementInterface, error) {
	return w.db.PrepareContext(ctx, query)
}

func (w *DatabaseAdapterWrapper) Ping() error {
	return w.db.Ping()
}

func (w *DatabaseAdapterWrapper) PingContext(ctx context.Context) error {
	return w.db.PingContext(ctx)
}

func (w *DatabaseAdapterWrapper) Close() error {
	return w.db.Close()
}

// 批量操作的默认实现
func (w *DatabaseAdapterWrapper) BatchInsert(ctx context.Context, table string, data []map[string]interface{}) error {
	if len(data) == 0 {
		return nil
	}

	// 简单的批量插入实现
	for _, row := range data {
		var columns []string
		var placeholders []string
		var values []interface{}

		for col, val := range row {
			columns = append(columns, col)
			placeholders = append(placeholders, "?")
			values = append(values, val)
		}

		query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
			table,
			strings.Join(columns, ", "),
			strings.Join(placeholders, ", "))

		_, err := w.ExecContext(ctx, query, values...)
		if err != nil {
			return err
		}
	}

	return nil
}

func (w *DatabaseAdapterWrapper) BatchUpdate(ctx context.Context, table string, data []map[string]interface{}, whereColumns []string) error {
	if len(data) == 0 {
		return nil
	}

	// 简单的批量更新实现
	for _, row := range data {
		var setClauses []string
		var setValues []interface{}
		var whereClauses []string
		var whereValues []interface{}

		for col, val := range row {
			isWhereColumn := false
			for _, whereCol := range whereColumns {
				if col == whereCol {
					isWhereColumn = true
					break
				}
			}

			if isWhereColumn {
				whereClauses = append(whereClauses, col+" = ?")
				whereValues = append(whereValues, val)
			} else {
				setClauses = append(setClauses, col+" = ?")
				setValues = append(setValues, val)
			}
		}

		if len(setClauses) == 0 || len(whereClauses) == 0 {
			continue
		}

		query := fmt.Sprintf("UPDATE %s SET %s WHERE %s",
			table,
			strings.Join(setClauses, ", "),
			strings.Join(whereClauses, " AND "))

		args := append(setValues, whereValues...)
		_, err := w.ExecContext(ctx, query, args...)
		if err != nil {
			return err
		}
	}

	return nil
}

// NewUniversalAdapter 创建通用适配器
func NewUniversalAdapter(instance interface{}) (UniversalAdapterInterface, error) {
	switch db := instance.(type) {
	case *sqlx.DB:
		return NewSqlxAdapter(db), nil
	case *sqlx.Tx:
		return NewSqlxTxAdapter(db), nil
	case *gorm.DB:
		return NewGormAdapter(db), nil
	default:
		return nil, errors.NewErrorf(errors.ErrorCodeUnsupported, "unsupported database instance type: %T", instance)
	}
}

// WrapSqlxAsUniversal 包装sqlx.DB为通用适配器
func WrapSqlxAsUniversal(db *sqlx.DB) UniversalAdapterInterface {
	return NewSqlxAdapter(db)
}

// WrapGormAsUniversal 包装gorm.DB为通用适配器
func WrapGormAsUniversal(db *gorm.DB) UniversalAdapterInterface {
	return NewGormAdapter(db)
}

// AutoDetectAdapter 自动检测并创建适配器
func AutoDetectAdapter(instance interface{}) (UniversalAdapterInterface, error) {
	return NewUniversalAdapter(instance)
}
