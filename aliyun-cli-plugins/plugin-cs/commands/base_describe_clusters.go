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

package commands

import (
	"fmt"

	"github.com/aliyun/aliyun-cli/aliyun-cli-runtime/ayc"
	"github.com/aliyun/aliyun-cli/aliyun-cli-runtime/http"
)

type BaseDescribeClustersCommand struct {
	*ayc.DefaultExtensibleCommand
	BaseExtensibleCommand *ayc.BaseExtensibleCommand
}

func NewBaseDescribeClustersCommand() *BaseDescribeClustersCommand {
	cmd := &BaseDescribeClustersCommand{
		DefaultExtensibleCommand: &ayc.DefaultExtensibleCommand{},
	}

	baseCmd := ayc.NewBaseExtensibleCommand(cmd)
	baseCmd.Command.Name = "cs describe-clusters"
	baseCmd.Command.Description = "Describe clusters."
	baseCmd.Command.Examples = []string{
		"aliyun cs describe-clusters",
		"aliyun cs describe-clusters --cluster-type Standard --name MyCluster --resource-group-id rg-123",
	}

	cmd.BaseExtensibleCommand = baseCmd
	return cmd
}

func (c *BaseDescribeClustersCommand) GetCommand() *ayc.Command {
	return c.BaseExtensibleCommand.Command
}

func (c *BaseDescribeClustersCommand) Arguments() *ayc.Arguments {
	schema := ayc.NewArguments()

	schema.AddField("cluster_type", ayc.NewStringArg(
		[]string{"--cluster-type"},
		"The type of the cluster.",
		false,
	))

	schema.AddField("name", ayc.NewStringArg(
		[]string{"-n", "--name"},
		"The name of the cluster.",
		false,
	))

	schema.AddField("resource_group_id", ayc.NewStringArg(
		[]string{"-g", "--resource-group-id"},
		"The resource group ID.",
		false,
	))

	return schema
}

func (c *BaseDescribeClustersCommand) Execute(args map[string]any) error {
	ctx, err := http.NewOperationContext("", http.WithServiceEndpoint("cs", true))
	if err != nil {
		return fmt.Errorf("failed to create HTTP context: %w", err)
	}

	op := http.NewOperation(ctx)
	op.SetMethod("GET")
	op.SetURL("/clusters")
	op.SetVersion("2015-12-15")
	op.SetAction("DescribeClusters")
	op.SetProtocol("HTTPS")
	op.SetStyle("ROA")
	op.SetContentType("application/json")

	if clusterType, ok := args["cluster_type"]; ok && clusterType != nil {
		if err := op.SerializeQueryParam("clusterType", clusterType, false, false, ""); err != nil {
			return fmt.Errorf("failed to serialize cluster_type: %w", err)
		}
	}

	if name, ok := args["name"]; ok && name != nil {
		if err := op.SerializeQueryParam("name", name, false, false, ""); err != nil {
			return fmt.Errorf("failed to serialize name: %w", err)
		}
	}

	if resourceGroupID, ok := args["resource_group_id"]; ok && resourceGroupID != nil {
		if err := op.SerializeQueryParam("resource_group_id", resourceGroupID, false, false, ""); err != nil {
			return fmt.Errorf("failed to serialize resource_group_id: %w", err)
		}
	}

	// 使用 ExecuteOperation 自动处理 pager/waiter 和 output filter
	// 如果配置了 pager，会自动合并多页结果
	// 如果配置了 waiter，会轮询直到满足条件
	// 如果配置了 output filter，会格式化为表格
	output, err := ayc.ExecuteOperation(op, args)
	if err != nil {
		return fmt.Errorf("API request failed: %w", err)
	}

	fmt.Println(output)

	return nil
}
