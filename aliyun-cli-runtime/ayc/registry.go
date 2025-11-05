// Copyright (c) 2009-present, Alibaba Cloud All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ayc

import (
	"github.com/spf13/cobra"
)

type BaseCommand interface {
	GetCommand() *Command
}

type ExtensionCommand interface {
	GetCommand() *Command
}

func AutoRegisterCommands(root *cobra.Command, baseCommands []BaseCommand, extensionCommands []ExtensionCommand) error {
	extMap := make(map[string]ExtensionCommand)
	for _, extCmd := range extensionCommands {
		cmd := extCmd.GetCommand()
		if cmd != nil && cmd.Name != "" {
			extMap[cmd.Name] = extCmd
		}
	}

	for _, baseCmd := range baseCommands {
		cmd := baseCmd.GetCommand()
		if cmd == nil || cmd.Name == "" {
			continue
		}

		var finalCmd *Command
		if extCmd, hasExt := extMap[cmd.Name]; hasExt {
			finalCmd = extCmd.GetCommand()
		} else {
			finalCmd = cmd
		}
		RegisterCommand(root, finalCmd)
	}

	return nil
}
