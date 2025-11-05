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
	"encoding/json"
	"fmt"
	"os"
	"runtime"
)

const (
	configPath = "/.aliyun"
	configFile = "config.json"
)

// Configuration represents the aliyun-cli configuration file structure
type Configuration struct {
	CurrentProfile string    `json:"current"`
	Profiles       []Profile `json:"profiles"`
	MetaPath       string    `json:"meta_path"`
}

// Profile represents a single profile in the configuration
type Profile struct {
	Name                      string `json:"name"`
	Mode                      string `json:"mode"`
	AccessKeyId               string `json:"access_key_id,omitempty"`
	AccessKeySecret           string `json:"access_key_secret,omitempty"`
	StsToken                 string `json:"sts_token,omitempty"`
	StsRegion                string `json:"sts_region,omitempty"`
	RamRoleName              string `json:"ram_role_name,omitempty"`
	RamRoleArn               string `json:"ram_role_arn,omitempty"`
	RoleSessionName          string `json:"ram_session_name,omitempty"`
	ExternalId               string `json:"external_id,omitempty"`
	SourceProfile            string `json:"source_profile,omitempty"`
	PrivateKey               string `json:"private_key,omitempty"`
	KeyPairName              string `json:"key_pair_name,omitempty"`
	ExpiredSeconds           int    `json:"expired_seconds,omitempty"`
	RegionId                 string `json:"region_id,omitempty"`
	OutputFormat             string `json:"output_format,omitempty"`
	Language                 string `json:"language,omitempty"`
	Site                     string `json:"site,omitempty"`
	ReadTimeout              int    `json:"retry_timeout,omitempty"`
	ConnectTimeout           int    `json:"connect_timeout,omitempty"`
	RetryCount               int    `json:"retry_count,omitempty"`
	ProcessCommand           string `json:"process_command,omitempty"`
	CredentialsURI           string `json:"credentials_uri,omitempty"`
	OIDCProviderARN          string `json:"oidc_provider_arn,omitempty"`
	OIDCTokenFile            string `json:"oidc_token_file,omitempty"`
	CloudSSOSignInUrl        string `json:"cloud_sso_sign_in_url,omitempty"`
	CloudSSOAccessConfig     string `json:"cloud_sso_access_config,omitempty"`
	OAuthAccessToken         string `json:"oauth_access_token,omitempty"`
	OAuthRefreshToken        string `json:"oauth_refresh_token,omitempty"`
	OAuthSiteType            string `json:"oauth_site_type,omitempty"`
}

// GetConfigPath returns the path to the aliyun-cli configuration directory
func GetConfigPath() string {
	path := GetHomePath() + configPath
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, 0755)
		if err != nil {
			panic(err)
		}
	}
	return path
}

// GetConfigFilePath returns the full path to the config.json file
func GetConfigFilePath() string {
	return GetConfigPath() + "/" + configFile
}

// GetHomePath returns the user's home directory path
func GetHomePath() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}

// LoadConfiguration loads the configuration from the default path
func LoadConfiguration() (*Configuration, error) {
	return LoadConfigurationFromFile(GetConfigFilePath())
}

// LoadConfigurationFromFile loads the configuration from a specific file path
func LoadConfigurationFromFile(filePath string) (*Configuration, error) {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("reading config from '%s' failed: %w", filePath, err)
	}

	conf := &Configuration{}
	err = json.Unmarshal(bytes, conf)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return conf, nil
}

// GetProfile retrieves a profile by name from the configuration
func (c *Configuration) GetProfile(name string) (*Profile, error) {
	// If name is empty, use current profile
	if name == "" {
		name = c.CurrentProfile
	}
	// Default to "default" if current profile is empty
	if name == "" {
		name = "default"
	}

	// Check environment variable for profile name
	if name == "default" {
		if envProfile := os.Getenv("ALIBABACLOUD_PROFILE"); envProfile != "" {
			name = envProfile
		} else if envProfile := os.Getenv("ALIBABA_CLOUD_PROFILE"); envProfile != "" {
			name = envProfile
		} else if envProfile := os.Getenv("ALICLOUD_PROFILE"); envProfile != "" {
			name = envProfile
		}
	}

	for _, p := range c.Profiles {
		if p.Name == name {
			return &p, nil
		}
	}

	return nil, fmt.Errorf("profile '%s' not found in configuration", name)
}

// LoadProfile loads a specific profile by name from the default configuration file
func LoadProfile(profileName string) (*Profile, error) {
	config, err := LoadConfiguration()
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	profile, err := config.GetProfile(profileName)
	if err != nil {
		return nil, err
	}

	return profile, nil
}

