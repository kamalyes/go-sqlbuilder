/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-05-11 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-11 02:32:26
 * @FilePath: \go-sqlbuilder\scope\config_test.go
 * @Description: Scope 配置测试 - 覆盖 Config、Option 及默认值
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */
package scope

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// 场景：验证 DefaultConfig 返回的默认配置值
// 预期：TenantDomainValue=1, OpsDomainValue=2, GlobalValue=1, RegionValue=2, PlatformValue=3, TenantValue=4
// 字段映射默认值：TenantIDField="tenant_id", PlatformIDField="platform_id", RegionCodeField="region_code"
func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	assert.Equal(t, int32(1), cfg.Domain.TenantDomainValue)
	assert.Equal(t, int32(2), cfg.Domain.OpsDomainValue)
	assert.Equal(t, int32(1), cfg.ScopeType.GlobalValue)
	assert.Equal(t, int32(2), cfg.ScopeType.RegionValue)
	assert.Equal(t, int32(3), cfg.ScopeType.PlatformValue)
	assert.Equal(t, int32(4), cfg.ScopeType.TenantValue)
	assert.Equal(t, "tenant_id", cfg.Mapping.TenantIDField)
	assert.Equal(t, "platform_id", cfg.Mapping.PlatformIDField)
	assert.Equal(t, "region_code", cfg.Mapping.RegionCodeField)
}

// 场景：WithDomainConfig 自定义域配置
// 预期：TenantDomainValue=10, OpsDomainValue=20
func TestWithDomainConfig(t *testing.T) {
	cfg := DefaultConfig()
	WithDomainConfig(10, 20)(&cfg)
	assert.Equal(t, int32(10), cfg.Domain.TenantDomainValue)
	assert.Equal(t, int32(20), cfg.Domain.OpsDomainValue)
}

// 场景：WithScopeTypeConfig 自定义作用域类型配置
// 预期：GlobalValue=100, RegionValue=200, PlatformValue=300, TenantValue=400
func TestWithScopeTypeConfig(t *testing.T) {
	cfg := DefaultConfig()
	WithScopeTypeConfig(100, 200, 300, 400)(&cfg)
	assert.Equal(t, int32(100), cfg.ScopeType.GlobalValue)
	assert.Equal(t, int32(200), cfg.ScopeType.RegionValue)
	assert.Equal(t, int32(300), cfg.ScopeType.PlatformValue)
	assert.Equal(t, int32(400), cfg.ScopeType.TenantValue)
}

// 场景：WithTenantIDField 设置租户ID字段名
// 预期：TenantIDField = "tid"
func TestWithTenantIDField(t *testing.T) {
	cfg := DefaultConfig()
	WithTenantIDField("tid")(&cfg)
	assert.Equal(t, "tid", cfg.Mapping.TenantIDField)
}

// 场景：WithPlatformIDField 设置平台ID字段名
// 预期：PlatformIDField = "pid"
func TestWithPlatformIDField(t *testing.T) {
	cfg := DefaultConfig()
	WithPlatformIDField("pid")(&cfg)
	assert.Equal(t, "pid", cfg.Mapping.PlatformIDField)
}

// 场景：WithRegionCodeField 设置地区编码字段名
// 预期：RegionCodeField = "region"
func TestWithRegionCodeField(t *testing.T) {
	cfg := DefaultConfig()
	WithRegionCodeField("region")(&cfg)
	assert.Equal(t, "region", cfg.Mapping.RegionCodeField)
}

// 场景：WithFieldMapping 一次性设置所有字段映射
// 预期：三个字段映射均被正确设置
func TestWithFieldMapping(t *testing.T) {
	cfg := DefaultConfig()
	mapping := FieldMapping{TenantIDField: "t", PlatformIDField: "p", RegionCodeField: "r"}
	WithFieldMapping(mapping)(&cfg)
	assert.Equal(t, mapping, cfg.Mapping)
}

// 场景：NewScopeData 同时使用多个 Option
// 预期：所有 Option 均生效，Domain 和 Mapping 字段正确
func TestNewScopeData_MultipleOptions(t *testing.T) {
	data := NewScopeData(
		WithDomainConfig(10, 20),
		WithTenantIDField("tid"),
		WithRegionCodeField("rc"),
	)
	assert.Equal(t, int32(10), data.Config.Domain.TenantDomainValue)
	assert.Equal(t, int32(20), data.Config.Domain.OpsDomainValue)
	assert.Equal(t, "tid", data.Config.Mapping.TenantIDField)
	assert.Equal(t, "rc", data.Config.Mapping.RegionCodeField)
}

// 场景：FieldMapping.hasAnyField 判断是否有任何字段映射
// 预期：全空返回false，任一字段有值返回true
func TestFieldMapping_HasAnyField(t *testing.T) {
	t.Run("全部为空", func(t *testing.T) {
		m := FieldMapping{}
		assert.False(t, m.hasAnyField())
	})

	t.Run("仅有PlatformIDField", func(t *testing.T) {
		m := FieldMapping{PlatformIDField: "pid"}
		assert.True(t, m.hasAnyField())
	})

	t.Run("仅有RegionCodeField", func(t *testing.T) {
		m := FieldMapping{RegionCodeField: "rc"}
		assert.True(t, m.hasAnyField())
	})

	t.Run("都有值", func(t *testing.T) {
		m := FieldMapping{PlatformIDField: "pid", RegionCodeField: "rc"}
		assert.True(t, m.hasAnyField())
	})
}

// 场景：scopeFieldSet 根据字段映射构建需要移除的字段集合
// 预期：全空映射返回空集合，全有值映射返回3个元素的集合
func TestScopeFieldSet(t *testing.T) {
	t.Run("全部为空", func(t *testing.T) {
		set := scopeFieldSet(FieldMapping{})
		assert.Empty(t, set)
	})

	t.Run("全部有值", func(t *testing.T) {
		set := scopeFieldSet(FieldMapping{
			TenantIDField:   "tid",
			PlatformIDField: "pid",
			RegionCodeField: "rc",
		})
		assert.Len(t, set, 3)
		assert.Contains(t, set, "tid")
		assert.Contains(t, set, "pid")
		assert.Contains(t, set, "rc")
	})
}
