# Quick Start Guide

快速上手 Aliyun CLI Plugin Expert Skill。

## 5 分钟快速开始

### 1. 在 Claude Code 中使用

直接在对话中提问，skill 会自动激活：

```
User: 帮我列出杭州的 ECS 实例

Claude: 让我帮您使用阿里云 CLI 的 ECS 插件...
[自动检查插件状态]
[提供正确命令]

aliyun ecs describe-instances --biz-region-id cn-hangzhou
```

### 2. 验证插件安装

```bash
cd aliyun-cli/aliyun-cli-skill
./scripts/check-plugin.sh ecs
```

输出示例：
```
✓ Plugin 'ecs' is installed
  Version: 1.0.0
```

### 3. 查看示例命令

```bash
# ECS 示例
./scripts/examples/ecs-examples.sh

# Function Compute 示例
./scripts/examples/fc-examples.sh

# 多版本 API 示例
./scripts/examples/multi-version.sh
```

### 4. 运行第一个命令

```bash
# 列出已安装的插件
aliyun plugin list

# 列出 ECS 实例（如果有 ECS 插件）
aliyun ecs describe-instances --biz-region-id cn-hangzhou

# 启用调试模式查看详细日志
aliyun ecs describe-instances \
  --biz-region-id cn-hangzhou \
  --log-level=debug
```

## 常见使用场景

### 场景 1: 我想列出某个地区的资源

**提问：**
```
User: 列出北京地区的所有 RDS 实例
```

**Skill 响应：**
```bash
# 确保安装 RDS 插件 (也可用 aliyun-cli-rds 或 RDS，不区分大小写)
aliyun plugin install --names rds

# 列出实例
aliyun rds describe-db-instances --biz-region-id cn-beijing
```

### 场景 2: 我想创建资源并打标签

**提问：**
```
User: 创建一个 ECS 实例，打上 env=prod 和 app=web 的标签
```

**Skill 响应：**
```bash
aliyun ecs create-instance \
  --biz-region-id cn-hangzhou \
  --instance-type ecs.g6.large \
  --image-id ubuntu_20_04_x64 \
  --security-group-id sg-xxx \
  --v-switch-id vsw-xxx \
  --tag key=env value=prod \
  --tag key=app value=web
```

### 场景 3: 我遇到了错误

**提问：**
```
User: 我运行 aliyun fc CreateFunction 报错了
```

**Skill 响应：**
1. 识别这是旧版语法
2. 提供新版正确语法
3. 添加调试标志
4. 解释参数差异

```bash
# 新版 CLI 使用 kebab-case
aliyun fc create-function \
  --function-name my-func \
  --runtime python3.9 \
  --handler index.handler \
  --code zipFile=@./code.zip \
  --log-level=debug
```

### 场景 4: 我需要使用新版 API

**提问：**
```
User: ESS 的 2022-02-22 版本怎么用？
```

**Skill 响应：**
```bash
# 查看可用版本
aliyun ess list-api-versions

# 使用特定版本
aliyun ess describe-scaling-groups \
  --api-version 2022-02-22 \
  --biz-region-id cn-hangzhou

# 或设置环境变量
export ALIBABA_CLOUD_ESS_API_VERSION=2022-02-22
aliyun ess describe-scaling-groups --biz-region-id cn-hangzhou
```

## 参考资料优先级

遇到问题时，按以下顺序查阅（所有文档都在此 skill 目录内）：

1. **Quick Reference** (最快)
   - `references/plugin-advantages.md` - 特性对比
   - `references/command-syntax.md` - 语法速查

2. **Examples** (实用)
   - `scripts/examples/ecs-examples.sh` - ECS 示例
   - `scripts/examples/fc-examples.sh` - FC 示例
   - `scripts/examples/multi-version.sh` - 多版本

3. **Scenarios** (场景化)
   - `references/common-scenarios.md` - 10 个常见场景

4. **Quick Start** (入门)
   - `QUICKSTART.md` (本文档) - 快速开始指南

## 命令速查表

| 任务 | 命令 |
|------|------|
| 检查插件 | `aliyun plugin list` |
| 安装插件 | `aliyun plugin install --names <name>` (支持: ecs 或 aliyun-cli-ecs，不区分大小写) |
| 更新插件 | `aliyun plugin update --names <name>` |
| 删除插件 | `aliyun plugin uninstall --names <name>` |
| 列出命令 | `aliyun <product> --help` |
| 命令帮助 | `aliyun <product> <command> --help` |
| 调试模式 | `--log-level=debug` |
| 查看版本 | `aliyun <product> list-api-versions` |
| JSON 输出 | `--output json` |
| 表格输出 | `--output table` |
| 自定义列 | `--output cols=Col1,Col2` |
| 过滤输出 | `--cli-query "<jmespath>"` |

## 语法对比速查

### 命令名

| 旧版 | 新版 |
|------|------|
| `DescribeInstances` | `describe-instances` |
| `CreateFunction` | `create-function` |
| `ListBuckets` | `list-buckets` |

### 参数名

| 旧版 | 新版 |
|------|------|
| `--InstanceId` | `--instance-id` |
| `--FunctionName` | `--function-name` |
| `--PageSize` | `--page-size` |

### 数组参数

| 旧版 | 新版 |
|------|------|
| `--Tag.1.Key=k1 --Tag.1.Value=v1` | `--tag key=k1 value=v1` |
| `--InstanceId.1=i-1 --InstanceId.2=i-2` | `--instance-id i-1 i-2` |

### Body 参数

| 旧版 | 新版 |
|------|------|
| `--body '{"functionName":"test"}'` | `--function-name test` |

## 下一步

1. **浏览示例**: 运行 `scripts/examples/` 下的脚本
2. **尝试命令**: 使用你的阿里云账号测试命令
3. **启用调试**: 添加 `--log-level=debug` 查看详细过程
4. **阅读文档**: 查看 `references/` 下的参考文档

## 获取帮助

- 使用 Claude Code 直接提问
- 查看 `references/common-scenarios.md` 的场景示例
- 运行命令时添加 `--help` 查看详细说明
- 添加 `--log-level=debug` 查看执行细节

## 技巧

**自动补全**: 支持 tab 补全（需配置）

**环境变量**: 设置常用参数
```bash
export ALIBABA_CLOUD_REGION_ID=cn-hangzhou
export ALIBABA_CLOUD_PROFILE=production
```

**输出过滤**: 使用 JMESPath 查询
```bash
aliyun ecs describe-instances \
  --cli-query "Instances.Instance[?Status=='Running'].InstanceId"
```

**组合使用**: 配合 jq 处理 JSON
```bash
aliyun ecs describe-instances --output json | \
  jq '.Instances.Instance[] | {id, name: .InstanceName}'
```

**脚本化**: 所有命令都可以写入脚本自动化
```bash
#!/bin/bash
REGION="cn-hangzhou"
INSTANCES=$(aliyun ecs describe-instances \
  --biz-region-id $REGION \
  --output json)
echo $INSTANCES | jq '.Instances.Instance | length'
```

---

**开始使用 Aliyun CLI Plugin Expert Skill，让云资源管理更简单！**
