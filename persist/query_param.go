/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-11 21:13:15
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-13 11:07:41
 * @FilePath: \go-sqlbuilder\persist\query_param.go
 * @Description:
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package persist

import (
	"github.com/kamalyes/go-sqlbuilder/meta"
	"gorm.io/gorm"
)

type QueryParam struct {
	Filters Filters
	Page    *meta.Paging
	Options []Option
}

func NewQueryParam(filters []Filter, page *meta.Paging, opts ...Option) *QueryParam {
	return &QueryParam{Filters: filters, Page: page, Options: opts}
}

func (qa *QueryParam) GetOffset() int32 {
	if qa == nil {
		return 0
	}
	return int32(qa.Page.Offset)
}

func (qa *QueryParam) GetLimit() int32 {
	if qa == nil {
		return 0
	}
	return int32(qa.Page.Limit)
}

func (qa *QueryParam) Where(db *gorm.DB) *gorm.DB {
	if qa == nil {
		return db
	}
	ret := qa.Filters.Where(db)
	if qa.Page != nil {
		qa.Options = append(qa.Options, withPaging(qa.Page))
	}
	opts := qa.Options
	for _, opt := range opts {
		ret = opt(ret)
	}
	return ret
}

func WithOrders(orders ...Order) Option {
	return func(db *gorm.DB) *gorm.DB {
		ret := db
		for _, order := range orders {
			ret = order.Order(ret)
		}
		return ret
	}
}

func withPaging(page *meta.Paging) Option {
	return func(db *gorm.DB) *gorm.DB {
		ret := db
		if page != nil {
			ret = ret.Offset(int(page.Offset)).Limit(int(page.PageSize))
		}
		return ret
	}
}
