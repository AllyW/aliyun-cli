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
package cli

import (
	"fmt"
	"text/tabwriter"
)

func (c *Command) PrintHead(ctx *Context) {
	Printf(ctx.Stdout(), "%s\n", c.Short.Text())

	if ctx != nil {
		aiModeFlag := ctx.Flags().Get("ai-mode")
		if aiModeFlag != nil && aiModeFlag.IsAssigned() {
			Printf(ctx.Stdout(), "\n=== AI Mode: Additional Diagnostic Information ===\n")
			Printf(ctx.Stdout(), "\nRequestId-based Troubleshooting Guide:\n")
			Printf(ctx.Stdout(), "1. When encountering errors, look for the 'RequestId' field in the error response\n")
			Printf(ctx.Stdout(), "2. Use the RequestId to query logs and trace the request flow:\n")
			Printf(ctx.Stdout(), "   - Check API Gateway logs using the RequestId\n")
			Printf(ctx.Stdout(), "   - Query service logs filtered by RequestId\n")
			Printf(ctx.Stdout(), "   - Use RequestId in support tickets for faster resolution\n")
			Printf(ctx.Stdout(), "3. Common RequestId locations:\n")
			Printf(ctx.Stdout(), "   - HTTP response headers: 'x-acs-request-id'\n")
			Printf(ctx.Stdout(), "   - Error response body: 'RequestId' field\n")
			Printf(ctx.Stdout(), "   - Log files: search for RequestId pattern\n")
			Printf(ctx.Stdout(), "4. For distributed tracing:\n")
			Printf(ctx.Stdout(), "   - RequestId helps correlate logs across multiple services\n")
			Printf(ctx.Stdout(), "   - Use RequestId to track request lifecycle end-to-end\n")
			Printf(ctx.Stdout(), "\nAdditional Context:\n")
			Printf(ctx.Stdout(), "- All API responses include a RequestId for traceability\n")
			Printf(ctx.Stdout(), "- RequestId format: typically UUID or alphanumeric string\n")
			Printf(ctx.Stdout(), "- Store RequestId when reporting issues for faster support\n")
			Printf(ctx.Stdout(), "\n================================================\n")
		} else {
			// Make the hint more prominent and AI-friendly
			Printf(ctx.Stdout(), "\nAI Assistant Tip: Add --ai-mode flag to get enhanced diagnostic information,\n")
			Printf(ctx.Stdout(), "   including RequestId-based troubleshooting guidance for better error analysis.\n")
		}
	} else {
		Printf(ctx.Stdout(), "\nAI Assistant Tip: Add --ai-mode flag to get enhanced diagnostic information,\n")
		Printf(ctx.Stdout(), "   including RequestId-based troubleshooting guidance for better error analysis.\n")
	}
}

func (c *Command) PrintUsage(ctx *Context) {
	if c.Usage != "" {
		Printf(ctx.Stdout(), "\nUsage:\n  %s\n", c.GetUsageWithParent())
	} else {
		c.PrintSubCommands(ctx)
	}
}

func (c *Command) PrintSample(ctx *Context) {
	if c.Sample != "" {
		Printf(ctx.Stdout(), "\nSample:\n  %s\n", c.Sample)
	}
}

func (c *Command) PrintSubCommands(ctx *Context) {
	if len(c.subCommands) > 0 {
		Printf(ctx.Stdout(), "\nCommands:\n")
		w := tabwriter.NewWriter(ctx.Stdout(), 8, 0, 1, ' ', 0)
		for _, cmd := range c.subCommands {
			if cmd.Hidden {
				continue
			}
			fmt.Fprintf(w, "  %s\t%s\n", cmd.Name, cmd.Short.Text())
		}
		w.Flush()
	}
}

func (c *Command) PrintFlags(ctx *Context) {
	if len(c.Flags().Flags()) == 0 {
		return
	}
	Printf(ctx.Stdout(), "\nFlags:\n")
	w := tabwriter.NewWriter(ctx.Stdout(), 8, 0, 1, ' ', 0)
	fs := c.Flags()
	if ctx != nil {
		fs = ctx.Flags()
	}

	var aiModeFlag *Flag
	var otherFlags []*Flag

	for _, flag := range fs.Flags() {
		if flag.Hidden {
			continue
		}
		if flag.Name == "ai-mode" {
			aiModeFlag = flag
		} else {
			otherFlags = append(otherFlags, flag)
		}
	}

	if aiModeFlag != nil {
		s := "--" + aiModeFlag.Name
		if aiModeFlag.Shorthand != 0 {
			s = s + ",-" + string(aiModeFlag.Shorthand)
		}
		fmt.Fprintf(w, "  %s\t%s AI-friendly mode\n", s, aiModeFlag.Short.Text())
	}

	for _, flag := range otherFlags {
		s := "--" + flag.Name
		if flag.Shorthand != 0 {
			s = s + ",-" + string(flag.Shorthand)
		}
		fmt.Fprintf(w, "  %s\t%s\n", s, flag.Short.Text())
	}
	w.Flush()
}

func (c *Command) PrintFailed(ctx *Context, err error, suggestion string) {
	Errorf(ctx.Stderr(), "ERROR: %v\n", err)
	Printf(ctx.Stderr(), "%s\n", suggestion)
}

func (c *Command) PrintTail(ctx *Context) {
	Printf(ctx.Stdout(), "\nUse `%s --help` for more information.\n", c.Name)
}
