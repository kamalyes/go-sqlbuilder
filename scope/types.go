/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-05-11 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-11 00:00:00
 * @FilePath: \go-sqlbuilder\scope\types.go
 * @Description: Scope 数据类型 - 作用域条目、作用域数据及业务判断方法
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */
package scope

// RegionPlatformEntry 地区-平台绑定条目，表示某地区下可访问的平台列表
type RegionPlatformEntry struct {
	// RegionCode 地区编码，如 "MM"、"TH"
	RegionCode string
	// PlatformIds 该地区下可访问的平台ID列表
	PlatformIds []string
}

// ScopeEntry 作用域条目，描述一个作用域级别的访问范围
type ScopeEntry struct {
	// ScopeType 作用域类型值，对应 ScopeTypeConfig 中的配置
	ScopeType int32
	// RegionCodes 地区编码列表（地区级作用域使用）
	RegionCodes []string
	// RegionPlatforms 地区-平台绑定列表（平台级作用域使用）
	RegionPlatforms []*RegionPlatformEntry
	// TenantIds 租户ID列表（OPS 租户级作用域使用）
	TenantIds []string
}

// AllPlatformIds 收集该条目下所有不重复的平台ID
func (e *ScopeEntry) AllPlatformIds() []string {
	seen := make(map[string]struct{})
	var result []string
	for _, rp := range e.RegionPlatforms {
		for _, pid := range rp.PlatformIds {
			if _, ok := seen[pid]; !ok {
				seen[pid] = struct{}{}
				result = append(result, pid)
			}
		}
	}
	return result
}

// ScopeData 作用域数据，包含当前用户的域、租户信息和作用域条目
type ScopeData struct {
	// Domain 当前域值（对应 DomainConfig 中的配置）
	Domain int32
	// TenantID 当前租户ID
	TenantID string
	// IsOwner 是否为租户 Owner（可查看租户下所有数据）
	IsOwner bool
	// ScopeEntries 作用域条目列表
	ScopeEntries []*ScopeEntry
	// Config 作用域配置
	Config Config
}

// NewScopeData 使用选项创建 ScopeData 实例
func NewScopeData(opts ...Option) ScopeData {
	cfg := DefaultConfig()
	for _, o := range opts {
		o(&cfg)
	}
	return ScopeData{Config: cfg}
}

// IsOps 判断当前域是否为 OPS 域
func (d ScopeData) IsOps() bool {
	return d.Domain == d.Config.Domain.OpsDomainValue
}

// IsTenant 判断当前域是否为租户域
func (d ScopeData) IsTenant() bool {
	return d.Domain == d.Config.Domain.TenantDomainValue
}

// HasGlobalScope 判断是否拥有全局作用域（OPS 全局管理员或 Tenant Owner/全局用户）
func (d ScopeData) HasGlobalScope() bool {
	if d.IsOwner && d.IsTenant() {
		return true
	}
	for _, e := range d.ScopeEntries {
		if e == nil {
			continue
		}
		if e.ScopeType == d.Config.ScopeType.GlobalValue {
			return true
		}
	}
	return false
}

// AllRegionCodes 收集所有作用域条目中不重复的地区编码
func (d ScopeData) AllRegionCodes() []string {
	seen := make(map[string]struct{})
	var result []string
	for _, e := range d.ScopeEntries {
		if e == nil {
			continue
		}
		for _, rc := range e.RegionCodes {
			if _, ok := seen[rc]; !ok {
				seen[rc] = struct{}{}
				result = append(result, rc)
			}
		}
		for _, rp := range e.RegionPlatforms {
			if _, ok := seen[rp.RegionCode]; !ok {
				seen[rp.RegionCode] = struct{}{}
				result = append(result, rp.RegionCode)
			}
		}
	}
	return result
}

// AllPlatformIds 收集所有作用域条目中不重复的平台ID
func (d ScopeData) AllPlatformIds() []string {
	seen := make(map[string]struct{})
	var result []string
	for _, e := range d.ScopeEntries {
		if e == nil {
			continue
		}
		for _, pid := range e.AllPlatformIds() {
			if _, ok := seen[pid]; !ok {
				seen[pid] = struct{}{}
				result = append(result, pid)
			}
		}
	}
	return result
}

// IsGlobalScope 判断指定条目是否为全局作用域
func (d ScopeData) IsGlobalScope(entry *ScopeEntry) bool {
	if entry == nil {
		return false
	}
	return entry.ScopeType == d.Config.ScopeType.GlobalValue
}

// IsRegionScope 判断指定条目是否为地区级作用域
func (d ScopeData) IsRegionScope(entry *ScopeEntry) bool {
	if entry == nil {
		return false
	}
	return entry.ScopeType == d.Config.ScopeType.RegionValue
}

// IsPlatformScope 判断指定条目是否为平台级作用域
func (d ScopeData) IsPlatformScope(entry *ScopeEntry) bool {
	if entry == nil {
		return false
	}
	return entry.ScopeType == d.Config.ScopeType.PlatformValue
}

// IsTenantScope 判断指定条目是否为租户级作用域
func (d ScopeData) IsTenantScope(entry *ScopeEntry) bool {
	if entry == nil {
		return false
	}
	return entry.ScopeType == d.Config.ScopeType.TenantValue
}
