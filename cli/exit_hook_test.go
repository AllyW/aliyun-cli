// Copyright (c) 2009-present, Alibaba Cloud All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.

package cli

import (
	"bytes"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExitCodeForError(t *testing.T) {
	assert.Equal(t, 0, ExitCodeForError(nil))
	assert.Equal(t, 1, ExitCodeForError(errors.New("x")))

	ctx := NewCommandContext(new(bytes.Buffer), new(bytes.Buffer))
	ctx.EnterCommand(&Command{Name: "root", flags: NewFlagSet()})
	err := NewInvalidCommandError("bad", ctx)
	assert.Equal(t, 2, ExitCodeForError(err))

	tipErr := NewErrorWithTip(errors.New("inner"), "hint")
	assert.Equal(t, 3, ExitCodeForError(tipErr))
}
