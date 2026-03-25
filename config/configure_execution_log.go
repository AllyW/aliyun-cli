// Copyright (c) 2009-present, Alibaba Cloud All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.

package config

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aliyun/aliyun-cli/v3/cli"
	"github.com/aliyun/aliyun-cli/v3/i18n"
)

const (
	executionLogEnabledFlagName       = "enabled"
	executionLogLogDirFlagName        = "log-dir"
	executionLogMaxFilesFlagName      = "max-files"
	executionLogRecordResponseFlagName = "record-response"
	executionLogVerboseArgsFlagName   = "verbose-args"
	executionLogMaxResponseBytesName  = "max-response-bytes"
)

func NewConfigureExecutionLogCommand() *cli.Command {
	enabled := &cli.Flag{
		Name:         executionLogEnabledFlagName,
		Short:        i18n.T("enable or disable command execution logging", "启用或禁用命令执行日志"),
		AssignedMode: cli.AssignedDefault,
	}
	logDir := &cli.Flag{
		Name:         executionLogLogDirFlagName,
		Short:        i18n.T("directory for per-invocation JSON log files", "每次执行的 JSON 日志目录"),
		AssignedMode: cli.AssignedDefault,
	}
	maxFiles := &cli.Flag{
		Name:         executionLogMaxFilesFlagName,
		Short:        i18n.T("retention: max log files to keep (oldest trimmed in batches)", "保留的日志文件数量上限"),
		AssignedMode: cli.AssignedDefault,
	}
	recordResponse := &cli.Flag{
		Name:         executionLogRecordResponseFlagName,
		Short:        i18n.T("record API response body in execution log (local/testing)", "在执行日志中记录 API 响应体（本地/测试）"),
		AssignedMode: cli.AssignedDefault,
	}
	verboseArgs := &cli.Flag{
		Name:         executionLogVerboseArgsFlagName,
		Short:        i18n.T("log full argv without redaction (local/testing only)", "记录完整命令行参数（不脱敏，仅本地测试）"),
		AssignedMode: cli.AssignedDefault,
	}
	maxRespBytes := &cli.Flag{
		Name:         executionLogMaxResponseBytesName,
		Short:        i18n.T("max response_body bytes in log (0=default 512KiB)", "响应体写入日志的最大字节数（0 表示默认 512KiB）"),
		AssignedMode: cli.AssignedDefault,
	}

	getCmd := &cli.Command{
		Name:  "get",
		Short: i18n.T("show execution logging settings", "查看执行日志配置"),
		Usage: "get",
		Run: func(c *cli.Context, args []string) error {
			if len(args) > 0 {
				return cli.NewInvalidCommandError(args[0], c)
			}
			s, err := LoadExecutionLoggingSettings()
			if err != nil {
				return err
			}
			out, err := json.MarshalIndent(s, "", "\t")
			if err != nil {
				return err
			}
			cli.Println(c.Stdout(), string(out))
			cli.Printf(c.Stdout(), "\nsettings file: %s\n", ExecutionLoggingSettingsPath())
			cli.Printf(c.Stdout(), "effective log dir: %s\n", resolvedExecutionLogDir(s))
			return nil
		},
	}

	setCmd := &cli.Command{
		Name:  "set",
		Short: i18n.T("set execution logging options", "设置执行日志选项"),
		Usage: "set [--enabled true|false] [--log-dir <path>] [--max-files <n>] [--record-response true|false] [--verbose-args true|false] [--max-response-bytes <n>]",
		Run: func(c *cli.Context, args []string) error {
			if len(args) > 0 {
				return cli.NewInvalidCommandError(args[0], c)
			}
			s, err := LoadExecutionLoggingSettings()
			if err != nil {
				return err
			}
			fs := c.Flags()
			if enabledFlag := fs.Get(executionLogEnabledFlagName); enabledFlag != nil && enabledFlag.IsAssigned() {
				v, ok := enabledFlag.GetValue()
				if !ok {
					return fmt.Errorf("--%s requires a value", executionLogEnabledFlagName)
				}
				b, err := ParseExecutionLogEnabled(v)
				if err != nil {
					return err
				}
				s.Enabled = b
			}
			if logDirFlag := fs.Get(executionLogLogDirFlagName); logDirFlag != nil && logDirFlag.IsAssigned() {
				v, ok := logDirFlag.GetValue()
				if !ok {
					return fmt.Errorf("--%s requires a value", executionLogLogDirFlagName)
				}
				s.LogDir = strings.TrimSpace(v)
			}
			if mf := fs.Get(executionLogMaxFilesFlagName); mf != nil && mf.IsAssigned() {
				v, ok := mf.GetValue()
				if !ok {
					return fmt.Errorf("--%s requires a value", executionLogMaxFilesFlagName)
				}
				n, err := ParseExecutionMaxFiles(v)
				if err != nil {
					return err
				}
				if n > 0 {
					s.MaxFiles = n
				}
			}
			if rf := fs.Get(executionLogRecordResponseFlagName); rf != nil && rf.IsAssigned() {
				v, ok := rf.GetValue()
				if !ok {
					return fmt.Errorf("--%s requires a value", executionLogRecordResponseFlagName)
				}
				b, err := ParseExecutionLogEnabled(v)
				if err != nil {
					return err
				}
				s.RecordResponse = b
			}
			if vf := fs.Get(executionLogVerboseArgsFlagName); vf != nil && vf.IsAssigned() {
				v, ok := vf.GetValue()
				if !ok {
					return fmt.Errorf("--%s requires a value", executionLogVerboseArgsFlagName)
				}
				b, err := ParseExecutionLogEnabled(v)
				if err != nil {
					return err
				}
				s.VerboseArgs = b
			}
			if mb := fs.Get(executionLogMaxResponseBytesName); mb != nil && mb.IsAssigned() {
				v, ok := mb.GetValue()
				if !ok {
					return fmt.Errorf("--%s requires a value", executionLogMaxResponseBytesName)
				}
				n, err := ParseExecutionMaxResponseBytes(v)
				if err != nil {
					return err
				}
				s.MaxResponseBytes = n
			}
			if err := SaveExecutionLoggingSettings(&s); err != nil {
				return err
			}
			cli.Printf(c.Stdout(), "Saved execution logging settings to %s\n", ExecutionLoggingSettingsPath())
			return nil
		},
	}
	setCmd.Flags().Add(enabled)
	setCmd.Flags().Add(logDir)
	setCmd.Flags().Add(maxFiles)
	setCmd.Flags().Add(recordResponse)
	setCmd.Flags().Add(verboseArgs)
	setCmd.Flags().Add(maxRespBytes)

	parent := &cli.Command{
		Name:  "execution-log",
		Short: i18n.T("command execution audit logging (global, not in profile)", "命令执行审计日志（全局配置，不在 profile 内）"),
		Usage: "execution-log {get|set}",
		Run: func(c *cli.Context, args []string) error {
			if len(args) > 0 {
				return cli.NewInvalidCommandError(args[0], c)
			}
			return fmt.Errorf("usage: aliyun configure execution-log get | aliyun configure execution-log set [--enabled true|false] [--log-dir <path>] [--max-files <n>] [--record-response true|false] [--verbose-args true|false] [--max-response-bytes <n>]")
		},
	}
	parent.AddSubCommand(getCmd)
	parent.AddSubCommand(setCmd)
	return parent
}
