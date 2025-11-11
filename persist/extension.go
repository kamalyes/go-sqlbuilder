package persist

import (
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Order interface {
	Order(db *gorm.DB) *gorm.DB
}

type FieldOrder struct {
	Field    string
	Reversed bool
}

func (f *FieldOrder) Order(db *gorm.DB) *gorm.DB {
	return db.Order(clause.OrderByColumn{Column: clause.Column{Name: f.Field}, Desc: f.Reversed})
}

type Filter interface {
	Where(db *gorm.DB) *gorm.DB
}

type Filters []Filter

func (fs Filters) Where(db *gorm.DB) *gorm.DB {
	if len(fs) < 1 {
		return db
	}
	for _, f := range fs {
		db = f.Where(db)
	}
	return db
}

type EqFilter struct {
	Name  string
	Value any
}

func (eq *EqFilter) Where(db *gorm.DB) *gorm.DB {
	sql := fmt.Sprintf("%s = ?", eq.Name)
	return db.Where(sql, eq.Value)
}

type NeFilter struct {
	Name  string
	Value any
}

func (ne *NeFilter) Where(db *gorm.DB) *gorm.DB {
	sql := fmt.Sprintf("%s != ?", ne.Name)
	return db.Where(sql, ne.Value)
}

type OrFilter[T any] struct {
	Name   string
	Values []T
}

func (or *OrFilter[T]) Where(db *gorm.DB) *gorm.DB {
	if len(or.Values) < 1 {
		return db
	}
	sql := fmt.Sprintf("%s = ?", or.Name)
	return db.Or(sql, or.Values)
}

type InFilter[T any] struct {
	Name   string
	Values []T
}

func (in *InFilter[T]) Where(db *gorm.DB) *gorm.DB {
	sql := fmt.Sprintf("%s IN ?", in.Name)
	return db.Where(sql, in.Values)
}

type LikeFilter struct {
	Name  string
	Value string
}

func (f *LikeFilter) Where(db *gorm.DB) *gorm.DB {
	sql := fmt.Sprintf("%s LIKE ?", f.Name)
	return db.Where(sql, "%"+fmt.Sprintf("%v", f.Value)+"%")
}

type PrefixFilter struct {
	*LikeFilter
}

func (p *PrefixFilter) Where(db *gorm.DB) *gorm.DB {
	sql := fmt.Sprintf("%s LIKE ?", p.Name)
	return db.Where(sql, fmt.Sprintf("%v", p.Value)+"%")
}

type SuffixFilter struct {
	*LikeFilter
}

func (s *SuffixFilter) Where(db *gorm.DB) *gorm.DB {
	sql := fmt.Sprintf("%s LIKE ?", s.Name)
	return db.Where(sql, "%"+fmt.Sprintf("%v", s.Value))
}

type PeriodFilter struct {
	Name  string
	Start time.Time
	End   time.Time
}

func (p *PeriodFilter) Where(db *gorm.DB) *gorm.DB {
	if p.End.IsZero() {
		return db
	}
	if p.Start.IsZero() {
		sql := fmt.Sprintf("%s < ?", p.Name)
		return db.Where(sql, p.End)
	}
	sql := fmt.Sprintf("%s BETWEEN ? AND ?", p.Name)
	return db.Where(sql, p.Start, p.End)
}

type Number interface {
	~int | ~int32 | ~int64 | ~uint | ~uint32 | ~uint64 | ~float32 | ~float64
}
type RangeFilter[T Number] struct {
	Name  string
	Start T
	End   T
}

func (r RangeFilter[T]) Where(db *gorm.DB) *gorm.DB {
	if r.End == 0 {
		return db
	}
	if r.Start == 0 {
		sql := fmt.Sprintf("%s < ?", r.Name)
		return db.Where(sql, r.End)
	}
	return db.Where(fmt.Sprintf("%s >= ? AND %s < ?", r.Name, r.Name), r.Start, r.End)
}

type OriginalWhereFilter struct {
	Query string
	Args  []any
}

func (w *OriginalWhereFilter) Where(db *gorm.DB) *gorm.DB {
	return db.Where(w.Query, w.Args...)
}

type OriginalOrFilter struct {
	Filters Filters
}

func (o *OriginalOrFilter) Where(db *gorm.DB) *gorm.DB {
	return db.Or(o.Filters.Where(db))
}

type JSONArrayContainsFilter struct {
	Name  string
	Value any
}

func (jc *JSONArrayContainsFilter) Where(db *gorm.DB) *gorm.DB {
	var sql string
	switch db.Dialector.Name() {
	case "sqlite":
		sql = fmt.Sprintf("EXISTS (SELECT 1 FROM json_each(%s) WHERE value = ?)", jc.Name)
	case "mysql":
		sql = fmt.Sprintf("JSON_CONTAINS(CONVERT(%s USING utf8mb4), JSON_ARRAY(?))", jc.Name)
	case "postgres":
		sql = fmt.Sprintf("%s @> ARRAY[?]", jc.Name)
	default:
		sql = fmt.Sprintf("EXISTS (SELECT 1 FROM json_each(%s) WHERE value = ?)", jc.Name)
	}
	return db.Where(sql, jc.Value)
}

func NewOrder(field string, reversed bool) Order {
	return &FieldOrder{Field: field, Reversed: reversed}
}

func NewEqFilter(name string, value any) Filter {
	return &EqFilter{Name: name, Value: value}
}

func NewNeFilter(name string, value any) Filter {
	return &NeFilter{Name: name, Value: value}
}

func NewInFilter[T any](name string, values []T) Filter {
	return &InFilter[T]{Name: name, Values: values}
}

func NewOrFilter[T any](name string, values []T) Filter {
	return &OrFilter[T]{Name: name, Values: values}
}

func newLikeFilter(name string, value string) *LikeFilter {
	return &LikeFilter{Name: name, Value: value}
}

func NewLikeFilter(name string, value string) Filter {
	return newLikeFilter(name, value)
}

func NewPrefixFilter(name string, value string) Filter {
	return &PrefixFilter{LikeFilter: newLikeFilter(name, value)}
}

func NewSuffixFilter(name string, value string) Filter {
	return &SuffixFilter{LikeFilter: newLikeFilter(name, value)}
}

func NewPeriodFilter(name string, start, end time.Time) Filter {
	return &PeriodFilter{Name: name, Start: start, End: end}
}

func NewRangeFilter[T Number](name string, start, end T) Filter {
	return &RangeFilter[T]{Name: name, Start: start, End: end}
}

func NewOriginalWhereFilter(query string, args ...any) Filter {
	return &OriginalWhereFilter{Query: query, Args: args}
}

func NewOriginalOrFilter(filters Filters) Filter {
	return &OriginalOrFilter{Filters: filters}
}

func NewJSONArrayContainsFilter(name string, value any) Filter {
	return &JSONArrayContainsFilter{Name: name, Value: value}
}
