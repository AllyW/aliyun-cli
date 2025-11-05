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
	"strings"
)

type Operation struct {
	ctx *OperationContext

	Method string

	// URL path (supports placeholders like {resourceId})
	URL string

	// API configuration
	Version     string // API version (e.g., "2015-12-15", "2014-05-26")
	Action      string // API action name (e.g., "DescribeClusters", "DescribeRegions")
	Protocol    string // Protocol (e.g., "HTTPS", "HTTP")
	Style       string // API style (e.g., "ROA", "RPC")
	ContentType string // Content-Type header (e.g., "application/json")

	QueryParameters map[string]string

	HeaderParameters map[string]string

	Content interface{}

	FormContent map[string]string

	ErrorMap map[int]error
}

func NewOperation(ctx *OperationContext) *Operation {
	return &Operation{
		ctx:              ctx,
		QueryParameters:  make(map[string]string),
		HeaderParameters: make(map[string]string),
		FormContent:      make(map[string]string),
		ErrorMap:         make(map[int]error),
		// Set default values
		Protocol:    "HTTPS",
		Style:       "ROA",
		ContentType: "application/json",
	}
}

func (op *Operation) SetMethod(method string) *Operation {
	op.Method = strings.ToUpper(method)
	return op
}

func (op *Operation) SetURL(url string) *Operation {
	op.URL = url
	return op
}

func (op *Operation) AddQueryParam(name string, value string) *Operation {
	if value != "" {
		op.QueryParameters[name] = value
	}
	return op
}

func (op *Operation) AddHeaderParam(name string, value string) *Operation {
	if value != "" {
		op.HeaderParameters[name] = value
	}
	return op
}

func (op *Operation) SetContent(content interface{}) *Operation {
	op.Content = content
	return op
}

func (op *Operation) SetFormContent(formContent map[string]string) *Operation {
	op.FormContent = formContent
	return op
}

func (op *Operation) SetErrorMap(errorMap map[int]error) *Operation {
	op.ErrorMap = errorMap
	return op
}

func (op *Operation) SetVersion(version string) *Operation {
	op.Version = version
	return op
}

func (op *Operation) SetAction(action string) *Operation {
	op.Action = action
	return op
}

func (op *Operation) SetProtocol(protocol string) *Operation {
	op.Protocol = protocol
	return op
}

func (op *Operation) SetStyle(style string) *Operation {
	op.Style = style
	return op
}

func (op *Operation) SetContentType(contentType string) *Operation {
	op.ContentType = contentType
	return op
}

func (op *Operation) SerializeURLParam(name string, value interface{}, required bool, skipQuote bool) error {
	params, err := SerializeURLParam(name, value, required, skipQuote)
	if err != nil {
		return err
	}
	for key, val := range params {
		op.URL = ReplacePathPlaceholders(op.URL, map[string]string{key: val})
	}
	return nil
}

func (op *Operation) SerializeQueryParam(name string, value interface{}, required bool, skipQuote bool, div string) error {
	params, err := SerializeQueryParam(name, value, required, skipQuote, div)
	if err != nil {
		return err
	}
	for key, val := range params {
		op.AddQueryParam(key, val)
	}
	return nil
}

func (op *Operation) MakeRequest() (*Response, error) {

	httpRequest := NewRequest()
	httpRequest.SetMethod(op.Method)
	httpRequest.SetPath(op.URL)

	// Use fields from Operation instead of hardcoded values
	if op.Version != "" {
		httpRequest.SetVersion(op.Version)
	}
	if op.Action != "" {
		httpRequest.SetAction(op.Action)
	}
	if op.Protocol != "" {
		httpRequest.SetProtocol(op.Protocol)
	}
	if op.Style != "" {
		httpRequest.SetStyle(op.Style)
	}
	if op.ContentType != "" {
		httpRequest.AddHeader("Content-Type", op.ContentType)
	}

	for key, value := range op.QueryParameters {
		httpRequest.AddQuery(key, value)
	}

	for key, value := range op.HeaderParameters {
		httpRequest.AddHeader(key, value)
	}

	if op.Content != nil {
		httpRequest.SetBody(op.Content)
	}

	executor := NewExecutor(op.ctx.GetHTTPClient())
	response, err := executor.Execute(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("request execution failed: %w", err)
	}

	if err := op.OnError(response); err != nil {
		return nil, err
	}

	return response, nil
}

func (op *Operation) OnError(response *Response) error {
	if response.IsSuccess() {
		return nil
	}

	statusCode := response.GetStatusCode()

	if customErr, ok := op.ErrorMap[statusCode]; ok {
		return fmt.Errorf("%w: %s", customErr, response.GetBodyString())
	}

	return fmt.Errorf("API request failed with status code %d: %s",
		statusCode, response.GetBodyString())
}
