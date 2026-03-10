# CS 场景脚本

Web Application Firewall v3 (Web 应用防火墙 v3) 相关场景脚本。

**场景数量**: 3

## 可用场景

### Addon

**场景文件**: `Addon.sh`

**运行方式**:
```bash
bash scenarios/CS/Addon.sh
```

**说明**: 此场景从 `codesample-json-repo/CS/Addon.json` 自动生成。

---

### arms_prometheus_alert_rule

**场景文件**: `arms_prometheus_alert_rule.sh`

**运行方式**:
```bash
bash scenarios/CS/arms_prometheus_alert_rule.sh
```

**说明**: 此场景从 `codesample-json-repo/CS/arms_prometheus_alert_rule.json` 自动生成。

---

### EnvFeature

**场景文件**: `EnvFeature.sh`

**运行方式**:
```bash
bash scenarios/CS/EnvFeature.sh
```

**说明**: 此场景从 `codesample-json-repo/CS/EnvFeature.json` 自动生成。

---


## 前置条件

```bash
# 安装 CS 插件
aliyun plugin install --names cs

# 配置凭证
aliyun configure
```

## 参考资料

- [Plugin CLI 命令语法](../../references/command-syntax.md)
- [全局参数参考](../../references/global-flags.md)
- [场景索引](../../SCENARIOS_INDEX.md)

## 注意事项

这些脚本由自动转换工具生成，可能需要根据实际环境调整：
- 区域和可用区配置
- 实例规格和资源配额
- 网络配置（VPC、VSwitch 等）
- 认证和权限设置

建议先使用 `--cli-dry-run` 参数测试命令。
