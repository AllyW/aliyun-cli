// Copyright (c) 2009-present, Alibaba Cloud All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.

package cli

import "sync"

var (
	postExecuteHookMu sync.Mutex
	postExecuteHook   func(*Context, []string, error)
)

// SetPostExecuteHook registers a callback invoked after the root command finishes (success or error),
// before processError exits the process. Use for execution audit logging.
func SetPostExecuteHook(h func(*Context, []string, error)) {
	postExecuteHookMu.Lock()
	defer postExecuteHookMu.Unlock()
	postExecuteHook = h
}

func runPostExecuteHook(ctx *Context, args []string, err error) {
	postExecuteHookMu.Lock()
	h := postExecuteHook
	postExecuteHookMu.Unlock()
	if h != nil {
		h(ctx, args, err)
	}
}

// ExitCodeForError maps errors to the same exit codes as processError.
func ExitCodeForError(err error) int {
	if err == nil {
		return 0
	}
	if _, ok := err.(SuggestibleError); ok {
		return 2
	}
	if _, ok := err.(ErrorWithTip); ok {
		return 3
	}
	return 1
}
