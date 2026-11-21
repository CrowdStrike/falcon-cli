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
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

func TestNewConfig_NoProfileReadsValues(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)
	viper.Set("client_id", "id-123")

	cfg, err := NewConfig()
	if err != nil {
		t.Fatalf("NewConfig: %v", err)
	}
	if cfg.ClientID != "id-123" {
		t.Errorf("ClientID = %q, want %q", cfg.ClientID, "id-123")
	}
	if cfg.ClientSecret != "" {
		t.Errorf("ClientSecret = %q, want empty", cfg.ClientSecret)
	}
	if cfg.Cloud != "autodiscover" {
		t.Errorf("Cloud = %q, want autodiscover", cfg.Cloud)
	}
}

func TestSave_PreservesOtherProfiles(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".falcon", "config.yaml")

	if err := (Config{Profile: "default", ClientID: "a", ClientSecret: "b"}).Save(path); err != nil {
		t.Fatalf("Save default: %v", err)
	}
	if err := (Config{Profile: "prod", ClientID: "c", ClientSecret: "d", Cloud: "eu-1"}).Save(path); err != nil {
		t.Fatalf("Save prod: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat: %v", err)
	}
	if perm := info.Mode().Perm(); perm != 0600 {
		t.Errorf("file mode = %o, want 600", perm)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	var got map[string]profileEntry
	if err := yaml.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if got["default"].ClientID != "a" || got["prod"].ClientID != "c" || got["prod"].Cloud != "eu-1" {
		t.Errorf("unexpected profiles: %+v", got)
	}
}
