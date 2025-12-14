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
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/aliyun/aliyun-cli/v3/cli"
)

// AI Agent 结构体
type Agent struct {
	serverClient ServerClient
	cliCtx       *cli.Context
	sessionID    string
}

// 创建Agent
func NewAgent(serverClient ServerClient, cliCtx *cli.Context) *Agent {
	return &Agent{
		serverClient: serverClient,
		cliCtx:       cliCtx,
		sessionID:    "",
	}
}

// 处理用户输入（主循环）
func (a *Agent) Process(ctx context.Context, initialInput string) error {
	userInput := initialInput

	// 主循环：持续处理响应直到会话关闭
	isFirstRequest := true
	var lastResponseType ResponseType

	for {
		// 确定请求类型
		var requestType RequestType
		if isFirstRequest {
			requestType = RequestTypePrompt
			isFirstRequest = false
		} else {
			switch lastResponseType {
			case ResponseTypeCommand:
				requestType = RequestTypeExecute
			case ResponseTypeQuestion:
				requestType = RequestTypeAnswer
			case ResponseTypeChoose:
				requestType = RequestTypeSelected
			default:
				requestType = RequestTypePrompt
			}
		}

		// 发送请求到 Server 服务（显示 loading 效果）
		resp, err := a.sendRequestWithLoading(ctx, a.sessionID, userInput, requestType)
		if err != nil {
			return fmt.Errorf("failed to send request: %w", err)
		}

		// 更新 sessionID
		if resp.SessionID != "" {
			a.sessionID = resp.SessionID
		}

		// 保存响应类型，用于下次请求时确定类型
		lastResponseType = resp.LlmMsgType

		// 根据响应类型处理，并获取下一次请求的输入
		nextInput, shouldContinue, err := a.handleResponse(ctx, resp)
		if err != nil {
			return err
		}

		if !shouldContinue {
			// 会话关闭
			break
		}

		// 如果需要继续，使用返回的 nextInput 作为下一次请求的输入
		if nextInput != "" {
			userInput = nextInput
		} else {
			// 如果没有返回输入，退出循环
			break
		}
	}

	return nil
}

// 处理响应
// 返回值: (nextInput, shouldContinue, error)
// nextInput: 下一次请求的输入内容，如果为空则不再继续
// shouldContinue: 是否应该继续会话
func (a *Agent) handleResponse(ctx context.Context, resp *ServerResponse) (string, bool, error) {
	_ = ctx // 保留参数以备将来使用
	switch resp.LlmMsgType {
	case ResponseTypeShow:
		// 直接展示消息
		cli.Printf(a.cliCtx.Stdout(), "%s\n", resp.Message)
		return "", false, nil // show 类型后不需要继续

	case ResponseTypeCommand:
		// 展示消息并等待确认执行
		// 直接输出，覆盖清除 loading 的行
		cli.Printf(a.cliCtx.Stdout(), "%s\n", resp.Message)
		if resp.Command != "" {
			cli.Printf(a.cliCtx.Stdout(), "命令: %s\n", resp.Command)
		}
		confirmed := a.confirmCommand()
		if confirmed {
			// 执行命令
			if resp.Command == "" {
				cli.Printf(a.cliCtx.Stderr(), "错误: 未提供命令\n")
			} else if err := a.executeCommand(resp.Command); err != nil {
				cli.Printf(a.cliCtx.Stderr(), "执行命令失败: %v\n", err)
			}
			// 执行后，只发送确认结果 "y"，执行结果已在本地显示，不需要发送给服务端
			return "y", true, nil
		} else {
			// 用户拒绝执行，直接结束会话，不再发起请求
			cli.Printf(a.cliCtx.Stdout(), "已取消执行\n")
			return "", false, nil
		}

	case ResponseTypeChoose:
		// 展示选项供用户选择
		cli.Printf(a.cliCtx.Stdout(), "%s\n", resp.Message)
		selected := a.chooseOption(resp.ChooseItems)
		if selected == "" {
			return "", false, nil // 用户取消选择，退出
		}
		// 将选择结果作为输入发送回去
		return selected, true, nil

	case ResponseTypeQuestion:
		// 等待用户输入
		cli.Printf(a.cliCtx.Stdout(), "%s\n", resp.Message)
		answer := a.getUserInput()
		if answer == "" {
			return "", false, nil // 用户取消输入，退出
		}
		// 将答案作为输入发送回去
		return answer, true, nil

	case ResponseTypeClose:
		// 退出会话
		if resp.Message != "" {
			cli.Printf(a.cliCtx.Stdout(), "%s\n", resp.Message)
		}
		return "", false, nil // 结束会话

	default:
		return "", false, fmt.Errorf("unknown response type: %s", resp.LlmMsgType)
	}
}

// 确认命令执行
func (a *Agent) confirmCommand() bool {
	cli.Printf(a.cliCtx.Stdout(), "\n是否执行？(y/n): ")
	// 确保输出刷新
	if w, ok := a.cliCtx.Stdout().(interface{ Flush() error }); ok {
		w.Flush()
	}
	
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		cli.Printf(a.cliCtx.Stderr(), "读取输入失败: %v\n", err)
		return false
	}
	input = strings.TrimSpace(strings.ToLower(input))
	return input == "y" || input == "yes"
}

// 选择选项
func (a *Agent) chooseOption(items []string) string {
	if len(items) == 0 {
		return ""
	}

	cli.Printf(a.cliCtx.Stdout(), "\n请选择（输入序号）:\n")
	for i, item := range items {
		cli.Printf(a.cliCtx.Stdout(), "%d. %s\n", i+1, item)
	}
	cli.Printf(a.cliCtx.Stdout(), "请输入序号 (1-%d): ", len(items))
	// 确保输出刷新
	if w, ok := a.cliCtx.Stdout().(interface{ Flush() error }); ok {
		w.Flush()
	}

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		cli.Printf(a.cliCtx.Stderr(), "读取输入失败: %v\n", err)
		return ""
	}
	input = strings.TrimSpace(input)
	index, err := strconv.Atoi(input)
	if err != nil || index < 1 || index > len(items) {
		cli.Printf(a.cliCtx.Stderr(), "无效的选择\n")
		return ""
	}
	return items[index-1]
}

// 获取用户输入
func (a *Agent) getUserInput() string {
	cli.Printf(a.cliCtx.Stdout(), "> ")
	// 确保输出刷新
	if w, ok := a.cliCtx.Stdout().(interface{ Flush() error }); ok {
		w.Flush()
	}
	
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		cli.Printf(a.cliCtx.Stderr(), "读取输入失败: %v\n", err)
		return ""
	}
	return strings.TrimSpace(input)
}

// 执行命令（直接执行命令字符串）
func (a *Agent) executeCommand(commandStr string) error {
	// 接口返回的命令可以直接执行，原样执行
	if strings.TrimSpace(commandStr) == "" {
		return fmt.Errorf("命令为空")
	}

	// 使用 shell 来执行命令，这样可以正确处理引号、转义等
	// 在 macOS/Linux 上使用 sh，它会正确处理引号和参数解析
	cmd := exec.Command("/bin/sh", "-c", commandStr)

	// 捕获输出
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// 执行并记录时间
	start := time.Now()
	err := cmd.Run()
	duration := time.Since(start)

	// 显示执行结果
	cli.Printf(a.cliCtx.Stdout(), "命令: %s\n", commandStr)
	if err != nil {
		exitCode := 0
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode = exitError.ExitCode()
		}
		cli.Printf(a.cliCtx.Stdout(), "✗ 执行失败 (退出码: %d)\n", exitCode)
		if stderr.String() != "" {
			cli.Printf(a.cliCtx.Stderr(), "\n错误:\n%s\n", stderr.String())
		}
		if stdout.String() != "" {
			cli.Printf(a.cliCtx.Stdout(), "\n输出:\n%s\n", stdout.String())
		}
	} else {
		cli.Printf(a.cliCtx.Stdout(), "✓ 执行成功\n")
		if stdout.String() != "" {
			cli.Printf(a.cliCtx.Stdout(), "\n输出:\n%s\n", stdout.String())
		}
	}
	cli.Printf(a.cliCtx.Stdout(), "耗时: %v\n", duration)

	return err
}

// 发送请求并显示 loading 效果
func (a *Agent) sendRequestWithLoading(ctx context.Context, sessionID string, input string, requestType RequestType) (*ServerResponse, error) {
	// 创建用于控制 loading 的 channel
	stopLoading := make(chan bool)
	var wg sync.WaitGroup
	
	// 启动 loading 动画
	wg.Add(1)
	go func() {
		defer wg.Done()
		a.showLoading(stopLoading)
	}()
	
	// 发送请求
	resp, err := a.serverClient.SendRequest(ctx, sessionID, input, requestType)
	
	// 停止 loading
	close(stopLoading)
	wg.Wait()
	
	// 清除 loading 行（不换行，让后续输出覆盖这一行）
	cli.Printf(a.cliCtx.Stdout(), "\r%s\r", strings.Repeat(" ", 50))
	if w, ok := a.cliCtx.Stdout().(interface{ Flush() error }); ok {
		w.Flush()
	}
	
	return resp, err
}

// 显示 loading 动画
func (a *Agent) showLoading(stop chan bool) {
	frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	index := 0
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	
	for {
		select {
		case <-stop:
			return
		case <-ticker.C:
			// 使用 \r 回到行首，覆盖之前的输出
			cli.Printf(a.cliCtx.Stdout(), "\r%s 正在处理中...", frames[index])
			if w, ok := a.cliCtx.Stdout().(interface{ Flush() error }); ok {
				w.Flush()
			}
			index = (index + 1) % len(frames)
		}
	}
}

