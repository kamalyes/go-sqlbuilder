/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-05-11 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-11 13:00:00
 * @FilePath: \go-sqlbuilder\scope\config.go
 * @Description: Scope 配置 - 字段映射、域类型、作用域类型及选项函数
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */
package scope

// FieldMapping 数据库字段映射，将作用域逻辑字段映射到实际数据库列名
type FieldMapping struct {
	// TenantIDField 租户ID对应的数据库列名，如 "tenant_id"
	TenantIDField string
	// PlatformIDField 平台ID对应的数据库列名，如 "platform_id"
	PlatformIDField string
	// RegionCodeField 地区编码对应的数据库列名，如 "region_code"
	RegionCodeField string
}

// DomainConfig 域类型配置，定义哪些 int32 值代表 OPS 或 Tenant 域
type DomainConfig struct {
	// TenantDomainValue 租户域的 int32 值
	TenantDomainValue int32
	// OpsDomainValue OPS域的 int32 值
	OpsDomainValue int32
}

// ScopeTypeConfig 作用域类型配置，定义各作用域级别对应的 int32 值
type ScopeTypeConfig struct {
	// GlobalValue 全局作用域值（如 OPS 全局管理员、Tenant Owner）
	GlobalValue int32
	// RegionValue 地区级作用域值
	RegionValue int32
	// PlatformValue 平台级作用域值
	PlatformValue int32
	// TenantValue 租户级作用域值（OPS 管理指定租户）
	TenantValue int32
}

// Config 作用域完整配置，包含域类型、作用域类型和字段映射
type Config struct {
	Domain    DomainConfig
	ScopeType ScopeTypeConfig
	Mapping   FieldMapping
}

// DefaultConfig 返回默认配置（域值: Tenant=1, Ops=2; 作用域值: Global=1, Region=2, Platform=3, Tenant=4;
// 字段映射: tenant_id, platform_id, region_code）
func DefaultConfig() Config {
	return Config{
		Domain: DomainConfig{
			TenantDomainValue: 1,
			OpsDomainValue:    2,
		},
		ScopeType: ScopeTypeConfig{
			GlobalValue:   1,
			RegionValue:   2,
			PlatformValue: 3,
			TenantValue:   4,
		},
		Mapping: FieldMapping{
			TenantIDField:   "tenant_id",
			PlatformIDField: "platform_id",
			RegionCodeField: "region_code",
		},
	}
}

// Option 作用域配置选项函数
type Option func(*Config)

// WithDomainConfig 设置域类型值映射（tenantValue=租户域值, opsValue=OPS域值）
func WithDomainConfig(tenantValue, opsValue int32) Option {
	return func(c *Config) {
		c.Domain.TenantDomainValue = tenantValue
		c.Domain.OpsDomainValue = opsValue
	}
}

// WithScopeTypeConfig 设置作用域类型值映射
func WithScopeTypeConfig(global, region, platform, tenant int32) Option {
	return func(c *Config) {
		c.ScopeType.GlobalValue = global
		c.ScopeType.RegionValue = region
		c.ScopeType.PlatformValue = platform
		c.ScopeType.TenantValue = tenant
	}
}

// WithTenantIDField 设置租户ID对应的数据库列名
func WithTenantIDField(field string) Option {
	return func(c *Config) {
		c.Mapping.TenantIDField = field
	}
}

// WithPlatformIDField 设置平台ID对应的数据库列名
func WithPlatformIDField(field string) Option {
	return func(c *Config) {
		c.Mapping.PlatformIDField = field
	}
}

// WithRegionCodeField 设置地区编码对应的数据库列名
func WithRegionCodeField(field string) Option {
	return func(c *Config) {
		c.Mapping.RegionCodeField = field
	}
}

// WithFieldMapping 一次性设置完整的字段映射
func WithFieldMapping(mapping FieldMapping) Option {
	return func(c *Config) {
		c.Mapping = mapping
	}
}
