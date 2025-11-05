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

	openapiClient "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	openapiTeaUtils "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	credentialsv2 "github.com/aliyun/credentials-go/credentials"
)

type Client struct {
	openapiClient *openapiClient.Client
	runtime       *openapiTeaUtils.RuntimeOptions
	config        *ClientConfig
}

type ClientConfig struct {
	Credential credentialsv2.Credential

	RegionId         string
	Endpoint         string
	UserAgent        string
	ReadTimeout      int64
	ConnectTimeout   int64
	IgnoreSSL        bool
	AutoRetry        bool
	MaxRetryAttempts int
	CustomHeaders    map[string]string
}

func NewClient(config *ClientConfig) (*Client, error) {
	if config == nil {
		return nil, fmt.Errorf("client config is required")
	}

	if config.RegionId == "" {
		return nil, fmt.Errorf("region ID is required")
	}

	conf := openapiClient.Config{
		Credential: config.Credential,
		RegionId:   tea.String(config.RegionId),
	}

	if config.Endpoint != "" {
		conf.Endpoint = tea.String(config.Endpoint)
	}

	if config.UserAgent != "" {
		conf.SetUserAgent(config.UserAgent)
	}

	if config.ReadTimeout > 0 {
		conf.SetReadTimeout(int(config.ReadTimeout * 1000))
	}

	if config.ConnectTimeout > 0 {
		conf.SetConnectTimeout(int(config.ConnectTimeout * 1000))
	}

	openapiClient, err := openapiClient.NewClient(&conf)
	if err != nil {
		return nil, fmt.Errorf("failed to create OpenAPI client: %w", err)
	}

	runtime := &openapiTeaUtils.RuntimeOptions{}
	if config.IgnoreSSL {
		runtime.SetIgnoreSSL(true)
	}

	if config.AutoRetry {
		runtime.SetAutoretry(true)
		if config.MaxRetryAttempts > 0 {
			runtime.SetMaxAttempts(config.MaxRetryAttempts)
		}
	}

	return &Client{
		openapiClient: openapiClient,
		runtime:       runtime,
		config:        config,
	}, nil
}

func (c *Client) GetOpenAPIClient() *openapiClient.Client {
	return c.openapiClient
}

func (c *Client) GetRuntime() *openapiTeaUtils.RuntimeOptions {
	return c.runtime
}

func (c *Client) GetConfig() *ClientConfig {
	return c.config
}
