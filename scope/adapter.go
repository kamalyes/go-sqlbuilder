/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-05-11 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-11 10:20:30
 * @FilePath: \go-sqlbuilder\scope\adapter.go
 * @Description: Scope SQL 适配器 - 将 ScopeData 转换为 SQL 查询过滤条件
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */
package scope

import (
	"github.com/kamalyes/go-sqlbuilder/constants"
	"github.com/kamalyes/go-sqlbuilder/repository"
)

// SQLScopeAdapter SQL 作用域适配器，将 ScopeData 转换为 SQL 查询过滤条件
type SQLScopeAdapter struct {
	data ScopeData
}

// NewSQLScopeAdapter 创建 SQL 作用域适配器
func NewSQLScopeAdapter(data ScopeData) *SQLScopeAdapter {
	return &SQLScopeAdapter{data: data}
}

// 编译期检查 SQLScopeAdapter 实现了 ScopeHook 接口
var _ repository.ScopeHook = (*SQLScopeAdapter)(nil)

// ApplyScope 将作用域条件应用到查询中：先移除已有的作用域字段过滤，再添加新的作用域过滤组
func (a *SQLScopeAdapter) ApplyScope(query *repository.Query) *repository.Query {
	mapping := a.data.Config.Mapping
	a.removeScopeFields(query, mapping)
	if group := a.buildFilterGroup(); group != nil && !group.IsEmpty() {
		query.WithFilterGroup(group)
	}
	return query
}

// removeScopeFields 移除查询中已存在的作用域字段过滤，避免重复条件
func (a *SQLScopeAdapter) removeScopeFields(query *repository.Query, mapping FieldMapping) {
	fieldsToRemove := scopeFieldSet(mapping)
	if len(fieldsToRemove) == 0 || len(query.Filters) == 0 {
		return
	}
	filtered := make([]*repository.Filter, 0, len(query.Filters))
	for _, f := range query.Filters {
		if f == nil {
			continue
		}
		if _, exists := fieldsToRemove[f.Field]; exists {
			continue
		}
		filtered = append(filtered, f)
	}
	query.Filters = filtered
}

// scopeFieldSet 根据字段映射构建需要移除的字段集合
func scopeFieldSet(mapping FieldMapping) map[string]struct{} {
	set := make(map[string]struct{})
	if mapping.TenantIDField != "" {
		set[mapping.TenantIDField] = struct{}{}
	}
	if mapping.PlatformIDField != "" {
		set[mapping.PlatformIDField] = struct{}{}
	}
	if mapping.RegionCodeField != "" {
		set[mapping.RegionCodeField] = struct{}{}
	}
	return set
}

// buildFilterGroup 根据域类型构建过滤条件组
func (a *SQLScopeAdapter) buildFilterGroup() *repository.FilterGroup {
	if a.data.IsOps() {
		return a.buildOpsFilterGroup()
	}
	if a.data.IsTenant() {
		return a.buildTenantFilterGroup()
	}
	return denyAllFilterGroup()
}

// buildOpsFilterGroup 构建 OPS 域过滤条件：
// - 全局管理员：不加任何过滤
// - 租户管理员：添加 tenant_id IN (指定租户列表)
func (a *SQLScopeAdapter) buildOpsFilterGroup() *repository.FilterGroup {
	mapping := a.data.Config.Mapping
	if a.data.HasGlobalScope() {
		return nil
	}

	tenantIds := a.collectOpsTenantIds()
	if len(tenantIds) == 0 || mapping.TenantIDField == "" {
		return denyAllFilterGroup()
	}

	return newTenantIdFilterGroup(mapping.TenantIDField, tenantIds)
}

// collectOpsTenantIds 从作用域条目中收集 OPS 可访问的租户ID列表
func (a *SQLScopeAdapter) collectOpsTenantIds() []string {
	seen := make(map[string]struct{})
	var result []string
	for _, e := range a.data.ScopeEntries {
		if a.data.IsTenantScope(e) {
			for _, tid := range e.TenantIds {
				if _, ok := seen[tid]; !ok {
					seen[tid] = struct{}{}
					result = append(result, tid)
				}
			}
		}
	}
	return result
}

// buildTenantFilterGroup 构建租户域过滤条件：
// - 全局用户/Owner：tenant_id = 'xxx'
// - 受限用户：tenant_id = 'xxx' AND (地区/平台 OR 条件)
func (a *SQLScopeAdapter) buildTenantFilterGroup() *repository.FilterGroup {
	if a.data.TenantID == "" {
		return denyAllFilterGroup()
	}

	if a.data.HasGlobalScope() {
		return a.buildTenantGlobalGroup()
	}

	return a.buildTenantScopedGroup()
}

// buildTenantGlobalGroup 构建租户全局过滤：仅添加 tenant_id 条件
func (a *SQLScopeAdapter) buildTenantGlobalGroup() *repository.FilterGroup {
	mapping := a.data.Config.Mapping
	if mapping.TenantIDField == "" {
		return denyAllFilterGroup()
	}
	group := repository.NewFilterGroup(constants.LOGIC_AND)
	group.AddFilter(&repository.Filter{
		Field:    mapping.TenantIDField,
		Operator: constants.OP_EQ,
		Value:    a.data.TenantID,
	})
	return group
}

// buildTenantScopedGroup 构建租户受限过滤：tenant_id + 地区/平台 OR 条件组
func (a *SQLScopeAdapter) buildTenantScopedGroup() *repository.FilterGroup {
	mapping := a.data.Config.Mapping
	if mapping.TenantIDField == "" {
		return denyAllFilterGroup()
	}
	outerGroup := repository.NewFilterGroup(constants.LOGIC_AND)

	outerGroup.AddFilter(&repository.Filter{
		Field:    mapping.TenantIDField,
		Operator: constants.OP_EQ,
		Value:    a.data.TenantID,
	})

	scopeOrGroup := a.buildScopeOrGroup()
	if scopeOrGroup != nil && !scopeOrGroup.IsEmpty() {
		outerGroup.AddGroup(scopeOrGroup)
	} else {
		outerGroup.AddGroup(denyAllFilterGroup())
	}

	return outerGroup
}

// buildScopeOrGroup 构建作用域 OR 条件组：
// 多个地区/平台条目之间是 OR 关系（满足任一即可访问）
func (a *SQLScopeAdapter) buildScopeOrGroup() *repository.FilterGroup {
	orGroup := repository.NewFilterGroup(constants.LOGIC_OR)

	for _, entry := range a.data.ScopeEntries {
		a.addEntryToOrGroup(orGroup, entry)
	}

	if orGroup.IsEmpty() {
		return nil
	}
	return orGroup
}

// addEntryToOrGroup 将单个作用域条目添加到 OR 条件组中
func (a *SQLScopeAdapter) addEntryToOrGroup(orGroup *repository.FilterGroup, entry *ScopeEntry) {
	if a.data.IsRegionScope(entry) {
		a.addRegionEntry(orGroup, entry)
	} else if a.data.IsPlatformScope(entry) {
		a.addPlatformEntry(orGroup, entry)
	}
}

// addRegionEntry 将地区级条目添加到 OR 组
func (a *SQLScopeAdapter) addRegionEntry(orGroup *repository.FilterGroup, entry *ScopeEntry) {
	if sub := a.buildRegionCondition(entry); sub != nil {
		orGroup.AddGroup(sub)
	}
}

// addPlatformEntry 将平台级条目添加到 OR 组，若子组为 OR 逻辑则展平
func (a *SQLScopeAdapter) addPlatformEntry(orGroup *repository.FilterGroup, entry *ScopeEntry) {
	sub := a.buildPlatformCondition(entry)
	if sub == nil {
		return
	}
	if sub.LogicOp == constants.LOGIC_OR {
		for _, g := range sub.Groups {
			orGroup.AddGroup(g)
		}
	} else {
		orGroup.AddGroup(sub)
	}
}

// buildRegionCondition 构建地区级条件：region_code IN ('MM','TH')
func (a *SQLScopeAdapter) buildRegionCondition(entry *ScopeEntry) *repository.FilterGroup {
	mapping := a.data.Config.Mapping
	if len(entry.RegionCodes) == 0 || mapping.RegionCodeField == "" {
		return nil
	}

	group := repository.NewFilterGroup(constants.LOGIC_AND)
	addValueFilter(group, mapping.RegionCodeField, entry.RegionCodes)
	return group
}

// buildPlatformCondition 构建平台级条件：
// - 单地区：region_code = 'MM' AND platform_id IN ('P1','P2')
// - 多地区混合：(region_code = 'MM') OR (region_code = 'SG' AND platform_id IN ('P1','P2'))
func (a *SQLScopeAdapter) buildPlatformCondition(entry *ScopeEntry) *repository.FilterGroup {
	mapping := a.data.Config.Mapping
	if len(entry.RegionPlatforms) == 0 {
		return nil
	}

	if !mapping.hasAnyField() {
		return nil
	}

	if isSingleRegionPlatform(entry, mapping) {
		return a.buildSingleRegionPlatformGroup(mapping, entry.RegionPlatforms[0])
	}

	return a.buildMultiRegionPlatformGroup(mapping, entry.RegionPlatforms)
}

// isSingleRegionPlatform 判断是否为单地区+双字段映射的优化路径
func isSingleRegionPlatform(entry *ScopeEntry, mapping FieldMapping) bool {
	return len(entry.RegionPlatforms) == 1 && mapping.RegionCodeField != "" && mapping.PlatformIDField != ""
}

// buildSingleRegionPlatformGroup 构建单地区平台条件：region_code = 'MM' AND platform_id IN ('P1','P2')
func (a *SQLScopeAdapter) buildSingleRegionPlatformGroup(mapping FieldMapping, rp *RegionPlatformEntry) *repository.FilterGroup {
	group := repository.NewFilterGroup(constants.LOGIC_AND)
	group.AddFilter(&repository.Filter{
		Field:    mapping.RegionCodeField,
		Operator: constants.OP_EQ,
		Value:    rp.RegionCode,
	})
	addValueFilter(group, mapping.PlatformIDField, rp.PlatformIds)
	return group
}

// buildMultiRegionPlatformGroup 构建多地区混合平台条件：每个地区一个 AND 组，之间用 OR 连接
func (a *SQLScopeAdapter) buildMultiRegionPlatformGroup(mapping FieldMapping, rps []*RegionPlatformEntry) *repository.FilterGroup {
	regionOrGroup := repository.NewFilterGroup(constants.LOGIC_OR)
	for _, rp := range rps {
		andGroup := a.buildRegionPlatformAndGroup(mapping, rp)
		if !andGroup.IsEmpty() {
			regionOrGroup.AddGroup(andGroup)
		}
	}

	if regionOrGroup.IsEmpty() {
		return nil
	}
	return regionOrGroup
}

// buildRegionPlatformAndGroup 构建单个地区的 AND 条件组：region_code = 'X' AND platform_id IN (...)
func (a *SQLScopeAdapter) buildRegionPlatformAndGroup(mapping FieldMapping, rp *RegionPlatformEntry) *repository.FilterGroup {
	andGroup := repository.NewFilterGroup(constants.LOGIC_AND)
	if mapping.RegionCodeField != "" {
		andGroup.AddFilter(&repository.Filter{
			Field:    mapping.RegionCodeField,
			Operator: constants.OP_EQ,
			Value:    rp.RegionCode,
		})
	}
	if mapping.PlatformIDField != "" && len(rp.PlatformIds) > 0 {
		addValueFilter(andGroup, mapping.PlatformIDField, rp.PlatformIds)
	}
	return andGroup
}

// hasAnyField 检查字段映射中是否至少配置了一个字段
func (m FieldMapping) hasAnyField() bool {
	return m.PlatformIDField != "" || m.RegionCodeField != ""
}

// addValueFilter 根据值数量自动选择 EQ 或 IN 操作符添加过滤条件
func addValueFilter(group *repository.FilterGroup, field string, values []string) {
	if len(values) == 1 {
		group.AddFilter(&repository.Filter{
			Field:    field,
			Operator: constants.OP_EQ,
			Value:    values[0],
		})
	} else {
		group.AddFilter(&repository.Filter{
			Field:    field,
			Operator: constants.OP_IN,
			Value:    values,
		})
	}
}

// newTenantIdFilterGroup 根据租户ID数量构建 tenant_id 过滤条件组
func newTenantIdFilterGroup(field string, tenantIds []string) *repository.FilterGroup {
	group := repository.NewFilterGroup(constants.LOGIC_AND)
	addValueFilter(group, field, tenantIds)
	return group
}

func denyAllFilterGroup() *repository.FilterGroup {
	group := repository.NewFilterGroup(constants.LOGIC_AND)
	group.AddFilter(&repository.Filter{
		Field:    "1 = 0",
		Operator: constants.OP_RAW,
	})
	return group
}

// ApplySQLScope 便捷函数：将作用域条件应用到查询中
func ApplySQLScope(query *repository.Query, data ScopeData) *repository.Query {
	return NewSQLScopeAdapter(data).ApplyScope(query)
}
