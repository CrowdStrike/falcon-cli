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
	"fmt"

	"github.com/crowdstrike/falcon-cli/pkg/cmd/auth"
	"github.com/crowdstrike/falcon-cli/pkg/cmd/sensor"
	versionCmd "github.com/crowdstrike/falcon-cli/pkg/cmd/version"
	"github.com/crowdstrike/falcon-cli/pkg/config"
	"github.com/crowdstrike/falcon-cli/pkg/factory"
	"github.com/crowdstrike/falcon-cli/pkg/iostreams"
	"github.com/crowdstrike/falcon-cli/pkg/utils"
	"github.com/crowdstrike/falcon-cli/pkg/version"
	"github.com/spf13/cobra"
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
	IO      *iostreams.IOStreams
	Config  config.Config
	Version bool
}

// NewCmdRoot represents the base command when called without any subcommands
func NewCmdRoot(f *factory.Factory, version string) *cobra.Command {
	opts := &RootOptions{}

	cmd := &cobra.Command{
		Use:   "falcon <command> <subcommand> [flags]",
		Short: shortDesc,
		Long:  longDesc,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := &RootOptions{
				IO: f.IOStreams,
			}

			cfg, err := f.Config()
			if err != nil {
				return err
			}
			opts.Config = cfg

			return runRoot(cmd, opts)
		},
	}

	cmd.PersistentFlags().String("config", "", "config file (default is $HOME/.falcon/falcon.yaml)")
	cmd.PersistentFlags().BoolVar(&opts.Version, "version", false, "Show version")
	cmd.PersistentFlags().Bool("help", false, "Show help for command")
	cmd.PersistentFlags().StringP("cid", "f", "", "The Falcon Customer ID (CID)")
	cmd.PersistentFlags().StringP("client-id", "u", "", "The Falcon API Oauth client ID")
	cmd.PersistentFlags().StringP("client-secret", "s", "", "The Falcon API Oauth client secret")
	cmd.PersistentFlags().StringP("member-cid", "m", "", "The Falcon API member CID")
	cmd.PersistentFlags().StringP("cloud", "r", "autodiscover", "The Falcon API Cloud Region")
	cmd.PersistentFlags().StringP("profile", "p", "default", "Use a specific profile from your config file")

	// Add subcommands
	cmd.AddCommand(versionCmd.NewCmdVersion(f))
	cmd.AddCommand(sensor.NewSensorCmd(f))
	cmd.AddCommand(auth.NewAuthCmd(f))

	utils.DisableAuthCheck(cmd)

	return cmd
}

func runRoot(cmd *cobra.Command, opts *RootOptions) error {
	if cmd.Flags().Changed("version") {
		_, err := fmt.Fprint(opts.IO.Out, version.VersionString())
		return err
	}

	return cmd.Help()
}
