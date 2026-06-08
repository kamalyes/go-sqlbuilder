/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-05-11 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-11 00:00:00
 * @FilePath: \go-sqlbuilder\scope\types_test.go
 * @Description: Scope 数据类型测试 - 覆盖 ScopeEntry、ScopeData 判断方法
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */
package scope

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ==============================================================================
// ScopeEntry
// ==============================================================================

// 场景：ScopeEntry.AllPlatformIds 各种输入
// 预期：空条目返回空，多地区平台ID自动去重，单条目正常返回
func TestScopeEntry_AllPlatformIds(t *testing.T) {
	t.Run("空条目", func(t *testing.T) {
		e := &ScopeEntry{}
		assert.Empty(t, e.AllPlatformIds())
	})

	t.Run("去重", func(t *testing.T) {
		e := &ScopeEntry{
			RegionPlatforms: []*RegionPlatformEntry{
				{RegionCode: "MM", PlatformIds: []string{"P1", "P2"}},
				{RegionCode: "TH", PlatformIds: []string{"P2", "P3"}},
			},
		}
		ids := e.AllPlatformIds()
		assert.Len(t, ids, 3)
		assert.Contains(t, ids, "P1")
		assert.Contains(t, ids, "P2")
		assert.Contains(t, ids, "P3")
	})

	t.Run("单条目", func(t *testing.T) {
		e := &ScopeEntry{
			RegionPlatforms: []*RegionPlatformEntry{
				{RegionCode: "MM", PlatformIds: []string{"P1"}},
			},
		}
		ids := e.AllPlatformIds()
		assert.Len(t, ids, 1)
		assert.Contains(t, ids, "P1")
	})
}

// ==============================================================================
// ScopeData 判断方法
// ==============================================================================

// 场景：Domain=2（默认OPS域值）
// 预期：IsOps()=true, IsTenant()=false
func TestScopeData_IsOps(t *testing.T) {
	data := NewScopeData()
	data.Domain = 2
	assert.True(t, data.IsOps())
	assert.False(t, data.IsTenant())
}

// 场景：Domain=1（默认租户域值）
// 预期：IsTenant()=true, IsOps()=false
func TestScopeData_IsTenant(t *testing.T) {
	data := NewScopeData()
	data.Domain = 1
	assert.True(t, data.IsTenant())
	assert.False(t, data.IsOps())
}

// 场景：自定义域配置（TenantDomain=5, OpsDomain=99），Domain=99
// 预期：IsOps()=true, IsTenant()=false
func TestScopeData_IsOps_WithCustomDomain(t *testing.T) {
	data := NewScopeData(WithDomainConfig(5, 99))
	data.Domain = 99
	assert.True(t, data.IsOps())
	assert.False(t, data.IsTenant())
}

// 场景：自定义域配置（TenantDomain=5, OpsDomain=99），Domain=5
// 预期：IsTenant()=true, IsOps()=false
func TestScopeData_IsTenant_WithCustomDomain(t *testing.T) {
	data := NewScopeData(WithDomainConfig(5, 99))
	data.Domain = 5
	assert.True(t, data.IsTenant())
	assert.False(t, data.IsOps())
}

// 场景：Domain 值既不匹配租户也不匹配 OPS（Domain=999）
// 预期：IsOps()=false, IsTenant()=false
func TestScopeData_IsOps_NeitherDomain(t *testing.T) {
	data := NewScopeData()
	data.Domain = 999
	assert.False(t, data.IsOps())
	assert.False(t, data.IsTenant())
}

// 场景：HasGlobalScope 判断是否存在全局作用域条目
// 预期：有ScopeType=1的条目返回true，无则返回false，空条目返回false
// 自定义ScopeTypeConfig时按自定义值判断
func TestScopeData_HasGlobalScope(t *testing.T) {
	t.Run("有全局作用域", func(t *testing.T) {
		data := NewScopeData()
		data.ScopeEntries = []*ScopeEntry{{ScopeType: 1}}
		assert.True(t, data.HasGlobalScope())
	})

	t.Run("无全局作用域", func(t *testing.T) {
		data := NewScopeData()
		data.ScopeEntries = []*ScopeEntry{{ScopeType: 2}}
		assert.False(t, data.HasGlobalScope())
	})

	t.Run("空ScopeEntries", func(t *testing.T) {
		data := NewScopeData()
		data.ScopeEntries = []*ScopeEntry{}
		assert.False(t, data.HasGlobalScope())
	})

	t.Run("TenantOwner", func(t *testing.T) {
		data := NewScopeData()
		data.Domain = 1
		data.IsOwner = true
		assert.True(t, data.HasGlobalScope())
	})

	t.Run("OpsOwnerWithoutGlobalScope", func(t *testing.T) {
		data := NewScopeData()
		data.Domain = 2
		data.IsOwner = true
		assert.False(t, data.HasGlobalScope())
	})

	t.Run("自定义ScopeTypeConfig", func(t *testing.T) {
		data := NewScopeData(WithScopeTypeConfig(10, 20, 30, 40))
		data.ScopeEntries = []*ScopeEntry{{ScopeType: 10}}
		assert.True(t, data.HasGlobalScope())
	})
}

// 场景：AllRegionCodes 从所有作用域条目中收集地区编码并去重
// 预期：跨地区级和平台级条目收集，自动去重，空条目返回空
func TestScopeData_AllRegionCodes(t *testing.T) {
	t.Run("正常去重", func(t *testing.T) {
		data := NewScopeData()
		data.ScopeEntries = []*ScopeEntry{
			{ScopeType: 2, RegionCodes: []string{"MM", "TH"}},
			{ScopeType: 3, RegionPlatforms: []*RegionPlatformEntry{
				{RegionCode: "SG", PlatformIds: []string{"P1"}},
				{RegionCode: "MM", PlatformIds: []string{"P2"}},
			}},
		}
		codes := data.AllRegionCodes()
		assert.Len(t, codes, 3)
		assert.Contains(t, codes, "MM")
		assert.Contains(t, codes, "TH")
		assert.Contains(t, codes, "SG")
	})

	t.Run("空ScopeEntries", func(t *testing.T) {
		data := NewScopeData()
		data.ScopeEntries = []*ScopeEntry{}
		assert.Empty(t, data.AllRegionCodes())
	})
}

// 场景：AllPlatformIds 从所有平台级条目中收集平台ID并去重
// 预期：跨多个RegionPlatformEntry收集，自动去重，空条目返回空
func TestScopeData_AllPlatformIds(t *testing.T) {
	t.Run("正常去重", func(t *testing.T) {
		data := NewScopeData()
		data.ScopeEntries = []*ScopeEntry{
			{ScopeType: 3, RegionPlatforms: []*RegionPlatformEntry{
				{RegionCode: "MM", PlatformIds: []string{"P1", "P2"}},
				{RegionCode: "TH", PlatformIds: []string{"P2", "P3"}},
			}},
		}
		ids := data.AllPlatformIds()
		assert.Len(t, ids, 3)
	})

	t.Run("空ScopeEntries", func(t *testing.T) {
		data := NewScopeData()
		data.ScopeEntries = []*ScopeEntry{}
		assert.Empty(t, data.AllPlatformIds())
	})
}

// 场景：IsXxxScope 使用默认 ScopeTypeConfig 判断作用域类型
// 预期：ScopeType=1为Global, =2为Region, =3为Platform, =4为Tenant
// ScopeType=99 不匹配任何类型
func TestScopeData_IsXxxScope(t *testing.T) {
	data := NewScopeData()
	assert.True(t, data.IsGlobalScope(&ScopeEntry{ScopeType: 1}))
	assert.True(t, data.IsRegionScope(&ScopeEntry{ScopeType: 2}))
	assert.True(t, data.IsPlatformScope(&ScopeEntry{ScopeType: 3}))
	assert.True(t, data.IsTenantScope(&ScopeEntry{ScopeType: 4}))
	assert.False(t, data.IsGlobalScope(&ScopeEntry{ScopeType: 99}))
}

// 场景：IsXxxScope 使用自定义 ScopeTypeConfig 判断作用域类型
// 预期：按自定义值（10/20/30/40）正确判断各类型
func TestScopeData_IsXxxScope_WithCustomScopeTypeConfig(t *testing.T) {
	data := NewScopeData(WithScopeTypeConfig(10, 20, 30, 40))
	assert.True(t, data.IsGlobalScope(&ScopeEntry{ScopeType: 10}))
	assert.True(t, data.IsRegionScope(&ScopeEntry{ScopeType: 20}))
	assert.True(t, data.IsPlatformScope(&ScopeEntry{ScopeType: 30}))
	assert.True(t, data.IsTenantScope(&ScopeEntry{ScopeType: 40}))
}
