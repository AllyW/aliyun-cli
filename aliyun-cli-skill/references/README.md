# References Directory

快速参考文档目录。

## 文档列表

### 1. plugin-advantages.md
**插件 CLI 优势速查**

涵盖内容：
- 统一参数输入
- 一致的命名规范
- Body/Header 参数可见性
- 多版本 API 支持
- 增强的日志系统
- 模块化插件架构

适用场景：
- 想了解新旧 CLI 区别
- 需要快速对比特性
- 准备从旧版迁移

### 2. command-syntax.md
**命令语法完整指南**

涵盖内容：
- 基础命令结构
- 6 种参数类型详解
- 特殊参数处理（文件上传、查询过滤）
- 全局标志说明
- 多版本 API 使用
- 调试和错误处理

适用场景：
- 学习命令语法
- 查找特定参数用法
- 调试命令问题

### 3. global-flags.md
**全局参数完整参考**

涵盖内容：
- 所有可见全局参数（--log-level, --cli-dry-run, --cli-query等）
- 隐藏的高级参数（--pager, --waiter, --header等）
- 参数详细说明和使用场景
- 调试技巧和最佳实践

适用场景：
- 查找特定全局参数
- 了解调试选项（--cli-dry-run, --log-level）
- 学习高级功能（--pager, --waiter）

### 4. common-scenarios.md
**10 个常见场景**

涵盖场景：
1. 实例管理（列出、过滤、批量操作）
2. 资源创建与标签
3. 监控与诊断
4. 函数计算部署
5. RDS 数据库管理
6. OSS 存储操作
7. SLS 日志服务
8. 多版本 API 使用
9. 资源标签策略
10. 调试与故障排查

适用场景：
- 需要完整的端到端示例
- 学习最佳实践
- 复制粘贴快速开始

## 使用建议

### 快速查找
```bash
# 搜索特定关键词
grep -r "create-function" references/
grep -r "multi-version" references/
```

### 按需阅读
- **新手**: 从 plugin-advantages.md 开始
- **语法查询**: 直接看 command-syntax.md
- **全局参数**: 查看 global-flags.md
- **实战示例**: 参考 common-scenarios.md

### 配合使用
1. 先看 plugin-advantages.md 了解特性
2. 再看 command-syntax.md 学习语法
3. 查看 global-flags.md 了解全局参数（特别是 --cli-dry-run）
4. 最后看 common-scenarios.md 实践场景

## 文档更新

这些文档基于：
- Aliyun CLI v3.x 插件系统
- Plugin Runtime 最新版本
- 真实生产环境最佳实践

文档保持与 Aliyun CLI 官方实现同步更新。

## 反馈

发现问题或有改进建议？
1. 在 Claude Code 中直接反馈
2. 通过 Aliyun CLI 官方渠道反馈
3. 联系 skill 维护者

---

**快速跳转：**
- [插件优势](./plugin-advantages.md)
- [命令语法](./command-syntax.md)
- [全局参数](./global-flags.md) ⭐ 新增
- [常见场景](./common-scenarios.md)
- [返回 Skill 主页](../README.md)
