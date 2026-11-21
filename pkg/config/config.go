// Copyright (c) 2026 CrowdStrike, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package config

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"

	"github.com/crowdstrike/falcon-cli/internal/version"
	"github.com/crowdstrike/gofalcon/falcon"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

// Config holds persistent configuration for falcon.
type Config struct {
	// The Falcon Customer ID
	CID string `yaml:"cid,omitempty"`
	// The Falcon API client ID.
	ClientID string `yaml:"client_id"`
	// The Falcon API client secret.
	ClientSecret string `yaml:"client_secret"`
	// The Falcon API base URL.
	MemberCID string `yaml:"member_cid,omitempty"`
	// The Falcon API cloud region.
	Cloud string `yaml:"cloud,omitempty"`
	// BaseURL overrides the Falcon API host (e.g. https://api.eu-1.crowdstrike.com).
	BaseURL string `yaml:"base_url,omitempty"`
	// The OAuth token returned from the Falcon API.
	OauthToken string `yaml:"oauth_token,omitempty"`
	// The Container Registry OAuth token returned from the Falcon API.
	RegistryToken string `yaml:"registry_token,omitempty"`
	// The Profile to use for the CLI
	Profile string `yaml:"profile,omitempty"`
}

var ConfigFile string

func NewConfig() (Config, error) {
	c := &Config{}

	profile := viper.GetViper().GetString("profile")

	c.CID = getViperKey("cid", profile)
	c.ClientID = getViperKey("client_id", profile)
	c.ClientSecret = getViperKey("client_secret", profile)
	c.MemberCID = getViperKey("member_cid", profile)
	c.Cloud = getViperKey("cloud", profile)
	c.BaseURL = getViperKey("base_url", profile)
	c.Profile = profile

	if c.Cloud == "" {
		c.Cloud = "autodiscover"
	}

	return *c, nil
}

func (c Config) APIConfig(_ string) *falcon.ApiConfig {
	cfg := &falcon.ApiConfig{
		ClientId:          c.ClientID,
		ClientSecret:      c.ClientSecret,
		MemberCID:         c.MemberCID,
		Cloud:             falcon.Cloud(c.Cloud),
		Context:           context.Background(),
		UserAgentOverride: fmt.Sprintf("falcon-cli/%s", version.Version),
		Debug:             viper.GetViper().GetBool("verbose"),
	}
	if c.BaseURL != "" {
		if u, err := url.Parse(c.BaseURL); err == nil && u.Host != "" {
			cfg.HostOverride = u.Host
		} else {
			cfg.HostOverride = c.BaseURL
		}
	}
	return cfg
}

func getViperKey(key string, profile string) string {
	if profile == "" {
		return viper.GetViper().GetString(key)
	}

	profileKey := fmt.Sprintf("%s.%s", profile, key)
	return viper.GetViper().GetString(profileKey)
}

// profileEntry is the on-disk representation of a single profile.
type profileEntry struct {
	CID          string `yaml:"cid,omitempty"`
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
	MemberCID    string `yaml:"member_cid,omitempty"`
	Cloud        string `yaml:"cloud,omitempty"`
	BaseURL      string `yaml:"base_url,omitempty"`
}

// DefaultConfigFile returns the default config file path ($HOME/.falcon/config.yaml).
func DefaultConfigFile() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".falcon", "config.yaml"), nil
}

// Save writes the config as the named profile into the YAML file at path,
// preserving any other profiles already stored there.
func (c Config) Save(path string) error {
	if c.Profile == "" {
		return errors.New("profile name must not be empty")
	}

	profiles := map[string]any{}
	data, err := os.ReadFile(filepath.Clean(path))
	switch {
	case err == nil:
		if err := yaml.Unmarshal(data, &profiles); err != nil {
			return fmt.Errorf("parsing %s: %w", path, err)
		}
		if profiles == nil {
			profiles = map[string]any{}
		}
	case !errors.Is(err, fs.ErrNotExist):
		return fmt.Errorf("reading %s: %w", path, err)
	}

	profiles[c.Profile] = profileEntry{
		CID:          c.CID,
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
		MemberCID:    c.MemberCID,
		Cloud:        c.Cloud,
		BaseURL:      c.BaseURL,
	}

	out, err := yaml.Marshal(profiles)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	return os.WriteFile(filepath.Clean(path), out, 0600)
}
