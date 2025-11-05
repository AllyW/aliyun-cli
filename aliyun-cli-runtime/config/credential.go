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

package config

import (
	"fmt"
	"os"
	"strings"

	credentialsv2 "github.com/aliyun/credentials-go/credentials"
)

// LoadCredential loads a credential from the aliyun-cli configuration file
// It supports the same profile name resolution as the main aliyun-cli:
// 1. Use the provided profileName if not empty
// 2. Use ALIBABACLOUD_PROFILE, ALIBABA_CLOUD_PROFILE, or ALICLOUD_PROFILE environment variable
// 3. Use the "current" profile from config.json
// 4. Default to "default" profile
func LoadCredential(profileName string) (credentialsv2.Credential, *Profile, error) {
	// Load configuration
	cfg, err := LoadConfiguration()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	// Get profile
	profile, err := cfg.GetProfile(profileName)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get profile: %w", err)
	}

	// Load environment variables as fallback
	profile.LoadFromEnv()

	// Create credential from profile
	cred, err := profile.GetCredential()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create credential: %w", err)
	}

	return cred, profile, nil
}

// LoadCredentialFromProfile loads a credential from a Profile object
func LoadCredentialFromProfile(profile *Profile) (credentialsv2.Credential, error) {
	profile.LoadFromEnv()
	return profile.GetCredential()
}

// GetCredential creates a credentials.Credential from the Profile
func (p *Profile) GetCredential() (credentialsv2.Credential, error) {
	credConfig := &credentialsv2.Config{}

	// Auto-detect mode if not set
	if p.Mode == "" {
		p.AutoDetectMode()
	}

	switch strings.ToUpper(p.Mode) {
	case "AK", "ACCESS_KEY":
		if p.AccessKeyId == "" || p.AccessKeySecret == "" {
			return nil, fmt.Errorf("AccessKeyId/AccessKeySecret is empty, run 'aliyun configure' first")
		}
		credConfig.SetType("access_key").
			SetAccessKeyId(p.AccessKeyId).
			SetAccessKeySecret(p.AccessKeySecret)

	case "STS", "STSTOKEN":
		if p.AccessKeyId == "" || p.AccessKeySecret == "" || p.StsToken == "" {
			return nil, fmt.Errorf("AccessKeyId/AccessKeySecret/StsToken is empty, run 'aliyun configure' first")
		}
		credConfig.SetType("sts").
			SetAccessKeyId(p.AccessKeyId).
			SetAccessKeySecret(p.AccessKeySecret).
			SetSecurityToken(p.StsToken)

	case "RAMROLEARN":
		if p.RamRoleArn == "" {
			return nil, fmt.Errorf("RamRoleArn is empty, run 'aliyun configure' first")
		}
		credConfig.SetType("ram_role_arn").
			SetAccessKeyId(p.AccessKeyId).
			SetAccessKeySecret(p.AccessKeySecret).
			SetRoleArn(p.RamRoleArn).
			SetRoleSessionName(p.RoleSessionName).
			SetRoleSessionExpiration(p.ExpiredSeconds).
			SetExternalId(p.ExternalId).
			SetSTSEndpoint(getSTSEndpoint(p.StsRegion))

		if p.StsToken != "" {
			credConfig.SetSecurityToken(p.StsToken)
		}

	case "ECSRAMROLE":
		if p.RamRoleName == "" {
			return nil, fmt.Errorf("RamRoleName is empty, run 'aliyun configure' first")
		}
		credConfig.SetType("ecs_ram_role").
			SetRoleName(p.RamRoleName)

	case "RSAKEYPAIR":
		if p.PrivateKey == "" || p.KeyPairName == "" {
			return nil, fmt.Errorf("PrivateKey/KeyPairName is empty, run 'aliyun configure' first")
		}
		credConfig.SetType("rsa_key_pair").
			SetPrivateKeyFile(p.PrivateKey).
			SetPublicKeyId(p.KeyPairName).
			SetSessionExpiration(p.ExpiredSeconds).
			SetSTSEndpoint(getSTSEndpoint(p.StsRegion))

	case "CREDENTIALSURI":
		// CredentialsURI mode is not directly supported by credentials-go library
		// This would require custom HTTP implementation to fetch credentials from URI
		// For now, we'll return an error suggesting to use other authentication modes
		return nil, fmt.Errorf("CredentialsURI mode is not yet supported in aliyun-cli-runtime, please use AK, STS, or RamRoleArn mode")

	case "OIDC":
		if p.OIDCProviderARN == "" || p.OIDCTokenFile == "" {
			return nil, fmt.Errorf("OIDCProviderARN/OIDCTokenFile is empty, run 'aliyun configure' first")
		}
		credConfig.SetType("oidc_role_arn").
			SetOIDCProviderArn(p.OIDCProviderARN).
			SetOIDCTokenFilePath(p.OIDCTokenFile).
			SetRoleArn(p.RamRoleArn).
			SetRoleSessionName(p.RoleSessionName).
			SetRoleSessionExpiration(p.ExpiredSeconds)

	default:
		return nil, fmt.Errorf("unsupported authentication mode: %s", p.Mode)
	}

	return credentialsv2.NewCredential(credConfig)
}

// LoadFromEnv loads credential values from environment variables if they are not set in profile
func (p *Profile) LoadFromEnv() {
	if p.AccessKeyId == "" {
		p.AccessKeyId = getFromEnv("ALIBABA_CLOUD_ACCESS_KEY_ID", "ALIBABACLOUD_ACCESS_KEY_ID", "ALICLOUD_ACCESS_KEY_ID", "ACCESS_KEY_ID")
	}

	if p.AccessKeySecret == "" {
		p.AccessKeySecret = getFromEnv("ALIBABA_CLOUD_ACCESS_KEY_SECRET", "ALIBABACLOUD_ACCESS_KEY_SECRET", "ALICLOUD_ACCESS_KEY_SECRET", "ACCESS_KEY_SECRET")
	}

	if p.StsToken == "" {
		p.StsToken = getFromEnv("ALIBABA_CLOUD_SECURITY_TOKEN", "ALIBABACLOUD_SECURITY_TOKEN", "ALICLOUD_SECURITY_TOKEN", "SECURITY_TOKEN")
	}

	if p.RegionId == "" {
		p.RegionId = getFromEnv("ALIBABA_CLOUD_REGION_ID", "ALIBABACLOUD_REGION_ID", "ALICLOUD_REGION_ID", "REGION_ID", "REGION")
	}

	if p.RamRoleArn == "" {
		p.RamRoleArn = getFromEnv("ALIBABACLOUD_ROLE_ARN", "ALIBABA_CLOUD_ROLE_ARN")
	}

	if p.ExternalId == "" {
		p.ExternalId = getFromEnv("ALIBABACLOUD_EXTERNAL_ID", "ALIBABA_CLOUD_EXTERNAL_ID")
	}

	if p.CredentialsURI == "" {
		p.CredentialsURI = os.Getenv("ALIBABA_CLOUD_CREDENTIALS_URI")
	}

	if p.OIDCProviderARN == "" {
		p.OIDCProviderARN = getFromEnv("ALIBABACLOUD_OIDC_PROVIDER_ARN", "ALIBABA_CLOUD_OIDC_PROVIDER_ARN")
	}

	if p.OIDCTokenFile == "" {
		p.OIDCTokenFile = getFromEnv("ALIBABACLOUD_OIDC_TOKEN_FILE", "ALIBABA_CLOUD_OIDC_TOKEN_FILE")
	}
}

// AutoDetectMode automatically detects the authentication mode based on available fields
func (p *Profile) AutoDetectMode() {
	if p.CredentialsURI != "" {
		p.Mode = "CredentialsURI"
		return
	}

	if p.OIDCProviderARN != "" && p.OIDCTokenFile != "" {
		p.Mode = "OIDC"
		return
	}

	if p.RamRoleName != "" {
		p.Mode = "EcsRamRole"
		return
	}

	if p.RamRoleArn != "" {
		p.Mode = "RamRoleArn"
		return
	}

	if p.StsToken != "" {
		p.Mode = "StsToken"
		return
	}

	if p.AccessKeyId != "" && p.AccessKeySecret != "" {
		p.Mode = "AK"
		return
	}

	// Default to AK if we have at least access key from env
	if p.AccessKeyId != "" || p.AccessKeySecret != "" {
		p.Mode = "AK"
		return
	}

	// Default to AK mode
	p.Mode = "AK"
}

// getFromEnv tries to get a value from environment variables, trying each key in order
func getFromEnv(keys ...string) string {
	for _, key := range keys {
		if val := os.Getenv(key); val != "" {
			return val
		}
	}
	return ""
}

// getSTSEndpoint returns the STS endpoint based on region
func getSTSEndpoint(region string) string {
	if region == "" {
		return "sts.aliyuncs.com"
	}

	// Handle special regions
	if region == "cn-hangzhou" || region == "cn-shanghai" || region == "cn-beijing" || region == "cn-shenzhen" {
		return "sts.cn-hangzhou.aliyuncs.com"
	}

	// For other regions, use the standard format
	return fmt.Sprintf("sts.%s.aliyuncs.com", region)
}

