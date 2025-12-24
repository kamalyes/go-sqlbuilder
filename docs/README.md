# go-sqlbuilder 文档中心

欢迎使用 go-sqlbuilder - 一个强大、灵活的 Go 语言 SQL 构建器和仓储层框架

## 📚 快速开始

- [快速入门](./QUICKSTART.md) - 5分钟快速上手
- [模型定义](./MODELS.md) - 定义数据模型
- [CRUD 操作](./CRUD-OPERATIONS.md) - 基础 CRUD 操作

## 🔧 核心功能

### 数据操作

- [CRUD 操作](./CRUD-OPERATIONS.md) - 创建、读取、更新、删除
- [便捷查询方法](./CONVENIENCE-METHODS.md) - 简化的查询 API
- [工具方法](./UTILITY-METHODS.md) - Count、Exists、Pluck 等

### 查询构建

- [过滤条件](./FILTERS.md) - 等值、范围、模糊查询等
- [FilterGroup](./FILTERGROUP.md) - 复杂 AND/OR 条件组合
- [排序和分页](./SORTING-AND-PAGINATION.md) - ORDER BY、LIMIT、OFFSET
- [高级查询](./ADVANCED-QUERIES.md) - 子查询、JOIN、聚合等

### 特性功能

- [时间查询](./TIME-QUERIES.md) - 今日、本周、本月等便捷方法
- [字段选择](./FIELD-SELECTION.md) - SELECT、DISTINCT
- [自动字段填充](./AUTO-FIELD-SELECTION.md) - 创建时间、更新时间等
- [并发查询](./CONCURRENT-STATS.md) - 并行执行多个查询

## 📖 高级特性

- [EnhancedRepository](./ENHANCED-REPOSITORY.md) - 扩展功能
- [上下文使用](./CONTEXT-USAGE.md) - Context 最佳实践
- [错误处理](./ERROR-HANDLING.md) - 错误检查和处理
- [数据库迁移](./MIGRATOR.md) - Schema 迁移工具

## 文档结构

```bash
docs/
├── README.md                      # 本文档（文档中心）
├── QUICKSTART.md                  # 快速入门
├── MODELS.md                      # 模型定义
├── CRUD-OPERATIONS.md             # CRUD 操作
├── CONVENIENCE-METHODS.md         # 便捷查询方法
├── UTILITY-METHODS.md             # 工具方法
├── FILTERS.md                     # 过滤条件
├── FILTERGROUP.md                 # 条件组合
├── SORTING-AND-PAGINATION.md      # 排序和分页
├── TIME-QUERIES.md                # 时间查询
├── FIELD-SELECTION.md             # 字段选择
├── ADVANCED-QUERIES.md            # 高级查询
├── AUTO-FIELD-SELECTION.md        # 自动字段填充
├── ENHANCED-REPOSITORY.md         # 扩展功能
├── CONCURRENT-STATS.md            # 并发查询
├── CONTEXT-USAGE.md               # 上下文使用
├── ERROR-HANDLING.md              # 错误处理
└── MIGRATOR.md                    # 数据库迁移
```