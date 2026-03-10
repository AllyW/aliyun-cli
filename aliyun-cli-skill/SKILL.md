---
name: aliyun-cli
version: 1.0.0
description: Expert assistance for using the Aliyun CLI plugin system to manage cloud resources
trigger: auto
keywords: [aliyun, alibaba-cloud, 阿里云, ecs, fc, rds, oss, sls, ess, plugin, cloud]
---

# Aliyun CLI Plugin Expert

Help users effectively use the Aliyun CLI plugin system to manage Aliyun cloud resources.

## Description

This skill provides expert assistance for using the next-generation Aliyun CLI plugin system. It helps users:
- Understand the plugin-based architecture and its advantages
- Execute plugin commands with proper syntax
- Troubleshoot issues and optimize workflows
- Leverage advanced features like multi-version API support, structured parameters, and workflow orchestration

## Triggers

This skill should be invoked when the user:
- Asks about Aliyun CLI or cloud resource management
- Needs help with ECS, FC, RDS, OSS, SLS, or other Aliyun services
- Wants to perform operations on Aliyun resources
- Mentions "aliyun", "阿里云", plugin names like "ecs", "fc", "rds", etc.
- Encounters errors with Aliyun CLI commands

## Instructions

When helping users with Aliyun CLI:

1. **Always check if plugins are installed first**
   - Use `aliyun plugin list` to check installed plugins
   - Use `aliyun plugin install --names <plugin-name>` if needed

2. **Use the correct command syntax**
   - Modern syntax: `aliyun <product> <command> --param-name value`
   - Commands use kebab-case: `describe-instances`, `create-function`
   - Parameters use kebab-case: `--instance-id`, `--region-id`

3. **Leverage structured parameter input**
   - For array parameters: `--tag key=env value=prod --tag key=app value=web`
   - For object parameters: `--data-disk size=100 category=cloud_ssd`
   - Framework handles serialization automatically

4. **Use multi-version API support when needed**
   - Check available versions: `aliyun <product> list-api-versions`
   - Specify version: `--api-version 2022-02-22`
   - Set default via env: `export ALIBABA_CLOUD_ESS_API_VERSION=2022-02-22`

5. **Enable debugging when troubleshooting**
   - Use `--log-level=debug` for detailed logs
   - Use `--log-level=dev` for development mode
   - Use `--cli-dry-run` to test commands without executing (validates parameters and shows request details)

6. **Reference documentation**
   - Check ./references/plugin-advantages.md for feature comparisons
   - Check ./references/command-syntax.md for complete syntax guide
   - Check ./references/global-flags.md for all global flags (including --cli-dry-run)
   - Check ./references/common-scenarios.md for practical examples
   - Run scripts in ./scripts/examples/ directory for executable demos

## Key Concepts

### Plugin-Based Architecture
- Each Aliyun service has its own plugin (e.g., plugin-ecs, plugin-fc)
- Plugins are installed on-demand
- Independent versioning and updates

### Unified Parameter Style
- **Old CLI**: Mixed naming styles (PascalCase, camelCase, snake_case)
- **Plugin CLI**: Consistent kebab-case for all commands and parameters
- **Structured input**: Framework handles complex parameter serialization

### Multi-Version Support
- Single plugin can support multiple API versions
- Seamless switching between versions
- Complete help documentation for all versions

### Enhanced Developer Experience
- Detailed logging at multiple levels
- Better error messages with suggestions
- Type validation and auto-completion support

## Common Workflows

### Working with ECS Instances

```bash
# List instances
aliyun ecs describe-instances --region-id cn-hangzhou

# Create instance with structured parameters
aliyun ecs create-instance \
  --region-id cn-hangzhou \
  --instance-type ecs.g6.large \
  --image-id ubuntu_20_04_x64 \
  --system-disk category=cloud_essd size=40 \
  --data-disk category=cloud_essd size=100 \
  --tag key=env value=prod key=app value=web
```

### Working with Function Compute

```bash
# Create function with body parameters flattened
aliyun fc create-function \
  --function-name my-function \
  --runtime python3.9 \
  --handler index.handler \
  --code zipFile=@./code.zip \
  --memory-size 512 \
  --timeout 60
```

### Multi-Version ESS Operations

```bash
# Check available versions
aliyun ess list-api-versions

# Use specific version
aliyun ess describe-scaling-groups \
  --api-version 2022-02-22 \
  --region-id cn-hangzhou
```

## Response Format

When providing command suggestions:
1. Explain what the command does
2. Show the complete command with all required parameters
3. Explain any non-obvious parameter values
4. Mention relevant flags or options
5. If debugging is needed, include `--log-level=debug`

## Examples

### Example 1: User asks to list ECS instances

**User**: "Show me all my ECS instances in Hangzhou"

**Response**:
```bash
# First, ensure the ECS plugin is installed
aliyun plugin list | grep ecs
# If not installed: aliyun plugin install --names ecs

# List all ECS instances in Hangzhou region
aliyun ecs describe-instances --region-id cn-hangzhou

# For more detailed output with specific fields
aliyun ecs describe-instances \
  --region-id cn-hangzhou \
  --output cols=InstanceId,InstanceName,Status,PublicIpAddress
```

### Example 2: User encounters parameter issues

**User**: "I'm getting an error with Tag.1.Key syntax"

**Response**:
The new plugin CLI uses a simplified structured parameter syntax. Instead of:
```bash
# Old style (don't use this)
--Tag.1.Key=env --Tag.1.Value=prod
```

Use the modern structured syntax:
```bash
# New plugin style (use this)
--tag key=env value=prod
```

The framework automatically serializes this to the correct API format based on the API style (RPC/ROA).

## References

All documentation is included within this skill:
- ./references/plugin-advantages.md: Plugin CLI advantages and comparisons
- ./references/command-syntax.md: Complete command syntax guide
- ./references/global-flags.md: Complete global flags reference (--log-level, --cli-dry-run, --cli-query, etc.)
- ./references/common-scenarios.md: 10 practical usage scenarios
- ./scripts/examples/: Executable example scripts (ecs, fc, multi-version)
- ./QUICKSTART.md: 5-minute quick start guide

For more information about Aliyun CLI internals, refer to the official Aliyun CLI documentation and source code repository.
