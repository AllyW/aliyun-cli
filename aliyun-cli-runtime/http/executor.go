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

package http

import (
	"fmt"
)

type Executor struct {
	client *Client
}

func NewExecutor(client *Client) *Executor {
	return &Executor{
		client: client,
	}
}

func (e *Executor) Execute(req *Request) (*Response, error) {
	if e.client == nil {
		return nil, fmt.Errorf("client is required")
	}

	if e.client.openapiClient == nil {
		return nil, fmt.Errorf("openapi client is not initialized")
	}

	if e.client.runtime == nil {
		return nil, fmt.Errorf("runtime options are not initialized")
	}

	if req == nil {
		return nil, fmt.Errorf("request is required")
	}

	params := req.GetOpenAPIParams()
	if params == nil {
		return nil, fmt.Errorf("request params are not initialized")
	}

	openapiReq := req.GetOpenAPIRequest()
	if openapiReq == nil {
		return nil, fmt.Errorf("openapi request is not initialized")
	}

	// Validate required fields before execution
	if params.Method == nil {
		return nil, fmt.Errorf("request method is required but not set")
	}
	if params.Pathname == nil {
		return nil, fmt.Errorf("request pathname is required but not set")
	}
	if params.Version == nil {
		return nil, fmt.Errorf("request version is required but not set")
	}

	// Execute the request using the OpenAPI client
	openapiResponse, err := e.client.openapiClient.CallApi(
		params,
		openapiReq,
		e.client.runtime,
	)

	if err != nil {
		return nil, fmt.Errorf("request execution failed: %w", err)
	}

	return NewResponse(openapiResponse), nil
}

func (e *Executor) ExecuteWithRequest(method, path string, options ...RequestOption) (*Response, error) {
	req := NewRequest()
	req.SetMethod(method).SetPath(path)

	for _, opt := range options {
		opt(req)
	}

	return e.Execute(req)
}

type RequestOption func(*Request)

func WithBody(body interface{}) RequestOption {
	return func(req *Request) {
		req.SetBody(body)
	}
}

func WithHeaders(headers map[string]string) RequestOption {
	return func(req *Request) {
		req.SetHeaders(headers)
	}
}

func WithQuery(query map[string]string) RequestOption {
	return func(req *Request) {
		req.SetQuery(query)
	}
}

func WithAction(action string) RequestOption {
	return func(req *Request) {
		req.SetAction(action)
	}
}

func WithVersion(version string) RequestOption {
	return func(req *Request) {
		req.SetVersion(version)
	}
}
