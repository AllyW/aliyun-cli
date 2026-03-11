---
name: aliyun-cli
version: 1.0.0
description: Expert assistance for using Aliyun CLI to manage cloud resources
trigger: auto
keywords: [aliyun, alibaba-cloud, 阿里云, ecs, fc, rds, oss, sls, ess, plugin, cloud, cli]
---

# Aliyun CLI Expert

Help users effectively use Aliyun CLI to manage Aliyun cloud resources.

## Description

This skill provides expert assistance for using Aliyun CLI (Command Line Interface). It helps users:
- Install and configure Aliyun CLI with proper authentication
- Understand the modern plugin-based architecture and its advantages
- Execute commands with correct syntax (using plugins when available)
- Troubleshoot issues and optimize workflows
- Leverage advanced features like multi-version API support, structured parameters, and workflow orchestration
- Handle special cases like OSS custom commands

## Getting Started

### Installation

**macOS (Homebrew)**
```bash
brew install aliyun-cli
```

**Linux (Binary)**
```bash
# Download and install
wget https://aliyuncli.alicdn.com/aliyun-cli-linux-latest-amd64.tgz
tar -xzf aliyun-cli-linux-latest-amd64.tgz
sudo mv aliyun /usr/local/bin/

# Verify installation
aliyun version
```

**Windows (Binary)**
```powershell
# Download from: https://aliyuncli.alicdn.com/aliyun-cli-windows-latest-amd64.zip
# Extract and add to PATH
```

### Configuration

**Interactive Configuration (Recommended)**
```bash
aliyun configure

# You'll be prompted for:
# - Access Key ID
# - Access Key Secret
# - Region (e.g., cn-hangzhou)
# - Output format (json, table, cols)
# - Language (en, zh)
```

**Configuration Modes**

1. **AK Mode (Access Key)**
```bash
aliyun configure --mode AK
# Enter Access Key ID and Secret
# Most common for personal accounts
```

2. **StsToken Mode (Temporary credentials)**
```bash
aliyun configure --mode StsToken
# Enter Access Key ID, Secret, and STS Token
# For temporary access or RAM role assumption
```

3. **RamRoleArn Mode (Assume RAM role)**
```bash
aliyun configure --mode RamRoleArn
# Enter Access Key ID, Secret, RAM Role ARN, and Role Session Name
# For cross-account or elevated privileges
```

4. **EcsRamRole Mode (Instance RAM role)**
```bash
aliyun configure --mode EcsRamRole
# Enter RAM Role Name attached to ECS instance
# No credentials needed - uses instance metadata
```

5. **RsaKeyPair Mode (RSA key pair)**
```bash
aliyun configure --mode RsaKeyPair
# Enter Key Pair Name and Private Key file path
# For RSA key-based authentication
```

6. **RamRoleArnWithEcs Mode (ECS instance + RAM role)**
```bash
aliyun configure --mode RamRoleArnWithEcs
# Combines ECS instance role with RAM role assumption
```

**Environment Variables**

You can also configure credentials via environment variables (highest priority):

```bash
# Access Key Mode
export ALIBABA_CLOUD_ACCESS_KEY_ID=your_access_key_id
export ALIBABA_CLOUD_ACCESS_KEY_SECRET=your_access_key_secret
export ALIBABA_CLOUD_REGION_ID=cn-hangzhou

# STS Token Mode (temporary credentials)
export ALIBABA_CLOUD_ACCESS_KEY_ID=your_access_key_id
export ALIBABA_CLOUD_ACCESS_KEY_SECRET=your_access_key_secret
export ALIBABA_CLOUD_SECURITY_TOKEN=your_sts_token

# ECS RAM Role Mode
export ALIBABA_CLOUD_ECS_METADATA=role_name
```

**Credential Priority**

Credentials are loaded in the following order (highest to lowest):
1. Environment variables
2. Configuration file (`~/.aliyun/config.json`)
3. ECS instance RAM role (if running on ECS)

**Managing Profiles**

```bash
# Create named profile
aliyun configure --profile project1 --mode AK --access-key-id xxxx --access-key-secret xxxx  --region cn-hangzhou

# Use specific profile
aliyun ecs describe-instances --profile project1

# Set default profile
export ALIBABA_CLOUD_PROFILE=project1

# List all profiles
aliyun configure list
```

**Quick Test**

```bash
# Test authentication
aliyun ecs describe-regions

# If successful, you'll see a list of available regions
```

## Triggers

This skill should be invoked when the user:
- Asks about Aliyun CLI or cloud resource management
- Needs help with ECS, FC, RDS, OSS, SLS, or other Aliyun services
- Wants to perform operations on Aliyun resources
- Mentions "aliyun", "阿里云", plugin names like "ecs", "fc", "rds", etc.
- Encounters errors with Aliyun CLI commands

## Instructions

When helping users with Aliyun CLI:

0. **Ensure CLI is installed and configured**
   - If user hasn't installed CLI, guide them through installation (see Getting Started)
   - If authentication fails, help configure credentials: `aliyun configure`
   - Test with: `aliyun ecs describe-regions`

1. **Always use --help to discover commands and parameters**
   - **Critical**: Do NOT guess command or parameter names
   - **Step 1**: Check available subcommands: `aliyun <product-code> --help`
   - **Step 2**: Check command parameters: `aliyun <product-code> <subcommand> --help`
   - Help output is the **authoritative source** for both built-in and plugin commands
   - Plugin help includes full structure information

1a. **Check if required service plugins are installed**
   - Use `aliyun plugin list` to check installed plugins
   - Use `aliyun plugin install --names <plugin-name>` to install
   - **Plugin name formats** (both work, case-insensitive):
     - Full name: `aliyun-cli-ecs` or `aliyun-cli-ECS`
     - Product code: `ecs` or `ECS` or `Ecs`
     - Examples:
       ```bash
       aliyun plugin install --names ecs           # Recommended (short)
       aliyun plugin install --names aliyun-cli-ecs  # Full name also works
       aliyun plugin install --names ECS VPC RDS   # Multiple plugins, any case
       ```
   - Other plugin operations also support both formats:
     ```bash
     aliyun plugin update --name ecs
     aliyun plugin uninstall --name aliyun-cli-ecs
     ```
   - Most services use plugins, but some (like OSS) have custom implementations

2. **Use the correct command style (but always verify with --help)**
   - **Prefer CLI Native style (plugin)**: `aliyun <product> <command> --param-name value`
     - Commands: kebab-case (`describe-instances`, `create-function`)
     - Parameters: kebab-case (`--instance-id`, `--region-id`)
     - **Consistent naming**: Unified style across all commands
   - **OpenAPI style (built-in)**: Only when plugin not available
     - Commands: PascalCase (`DescribeInstances`, `CreateFunction`)
     - Parameters: **Mixed/Inconsistent** (varies by API - must check `--help`)
   - **Command detection**: All-lowercase commands → plugin; contains uppercase → built-in
   - **Always verify**: Use `--help` to get exact command and parameter names
   - **Exception - OSS**: Uses custom commands, NOT API patterns. Always check `aliyun oss --help` and `aliyun ossutil --help`

3. **Leverage structured parameter input**
   - For array parameters: `--tag key=env value=prod --tag key=app value=web`
   - For object parameters: `--data-disk size=100 category=cloud_ssd`
   - Framework handles serialization automatically

3a. **Understand global parameter naming in plugins**
   - CLI has global parameters like `--region`, `--region-id`, `--profile`, etc.
   - When product API has conflicting parameter names, plugin renames them to avoid conflicts:
     - **First choice**: Add `biz-` prefix → `--biz-region`, `--biz-region-id`
     - **If still conflicts**: Add product-code prefix → `--<product-code>-region-id` (e.g., `--ecs-region-id`)
   - **Always check `--help`** to see the actual parameter name used in the plugin
   - Example:
     ```bash
     # Original API parameter: regionId
     # CLI global parameter: --region-id
     # Plugin renames to: --biz-region-id (or --ecs-region-id if biz-region-id also conflicts)

     aliyun ecs some-command --biz-region-id cn-hangzhou  # Not --region-id
     ```
   - This applies to all CLI global parameters, not just region

4. **Use multi-version API support when needed**
   - Check available versions: `aliyun <product> list-api-versions`
   - Specify version: `--api-version 2022-02-22`
   - Set default via env: `export ALIBABA_CLOUD_ESS_API_VERSION=2022-02-22`

5. **Enable debugging when troubleshooting**
   - Use `--log-level debug` for detailed logs
   - Use `--log-level dev` for development mode
   - Use `--cli-dry-run` to test commands without executing (validates parameters and shows request details)

6. **Leverage the help system**
   - Use `aliyun <product> --help` to see available commands
   - Plugin help (if installed): Shows CLI Native style with full parameter info
   - Built-in help: Set `ALIBABA_CLOUD_ORIGINAL_PRODUCT_HELP=true` to view
   - Use `aliyun <product> <command> --help` for detailed command help
   - Plugin commands show structure field for complex parameters

7. **Reference documentation**
   - Check ./references/installation-guide.md for installation and authentication setup
   - Check ./references/plugin-advantages.md for feature comparisons
   - Check ./references/command-syntax.md for complete syntax guide
   - Check ./references/global-flags.md for all global flags (including --cli-dry-run)
   - Check ./references/common-scenarios.md for practical examples
   - Run scripts in ./scripts/examples/ directory for executable demos

## Key Concepts

### Modern CLI Architecture
- **Plugin-Based**: Each Aliyun service has its own plugin (e.g., plugin-ecs, plugin-fc)
- **On-Demand Installation**: Install only the plugins you need
- **Independent Versioning**: Each plugin can be updated independently
- **Special Implementations**: Some services like OSS use custom implementations

### Unified Command Style
- **Consistent Naming**: kebab-case for all commands and parameters
- **Structured Input**: Framework handles complex parameter serialization
- **Easy to Remember**: Predictable command patterns

### Multi-Version Support
- Single plugin can support multiple API versions
- Seamless switching between versions
- Complete help documentation for all versions

### Enhanced Developer Experience
- Detailed logging at multiple levels
- Better error messages with suggestions
- Type validation and auto-completion support

## Command Styles and Help System

### Two Command Styles in Aliyun CLI

Aliyun CLI supports **two command styles** that coexist:

#### 1. OpenAPI Style (Legacy/Built-in)

**Characteristics:**
- **Subcommand naming**: PascalCase (e.g., `CreateInstance`, `DescribeInstances`)
- **Parameter naming**: **Mixed/Inconsistent** (could be PascalCase, camelCase, snake_case, etc.)
- **Direct API mapping**: Maps 1:1 to OpenAPI specifications
- **ROA APIs**: Create commands often only expose `--body json` parameter with minimal internal information

**Example:**
```bash
# OpenAPI style command - note inconsistent parameter naming
aliyun Ecs CreateInstance \
  --InstanceName my-instance \
  --ImageId ubuntu_20_04 \
  --InstanceType ecs.g6.large
```

**Limitations:**
- Less discoverable (need to know exact API names)
- **Inconsistent parameter naming** across different APIs
- ROA API body parameters not expanded
- Requires external API documentation for complex structures

#### 2. CLI Native Style (Plugin-based, Modern)

**Characteristics:**
- **Subcommand naming**: kebab-case (e.g., `create-instance`, `describe-instances`)
- **Parameter naming**: kebab-case (e.g., `--instance-name`, `--image-id`)
- **Consistent naming**: **Unified, predictable naming style** across all commands and parameters
- **CLI-friendly design**: Intuitive, follows Unix CLI conventions
- **Fully expanded parameters**: ROA API body structures fully expanded in help
- **Self-contained**: All parameter information available via `--help`, no need for external docs
- **Structure field**: Each parameter shows its internal structure in help output
- **Format hints**: Help output shows input format for different parameter types

**Example:**
```bash
# CLI Native style command (plugin)
aliyun ecs create-instance \
  --instance-name my-instance \
  --image-id ubuntu_20_04 \
  --instance-type ecs.g6.large
```

**Parameter Format Guidelines (from help output):**

Plugin help displays format hints for different parameter types:

1. **Basic types (string, int, bool)**: Direct value assignment
   ```bash
   --instance-name my-instance    # string
   --port 8080                     # int
   --dry-run true                  # bool
   ```

2. **Simple lists**: Space-separated values (format hint: `--param value1 value2 value3`)
   ```bash
   # Format shown in help: --additional-attributes value1 value2 value3
   aliyun ecs create-instance \
     --security-group-ids sg-001 sg-002 sg-003
   ```

3. **Key-value pairs**: Use `Key=value` syntax (format hint: `--param Key=a Value=b`)
   ```bash
   # Format shown in help: --tag Key=a Value=b
   aliyun ecs create-instance \
     --tag Key=env Value=prod \
     --tag Key=app Value=web
   ```

4. **Complex structures**: Use single quotes to wrap structured data
   ```bash
   # Help shows the structure field for this parameter
   # Example: structure: {DiskName: string, Size: int, Category: string}
   aliyun ecs create-instance \
     --data-disk '{"DiskName": "data1", "Size": 100, "Category": "cloud_ssd"}' \
     --data-disk '{"DiskName": "data2", "Size": 200, "Category": "cloud_efficiency"}'
   ```

5. **Global parameter conflict handling**: Renamed to avoid CLI global parameters
   ```bash
   # CLI has global parameters: --region, --region-id, --profile, etc.
   # If API has same parameter name, plugin renames it:

   # Priority 1: Add biz- prefix
   aliyun ecs some-command --biz-region-id cn-hangzhou  # Not --region-id

   # Priority 2: If biz- conflicts, use product-code prefix
   aliyun ecs another-command --ecs-region-id cn-hangzhou  # If --biz-region-id also exists

   # Always check help to see actual parameter name
   aliyun ecs some-command --help | grep region
   ```

**Always check `--help` output** - it provides the exact format and actual parameter names (including any renamed ones).

**Advantages:**
- **Discoverable**: Use `aliyun ecs --help` to see all commands
- **Self-documenting**: Each parameter has detailed help with structure information
- **Body expansion**: ROA API bodies fully expanded into individual parameters
- **No external docs needed**: Complete information in CLI help
- **Type information**: Shows data types and constraints

### Command Style Detection

**How Aliyun CLI determines which system to use:**

```bash
# All-lowercase subcommand → Plugin system (CLI Native)
aliyun ecs describe-instances

# Mixed/PascalCase subcommand → Built-in system (OpenAPI)
aliyun Ecs DescribeInstances
```

**Rule**: If the subcommand is **all lowercase** (including hyphens), CLI routes to the plugin. Otherwise, it uses the built-in OpenAPI system.

### Help System Intelligence

#### Product-level Help

**Case 1: No plugin installed**
```bash
$ aliyun ecs --help
# Shows built-in product help
# Displays message: "Plugin available: aliyun plugin install --names ecs"
```

**Case 2: Plugin installed**
```bash
$ aliyun ecs --help
# Automatically shows plugin help (CLI Native commands)
# Lists kebab-case commands: create-instance, describe-instances, etc.
```

**Case 3: View original built-in help**
```bash
# Override to see built-in OpenAPI help
$ ALIBABA_CLOUD_ORIGINAL_PRODUCT_HELP=true aliyun ecs --help
# Shows PascalCase commands: CreateInstance, DescribeInstances, etc.
```

#### Command-level Help

**Plugin command (CLI Native):**
```bash
$ aliyun ecs create-instance --help

# Shows:
# - All parameters with kebab-case names
# - Type information for each parameter
# - Structure field showing nested object layout
# - Default values and constraints
# - Examples

# Example parameter help:
# --data-disk
#   Type: list
#   Structure: category=string,size=integer,encrypted=boolean
#   Description: Data disks to attach to the instance
```

**Built-in command (OpenAPI):**
```bash
$ aliyun Ecs CreateInstance --help

# Shows:
# - Parameters with PascalCase names
# - Basic type information
# - For ROA APIs: Often just --body parameter
# - May need external API docs for complex structures
```

### ROA API Body Expansion

**Key Difference for ROA-style APIs:**

#### OpenAPI Style (Built-in)
```bash
# Create command with ROA body - minimal help
aliyun FC CreateFunction --body '{
  "functionName": "my-func",
  "runtime": "python3.9",
  "handler": "index.handler",
  "memorySize": 512
}'

# Help shows:
# --body string   Request body in JSON format
# (No information about what fields are available!)
```

#### CLI Native Style (Plugin)
```bash
# Same operation - fully expanded parameters
aliyun fc create-function \
  --function-name my-func \
  --runtime python3.9 \
  --handler index.handler \
  --memory-size 512

# Help shows each parameter:
# --function-name string     Function name
# --runtime string           Runtime (python3.9, nodejs14, etc.)
# --handler string           Handler in format file.method
# --memory-size integer      Memory in MB (128-32768)
#   Structure: None (scalar)
```

**Plugin Advantage**: All body fields are expanded as individual CLI parameters with full documentation.

### Style Combination Rules

**IMPORTANT**: Understand the naming conventions!

**CLI Native (Plugin) - Consistent:**
```bash
# Plugin: all-lowercase subcommand + kebab-case params (always consistent)
# product-code case doesn't matter (ecs, Ecs, ECS all work the same)
# Note: --biz-region-id (not --region-id) due to global parameter conflict
aliyun ecs describe-instances --instance-id i-xxx --biz-region-id cn-hangzhou
aliyun ECS describe-instances --instance-id i-xxx --biz-region-id cn-hangzhou  # Same as above
```

**Built-in (OpenAPI) - Inconsistent:**
```bash
# Built-in: Subcommand with uppercase + MIXED parameter naming
# product-code case doesn't matter, subcommand case determines the system
aliyun ecs DescribeInstances --InstanceId i-xxx --RegionId cn-hangzhou
aliyun ECS DescribeInstances --InstanceId i-xxx --RegionId cn-hangzhou  # Same as above
# Note: Parameters might be PascalCase, camelCase, or other styles
# Varies by API - not consistent across services
```

**WRONG - Don't mix command styles:**
```bash
# DON'T mix plugin subcommands with built-in parameter styles
aliyun ecs describe-instances --InstanceId i-xxx    # Wrong! kebab-case subcommand should use kebab-case params
aliyun ecs DescribeInstances --instance-id i-xxx    # Wrong! PascalCase subcommand should use its API params
```

**Key Point - Command Detection Logic**:
- Command structure: `aliyun <product-code> <subcommand> --parameters`
- **product-code** is always case-insensitive (ecs, Ecs, ECS are all the same)
- **subcommand** determines which system is used:
  - **All lowercase subcommand** (e.g., `describe-instances`) → Uses plugin → kebab-case parameters
  - **Contains uppercase** (e.g., `DescribeInstances`) → Uses built-in → API-specific parameters (inconsistent naming)

### When Constructing Commands

**For assistants/agents generating Aliyun CLI commands:**

**Golden Rule: Always use `--help` to discover commands and parameters**

1. **Step 1: Discover available subcommands**
   ```bash
   aliyun <product-code> --help
   ```
   - Shows all available subcommands for the product
   - Indicates if it's plugin-based or built-in
   - Lists command naming style

2. **Step 2: Get detailed parameter information**
   ```bash
   aliyun <product-code> <subcommand> --help
   ```
   - Shows ALL parameters with exact names
   - For plugins: includes type, structure, constraints
   - For built-in: shows parameter names (which may be inconsistent in format)

3. **Step 3: Construct command based on help output**
   - Use exact command names from help
   - Use exact parameter names from help
   - Do NOT guess or infer parameter names

4. **Why this workflow is critical:**
   - Built-in commands have inconsistent parameter naming
   - Avoids guessing command names
   - Gets authoritative, up-to-date information
   - Plugin help includes full structure information
   - Ensures compatibility with installed versions

5. **Workflow example:**
   ```bash
   # Step 1: Check what commands are available
   $ aliyun ecs --help
   # Output shows: create-instance, describe-instances, etc.

   # Step 2: Check specific command parameters
   $ aliyun ecs create-instance --help
   # Output shows: --instance-name, --image-id, --instance-type, etc.

   # Step 3: Construct command with exact names from help
   $ aliyun ecs create-instance \
       --instance-name my-instance \
       --image-id ubuntu_20_04 \
       --instance-type ecs.g6.large
   ```

### Quick Reference

| Aspect | OpenAPI Style (Built-in) | CLI Native Style (Plugin) |
|--------|--------------------------|---------------------------|
| **Command** | PascalCase | kebab-case |
| **Parameters** | **Mixed/Inconsistent** | **kebab-case (Consistent)** |
| **Naming** | Varies by API | Unified across all |
| **Detection** | Contains uppercase | All lowercase |
| **ROA Body** | Single `--body` | Expanded parameters |
| **Help Quality** | Basic | Comprehensive with structure |
| **Self-contained** | No (needs API docs) | Yes (complete info in CLI) |
| **Example** | `CreateInstance` (params vary) | `create-instance --instance-name` |
| **When to use** | No plugin available | Plugin installed (preferred) |

### Examples in Practice

**Checking what's available:**
```bash
# Check if plugin installed
aliyun plugin list | grep ecs

# If not, install it (use short name or full name, case-insensitive)
aliyun plugin install --names ecs              # Recommended
aliyun plugin install --names aliyun-cli-ecs   # Also works
aliyun plugin install --names ECS              # Case doesn't matter

# View plugin help
aliyun ecs --help

# View original built-in help
ALIBABA_CLOUD_ORIGINAL_PRODUCT_HELP=true aliyun ecs --help

# Get detailed command help (plugin)
aliyun ecs create-instance --help
# Shows all parameters with structure information

# Get detailed command help (built-in)
aliyun Ecs CreateInstance --help
# Shows basic parameter info
```

**Using the plugin (preferred):**
```bash
# List instances with plugin
aliyun ecs describe-instances \
  --biz-region-id cn-hangzhou \
  --instance-ids '["i-xxx"]' \
  --output json

# Every parameter documented in help:
# aliyun ecs describe-instances --help
```

**Fallback to built-in:**
```bash
# If plugin not available, use built-in
aliyun Ecs DescribeInstances \
  --RegionId cn-hangzhou \
  --InstanceIds '["i-xxx"]' \
  --output json
```

## Special Product Notes

### OSS (Object Storage Service) - Important

**OSS uses a custom implementation, NOT standard API-based commands.**

Unlike other Aliyun products, OSS in the CLI is implemented with:
- Hand-written code and external binary executables
- Custom command structure (NOT API triple mapping)
- Two command namespaces: `aliyun oss` and `aliyun ossutil`

**Critical points:**
- **NO API-based commands** like `PutBucket`, `GetObject`, etc.
- **Use custom OSS commands** found in help docs
- **Always check help** for correct command syntax:
  ```bash
  aliyun oss --help        # Basic OSS operations
  aliyun ossutil --help    # Advanced OSS utilities
  ```

**Example - Wrong vs Right:**
```bash
# WRONG - This does NOT exist
aliyun oss put-bucket --bucket-name my-bucket

# CORRECT - Use OSS custom commands
aliyun oss mb oss://my-bucket              # Make bucket
aliyun oss cp local.txt oss://bucket/      # Upload file
aliyun ossutil ls oss://bucket/            # List objects
```

**When users ask about OSS operations:**
1. Check `aliyun oss --help` first for basic operations
2. Check `aliyun ossutil --help` for advanced features
3. Do NOT assume API command patterns work for OSS
4. Always verify the actual command exists in the help docs

## Common Workflows

### Working with ECS Instances

```bash
# List instances (note: --biz-region-id, not --region-id)
aliyun ecs describe-instances --biz-region-id cn-hangzhou

# Create instance with structured parameters
aliyun ecs create-instance \
  --biz-region-id cn-hangzhou \
  --instance-type ecs.g6.large \
  --image-id ubuntu_20_04_x64 \
  --data-disk Category=cloud_essd Size=100 \
  --tag key=env value=prod --tag key=app value=web
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
  --biz-region-id cn-hangzhou
```

### Working with OSS (Object Storage)

**Important: OSS uses custom commands, not API-based patterns**

```bash
# Bucket operations
aliyun oss mb oss://my-bucket                    # Create bucket
aliyun oss ls                                     # List all buckets
aliyun oss ls oss://my-bucket/                   # List objects in bucket

# File operations
aliyun oss cp local-file.txt oss://bucket/path/  # Upload file
aliyun oss cp oss://bucket/file.txt ./           # Download file
aliyun oss rm oss://bucket/file.txt              # Delete object

# Advanced operations with ossutil
aliyun ossutil ls oss://bucket/ -d               # List with details
aliyun ossutil stat oss://bucket/file.txt        # Get object info
aliyun ossutil sync ./local/ oss://bucket/       # Sync directory

# Always check help for available commands
aliyun oss --help
aliyun ossutil --help
```

## Response Format

When providing command suggestions:
1. Explain what the command does
2. Show the complete command with all required parameters
3. Explain any non-obvious parameter values
4. Mention relevant flags or options
5. If debugging is needed, include `--log-level debug`

## Examples

### Example 1: User asks to list ECS instances

**User**: "Show me all my ECS instances in Hangzhou"

**Response**:
```bash
# First, ensure the ECS plugin is installed
aliyun plugin list | grep ecs
# If not installed: aliyun plugin install --names ecs

# List all ECS instances in Hangzhou region
# Note: Use --biz-region-id (not --region-id) to avoid global parameter conflict
aliyun ecs describe-instances --biz-region-id cn-hangzhou

# For more detailed output with specific fields
aliyun ecs describe-instances \
  --biz-region-id cn-hangzhou \
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
