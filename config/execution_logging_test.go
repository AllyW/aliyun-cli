// Copyright (c) 2009-present, Alibaba Cloud All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.

package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRedactExecutionArgs(t *testing.T) {
	in := []string{"ecs", "DescribeRegions", "--access-key-id", "SECRETID", "--region", "cn-hangzhou"}
	out := RedactExecutionArgs(in)
	assert.Equal(t, []string{"ecs", "DescribeRegions", "--access-key-id", "{}", "--region", "cn-hangzhou"}, out)

	eq := []string{"configure", "set", "--access-key-secret=topsecret"}
	out2 := RedactExecutionArgs(eq)
	assert.Equal(t, []string{"configure", "set", "--access-key-secret={}"}, out2)
}

func TestCommandKeyFromArgs(t *testing.T) {
	assert.Equal(t, "ecs DescribeRegions", commandKeyFromArgs([]string{"ecs", "DescribeRegions", "--region", "x"}))
	assert.Equal(t, "", commandKeyFromArgs([]string{}))
}

func TestLoadExecutionLoggingSettings_defaultFile(t *testing.T) {
	dir := t.TempDir()
	old := hookGetHomePath
	hookGetHomePath = func(fn func() string) func() string {
		return func() string { return dir }
	}
	t.Cleanup(func() { hookGetHomePath = old })

	s, err := LoadExecutionLoggingSettings()
	assert.NoError(t, err)
	assert.False(t, s.Enabled)
	assert.Equal(t, defaultExecutionMaxFiles, s.MaxFiles)
}

func TestTruncateResponseBody(t *testing.T) {
	s := strings.Repeat("a", 100)
	out := truncateResponseBody(s, 50)
	assert.Len(t, out, len(s[:50])+len("\n... [truncated by max_response_bytes]"))
	assert.Contains(t, out, "[truncated")
}

func TestSaveAndLoadExecutionLoggingSettings(t *testing.T) {
	dir := t.TempDir()
	old := hookGetHomePath
	hookGetHomePath = func(fn func() string) func() string {
		return func() string { return dir }
	}
	t.Cleanup(func() { hookGetHomePath = old })
	_ = os.MkdirAll(filepath.Join(dir, ".aliyun"), 0755)

	s := ExecutionLoggingSettings{Enabled: true, LogDir: "/tmp/logs", MaxFiles: 10}
	err := SaveExecutionLoggingSettings(&s)
	assert.NoError(t, err)

	s2, err := LoadExecutionLoggingSettings()
	assert.NoError(t, err)
	assert.True(t, s2.Enabled)
	assert.Equal(t, "/tmp/logs", s2.LogDir)
	assert.Equal(t, 10, s2.MaxFiles)
}
