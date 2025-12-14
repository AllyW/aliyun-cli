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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Server服务客户端接口
type ServerClient interface {
	SendRequest(ctx context.Context, sessionID string, input string, requestType RequestType) (*ServerResponse, error)
}

// HTTP Server服务客户端实现
type HTTPServerClient struct {
	baseURL    string
	httpClient *http.Client
}

// 创建HTTP Server服务客户端
func NewHTTPServerClient(baseURL string) *HTTPServerClient {
	return &HTTPServerClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// 发送请求到 Server 端点
func (c *HTTPServerClient) SendRequest(ctx context.Context, sessionID string, input string, requestType RequestType) (*ServerResponse, error) {
	reqBody := ServerRequest{
		SessionID:  sessionID,
		Input:      input,
		CliMsgType: requestType,
	}

	reqData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := c.baseURL + "/ai/chat2"
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(reqData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call aggregate service: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应体（需要先读取，因为后面可能还要读取）
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("aggregate service returned error: %d, body: %s", resp.StatusCode, string(body))
	}

	var result ServerResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}
