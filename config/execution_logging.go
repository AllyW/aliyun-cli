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
	defaultMaxResponseBytes      = 512 * 1024 // 512 KiB per log when record_response is on
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
// Add new JSON fields here as needed for future options.
type ExecutionLoggingSettings struct {
	Enabled          bool   `json:"enabled"`
	LogDir           string `json:"log_dir,omitempty"`
	MaxFiles         int    `json:"max_files,omitempty"`
	RecordResponse   bool   `json:"record_response,omitempty"`
	VerboseArgs      bool   `json:"verbose_args,omitempty"`
	MaxResponseBytes int    `json:"max_response_bytes,omitempty"`
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
	ArgsRedacted []string               `json:"args_redacted,omitempty"`
	Args         []string               `json:"args,omitempty"` // full argv when verbose_args is true
	ResponseBody string                 `json:"response_body,omitempty"`
	PluginStderr string                 `json:"plugin_stderr,omitempty"`
	Error        string                 `json:"error,omitempty"`
	CommandKey   string                 `json:"command_key,omitempty"`
	RunID        string                 `json:"run_id,omitempty"`  // from ALIYUN_EXECUTION_LOG_RUN_ID
	JobID        string                 `json:"job_id,omitempty"`  // from ALIYUN_EXECUTION_LOG_JOB_ID (e.g. parallel slot)
	Extra        map[string]interface{} `json:"extra,omitempty"`
}

var sensitiveFlagNames = map[string]struct{}{
	"--" + AccessKeyIdFlagName:     {},
	"--" + AccessKeySecretFlagName: {},
	"--" + StsTokenFlagName:        {},
	"--" + PrivateKeyFlagName:      {},
	"--" + ProcessCommandFlagName:  {},
	"--" + OIDCTokenFileFlagName:   {},
	"--password":                    {},
	"--secret":                      {},
	"--token":                       {},
}

// RedactExecutionArgs returns a copy of argv with flag values redacted where appropriate.
func RedactExecutionArgs(argv []string) []string {
	if len(argv) == 0 {
		return nil
	}
	out := make([]string, 0, len(argv))
	for i := 0; i < len(argv); i++ {
		a := argv[i]
		if eq := strings.IndexByte(a, '='); eq > 0 {
			prefix := a[:eq]
			if isSensitiveFlagToken(prefix) {
				out = append(out, prefix+"={}")
				continue
			}
		}
		out = append(out, a)
		if isSensitiveFlagToken(a) {
			if i+1 < len(argv) && !strings.HasPrefix(argv[i+1], "-") {
				out = append(out, "{}")
				i++
			}
		}
	}
	return out
}

func isSensitiveFlagToken(tok string) bool {
	tok = strings.TrimSpace(tok)
	if _, ok := sensitiveFlagNames[tok]; ok {
		return true
	}
	lt := strings.ToLower(tok)
	if strings.HasPrefix(lt, "--access-key") && strings.Contains(lt, "secret") {
		return true
	}
	if strings.HasPrefix(lt, "--access-key-id") || lt == "--access-key-id" {
		return true
	}
	return false
}

func commandKeyFromArgs(argv []string) string {
	if len(argv) == 0 {
		return ""
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
	return strings.Join(parts, " ")
}

func truncateResponseBody(s string, max int) string {
	if max <= 0 {
		max = defaultMaxResponseBytes
	}
	if len(s) <= max {
		return s
	}
	return s[:max] + "\n... [truncated by max_response_bytes]"
}

// CaptureResponseForExecutionLog stores response text on ctx when record_response is enabled (local/testing).
func CaptureResponseForExecutionLog(ctx *cli.Context, responseBody string) {
	if ctx == nil || responseBody == "" {
		return
	}
	s, err := LoadExecutionLoggingSettings()
	if err != nil || !s.Enabled || !s.RecordResponse {
		return
	}
	maxB := s.MaxResponseBytes
	if maxB <= 0 {
		maxB = defaultMaxResponseBytes
	}
	ctx.SetExecutionLogResponse(truncateResponseBody(responseBody, maxB))
}

// CapturePluginStreamsForExecutionLog stores plugin stdout (as response) and stderr when record_response is enabled.
func CapturePluginStreamsForExecutionLog(ctx *cli.Context, stdoutText, stderrText string) {
	if ctx == nil {
		return
	}
	s, err := LoadExecutionLoggingSettings()
	if err != nil || !s.Enabled || !s.RecordResponse {
		return
	}
	maxB := s.MaxResponseBytes
	if maxB <= 0 {
		maxB = defaultMaxResponseBytes
	}
	ctx.SetExecutionLogResponse(truncateResponseBody(stdoutText, maxB))
	ctx.SetExecutionLogPluginStderr(truncateResponseBody(stderrText, maxB))
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
		Time:       nowForExecutionLog().Format(time.RFC3339Nano),
		PID:        os.Getpid(),
		ExitCode:   exitCode,
		Success:    exitCode == 0,
		CommandKey: commandKeyFromArgs(argv),
		Extra:      extra,
	}
	if s.VerboseArgs {
		rec.Args = append([]string(nil), argv...)
	} else {
		rec.ArgsRedacted = RedactExecutionArgs(argv)
	}
	if ctx != nil && s.RecordResponse {
		rec.ResponseBody = ctx.ExecutionLogResponse()
		rec.PluginStderr = ctx.ExecutionLogPluginStderr()
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
	// Oldest names first with sort.Strings on our timestamp format
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

// ParseExecutionMaxResponseBytes parses a non-negative int; 0 means use default limit in CaptureResponseForExecutionLog.
func ParseExecutionMaxResponseBytes(s string) (int, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, fmt.Errorf("empty value")
	}
	n, err := strconv.Atoi(s)
	if err != nil || n < 0 {
		return 0, fmt.Errorf("max-response-bytes must be a non-negative integer")
	}
	return n, nil
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
