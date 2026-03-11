# Aliyun CLI Installation & Configuration Guide

Complete guide for installing and configuring Aliyun CLI.

## Installation

### macOS

**Using Homebrew (Recommended)**
```bash
brew install aliyun-cli
```

**Using Binary**
```bash
# Download
wget https://aliyuncli.alicdn.com/aliyun-cli-macosx-latest-amd64.tgz

# Extract
tar -xzf aliyun-cli-macosx-latest-amd64.tgz

# Move to PATH
sudo mv aliyun /usr/local/bin/

# Verify
aliyun version
```

### Linux

**Debian/Ubuntu**
```bash
# Download
wget https://aliyuncli.alicdn.com/aliyun-cli-linux-latest-amd64.tgz

# Extract and install
tar -xzf aliyun-cli-linux-latest-amd64.tgz
sudo mv aliyun /usr/local/bin/

# Verify
aliyun version
```

**CentOS/RHEL**
```bash
# Download
wget https://aliyuncli.alicdn.com/aliyun-cli-linux-latest-amd64.tgz

# Extract and install
tar -xzf aliyun-cli-linux-latest-amd64.tgz
sudo mv aliyun /usr/local/bin/

# Verify
aliyun version
```

**ARM64 Architecture**
```bash
# Download ARM64 version
wget https://aliyuncli.alicdn.com/aliyun-cli-linux-latest-arm64.tgz

# Extract and install
tar -xzf aliyun-cli-linux-latest-arm64.tgz
sudo mv aliyun /usr/local/bin/
```

### Windows

**Using Binary**
1. Download from: https://aliyuncli.alicdn.com/aliyun-cli-windows-latest-amd64.zip
2. Extract the ZIP file
3. Add the directory to your PATH environment variable
4. Open new Command Prompt or PowerShell
5. Verify: `aliyun version`

**Using PowerShell**
```powershell
# Download
Invoke-WebRequest -Uri "https://aliyuncli.alicdn.com/aliyun-cli-windows-latest-amd64.zip" -OutFile "aliyun-cli.zip"

# Extract
Expand-Archive -Path aliyun-cli.zip -DestinationPath C:\aliyun-cli

# Add to PATH (requires admin privileges)
$env:Path += ";C:\aliyun-cli"
[Environment]::SetEnvironmentVariable("Path", $env:Path, [System.EnvironmentVariableTarget]::Machine)

# Verify
aliyun version
```

## Configuration

### Quick Start

**Interactive Configuration**
```bash
aliyun configure
```

This will prompt you for:
- **Access Key ID**: Your Aliyun access key ID
- **Access Key Secret**: Your Aliyun access key secret
- **Default Region**: e.g., `cn-hangzhou`, `cn-beijing`, `cn-shanghai`
- **Default Output Format**: `json` (recommended), `table`, or `cols`
- **Language**: `en` or `zh`

**Where to Get Access Keys**

1. Log in to Aliyun Console: https://ram.console.aliyun.com/
2. Navigate to: AccessKey Management
3. Create a new AccessKey pair
4. **Important**: Save the secret immediately - it's only shown once

### Configuration Modes

Aliyun CLI supports 6 authentication modes:

#### 1. AK Mode (Access Key)

**Most common mode for personal accounts**

```bash
aliyun configure --mode AK

# You'll be prompted for:
# - Access Key ID
# - Access Key Secret
# - Region
```

**Use Case**: Personal development, testing, scripts

**Configuration File** (`~/.aliyun/config.json`):
```json
{
  "current": "default",
  "profiles": [
    {
      "name": "default",
      "mode": "AK",
      "access_key_id": "your_access_key_id",
      "access_key_secret": "your_access_key_secret",
      "region_id": "cn-hangzhou",
      "output_format": "json",
      "language": "en"
    }
  ]
}
```

#### 2. StsToken Mode (Temporary Credentials)

**For temporary access or short-lived credentials**

```bash
aliyun configure --mode StsToken

# You'll be prompted for:
# - Access Key ID
# - Access Key Secret
# - STS Token
# - Region
```

**Use Case**:
- Temporary access for external contractors
- CI/CD pipelines with short-lived tokens
- Cross-account access

**STS Token expires** (usually 1-12 hours) - you'll need to refresh it.

#### 3. RamRoleArn Mode (Assume RAM Role)

**Assume a RAM role for elevated or cross-account access**

```bash
aliyun configure --mode RamRoleArn

# You'll be prompted for:
# - Access Key ID (of the user assuming the role)
# - Access Key Secret
# - RAM Role ARN (e.g., acs:ram::123456789012:role/MyRole)
# - Role Session Name (any identifier for your session)
# - Region
```

**Use Case**:
- Cross-account resource access
- Temporary elevated privileges
- Role-based access control

**Example RAM Role ARN**:
```
acs:ram::123456789012:role/AdminRole
```

#### 4. EcsRamRole Mode (ECS Instance RAM Role)

**Use RAM role attached to ECS instance - no credentials needed**

```bash
aliyun configure --mode EcsRamRole

# You'll be prompted for:
# - RAM Role Name (attached to your ECS instance)
# - Region
```

**Use Case**:
- Scripts running on ECS instances
- Automated workflows on ECS
- No need to store credentials

**Requirements**:
- Must be running on an ECS instance
- Instance must have a RAM role attached
- Role must have necessary permissions

**Example**:
```bash
# On ECS instance with role "MyEcsRole"
aliyun configure --mode EcsRamRole
# Enter: MyEcsRole
```

#### 5. RsaKeyPair Mode (RSA Key Pair)

**Use RSA key pair for authentication**

```bash
aliyun configure --mode RsaKeyPair

# You'll be prompted for:
# - Key Pair Name
# - Private Key File Path (e.g., ~/.ssh/aliyun_rsa)
# - Region
```

**Use Case**: Enhanced security with key-based auth

**Setup**:
1. Generate RSA key pair in Aliyun Console
2. Download private key file
3. Configure CLI with key pair name and file path

#### 6. RamRoleArnWithEcs Mode (ECS + RAM Role)

**Combine ECS instance role with RAM role assumption**

```bash
aliyun configure --mode RamRoleArnWithEcs

# You'll be prompted for:
# - RAM Role Name (on the ECS instance)
# - RAM Role ARN (to assume)
# - Role Session Name
# - Region
```

**Use Case**:
- ECS instance needs to assume another role
- Cross-account access from ECS

### Environment Variables

**Highest priority** - overrides config file

**Access Key Mode**
```bash
export ALIBABA_CLOUD_ACCESS_KEY_ID=your_access_key_id
export ALIBABA_CLOUD_ACCESS_KEY_SECRET=your_access_key_secret
export ALIBABA_CLOUD_REGION_ID=cn-hangzhou
```

**STS Token Mode**
```bash
export ALIBABA_CLOUD_ACCESS_KEY_ID=your_access_key_id
export ALIBABA_CLOUD_ACCESS_KEY_SECRET=your_access_key_secret
export ALIBABA_CLOUD_SECURITY_TOKEN=your_sts_token
export ALIBABA_CLOUD_REGION_ID=cn-hangzhou
```

**ECS RAM Role Mode**
```bash
export ALIBABA_CLOUD_ECS_METADATA=role_name
```

**Use Case**:
- CI/CD pipelines
- Docker containers
- Temporary credential override

### Managing Multiple Profiles

**Create Named Profiles**

```bash
# Create profile for project A
aliyun configure --profile projectA
# Enter credentials for project A

# Create profile for project B
aliyun configure --profile projectB
# Enter credentials for project B
```

**Use Specific Profile**

```bash
# Use profile in command
aliyun ecs describe-instances --profile projectA

# Set default profile via environment
export ALIBABA_CLOUD_PROFILE=projectA

# Now all commands use projectA
aliyun ecs describe-instances
```

**List Profiles**

```bash
# View all profiles
cat ~/.aliyun/config.json

# Current profile
aliyun configure get current
```

**Switch Default Profile**

Edit `~/.aliyun/config.json` and change the `"current"` field:
```json
{
  "current": "projectA",
  "profiles": [...]
}
```

### Credential Priority

Credentials are loaded in this order (first found wins):

1. **Command-line flag**: `--profile <name>`
2. **Environment variable**: `ALIBABA_CLOUD_PROFILE`
3. **Environment credentials**: `ALIBABA_CLOUD_ACCESS_KEY_ID`, etc.
4. **Configuration file**: `~/.aliyun/config.json` (current profile)
5. **ECS Instance RAM Role**: If running on ECS with attached role

## Verification

### Test Authentication

```bash
# Basic test - list regions
aliyun ecs describe-regions

# Expected output: JSON array of regions
```

**If successful**, you'll see:
```json
{
  "Regions": {
    "Region": [
      {
        "RegionId": "cn-hangzhou",
        "RegionEndpoint": "ecs.cn-hangzhou.aliyuncs.com",
        "LocalName": "华东 1（杭州）"
      },
      ...
    ]
  },
  "RequestId": "..."
}
```

**If failed**, you'll see error messages:
- `InvalidAccessKeyId.NotFound` - Wrong Access Key ID
- `SignatureDoesNotMatch` - Wrong Access Key Secret
- `InvalidSecurityToken.Expired` - STS token expired (for StsToken mode)
- `Forbidden.RAM` - Insufficient permissions

### Debug Configuration

```bash
# Show current configuration
aliyun configure get

# Test with debug logging
aliyun ecs describe-regions --log-level=debug

# Check credential provider
aliyun configure get mode
```

## Security Best Practices

### 1. Use RAM Users (Not Root Account)

❌ **Don't**: Use Aliyun root account credentials
✅ **Do**: Create RAM users with specific permissions

```bash
# Create RAM user in console
# Attach only necessary policies
# Use RAM user's access keys
```

### 2. Principle of Least Privilege

Grant only the minimum permissions needed:

```bash
# Example: Read-only ECS access
# Attach policy: AliyunECSReadOnlyAccess
```

### 3. Rotate Access Keys Regularly

```bash
# Create new access key
# Update configuration
aliyun configure

# Delete old access key from console
```

### 4. Use STS Tokens for Temporary Access

```bash
# Issue short-lived STS tokens instead of long-term keys
aliyun configure --mode StsToken
```

### 5. Use ECS RAM Roles When Possible

```bash
# No credentials in code or config
aliyun configure --mode EcsRamRole
```

### 6. Never Commit Credentials

```bash
# Add to .gitignore
echo "~/.aliyun/config.json" >> .gitignore

# Use environment variables in CI/CD instead
```

### 7. Secure Config File

```bash
# Restrict permissions
chmod 600 ~/.aliyun/config.json
```

## Troubleshooting

### Issue: Command Not Found

```bash
# Check installation
which aliyun

# Check PATH
echo $PATH

# Reinstall or add to PATH
```

### Issue: Authentication Failed

```bash
# Verify configuration
aliyun configure get

# Test with debug
aliyun ecs describe-regions --log-level=debug

# Check credentials in console
# Verify access key is active
```

### Issue: Permission Denied

```bash
# Error: Forbidden.RAM

# Check RAM user permissions
# Attach necessary policies in RAM console
# Example: AliyunECSFullAccess for ECS operations
```

### Issue: STS Token Expired

```bash
# Error: InvalidSecurityToken.Expired

# Refresh STS token
# Reconfigure with new token
aliyun configure --mode StsToken
```

### Issue: Wrong Region

```bash
# Some resources may not exist in the specified region

# Check available regions
aliyun ecs describe-regions

# Update default region
aliyun configure set region cn-shanghai
```

## Advanced Configuration

### Custom Endpoint

```bash
# Use custom or private endpoint
export ALIBABA_CLOUD_ECS_ENDPOINT=ecs-vpc.cn-hangzhou.aliyuncs.com
```

### Proxy Settings

```bash
# HTTP proxy
export HTTP_PROXY=http://proxy.example.com:8080
export HTTPS_PROXY=http://proxy.example.com:8080

# No proxy for specific domains
export NO_PROXY=localhost,127.0.0.1,.aliyuncs.com
```

### Timeout Settings

```bash
# Connection timeout (default: 10s)
export ALIBABA_CLOUD_CONNECT_TIMEOUT=30

# Read timeout (default: 10s)
export ALIBABA_CLOUD_READ_TIMEOUT=30
```

## Next Steps

After installation and configuration:

1. **Install plugins** for services you need:
   ```bash
   aliyun plugin install --names ecs,vpc,rds,oss
   ```

2. **Explore commands**:
   ```bash
   aliyun ecs --help
   aliyun fc --help
   ```

3. **Try example scenarios**:
   ```bash
   # See ./scenarios/ directory
   bash scenarios/Ecs/disk-lifecycle.sh
   ```

4. **Read documentation**:
   - [Command Syntax Guide](./command-syntax.md)
   - [Global Flags Reference](./global-flags.md)
   - [Common Scenarios](./common-scenarios.md)

## References

- Official Documentation: https://help.aliyun.com/zh/cli/
- RAM Console: https://ram.console.aliyun.com/
- Access Key Management: https://ram.console.aliyun.com/manage/ak
- Plugin Repository: https://github.com/aliyun/aliyun-cli
