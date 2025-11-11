/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-11 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-11 00:00:00
 * @FilePath: \go-sqlbuilder\builder.go
 * @Description: 高效、完整的SQL查询构建器 - 支持所有主流框架
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */

package sqlbuilder

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"time"
)

// Builder 通用SQL查询构建器
type Builder struct {
	adapter UniversalAdapterInterface
	ctx     context.Context
	timeout time.Duration

	// SQL构建组件
	table       string
	tableAlias  string
	distinct    bool
	columns     []string
	joins       []string
	wheres      []string
	havings     []string
	groupByCols []string
	orderByCols []string
	limitVal    int64
	offsetVal   int64

	// 操作数据
	insertData  map[string]interface{}
	updateData  map[string]interface{}
	deleteWhere bool

	// 参数
	args      []interface{}
	queryType string // select, insert, update, delete
}

// New 创建新的查询构建器
func New(dbInstance interface{}) (*Builder, error) {
	adapter, err := AutoDetectAdapter(dbInstance)
	if err != nil {
		return nil, err
	}

	return &Builder{
		adapter:     adapter,
		ctx:         context.Background(),
		timeout:     30 * time.Second,
		columns:     []string{},
		joins:       []string{},
		wheres:      []string{},
		havings:     []string{},
		groupByCols: []string{},
		orderByCols: []string{},
		insertData:  make(map[string]interface{}),
		updateData:  make(map[string]interface{}),
		args:        []interface{}{},
		queryType:   "select",
	}, nil
}

// WithContext 设置上下文
func (b *Builder) WithContext(ctx context.Context) *Builder {
	b.ctx = ctx
	return b
}

// WithTimeout 设置超时时间
func (b *Builder) WithTimeout(timeout time.Duration) *Builder {
	b.timeout = timeout
	return b
}

// ==================== 表操作 ====================

// Table 设置表名
func (b *Builder) Table(table string) *Builder {
	b.table = table
	b.queryType = "select"
	return b
}

// As 设置表别名
func (b *Builder) As(alias string) *Builder {
	b.tableAlias = alias
	return b
}

// ==================== SELECT ====================

// Select 选择列
func (b *Builder) Select(columns ...string) *Builder {
	b.queryType = "select"
	if len(columns) == 0 {
		b.columns = []string{"*"}
	} else {
		b.columns = columns
	}
	return b
}

// SelectRaw 选择原始SQL
func (b *Builder) SelectRaw(sql string, args ...interface{}) *Builder {
	b.queryType = "select"
	b.columns = []string{sql}
	b.args = append(b.args, args...)
	return b
}

// Distinct 去重
func (b *Builder) Distinct() *Builder {
	b.distinct = true
	return b
}

// ==================== JOIN ====================

// Join 内连接
func (b *Builder) Join(table, on string, args ...interface{}) *Builder {
	return b.addJoin("INNER", table, on, args...)
}

// LeftJoin 左连接
func (b *Builder) LeftJoin(table, on string, args ...interface{}) *Builder {
	return b.addJoin("LEFT", table, on, args...)
}

// RightJoin 右连接
func (b *Builder) RightJoin(table, on string, args ...interface{}) *Builder {
	return b.addJoin("RIGHT", table, on, args...)
}

// FullJoin 全连接
func (b *Builder) FullJoin(table, on string, args ...interface{}) *Builder {
	return b.addJoin("FULL", table, on, args...)
}

// CrossJoin 交叉连接
func (b *Builder) CrossJoin(table string) *Builder {
	b.joins = append(b.joins, fmt.Sprintf("CROSS JOIN %s", table))
	return b
}

func (b *Builder) addJoin(joinType, table, on string, args ...interface{}) *Builder {
	b.joins = append(b.joins, fmt.Sprintf("%s JOIN %s ON %s", joinType, table, on))
	b.args = append(b.args, args...)
	return b
}

// ==================== WHERE ====================

// Where WHERE条件
func (b *Builder) Where(column string, operator string, value interface{}) *Builder {
	return b.addWhere("AND", column, operator, value)
}

// OrWhere OR WHERE条件
func (b *Builder) OrWhere(column string, operator string, value interface{}) *Builder {
	return b.addWhere("OR", column, operator, value)
}

// WhereRaw 原始WHERE
func (b *Builder) WhereRaw(sql string, args ...interface{}) *Builder {
	if len(b.wheres) > 0 {
		b.wheres = append(b.wheres, fmt.Sprintf("AND %s", sql))
	} else {
		b.wheres = append(b.wheres, sql)
	}
	b.args = append(b.args, args...)
	return b
}

// OrWhereRaw 原始OR WHERE
func (b *Builder) OrWhereRaw(sql string, args ...interface{}) *Builder {
	if len(b.wheres) > 0 {
		b.wheres = append(b.wheres, fmt.Sprintf("OR %s", sql))
	} else {
		b.wheres = append(b.wheres, sql)
	}
	b.args = append(b.args, args...)
	return b
}

// WhereIn IN条件
func (b *Builder) WhereIn(column string, values ...interface{}) *Builder {
	if len(values) == 0 {
		return b
	}
	placeholders := strings.Repeat("?,", len(values))
	placeholders = placeholders[:len(placeholders)-1]
	sql := fmt.Sprintf("%s IN (%s)", column, placeholders)
	return b.WhereRaw(sql, values...)
}

// WhereNotIn NOT IN条件
func (b *Builder) WhereNotIn(column string, values ...interface{}) *Builder {
	if len(values) == 0 {
		return b
	}
	placeholders := strings.Repeat("?,", len(values))
	placeholders = placeholders[:len(placeholders)-1]
	sql := fmt.Sprintf("%s NOT IN (%s)", column, placeholders)
	return b.WhereRaw(sql, values...)
}

// WhereBetween BETWEEN条件
func (b *Builder) WhereBetween(column string, min, max interface{}) *Builder {
	return b.WhereRaw(fmt.Sprintf("%s BETWEEN ? AND ?", column), min, max)
}

// WhereNull NULL条件
func (b *Builder) WhereNull(column string) *Builder {
	return b.WhereRaw(fmt.Sprintf("%s IS NULL", column))
}

// WhereNotNull NOT NULL条件
func (b *Builder) WhereNotNull(column string) *Builder {
	return b.WhereRaw(fmt.Sprintf("%s IS NOT NULL", column))
}

// WhereLike LIKE条件
func (b *Builder) WhereLike(column, value string) *Builder {
	return b.Where(column, "LIKE", value)
}

func (b *Builder) addWhere(boolean string, column, operator string, value interface{}) *Builder {
	if len(b.wheres) > 0 {
		b.wheres = append(b.wheres, fmt.Sprintf("%s %s %s ?", boolean, column, operator))
	} else {
		b.wheres = append(b.wheres, fmt.Sprintf("%s %s ?", column, operator))
	}
	b.args = append(b.args, value)
	return b
}

// ==================== GROUP BY / HAVING ====================

// GroupBy 分组
func (b *Builder) GroupBy(columns ...string) *Builder {
	b.groupByCols = append(b.groupByCols, columns...)
	return b
}

// Having HAVING条件
func (b *Builder) Having(column, operator string, value interface{}) *Builder {
	b.havings = append(b.havings, fmt.Sprintf("%s %s ?", column, operator))
	b.args = append(b.args, value)
	return b
}

// HavingRaw 原始HAVING
func (b *Builder) HavingRaw(sql string, args ...interface{}) *Builder {
	b.havings = append(b.havings, sql)
	b.args = append(b.args, args...)
	return b
}

// ==================== ORDER BY ====================

// OrderBy 排序（升序）
func (b *Builder) OrderBy(column string) *Builder {
	b.orderByCols = append(b.orderByCols, fmt.Sprintf("%s ASC", column))
	return b
}

// OrderByDesc 排序（降序）
func (b *Builder) OrderByDesc(column string) *Builder {
	b.orderByCols = append(b.orderByCols, fmt.Sprintf("%s DESC", column))
	return b
}

// OrderByRaw 原始ORDER BY
func (b *Builder) OrderByRaw(sql string) *Builder {
	b.orderByCols = append(b.orderByCols, sql)
	return b
}

// ==================== LIMIT / OFFSET ====================

// Limit 限制结果数
func (b *Builder) Limit(limit int64) *Builder {
	b.limitVal = limit
	return b
}

// Offset 偏移
func (b *Builder) Offset(offset int64) *Builder {
	b.offsetVal = offset
	return b
}

// Paginate 分页
func (b *Builder) Paginate(page, pageSize int64) *Builder {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	b.offsetVal = (page - 1) * pageSize
	b.limitVal = pageSize
	return b
}

// ==================== INSERT ====================

// Insert 插入
func (b *Builder) Insert(data map[string]interface{}) *Builder {
	b.queryType = "insert"
	b.insertData = data
	return b
}

// InsertGetID 插入并返回ID
func (b *Builder) InsertGetID(data map[string]interface{}) (int64, error) {
	b.queryType = "insert"
	b.insertData = data
	result, err := b.Exec()
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// ==================== UPDATE ====================

// Update 更新
func (b *Builder) Update(data map[string]interface{}) (*Builder, error) {
	b.queryType = "update"
	b.updateData = data
	return b, nil
}

// Set 单个字段更新
func (b *Builder) Set(column string, value interface{}) *Builder {
	b.queryType = "update"
	if b.updateData == nil {
		b.updateData = make(map[string]interface{})
	}
	b.updateData[column] = value
	return b
}

// Increment 增加
func (b *Builder) Increment(column string, value int64) error {
	b.queryType = "update"
	b.updateData[column] = struct {
		Op    string
		Value int64
	}{Op: "+", Value: value}
	_, err := b.Exec()
	return err
}

// Decrement 减少
func (b *Builder) Decrement(column string, value int64) error {
	b.queryType = "update"
	b.updateData[column] = struct {
		Op    string
		Value int64
	}{Op: "-", Value: value}
	_, err := b.Exec()
	return err
}

// ==================== DELETE ====================

// Delete 删除
func (b *Builder) Delete() *Builder {
	b.queryType = "delete"
	b.deleteWhere = true
	return b
}

// ==================== 执行方法 ====================

// ToSQL 生成SQL
func (b *Builder) ToSQL() (string, []interface{}) {
	var sql strings.Builder

	switch b.queryType {
	case "select":
		sql.WriteString(b.buildSelect())
	case "insert":
		sql.WriteString(b.buildInsert())
	case "update":
		sql.WriteString(b.buildUpdate())
	case "delete":
		sql.WriteString(b.buildDelete())
	}

	return sql.String(), b.args
}

// First 获取第一条记录
func (b *Builder) First(dest interface{}) error {
	b.Limit(1)
	sql, args := b.ToSQL()
	row := b.adapter.QueryRowContext(b.ctx, sql, args...)
	return row.Scan(dest)
}

// Get 获取结果集
func (b *Builder) Get(dest interface{}) error {
	sql, args := b.ToSQL()
	rows, err := b.adapter.QueryContext(b.ctx, sql, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	// 使用反射处理扫描
	destVal := reflect.ValueOf(dest)
	if destVal.Kind() != reflect.Ptr {
		return fmt.Errorf("dest must be a pointer")
	}

	sliceVal := destVal.Elem()
	if sliceVal.Kind() != reflect.Slice {
		return fmt.Errorf("dest must be a slice pointer")
	}

	elemType := sliceVal.Type().Elem()

	for rows.Next() {
		elem := reflect.New(elemType).Elem()

		// 扫描行到结构体
		fields := make([]interface{}, 0)
		for i := 0; i < elem.NumField(); i++ {
			fields = append(fields, elem.Field(i).Addr().Interface())
		}

		if err := rows.Scan(fields...); err != nil {
			return err
		}

		sliceVal.Set(reflect.Append(sliceVal, elem))
	}

	return rows.Err()
}

// Exec 执行SQL
func (b *Builder) Exec() (sql.Result, error) {
	sql, args := b.ToSQL()
	return b.adapter.ExecContext(b.ctx, sql, args...)
}

// Count 获取计数
func (b *Builder) Count() (int64, error) {
	oldCols := b.columns
	oldLimit := b.limitVal
	oldOffset := b.offsetVal

	b.columns = []string{"COUNT(*) as count"}
	b.limitVal = 0
	b.offsetVal = 0

	sql, args := b.ToSQL()
	row := b.adapter.QueryRowContext(b.ctx, sql, args...)

	var count int64
	err := row.Scan(&count)

	b.columns = oldCols
	b.limitVal = oldLimit
	b.offsetVal = oldOffset

	return count, err
}

// Exists 检查是否存在
func (b *Builder) Exists() (bool, error) {
	count, err := b.Count()
	return count > 0, err
}

// ==================== 批量操作 ====================

// BatchInsert 批量插入
func (b *Builder) BatchInsert(data []map[string]interface{}) error {
	if len(data) == 0 {
		return nil
	}
	return b.adapter.BatchInsert(b.ctx, b.table, data)
}

// BatchUpdate 批量更新
func (b *Builder) BatchUpdate(data []map[string]interface{}, whereColumns []string) error {
	if len(data) == 0 {
		return nil
	}
	return b.adapter.BatchUpdate(b.ctx, b.table, data, whereColumns)
}

// ==================== 事务支持 ====================

// Transaction 事务
func (b *Builder) Transaction(fn func(*Builder) error) error {
	tx, err := b.adapter.BeginTx(b.ctx, nil)
	if err != nil {
		return err
	}

	txBuilder := &Builder{
		adapter: tx.(UniversalAdapterInterface),
		ctx:     b.ctx,
		timeout: b.timeout,
		table:   b.table,
		columns: b.columns,
	}

	if err := fn(txBuilder); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// ==================== 工具方法 ====================

// GetAdapter 获取适配器
func (b *Builder) GetAdapter() UniversalAdapterInterface {
	return b.adapter
}

// Close 关闭连接
func (b *Builder) Close() error {
	return b.adapter.Close()
}

// Ping 检查连接
func (b *Builder) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return b.adapter.PingContext(ctx)
}

// ==================== 私有方法 ====================

func (b *Builder) buildSelect() string {
	var sql strings.Builder

	sql.WriteString("SELECT")
	if b.distinct {
		sql.WriteString(" DISTINCT")
	}

	if len(b.columns) == 0 {
		sql.WriteString(" *")
	} else {
		sql.WriteString(" ")
		sql.WriteString(strings.Join(b.columns, ", "))
	}

	sql.WriteString(fmt.Sprintf(" FROM %s", b.table))
	if b.tableAlias != "" {
		sql.WriteString(fmt.Sprintf(" AS %s", b.tableAlias))
	}

	// JOINs
	for _, j := range b.joins {
		sql.WriteString(" ")
		sql.WriteString(j)
	}

	// WHERE
	if len(b.wheres) > 0 {
		sql.WriteString(" WHERE ")
		sql.WriteString(strings.Join(b.wheres, " "))
	}

	// GROUP BY
	if len(b.groupByCols) > 0 {
		sql.WriteString(fmt.Sprintf(" GROUP BY %s", strings.Join(b.groupByCols, ", ")))
	}

	// HAVING
	if len(b.havings) > 0 {
		sql.WriteString(fmt.Sprintf(" HAVING %s", strings.Join(b.havings, " AND ")))
	}

	// ORDER BY
	if len(b.orderByCols) > 0 {
		sql.WriteString(fmt.Sprintf(" ORDER BY %s", strings.Join(b.orderByCols, ", ")))
	}

	// LIMIT
	if b.limitVal > 0 {
		sql.WriteString(fmt.Sprintf(" LIMIT %d", b.limitVal))
	}

	// OFFSET
	if b.offsetVal >= 0 {
		sql.WriteString(fmt.Sprintf(" OFFSET %d", b.offsetVal))
	}

	return sql.String()
}

func (b *Builder) buildInsert() string {
	var cols []string
	var placeholders []string

	for k, v := range b.insertData {
		cols = append(cols, k)
		placeholders = append(placeholders, "?")
		b.args = append(b.args, v)
	}

	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		b.table,
		strings.Join(cols, ", "),
		strings.Join(placeholders, ", "))
}

func (b *Builder) buildUpdate() string {
	if len(b.updateData) == 0 {
		return ""
	}

	var sql strings.Builder
	sql.WriteString(fmt.Sprintf("UPDATE %s SET ", b.table))

	var setParts []string
	for k, v := range b.updateData {
		setParts = append(setParts, fmt.Sprintf("%s = ?", k))
		b.args = append(b.args, v)
	}

	sql.WriteString(strings.Join(setParts, ", "))

	// WHERE
	if len(b.wheres) > 0 {
		sql.WriteString(" WHERE ")
		sql.WriteString(strings.Join(b.wheres, " "))
	}

	return sql.String()
}

func (b *Builder) buildDelete() string {
	var sql strings.Builder
	sql.WriteString(fmt.Sprintf("DELETE FROM %s", b.table))

	// WHERE
	if len(b.wheres) > 0 {
		sql.WriteString(" WHERE ")
		sql.WriteString(strings.Join(b.wheres, " "))
	}

	return sql.String()
}
