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

package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/aliyun/aliyun-cli/aliyun-cli-plugins/plugin-cs/commands"
	"github.com/aliyun/aliyun-cli/aliyun-cli-plugins/plugin-cs/extension"
	"github.com/aliyun/aliyun-cli/aliyun-cli-runtime/ayc"
)

func main() {
	// 创建根命令
	rootCmd := &cobra.Command{
		Use:   "aliyun cs <command>",
		Short: "Aliyun CLI - CS Plugin",
		Long:  "Aliyun CLI plugin for Container Service (CS) operations.",
	}

	baseCmds := commands.AutoRegister()
	extCmds := extension.AutoRegister()

	if err := ayc.AutoRegisterCommands(rootCmd, baseCmds, extCmds); err != nil {
		fmt.Fprintf(os.Stderr, "Error registering commands: %v\n", err)
		os.Exit(1)
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
