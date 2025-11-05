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
	openapiClient "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	openapiutil "github.com/alibabacloud-go/darabonba-openapi/v2/utils"
	"github.com/alibabacloud-go/tea/tea"
)

type Request struct {
	openapiRequest *openapiutil.OpenApiRequest
	openapiParams  *openapiClient.Params
}

func NewRequest() *Request {
	params := &openapiClient.Params{}
	// Set default values (can be overridden by Set methods)
	params.AuthType = tea.String("AK")
	params.Style = tea.String("ROA")
	params.ReqBodyType = tea.String("json")
	params.BodyType = tea.String("json")
	params.Protocol = tea.String("HTTPS")

	return &Request{
		openapiRequest: &openapiutil.OpenApiRequest{
			Query:   make(map[string]*string),
			Headers: make(map[string]*string),
			HostMap: make(map[string]*string),
		},
		openapiParams: params,
	}
}

func (r *Request) SetMethod(method string) *Request {
	r.openapiParams.Method = tea.String(method)
	return r
}

func (r *Request) SetPath(path string) *Request {
	r.openapiParams.Pathname = tea.String(path)
	return r
}

func (r *Request) SetAction(action string) *Request {
	r.openapiParams.Action = tea.String(action)
	return r
}

func (r *Request) SetVersion(version string) *Request {
	r.openapiParams.Version = tea.String(version)
	return r
}

func (r *Request) SetProtocol(protocol string) *Request {
	r.openapiParams.Protocol = tea.String(protocol)
	return r
}

func (r *Request) SetStyle(style string) *Request {
	r.openapiParams.Style = tea.String(style)
	return r
}

func (r *Request) SetBody(body interface{}) *Request {
	switch v := body.(type) {
	case string:
		r.openapiRequest.SetBody([]byte(v))
	case []byte:
		r.openapiRequest.SetBody(v)
	case map[string]interface{}:
		r.openapiRequest.Body = v
	}
	return r
}

func (r *Request) SetBodyType(bodyType string) *Request {
	r.openapiParams.ReqBodyType = tea.String(bodyType)
	return r
}

func (r *Request) SetResponseBodyType(bodyType string) *Request {
	r.openapiParams.BodyType = tea.String(bodyType)
	return r
}

func (r *Request) AddHeader(key, value string) *Request {
	r.openapiRequest.Headers[key] = tea.String(value)
	return r
}

func (r *Request) SetHeaders(headers map[string]string) *Request {
	for k, v := range headers {
		r.openapiRequest.Headers[k] = tea.String(v)
	}
	return r
}

func (r *Request) AddQuery(key, value string) *Request {
	r.openapiRequest.Query[key] = tea.String(value)
	return r
}

func (r *Request) SetQuery(query map[string]string) *Request {
	for k, v := range query {
		r.openapiRequest.Query[k] = tea.String(v)
	}
	return r
}

func (r *Request) AddHostParam(key, value string) *Request {
	r.openapiRequest.HostMap[key] = tea.String(value)
	return r
}

func (r *Request) SetHostParams(params map[string]string) *Request {
	for k, v := range params {
		r.openapiRequest.HostMap[k] = tea.String(v)
	}
	return r
}

func (r *Request) SetEndpointOverride(endpoint string) *Request {
	r.openapiRequest.EndpointOverride = tea.String(endpoint)
	return r
}

func (r *Request) GetOpenAPIRequest() *openapiutil.OpenApiRequest {
	return r.openapiRequest
}

func (r *Request) GetOpenAPIParams() *openapiClient.Params {
	return r.openapiParams
}
