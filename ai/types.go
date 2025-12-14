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

// 对话中的一条消息
type Message struct {
	Role    string `json:"role"`    // system/user/assistant
	Content string `json:"content"` // 消息内容
}

// 响应类型
type ResponseType string

// 请求类型
type RequestType string

const (
	RequestTypePrompt  RequestType = "prompt"   // 用户首次运行
	RequestTypeExecute RequestType = "execute"  // 执行命令
	RequestTypeAnswer  RequestType = "answer"   // 回答question
	RequestTypeSelected RequestType = "selected" // 处理choose
)

const (
	ResponseTypeShow    ResponseType = "show"     // 直接展示消息
	ResponseTypeCommand ResponseType = "command"  // 展示消息并等待确认执行
	ResponseTypeChoose  ResponseType = "choose"   // 展示选项供用户选择
	ResponseTypeQuestion ResponseType = "question" // 等待用户输入
	ResponseTypeClose   ResponseType = "close"    // 退出会话
)

// Server API 响应
type ServerResponse struct {
	LlmMsgType  ResponseType `json:"llmMsgType"`   // 响应类型
	Message     string       `json:"message"`      // 消息内容
	Command     string       `json:"command"`      // 命令内容（用于 command 类型）
	ChooseItems []string     `json:"chooseItems"`  // 选项列表（用于 choose 类型）
	SessionID   string       `json:"sessionId"`    // 会话ID
}

// Server API 请求
type ServerRequest struct {
	SessionID string      `json:"sessionId"`     // 会话ID
	Input     string      `json:"userInput"`     // 用户输入
	CliMsgType RequestType `json:"cliMsgType"`   // 请求类型
}