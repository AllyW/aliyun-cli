// Copyright (c) 2009-present, Alibaba Cloud All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.

package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/aliyun/aliyun-cli/v3/cli"
)

const (
	executionLoggingSettingsFile = "execution_logging.json"
	defaultExecutionLogSubdir    = "logs/commands"
	defaultExecutionMaxFiles     = 500
	executionLogTrimBatch        = 5
	argValuePlaceholder          = "{}"
	// Optional env for batch tests: same value on all spawned aliyun processes to grep/collect logs together.
	envExecutionLogRunID = "ALIYUN_EXECUTION_LOG_RUN_ID"
	envExecutionLogJobID = "ALIYUN_EXECUTION_LOG_JOB_ID"
)

// Beijing time (Asia/Shanghai, UTC+8) for execution log timestamps and file names.
var beijingLocation = sync.OnceValue(func() *time.Location {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return time.FixedZone("CST", 8*3600)
	}
	return loc
})

func nowForExecutionLog() time.Time {
	return time.Now().In(beijingLocation())
}

// ExecutionLoggingSettings is stored in ~/.aliyun/execution_logging.json (not in profile config.json).
type ExecutionLoggingSettings struct {
	Enabled  bool   `json:"enabled"`
	LogDir   string `json:"log_dir,omitempty"`
	MaxFiles int    `json:"max_files,omitempty"`
}

func defaultExecutionLoggingSettings() ExecutionLoggingSettings {
	return ExecutionLoggingSettings{
		Enabled:  false,
		LogDir:   "",
		MaxFiles: defaultExecutionMaxFiles,
	}
}

func executionLoggingSettingsPath() string {
	return filepath.Join(GetConfigPath(), executionLoggingSettingsFile)
}

// LoadExecutionLoggingSettings reads global execution-log settings; missing file returns defaults.
func LoadExecutionLoggingSettings() (ExecutionLoggingSettings, error) {
	path := executionLoggingSettingsPath()
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return defaultExecutionLoggingSettings(), nil
		}
		return ExecutionLoggingSettings{}, err
	}
	var s ExecutionLoggingSettings
	if err := json.Unmarshal(data, &s); err != nil {
		return ExecutionLoggingSettings{}, err
	}
	if s.MaxFiles <= 0 {
		s.MaxFiles = defaultExecutionMaxFiles
	}
	return s, nil
}

// SaveExecutionLoggingSettings writes settings to ~/.aliyun/execution_logging.json.
func SaveExecutionLoggingSettings(s *ExecutionLoggingSettings) error {
	if s.MaxFiles <= 0 {
		s.MaxFiles = defaultExecutionMaxFiles
	}
	data, err := json.MarshalIndent(s, "", "\t")
	if err != nil {
		return err
	}
	return os.WriteFile(executionLoggingSettingsPath(), data, 0600)
}

func resolvedExecutionLogDir(s ExecutionLoggingSettings) string {
	if strings.TrimSpace(s.LogDir) != "" {
		return filepath.Clean(os.ExpandEnv(s.LogDir))
	}
	return filepath.Join(GetConfigPath(), defaultExecutionLogSubdir)
}

// ExecutionLogRecord is one command invocation; Extra holds forward-compatible fields.
type ExecutionLogRecord struct {
	Time         string                 `json:"time"`
	PID          int                    `json:"pid"`
	ExitCode     int                    `json:"exit_code"`
	Success      bool                   `json:"success"`
	DurationMs   int64                  `json:"duration_ms,omitempty"`
	ArgsRedacted []string               `json:"args_redacted,omitempty"`
	Error        string                 `json:"error,omitempty"`
	CommandKey   string                 `json:"command_key,omitempty"`
	RunID        string                 `json:"run_id,omitempty"` // from ALIYUN_EXECUTION_LOG_RUN_ID
	JobID        string                 `json:"job_id,omitempty"` // from ALIYUN_EXECUTION_LOG_JOB_ID (e.g. parallel slot)
	Extra        map[string]interface{} `json:"extra,omitempty"`
}

func isFlagToken(tok string) bool {
	return strings.HasPrefix(tok, "-") && !strings.HasPrefix(tok, "---") && len(tok) > 1
}

func commandPartsFromArgs(argv []string) []string {
	if len(argv) == 0 {
		return nil
	}
	var parts []string
	for _, a := range argv {
		if strings.HasPrefix(a, "-") {
			break
		}
		parts = append(parts, a)
		if len(parts) >= 3 {
			break
		}
	}
	return parts
}

// RedactExecutionArgs returns argv with command tokens preserved and all parameter values replaced by {}.
// Aligned with Azure CLI azlogging._get_clean_args (per-command local audit log).
func RedactExecutionArgs(argv []string) []string {
	if len(argv) == 0 {
		return nil
	}
	cmdLen := len(commandPartsFromArgs(argv))
	out := make([]string, 0, len(argv))
	for i, arg := range argv {
		if i < cmdLen {
			out = append(out, arg)
			continue
		}
		if !isFlagToken(arg) {
			out = append(out, argValuePlaceholder)
			continue
		}
		if !strings.HasPrefix(arg, "--") {
			opt := arg[:2]
			if len(arg) > 2 && arg[2] == '=' {
				opt += "=" + argValuePlaceholder
			} else if len(arg) > 2 {
				opt += argValuePlaceholder
			}
			out = append(out, opt)
			continue
		}
		if eq := strings.IndexByte(arg, '='); eq > 0 {
			out = append(out, arg[:eq]+"="+argValuePlaceholder)
			continue
		}
		out = append(out, arg)
	}
	return out
}

func commandKeyFromArgs(argv []string) string {
	return strings.Join(commandPartsFromArgs(argv), " ")
}

func truncateErr(s string, max int) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	if len(s) <= max {
		return s
	}
	if idx := strings.IndexByte(s, '\n'); idx > 0 && idx < max {
		return s[:idx]
	}
	return s[:max] + "..."
}

// LogExecutionIfEnabled writes one JSON record per invocation when enabled in execution_logging.json.
// Safe to call on every root run; failures are ignored.
func LogExecutionIfEnabled(ctx *cli.Context, argv []string, err error) {
	exitCode := cli.ExitCodeForError(err)
	writeExecutionLog(ctx, argv, exitCode, err, nil)
}

// LogExecutionIfEnabledWithExitCode writes one record with an explicit exit code (e.g. plugin subprocess).
// Use the same argv as the root audit trail, typically os.Args[1:]. Safe to call before os.Exit; failures are ignored.
func LogExecutionIfEnabledWithExitCode(ctx *cli.Context, argv []string, exitCode int, err error) {
	extra := map[string]interface{}{"invoker": "plugin"}
	writeExecutionLog(ctx, argv, exitCode, err, extra)
}

func writeExecutionLog(ctx *cli.Context, argv []string, exitCode int, err error, extra map[string]interface{}) {
	if ctx != nil && ctx.Completion() != nil {
		return
	}
	s, e := LoadExecutionLoggingSettings()
	if e != nil || !s.Enabled {
		return
	}
	logDir := resolvedExecutionLogDir(s)
	if e := os.MkdirAll(logDir, 0755); e != nil {
		return
	}

	rec := ExecutionLogRecord{
		Time:         nowForExecutionLog().Format(time.RFC3339Nano),
		PID:          os.Getpid(),
		ExitCode:     exitCode,
		Success:      exitCode == 0,
		CommandKey:   commandKeyFromArgs(argv),
		ArgsRedacted: RedactExecutionArgs(argv),
		Extra:        extra,
	}
	if ctx != nil {
		if ms := ctx.ExecutionDurationMs(); ms >= 0 {
			rec.DurationMs = ms
		}
	}
	if err != nil {
		rec.Error = truncateErr(err.Error(), 2000)
	} else if exitCode != 0 {
		rec.Error = fmt.Sprintf("exit status %d", exitCode)
	}

	if v := strings.TrimSpace(os.Getenv(envExecutionLogRunID)); v != "" {
		rec.RunID = v
	}
	if v := strings.TrimSpace(os.Getenv(envExecutionLogJobID)); v != "" {
		rec.JobID = v
	}

	data, e := json.Marshal(rec)
	if e != nil {
		return
	}

	ts := nowForExecutionLog().Format("20060102-150405")
	name := fmt.Sprintf("%s.%d.json", ts, os.Getpid())
	path := filepath.Join(logDir, name)
	if e := os.WriteFile(path, append(data, '\n'), 0600); e != nil {
		return
	}
	trimOldExecutionLogs(logDir, s.MaxFiles)
}

func trimOldExecutionLogs(logDir string, maxKeep int) {
	if maxKeep <= 0 {
		maxKeep = defaultExecutionMaxFiles
	}
	entries, err := os.ReadDir(logDir)
	if err != nil {
		return
	}
	var files []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		n := e.Name()
		if strings.HasSuffix(n, ".json") {
			files = append(files, n)
		}
	}
	if len(files) <= maxKeep {
		return
	}
	sort.Strings(files)
	toRemove := len(files) - maxKeep
	if toRemove > executionLogTrimBatch {
		toRemove = executionLogTrimBatch
	}
	for i := 0; i < toRemove && i < len(files); i++ {
		_ = os.Remove(filepath.Join(logDir, files[i]))
	}
}

// ExecutionLoggingSettingsPath returns the path to execution_logging.json (for configure / tests).
func ExecutionLoggingSettingsPath() string {
	return executionLoggingSettingsPath()
}

// ParseExecutionLogEnabled parses "true"/"false" or "1"/"0" for configure set.
func ParseExecutionLogEnabled(s string) (bool, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "true", "1", "yes", "on":
		return true, nil
	case "false", "0", "no", "off":
		return false, nil
	default:
		return false, fmt.Errorf("invalid boolean: %q (use true or false)", s)
	}
}

// ParseExecutionMaxFiles parses positive int; empty returns 0 (caller uses default).
func ParseExecutionMaxFiles(s string) (int, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, nil
	}
	n, err := strconv.Atoi(s)
	if err != nil || n < 1 {
		return 0, fmt.Errorf("max-files must be a positive integer")
	}
	return n, nil
}
