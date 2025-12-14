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

package ai

import (
	"context"
	"fmt"
	"strings"

	"github.com/aliyun/aliyun-cli/v3/cli"
	"github.com/aliyun/aliyun-cli/v3/i18n"
)

// 创建AI命令
func NewAICommand() *cli.Command {
	cmd := &cli.Command{
		Name:  "ai",
		Short: i18n.T("AI assistant for Alibaba Cloud CLI", "阿里云CLI AI助手"),
		Long: i18n.T(
			"Use natural language to interact with Alibaba Cloud CLI. "+
				"The AI assistant will understand your requirements, "+
				"ask for missing information, generate execution plans, "+
				"and execute commands with your confirmation.",
			"使用自然语言与阿里云CLI交互。"+
				"AI助手将理解您的需求，"+
				"询问缺失信息，生成执行计划，"+
				"并在您确认后执行命令。",
		),
		Usage:  "aliyun ai <query>",
		Sample: "aliyun ai \"帮我创建一个ECS实例\"",
		Run: func(ctx *cli.Context, args []string) error {
			return runAI(ctx, args)
		},
	}

	return cmd
}

func runAI(ctx *cli.Context, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("please provide a query, e.g., aliyun ai \"help me create an ECS instance\"")
	}

	// 合并所有参数作为查询
	query := strings.Join(args, " ")

	// 创建 Aggregate 服务客户端
	aggregateURL := "http://localhost:8000"
	serverClient := NewHTTPServerClient(aggregateURL)
	
	agent := NewAgent(serverClient, ctx)

	// 处理用户查询
	agentCtx := context.Background()
	if err := agent.Process(agentCtx, query); err != nil {
		return fmt.Errorf("failed to process query: %w", err)
	}

	return nil
}

