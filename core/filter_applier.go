/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-11 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-13 11:03:15
 * @FilePath: \go-sqlbuilder\constant\error.go
 * @Description:
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package core

import (
	"fmt"

	"github.com/kamalyes/go-sqlbuilder/constant"
	"github.com/kamalyes/go-sqlbuilder/errors"
	"gorm.io/gorm"
)

type FilterApplier struct {
	translator *OperatorTranslator
}

func NewFilterApplier() *FilterApplier {
	return &FilterApplier{
		translator: NewOperatorTranslator(),
	}
}

func (fa *FilterApplier) Apply(query *gorm.DB, filters []Filter) (*gorm.DB, error) {
	if query == nil {
		return nil, errors.NewError(
			errors.ErrorCodeInvalidInput,
			constant.ErrQueryNil,
		)
	}

	for _, filter := range filters {
		if err := fa.applyFilter(query, filter); err != nil {
			return nil, err
		}
	}

	return query, nil
}

func (fa *FilterApplier) ApplyFilters(query *gorm.DB, filters ...Filter) (*gorm.DB, error) {
	return fa.Apply(query, filters)
}

func (fa *FilterApplier) applyFilter(query *gorm.DB, filter Filter) error {
	if filter.Field == "" {
		return errors.NewError(
			errors.ErrorCodeInvalidFilterValue,
			constant.ErrFilterFieldEmpty,
		)
	}

	op := filter.Operator

	switch op {
	case constant.OP_EQ:
		return fa.applyEQ(query, filter)
	case constant.OP_NEQ:
		return fa.applyNEQ(query, filter)
	case constant.OP_GT:
		return fa.applyGT(query, filter)
	case constant.OP_GTE:
		return fa.applyGTE(query, filter)
	case constant.OP_LT:
		return fa.applyLT(query, filter)
	case constant.OP_LTE:
		return fa.applyLTE(query, filter)
	case constant.OP_LIKE:
		return fa.applyLIKE(query, filter)
	case constant.OP_NOT_LIKE:
		return fa.applyNOT_LIKE(query, filter)
	case constant.OP_IN:
		return fa.applyIN(query, filter)
	case constant.OP_NOT_IN:
		return fa.applyNOT_IN(query, filter)
	case constant.OP_BETWEEN:
		return fa.applyBETWEEN(query, filter)
	case constant.OP_IS_NULL:
		return fa.applyIS_NULL(query, filter)
	case constant.OP_IS_NOT_NULL:
		return fa.applyIS_NOT_NULL(query, filter)
	case constant.OP_FIND_IN_SET:
		return fa.applyFIND_IN_SET(query, filter)
	default:
		return errors.NewError(
			errors.ErrorCodeInvalidOperator,
			fmt.Sprintf(constant.ErrUnknownOperator, op),
		)
	}
}

func (fa *FilterApplier) applyEQ(query *gorm.DB, filter Filter) error {
	if filter.Value == nil {
		return errors.NewError(
			errors.ErrorCodeInvalidFilterValue,
			"EQ value cannot be nil",
		)
	}
	query.Where(fmt.Sprintf("%s = ?", filter.Field), filter.Value)
	return nil
}

func (fa *FilterApplier) applyNEQ(query *gorm.DB, filter Filter) error {
	if filter.Value == nil {
		return errors.NewError(
			errors.ErrorCodeInvalidFilterValue,
			"NEQ value cannot be nil",
		)
	}
	query.Where(fmt.Sprintf("%s != ?", filter.Field), filter.Value)
	return nil
}

func (fa *FilterApplier) applyGT(query *gorm.DB, filter Filter) error {
	if filter.Value == nil {
		return errors.NewError(
			errors.ErrorCodeInvalidFilterValue,
			"GT value cannot be nil",
		)
	}
	query.Where(fmt.Sprintf("%s > ?", filter.Field), filter.Value)
	return nil
}

func (fa *FilterApplier) applyGTE(query *gorm.DB, filter Filter) error {
	if filter.Value == nil {
		return errors.NewError(
			errors.ErrorCodeInvalidFilterValue,
			"GTE value cannot be nil",
		)
	}
	query.Where(fmt.Sprintf("%s >= ?", filter.Field), filter.Value)
	return nil
}

func (fa *FilterApplier) applyLT(query *gorm.DB, filter Filter) error {
	if filter.Value == nil {
		return errors.NewError(
			errors.ErrorCodeInvalidFilterValue,
			"LT value cannot be nil",
		)
	}
	query.Where(fmt.Sprintf("%s < ?", filter.Field), filter.Value)
	return nil
}

func (fa *FilterApplier) applyLTE(query *gorm.DB, filter Filter) error {
	if filter.Value == nil {
		return errors.NewError(
			errors.ErrorCodeInvalidFilterValue,
			"LTE value cannot be nil",
		)
	}
	query.Where(fmt.Sprintf("%s <= ?", filter.Field), filter.Value)
	return nil
}

func (fa *FilterApplier) applyLIKE(query *gorm.DB, filter Filter) error {
	if filter.Value == nil {
		return errors.NewError(
			errors.ErrorCodeInvalidFilterValue,
			"LIKE value cannot be nil",
		)
	}
	query.Where(fmt.Sprintf("%s LIKE ?", filter.Field), filter.Value)
	return nil
}

func (fa *FilterApplier) applyNOT_LIKE(query *gorm.DB, filter Filter) error {
	if filter.Value == nil {
		return errors.NewError(
			errors.ErrorCodeInvalidFilterValue,
			"NOT_LIKE value cannot be nil",
		)
	}
	query.Where(fmt.Sprintf("%s NOT LIKE ?", filter.Field), filter.Value)
	return nil
}

func (fa *FilterApplier) applyIN(query *gorm.DB, filter Filter) error {
	if filter.Value == nil {
		return errors.NewError(
			errors.ErrorCodeInvalidFilterValue,
			"IN value cannot be nil",
		)
	}
	query.Where(fmt.Sprintf("%s IN ?", filter.Field), filter.Value)
	return nil
}

func (fa *FilterApplier) applyNOT_IN(query *gorm.DB, filter Filter) error {
	if filter.Value == nil {
		return errors.NewError(
			errors.ErrorCodeInvalidFilterValue,
			"NOT_IN value cannot be nil",
		)
	}
	query.Where(fmt.Sprintf("%s NOT IN ?", filter.Field), filter.Value)
	return nil
}

func (fa *FilterApplier) applyBETWEEN(query *gorm.DB, filter Filter) error {
	if filter.Value == nil {
		return errors.NewError(
			errors.ErrorCodeInvalidFilterValue,
			"BETWEEN value cannot be nil",
		)
	}
	query.Where(fmt.Sprintf("%s BETWEEN ? AND ?", filter.Field), filter.Value)
	return nil
}

func (fa *FilterApplier) applyIS_NULL(query *gorm.DB, filter Filter) error {
	query.Where(fmt.Sprintf("%s IS NULL", filter.Field))
	return nil
}

func (fa *FilterApplier) applyIS_NOT_NULL(query *gorm.DB, filter Filter) error {
	query.Where(fmt.Sprintf("%s IS NOT NULL", filter.Field))
	return nil
}

func (fa *FilterApplier) applyFIND_IN_SET(query *gorm.DB, filter Filter) error {
	if filter.Value == nil {
		return errors.NewError(
			errors.ErrorCodeInvalidFilterValue,
			"FIND_IN_SET value cannot be nil",
		)
	}
	query.Where(fmt.Sprintf("FIND_IN_SET(?, %s)", filter.Field), filter.Value)
	return nil
}

type OrderApplier struct{}

func NewOrderApplier() *OrderApplier {
	return &OrderApplier{}
}

func (oa *OrderApplier) Apply(query *gorm.DB, orders []OrderBy) (*gorm.DB, error) {
	if query == nil {
		return nil, errors.NewError(
			errors.ErrorCodeInvalidInput,
			constant.ErrQueryNil,
		)
	}

	for _, order := range orders {
		if err := oa.applyOrder(query, order); err != nil {
			return nil, err
		}
	}

	return query, nil
}

func (oa *OrderApplier) ApplyOrders(query *gorm.DB, orders ...OrderBy) (*gorm.DB, error) {
	return oa.Apply(query, orders)
}

func (oa *OrderApplier) applyOrder(query *gorm.DB, order OrderBy) error {
	if order.Field == "" {
		return errors.NewError(
			errors.ErrorCodeInvalidInput,
			constant.ErrOrderFieldEmpty,
		)
	}

	orderStr := order.Order
	if orderStr == "" {
		orderStr = constant.OrderASC
	}

	query.Order(fmt.Sprintf("%s %s", order.Field, orderStr))
	return nil
}

type PaginationApplier struct{}

func NewPaginationApplier() *PaginationApplier {
	return &PaginationApplier{}
}

func (pa *PaginationApplier) Apply(query *gorm.DB, pagination *Pagination) (*gorm.DB, error) {
	if query == nil {
		return nil, errors.NewError(
			errors.ErrorCodeInvalidInput,
			constant.ErrQueryNil,
		)
	}

	if pagination == nil {
		return query, nil
	}

	if pagination.Limit <= 0 {
		return nil, errors.NewError(
			errors.ErrorCodePageSizeInvalid,
			constant.ErrPaginationLimitInvalid,
		)
	}

	if pagination.Offset < 0 {
		return nil, errors.NewError(
			errors.ErrorCodePageNumberInvalid,
			constant.ErrPaginationOffsetInvalid,
		)
	}

	query.Offset(pagination.Offset).Limit(pagination.Limit)
	return query, nil
}

type QueryApplier struct {
	filterApplier     *FilterApplier
	orderApplier      *OrderApplier
	paginationApplier *PaginationApplier
}

func NewQueryApplier() *QueryApplier {
	return &QueryApplier{
		filterApplier:     NewFilterApplier(),
		orderApplier:      NewOrderApplier(),
		paginationApplier: NewPaginationApplier(),
	}
}

func (qa *QueryApplier) ApplyCondition(query *gorm.DB, condition *QueryCondition) (*gorm.DB, error) {
	if query == nil {
		return nil, errors.NewError(
			errors.ErrorCodeInvalidInput,
			constant.ErrQueryNil,
		)
	}

	if condition == nil {
		return query, nil
	}

	var err error

	if len(condition.Filters) > 0 {
		query, err = qa.filterApplier.Apply(query, condition.Filters)
		if err != nil {
			return nil, err
		}
	}

	if len(condition.Orders) > 0 {
		query, err = qa.orderApplier.Apply(query, condition.Orders)
		if err != nil {
			return nil, err
		}
	}

	if condition.Pagination != nil {
		query, err = qa.paginationApplier.Apply(query, condition.Pagination)
		if err != nil {
			return nil, err
		}
	}

	return query, nil
}

func (qa *QueryApplier) ApplyFilters(query *gorm.DB, filters []Filter) (*gorm.DB, error) {
	return qa.filterApplier.Apply(query, filters)
}

func (qa *QueryApplier) ApplyOrders(query *gorm.DB, orders []OrderBy) (*gorm.DB, error) {
	return qa.orderApplier.Apply(query, orders)
}

func (qa *QueryApplier) ApplyPagination(query *gorm.DB, pagination *Pagination) (*gorm.DB, error) {
	return qa.paginationApplier.Apply(query, pagination)
}
