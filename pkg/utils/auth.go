// Copyright (c) 2022 CrowdStrike, Inc.
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

package utils

import (
	"github.com/crowdstrike/falcon-cli/pkg/config"
	"github.com/spf13/cobra"
)

// CheckAuth checks if required authentication is configured
func CheckAuth(cfg config.Config) bool {
	// Verify required variables are set
	if cfg.ClientID == "" || cfg.ClientSecret == "" {
		return false
	}
	return true
}

// DisableAuthCheck disables the auth check for a command
func DisableAuthCheck(cmd *cobra.Command) {
	if cmd.Annotations == nil {
		cmd.Annotations = map[string]string{}
	}

	cmd.Annotations["skipAuthCheck"] = "true"
}

// IsAuthCheckEnabled checks if the auth check is enabled for a command
func IsAuthCheckEnabled(cmd *cobra.Command) bool {
	switch cmd.Name() {
	case "help", cobra.ShellCompRequestCmd, cobra.ShellCompNoDescRequestCmd, "falcon":
		return false
	}

	for c := cmd; c.Parent() != nil; c = c.Parent() {
		if c.Annotations != nil && c.Annotations["skipAuthCheck"] == "true" {
			return false
		}
	}

	return true
}
