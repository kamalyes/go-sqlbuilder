/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-05-11 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-11 10:01:52
 * @FilePath: \go-sqlbuilder\scope\adapter_test.go
 * @Description: Scope SQL 适配器测试 - 覆盖 OPS/租户域各种作用域场景
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */
package scope

import (
	"testing"

	sqlconstants "github.com/kamalyes/go-sqlbuilder/constants"
	sqlrepo "github.com/kamalyes/go-sqlbuilder/repository"
	"github.com/stretchr/testify/assert"
)

func assertDenyAllScope(t *testing.T, group *sqlrepo.FilterGroup) {
	t.Helper()
	if !assert.NotNil(t, group) {
		return
	}
	assert.Len(t, group.Filters, 1)
	assert.Equal(t, sqlconstants.OP_RAW, group.Filters[0].Operator)
	assert.Equal(t, "1 = 0", group.Filters[0].Field)
}

// ==============================================================================
// OPS 域
// ==============================================================================

// 场景：OPS全局管理员（Domain=2, ScopeType=Global）
// 预期：不加任何过滤条件，FilterGroup 为 nil
func TestSQLScopeAdapter_OpsGlobal(t *testing.T) {
	data := NewScopeData()
	data.Domain = 2
	data.ScopeEntries = []*ScopeEntry{{ScopeType: 1}}

	query := sqlrepo.NewQuery()
	result := ApplySQLScope(query, data)
	assert.Nil(t, result.FilterGroup)
}

// 场景：OPS租户管理员，仅管理单个租户T1
// 预期：生成 AND { tenant_id = 'T1' }，使用 EQ 操作符
func TestSQLScopeAdapter_OpsTenantSingle(t *testing.T) {
	data := NewScopeData()
	data.Domain = 2
	data.ScopeEntries = []*ScopeEntry{
		{ScopeType: 4, TenantIds: []string{"T1"}},
	}

	query := sqlrepo.NewQuery()
	result := ApplySQLScope(query, data)
	assert.NotNil(t, result.FilterGroup)

	fg := result.FilterGroup
	assert.Equal(t, sqlconstants.LOGIC_AND, fg.LogicOp)
	assert.Len(t, fg.Filters, 1)
	assert.Equal(t, "tenant_id", fg.Filters[0].Field)
	assert.Equal(t, sqlconstants.OP_EQ, fg.Filters[0].Operator)
	assert.Equal(t, "T1", fg.Filters[0].Value)
}

// 场景：OPS租户管理员，管理多个租户T1、T2
// 预期：生成 AND { tenant_id IN ('T1','T2') }，多值自动切换为 IN 操作符
func TestSQLScopeAdapter_OpsTenantMultiple(t *testing.T) {
	data := NewScopeData()
	data.Domain = 2
	data.ScopeEntries = []*ScopeEntry{
		{ScopeType: 4, TenantIds: []string{"T1", "T2"}},
	}

	query := sqlrepo.NewQuery()
	result := ApplySQLScope(query, data)
	assert.NotNil(t, result.FilterGroup)

	fg := result.FilterGroup
	assert.Equal(t, sqlconstants.OP_IN, fg.Filters[0].Operator)
}

// 场景：OPS域但 ScopeEntries 为空
// 预期：权限缺失时默认拒绝，避免退化成全量查询
func TestSQLScopeAdapter_OpsNoEntries(t *testing.T) {
	data := NewScopeData()
	data.Domain = 2
	data.ScopeEntries = []*ScopeEntry{}

	query := sqlrepo.NewQuery()
	result := ApplySQLScope(query, data)
	assertDenyAllScope(t, result.FilterGroup)
}

func TestSQLScopeAdapter_UnknownDomainDenyAll(t *testing.T) {
	data := NewScopeData()
	data.Domain = 0

	query := sqlrepo.NewQuery()
	result := ApplySQLScope(query, data)
	assertDenyAllScope(t, result.FilterGroup)
}

// 场景：OPS域有租户条目但显式清空 TenantIDField（模拟未配置字段映射）
// 预期：缺少字段映射，无法生成条件，FilterGroup 为 nil
func TestSQLScopeAdapter_OpsNoTenantIdField(t *testing.T) {
	data := NewScopeData(WithFieldMapping(FieldMapping{}))
	data.Domain = 2
	data.ScopeEntries = []*ScopeEntry{
		{ScopeType: 4, TenantIds: []string{"T1"}},
	}

	query := sqlrepo.NewQuery()
	result := ApplySQLScope(query, data)
	assertDenyAllScope(t, result.FilterGroup)
}

// 场景：OPS域多条租户条目存在重复的租户ID（T1,T2 和 T2,T3）
// 预期：自动去重，最终生成 tenant_id IN ('T1','T2','T3')，共3个值
func TestSQLScopeAdapter_OpsDuplicateTenantIds(t *testing.T) {
	data := NewScopeData()
	data.Domain = 2
	data.ScopeEntries = []*ScopeEntry{
		{ScopeType: 4, TenantIds: []string{"T1", "T2"}},
		{ScopeType: 4, TenantIds: []string{"T2", "T3"}},
	}

	query := sqlrepo.NewQuery()
	result := ApplySQLScope(query, data)
	assert.NotNil(t, result.FilterGroup)

	fg := result.FilterGroup
	assert.Equal(t, sqlconstants.OP_IN, fg.Filters[0].Operator)
	ids := fg.Filters[0].Value.([]string)
	assert.Len(t, ids, 3)
	assert.Contains(t, ids, "T1")
	assert.Contains(t, ids, "T2")
	assert.Contains(t, ids, "T3")
}

// ==============================================================================
// 租户域（全局）
// ==============================================================================

// 场景：租户全局管理员（Domain=1, ScopeType=Global, TenantID=T001）
// 预期：生成 AND { tenant_id = 'T001' }，仅限制租户维度
func TestSQLScopeAdapter_TenantGlobal(t *testing.T) {
	data := NewScopeData()
	data.Domain = 1
	data.TenantID = "T001"
	data.ScopeEntries = []*ScopeEntry{{ScopeType: 1}}

	query := sqlrepo.NewQuery()
	result := ApplySQLScope(query, data)
	assert.NotNil(t, result.FilterGroup)

	fg := result.FilterGroup
	assert.Equal(t, sqlconstants.LOGIC_AND, fg.LogicOp)
	assert.Len(t, fg.Filters, 1)
	assert.Equal(t, "tenant_id", fg.Filters[0].Field)
	assert.Equal(t, "T001", fg.Filters[0].Value)
}

// 场景：租户全局管理员但显式清空 TenantIDField
// 预期：无法映射租户字段时默认拒绝，避免跨租户查询
func TestSQLScopeAdapter_TenantGlobalNoTenantIdField(t *testing.T) {
	data := NewScopeData(WithFieldMapping(FieldMapping{}))
	data.Domain = 1
	data.TenantID = "T001"
	data.ScopeEntries = []*ScopeEntry{{ScopeType: 1}}

	query := sqlrepo.NewQuery()
	result := ApplySQLScope(query, data)
	assertDenyAllScope(t, result.FilterGroup)
}

func TestSQLScopeAdapter_TenantNoTenantIDDenyAll(t *testing.T) {
	data := NewScopeData()
	data.Domain = 1
	data.ScopeEntries = []*ScopeEntry{{ScopeType: 1}}

	query := sqlrepo.NewQuery()
	result := ApplySQLScope(query, data)
	assertDenyAllScope(t, result.FilterGroup)
}

// ==============================================================================
// 租户域（地区级）
// ==============================================================================

// 场景：租户地区级权限，单个地区MM
// 预期：生成 AND { tenant_id='T001', OR { AND { region_code = 'MM' } } }
// 结构：外层AND包含tenant_id + OR组（内含一个地区AND条目）
func TestSQLScopeAdapter_TenantRegionSingle(t *testing.T) {
	data := NewScopeData()
	data.Domain = 1
	data.TenantID = "T001"
	data.ScopeEntries = []*ScopeEntry{
		{ScopeType: 2, RegionCodes: []string{"MM"}},
	}

	query := sqlrepo.NewQuery()
	result := ApplySQLScope(query, data)
	assert.NotNil(t, result.FilterGroup)

	fg := result.FilterGroup
	assert.Equal(t, sqlconstants.LOGIC_AND, fg.LogicOp)
	assert.Len(t, fg.Filters, 1)
	assert.Equal(t, "tenant_id", fg.Filters[0].Field)
	assert.Len(t, fg.Groups, 1)

	regionOr := fg.Groups[0]
	assert.Equal(t, sqlconstants.LOGIC_OR, regionOr.LogicOp)
	assert.Len(t, regionOr.Groups, 1)

	innerAnd := regionOr.Groups[0]
	assert.Equal(t, sqlconstants.LOGIC_AND, innerAnd.LogicOp)
	assert.Len(t, innerAnd.Filters, 1)
	assert.Equal(t, "region_code", innerAnd.Filters[0].Field)
	assert.Equal(t, "MM", innerAnd.Filters[0].Value)
}

// 场景：租户地区级权限，多个地区MM、TH
// 预期：地区码多于1个时自动使用 IN 操作符：region_code IN ('MM','TH')
func TestSQLScopeAdapter_TenantRegionMultiple(t *testing.T) {
	data := NewScopeData()
	data.Domain = 1
	data.TenantID = "T001"
	data.ScopeEntries = []*ScopeEntry{
		{ScopeType: 2, RegionCodes: []string{"MM", "TH"}},
	}

	query := sqlrepo.NewQuery()
	result := ApplySQLScope(query, data)
	assert.NotNil(t, result.FilterGroup)

	fg := result.FilterGroup
	regionOr := fg.Groups[0]
	innerAnd := regionOr.Groups[0]
	assert.Equal(t, sqlconstants.OP_IN, innerAnd.Filters[0].Operator)
}

// ==============================================================================
// 租户域（平台级 - 单地区）
// ==============================================================================

// 场景：租户平台级权限，单地区MM下有多个平台P1、P2
// 预期：走单地区优化路径，生成 AND { region_code = 'MM', platform_id IN ('P1','P2') }
// 单地区+双字段映射时，地区和平台条件合并到同一个AND组
func TestSQLScopeAdapter_TenantPlatformSingleRegion(t *testing.T) {
	data := NewScopeData()
	data.Domain = 1
	data.TenantID = "T001"
	data.ScopeEntries = []*ScopeEntry{
		{
			ScopeType: 3,
			RegionPlatforms: []*RegionPlatformEntry{
				{RegionCode: "MM", PlatformIds: []string{"P1", "P2"}},
			},
		},
	}

	query := sqlrepo.NewQuery()
	result := ApplySQLScope(query, data)
	assert.NotNil(t, result.FilterGroup)

	fg := result.FilterGroup
	assert.Equal(t, "tenant_id", fg.Filters[0].Field)

	scopeOr := fg.Groups[0]
	platformAnd := scopeOr.Groups[0]
	assert.Equal(t, sqlconstants.LOGIC_AND, platformAnd.LogicOp)
	assert.Len(t, platformAnd.Filters, 2)
	assert.Equal(t, "region_code", platformAnd.Filters[0].Field)
	assert.Equal(t, "MM", platformAnd.Filters[0].Value)
	assert.Equal(t, "platform_id", platformAnd.Filters[1].Field)
	assert.Equal(t, sqlconstants.OP_IN, platformAnd.Filters[1].Operator)
}

// ==============================================================================
// 租户域（平台级 - 多地区混合）
// ==============================================================================

// 场景：租户平台级权限，多地区混合（MM有P1，SG有P1和P2）
// 预期：走多地区路径，生成 OR { AND { region_code='MM', platform_id='P1' }, AND { region_code='SG', platform_id IN ('P1','P2') } }
// MM单平台用EQ，SG多平台用IN
func TestSQLScopeAdapter_TenantPlatformMultiRegion(t *testing.T) {
	data := NewScopeData()
	data.Domain = 1
	data.TenantID = "T001"
	data.ScopeEntries = []*ScopeEntry{
		{
			ScopeType: 3,
			RegionPlatforms: []*RegionPlatformEntry{
				{RegionCode: "MM", PlatformIds: []string{"P1"}},
				{RegionCode: "SG", PlatformIds: []string{"P1", "P2"}},
			},
		},
	}

	query := sqlrepo.NewQuery()
	result := ApplySQLScope(query, data)
	assert.NotNil(t, result.FilterGroup)

	fg := result.FilterGroup
	scopeOr := fg.Groups[0]
	assert.Equal(t, sqlconstants.LOGIC_OR, scopeOr.LogicOp)
	assert.Len(t, scopeOr.Groups, 2)

	mmGroup := scopeOr.Groups[0]
	assert.Equal(t, sqlconstants.LOGIC_AND, mmGroup.LogicOp)
	assert.Len(t, mmGroup.Filters, 2)
	assert.Equal(t, "region_code", mmGroup.Filters[0].Field)
	assert.Equal(t, "MM", mmGroup.Filters[0].Value)
	assert.Equal(t, "platform_id", mmGroup.Filters[1].Field)
	assert.Equal(t, "P1", mmGroup.Filters[1].Value)

	sgGroup := scopeOr.Groups[1]
	assert.Equal(t, sqlconstants.LOGIC_AND, sgGroup.LogicOp)
	assert.Equal(t, "SG", sgGroup.Filters[0].Value)
	assert.Equal(t, sqlconstants.OP_IN, sgGroup.Filters[1].Operator)
}

// ==============================================================================
// 租户域（混合地区+平台条目）
// ==============================================================================

// 场景：租户同时拥有地区级条目（MM）和平台级条目（SG+P1）
// 预期：两种条目合并到同一个 OR 组，生成 OR { AND { region_code IN ('MM') }, AND { region_code='SG', platform_id='P1' } }
func TestSQLScopeAdapter_TenantMixedScopeEntries(t *testing.T) {
	data := NewScopeData()
	data.Domain = 1
	data.TenantID = "T001"
	data.ScopeEntries = []*ScopeEntry{
		{ScopeType: 2, RegionCodes: []string{"MM"}},
		{
			ScopeType: 3,
			RegionPlatforms: []*RegionPlatformEntry{
				{RegionCode: "SG", PlatformIds: []string{"P1"}},
			},
		},
	}

	query := sqlrepo.NewQuery()
	result := ApplySQLScope(query, data)
	assert.NotNil(t, result.FilterGroup)

	fg := result.FilterGroup
	scopeOr := fg.Groups[0]
	assert.Equal(t, sqlconstants.LOGIC_OR, scopeOr.LogicOp)
	assert.Len(t, scopeOr.Groups, 2)
}

// ==============================================================================
// 自定义字段名
// ==============================================================================

// 场景：使用自定义字段名（TenantIDField="tid", RegionCodeField="rc", PlatformIDField="pid"）
// 预期：生成的过滤条件使用自定义字段名而非默认值
func TestSQLScopeAdapter_CustomFieldNames(t *testing.T) {
	data := NewScopeData(
		WithTenantIDField("tid"),
		WithRegionCodeField("rc"),
		WithPlatformIDField("pid"),
	)
	data.Domain = 1
	data.TenantID = "T001"
	data.ScopeEntries = []*ScopeEntry{
		{
			ScopeType: 3,
			RegionPlatforms: []*RegionPlatformEntry{
				{RegionCode: "MM", PlatformIds: []string{"P1"}},
			},
		},
	}

	query := sqlrepo.NewQuery()
	result := ApplySQLScope(query, data)
	assert.NotNil(t, result.FilterGroup)

	fg := result.FilterGroup
	assert.Equal(t, "tid", fg.Filters[0].Field)

	scopeOr := fg.Groups[0]
	platformAnd := scopeOr.Groups[0]
	assert.Equal(t, "rc", platformAnd.Filters[0].Field)
	assert.Equal(t, "pid", platformAnd.Filters[1].Field)
}

// ==============================================================================
// removeScopeFields
// ==============================================================================

// 场景：查询中已有 tenant_id 和 region_code 的旧过滤条件，应用新的作用域
// 预期：旧的作用域字段过滤被移除，仅保留非作用域字段（name）
func TestSQLScopeAdapter_RemoveScopeFields(t *testing.T) {
	data := NewScopeData()
	data.Domain = 1
	data.TenantID = "T001"
	data.ScopeEntries = []*ScopeEntry{{ScopeType: 1}}

	query := sqlrepo.NewQuery()
	query.Filters = []*sqlrepo.Filter{
		{Field: "tenant_id", Operator: sqlconstants.OP_EQ, Value: "old"},
		{Field: "region_code", Operator: sqlconstants.OP_EQ, Value: "old"},
		{Field: "name", Operator: sqlconstants.OP_EQ, Value: "keep"},
	}

	result := ApplySQLScope(query, data)
	assert.Len(t, result.Filters, 1)
	assert.Equal(t, "name", result.Filters[0].Field)
}

// 场景：查询 Filters 中包含 nil 元素
// 预期：nil 元素被安全跳过，仅移除作用域字段，保留非作用域字段
func TestSQLScopeAdapter_RemoveScopeFieldsWithNilFilter(t *testing.T) {
	data := NewScopeData()
	data.Domain = 1
	data.TenantID = "T001"
	data.ScopeEntries = []*ScopeEntry{{ScopeType: 1}}

	query := sqlrepo.NewQuery()
	query.Filters = []*sqlrepo.Filter{
		nil,
		{Field: "tenant_id", Operator: sqlconstants.OP_EQ, Value: "old"},
		{Field: "name", Operator: sqlconstants.OP_EQ, Value: "keep"},
	}

	result := ApplySQLScope(query, data)
	assert.Len(t, result.Filters, 1)
	assert.Equal(t, "name", result.Filters[0].Field)
}

// 场景：查询 Filters 为空切片
// 预期：不报错，正常生成作用域 FilterGroup
func TestSQLScopeAdapter_RemoveScopeFieldsEmptyFilters(t *testing.T) {
	data := NewScopeData()
	data.Domain = 1
	data.TenantID = "T001"
	data.ScopeEntries = []*ScopeEntry{{ScopeType: 1}}

	query := sqlrepo.NewQuery()
	query.Filters = []*sqlrepo.Filter{}

	result := ApplySQLScope(query, data)
	assert.NotNil(t, result.FilterGroup)
}

// 场景：显式清空所有字段映射时，查询中有非作用域字段
// 预期：没有需要移除的字段，所有原有过滤条件保持不变
func TestSQLScopeAdapter_RemoveScopeFieldsNoMapping(t *testing.T) {
	data := NewScopeData(WithFieldMapping(FieldMapping{}))
	data.Domain = 1
	data.TenantID = "T001"
	data.ScopeEntries = []*ScopeEntry{{ScopeType: 1}}

	query := sqlrepo.NewQuery()
	query.Filters = []*sqlrepo.Filter{
		{Field: "name", Operator: sqlconstants.OP_EQ, Value: "keep"},
	}

	result := ApplySQLScope(query, data)
	assert.Len(t, result.Filters, 1)
	assert.Equal(t, "name", result.Filters[0].Field)
}

// ==============================================================================
// 边界情况
// ==============================================================================

// 场景：Domain 值既不是租户也不是 OPS（Domain=99）
// 预期：无法识别域类型，FilterGroup 为 nil
func TestSQLScopeAdapter_UnknownDomain(t *testing.T) {
	data := NewScopeData()
	data.Domain = 99
	data.TenantID = "T001"

	query := sqlrepo.NewQuery()
	result := ApplySQLScope(query, data)
	assertDenyAllScope(t, result.FilterGroup)
}

// 场景：租户域但 TenantID 为空字符串
// 预期：无法构建有效的租户条件，FilterGroup 为 nil
func TestSQLScopeAdapter_EmptyTenantID(t *testing.T) {
	data := NewScopeData()
	data.Domain = 1
	data.TenantID = ""

	query := sqlrepo.NewQuery()
	result := ApplySQLScope(query, data)
	assertDenyAllScope(t, result.FilterGroup)
}

// 场景：租户域有地区级条目但显式清空所有字段映射
// 预期：没有可映射的字段，无法生成条件，FilterGroup 为 nil
func TestSQLScopeAdapter_NoFieldMapping(t *testing.T) {
	data := NewScopeData(WithFieldMapping(FieldMapping{}))
	data.Domain = 1
	data.TenantID = "T001"
	data.ScopeEntries = []*ScopeEntry{{ScopeType: 2, RegionCodes: []string{"MM"}}}

	query := sqlrepo.NewQuery()
	result := ApplySQLScope(query, data)
	assertDenyAllScope(t, result.FilterGroup)
}

// 场景：平台级条目，只配置了 PlatformIDField 没有配置 RegionCodeField
// 预期：只生成 platform_id 条件，不生成 region_code 条件
// AND { platform_id = 'P1' }（仅有1个filter）
func TestSQLScopeAdapter_PlatformNoRegionField(t *testing.T) {
	data := NewScopeData(
		WithRegionCodeField(""),
	)
	data.Domain = 1
	data.TenantID = "T001"
	data.ScopeEntries = []*ScopeEntry{
		{
			ScopeType: 3,
			RegionPlatforms: []*RegionPlatformEntry{
				{RegionCode: "MM", PlatformIds: []string{"P1"}},
			},
		},
	}

	query := sqlrepo.NewQuery()
	result := ApplySQLScope(query, data)
	assert.NotNil(t, result.FilterGroup)

	fg := result.FilterGroup
	scopeOr := fg.Groups[0]
	platformAnd := scopeOr.Groups[0]
	assert.Len(t, platformAnd.Filters, 1)
	assert.Equal(t, "platform_id", platformAnd.Filters[0].Field)
}

// 场景：平台级条目，单地区单平台（MM + P1）
// 预期：单平台ID自动使用 EQ 操作符而非 IN：platform_id = 'P1'
func TestSQLScopeAdapter_PlatformSinglePlatformId(t *testing.T) {
	data := NewScopeData()
	data.Domain = 1
	data.TenantID = "T001"
	data.ScopeEntries = []*ScopeEntry{
		{
			ScopeType: 3,
			RegionPlatforms: []*RegionPlatformEntry{
				{RegionCode: "MM", PlatformIds: []string{"P1"}},
			},
		},
	}

	query := sqlrepo.NewQuery()
	result := ApplySQLScope(query, data)
	assert.NotNil(t, result.FilterGroup)

	fg := result.FilterGroup
	scopeOr := fg.Groups[0]
	platformAnd := scopeOr.Groups[0]
	assert.Equal(t, sqlconstants.OP_EQ, platformAnd.Filters[1].Operator)
	assert.Equal(t, "P1", platformAnd.Filters[1].Value)
}

// 场景：平台级条目但 RegionPlatforms 为空切片
// 预期：受限作用域无有效平台范围时默认拒绝
func TestSQLScopeAdapter_PlatformEmptyRegionPlatforms(t *testing.T) {
	data := NewScopeData()
	data.Domain = 1
	data.TenantID = "T001"
	data.ScopeEntries = []*ScopeEntry{
		{ScopeType: 3, RegionPlatforms: []*RegionPlatformEntry{}},
	}

	query := sqlrepo.NewQuery()
	result := ApplySQLScope(query, data)
	assert.NotNil(t, result.FilterGroup)

	fg := result.FilterGroup
	assert.Len(t, fg.Filters, 1)
	assert.Equal(t, "tenant_id", fg.Filters[0].Field)
	assert.Len(t, fg.Groups, 1)
	assertDenyAllScope(t, fg.Groups[0])
}

// 场景：平台级条目但显式清空所有字段映射（hasAnyField=false）
// 预期：无法映射租户字段时默认拒绝，避免跨租户查询
func TestSQLScopeAdapter_PlatformNoFieldMappingAtAll(t *testing.T) {
	data := NewScopeData(WithFieldMapping(FieldMapping{}))
	data.Domain = 1
	data.TenantID = "T001"
	data.ScopeEntries = []*ScopeEntry{
		{
			ScopeType: 3,
			RegionPlatforms: []*RegionPlatformEntry{
				{RegionCode: "MM", PlatformIds: []string{"P1"}},
			},
		},
	}

	query := sqlrepo.NewQuery()
	result := ApplySQLScope(query, data)
	assertDenyAllScope(t, result.FilterGroup)
}

// 场景：租户受限作用域但显式清空 TenantIDField，仅保留 RegionCodeField
// 预期：无法映射租户字段时默认拒绝，避免只按地区跨租户查询
func TestSQLScopeAdapter_TenantScopedNoTenantIdField(t *testing.T) {
	data := NewScopeData(
		WithTenantIDField(""),
	)
	data.Domain = 1
	data.TenantID = "T001"
	data.ScopeEntries = []*ScopeEntry{
		{ScopeType: 2, RegionCodes: []string{"MM"}},
	}

	query := sqlrepo.NewQuery()
	result := ApplySQLScope(query, data)
	assertDenyAllScope(t, result.FilterGroup)
}

// 场景：租户域但 ScopeEntries 中只有全局条目（ScopeType=1）
// 预期：走 buildTenantGlobalGroup 路径，仅生成 tenant_id = 'T001'
func TestSQLScopeAdapter_TenantScopedOnlyGlobalEntries(t *testing.T) {
	data := NewScopeData()
	data.Domain = 1
	data.TenantID = "T001"
	data.ScopeEntries = []*ScopeEntry{
		{ScopeType: 1},
	}

	query := sqlrepo.NewQuery()
	result := ApplySQLScope(query, data)
	assert.NotNil(t, result.FilterGroup)

	fg := result.FilterGroup
	assert.Equal(t, sqlconstants.LOGIC_AND, fg.LogicOp)
	assert.Len(t, fg.Filters, 1)
	assert.Equal(t, "tenant_id", fg.Filters[0].Field)
}

// 场景：多地区平台级，其中MM的PlatformIds为空（该地区全部平台），SG有指定平台P1
// 预期：MM条目仅生成 region_code = 'MM'（1个filter），SG条目生成 region_code='SG' + platform_id='P1'（2个filter）
// 两个条目通过 OR 连接
func TestSQLScopeAdapter_MultiRegionPlatformWithEmptyPlatformIds(t *testing.T) {
	data := NewScopeData()
	data.Domain = 1
	data.TenantID = "T001"
	data.ScopeEntries = []*ScopeEntry{
		{
			ScopeType: 3,
			RegionPlatforms: []*RegionPlatformEntry{
				{RegionCode: "MM", PlatformIds: []string{}},
				{RegionCode: "SG", PlatformIds: []string{"P1"}},
			},
		},
	}

	query := sqlrepo.NewQuery()
	result := ApplySQLScope(query, data)
	assert.NotNil(t, result.FilterGroup)

	fg := result.FilterGroup
	scopeOr := fg.Groups[0]
	assert.Equal(t, sqlconstants.LOGIC_OR, scopeOr.LogicOp)
	assert.Len(t, scopeOr.Groups, 2)

	mmGroup := scopeOr.Groups[0]
	assert.Len(t, mmGroup.Filters, 1)
	assert.Equal(t, "region_code", mmGroup.Filters[0].Field)
}

// 场景：地区级条目但 RegionCodes 为空切片
// 预期：受限作用域无有效地区范围时默认拒绝
func TestSQLScopeAdapter_RegionConditionEmptyRegionCodes(t *testing.T) {
	data := NewScopeData()
	data.Domain = 1
	data.TenantID = "T001"
	data.ScopeEntries = []*ScopeEntry{
		{ScopeType: 2, RegionCodes: []string{}},
	}

	query := sqlrepo.NewQuery()
	result := ApplySQLScope(query, data)
	assert.NotNil(t, result.FilterGroup)

	fg := result.FilterGroup
	assert.Len(t, fg.Filters, 1)
	assert.Equal(t, "tenant_id", fg.Filters[0].Field)
	assert.Len(t, fg.Groups, 1)
	assertDenyAllScope(t, fg.Groups[0])
}

// 场景：单地区平台级条目，只配置了 RegionCodeField 没有配置 PlatformIDField
// 预期：走多地区路径（isSingleRegionPlatform=false），仅生成 region_code = 'MM'（1个filter）
func TestSQLScopeAdapter_SingleRegionPlatformOnlyRegionField(t *testing.T) {
	data := NewScopeData(
		WithPlatformIDField(""),
	)
	data.Domain = 1
	data.TenantID = "T001"
	data.ScopeEntries = []*ScopeEntry{
		{
			ScopeType: 3,
			RegionPlatforms: []*RegionPlatformEntry{
				{RegionCode: "MM", PlatformIds: []string{"P1"}},
			},
		},
	}

	query := sqlrepo.NewQuery()
	result := ApplySQLScope(query, data)
	assert.NotNil(t, result.FilterGroup)

	fg := result.FilterGroup
	scopeOr := fg.Groups[0]
	platformAnd := scopeOr.Groups[0]
	assert.Len(t, platformAnd.Filters, 1)
	assert.Equal(t, "region_code", platformAnd.Filters[0].Field)
	assert.Equal(t, "MM", platformAnd.Filters[0].Value)
}

// 场景：单地区平台级条目，只配置了 PlatformIDField 没有配置 RegionCodeField
// 预期：走多地区路径（isSingleRegionPlatform=false），仅生成 platform_id = 'P1'（1个filter）
func TestSQLScopeAdapter_SingleRegionPlatformOnlyPlatformField(t *testing.T) {
	data := NewScopeData(
		WithRegionCodeField(""),
	)
	data.Domain = 1
	data.TenantID = "T001"
	data.ScopeEntries = []*ScopeEntry{
		{
			ScopeType: 3,
			RegionPlatforms: []*RegionPlatformEntry{
				{RegionCode: "MM", PlatformIds: []string{"P1"}},
			},
		},
	}

	query := sqlrepo.NewQuery()
	result := ApplySQLScope(query, data)
	assert.NotNil(t, result.FilterGroup)

	fg := result.FilterGroup
	scopeOr := fg.Groups[0]
	platformAnd := scopeOr.Groups[0]
	assert.Len(t, platformAnd.Filters, 1)
	assert.Equal(t, "platform_id", platformAnd.Filters[0].Field)
}

// 场景：addValueFilter 工具函数
// 预期：单值时使用 EQ 操作符，多值时使用 IN 操作符
func TestSQLScopeAdapter_AddValueFilter(t *testing.T) {
	t.Run("单值用EQ", func(t *testing.T) {
		group := sqlrepo.NewFilterGroup(sqlconstants.LOGIC_AND)
		addValueFilter(group, "field", []string{"V1"})
		assert.Len(t, group.Filters, 1)
		assert.Equal(t, sqlconstants.OP_EQ, group.Filters[0].Operator)
		assert.Equal(t, "V1", group.Filters[0].Value)
	})

	t.Run("多值用IN", func(t *testing.T) {
		group := sqlrepo.NewFilterGroup(sqlconstants.LOGIC_AND)
		addValueFilter(group, "field", []string{"V1", "V2"})
		assert.Len(t, group.Filters, 1)
		assert.Equal(t, sqlconstants.OP_IN, group.Filters[0].Operator)
	})
}

// 场景：newTenantIdFilterGroup 工具函数
// 预期：单个租户ID时使用 EQ，多个租户ID时使用 IN，逻辑运算均为 AND
func TestSQLScopeAdapter_NewTenantIdFilterGroup(t *testing.T) {
	t.Run("单个租户ID", func(t *testing.T) {
		group := newTenantIdFilterGroup("tenant_id", []string{"T1"})
		assert.Equal(t, sqlconstants.LOGIC_AND, group.LogicOp)
		assert.Len(t, group.Filters, 1)
		assert.Equal(t, sqlconstants.OP_EQ, group.Filters[0].Operator)
	})

	t.Run("多个租户ID", func(t *testing.T) {
		group := newTenantIdFilterGroup("tenant_id", []string{"T1", "T2"})
		assert.Equal(t, sqlconstants.LOGIC_AND, group.LogicOp)
		assert.Len(t, group.Filters, 1)
		assert.Equal(t, sqlconstants.OP_IN, group.Filters[0].Operator)
	})
}

// 场景：多地区平台级，所有地区的 PlatformIds 均为空，且只配置了 PlatformIDField（无 RegionCodeField）
// 预期：每个 buildRegionPlatformAndGroup 都生成空 AND 组，buildMultiRegionPlatformGroup 返回 nil
// 最终追加拒绝条件，避免仅凭 tenant_id 放大平台受限权限
func TestSQLScopeAdapter_MultiRegionPlatformAllEmptyPlatformIds(t *testing.T) {
	data := NewScopeData(
		WithRegionCodeField(""),
	)
	data.Domain = 1
	data.TenantID = "T001"
	data.ScopeEntries = []*ScopeEntry{
		{
			ScopeType: 3,
			RegionPlatforms: []*RegionPlatformEntry{
				{RegionCode: "MM", PlatformIds: []string{}},
				{RegionCode: "SG", PlatformIds: []string{}},
			},
		},
	}

	query := sqlrepo.NewQuery()
	result := ApplySQLScope(query, data)
	assert.NotNil(t, result.FilterGroup)

	fg := result.FilterGroup
	assert.Len(t, fg.Filters, 1)
	assert.Equal(t, "tenant_id", fg.Filters[0].Field)
	assert.Len(t, fg.Groups, 1)
	assertDenyAllScope(t, fg.Groups[0])
}
