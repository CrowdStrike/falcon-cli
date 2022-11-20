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

package root

import (
	"log"
	"strings"

	sensorCmd "github.com/crowdstrike/falcon-cli/pkg/cmd/sensor"
	versionCmd "github.com/crowdstrike/falcon-cli/pkg/cmd/version"
	"github.com/crowdstrike/falcon-cli/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"k8s.io/kubectl/pkg/util/templates"
)

var (
	shortDesc = `The CrowdStrike Falcon CLI`

	longDesc = templates.LongDesc(
		`The CrowdStrike Falcon CLI allows you to work effortlessly
        with the CrowdStrike Falcon platform.
    `)
)

type RootOptions struct {
	CfgFile      string
	CID          string
	ClientID     string
	ClientSecret string
	MemberCID    string
	Cloud        string
	Version      bool
	Help         bool
}

// NewCmdRoot represents the base command when called without any subcommands
func NewCmdRoot(f *utils.Factory, version string) *cobra.Command {
	opts := &RootOptions{}

	cmd := &cobra.Command{
		Use:   "falcon <command> <subcommand> [flags]",
		Short: shortDesc,
		Long:  longDesc,
		RunE:  runHelp,
	}

	cmd.PersistentFlags().StringVar(&opts.CfgFile, "config", "", "config file (default is $HOME/.falcon/falcon.yaml)")
	cmd.Flags().BoolVar(&opts.Version, "version", false, "Show version")
	cmd.PersistentFlags().BoolVar(&opts.Help, "help", false, "Show help for command")
	cmd.PersistentFlags().StringVarP(&opts.CID, "cid", "f", "", "The Falcon Customer ID (CID)")
	cmd.PersistentFlags().StringVarP(&opts.ClientID, "client-id", "u", "", "The Falcon API Oauth client ID")
	cmd.PersistentFlags().StringVarP(&opts.ClientSecret, "client-secret", "s", "", "The Falcon API Oauth client secret")
	cmd.PersistentFlags().StringVarP(&opts.MemberCID, "member-cid", "m", "", "The Falcon API member CID")
	cmd.PersistentFlags().StringVarP(&opts.Cloud, "cloud", "r", "autodiscover", "The Falcon API Cloud Region")

	pf := cmd.PersistentFlags()
	normalizeFunc := pf.GetNormalizeFunc()
	pf.SetNormalizeFunc(func(f *pflag.FlagSet, name string) pflag.NormalizedName {
		result := normalizeFunc(f, name)
		name = strings.ReplaceAll(string(result), "-", "_")
		return pflag.NormalizedName(name)
	})

	// Bind flags to viper
	if err := viper.GetViper().BindPFlags(cmd.PersistentFlags()); err != nil {
		log.Fatalf("Error binding flags to viper: %v", err)
	}

	// Add subcommands
	cmd.AddCommand(versionCmd.NewCmdVersion(f))
	cmd.AddCommand(sensorCmd.NewSensorCmd(f))

	return cmd
}

func runHelp(cmd *cobra.Command, args []string) error {
	return cmd.Help()
}
