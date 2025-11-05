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

	"github.com/aliyun/aliyun-cli/aliyun-cli-runtime/config"
	credentialsv2 "github.com/aliyun/credentials-go/credentials"
)

type OperationContext struct {
	credential credentialsv2.Credential

	profile *config.Profile

	client *Client

	nextLink string
}

type OperationContextOption func(*ClientConfig) error

// WithEndpoint sets a custom endpoint for the HTTP client
func WithEndpoint(endpoint string) OperationContextOption {
	return func(config *ClientConfig) error {
		config.Endpoint = endpoint
		return nil
	}
}

// WithServiceEndpoint generates an endpoint from service name
// Examples:
//   - "ecs" with useRegion=false → "ecs.aliyuncs.com"
//   - "cs" with useRegion=true → "cs.{region}.aliyuncs.com"
//   - "alert" with useRegion=true → "alert.{region}.aliyuncs.com"
func WithServiceEndpoint(serviceName string, useRegion bool) OperationContextOption {
	return func(config *ClientConfig) error {
		if serviceName == "" {
			return nil
		}
		serviceName = strings.ToLower(serviceName)
		if useRegion {
			config.Endpoint = fmt.Sprintf("%s.{region}.aliyuncs.com", serviceName)
		} else {
			// Format: {service}.aliyuncs.com
			config.Endpoint = fmt.Sprintf("%s.aliyuncs.com", serviceName)
		}
		return nil
	}
}

func WithEndpointTemplate(template string, params map[string]string) OperationContextOption {
	return WithEndpointTemplateFormat(template, params, "{", "}")
}

// WithEndpointTemplateFormat generates an endpoint from a template with custom placeholder format
// The template can contain placeholders with custom format like [region], [service], etc.
// Parameters are provided as a map and will be replaced in the template.
//
// Parameters:
//   - template: Endpoint template string with placeholders
//   - params: Map of parameter values to replace placeholders
//   - leftDelim: Left delimiter for placeholders (e.g., "{", "[", "{{"), defaults to "{" if empty
//   - rightDelim: Right delimiter for placeholders (e.g., "}", "]", "}}"), defaults to "}" if empty
//
// Examples:
//
//	// Using default {} format
//	ctx, err := http.NewOperationContext("", http.WithEndpointTemplateFormat(
//	    "{region}.alert.aliyuncs.com",
//	    nil,
//	    "{", "}",
//	))
//
//	// Using [] format
//	ctx, err := http.NewOperationContext("", http.WithEndpointTemplateFormat(
//	    "[region].alert.aliyuncs.com",
//	    nil,
//	    "[", "]",
//	))
//
//	// Using [] format with parameters
//	ctx, err := http.NewOperationContext("", http.WithEndpointTemplateFormat(
//	    "[region].[service].aliyuncs.com",
//	    map[string]string{
//	        "service": "cs",
//	    },
//	    "[", "]",
//	))
//
//	// Using {{}} format
//	ctx, err := http.NewOperationContext("", http.WithEndpointTemplateFormat(
//	    "{{subdomain}}.{{service}}.{{region}}.aliyuncs.com",
//	    map[string]string{
//	        "subdomain": "api",
//	        "service":   "ecs",
//	    },
//	    "{{", "}}",
//	))
//
// Supported placeholders:
//   - {region} or [region]: Automatically replaced with profile.RegionId (can also be provided in params)
//   - Any other placeholder: Must be provided in the params map
//
// If a placeholder is not provided in params and is not {region} or [region], it will remain as-is in the endpoint.
func WithEndpointTemplateFormat(template string, params map[string]string, leftDelim, rightDelim string) OperationContextOption {
	return func(config *ClientConfig) error {
		if template == "" {
			return nil
		}

		if leftDelim == "" {
			leftDelim = "{"
		}
		if rightDelim == "" {
			rightDelim = "}"
		}

		endpoint := template

		for key, value := range params {
			placeholder := fmt.Sprintf("%s%s%s", leftDelim, key, rightDelim)
			endpoint = strings.ReplaceAll(endpoint, placeholder, value)
		}

		config.Endpoint = endpoint
		return nil
	}
}

// NewOperationContext creates a new HTTP operation context
// Parameters:
//   - profileName: Profile name to load credentials from (empty string for default)
//   - options: Optional configuration functions (e.g., WithEndpoint, WithServiceEndpoint, WithEndpointTemplate)
func NewOperationContext(profileName string, options ...OperationContextOption) (*OperationContext, error) {
	cred, profile, err := config.LoadCredential(profileName)
	if err != nil {
		return nil, err
	}

	if profile.RegionId == "" {
		return nil, fmt.Errorf("region ID is required")
	}

	clientConfig := &ClientConfig{
		Credential:       cred,
		RegionId:         profile.RegionId,
		AutoRetry:        true,
		MaxRetryAttempts: 3,
		ReadTimeout:      30,
		ConnectTimeout:   10,
	}

	for _, option := range options {
		if err := option(clientConfig); err != nil {
			return nil, fmt.Errorf("failed to apply option: %w", err)
		}
	}
	if clientConfig.Endpoint != "" {
		clientConfig.Endpoint = strings.ReplaceAll(clientConfig.Endpoint, "{region}", profile.RegionId)
		clientConfig.Endpoint = strings.ReplaceAll(clientConfig.Endpoint, "[region]", profile.RegionId)
	}

	httpClient, err := NewClient(clientConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP client: %w", err)
	}

	return &OperationContext{
		credential: cred,
		profile:    profile,
		client:     httpClient,
	}, nil
}

func (ctx *OperationContext) GetCredential() credentialsv2.Credential {
	return ctx.credential
}

func (ctx *OperationContext) GetProfile() *config.Profile {
	return ctx.profile
}

func (ctx *OperationContext) GetHTTPClient() *Client {
	return ctx.client
}

func (ctx *OperationContext) GetRegionId() string {
	return ctx.profile.RegionId
}

func (ctx *OperationContext) SetNextLink(link string) {
	ctx.nextLink = link
}

func (ctx *OperationContext) GetNextLink() string {
	return ctx.nextLink
}
